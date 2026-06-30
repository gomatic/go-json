package value

import (
	"encoding/json"
	"testing"
)

// maxFuzzDepth bounds the recursive walk so a pathologically nested fuzz input
// cannot overflow the harness stack (a property of the input, not of the code
// under test).
const maxFuzzDepth = 64

// decodeSeeds are valid JSON documents spanning every Value kind plus the
// awkward edges: empty composites, the empty string, negative and very large
// numbers, an integer past float64's exact range, unicode/escapes, and deep
// nesting. A decoded JSON document is already a Value, so these exercise the
// decode-to-accessor seam directly.
var decodeSeeds = []string{
	`null`,
	`true`,
	`false`,
	`0`,
	`-1`,
	`3.14`,
	`-0`,
	`1e308`,
	`""`,
	`"hello"`,
	"\"é \U0001f600 \\u0041\"",
	`[]`,
	`[1,2,3]`,
	`{}`,
	`{"name":"ada","age":36}`,
	`{"a":[1,{"b":null}],"c":true,"d":""}`,
	`[[[[[1]]]]]`,
	`123456789012345678901234567890`,
}

// FuzzDecodeAccessors asserts that for any value decodable from JSON the
// accessor and comparison contract holds: the accessor matching KindOf accepts
// its own value, Equal is reflexive, Compare(v,v) is 0 when ordered, and the
// value survives a JSON round-trip unchanged — and that none of this panics.
func FuzzDecodeAccessors(f *testing.F) {
	for _, s := range decodeSeeds {
		f.Add(s)
	}
	f.Fuzz(func(t *testing.T, in string) {
		var v Value
		if json.Unmarshal([]byte(in), &v) != nil {
			return
		}
		assertRoundTrip(t, v)
		walk(t, v, 0)
	})
}

// pairSeeds are JSON scalars used to fuzz the binary operators across kinds.
var pairSeeds = []string{`null`, `true`, `0`, `-2.5`, `"a"`, `"b"`, `[1]`, `{"k":1}`}

// FuzzCompareAdd asserts the algebraic invariants of the binary operators over
// any pair of JSON-decodable values: Equal is symmetric, Compare is
// antisymmetric where both directions are ordered, and Add never panics.
func FuzzCompareAdd(f *testing.F) {
	for _, a := range pairSeeds {
		f.Add(a, a)
		f.Add(`1`, a)
	}
	f.Fuzz(func(t *testing.T, ja, jb string) {
		var a, b Value
		if json.Unmarshal([]byte(ja), &a) != nil || json.Unmarshal([]byte(jb), &b) != nil {
			return
		}
		assertEqualSymmetric(t, a, b)
		assertCompareAntisymmetric(t, a, b)
		_, _ = Add(a, b) // must not panic for any pairing
	})
}

// walk asserts the per-value invariants for v and descends into composites up
// to maxFuzzDepth.
func walk(t *testing.T, v Value, depth int) {
	t.Helper()
	assertKindConsistency(t, v)
	assertReflexive(t, v)
	if depth >= maxFuzzDepth {
		return
	}
	descend(t, v, depth+1)
}

// descend walks the elements of a list or the values of an object.
func descend(t *testing.T, v Value, depth int) {
	t.Helper()
	switch c := v.(type) {
	case []Value:
		for _, e := range c {
			walk(t, e, depth)
		}
	case map[string]Value:
		for _, e := range c {
			walk(t, e, depth)
		}
	}
}

// kindAccessors maps each non-null Kind to the accessor that must accept a value
// of that Kind without error.
var kindAccessors = map[Kind]func(Value) error{
	KindBool:   func(v Value) error { _, e := AsBool(v); return e },
	KindInt:    func(v Value) error { _, e := AsInt(v); return e },
	KindFloat:  func(v Value) error { _, e := AsFloat(v); return e },
	KindString: func(v Value) error { _, e := AsString(v); return e },
	KindList:   func(v Value) error { _, e := AsList(v); return e },
	KindObject: func(v Value) error { _, e := AsObject(v); return e },
}

// assertKindConsistency verifies the accessor matching KindOf accepts v, and
// that Truthy does not panic.
func assertKindConsistency(t *testing.T, v Value) {
	t.Helper()
	_ = Truthy(v)
	if acc, ok := kindAccessors[KindOf(v)]; ok && acc(v) != nil {
		t.Fatalf("accessor for kind %v rejected its own value %#v", KindOf(v), v)
	}
}

// assertReflexive verifies Equal(v,v) and that an ordered Compare(v,v) is 0.
func assertReflexive(t *testing.T, v Value) {
	t.Helper()
	if !Equal(v, v) {
		t.Fatalf("Equal not reflexive for %#v", v)
	}
	if c, err := Compare(v, v); err == nil && c != 0 {
		t.Fatalf("Compare(v,v) = %d, want 0 for %#v", c, v)
	}
}

// assertRoundTrip verifies v re-marshals and re-decodes to an Equal value.
func assertRoundTrip(t *testing.T, v Value) {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("re-marshal failed for %#v: %v", v, err)
	}
	var v2 Value
	if err := json.Unmarshal(b, &v2); err != nil {
		t.Fatalf("re-unmarshal failed for %q: %v", b, err)
	}
	if !Equal(v, v2) {
		t.Fatalf("round-trip changed value: %#v -> %#v", v, v2)
	}
}

// assertEqualSymmetric verifies Equal(a,b) == Equal(b,a).
func assertEqualSymmetric(t *testing.T, a, b Value) {
	t.Helper()
	if Equal(a, b) != Equal(b, a) {
		t.Fatalf("Equal asymmetric for %#v, %#v", a, b)
	}
}

// assertCompareAntisymmetric verifies Compare(a,b) == -Compare(b,a) whenever
// both directions are ordered.
func assertCompareAntisymmetric(t *testing.T, a, b Value) {
	t.Helper()
	c, cerr := Compare(a, b)
	d, derr := Compare(b, a)
	if cerr == nil && derr == nil && c != -d {
		t.Fatalf("Compare not antisymmetric: Compare(%#v,%#v)=%d, reverse=%d", a, b, c, d)
	}
}

package value

import (
	"errors"
	"testing"
)

func TestError_Message(t *testing.T) {
	if ErrNotObject.Error() != "value: not an object" {
		t.Fatalf("unexpected message %q", ErrNotObject.Error())
	}
}

func TestKindOf(t *testing.T) {
	cases := []struct {
		v    Value
		name string
		want Kind
	}{
		{name: "null", v: nil, want: KindNull},
		{name: "bool", v: true, want: KindBool},
		{name: "int", v: int64(1), want: KindInt},
		{name: "float", v: 3.14, want: KindFloat},
		{name: "string", v: "s", want: KindString},
		{name: "list", v: []Value{int64(1)}, want: KindList},
		{name: "object", v: map[string]Value{"a": int64(1)}, want: KindObject},
		{name: "unknown", v: struct{}{}, want: KindNull},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := KindOf(c.v); got != c.want {
				t.Errorf("KindOf(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}

func TestAsObject(t *testing.T) {
	m, err := AsObject(map[string]Value{"a": int64(1)})
	if err != nil || m["a"] != int64(1) {
		t.Fatalf("AsObject ok: %v %v", m, err)
	}
	if _, err := AsObject("x"); !errors.Is(err, ErrNotObject) {
		t.Fatalf("AsObject err: got %v want ErrNotObject", err)
	}
}

func TestAsList(t *testing.T) {
	l, err := AsList([]Value{int64(1)})
	if err != nil || len(l) != 1 {
		t.Fatalf("AsList ok: %v %v", l, err)
	}
	if _, err := AsList("x"); !errors.Is(err, ErrNotList) {
		t.Fatalf("AsList err: got %v want ErrNotList", err)
	}
}

func TestAsString(t *testing.T) {
	s, err := AsString("hi")
	if err != nil || s != "hi" {
		t.Fatalf("AsString ok: %v %v", s, err)
	}
	if _, err := AsString(int64(1)); !errors.Is(err, ErrNotString) {
		t.Fatalf("AsString err: got %v want ErrNotString", err)
	}
}

func TestAsBool(t *testing.T) {
	b, err := AsBool(true)
	if err != nil || !b {
		t.Fatalf("AsBool ok: %v %v", b, err)
	}
	if _, err := AsBool("x"); !errors.Is(err, ErrNotBool) {
		t.Fatalf("AsBool err: got %v want ErrNotBool", err)
	}
}

func TestAsInt(t *testing.T) {
	if n, err := AsInt(int64(5)); err != nil || n != 5 {
		t.Fatalf("AsInt(int): %v %v", n, err)
	}
	if n, err := AsInt(5.9); err != nil || n != 5 {
		t.Fatalf("AsInt(float) trunc: %v %v", n, err)
	}
	if _, err := AsInt("x"); !errors.Is(err, ErrNotNumber) {
		t.Fatalf("AsInt err: got %v want ErrNotNumber", err)
	}
}

func TestAsFloat(t *testing.T) {
	if n, err := AsFloat(2.5); err != nil || n != 2.5 {
		t.Fatalf("AsFloat(float): %v %v", n, err)
	}
	if n, err := AsFloat(int64(3)); err != nil || n != 3.0 {
		t.Fatalf("AsFloat(int): %v %v", n, err)
	}
	if _, err := AsFloat("x"); !errors.Is(err, ErrNotNumber) {
		t.Fatalf("AsFloat err: got %v want ErrNotNumber", err)
	}
}

func TestTruthy(t *testing.T) {
	cases := map[string]struct {
		v    Value
		want bool
	}{
		"nil":      {nil, false},
		"false":    {false, false},
		"true":     {true, true},
		"zeroInt":  {int64(0), true},
		"emptyStr": {"", true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			if got := Truthy(c.v); got != c.want {
				t.Errorf("Truthy(%#v) = %v, want %v", c.v, got, c.want)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	if !Equal(int64(2), 2.0) {
		t.Error("2 should equal 2.0 across numeric types")
	}
	if !Equal("a", "a") {
		t.Error("equal strings")
	}
	if Equal("a", "b") {
		t.Error("unequal strings")
	}
	if Equal(true, int64(1)) {
		t.Error("bool and int are not equal")
	}
}

func TestCompare(t *testing.T) {
	if c, err := Compare(int64(1), 2.0); err != nil || c != -1 {
		t.Fatalf("Compare(1,2.0)=%d,%v want -1", c, err)
	}
	if c, err := Compare(2.0, 2.0); err != nil || c != 0 {
		t.Fatalf("Compare(2,2)=%d,%v want 0", c, err)
	}
	if c, err := Compare(3.0, 2.0); err != nil || c != 1 {
		t.Fatalf("Compare(3,2)=%d,%v want 1", c, err)
	}
	if c, err := Compare("a", "b"); err != nil || c != -1 {
		t.Fatalf("Compare(a,b)=%d,%v want -1", c, err)
	}
	if c, err := Compare("b", "a"); err != nil || c != 1 {
		t.Fatalf("Compare(b,a)=%d,%v want 1", c, err)
	}
	if c, err := Compare("a", "a"); err != nil || c != 0 {
		t.Fatalf("Compare(a,a)=%d,%v want 0", c, err)
	}
	if _, err := Compare("a", int64(1)); !errors.Is(err, ErrIncomparable) {
		t.Fatalf("Compare(str,int)=%v want ErrIncomparable", err)
	}
}

func TestAdd(t *testing.T) {
	if v, err := Add(int64(1), int64(2)); err != nil || v != int64(3) {
		t.Fatalf("Add(int,int)=%v,%v want int64(3)", v, err)
	}
	if v, err := Add(int64(1), 2.5); err != nil || v != 3.5 {
		t.Fatalf("Add(int,float)=%v,%v want 3.5", v, err)
	}
	if v, err := Add("a", "b"); err != nil || v != "ab" {
		t.Fatalf("Add(str,str)=%v,%v want ab", v, err)
	}
	if v, err := Add("n=", int64(5)); err != nil || v != "n=5" {
		t.Fatalf("Add(str,int)=%v,%v want n=5", v, err)
	}
	if v, err := Add(2.5, "x"); err != nil || v != "2.5x" {
		t.Fatalf("Add(float,str)=%v,%v want 2.5x", v, err)
	}
	if v, err := Add(true, "!"); err != nil || v != "true!" {
		t.Fatalf("Add(bool,str)=%v,%v want true!", v, err)
	}
	if _, err := Add([]Value{}, int64(1)); !errors.Is(err, ErrNotNumber) {
		t.Fatalf("Add(list,int)=%v want ErrNotNumber", err)
	}
}

func TestAdd_ConcatRejectsComposite(t *testing.T) {
	if _, err := Add("a", []Value{}); !errors.Is(err, ErrNotString) {
		t.Fatalf("Add(str,list)=%v want ErrNotString", err)
	}
	if _, err := Add([]Value{}, "a"); !errors.Is(err, ErrNotString) {
		t.Fatalf("Add(list,str)=%v want ErrNotString", err)
	}
}

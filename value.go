// Package value is a dynamic, JSON-compatible value model: a value space with
// typed, constant-error accessors and the coercion rules a small expression
// language needs. A decoded JSON document is already a [Value], so the package
// interoperates directly with encoding/json while supplying the typed accessors
// and arithmetic that raw any values lack.
package value

import (
	"reflect"
	"strconv"
)

// Value is the dynamic value: nil | bool | int64 | float64 | string | []Value |
// map[string]Value. The alias keeps encoding/json interop direct (a decoded
// JSON document is already a Value); the accessors below supply the typed,
// constant-error discipline. Numbers decoded from JSON are float64; numeric
// literals may be int64. The accessors and arithmetic treat both uniformly.
type Value = any

// Kind names the dynamic type of a Value.
type Kind int

const (
	// KindNull is the kind of nil.
	KindNull Kind = iota
	// KindBool is the kind of a bool.
	KindBool
	// KindInt is the kind of an int64.
	KindInt
	// KindFloat is the kind of a float64.
	KindFloat
	// KindString is the kind of a string.
	KindString
	// KindList is the kind of a []Value.
	KindList
	// KindObject is the kind of a map[string]Value.
	KindObject
)

// KindOf reports the Kind of v. An unrecognized concrete type reports KindNull.
func KindOf(v Value) Kind {
	switch v.(type) {
	case bool:
		return KindBool
	case int64:
		return KindInt
	case float64:
		return KindFloat
	case string:
		return KindString
	case []Value:
		return KindList
	case map[string]Value:
		return KindObject
	}
	return KindNull
}

// AsObject returns v as an object, or ErrNotObject.
func AsObject(v Value) (map[string]Value, error) {
	if o, ok := v.(map[string]Value); ok {
		return o, nil
	}
	return nil, ErrNotObject
}

// AsList returns v as a list, or ErrNotList.
func AsList(v Value) ([]Value, error) {
	if l, ok := v.([]Value); ok {
		return l, nil
	}
	return nil, ErrNotList
}

// AsString returns v as a string, or ErrNotString.
func AsString(v Value) (string, error) {
	if s, ok := v.(string); ok {
		return s, nil
	}
	return "", ErrNotString
}

// AsBool returns v as a bool, or ErrNotBool.
func AsBool(v Value) (bool, error) {
	if b, ok := v.(bool); ok {
		return b, nil
	}
	return false, ErrNotBool
}

// AsInt returns v as an int64 (truncating a float), or ErrNotNumber.
func AsInt(v Value) (int64, error) {
	switch n := v.(type) {
	case int64:
		return n, nil
	case float64:
		return int64(n), nil
	}
	return 0, ErrNotNumber
}

// AsFloat returns v as a float64, or ErrNotNumber.
func AsFloat(v Value) (float64, error) {
	switch n := v.(type) {
	case float64:
		return n, nil
	case int64:
		return float64(n), nil
	}
	return 0, ErrNotNumber
}

// Truthy reports truthiness: nil and false are falsey; everything else
// (including zero and the empty string) is truthy.
func Truthy(v Value) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	}
	return true
}

// bothNumbers reports whether a and b are both numeric, returning float views.
func bothNumbers(a, b Value) (x, y float64, isOk bool) {
	af, aerr := AsFloat(a)
	bf, berr := AsFloat(b)
	if aerr != nil || berr != nil {
		return 0, 0, false
	}
	return af, bf, true
}

// Equal reports value equality, treating int and float as comparable numbers.
func Equal(a, b Value) bool {
	if x, y, ok := bothNumbers(a, b); ok {
		return x == y
	}
	// reflect.DeepEqual handles every kind, including the uncomparable []Value and
	// map[string]Value (a direct == on those panics) and distinct types.
	return reflect.DeepEqual(a, b)
}

// Compare orders a and b, returning -1, 0, or 1. Numbers compare across
// int/float; strings compare lexically; any other pairing is ErrIncomparable.
func Compare(a, b Value) (int, error) {
	if x, y, ok := bothNumbers(a, b); ok {
		return sign(delta(x - y)), nil
	}
	as, aerr := AsString(a)
	bs, berr := AsString(b)
	if aerr == nil && berr == nil {
		return sign3(as < bs, as > bs), nil
	}
	return 0, ErrIncomparable
}

// delta is the signed difference between two compared numbers; only its sign
// is meaningful.
type delta float64

// sign returns the sign of d as -1, 0, or 1.
func sign(d delta) int {
	return sign3(float64(d) < 0, float64(d) > 0)
}

// sign3 collapses a isLess/isGreater pair into -1, 0, or 1.
func sign3(isLess, isGreater bool) int {
	switch {
	case isLess:
		return -1
	case isGreater:
		return 1
	}
	return 0
}

// Add adds two numbers, or concatenates when either operand is a string. Any
// other pairing is ErrNotNumber.
func Add(a, b Value) (Value, error) {
	if KindOf(a) == KindString || KindOf(b) == KindString {
		return concat(a, b)
	}
	return numericAdd(a, b)
}

// concat string-coerces both operands and joins them.
func concat(a, b Value) (Value, error) {
	as, err := coerceString(a)
	if err != nil {
		return nil, err
	}
	bs, err := coerceString(b)
	if err != nil {
		return nil, err
	}
	return as + bs, nil
}

// numericAdd keeps int when both operands are int, else promotes to float.
func numericAdd(a, b Value) (Value, error) {
	if ai, aok := a.(int64); aok {
		if bi, bok := b.(int64); bok {
			return ai + bi, nil
		}
	}
	x, y, ok := bothNumbers(a, b)
	if !ok {
		return nil, ErrNotNumber
	}
	return x + y, nil
}

// coerceString renders a scalar as a string; a composite is ErrNotString.
func coerceString(v Value) (string, error) {
	switch t := v.(type) {
	case string:
		return t, nil
	case int64:
		return strconv.FormatInt(t, 10), nil
	case float64:
		return strconv.FormatFloat(t, 'g', -1, 64), nil
	case bool:
		return strconv.FormatBool(t), nil
	}
	return "", ErrNotString
}

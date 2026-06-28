package value

// Error is the sentinel error type for the value package. Every error the
// package emits is a constant of this type, matched with errors.Is.
type Error string

// Error renders the sentinel message.
func (e Error) Error() string { return string(e) }

const (
	// ErrNotObject is returned when a value is required to be an object.
	ErrNotObject Error = "value: not an object"
	// ErrNotList is returned when a value is required to be a list.
	ErrNotList Error = "value: not a list"
	// ErrNotString is returned when a value is required to be a string.
	ErrNotString Error = "value: not a string"
	// ErrNotNumber is returned when a value is required to be numeric.
	ErrNotNumber Error = "value: not a number"
	// ErrNotBool is returned when a value is required to be a bool.
	ErrNotBool Error = "value: not a bool"
	// ErrIncomparable is returned when two values cannot be ordered.
	ErrIncomparable Error = "value: incomparable types"
)

package errors

import "fmt"

// Kind classifies configuration failures.
type Kind string

const (
	KindInvalidTarget Kind = "invalid_target"
	KindRequired      Kind = "required"
	KindParse         Kind = "parse"
	KindUnsupported   Kind = "unsupported"
)

// Error is a human-friendly configuration error.
// It contains enough structured data for tests and for rich formatting.
type Error struct {
	Kind     Kind   // required, parse, etc.
	Var      string // env var name (e.g. PORT)
	Field    string // struct field name (e.g. Port)
	Expected string // expected type (e.g. int)
	Value    string // raw string value (e.g. "abc")
	Cause    error  // underlying error (optional)
}

func (e *Error) Unwrap() error { return e.Cause }

// Code-like constructors (nice for internal use).
func Required(varName, field string) *Error {
	return &Error{Kind: KindRequired, Var: varName, Field: field}
}

func InvalidTarget(msg string) *Error {
	return &Error{Kind: KindInvalidTarget, Cause: fmt.Errorf(msg)}
}

func Parse(varName, field, expected, value string, cause error) *Error {
	return &Error{
		Kind:     KindParse,
		Var:      varName,
		Field:    field,
		Expected: expected,
		Value:    value,
		Cause:    cause,
	}
}

func Unsupported(varName, field, expected string) *Error {
	return &Error{
		Kind:     KindUnsupported,
		Var:      varName,
		Field:    field,
		Expected: expected,
	}
}

package model

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidRequirement is returned (wrapped) for a requirement that can never be validated.
	// The validations of a struct are checked on extraction (on the first call for a struct type), eg. a requirement
	// that can not be parsed, an invalid group, `equ` with a value that is not a number on an int field
	// or `cmn` on a field that is not a struct. While validating it is returned for an unknown function (`fun`),
	// an unknown condition type or `rel` or `$` used without the other fields (eg. in validator.ValidateValueWithParser).
	// It is a mistake in the code, not in the input, so a handler can answer it with 500.
	ErrInvalidRequirement = errors.New("invalid requirement")
	// ErrUnsupportedMediaType is returned by validator.UnpackAndValidate (wrapped) if the Content-Type of the request
	// is not supported. A handler can answer it with 415.
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	// ErrMalformedInput is returned by all unpack functions (wrapped) if the input can not be read or parsed,
	// eg. invalid json or yaml, an invalid query or a body that is larger than allowed by http.MaxBytesReader
	// (the *http.MaxBytesError is wrapped too). A handler can answer it with 400 (or 413).
	ErrMalformedInput = errors.New("malformed input")
)

// FieldErrors holds the errors of all invalid fields of a validation.
// All errors that are not one of the sentinel errors above are errors of the validation.
type FieldErrors []error

// Error joins the messages of all field errors.
func (e FieldErrors) Error() string {
	messages := make([]string, len(e))
	for i, err := range e {
		messages[i] = err.Error()
	}
	return strings.Join(messages, ", ")
}

// Unwrap returns the single field errors, so errors.Is and errors.As work on them.
func (e FieldErrors) Unwrap() []error {
	return e
}

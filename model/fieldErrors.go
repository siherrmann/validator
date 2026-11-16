package model

import "strings"

// FieldErrors holds the errors of all invalid fields of a validation.
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

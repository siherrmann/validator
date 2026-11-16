package model

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldErrorsError(t *testing.T) {
	tests := []struct {
		name     string
		input    FieldErrors
		expected string
	}{
		{
			name:     "Single field error",
			input:    FieldErrors{errors.New("field a invalid")},
			expected: "field a invalid",
		},
		{
			name:     "Multiple field errors",
			input:    FieldErrors{errors.New("field a invalid"), errors.New("field b invalid")},
			expected: "field a invalid, field b invalid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.input.Error(), "Expected error message to match")
		})
	}
}

func TestFieldErrorsUnwrap(t *testing.T) {
	errA := errors.New("field a invalid")
	errB := errors.New("field b invalid")

	t.Run("Unwrap returns all field errors", func(t *testing.T) {
		assert.Equal(t, []error{errA, errB}, FieldErrors{errA, errB}.Unwrap(), "Expected unwrapped errors to match")
	})

	t.Run("errors.Is finds a single field error", func(t *testing.T) {
		err := error(FieldErrors{errA, errB})
		assert.ErrorIs(t, err, errB, "Expected errors.Is to find a single field error")
	})
}

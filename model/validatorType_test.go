package model

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatorTypeToReflectKind(t *testing.T) {
	tests := []struct {
		name     string
		input    ValidatorType
		expected reflect.Kind
	}{
		{"String", String, reflect.String},
		{"Int", Int, reflect.Int},
		{"Float", Float, reflect.Float64},
		{"Bool", Bool, reflect.Bool},
		{"Array", Array, reflect.Array},
		{"Map", Map, reflect.Map},
		{"Struct", Struct, reflect.Struct},
		{"Unknown", "unknown", reflect.Struct},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.ToReflectKind()
			assert.Equal(t, test.expected, result, "Expected %v for %v, got %v", test.expected, test.input, result)
		})
	}
}

func TestReflectKindToValidatorType(t *testing.T) {
	tests := []struct {
		name     string
		input    reflect.Kind
		expected ValidatorType
	}{
		{"String", reflect.String, String},
		{"Int", reflect.Int, Int},
		{"Float", reflect.Float64, Float},
		{"Bool", reflect.Bool, Bool},
		{"Array", reflect.Array, Array},
		{"Map", reflect.Map, Map},
		{"Struct", reflect.Struct, Struct},
		{"Unknown", reflect.Invalid, Struct},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ReflectKindToValidatorType(test.input)
			assert.Equal(t, test.expected, result, "Expected %v for %v, got %v", test.expected, test.input, result)
		})
	}
}

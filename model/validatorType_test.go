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
		expected reflect.Type
	}{
		{"String", String, reflect.TypeOf("")},
		{"Int", Int, reflect.TypeOf(int(0))},
		{"Float", Float, reflect.TypeOf(float64(0))},
		{"Bool", Bool, reflect.TypeOf(false)},
		{"Array", Array, reflect.TypeOf([]string{})},
		{"Map", Map, reflect.TypeOf(map[string]string{})},
		{"Struct", Struct, reflect.TypeOf(struct{}{})},
		{"Unknown", "unknown", reflect.TypeOf(struct{}{})},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.input.ToReflectReflectType()
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

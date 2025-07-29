package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArrayToArrayOfAny(t *testing.T) {
	t.Run("Successfully extract array values as array", func(t *testing.T) {
		value := []int{1, 2}
		mapKeys, err := ArrayToArrayOfAny(value)
		assert.NoError(t, err, "Expected no error on map key extraction")
		assert.Equal(t, []any{1, 2}, mapKeys, "Expected array to contain map keys")
	})

	t.Run("Invalid type for array value extraction", func(t *testing.T) {
		value := map[int]string{1: "apple", 2: "banana"}
		mapKeys, err := ArrayToArrayOfAny(value)
		assert.Error(t, err, "Expected error on map key extraction")
		assert.Equal(t, []interface{}(nil), mapKeys, "Expected array to contain map keys")
	})
}

func TestArrayToArrayOfType(t *testing.T) {
	t.Run("Successfully convert int array", func(t *testing.T) {
		value := []any{1, 2, 3}
		expectedType := reflect.TypeOf(0)

		arrayValue, err := ArrayToArrayOfType(value, expectedType)
		assert.NoError(t, err, "Expected no error on conversion")
		assert.Equal(t, reflect.TypeOf([]int{1, 2, 3}), arrayValue.Type(), "Expected converted array to be of type []int")
	})

	t.Run("Successfully convert string to int array", func(t *testing.T) {
		value := []any{"1", "2", "3"}
		expectedType := reflect.TypeOf(0)

		arrayValue, err := ArrayToArrayOfType(value, expectedType)
		assert.NoError(t, err, "Expected no error on conversion")
		assert.Equal(t, reflect.TypeOf([]int{1, 2, 3}), arrayValue.Type(), "Expected converted array to be of type []int")
	})

	t.Run("Error on conversion string to int array", func(t *testing.T) {
		value := []any{"apple", "banana"}
		expectedType := reflect.TypeOf(0)

		arrayValue, err := ArrayToArrayOfType(value, expectedType)
		assert.Error(t, err, "Expected error on conversion with incompatible type")
		assert.Equal(t, reflect.Value{}, arrayValue, "Expected empty value on error")
	})
}

package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapKeysToArrayOfAny(t *testing.T) {
	t.Run("Successfully extract map keys as array", func(t *testing.T) {
		value := map[int]string{1: "apple", 2: "banana"}
		mapKeys, err := MapKeysToArrayOfAny(value)
		assert.NoError(t, err, "Expected no error on map key extraction")
		assert.ElementsMatch(t, []any{1, 2}, mapKeys, "Expected array to contain map keys")
	})

	t.Run("Invalid type for map key extraction", func(t *testing.T) {
		value := []int{1, 2}
		mapKeys, err := MapKeysToArrayOfAny(value)
		assert.Error(t, err, "Expected error on map key extraction")
		assert.Equal(t, []interface{}(nil), mapKeys, "Expected array to contain map keys")
	})

	t.Run("Invalid nil value for map key extraction", func(t *testing.T) {
		mapKeys, err := MapKeysToArrayOfAny(nil)
		assert.Error(t, err, "Expected error on nil value")
		assert.Equal(t, []interface{}(nil), mapKeys, "Expected nil result")
	})
}

func TestMapValue(t *testing.T) {
	value, ok := MapValue(map[string]any{"a": 1}, "a")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = MapValue(map[string]string{"a": "b"}, "a")
	assert.True(t, ok)
	assert.Equal(t, "b", value)

	_, ok = MapValue(map[string]any{"a": 1}, "b")
	assert.False(t, ok, "Expected missing key")

	_, ok = MapValue(map[int]string{1: "a"}, "1")
	assert.False(t, ok, "Expected no value for a map without string keys")

	_, ok = MapValue([]string{"a"}, "a")
	assert.False(t, ok, "Expected no value for a non map")
}

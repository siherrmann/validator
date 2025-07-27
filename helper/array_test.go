package helper

import (
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

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
}

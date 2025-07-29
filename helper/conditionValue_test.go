package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConditionValueToT(t *testing.T) {
	t.Run("Successfully run condition conversion", func(t *testing.T) {
		value := 1
		condition := "2"
		conditionConverted, err := ConditionValueToT(value, condition)
		assert.NoError(t, err, "Expected no error for condition conversion")
		assert.Equal(t, 2, conditionConverted, "Expected converted condition to be 2")
	})

	t.Run("Invalid type for condition conversion", func(t *testing.T) {
		value := 1
		condition := "apple"
		conditionConverted, err := ConditionValueToT(value, condition)
		assert.Error(t, err, "Expected error for condition conversion")
		assert.Equal(t, 1, conditionConverted, "Expected converted condition to be 0")
	})
}

func TestConditionValueToArrayOfAny(t *testing.T) {
	t.Run("Successfully run condition conversion", func(t *testing.T) {
		value := 1
		condition := "2,3"
		conditionConverted, err := ConditionValueToArrayOfAny(condition, reflect.TypeOf(value))
		assert.NoError(t, err, "Expected no error for condition conversion")
		assert.Equal(t, []interface{}{2, 3}, conditionConverted, "Expected converted condition to be 2")
	})

	t.Run("Empty condition value", func(t *testing.T) {
		value := 1
		condition := ""
		conditionConverted, err := ConditionValueToArrayOfAny(condition, reflect.TypeOf(value))
		assert.Error(t, err, "Expected error for condition conversion")
		assert.Equal(t, []interface{}{}, conditionConverted, "Expected converted condition to be nil array")
	})

	t.Run("Invalid type for condition conversion", func(t *testing.T) {
		value := 1
		condition := "apple,banana"
		conditionConverted, err := ConditionValueToArrayOfAny(condition, reflect.TypeOf(value))
		assert.Error(t, err, "Expected error for condition conversion")
		assert.Equal(t, []interface{}(nil), conditionConverted, "Expected converted condition to be nil array")
	})
}

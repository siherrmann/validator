package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsString(t *testing.T) {
	testString := IsString("test")
	assert.True(t, testString, "expected true for string, got false")

	testInt := IsString(123)
	assert.False(t, testInt, "expected false for int, got true")
}

func TestIsArray(t *testing.T) {
	testArray := IsArray([]int{1, 2, 3})
	assert.True(t, testArray, "expected true for array, got false")

	testSlice := IsArray([]string{"a", "b", "c"})
	assert.True(t, testSlice, "expected true for slice, got false")

	testInt := IsArray(123)
	assert.False(t, testInt, "expected false for int, got true")
}

func TestIsArrayOfStruct(t *testing.T) {
	type TestStruct struct {
		Field string
	}
	testArrayOfStruct := IsArrayOfStruct([]TestStruct{{Field: "value1"}, {Field: "value2"}})
	assert.True(t, testArrayOfStruct, "expected true for array of structs, got false")

	testArrayOfInt := IsArrayOfStruct([]int{1, 2, 3})
	assert.False(t, testArrayOfInt, "expected false for array of ints, got true")

	testString := IsArrayOfStruct("not an array")
	assert.False(t, testString, "expected false for string, got true")
}

func TestIsArrayOfMap(t *testing.T) {
	testArrayOfMap := IsArrayOfMap([]map[string]int{{"key": 1}, {"key": 2}})
	assert.True(t, testArrayOfMap, "expected true for array of maps, got false")

	testArrayOfStruct := IsArrayOfMap([]struct{ Key string }{{Key: "value1"}, {Key: "value2"}})
	assert.False(t, testArrayOfStruct, "expected false for array of structs, got true")
}

func TestCheckValidPointerToStruct(t *testing.T) {
	type TestStruct struct {
		Field string
	}

	validPointer := &TestStruct{}
	err := CheckValidPointerToStruct(validPointer)
	assert.NoError(t, err, "expected no error for valid pointer to struct")

	invalidPointer := TestStruct{}
	err = CheckValidPointerToStruct(invalidPointer)
	assert.Error(t, err, "expected error for invalid pointer to struct")

	invalidType := &[]int{1, 2, 3}
	err = CheckValidPointerToStruct(invalidType)
	assert.Error(t, err, "expected error for invalid type")
}

func TestArrayOfInterfaceToArrayOf(t *testing.T) {
	t.Run("Successfully convert to array of int", func(t *testing.T) {
		input := []any{1, 2, 3}
		expected := []int{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[int](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.Equal(t, expected[i], v, "expected value at index %d to match expected value", i)
		}

		input = []any{float64(1), float64(2), float64(3)}
		expected = []int{1, 2, 3}

		result, err = ArrayOfInterfaceToArrayOf[int](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of int with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[int](input)
		assert.Error(t, err, "expected error converting array of interface to array of int with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of int8", func(t *testing.T) {
		input := []any{int8(1), int8(2), int8(3)}
		expected := []int8{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[int8](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int8")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}

		input = []any{float64(1), float64(2), float64(3)}
		expected = []int8{1, 2, 3}

		result, err = ArrayOfInterfaceToArrayOf[int8](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int8")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of int8 with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[int8](input)
		assert.Error(t, err, "expected error converting array of interface to array of int8 with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of int16", func(t *testing.T) {
		input := []any{int16(1), int16(2), int16(3)}
		expected := []int16{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[int16](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int16")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}

		input = []any{float64(1), float64(2), float64(3)}
		expected = []int16{1, 2, 3}

		result, err = ArrayOfInterfaceToArrayOf[int16](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int16")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of int16 with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[int16](input)
		assert.Error(t, err, "expected error converting array of interface to array of int16 with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of int32", func(t *testing.T) {
		input := []any{int32(1), int32(2), int32(3)}
		expected := []int32{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[int32](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int32")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}

		input = []any{float64(1), float64(2), float64(3)}
		expected = []int32{1, 2, 3}

		result, err = ArrayOfInterfaceToArrayOf[int32](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int32")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of int32 with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[int32](input)
		assert.Error(t, err, "expected error converting array of interface to array of int32 with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of int64", func(t *testing.T) {
		input := []any{int64(1), int64(2), int64(3)}
		expected := []int64{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[int64](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int64")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}

		input = []any{float64(1), float64(2), float64(3)}
		expected = []int64{1, 2, 3}

		result, err = ArrayOfInterfaceToArrayOf[int64](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of int64")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of int64 with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[int64](input)
		assert.Error(t, err, "expected error converting array of interface to array of int64 with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of float64", func(t *testing.T) {
		input := []any{1.1, 2.2, 3.3}
		expected := []float64{1.1, 2.2, 3.3}

		result, err := ArrayOfInterfaceToArrayOf[float64](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of float64")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.EqualValues(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of float64 with invalid input type", func(t *testing.T) {
		input := []any{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[float64](input)
		assert.Error(t, err, "expected error converting array of interface to array of float64 with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})

	t.Run("Successfully convert to array of string", func(t *testing.T) {
		input := []any{"one", "two", "three"}
		expected := []string{"one", "two", "three"}

		result, err := ArrayOfInterfaceToArrayOf[string](input)
		assert.NoError(t, err, "expected no error converting array of interface to array of string")
		assert.Equal(t, len(result), len(expected), "expected length of result to match expected length")
		for i, v := range result {
			assert.Equal(t, expected[i], v, "expected value at index %d to match expected value", i)
		}
	})

	t.Run("Error converting to array of string with invalid input type", func(t *testing.T) {
		input := []any{1, 2, 3}

		result, err := ArrayOfInterfaceToArrayOf[string](input)
		assert.Error(t, err, "expected error converting array of interface to array of string with invalid input type")
		assert.Empty(t, result, "expected result to be empty on error")
	})
}

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

func TestIsStruct(t *testing.T) {
	type TestStruct struct {
		Field string
	}

	testStruct := IsStruct(TestStruct{Field: "value"})
	assert.True(t, testStruct, "expected true for struct, got false")

	testInt := IsStruct(123)
	assert.False(t, testInt, "expected false for int, got true")

	testString := IsStruct("not a struct")
	assert.False(t, testString, "expected false for string, got true")

	testPointer := IsStruct(&TestStruct{Field: "value"})
	assert.False(t, testPointer, "expected false for pointer to struct, got true")

	testSlice := IsStruct([]TestStruct{{Field: "value"}})
	assert.False(t, testSlice, "expected false for slice of structs, got true")
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

package helper

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAnyToFloat(t *testing.T) {
	tests := []struct {
		name          string
		arg           any
		expected      float64
		expectedError bool
	}{
		{
			name:          "Valid string",
			arg:           "apple",
			expected:      5.0,
			expectedError: false,
		},
		{
			name:          "Valid time",
			arg:           time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected:      1.735776e+09,
			expectedError: false,
		},
		{
			name:          "Valid int",
			arg:           int(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid int8",
			arg:           int8(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid int16",
			arg:           int16(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid int32",
			arg:           int32(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid int64",
			arg:           int64(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid uint",
			arg:           uint(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid uint8",
			arg:           uint8(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid uint16",
			arg:           uint16(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid uint32",
			arg:           uint32(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid uint64",
			arg:           uint64(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid float32",
			arg:           float32(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid float64",
			arg:           float64(1),
			expected:      1,
			expectedError: false,
		},
		{
			name:          "Valid array",
			arg:           []float64{1, 2, 3},
			expected:      3,
			expectedError: false,
		},
		{
			name:          "Invalid type",
			arg:           struct{}{},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := AnyToFloat(test.arg)
			if test.expectedError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, test.expected, result, "Expected result to be equal to expected value")
			}
		})
	}
}

func TestAnyToString(t *testing.T) {
	tests := []struct {
		name          string
		arg           any
		expected      string
		expectedError bool
	}{
		{
			name:          "Valid string",
			arg:           "apple",
			expected:      "apple",
			expectedError: false,
		},
		{
			name:          "Valid bool",
			arg:           true,
			expected:      "true",
			expectedError: false,
		},
		{
			name:          "Valid time",
			arg:           time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			expected:      "1735776000",
			expectedError: false,
		},
		{
			name:          "Valid int",
			arg:           int(1),
			expected:      "1",
			expectedError: false,
		},
		{
			name:          "Valid uint",
			arg:           uint(1),
			expected:      "1",
			expectedError: false,
		},
		{
			name:          "Valid float32",
			arg:           float32(1),
			expected:      "1.000000",
			expectedError: false,
		},
		{
			name:          "Invalid type",
			arg:           struct{}{},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := AnyToString(test.arg)
			if test.expectedError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, test.expected, result, "Expected result to be equal to expected value")
			}
		})
	}
}

func TestAnyToArrayOfString(t *testing.T) {
	tests := []struct {
		name          string
		arg           any
		expected      []string
		expectedError bool
	}{
		{
			name:          "Valid array string",
			arg:           []string{"apple"},
			expected:      []string{"apple"},
			expectedError: false,
		},
		{
			name:          "Valid array bool",
			arg:           []bool{true},
			expected:      []string{"true"},
			expectedError: false,
		},
		{
			name:          "Valid array time",
			arg:           []time.Time{time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
			expected:      []string{"1735776000"},
			expectedError: false,
		},
		{
			name:          "Valid array int",
			arg:           []int{1},
			expected:      []string{"1"},
			expectedError: false,
		},
		{
			name:          "Valid array uint",
			arg:           []uint{1},
			expected:      []string{"1"},
			expectedError: false,
		},
		{
			name:          "Valid array float32",
			arg:           []float32{1},
			expected:      []string{"1.000000"},
			expectedError: false,
		},
		{
			name:          "Invalid type array",
			arg:           []struct{}{{}},
			expectedError: true,
		},
		{
			name:          "Valid map string string",
			arg:           map[string]string{"1": "apple"},
			expected:      []string{"1"},
			expectedError: false,
		},
		{
			name:          "Invalid type map",
			arg:           map[struct{}]string{{}: "apple"},
			expectedError: true,
		},
		{
			name:          "Invalid type",
			arg:           struct{}{},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := AnyToArrayOfString(test.arg)
			if test.expectedError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, test.expected, result, "Expected result to be equal to expected value")
			}
		})
	}
}

func TestAnyToType(t *testing.T) {
	type testStruct struct {
		Fruit string `json:"fruit"`
	}
	type nestedStruct struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}
	type structWithNestedStruct struct {
		Title         string        `json:"title"`
		Nested        nestedStruct  `json:"nested"`
		NestedPointer *nestedStruct `json:"nested_pointer"`
	}
	type args struct {
		v        any
		expected reflect.Type
	}
	tests := []struct {
		name          string
		args          args
		expected      any
		expectedError bool
	}{
		// string
		{
			name: "Valid string to string",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(string("")),
			},
			expected:      any("apple"),
			expectedError: false,
		},
		// bool
		{
			name: "Valid bool to bool",
			args: args{
				v:        true,
				expected: reflect.TypeOf(bool(false)),
			},
			expected:      any(true),
			expectedError: false,
		},
		{
			name: "Valid string on to bool",
			args: args{
				v:        "on",
				expected: reflect.TypeOf(bool(false)),
			},
			expected:      any(true),
			expectedError: false,
		},
		{
			name: "Valid string off to bool",
			args: args{
				v:        "off",
				expected: reflect.TypeOf(bool(false)),
			},
			expected:      any(false),
			expectedError: false,
		},
		{
			name: "Valid string to bool",
			args: args{
				v:        "true",
				expected: reflect.TypeOf(bool(false)),
			},
			expected:      any(true),
			expectedError: false,
		},
		{
			name: "Invalid string to bool",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(bool(false)),
			},
			expectedError: true,
		},
		// int
		{
			name: "Valid int to int",
			args: args{
				v:        int(1),
				expected: reflect.TypeOf(int(0)),
			},
			expected:      any(1),
			expectedError: false,
		},
		{
			name: "Valid float to int",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(int(0)),
			},
			expected:      any(1),
			expectedError: false,
		},
		{
			name: "Valid string to int",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(int(0)),
			},
			expected:      any(1),
			expectedError: false,
		},
		{
			name: "Invalid string to int",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(int(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int8 to int",
			args: args{
				v:        int8(1),
				expected: reflect.TypeOf(int(0)),
			},
			expected:      any(1),
			expectedError: false,
		},
		// int8
		{
			name: "Valid int8 to int8",
			args: args{
				v:        int8(1),
				expected: reflect.TypeOf(int8(0)),
			},
			expected:      any(int8(1)),
			expectedError: false,
		},
		{
			name: "Valid float to int8",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(int8(0)),
			},
			expected:      any(int8(1)),
			expectedError: false,
		},
		{
			name: "Valid string to int8",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(int8(0)),
			},
			expected:      any(int8(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to int8",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(int8(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int16 to int8",
			args: args{
				v:        int16(1),
				expected: reflect.TypeOf(int8(0)),
			},
			expected:      any(int8(1)),
			expectedError: false,
		},
		// int16
		{
			name: "Valid int16 to int16",
			args: args{
				v:        int16(1),
				expected: reflect.TypeOf(int16(0)),
			},
			expected:      any(int16(1)),
			expectedError: false,
		},
		{
			name: "Valid float to int16",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(int16(0)),
			},
			expected:      any(int16(1)),
			expectedError: false,
		},
		{
			name: "Valid string to int16",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(int16(0)),
			},
			expected:      any(int16(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to int16",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(int16(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int32 to int16",
			args: args{
				v:        int32(1),
				expected: reflect.TypeOf(int16(0)),
			},
			expected:      any(int16(1)),
			expectedError: false,
		},
		// int32
		{
			name: "Valid int32 to int32",
			args: args{
				v:        int32(1),
				expected: reflect.TypeOf(int32(0)),
			},
			expected:      any(int32(1)),
			expectedError: false,
		},
		{
			name: "Valid float to int32",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(int32(0)),
			},
			expected:      any(int32(1)),
			expectedError: false,
		},
		{
			name: "Valid string to int32",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(int32(0)),
			},
			expected:      any(int32(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to int32",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(int32(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int64 to int32",
			args: args{
				v:        int64(1),
				expected: reflect.TypeOf(int32(0)),
			},
			expected:      any(int32(1)),
			expectedError: false,
		},
		// int64
		{
			name: "Valid int64 to int64",
			args: args{
				v:        int64(1),
				expected: reflect.TypeOf(int64(0)),
			},
			expected:      any(int64(1)),
			expectedError: false,
		},
		{
			name: "Valid float to int64",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(int64(0)),
			},
			expected:      any(int64(1)),
			expectedError: false,
		},
		{
			name: "Valid string to int64",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(int64(0)),
			},
			expected:      any(int64(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to int64",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(int64(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int to int64",
			args: args{
				v:        int(1),
				expected: reflect.TypeOf(int64(0)),
			},
			expected:      any(int64(1)),
			expectedError: false,
		},
		// uint
		{
			name: "Valid uint to uint",
			args: args{
				v:        uint(1),
				expected: reflect.TypeOf(uint(0)),
			},
			expected:      any(uint(1)),
			expectedError: false,
		},
		{
			name: "Valid float to uint",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(uint(0)),
			},
			expected:      any(uint(1)),
			expectedError: false,
		},
		{
			name: "Valid string to uint",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(uint(0)),
			},
			expected:      any(uint(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to uint",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(uint(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid uint8 to uint",
			args: args{
				v:        uint8(1),
				expected: reflect.TypeOf(uint(0)),
			},
			expected:      any(uint(1)),
			expectedError: false,
		},
		// uint8
		{
			name: "Valid uint8 to uint8",
			args: args{
				v:        uint8(1),
				expected: reflect.TypeOf(uint8(0)),
			},
			expected:      any(uint8(1)),
			expectedError: false,
		},
		{
			name: "Valid float to uint8",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(uint8(0)),
			},
			expected:      any(uint8(1)),
			expectedError: false,
		},
		{
			name: "Valid string to uint8",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(uint8(0)),
			},
			expected:      any(uint8(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to uint8",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(uint8(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid uint16 to uint8",
			args: args{
				v:        uint16(1),
				expected: reflect.TypeOf(uint8(0)),
			},
			expected:      any(uint8(1)),
			expectedError: false,
		},
		// uint16
		{
			name: "Valid uint16 to uint16",
			args: args{
				v:        uint16(1),
				expected: reflect.TypeOf(uint16(0)),
			},
			expected:      any(uint16(1)),
			expectedError: false,
		},
		{
			name: "Valid float to uint16",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(uint16(0)),
			},
			expected:      any(uint16(1)),
			expectedError: false,
		},
		{
			name: "Valid string to uint16",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(uint16(0)),
			},
			expected:      any(uint16(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to uint16",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(uint16(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid uint32 to uint16",
			args: args{
				v:        uint32(1),
				expected: reflect.TypeOf(uint16(0)),
			},
			expected:      any(uint16(1)),
			expectedError: false,
		},
		// uint32
		{
			name: "Valid uint32 to uint32",
			args: args{
				v:        uint32(1),
				expected: reflect.TypeOf(uint32(0)),
			},
			expected:      any(uint32(1)),
			expectedError: false,
		},
		{
			name: "Valid float to uint32",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(uint32(0)),
			},
			expected:      any(uint32(1)),
			expectedError: false,
		},
		{
			name: "Valid string to uint32",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(uint32(0)),
			},
			expected:      any(uint32(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to uint32",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(uint32(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid uint64 to uint32",
			args: args{
				v:        uint64(1),
				expected: reflect.TypeOf(uint32(0)),
			},
			expected:      any(uint32(1)),
			expectedError: false,
		},
		// uint64
		{
			name: "Valid uint64 to uint64",
			args: args{
				v:        uint64(1),
				expected: reflect.TypeOf(uint64(0)),
			},
			expected:      any(uint64(1)),
			expectedError: false,
		},
		{
			name: "Valid float to uint64",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(uint64(0)),
			},
			expected:      any(uint64(1)),
			expectedError: false,
		},
		{
			name: "Valid string to uint64",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(uint64(0)),
			},
			expected:      any(uint64(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to uint64",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(uint64(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid uint32 to uint64",
			args: args{
				v:        uint32(1),
				expected: reflect.TypeOf(uint64(0)),
			},
			expected:      any(uint64(1)),
			expectedError: false,
		},
		// float32
		{
			name: "Valid float32 to float32",
			args: args{
				v:        float32(1),
				expected: reflect.TypeOf(float32(0)),
			},
			expected:      any(float32(1)),
			expectedError: false,
		},
		{
			name: "Valid float to float32",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(float32(0)),
			},
			expected:      any(float32(1)),
			expectedError: false,
		},
		{
			name: "Valid string to float32",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(float32(0)),
			},
			expected:      any(float32(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to float32",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(float32(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int8 to float32",
			args: args{
				v:        int8(1),
				expected: reflect.TypeOf(float32(0)),
			},
			expected:      any(float32(1)),
			expectedError: false,
		},
		// float64
		{
			name: "Valid float64 to float64",
			args: args{
				v:        float64(1),
				expected: reflect.TypeOf(float64(0)),
			},
			expected:      any(float64(1)),
			expectedError: false,
		},
		{
			name: "Valid float to float64",
			args: args{
				v:        float32(1),
				expected: reflect.TypeOf(float64(0)),
			},
			expected:      any(float64(1)),
			expectedError: false,
		},
		{
			name: "Valid string to float64",
			args: args{
				v:        "1",
				expected: reflect.TypeOf(float64(0)),
			},
			expected:      any(float64(1)),
			expectedError: false,
		},
		{
			name: "Invalid string to float64",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf(float64(0)),
			},
			expectedError: true,
		},
		{
			name: "Valid int8 to float64",
			args: args{
				v:        int8(1),
				expected: reflect.TypeOf(float64(0)),
			},
			expected:      any(float64(1)),
			expectedError: false,
		},
		{
			name: "Valid time.Time to time.Time",
			args: args{
				v:        time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				expected: reflect.TypeOf(time.Time{}),
			},
			expected:      any(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
			expectedError: false,
		},
		{
			name: "Valid float to time.Time",
			args: args{
				v:        float64(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC).Unix()),
				expected: reflect.TypeOf(time.Time{}),
			},
			expected:      any(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
			expectedError: false,
		},
		{
			name: "Valid string to time.Time",
			args: args{
				v:        `2025-01-02T00:00:00.000Z`,
				expected: reflect.TypeOf(time.Time{}),
			},
			expected:      any(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
			expectedError: false,
		},
		{
			name: "Valid string to struct",
			args: args{
				v:        `{"fruit": "apple"}`,
				expected: reflect.TypeOf(testStruct{}),
			},
			expected:      any(testStruct{Fruit: "apple"}),
			expectedError: false,
		},
		{
			name: "Valid map to nested struct",
			args: args{
				v: map[string]any{
					"title": "Main Title",
					"nested": map[string]any{
						"name":  "Nested Name",
						"count": float64(42),
					},
				},
				expected: reflect.TypeOf(structWithNestedStruct{}),
			},
			expected:      any(structWithNestedStruct{Title: "Main Title", Nested: nestedStruct{Name: "Nested Name", Count: 42}}),
			expectedError: false,
		},
		{
			name: "Valid string JSON to nested struct",
			args: args{
				v:        `{"title": "Main Title", "nested": {"name": "Nested Name", "count": 42}}`,
				expected: reflect.TypeOf(structWithNestedStruct{}),
			},
			expected:      any(structWithNestedStruct{Title: "Main Title", Nested: nestedStruct{Name: "Nested Name", Count: 42}}),
			expectedError: false,
		},
		{
			name: "Valid map to struct with nested pointer",
			args: args{
				v: map[string]any{
					"title": "Main Title",
					"nested_pointer": map[string]any{
						"name":  "Nested Name",
						"count": float64(42),
					},
				},
				expected: reflect.TypeOf(structWithNestedStruct{}),
			},
			expected:      any(structWithNestedStruct{Title: "Main Title", NestedPointer: &nestedStruct{Name: "Nested Name", Count: 42}}),
			expectedError: false,
		},
		{
			name: "Valid string JSON to struct with nested pointer",
			args: args{
				v:        `{"title": "Main Title", "nested_pointer": {"name": "Nested Name", "count": 42}}`,
				expected: reflect.TypeOf(structWithNestedStruct{}),
			},
			expected:      any(structWithNestedStruct{Title: "Main Title", NestedPointer: &nestedStruct{Name: "Nested Name", Count: 42}}),
			expectedError: false,
		},
		{
			name: "Valid map to struct with nested pointer (nested_pointer is nil)",
			args: args{
				v: map[string]any{
					"title": "Main Title",
				},
				expected: reflect.TypeOf(structWithNestedStruct{}),
			},
			expected:      any(structWithNestedStruct{Title: "Main Title", NestedPointer: nil}),
			expectedError: false,
		},
		{
			name: "Invalid type",
			args: args{
				v:        args{},
				expected: reflect.TypeOf(struct{}{}),
			},
			expectedError: true,
		},
		// Pointer types
		{
			name: "Valid string to *string",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf((*string)(nil)),
			},
			expected:      func() any { s := "apple"; return &s }(),
			expectedError: false,
		},
		{
			name: "Valid int to *int",
			args: args{
				v:        42,
				expected: reflect.TypeOf((*int)(nil)),
			},
			expected:      func() any { i := 42; return &i }(),
			expectedError: false,
		},
		{
			name: "Valid float to *int",
			args: args{
				v:        float64(42),
				expected: reflect.TypeOf((*int)(nil)),
			},
			expected:      func() any { i := 42; return &i }(),
			expectedError: false,
		},
		{
			name: "Valid string to *int",
			args: args{
				v:        "42",
				expected: reflect.TypeOf((*int)(nil)),
			},
			expected:      func() any { i := 42; return &i }(),
			expectedError: false,
		},
		{
			name: "Valid bool to *bool",
			args: args{
				v:        true,
				expected: reflect.TypeOf((*bool)(nil)),
			},
			expected:      func() any { b := true; return &b }(),
			expectedError: false,
		},
		{
			name: "Valid string to *bool",
			args: args{
				v:        "true",
				expected: reflect.TypeOf((*bool)(nil)),
			},
			expected:      func() any { b := true; return &b }(),
			expectedError: false,
		},
		{
			name: "Valid float64 to *float64",
			args: args{
				v:        3.14,
				expected: reflect.TypeOf((*float64)(nil)),
			},
			expected:      func() any { f := 3.14; return &f }(),
			expectedError: false,
		},
		{
			name: "Valid string to *float64",
			args: args{
				v:        "3.14",
				expected: reflect.TypeOf((*float64)(nil)),
			},
			expected:      func() any { f := 3.14; return &f }(),
			expectedError: false,
		},
		{
			name: "Valid float32 to *float32",
			args: args{
				v:        float32(3.14),
				expected: reflect.TypeOf((*float32)(nil)),
			},
			expected:      func() any { f := float32(3.14); return &f }(),
			expectedError: false,
		},
		{
			name: "Valid int8 to *int8",
			args: args{
				v:        int8(42),
				expected: reflect.TypeOf((*int8)(nil)),
			},
			expected:      func() any { i := int8(42); return &i }(),
			expectedError: false,
		},
		{
			name: "Valid int16 to *int16",
			args: args{
				v:        int16(42),
				expected: reflect.TypeOf((*int16)(nil)),
			},
			expected:      func() any { i := int16(42); return &i }(),
			expectedError: false,
		},
		{
			name: "Valid int32 to *int32",
			args: args{
				v:        int32(42),
				expected: reflect.TypeOf((*int32)(nil)),
			},
			expected:      func() any { i := int32(42); return &i }(),
			expectedError: false,
		},
		{
			name: "Valid int64 to *int64",
			args: args{
				v:        int64(42),
				expected: reflect.TypeOf((*int64)(nil)),
			},
			expected:      func() any { i := int64(42); return &i }(),
			expectedError: false,
		},
		{
			name: "Valid uint to *uint",
			args: args{
				v:        uint(42),
				expected: reflect.TypeOf((*uint)(nil)),
			},
			expected:      func() any { u := uint(42); return &u }(),
			expectedError: false,
		},
		{
			name: "Valid uint8 to *uint8",
			args: args{
				v:        uint8(42),
				expected: reflect.TypeOf((*uint8)(nil)),
			},
			expected:      func() any { u := uint8(42); return &u }(),
			expectedError: false,
		},
		{
			name: "Valid uint16 to *uint16",
			args: args{
				v:        uint16(42),
				expected: reflect.TypeOf((*uint16)(nil)),
			},
			expected:      func() any { u := uint16(42); return &u }(),
			expectedError: false,
		},
		{
			name: "Valid uint32 to *uint32",
			args: args{
				v:        uint32(42),
				expected: reflect.TypeOf((*uint32)(nil)),
			},
			expected:      func() any { u := uint32(42); return &u }(),
			expectedError: false,
		},
		{
			name: "Valid uint64 to *uint64",
			args: args{
				v:        uint64(42),
				expected: reflect.TypeOf((*uint64)(nil)),
			},
			expected:      func() any { u := uint64(42); return &u }(),
			expectedError: false,
		},
		{
			name: "Valid time.Time to *time.Time",
			args: args{
				v:        time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				expected: reflect.TypeOf((*time.Time)(nil)),
			},
			expected:      func() any { t := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC); return &t }(),
			expectedError: false,
		},
		{
			name: "Valid string to *time.Time",
			args: args{
				v:        `2025-01-02T00:00:00.000Z`,
				expected: reflect.TypeOf((*time.Time)(nil)),
			},
			expected:      func() any { t := time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC); return &t }(),
			expectedError: false,
		},
		{
			name: "Valid struct to *struct",
			args: args{
				v:        `{"fruit": "apple"}`,
				expected: reflect.TypeOf((*testStruct)(nil)),
			},
			expected:      func() any { s := testStruct{Fruit: "apple"}; return &s }(),
			expectedError: false,
		},
		{
			name: "Nil to *string",
			args: args{
				v:        nil,
				expected: reflect.TypeOf((*string)(nil)),
			},
			expected:      (*string)(nil),
			expectedError: false,
		},
		{
			name: "Nil to *int",
			args: args{
				v:        nil,
				expected: reflect.TypeOf((*int)(nil)),
			},
			expected:      (*int)(nil),
			expectedError: false,
		},
		{
			name: "Invalid string to *int",
			args: args{
				v:        "apple",
				expected: reflect.TypeOf((*int)(nil)),
			},
			expectedError: true,
		},
		{
			name: "Invalid string to *bool",
			args: args{
				v:        "invalid",
				expected: reflect.TypeOf((*bool)(nil)),
			},
			expectedError: true,
		},
		// Array and Slice types
		{
			name: "Valid string to []byte",
			args: args{
				v:        "hello",
				expected: reflect.TypeOf([]byte{}),
			},
			expected:      any([]byte("hello")),
			expectedError: false,
		},
		{
			name: "Valid string to [5]byte",
			args: args{
				v:        "hello",
				expected: reflect.TypeOf([5]byte{}),
			},
			expected:      any([5]byte{'h', 'e', 'l', 'l', 'o'}),
			expectedError: false,
		},
		{
			name: "Valid string to [3]byte (truncated)",
			args: args{
				v:        "hello",
				expected: reflect.TypeOf([3]byte{}),
			},
			expected:      any([3]byte{'h', 'e', 'l'}),
			expectedError: false,
		},
		// Array/Slice from []any
		{
			name: "Valid []any to []int",
			args: args{
				v:        []any{1.0, 2.0, 3.0},
				expected: reflect.TypeOf([]int{}),
			},
			expected:      any([]int{1, 2, 3}),
			expectedError: false,
		},
		{
			name: "Valid []any to []string",
			args: args{
				v:        []any{"hello", "world"},
				expected: reflect.TypeOf([]string{}),
			},
			expected:      any([]string{"hello", "world"}),
			expectedError: false,
		},
		{
			name: "Valid []any with mixed convertible types to []int",
			args: args{
				v:        []any{"42", 100.0, int64(50)},
				expected: reflect.TypeOf([]int{}),
			},
			expected:      any([]int{42, 100, 50}),
			expectedError: false,
		},
		{
			name: "Valid []any to []float64",
			args: args{
				v:        []any{1.5, "2.5", int(3)},
				expected: reflect.TypeOf([]float64{}),
			},
			expected:      any([]float64{1.5, 2.5, 3.0}),
			expectedError: false,
		},
		{
			name: "Valid []any to []bool",
			args: args{
				v:        []any{true, "false", "on"},
				expected: reflect.TypeOf([]bool{}),
			},
			expected:      any([]bool{true, false, true}),
			expectedError: false,
		},
		{
			name: "Invalid []any to []int (contains non-convertible)",
			args: args{
				v:        []any{"hello", "world"},
				expected: reflect.TypeOf([]int{}),
			},
			expectedError: true,
		},
		{
			name: "Valid empty []any to []string",
			args: args{
				v:        []any{},
				expected: reflect.TypeOf([]string{}),
			},
			expected:      any([]string{}),
			expectedError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := AnyToType(test.args.v, test.args.expected)
			if test.expectedError {
				assert.Error(t, err, "Expected error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, test.expected, result, "Expected result to be equal to expected value")
			}
		})
	}
}

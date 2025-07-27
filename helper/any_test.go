package helper

import (
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
			expected:      0,
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
			expected:      "",
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
			expected:      []string{},
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
			expected:      []string{},
			expectedError: true,
		},
		{
			name:          "Invalid type",
			arg:           struct{}{},
			expected:      []string{},
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

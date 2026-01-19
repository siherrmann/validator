package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetValidationsFromStruct(t *testing.T) {
	type args struct {
		input   any
		tagType string
	}
	tests := []struct {
		name          string
		args          args
		expected      []Validation
		expectedError bool
	}{
		{
			name: "Valid struct with conditions and groups",
			args: args{
				input: &struct {
					Field1 string `vld:"equ1, gr1min1"`
					Field2 int    `vld:"min2, gr1min1"`
				}{},
				tagType: VLD,
			},
			expected: []Validation{
				{Key: "Field1", Type: String, Requirement: "equ1", Groups: []*Group{{Name: "gr1", ConditionType: "min", ConditionValue: "1"}}},
				{Key: "Field2", Type: Int, Requirement: "min2", Groups: []*Group{{Name: "gr1", ConditionType: "min", ConditionValue: "1"}}},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with custom tag",
			args: args{
				input: &struct {
					Field1 string `upd:"equ1"`
					Field2 int    `upd:"min2"`
				}{},
				tagType: "upd",
			},
			expected: []Validation{
				{Key: "Field1", Type: String, Requirement: "equ1"},
				{Key: "Field2", Type: Int, Requirement: "min2"},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with json tag",
			args: args{
				input: &struct {
					Field1 string `json:"field1" vld:"equ1"`
					Field2 int    `json:"field2" vld:"min2"`
				}{},
				tagType: VLD,
			},
			expected: []Validation{
				{Key: "field1", Type: String, Requirement: "equ1"},
				{Key: "field2", Type: Int, Requirement: "min2"},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with inner struct",
			args: args{
				input: &struct {
					Field1 struct {
						Name string `json:"name" vld:"equ1"`
					} `json:"field1" vld:"-"`
				}{},
				tagType: VLD,
			},
			expected: []Validation{
				{Key: "field1", Type: Struct, Requirement: "-", InnerValidation: []Validation{
					{Key: "name", Type: String, Requirement: "equ1"},
				}},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with array of structs",
			args: args{
				input: &struct {
					Field1 []struct {
						Name string `json:"name" vld:"equ1"`
					} `json:"field1" vld:"min1"`
				}{},
				tagType: VLD,
			},
			expected: []Validation{
				{Key: "field1", Type: Array, Requirement: "min1", InnerValidation: []Validation{
					{Key: "name", Type: String, Requirement: "equ1"},
				}},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with invalid struct validation",
			args: args{
				input: &struct {
					Field1 []struct {
						Name string `json:"name" vld:"equ"`
					} `json:"field1" vld:"min1"`
				}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: true,
		},
		{
			name: "Valid struct with invalid inner group",
			args: args{
				input: &struct {
					Field1 []struct {
						Name string `json:"name" vld:"equ1, gr"`
					} `json:"field1" vld:"min1"`
				}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: true,
		},
		{
			name: "Valid struct with invalid inner struct validation",
			args: args{
				input: &struct {
					Field1 struct {
						Name string `json:"name" vld:"equ"`
					} `json:"field1" vld:"-"`
				}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: true,
		},
		{
			name: "Empty struct",
			args: args{
				input:   &struct{}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: false,
		},
		{
			name: "Invalid struct with invalid group",
			args: args{
				input: &struct {
					Field1 string `vld:"equ1, gr"`
				}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: true,
		},
		{
			name: "Invalid struct type",
			args: args{
				input:   struct{}{},
				tagType: VLD,
			},
			expected:      []Validation{},
			expectedError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			validations, err := GetValidationsFromStruct(test.args.input, test.args.tagType)
			if test.expectedError {
				assert.Error(t, err, "Expected error when getting validations")
			} else {
				assert.NoError(t, err, "Expected no error when getting validations")
				assert.Equal(t, test.expected, validations, "Expected validations to match")
			}
		})
	}
}

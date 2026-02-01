package validator

import (
	"testing"

	"github.com/siherrmann/validator/model"
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
		expected      []model.Validation
		expectedError bool
	}{
		{
			name: "Valid struct with conditions and groups",
			args: args{
				input: &struct {
					Field1 string `vld:"equ1, gr1min1"`
					Field2 int    `vld:"min2, gr1min1"`
				}{},
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "Field1", Type: model.String, Requirement: "equ1", Groups: []*model.Group{{Name: "gr1", ConditionType: "min", ConditionValue: "1"}}},
				{Key: "Field2", Type: model.Int, Requirement: "min2", Groups: []*model.Group{{Name: "gr1", ConditionType: "min", ConditionValue: "1"}}},
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
			expected: []model.Validation{
				{Key: "Field1", Type: model.String, Requirement: "equ1"},
				{Key: "Field2", Type: model.Int, Requirement: "min2"},
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
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "field1", Type: model.String, Requirement: "equ1"},
				{Key: "field2", Type: model.Int, Requirement: "min2"},
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
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "field1", Type: model.Struct, Requirement: "-", InnerValidation: []model.Validation{
					{Key: "name", Type: model.String, Requirement: "equ1"},
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
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "field1", Type: model.Array, Requirement: "min1", InnerValidation: []model.Validation{
					{Key: "name", Type: model.String, Requirement: "equ1"},
				}},
			},
			expectedError: false,
		},
		{
			name: "Valid struct with invalid inner group",
			args: args{
				input: &struct {
					Field1 []struct {
						Name string `json:"name" vld:"equ1, gr"`
					} `json:"field1" vld:"min1"`
				}{},
				tagType: model.VLD,
			},
			expected:      []model.Validation{},
			expectedError: true,
		},
		{
			name: "Valid struct with valid ignored field",
			args: args{
				input: &struct {
					Name   string `json:"name"`
					Field1 string `json:"field1" vld:"-"`
				}{},
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "field1", Type: model.String, Requirement: "-"},
			},
			expectedError: false,
		},
		{
			name: "Empty struct",
			args: args{
				input:   &struct{}{},
				tagType: model.VLD,
			},
			expected:      []model.Validation{},
			expectedError: false,
		},
		{
			name: "Invalid struct with invalid group",
			args: args{
				input: &struct {
					Field1 string `vld:"equ1, gr"`
				}{},
				tagType: model.VLD,
			},
			expected:      []model.Validation{},
			expectedError: true,
		},
		{
			name: "Invalid struct type",
			args: args{
				input:   struct{}{},
				tagType: model.VLD,
			},
			expected:      []model.Validation{},
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

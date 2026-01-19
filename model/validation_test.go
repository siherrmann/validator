package model

import (
	"encoding/json"
	"testing"

	"github.com/siherrmann/validator/helper"
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

func TestValidationJSONMarshaling(t *testing.T) {
	t.Run("Marshal and unmarshal Validation with nil Groups", func(t *testing.T) {
		v := Validation{
			Key:         "TestField",
			Type:        String,
			Requirement: "min5",
			Groups:      nil,
			Default:     "",
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(v)
		assert.NoError(t, err, "Expected no error marshaling")

		// Unmarshal back
		var result Validation
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err, "Expected no error unmarshaling")
		assert.Equal(t, v, result, "Expected validation to match after round-trip")
	})

	t.Run("Marshal and unmarshal Validation with Groups", func(t *testing.T) {
		v := Validation{
			Key:         "TestField",
			Type:        String,
			Requirement: "min5",
			Groups:      []*Group{{Name: "gr1", ConditionType: "min", ConditionValue: "1"}},
			Default:     "",
		}

		// Marshal to JSON
		jsonBytes, err := json.Marshal(v)
		assert.NoError(t, err, "Expected no error marshaling")

		// Unmarshal back
		var result Validation
		err = json.Unmarshal(jsonBytes, &result)
		assert.NoError(t, err, "Expected no error unmarshaling")
		assert.Equal(t, v, result, "Expected validation to match after round-trip")
	})

	t.Run("Unmarshal Validation from map[string]any", func(t *testing.T) {
		jsonMap := map[string]any{
			"Key":         "TestField",
			"Type":        String,
			"Requirement": "min5",
			"Groups":      nil,
			"Default":     "",
		}

		var result Validation
		err := helper.MapJsonMapToStruct(jsonMap, &result)
		assert.NoError(t, err, "Expected no error mapping to struct")
		assert.Equal(t, "TestField", result.Key)
		assert.Equal(t, String, result.Type)
		assert.Equal(t, "min5", result.Requirement)
		assert.Nil(t, result.Groups)
	})

	t.Run("Unmarshal Validation with Groups from map[string]any", func(t *testing.T) {
		jsonMap := map[string]any{
			"Key":         "TestField",
			"Type":        String,
			"Requirement": "min5",
			"Groups": []any{
				map[string]any{"Name": "gr1", "ConditionType": "min", "ConditionValue": "1"},
			},
			"Default": "",
		}

		var result Validation
		err := helper.MapJsonMapToStruct(jsonMap, &result)
		assert.NoError(t, err, "Expected no error mapping to struct")
		assert.Equal(t, "TestField", result.Key)
		assert.Equal(t, String, result.Type)
		assert.Equal(t, "min5", result.Requirement)
		assert.NotNil(t, result.Groups)
		assert.Len(t, result.Groups, 1)
		assert.Equal(t, "gr1", result.Groups[0].Name)
	})

	t.Run("Unmarshal Validation with empty Groups from map[string]any", func(t *testing.T) {
		jsonMap := map[string]any{
			"Key":         "TestField",
			"Type":        String,
			"Requirement": "min5",
			"Groups":      []any{},
			"Default":     "",
		}

		var result Validation
		err := helper.MapJsonMapToStruct(jsonMap, &result)
		assert.NoError(t, err, "Expected no error mapping to struct with empty Groups")
		assert.Equal(t, "TestField", result.Key)
		assert.NotNil(t, result.Groups)
		assert.Len(t, result.Groups, 0, "Expected Groups to be empty slice")
	})

	t.Run("Unmarshal Validation with InnerValidation", func(t *testing.T) {
		jsonMap := map[string]any{
			"Key":         "TestField",
			"Type":        Struct,
			"Requirement": "-",
			"Groups":      nil,
			"Default":     "",
			"InnerValidation": []any{
				map[string]any{
					"Key":         "InnerField",
					"Type":        String,
					"Requirement": "min3",
					"Groups":      nil,
					"Default":     "",
				},
			},
		}

		var result Validation
		err := helper.MapJsonMapToStruct(jsonMap, &result)
		assert.NoError(t, err, "Expected no error mapping to struct with InnerValidation")
		assert.Equal(t, "TestField", result.Key)
		assert.Equal(t, Struct, result.Type)
		assert.NotNil(t, result.InnerValidation)
		assert.Len(t, result.InnerValidation, 1)
		assert.Equal(t, "InnerField", result.InnerValidation[0].Key)
		assert.Equal(t, String, result.InnerValidation[0].Type)
		assert.Equal(t, "min3", result.InnerValidation[0].Requirement)
	})
}

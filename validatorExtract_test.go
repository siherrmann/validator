package validator

import (
	"reflect"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "Valid struct with omitempty and inline tags",
			args: args{
				input: &struct {
					Inner struct {
						Kind string `json:"kind" vld:"equ1"`
					} `json:",inline" vld:"-"`
					Field1 string `json:"field1,omitempty" vld:"equ1"`
				}{},
				tagType: model.VLD,
			},
			expected: []model.Validation{
				{Key: "kind", Type: model.String, Requirement: "equ1"},
				{Key: "field1", Type: model.String, Requirement: "equ1", OmitEmpty: true},
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

				var clearASTs func(vs []model.Validation)
				clearASTs = func(vs []model.Validation) {
					for i := range vs {
						vs[i].RequirementAST = nil
						vs[i].NoFieldUsage = false
						clearASTs(vs[i].InnerValidation)
					}
				}
				clearASTs(validations)

				assert.Equal(t, test.expected, validations, "Expected validations to match")
			}
		})
	}

	t.Run("Field usage of rel and references", func(t *testing.T) {
		type Inner struct {
			Value int `json:"value" vld:"min$../min"`
		}
		type Plain struct {
			Value int `json:"value" vld:"min1"`
		}
		type Request struct {
			Min   int    `json:"min" vld:"min1"`
			Name  string `json:"name" vld:"rel(min equ1) => min1"`
			Inner Inner  `json:"inner" vld:"-"`
			Plain Plain  `json:"plain" vld:"-"`
		}

		validations, err := GetValidationsFromStruct(&Request{}, model.VLD)
		require.NoError(t, err)
		usage := map[string]bool{}
		for _, v := range validations {
			usage[v.Key] = v.UsesFields()
		}
		assert.Equal(t, map[string]bool{"min": false, "name": true, "inner": true, "plain": false}, usage)
	})

	t.Run("Embedded structs", func(t *testing.T) {
		type Base struct {
			ID string `json:"id" vld:"min1"`
		}
		type Request struct {
			Base
			Name string `json:"name" vld:"min1"`
		}

		validations, err := GetValidationsFromStruct(&Request{}, model.VLD)
		require.NoError(t, err)
		require.Len(t, validations, 2)
		assert.Equal(t, "id", validations[0].Key, "Expected fields of an embedded struct on the level of the struct")
		assert.Equal(t, "name", validations[1].Key)
	})

	t.Run("Recursive types", func(t *testing.T) {
		validations, err := GetValidationsFromStruct(&extractNode{}, model.VLD)
		require.NoError(t, err)
		require.Len(t, validations, 3)

		children := validations[1]
		assert.Equal(t, "children", children.Key)
		assert.Nil(t, children.InnerValidation)
		require.NotNil(t, children.InnerValidationRef)
		assert.Same(t, &validations[0], &children.Inner()[0], "Expected a reference to the validations of the struct")

		assert.True(t, validations[0].UsesFields())
		assert.True(t, children.UsesFields(), "Expected field usage through the reference")
		assert.True(t, validations[2].UsesFields(), "Expected field usage through the cycle")
	})
}

func TestGetValidationFromStructField(t *testing.T) {
	type Inner struct {
		Name string `json:"name" vld:"min1"`
	}
	type Request struct {
		Name     string  `json:"name,omitempty" vld:"min3 def'abc', gr1min1"`
		Short    string  `vld:"unq"`
		Ignored  string  `json:"ignored"`
		Inner    Inner   `json:"inner" vld:"-"`
		List     []Inner `json:"list" vld:"min1"`
		Relation string  `json:"relation" vld:"rel(name equabc) => min1"`
		Invalid  string  `vld:"invalid"`
	}

	extract := func(name string) (*model.Validation, error) {
		value := reflect.ValueOf(&Request{}).Elem()
		structField, ok := value.Type().FieldByName(name)
		require.True(t, ok)
		return GetValidationFromStructField(model.VLD, value.FieldByName(name), structField)
	}

	t.Run("Requirement, default, omitempty and groups", func(t *testing.T) {
		validation, err := extract("Name")
		require.NoError(t, err)
		require.NotNil(t, validation)
		assert.Equal(t, "name", validation.Key)
		assert.Equal(t, model.String, validation.Type)
		assert.Equal(t, "min3 def'abc'", validation.Requirement)
		assert.Equal(t, "abc", validation.Default)
		assert.True(t, validation.OmitEmpty)
		require.Len(t, validation.Groups, 1)
		assert.Equal(t, "gr1", validation.Groups[0].Name)
		assert.NotNil(t, validation.RequirementAST)
		assert.False(t, validation.UsesFields())
	})

	t.Run("Short requirement", func(t *testing.T) {
		validation, err := extract("Short")
		require.NoError(t, err)
		require.NotNil(t, validation)
		assert.Equal(t, "unq", validation.Requirement)
	})

	t.Run("Field without tag", func(t *testing.T) {
		validation, err := extract("Ignored")
		assert.NoError(t, err)
		assert.Nil(t, validation)
	})

	t.Run("Inner validations of struct and list", func(t *testing.T) {
		validation, err := extract("Inner")
		require.NoError(t, err)
		require.Len(t, validation.InnerValidation, 1)
		assert.Equal(t, "name", validation.InnerValidation[0].Key)

		validation, err = extract("List")
		require.NoError(t, err)
		assert.Equal(t, model.Array, validation.Type)
		require.Len(t, validation.InnerValidation, 1)
	})

	t.Run("Field usage", func(t *testing.T) {
		validation, err := extract("Relation")
		require.NoError(t, err)
		assert.True(t, validation.UsesFields())
	})

	t.Run("Invalid requirement", func(t *testing.T) {
		_, err := extract("Invalid")
		assert.Error(t, err)
	})
}

type extractNode struct {
	Value    int             `json:"value" vld:"rel(children min1) => min1"`
	Children []extractNode   `json:"children" vld:"-"`
	Parent   *extractWrapper `json:"parent" vld:"-"`
}

type extractWrapper struct {
	Node *extractNode `json:"node" vld:"-"`
}

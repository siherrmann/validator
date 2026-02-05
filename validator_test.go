package validator

import (
	"fmt"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {
	t.Run("Valid struct", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" vld:"equapple"`
		}
		testStruct := &TestStruct{
			Fruit: "apple",
		}
		r := NewValidator()
		err := r.Validate(testStruct)
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Invalid struct", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" vld:"equapple"`
		}
		testStruct := &TestStruct{
			Fruit: "banana",
		}
		r := NewValidator()
		err := r.Validate(testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field fruit invalid: value not equal condition apple", "Expected error to contain 'field fruit invalid: value not equal condition apple'")
	})

	t.Run("Invalid struct pointer", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" vld:"equapple"`
		}
		testStruct := TestStruct{
			Fruit: "banana",
		}
		r := NewValidator()
		err := r.Validate(testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain 'value has to be of kind pointer'")
	})

	t.Run("Valid struct with custom tag type", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" update:"equapple"`
		}
		testStruct := &TestStruct{
			Fruit: "apple",
		}
		r := NewValidator()
		err := r.Validate(testStruct, "update")
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Invalid struct validation", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" update:"equapple, gp1min1"`
		}
		testStruct := &TestStruct{
			Fruit: "apple",
		}
		r := NewValidator()
		err := r.Validate(testStruct, "update")
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "invalid group name: gp1", "Expected error to contain 'invalid group name: gp1'")
	})
}

func TestValidateAndUpdate(t *testing.T) {
	r := NewValidator()

	t.Run("Valid struct", func(t *testing.T) {
		testStruct := &struct {
			Fruit string `json:"fruit" vld:"equapple"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": "apple"}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruit, "Expected output to match input")
	})

	t.Run("Invalid struct", func(t *testing.T) {
		testStruct := &struct {
			Fruit string `json:"fruit" vld:"equapple"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": "banana"}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error with invalid input")
		assert.Equal(t, "", testStruct.Fruit, "Expected output to be unchanged")
	})

	t.Run("Valid struct with inner struct", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruit" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{"name": "apple"}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruit.Name, "Expected output to match input")
	})

	t.Run("Invalid struct with inner struct", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruit" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{"name": "banana"}}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error but got none")
		assert.Equal(t, "", testStruct.Fruit.Name, "Expected output to be unchanged")
	})

	t.Run("Valid struct with array of structs", func(t *testing.T) {
		testStruct := &struct {
			Fruits []struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruits" vld:"min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruits": []any{map[string]any{"name": "apple"}}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruits[0].Name, "Expected output to match input")
	})

	t.Run("Invalid struct with array of structs", func(t *testing.T) {
		testStruct := &struct {
			Fruits []struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruits" vld:"min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruits": []any{map[string]any{"name": "banana"}}}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error but got none")
		assert.Empty(t, testStruct.Fruits, "Expected output to be unchanged")
	})

	t.Run("Valid struct with complex condition", func(t *testing.T) {
		testStruct := &struct {
			Fruits []struct {
				Name string `json:"name" vld:"(equapple || min1) && neqbanana"`
			} `json:"fruits" vld:"min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruits": []any{map[string]any{"name": "apple"}}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruits[0].Name, "Expected output to match input")

		err = r.ValidateAndUpdate(map[string]any{"fruits": []any{map[string]any{"name": "a"}}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "a", testStruct.Fruits[0].Name, "Expected output to match input")

		err = r.ValidateAndUpdate(map[string]any{"fruits": []any{map[string]any{"name": "banana"}}}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error but got none")
		assert.Equal(t, "a", testStruct.Fruits[0].Name, "Expected output to match input")
	})

	t.Run("Invalid struct pointer", func(t *testing.T) {
		type TestStructInvalid struct {
			Fruit string `json:"fruit" vld:"equapple, gp1min1"`
		}
		testStruct := TestStructInvalid{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": "apple"}, testStruct)
		require.Error(t, err, "Expected an error for invalid validation")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain 'value has to be of kind pointer'")
	})

	t.Run("Invalid struct tag", func(t *testing.T) {
		type TestStructInvalid struct {
			Fruit string `json:"fruit" vld:"equapple, gp1min1"`
		}
		testStruct := &TestStructInvalid{}
		r := NewValidator()
		err := r.ValidateAndUpdate(map[string]any{"fruit": "apple"}, testStruct)
		require.Error(t, err, "Expected an error for invalid validation")
		assert.Contains(t, err.Error(), "invalid group name: gp1", "Expected error to contain 'invalid group name: gp1'")
	})

	t.Run("Tag filtering in arrays of structs with custom tag", func(t *testing.T) {
		testStruct := &struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name" upd:"-"`
				Age  int    `json:"age" upd:"min18"`
			} `json:"items" upd:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"items": []any{map[string]any{"id": "123", "name": "test", "age": 19}}}, testStruct, "upd")
		assert.NoError(t, err, "Expected no error")
		assert.Len(t, testStruct.Items, 1, "Items should be updated")
		assert.Equal(t, "", testStruct.Items[0].ID, "ID should not be updated (no upd tag)")
		assert.Equal(t, "test", testStruct.Items[0].Name, "Name should not be updated (has ignore upd tag)")
		assert.Equal(t, 19, testStruct.Items[0].Age, "Age should be updated (has full upd tag)")
	})
}

func TestValidateAndUpdateWithValidation(t *testing.T) {
	type args struct {
		jsonMap     map[string]any
		validations []model.Validation
	}
	tests := []struct {
		name     string
		args     args
		expected any
		wantErr  bool
	}{
		{
			name: "Valid validation",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid validation",
			args: args{
				jsonMap: map[string]any{
					"fruit": "banana",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
				},
			},
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewValidator()

			mapOut := &map[string]any{}
			err := r.ValidateAndUpdateWithValidation(test.args.jsonMap, mapOut, test.args.validations)
			if test.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				if test.expected != nil {
					assert.Equal(t, &test.expected, mapOut, "Expected output to match input")
				} else {
					assert.Equal(t, &test.args.jsonMap, mapOut, "Expected output to match input")
				}
			}
		})
	}
}

func TestValidateWithValidation(t *testing.T) {
	type args struct {
		jsonMap     map[string]any
		validations []model.Validation
	}
	tests := []struct {
		name     string
		args     args
		expected any
		wantErr  bool
	}{
		{
			name: "Valid validation string",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid validation string",
			args: args{
				jsonMap: map[string]any{
					"fruit": "banana",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid validation string with group",
			args: args{
				jsonMap: map[string]any{
					"fruit": "banana",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
						Groups: []*model.Group{{
							Name:           "gr1",
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "1",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid duplicate key",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equbanana",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid validation with group",
			args: args{
				jsonMap: map[string]any{
					"fruit":  "apple",
					"fruit2": "banana",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
						Groups: []*model.Group{{
							Name:           "gr1",
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "2",
						}},
					},
					{
						Key:         "fruit2",
						Type:        model.String,
						Requirement: "equbanana",
						Groups: []*model.Group{{
							Name:           "gr1",
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "2",
						}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid with missing key",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:         "missing",
						Type:        model.String,
						Requirement: "-",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid with missing key",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:         "missing",
						Type:        model.String,
						Requirement: "equbanana",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid with missing key with group",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:         "missing",
						Type:        model.String,
						Requirement: "equbanana",
						Groups: []*model.Group{{
							Name:           "gr1",
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "2",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid nested struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
					"nested": map[string]any{
						"fruit": "banana",
					},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:  "nested",
						Type: model.Struct,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equbanana",
						}},
					},
				},
			},
			expected: map[string]any{"fruit": "apple", "nested": map[string]any{"fruit": "banana"}},
			wantErr:  false,
		},
		{
			name: "Invalid nested struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
					"nested": map[string]any{
						"fruit": "apple",
					},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:  "nested",
						Type: model.Struct,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equbanana",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid nested struct",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
					"nested": map[string]any{
						"fruit": 1,
					},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
					},
					{
						Key:  "nested",
						Type: model.Struct,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equbanana",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid array validation",
			args: args{
				jsonMap: map[string]any{
					"fruits": []string{"apple", "banana"},
				},
				validations: []model.Validation{
					{
						Key:         "fruits",
						Type:        model.Array,
						Requirement: "frmapple,banana",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid array validation from url values with type string",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Array,
						Requirement: "frmapple,banana",
					},
				},
			},
			expected: map[string]any{"fruit": []string{"apple"}},
			wantErr:  false,
		},
		{
			name: "Valid array of struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruits": []any{
						map[string]any{
							"fruit": "apple",
						},
					},
				},
				validations: []model.Validation{
					{
						Key:  "fruits",
						Type: model.Array,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equapple",
						}},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid array of struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruits": []any{
						map[string]any{
							"fruit": "banana",
						},
					},
				},
				validations: []model.Validation{
					{
						Key:  "fruits",
						Type: model.Array,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equapple",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid array type of struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruits": []map[string]any{{
						"fruit": "apple",
					}},
				},
				validations: []model.Validation{
					{
						Key:  "fruits",
						Type: model.Array,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.String,
							Requirement: "equapple",
						}},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid array element type of struct validation",
			args: args{
				jsonMap: map[string]any{
					"fruits": []any{
						map[string]int{
							"fruit": 1,
						},
					},
				},
				validations: []model.Validation{
					{
						Key:  "fruits",
						Type: model.Array,
						InnerValidation: []model.Validation{{
							Key:         "fruit",
							Type:        model.Int,
							Requirement: "equ2",
						}},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewValidator()
			out, err := r.ValidateWithValidation(test.args.jsonMap, test.args.validations)
			if test.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				if test.expected != nil {
					assert.Equal(t, test.expected, out, "Expected output to match input")
				} else {
					assert.Equal(t, test.args.jsonMap, out, "Expected output to match input")
				}
			}
		})
	}
}

func TestValidateValueWithParser(t *testing.T) {
	type args struct {
		input      any
		validation *model.Validation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid value EQUAL",
			args: args{
				input: "apple",
				validation: &model.Validation{
					Key:         "fruit",
					Type:        model.String,
					Requirement: "equapple",
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value EQUAL",
			args: args{
				input: "banana",
				validation: &model.Validation{
					Key:         "fruit",
					Type:        model.String,
					Requirement: "equapple",
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid validation EQUAL",
			args: args{
				input: "banana",
				validation: &model.Validation{
					Key:         "fruit",
					Type:        model.String,
					Requirement: "(equapple &| equbanana)",
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := NewValidator()
			err := r.ValidateValueWithParser(test.args.input, test.args.validation)
			if test.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

func TestRunValidatorsOnConditionGroup(t *testing.T) {
	type args struct {
		input    any
		astValue *model.AstValue
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid condition group with empty condition",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:     model.EMPTY,
							Operator: model.AND,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid condition group with NONE condition type",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:          model.CONDITION,
							ConditionType: model.NONE,
							Operator:      model.AND,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid condition group with invalid condition type",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:          model.CONDITION,
							ConditionType: model.ConditionType("Uff"),
							Operator:      model.AND,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid condition group with inner group",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type: model.GROUP,
							ConditionGroup: []*model.AstValue{
								{
									Type:           model.CONDITION,
									ConditionType:  model.EQUAL,
									ConditionValue: "apple",
									Operator:       model.AND,
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid with OR operator",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "banana",
							Operator:       model.OR,
						},
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "cherry",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.EQUAL,
							ConditionValue: "banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_EQUAL,
							ConditionValue: "banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_EQUAL",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_EQUAL,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value MIN_VALUE",
			args: args{
				input: 5,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "3",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value MIN_VALUE",
			args: args{
				input: 3,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MIN_VALUE,
							ConditionValue: "5",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value MAX_VALUE",
			args: args{
				input: 3,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MAX_VALUE,
							ConditionValue: "5",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value MAX_VALUE",
			args: args{
				input: 5,
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.MAX_VALUE,
							ConditionValue: "3",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value CONTAINS",
			args: args{
				input: "pineapple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value CONTAINS",
			args: args{
				input: "banana",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_CONTAINS",
			args: args{
				input: "banana",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_CONTAINS",
			args: args{
				input: "pineapple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_CONTAINS,
							ConditionValue: "apple",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value FROM",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value FROM",
			args: args{
				input: "cherry",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value NOT_FROM",
			args: args{
				input: "cherry",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value NOT_FROM",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.NOT_FROM,
							ConditionValue: "apple,banana",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value REGX",
			args: args{
				input: "abc123",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.REGX,
							ConditionValue: "^[a-z]+[0-9]+$",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value REGX",
			args: args{
				input: "abc@123",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.REGX,
							ConditionValue: "^[a-z]+[0-9]+$",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid value FUNC",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FUNC,
							ConditionValue: "testFunc",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid value FUNC",
			args: args{
				input: "banana",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FUNC,
							ConditionValue: "testFunc",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid function name FUNC",
			args: args{
				input: "apple",
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.FUNC,
							ConditionValue: "unknownFunc",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			testFunc := func(input any, astValue *model.AstValue) error {
				if input != "apple" {
					return fmt.Errorf("testFunc validation failed")
				}
				return nil
			}

			r := NewValidator()
			r.AddValidationFunc(testFunc, "testFunc")
			require.NotNil(t, r, "Validator should not be nil")
			require.Contains(t, r.ValidationFuncs, "testFunc", "Validation function 'testFunc' should be registered")

			err := r.RunValidatorsOnConditionGroup(test.args.input, test.args.astValue)
			if test.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}

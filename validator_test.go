package validator

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
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

	t.Run("Valid struct with istuuid", func(t *testing.T) {
		type TestStruct struct {
			ID string `json:"id" vld:"istuuid"`
		}
		testStruct := &TestStruct{
			ID: "123e4567-e89b-12d3-a456-426614174000",
		}
		r := NewValidator()
		err := r.Validate(testStruct)
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Invalid struct with istuuid", func(t *testing.T) {
		type TestStruct struct {
			ID string `json:"id" vld:"istuuid"`
		}
		testStruct := &TestStruct{
			ID: "invalid-uuid",
		}
		r := NewValidator()
		err := r.Validate(testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field id invalid: value is not a valid uuid", "Expected error to contain 'value is not a valid uuid'")
	})

	t.Run("Valid struct with istpassword", func(t *testing.T) {
		type TestStruct struct {
			Password string `json:"password" vld:"istpassword"`
		}
		testStruct := &TestStruct{
			Password: "StrongPassword1!",
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

	t.Run("Valid struct with omitempty json tag", func(t *testing.T) {
		type TestStruct struct {
			Name string `json:"name,omitempty" vld:"min1"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{Name: "apple"})
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Invalid struct with yaml tag only", func(t *testing.T) {
		type TestStruct struct {
			Name string `yaml:"name" vld:"min3"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{Name: "ap"})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field name invalid", "Expected error to use yaml key")
	})

	t.Run("Valid struct with nil interface value", func(t *testing.T) {
		type TestStruct struct {
			Value any `json:"value" vld:"max0"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{})
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Valid struct with uuid and istuuid", func(t *testing.T) {
		type TestStruct struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional" vld:"-"`
		}
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		r := NewValidator()
		err := r.Validate(&TestStruct{ID: id, Optional: &id})
		assert.NoError(t, err, "Expected no error but got one")

		err = r.Validate(&TestStruct{ID: id})
		assert.NoError(t, err, "Expected no error for nil uuid pointer but got one")
	})

	t.Run("Valid struct with uuid and min max", func(t *testing.T) {
		type TestStruct struct {
			ID uuid.UUID `json:"id" vld:"min36 && max36"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{ID: uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")})
		assert.NoError(t, err, "Expected no error but got one")
	})

	t.Run("Invalid struct with nil uuid", func(t *testing.T) {
		type TestStruct struct {
			ID uuid.UUID `json:"id" vld:"neq00000000-0000-0000-0000-000000000000"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field id invalid", "Expected error for nil uuid")
	})

	t.Run("Valid struct with net.IP", func(t *testing.T) {
		type TestStruct struct {
			IP net.IP `json:"ip" vld:"istipv4"`
		}
		r := NewValidator()
		err := r.Validate(&TestStruct{IP: net.ParseIP("192.168.0.1")})
		assert.NoError(t, err, "Expected no error but got one")
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
	t.Run("Valid struct with missing field having def tag", func(t *testing.T) {
		testStruct := &struct {
			Fruit string `json:"fruit" vld:"equapple defunasigned"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "unasigned", testStruct.Fruit, "Expected output to match default")
	})

	t.Run("Valid struct with invalid field having def tag", func(t *testing.T) {
		testStruct := &struct {
			Fruit string `json:"fruit" vld:"equapple defunasigned"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": "banana"}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "unasigned", testStruct.Fruit, "Expected output to match default")
	})

	t.Run("Valid struct with valid field having def tag", func(t *testing.T) {
		testStruct := &struct {
			Fruit string `json:"fruit" vld:"equapple defunasigned"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": "apple"}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruit, "Expected output to match input")
	})

	t.Run("Valid struct with def tag for int", func(t *testing.T) {
		testStruct := &struct {
			Count int `json:"count" vld:"min1 def10"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, 10, testStruct.Count, "Expected output to match default")

		testStruct.Count = 0
		err = r.ValidateAndUpdate(map[string]any{"count": 0}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, 10, testStruct.Count, "Expected output to match default")
	})

	t.Run("Valid struct with yaml tag only", func(t *testing.T) {
		testStruct := &struct {
			Name string `yaml:"name" vld:"min3"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"name": "apple"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Name, "Expected name to be updated")
	})

	t.Run("Valid struct with null values", func(t *testing.T) {
		type Inner struct {
			Name string `json:"name" vld:"min3"`
		}
		testStruct := &struct {
			Name     string   `json:"name" vld:"max0 || min3"`
			Tags     []string `json:"tags" vld:"max0 || min1"`
			Items    []Inner  `json:"items" vld:"-"`
			Inner    Inner    `json:"inner" vld:"max0"`
			Optional *Inner   `json:"optional" vld:"-"`
		}{Name: "old"}
		input := map[string]any{"name": nil, "tags": nil, "items": nil, "inner": nil, "optional": nil}
		err := r.ValidateAndUpdate(input, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "", testStruct.Name, "Expected name to be set to zero value")
		assert.Empty(t, testStruct.Tags, "Expected tags to be empty")
		assert.Empty(t, testStruct.Items, "Expected items to be empty")
		assert.Nil(t, testStruct.Optional, "Expected optional to be nil")
	})

	t.Run("Invalid struct with null values and min", func(t *testing.T) {
		type TestStruct struct {
			Name string   `json:"name" vld:"min1"`
			Tags []string `json:"tags" vld:"min1"`
		}
		for _, key := range []string{"name", "tags"} {
			input := map[string]any{"name": "apple", "tags": []any{"a"}}
			input[key] = nil
			err := r.ValidateAndUpdate(input, &TestStruct{})
			assert.Error(t, err, "Expected an error but got none")
			assert.Contains(t, err.Error(), "field "+key+" invalid", "Expected error for null field")
		}
	})

	t.Run("Valid struct with null value and def tag", func(t *testing.T) {
		testStruct := &struct {
			Name string `json:"name" vld:"min1 defapple"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"name": nil}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Name, "Expected default to be used")
	})

	t.Run("Valid struct with uuid strings", func(t *testing.T) {
		testStruct := &struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional,omitempty" vld:"-"`
		}{}
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		err := r.ValidateAndUpdate(map[string]any{"id": id.String(), "optional": id.String()}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, id, testStruct.ID, "Expected id to be updated")
		require.NotNil(t, testStruct.Optional, "Expected optional to be set")
		assert.Equal(t, id, *testStruct.Optional, "Expected optional to be updated")
	})

	t.Run("Valid struct with uuid value", func(t *testing.T) {
		testStruct := &struct {
			ID uuid.UUID `json:"id" vld:"istuuid"`
		}{}
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		err := r.ValidateAndUpdate(map[string]any{"id": id}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, id, testStruct.ID, "Expected id to be updated")
	})

	t.Run("Valid struct with missing and empty optional uuid", func(t *testing.T) {
		type TestStruct struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional,omitempty" vld:"-"`
		}
		id := "550e8400-e29b-41d4-a716-446655440000"
		for _, input := range []map[string]any{{"id": id}, {"id": id, "optional": ""}} {
			testStruct := &TestStruct{}
			err := r.ValidateAndUpdate(input, testStruct)
			assert.NoError(t, err, "Expected no error but got one")
			assert.Nil(t, testStruct.Optional, "Expected optional to stay nil")
		}
	})

	t.Run("Valid struct with empty uuid pointer without omitempty", func(t *testing.T) {
		testStruct := &struct {
			ID *uuid.UUID `json:"id" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"id": ""}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Nil(t, testStruct.ID, "Expected id to stay nil")
	})

	t.Run("Invalid struct with uuid and istuuid", func(t *testing.T) {
		testStruct := &struct {
			ID uuid.UUID `json:"id" vld:"istuuid"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"id": "not-a-uuid"}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value is not a valid uuid", "Expected uuid validation error")
	})

	t.Run("Invalid struct with uuid without validation", func(t *testing.T) {
		testStruct := &struct {
			ID uuid.UUID `json:"id" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"id": "not-a-uuid"}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "error parsing string to uuid.UUID", "Expected uuid parse error")
	})

	t.Run("Valid struct with net.IP", func(t *testing.T) {
		type TestStruct struct {
			IP net.IP `json:"ip" vld:"istip"`
		}
		for _, input := range []string{"192.168.0.1", "2001:db8::1"} {
			testStruct := &TestStruct{}
			err := r.ValidateAndUpdate(map[string]any{"ip": input}, testStruct)
			assert.NoError(t, err, "Expected no error but got one")
			assert.Equal(t, input, testStruct.IP.String(), "Expected ip to be updated")
		}
	})

	t.Run("Invalid struct with net.IP", func(t *testing.T) {
		testStruct := &struct {
			IP net.IP `json:"ip" vld:"istip"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"ip": "not-an-ip"}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value is not a valid ip address", "Expected ip validation error")
	})

	t.Run("Valid struct with time.Time", func(t *testing.T) {
		testStruct := &struct {
			CreatedAt time.Time `json:"created_at" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"created_at": "2025-01-02T00:00:00Z"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), testStruct.CreatedAt.UTC(), "Expected created_at to be updated")

		err = r.ValidateAndUpdate(map[string]any{"created_at": "1735776000"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, int64(1735776000), testStruct.CreatedAt.Unix(), "Expected created_at to be updated from unix")
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
			name: "Valid omitempty logic with empty string",
			args: args{
				jsonMap: map[string]any{
					"fruit": "",
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.String,
						Requirement: "equapple",
						OmitEmpty:   true,
					},
				},
			},
			expected: map[string]any{},
			wantErr:  false,
		},
		{
			name: "Valid defFalse logic for missing boolean",
			args: args{
				jsonMap: map[string]any{},
				validations: []model.Validation{
					{
						Key:         "isActive",
						Type:        model.Bool,
						Requirement: "deffalse",
						Default:     "false",
					},
				},
			},
			expected: map[string]any{"isActive": false},
			wantErr:  false,
		},
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
		{
			name: "Valid null value for string with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.String,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for string with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.String,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for int with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Int,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for int with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Int,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for float with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Float,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for float with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Float,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for bool with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Bool,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for bool with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Bool,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for array with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Array,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for array with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Array,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for map with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Map,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for map with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Map,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for struct with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Struct,
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for struct with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        model.Struct,
						Requirement: "min1",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid null value for untyped with max0",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        "",
						Requirement: "max0",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid null value for untyped with min1",
			args: args{
				jsonMap: map[string]any{
					"value": nil,
				},
				validations: []model.Validation{
					{
						Key:         "value",
						Type:        "",
						Requirement: "min1",
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

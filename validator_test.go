package validator

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValidator(t *testing.T) {
	r := NewValidator()
	require.NotNil(t, r)
	assert.NotNil(t, r.ValidationFuncs)
	assert.Empty(t, r.ValidationFuncs)
}

func TestAddValidationFunc(t *testing.T) {
	type Request struct {
		Fruit string `json:"fruit" vld:"funisApple"`
	}
	isApple := func(input any, astValue *model.AstValue) error {
		if input != "apple" {
			return fmt.Errorf("value is not apple")
		}
		return nil
	}

	r := NewValidator()
	r.AddValidationFunc(isApple, "isApple")
	require.Contains(t, r.ValidationFuncs, "isApple")
	assert.NoError(t, r.Validate(&Request{Fruit: "apple"}))
	assert.ErrorContains(t, r.Validate(&Request{Fruit: "banana"}), "value is not apple")

	t.Run("Zero validator", func(t *testing.T) {
		zero := &Validator{}
		zero.AddValidationFunc(isApple, "isApple")
		assert.Contains(t, zero.ValidationFuncs, "isApple")
	})
}

func TestAddValidationFuncContext(t *testing.T) {
	type Request struct {
		Token string `json:"token" vld:"funsigned"`
	}

	r := NewValidator()
	r.AddValidationFuncContext(func(ctx context.Context, input any, astValue *model.AstValue) error {
		secret, ok := Get[string](ctx, "signing_key")
		if !ok {
			return fmt.Errorf("signing key not injected")
		}
		if input != "payload."+secret {
			return fmt.Errorf("invalid signature")
		}
		return nil
	}, "signed")

	ctx := Inject(context.Background(), "signing_key", "k1")
	assert.NoError(t, r.ValidateContext(ctx, &Request{Token: "payload.k1"}))
	assert.ErrorContains(t, r.ValidateContext(ctx, &Request{Token: "payload.k2"}), "invalid signature")
	assert.ErrorContains(t, r.Validate(&Request{Token: "payload.k1"}), "signing key not injected")

	t.Run("Functions with the same name replace each other", func(t *testing.T) {
		r.AddValidationFunc(func(input any, astValue *model.AstValue) error { return nil }, "signed")
		assert.NoError(t, r.Validate(&Request{Token: "anything"}))

		r.AddValidationFuncContext(func(ctx context.Context, input any, astValue *model.AstValue) error {
			return fmt.Errorf("context function")
		}, "signed")
		assert.ErrorContains(t, r.Validate(&Request{Token: "anything"}), "context function")
	})

	t.Run("Zero validator", func(t *testing.T) {
		zero := &Validator{}
		zero.AddValidationFuncContext(func(ctx context.Context, input any, astValue *model.AstValue) error { return nil }, "ok")
		assert.Contains(t, zero.ValidationFuncs, "ok")
	})
}

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
	})

	t.Run("Valid struct with nil uuid pointer", func(t *testing.T) {
		type TestStruct struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional" vld:"-"`
		}
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		r := NewValidator()
		err := r.Validate(&TestStruct{ID: id})
		assert.NoError(t, err, "Expected no error but got one")
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

	t.Run("Struct with custom validation function", func(t *testing.T) {
		type TestStruct struct {
			Fruit string `json:"fruit" vld:"funisApple"`
		}
		r := NewValidator()
		r.AddValidationFunc(func(input any, astValue *model.AstValue) error {
			if input != "apple" {
				return fmt.Errorf("value is not apple")
			}
			return nil
		}, "isApple")

		err := r.Validate(&TestStruct{Fruit: "apple"})
		assert.NoError(t, err, "Expected no error but got one")

		err = r.Validate(&TestStruct{Fruit: "banana"})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value is not apple")
	})

	type NestedInner struct {
		Name string `json:"name" vld:"min3"`
	}
	type NestedInline struct {
		Code string `json:"code" vld:"min3"`
	}
	type NestedOuter struct {
		Inner        NestedInner   `json:"inner" vld:"-"`
		InnerPointer *NestedInner  `json:"inner_pointer" vld:"-"`
		InnerList    []NestedInner `json:"inner_list" vld:"-"`
		NestedInline `json:",inline" vld:"-"`
		CreatedAt    time.Time `json:"created_at" vld:"-"`
		unexported   string
	}
	validNested := func() *NestedOuter {
		return &NestedOuter{
			Inner:        NestedInner{Name: "abc"},
			InnerPointer: &NestedInner{Name: "abc"},
			InnerList:    []NestedInner{{Name: "abc"}, {Name: "def"}},
			NestedInline: NestedInline{Code: "abc"},
			CreatedAt:    time.Now(),
			unexported:   "x",
		}
	}

	t.Run("Valid struct with nested structs", func(t *testing.T) {
		r := NewValidator()
		assert.NoError(t, r.Validate(validNested()), "Expected no error but got one")

		withNilPointer := validNested()
		withNilPointer.InnerPointer = nil
		withNilPointer.InnerList = nil
		assert.NoError(t, r.Validate(withNilPointer), "Expected no error for nil pointer and nil list")
	})

	t.Run("Invalid nested struct", func(t *testing.T) {
		in := validNested()
		in.Inner.Name = "a"
		err := NewValidator().Validate(in)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field inner invalid")
	})

	t.Run("Invalid nested struct pointer", func(t *testing.T) {
		in := validNested()
		in.InnerPointer.Name = "a"
		err := NewValidator().Validate(in)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field inner_pointer invalid")
	})

	t.Run("Invalid nested struct in list", func(t *testing.T) {
		in := validNested()
		in.InnerList[1].Name = "a"
		err := NewValidator().Validate(in)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field inner_list invalid")
	})

	t.Run("Invalid inline struct", func(t *testing.T) {
		in := validNested()
		in.Code = "a"
		err := NewValidator().Validate(in)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field code invalid")
	})

	t.Run("Embedded struct", func(t *testing.T) {
		type Base struct {
			ID string `json:"id" vld:"min3"`
		}
		type Request struct {
			*Base
			Name string `json:"name" vld:"min1"`
		}

		r := NewValidator()
		assert.NoError(t, r.Validate(&Request{Base: &Base{ID: "abc"}, Name: "n"}))
		assert.ErrorContains(t, r.Validate(&Request{Base: &Base{ID: "a"}, Name: "n"}), "field id invalid")
		assert.ErrorContains(t, r.Validate(&Request{Name: "n"}), "json id key not in map", "Expected fields of a nil embedded pointer to be missing")
	})

	t.Run("Recursive struct", func(t *testing.T) {
		type Node struct {
			Name     string `json:"name" vld:"min1"`
			Children []Node `json:"children" vld:"-"`
			Next     *Node  `json:"next" vld:"-"`
		}

		r := NewValidator()
		valid := &Node{Name: "a", Children: []Node{{Name: "b", Children: []Node{{Name: "c"}}}}, Next: &Node{Name: "d"}}
		assert.NoError(t, r.Validate(valid))

		invalid := &Node{Name: "a", Children: []Node{{Name: "b", Children: []Node{{Name: ""}}}}}
		assert.ErrorContains(t, r.Validate(invalid), "field children invalid")
		assert.ErrorContains(t, r.Validate(&Node{Name: "a", Next: &Node{}}), "field next invalid")
	})

	t.Run("Requirement of a list of structs", func(t *testing.T) {
		type Item struct {
			Name string `json:"name" vld:"min1"`
		}
		type Request struct {
			Items []Item   `json:"items" vld:"min1 && unqname"`
			Tags  []string `json:"tags" vld:"unq"`
		}

		r := NewValidator()
		assert.NoError(t, r.Validate(&Request{Items: []Item{{Name: "a"}, {Name: "b"}}, Tags: []string{"a"}}))

		err := r.Validate(&Request{Items: []Item{}, Tags: []string{"a"}})
		require.Error(t, err, "The requirement of a list of structs has to be checked")
		assert.Contains(t, err.Error(), "field items invalid")

		err = r.Validate(&Request{Items: []Item{{Name: "a"}, {Name: "a"}}})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "not unique")

		err = r.Validate(&Request{Items: []Item{{Name: "a"}}, Tags: []string{"a", "a"}})
		require.Error(t, err, "A tag with only `unq` has to be validated")
		assert.Contains(t, err.Error(), "field tags invalid")
	})
}

func TestValidateContext(t *testing.T) {
	type Request struct {
		Owner string `json:"owner" vld:"equ${user_id}"`
	}
	ctx := Inject(context.Background(), "user_id", "u1")

	r := NewValidator()
	assert.NoError(t, r.ValidateContext(ctx, &Request{Owner: "u1"}))
	assert.ErrorContains(t, r.ValidateContext(ctx, &Request{Owner: "u2"}), "field owner invalid")
	assert.ErrorContains(t, r.ValidateContext(context.Background(), &Request{Owner: "u1"}), "is not given")
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

	t.Run("Valid struct with empty inner struct and max0", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruit" vld:"max0 || min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "", testStruct.Fruit.Name, "Expected output to be unchanged")
	})

	t.Run("Invalid struct with filled inner struct and max0", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name string `json:"name" vld:"equapple"`
			} `json:"fruit" vld:"max0 || min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{"name": "banana"}}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error but got none")
		assert.Equal(t, "", testStruct.Fruit.Name, "Expected output to be unchanged")
	})

	t.Run("Valid struct with inner struct and cmn", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name  string `json:"name" vld:"equapple"`
				Color string `json:"color" vld:"equred"`
				Size  int    `json:"size" vld:"min1"`
			} `json:"fruit" vld:"cmn2"`
		}{}
		testStruct.Fruit.Color = "green"
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{"name": "apple", "color": "blue", "size": 3}}, testStruct, model.VLD)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "apple", testStruct.Fruit.Name, "Expected name to be updated")
		assert.Equal(t, "green", testStruct.Fruit.Color, "Expected invalid color to be left out")
		assert.Equal(t, 3, testStruct.Fruit.Size, "Expected size to be updated")
	})

	t.Run("Invalid struct with inner struct and cmn", func(t *testing.T) {
		testStruct := &struct {
			Fruit struct {
				Name  string `json:"name" vld:"equapple"`
				Color string `json:"color" vld:"equred"`
				Size  int    `json:"size" vld:"min1"`
			} `json:"fruit" vld:"cmn2"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"fruit": map[string]any{"name": "apple", "color": "blue", "size": 0}}, testStruct, model.VLD)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "less than 2 correct inner fields", "Expected cmn error")
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

	t.Run("Invalid struct with null string and min", func(t *testing.T) {
		testStruct := &struct {
			Name string `json:"name" vld:"min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"name": nil}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field name invalid", "Expected error for null field")
	})

	t.Run("Invalid struct with null array and min", func(t *testing.T) {
		testStruct := &struct {
			Tags []string `json:"tags" vld:"min1"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"tags": nil}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field tags invalid", "Expected error for null field")
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

	t.Run("Valid struct with missing optional uuid", func(t *testing.T) {
		testStruct := &struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional,omitempty" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"id": "550e8400-e29b-41d4-a716-446655440000"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Nil(t, testStruct.Optional, "Expected optional to stay nil")
	})

	t.Run("Valid struct with empty optional uuid", func(t *testing.T) {
		testStruct := &struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional,omitempty" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"id": "550e8400-e29b-41d4-a716-446655440000", "optional": ""}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Nil(t, testStruct.Optional, "Expected optional to stay nil")
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

	t.Run("Valid struct with net.IP v4", func(t *testing.T) {
		testStruct := &struct {
			IP net.IP `json:"ip" vld:"istip"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"ip": "192.168.0.1"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "192.168.0.1", testStruct.IP.String(), "Expected ip to be updated")
	})

	t.Run("Valid struct with net.IP v6", func(t *testing.T) {
		testStruct := &struct {
			IP net.IP `json:"ip" vld:"istip"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"ip": "2001:db8::1"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, "2001:db8::1", testStruct.IP.String(), "Expected ip to be updated")
	})

	t.Run("Invalid struct with net.IP", func(t *testing.T) {
		testStruct := &struct {
			IP net.IP `json:"ip" vld:"istip"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"ip": "not-an-ip"}, testStruct)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value is not a valid ip address", "Expected ip validation error")
	})

	t.Run("Valid struct with time.Time from iso8601", func(t *testing.T) {
		testStruct := &struct {
			CreatedAt time.Time `json:"created_at" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"created_at": "2025-01-02T00:00:00Z"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), testStruct.CreatedAt.UTC(), "Expected created_at to be updated")
	})

	t.Run("Embedded struct", func(t *testing.T) {
		type Base struct {
			ID string `json:"id" vld:"min1"`
		}
		type Request struct {
			*Base
			Name string `json:"name" vld:"min1"`
		}

		result := &Request{}
		require.NoError(t, r.ValidateAndUpdate(map[string]any{"id": "1", "name": "n"}, result))
		require.NotNil(t, result.Base)
		assert.Equal(t, "1", result.ID)
		assert.ErrorContains(t, r.ValidateAndUpdate(map[string]any{"name": "n"}, &Request{}), "json id key not in map")
	})

	t.Run("Recursive struct", func(t *testing.T) {
		type Node struct {
			Name     string `json:"name" vld:"min1"`
			Children []Node `json:"children" vld:"-"`
		}

		result := &Node{}
		input := map[string]any{"name": "a", "children": []any{map[string]any{"name": "b", "children": []any{}}}}
		require.NoError(t, r.ValidateAndUpdate(input, result))
		require.Len(t, result.Children, 1)
		assert.Equal(t, "b", result.Children[0].Name)

		input = map[string]any{"name": "a", "children": []any{map[string]any{"name": "", "children": []any{}}}}
		assert.ErrorContains(t, r.ValidateAndUpdate(input, &Node{}), "field children invalid")
	})

	t.Run("Valid struct with time.Time from unix", func(t *testing.T) {
		testStruct := &struct {
			CreatedAt time.Time `json:"created_at" vld:"-"`
		}{}
		err := r.ValidateAndUpdate(map[string]any{"created_at": "1735776000"}, testStruct)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, int64(1735776000), testStruct.CreatedAt.Unix(), "Expected created_at to be updated from unix")
	})
}

func TestValidateAndUpdateContext(t *testing.T) {
	type Request struct {
		Owner string `json:"owner" vld:"equ${user_id}"`
	}
	ctx := Inject(context.Background(), "user_id", "u1")

	r := NewValidator()
	result := &Request{}
	require.NoError(t, r.ValidateAndUpdateContext(ctx, map[string]any{"owner": "u1"}, result))
	assert.Equal(t, "u1", result.Owner)

	result = &Request{}
	assert.Error(t, r.ValidateAndUpdateContext(ctx, map[string]any{"owner": "u2"}, result))
	assert.Empty(t, result.Owner)
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

func TestValidateAndUpdateWithValidationContext(t *testing.T) {
	validations := []model.Validation{{Key: "owner", Type: model.String, Requirement: "equ${user_id}"}}
	ctx := Inject(context.Background(), "user_id", "u1")

	r := NewValidator()
	result := map[string]any{}
	require.NoError(t, r.ValidateAndUpdateWithValidationContext(ctx, map[string]any{"owner": "u1"}, &result, validations))
	assert.Equal(t, map[string]any{"owner": "u1"}, result)

	result = map[string]any{}
	assert.Error(t, r.ValidateAndUpdateWithValidationContext(ctx, map[string]any{"owner": "u2"}, &result, validations))
	assert.Empty(t, result)
}

func TestValidateWithValidation(t *testing.T) {
	type args struct {
		jsonMap     map[string]any
		validations []model.Validation
	}
	tests := []struct {
		name                string
		args                args
		expected            any
		expectedFieldErrors int
		wantErr             bool
	}{
		{
			name: "Invalid fields with collected field errors",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
					"color": "blue",
				},
				validations: []model.Validation{
					{Key: "fruit", Type: model.String, Requirement: "equapple"},
					{Key: "color", Type: model.String, Requirement: "equred"},
					{Key: "size", Type: model.Int, Requirement: "min1"},
				},
			},
			expected:            map[string]any{"fruit": "apple"},
			expectedFieldErrors: 2,
			wantErr:             true,
		},
		{
			name: "Invalid duplicate key",
			args: args{
				jsonMap: map[string]any{
					"fruit": "apple",
				},
				validations: []model.Validation{
					{Key: "fruit", Type: model.String, Requirement: "equapple"},
					{Key: "fruit", Type: model.String, Requirement: "equapple"},
				},
			},
			wantErr: true,
		},
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
			name: "Invalid duplicate key with different requirements",
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
			name: "Valid empty nested struct with max0",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "max0 || min1",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{}},
			wantErr:  false,
		},
		{
			name: "Invalid empty nested struct with min1",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "min1",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid empty nested struct without requirement",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{},
				},
				validations: []model.Validation{
					{
						Key:  "fruit",
						Type: model.Struct,
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid filled nested struct with max0",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "banana", "color": "red", "size": 1},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "max0 || min1",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Invalid nested struct with too many keys for max2",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "red", "size": 1},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "max2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid nested struct with cmn2 and one invalid field",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "blue", "size": 1},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "cmn2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{"name": "apple", "size": 1}},
			wantErr:  false,
		},
		{
			name: "Valid nested struct with cmn2 and one missing field",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "red"},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "cmn2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{"name": "apple", "color": "red"}},
			wantErr:  false,
		},
		{
			name: "Invalid nested struct with cmn2 and two invalid fields",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "blue", "size": 0},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "cmn2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid nested struct with cmx1 and one correct field",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "blue", "size": 0},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "cmx1",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{"name": "apple"}},
			wantErr:  false,
		},
		{
			name: "Invalid nested struct with cmx1 and two correct fields",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "red", "size": 0},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "cmx1",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Valid empty nested struct with max0 || cmn2",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "max0 || cmn2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{}},
			wantErr:  false,
		},
		{
			name: "Valid nested struct with max0 || cmn2",
			args: args{
				jsonMap: map[string]any{
					"fruit": map[string]any{"name": "apple", "color": "red", "size": 0},
				},
				validations: []model.Validation{
					{
						Key:         "fruit",
						Type:        model.Struct,
						Requirement: "max0 || cmn2",
						InnerValidation: []model.Validation{
							{Key: "name", Type: model.String, Requirement: "equapple"},
							{Key: "color", Type: model.String, Requirement: "equred"},
							{Key: "size", Type: model.Int, Requirement: "min1"},
						},
					},
				},
			},
			expected: map[string]any{"fruit": map[string]any{"name": "apple", "color": "red"}},
			wantErr:  false,
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
				if test.expectedFieldErrors > 0 {
					var fieldErrors model.FieldErrors
					require.ErrorAs(t, err, &fieldErrors, "Expected field errors")
					assert.Len(t, fieldErrors, test.expectedFieldErrors, "Expected number of field errors to match")
					assert.Equal(t, test.expected, out, "Expected output to contain the valid fields")
				}
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

func TestValidateWithValidationContext(t *testing.T) {
	validations := []model.Validation{{Key: "amount", Type: model.Int, Requirement: "max${limit}"}}
	ctx := Inject(context.Background(), "limit", 100)

	r := NewValidator()
	validated, err := r.ValidateWithValidationContext(ctx, map[string]any{"amount": 50.0}, validations)
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"amount": 50}, validated)

	_, err = r.ValidateWithValidationContext(ctx, map[string]any{"amount": 150.0}, validations)
	assert.ErrorContains(t, err, "value does not fulfill max${limit}")
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
	correctTwo := 2
	type args struct {
		input    any
		astValue *model.AstValue
		correct  *int
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
			name: "Valid cmn with correct inner fields",
			args: args{
				input: map[string]any{"name": "apple", "color": "red"},
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CORRECT_MIN,
							ConditionValue: "2",
						},
					},
				},
				correct: &correctTwo,
			},
			wantErr: false,
		},
		{
			name: "Invalid cmx with correct inner fields",
			args: args{
				input: map[string]any{"name": "apple", "color": "red"},
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CORRECT_MAX,
							ConditionValue: "1",
						},
					},
				},
				correct: &correctTwo,
			},
			wantErr: true,
		},
		{
			name: "Invalid cmn without inner validations",
			args: args{
				input: map[string]any{"name": "apple"},
				astValue: &model.AstValue{
					ConditionGroup: []*model.AstValue{
						{
							Type:           model.CONDITION,
							ConditionType:  model.CORRECT_MIN,
							ConditionValue: "1",
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

			err := r.RunValidatorsOnConditionGroup(test.args.input, test.args.astValue, test.args.correct)
			if test.wantErr {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}

	t.Run("Implication", func(t *testing.T) {
		tests := []struct {
			requirement string
			value       any
			wantErr     bool
		}{
			{requirement: "equoci => min5", value: "oci", wantErr: true},
			{requirement: "min3 => max5", value: "ab"},
			{requirement: "min3 => max5", value: "abcd"},
			{requirement: "min3 => max5", value: "abcdef", wantErr: true},
			{requirement: "min1 && max3 => equab || equabc", value: "abc"},
			{requirement: "min1 && max3 => equab || equabc", value: "a", wantErr: true},
			{requirement: "min1 && max3 => equab || equabc", value: "abcde"},
			{requirement: "min1 => min2 => min3", value: "ab", wantErr: true},
			{requirement: "min1 => min2 => min3", value: "a"},
			{requirement: "(min3 => max5) && conb", value: "abcdef", wantErr: true},
			{requirement: "(min3 => max5) && conb", value: "ab"},
		}

		r := NewValidator()
		for _, test := range tests {
			t.Run(fmt.Sprintf("%v with %v", test.requirement, test.value), func(t *testing.T) {
				err := r.ValidateValueWithParser(test.value, &model.Validation{Requirement: test.requirement})
				if test.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}

		t.Run("Error names the premise", func(t *testing.T) {
			err := r.ValidateValueWithParser("abcdef", &model.Validation{Requirement: "min3 => max5"})
			require.Error(t, err)
			assert.Contains(t, err.Error(), "required because of min'3'")
		})
	})

	t.Run("OR group stops at the first fulfilled condition", func(t *testing.T) {
		r := NewValidator()
		called := 0
		r.AddValidationFunc(func(input any, astValue *model.AstValue) error {
			called++
			return nil
		}, "count")

		assert.NoError(t, r.ValidateValueWithParser("abc", &model.Validation{Requirement: "min1 || funcount"}))
		assert.Equal(t, 0, called, "The second condition of a fulfilled OR group does not have to be checked")

		assert.NoError(t, r.ValidateValueWithParser("abc", &model.Validation{Requirement: "max1 || funcount"}))
		assert.Equal(t, 1, called)
	})
}

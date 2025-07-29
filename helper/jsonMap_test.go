package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetStructValueByJson(t *testing.T) {
	t.Run("Set string", func(t *testing.T) {
		type TestStruct struct {
			String string `json:"string"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("String")

		t.Run("Set valid string", func(t *testing.T) {
			err := SetStructValueByJson(fv, "apple")
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, "apple", result.String, "Expected String to be set correctly")
		})

		t.Run("Set invalid string", func(t *testing.T) {
			err := SetStructValueByJson(fv, 1)
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, "apple", result.String, "Expected String to remain unchanged")
		})
	})

	t.Run("Set bool", func(t *testing.T) {
		type TestStruct struct {
			Bool bool `json:"bool"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Bool")

		t.Run("Set valid bool", func(t *testing.T) {
			err := SetStructValueByJson(fv, true)
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, true, result.Bool, "Expected Bool to be set correctly")
		})

		t.Run("Set invalid bool", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not a bool")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, true, result.Bool, "Expected Bool to remain unchanged")
		})
	})

	t.Run("Set int", func(t *testing.T) {
		type TestStruct struct {
			Int int `json:"int"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Int")

		t.Run("Set valid int", func(t *testing.T) {
			err := SetStructValueByJson(fv, 2)
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, 2, result.Int, "Expected Int to be set correctly")
		})

		t.Run("Set invalid int", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not an int")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, 2, result.Int, "Expected Int to remain unchanged")
		})
	})

	t.Run("Set uint", func(t *testing.T) {
		type TestStruct struct {
			Uint uint `json:"uint"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Uint")

		t.Run("Set valid uint", func(t *testing.T) {
			err := SetStructValueByJson(fv, 2)
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, uint(2), result.Uint, "Expected Uint to be set correctly")
		})

		t.Run("Set invalid uint", func(t *testing.T) {
			err := SetStructValueByJson(fv, "apple")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, uint(2), result.Uint, "Expected Uint to remain unchanged")
		})
	})

	t.Run("Set float", func(t *testing.T) {
		type TestStruct struct {
			Float float64 `json:"float"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Float")

		t.Run("Set valid float", func(t *testing.T) {
			err := SetStructValueByJson(fv, 3.14)
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, 3.14, result.Float, "Expected Float to be set correctly")
		})

		t.Run("Set invalid float", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not a float")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, 3.14, result.Float, "Expected Float to remain unchanged")
		})
	})

	t.Run("Set struct", func(t *testing.T) {
		type InnerStruct struct {
			Field string `json:"field"`
		}
		type TestStruct struct {
			Inner InnerStruct `json:"inner"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Inner")

		t.Run("Set valid struct", func(t *testing.T) {
			err := SetStructValueByJson(fv, map[string]any{"field": "apple"})
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, "apple", result.Inner.Field, "Expected Inner.Field to be set correctly")
		})

		t.Run("Set invalid struct", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not a struct")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, "apple", result.Inner.Field, "Expected Inner.Field to remain unchanged")
		})
	})

	t.Run("Set map", func(t *testing.T) {
		type TestStruct struct {
			Map map[string]string `json:"map"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Map")

		t.Run("Set valid map", func(t *testing.T) {
			err := SetStructValueByJson(fv, map[string]any{"key1": "value1", "key2": "value2"})
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, result.Map, "Expected Map to be set correctly")
		})

		t.Run("Set invalid map", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not a map")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, result.Map, "Expected Map to remain unchanged")
		})

		t.Run("Set invalid map value type", func(t *testing.T) {
			err := SetStructValueByJson(fv, map[string]any{"key1": true, "key2": false})
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, map[string]string{"key1": "value1", "key2": "value2"}, result.Map, "Expected Map to remain unchanged")
		})
	})

	t.Run("Set array", func(t *testing.T) {
		type TestStruct struct {
			Array []string `json:"array"`
		}
		result := &TestStruct{}
		fv := reflect.ValueOf(result).Elem().FieldByName("Array")

		t.Run("Set valid array", func(t *testing.T) {
			err := SetStructValueByJson(fv, []string{"apple", "banana"})
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, []string{"apple", "banana"}, result.Array, "Expected Array to be set correctly")
		})

		t.Run("Set valid array json", func(t *testing.T) {
			err := SetStructValueByJson(fv, []any{"apple", "banana"})
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, []string{"apple", "banana"}, result.Array, "Expected Array to be set correctly")
		})

		t.Run("Set invalid array", func(t *testing.T) {
			err := SetStructValueByJson(fv, []int{1, 2})
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, []string{"apple", "banana"}, result.Array, "Expected Array to remain unchanged")
		})

		t.Run("Set invalid json", func(t *testing.T) {
			err := SetStructValueByJson(fv, "not an array")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, []string{"apple", "banana"}, result.Array, "Expected Array to remain unchanged")
		})

		t.Run("Set invalid json array", func(t *testing.T) {
			err := SetStructValueByJson(fv, []any{true, true})
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, []string{"apple", "banana"}, result.Array, "Expected Array to remain unchanged")
		})
	})

	t.Run("Set array of struct", func(t *testing.T) {
		type InnerStruct struct {
			Name string `json:"name"`
		}
		type TestStructWithArray struct {
			Array []InnerStruct `json:"array"`
		}
		resultWithArray := &TestStructWithArray{}
		fvWithArray := reflect.ValueOf(resultWithArray).Elem().FieldByName("Array")

		t.Run("Set valid array of struct", func(t *testing.T) {
			err := SetStructValueByJson(fvWithArray, []any{
				map[string]any{"name": "apple"},
				map[string]any{"name": "banana"},
			})
			assert.NoError(t, err, "Expected no error setting struct value by json map")
			assert.Equal(t, []InnerStruct{{Name: "apple"}, {Name: "banana"}}, resultWithArray.Array, "Expected Array of structs to be set correctly")
		})

		t.Run("Set invalid array of struct", func(t *testing.T) {
			err := SetStructValueByJson(fvWithArray, []any{
				map[string]any{"name": true},
				map[string]any{"name": true},
			})
			assert.Error(t, err, "Expected error setting struct value by json map with invalid type")
			assert.Equal(t, []InnerStruct{{Name: "apple"}, {Name: "banana"}}, resultWithArray.Array, "Expected Array of structs to remain unchanged")
		})

		t.Run("Set invalid json for array of struct", func(t *testing.T) {
			err := SetStructValueByJson(fvWithArray, "not an array")
			assert.Error(t, err, "Expected error setting struct value by json map with invalid json")
			assert.Equal(t, []InnerStruct{{Name: "apple"}, {Name: "banana"}}, resultWithArray.Array, "Expected Array of structs to remain unchanged")
		})

		t.Run("Set invalid json array for array of struct", func(t *testing.T) {
			err := SetStructValueByJson(fvWithArray, []any{
				map[string]int{"name": 1},
				map[string]int{"name": 2},
			})
			assert.Error(t, err, "Expected error setting struct value by json map with invalid json")
			assert.Equal(t, []InnerStruct{{Name: "apple"}, {Name: "banana"}}, resultWithArray.Array, "Expected Array of structs to remain unchanged")
		})
	})
}

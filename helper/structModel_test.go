package helper

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type modelBase struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// modelOther has the key name in the yaml tag, because go vet reports the same json key in two embedded structs.
type modelOther struct {
	Name  string `yaml:"name"`
	Other string `json:"other"`
}

type modelDeep struct {
	modelBase
}

type modelNode struct {
	*modelNode
	Value    int         `json:"value"`
	Children []modelNode `json:"children"`
}

func modelKeys(t *testing.T, value any) map[string][]int {
	t.Helper()
	keys := map[string][]int{}
	for _, field := range StructModelOf(reflect.TypeOf(value)).Fields {
		keys[field.Key] = field.Index
	}
	return keys
}

func TestStructModelOf(t *testing.T) {
	t.Run("Keys and ignored fields", func(t *testing.T) {
		type Model struct {
			Name       string `json:"name,omitempty"`
			Yaml       string `yaml:"yaml"`
			NoTag      string
			Ignored    string `json:"-"`
			YamlIgnore string `yaml:"-"`
			unexported string
		}

		structModel := StructModelOf(reflect.TypeOf(Model{}))
		require.Len(t, structModel.Fields, 3)
		assert.Equal(t, "name", structModel.Fields[0].Key)
		assert.True(t, structModel.Fields[0].OmitEmpty)
		assert.Equal(t, "yaml", structModel.Fields[1].Key)
		assert.Equal(t, "NoTag", structModel.Fields[2].Key)
		assert.Same(t, structModel, StructModelOf(reflect.TypeOf(Model{})), "Expected the model to be cached")
	})

	t.Run("Embedded and inline structs are flattened", func(t *testing.T) {
		type Address struct {
			City string `json:"city"`
		}
		type Nested struct {
			Name string `json:"name"`
		}
		type Model struct {
			modelBase
			*Nested `json:"nested"`
			*modelOther
			Address *Address `json:",inline"`
			Count   int      `json:"count"`
		}

		assert.Equal(t, map[string][]int{
			"id":     {0, 0},
			"name":   {0, 1},
			"nested": {1},
			"city":   {3, 0},
			"count":  {4},
		}, modelKeys(t, Model{}), "Expected a pointer to an unexported struct to be left out, it can not be set")
	})

	t.Run("Conflicting keys", func(t *testing.T) {
		type SameDepth struct {
			modelBase
			modelOther
		}
		type Shallow struct {
			modelDeep
			modelOther
		}
		type Named struct {
			modelBase
			Other struct {
				Name string `json:"name"`
			} `json:",inline"`
		}
		type Untagged struct {
			Name string
		}
		type NamedWins struct {
			Untagged
			modelOther
		}

		assert.Equal(t, map[string][]int{"id": {0, 0}, "other": {1, 1}}, modelKeys(t, SameDepth{}), "Expected name to be dropped at the same depth")
		assert.Equal(t, map[string][]int{"id": {0, 0, 0}, "name": {1, 0}, "other": {1, 1}}, modelKeys(t, Shallow{}), "Expected the shallowest field to win")
		assert.Equal(t, map[string][]int{"id": {0, 0}}, modelKeys(t, Named{}), "Expected two named fields at the same depth to be dropped")
		assert.Equal(t, map[string][]int{"Name": {0, 0}, "name": {1, 0}, "other": {1, 1}}, modelKeys(t, NamedWins{}), "Expected keys to be case-sensitive")
	})

	t.Run("Recursive types", func(t *testing.T) {
		assert.Equal(t, map[string][]int{"value": {1}, "children": {2}}, modelKeys(t, modelNode{}))
	})

	t.Run("Fields decoded with encoding/json", func(t *testing.T) {
		type Model struct {
			Number   int          `json:"number,string"`
			Raw      rawJSON      `json:"raw"`
			Pointer  *rawJSON     `json:"pointer"`
			Duration int64        `json:"duration,format:nano"`
			Plain    int          `json:"plain"`
			Time     time.Time    `json:"time"`
			Struct   *plainStruct `json:"struct"`
		}

		decoded := map[string]bool{}
		for _, field := range StructModelOf(reflect.TypeOf(Model{})).Fields {
			decoded[field.Key] = field.decodeType != nil
		}
		assert.Equal(t, map[string]bool{"number": true, "raw": true, "pointer": true, "duration": true, "plain": false, "time": false, "struct": false}, decoded)
	})
}

// rawJSON keeps the json of its value, to test fields with an own UnmarshalJSON.
type rawJSON string

func (r *rawJSON) UnmarshalJSON(data []byte) error {
	*r = rawJSON(data)
	return nil
}

type plainStruct struct {
	Value string `json:"value"`
}

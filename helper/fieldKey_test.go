package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldKey(t *testing.T) {
	type TestStruct struct {
		NoTag        string
		JsonTag      string `json:"json_tag"`
		JsonOmit     string `json:"json_omit,omitempty"`
		JsonDash     string `json:"-"`
		JsonInline   string `json:",inline"`
		YamlTag      string `yaml:"yaml_tag"`
		YamlOptions  string `yaml:"yaml_options,omitempty,inline"`
		JsonOverYaml string `json:"json_key" yaml:"yaml_key,omitempty"`
	}

	tests := []struct {
		field     string
		key       string
		omitEmpty bool
		inline    bool
	}{
		{"NoTag", "NoTag", false, false},
		{"JsonTag", "json_tag", false, false},
		{"JsonOmit", "json_omit", true, false},
		{"JsonDash", "JsonDash", false, false},
		{"JsonInline", "JsonInline", false, true},
		{"YamlTag", "yaml_tag", false, false},
		{"YamlOptions", "yaml_options", true, true},
		{"JsonOverYaml", "json_key", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			field, ok := reflect.TypeOf(TestStruct{}).FieldByName(tt.field)
			assert.True(t, ok, "Expected field to exist")

			key, omitEmpty, inline := FieldKey(field)
			assert.Equal(t, tt.key, key, "Expected key to match")
			assert.Equal(t, tt.omitEmpty, omitEmpty, "Expected omitEmpty to match")
			assert.Equal(t, tt.inline, inline, "Expected inline to match")
		})
	}
}

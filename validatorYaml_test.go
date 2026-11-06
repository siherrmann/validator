package validator

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type yamlItem struct {
	Name  string `yaml:"name" vld:"min1"`
	Count int    `yaml:"count" vld:"min1"`
}

type yamlConfig struct {
	ID        uuid.UUID         `yaml:"id" vld:"istuuid"`
	Name      string            `yaml:"name" vld:"min3"`
	Port      int               `yaml:"port" vld:"min1 && max65535"`
	Ratio     float64           `yaml:"ratio" vld:"min0"`
	Debug     bool              `yaml:"debug" vld:"-"`
	CreatedAt time.Time         `yaml:"created_at" vld:"-"`
	Labels    map[string]string `yaml:"labels,omitempty" vld:"-"`
	Items     []yamlItem        `yaml:"items" vld:"min1"`
	Nested    yamlItem          `yaml:"nested" vld:"-"`
	JsonKey   string            `json:"json_key" yaml:"yaml_key" vld:"-"`
}

const validYamlConfig = `
id: 550e8400-e29b-41d4-a716-446655440000
name: service
port: 8080
ratio: 1
debug: true
created_at: 2025-01-02T00:00:00Z
labels:
  env: prod
items:
  - name: first
    count: 1
  - name: second
    count: 2
nested:
  name: inner
  count: 3
json_key: from json tag
`

func TestUnmarshalYamlValidateAndUpdate(t *testing.T) {
	t.Run("Valid yaml", func(t *testing.T) {
		config := &yamlConfig{}
		err := NewValidator().UnmarshalYamlValidateAndUpdate([]byte(validYamlConfig), config)
		require.NoError(t, err, "Expected no error but got one")

		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), config.ID, "Expected id to be updated")
		assert.Equal(t, "service", config.Name, "Expected name to be updated")
		assert.Equal(t, 8080, config.Port, "Expected port to be updated")
		assert.Equal(t, 1.0, config.Ratio, "Expected int yaml value to be set as float")
		assert.True(t, config.Debug, "Expected debug to be updated")
		assert.Equal(t, time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), config.CreatedAt, "Expected created_at to be updated")
		assert.Equal(t, map[string]string{"env": "prod"}, config.Labels, "Expected labels to be updated")
		assert.Equal(t, []yamlItem{{Name: "first", Count: 1}, {Name: "second", Count: 2}}, config.Items, "Expected items to be updated")
		assert.Equal(t, yamlItem{Name: "inner", Count: 3}, config.Nested, "Expected nested to be updated")
		assert.Equal(t, "from json tag", config.JsonKey, "Expected json tag to take precedence over yaml tag")
	})

	t.Run("Invalid value", func(t *testing.T) {
		err := NewValidator().UnmarshalYamlValidateAndUpdate([]byte("id: 550e8400-e29b-41d4-a716-446655440000\nname: ab\n"), &yamlConfig{})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field name invalid", "Expected error for name")
	})

	t.Run("Invalid inner value", func(t *testing.T) {
		input := "id: 550e8400-e29b-41d4-a716-446655440000\nname: service\nport: 80\nratio: 1\nitems:\n  - name: first\n    count: 0\n"
		err := NewValidator().UnmarshalYamlValidateAndUpdate([]byte(input), &yamlConfig{})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field items invalid", "Expected error for items")
	})

	t.Run("Invalid uuid", func(t *testing.T) {
		err := NewValidator().UnmarshalYamlValidateAndUpdate([]byte("id: not-a-uuid\n"), &yamlConfig{})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "value is not a valid uuid", "Expected error for id")
	})

	t.Run("Invalid yaml", func(t *testing.T) {
		err := NewValidator().UnmarshalYamlValidateAndUpdate([]byte("name: [unclosed\n"), &yamlConfig{})
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "error unmarshaling yaml", "Expected unmarshal error")
	})
}

func TestUnmarshalYamlValidateAndUpdateWithValidation(t *testing.T) {
	validations := []model.Validation{
		{Key: "name", Type: model.String, Requirement: "equapple"},
		{Key: "age", Type: model.Int, Requirement: "min2"},
	}

	t.Run("Valid yaml", func(t *testing.T) {
		mapToUpdate := map[string]any{}
		err := NewValidator().UnmarshalYamlValidateAndUpdateWithValidation([]byte("name: apple\nage: 3\n"), &mapToUpdate, validations)
		assert.NoError(t, err, "Expected no error but got one")
		assert.Equal(t, map[string]any{"name": "apple", "age": 3}, mapToUpdate, "Expected map to be updated")
	})

	t.Run("Invalid yaml value", func(t *testing.T) {
		mapToUpdate := map[string]any{}
		err := NewValidator().UnmarshalYamlValidateAndUpdateWithValidation([]byte("name: apple\nage: 1\n"), &mapToUpdate, validations)
		assert.Error(t, err, "Expected an error but got none")
		assert.Contains(t, err.Error(), "field age invalid", "Expected error for age")
		assert.Empty(t, mapToUpdate, "Expected map not to be updated")
	})
}

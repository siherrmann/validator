package helper

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshalRequestToJsonMap(t *testing.T) {
	t.Run("Valid JSON request", func(t *testing.T) {
		data := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(data))
		require.NoError(t, err, "Expected no error creating request")
		req.Header.Set("Content-Type", "application/json")

		mapOut, err := UnmarshalRequestToJsonMap(req)
		assert.NoError(t, err, "Expected no error unmarshaling request to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
	})

	t.Run("Invalid request", func(t *testing.T) {
		_, err := UnmarshalRequestToJsonMap(nil)
		assert.Error(t, err, "Expected error unmarshaling invalid JSON request")
		assert.Contains(t, err.Error(), "request is nil", "Expected error to contain nil error")
	})
}

func TestUnmapRequestToJsonMap(t *testing.T) {
	t.Run("Valid form request", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		require.NoError(t, err, "Expected no error creating request")
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		mapOut, err := UnmapRequestToJsonMap(req)
		assert.NoError(t, err, "Expected no error unmapping request to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
	})

	t.Run("Invalid request", func(t *testing.T) {
		_, err := UnmapRequestToJsonMap(nil)
		assert.Error(t, err, "Expected error unmapping invalid request")
		assert.Contains(t, err.Error(), "request is nil", "Expected error to contain nil error")
	})
}

func TestUnmarshalJsonToJsonMap(t *testing.T) {
	t.Run("Valid JSON", func(t *testing.T) {
		jsonData := []byte(`{"key1": "value1", "key2": 2, "key3": true}`)

		mapOut, err := UnmarshalJsonToJsonMap(jsonData)
		assert.NoError(t, err, "Expected no error when unmarshaling JSON to JsonMap")
		assert.Equal(t, "value1", mapOut["key1"], "Expected key1 to match")
		assert.Equal(t, float64(2), mapOut["key2"], "Expected key2 to match")
		assert.Equal(t, true, mapOut["key3"], "Expected key3 to match")
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		jsonData := []byte(`<html><body>Invalid JSON</body></html>`)

		_, err := UnmarshalJsonToJsonMap(jsonData)
		assert.Error(t, err, "Expected error when unmarshaling invalid JSON")
		assert.Contains(t, err.Error(), "error unmarshaling:", "Expected error to contain JSON parsing error")
	})
}

func TestUnmapUrlValuesToJsonMap(t *testing.T) {
	t.Run("Valid URL values", func(t *testing.T) {
		values := url.Values{}
		values.Set("name", "apple")
		values.Set("age", "2")
		values["array"] = []string{"fruit", "food"}
		values["map"] = []string{"{\"fruit\": \"apple\", \"food\": \"banana\"}"}

		mapOut, err := UnmapUrlValuesToJsonMap(values)
		assert.NoError(t, err, "Expected no error unmapping URL values to JsonMap")
		assert.Equal(t, "apple", mapOut["name"], "Expected name to be 'apple'")
		assert.Equal(t, float64(2), mapOut["age"], "Expected age to be 2")
		assert.ElementsMatch(t, []string{"fruit", "food"}, mapOut["array"], "Expected array to match")
		assert.Equal(t, map[string]any{"food": "banana", "fruit": "apple"}, mapOut["map"], "Expected map to match")
	})

	t.Run("Empty URL values", func(t *testing.T) {
		values := url.Values{}

		mapOut, err := UnmapUrlValuesToJsonMap(values)
		assert.NoError(t, err, "Expected no error unmapping empty URL values to JsonMap")
		assert.Empty(t, mapOut, "Expected empty JsonMap for empty URL values")
	})
}

package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name string `json:"name" vld:"equapple"`
	Age  int    `json:"age" vld:"min2"`
}

func newValidator() *Validator {
	return &Validator{}
}

func TestUnmapOrUnmarshalAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal and validate")
	})

	t.Run("Invalid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmap or unmarshal and validate")
		assert.Contains(t, err.Error(), "error validating struct", "Expected error to contain validation error")
	})

	t.Run("Valid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"apple","age":2}`))

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal and validate")
	})

	t.Run("Invalid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"banana","age":1}`))

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmap or unmarshal and validate")
		assert.Contains(t, err.Error(), "error validating struct", "Expected error to contain validation error")
	})
}

func TestUnmarshalAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid JSON input", func(t *testing.T) {
		jsonInput := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmarshal and validate")
	})

	t.Run("Invalid JSON values", func(t *testing.T) {
		jsonInput := []byte(`{"name":"banana","age":1}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal and validate")
		assert.Contains(t, err.Error(), "error validating struct", "Expected error to contain validation error")
	})

	t.Run("Invalid JSON input", func(t *testing.T) {
		jsonInput := []byte(`<html><body>Invalid JSON</body></html>`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal and validate with invalid JSON")
		assert.Contains(t, err.Error(), "error unmarshaling json:", "Expected error to contain unmarshaling error")
	})

	t.Run("Invalid struct type", func(t *testing.T) {
		ts := testStruct{}
		jsonInput := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal and validate with invalid struct type")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain validation error")
	})
}

func TestUnmapAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap and validate")
	})

	t.Run("Invalid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmap and validate")
		assert.Contains(t, err.Error(), "error validating struct", "Expected error to contain validation error")
	})

	t.Run("Invalid struct type", func(t *testing.T) {
		ts := testStruct{}
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapAndValidate(req, ts)
		assert.Error(t, err, "Expected error on unmap and validate")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain validation error")
	})
}
func TestUnmapOrUnmarshalValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update")
	})

	t.Run("Invalid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmap or unmarshal validate and update")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})

	t.Run("Valid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"apple","age":2}`))

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update")
	})

	t.Run("Invalid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"banana","age":1}`))

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmap or unmarshal validate and update")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})
}

func TestUnmarshalValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid JSON input", func(t *testing.T) {
		jsonInput := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmarshal validate and update")
	})

	t.Run("Invalid JSON values", func(t *testing.T) {
		jsonInput := []byte(`{"name":"banana","age":1}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal validate and update")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})

	t.Run("Invalid JSON input", func(t *testing.T) {
		jsonInput := []byte(`<html><body>Invalid JSON</body></html>`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal validate and update with invalid JSON")
		assert.Contains(t, err.Error(), "error unmarshaling request body", "Expected error to contain unmarshaling error")
	})

	t.Run("Invalid struct type", func(t *testing.T) {
		ts := testStruct{}
		jsonInput := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmarshal validate and update with invalid struct type")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain validation error")
	})
}

func TestUnmapValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}

	t.Run("Valid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap validate and update")
	})

	t.Run("Invalid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmap validate and update")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})

	t.Run("Invalid struct type", func(t *testing.T) {
		ts := testStruct{}
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapValidateAndUpdate(req, ts)
		assert.Error(t, err, "Expected error on unmap validate and update with invalid struct type")
		assert.Contains(t, err.Error(), "value has to be of kind pointer", "Expected error to contain validation error")
	})
}

func TestUnmapOrUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	mapToUpdate := &model.JsonMap{}

	t.Run("Valid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update with validation")
	})

	t.Run("Invalid with form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.Error(t, err, "Expected error on unmap or unmarshal validate and update with validation")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})

	t.Run("Valid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"apple","age":2}`))

		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update with validation")
	})

	t.Run("Invalid with request body", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"banana","age":1}`))

		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.Error(t, err, "Expected error on unmap or unmarshal validate and update with validation")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})
}

func TestUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	mapToUpdate := &model.JsonMap{}

	t.Run("Valid JSON input", func(t *testing.T) {
		jsonInput := []byte(`{"name":"apple","age":2}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.NoError(t, err, "Expected no error on unmarshal validate and update with validation")
	})

	t.Run("Invalid JSON values", func(t *testing.T) {
		jsonInput := []byte(`{"name":"banana","age":1}`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.Error(t, err, "Expected error on unmarshal validate and update with validation")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})

	t.Run("Invalid JSON input", func(t *testing.T) {
		jsonInput := []byte(`<html><body>Invalid JSON</body></html>`)
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(jsonInput))
		req.Header.Set("Content-Type", "application/json")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.Error(t, err, "Expected error on unmarshal validate and update with validation with invalid JSON")
		assert.Contains(t, err.Error(), "error unmarshaling", "Expected error to contain unmarshaling error")
	})
}

func TestUnmapValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	mapToUpdate := &model.JsonMap{}

	t.Run("Valid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "2")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.NoError(t, err, "Expected no error on unmap validate and update with validation")
	})

	t.Run("Invalid form values", func(t *testing.T) {
		form := url.Values{}
		form.Set("name", "banana")
		form.Set("age", "1")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmapValidateAndUpdateWithValidation(req, mapToUpdate, validations)
		assert.Error(t, err, "Expected error on unmap validate and update with validation")
		assert.Contains(t, err.Error(), "error updating struct", "Expected error to contain validation error")
	})
}

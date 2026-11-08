package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/uuid"
	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	Name string `json:"name" vld:"equapple"`
	Age  int    `json:"age" vld:"min2"`
}

type testKind string

type testSliceStruct struct {
	Tool  string     `json:"tool" vld:"-"`
	Tools []string   `json:"tools" vld:"-"`
	IDs   []int      `json:"ids" vld:"-"`
	Kinds []testKind `json:"kinds" vld:"-"`
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

	t.Run("Valid with one form value for slice fields", func(t *testing.T) {
		ts := &testSliceStruct{}
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"tool": {"miner"}, "tools": {"miner"}, "ids": {"7"}, "kinds": {"miner"}}

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal and validate")
		assert.Equal(t, &testSliceStruct{Tool: "miner", Tools: []string{"miner"}, IDs: []int{7}, Kinds: []testKind{"miner"}}, ts)
	})

	t.Run("Valid with two form values for slice fields", func(t *testing.T) {
		ts := &testSliceStruct{}
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"tool": {"miner"}, "tools": {"miner", "run"}, "ids": {"7", "8"}, "kinds": {"miner", "run"}}

		err := v.UnmapOrUnmarshalAndValidate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal and validate")
		assert.Equal(t, &testSliceStruct{Tool: "miner", Tools: []string{"miner", "run"}, IDs: []int{7, 8}, Kinds: []testKind{"miner", "run"}}, ts)
	})

	t.Run("Invalid with non-convertible form value for slice field", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"ids": {"abc"}}

		err := v.UnmapOrUnmarshalAndValidate(req, &testSliceStruct{})
		assert.Error(t, err, "Expected error on unmap or unmarshal and validate")
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

	t.Run("Valid JSON input with uuid", func(t *testing.T) {
		testStruct := &struct {
			ID uuid.UUID `json:"id" vld:"istuuid"`
		}{}
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(`{"id":"550e8400-e29b-41d4-a716-446655440000"}`))
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalAndValidate(req, testStruct)
		assert.NoError(t, err, "Expected no error on unmarshal and validate")
		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), testStruct.ID, "Expected id to be updated")
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

	t.Run("Valid with one form value for slice fields", func(t *testing.T) {
		ts := &testSliceStruct{}
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"tool": {"miner"}, "tools": {"miner"}, "ids": {"7"}, "kinds": {"miner"}}

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update")
		assert.Equal(t, &testSliceStruct{Tool: "miner", Tools: []string{"miner"}, IDs: []int{7}, Kinds: []testKind{"miner"}}, ts)
	})

	t.Run("Valid with two form values for slice fields", func(t *testing.T) {
		ts := &testSliceStruct{}
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"tool": {"miner"}, "tools": {"miner", "run"}, "ids": {"7", "8"}, "kinds": {"miner", "run"}}

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
		assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update")
		assert.Equal(t, &testSliceStruct{Tool: "miner", Tools: []string{"miner", "run"}, IDs: []int{7, 8}, Kinds: []testKind{"miner", "run"}}, ts)
	})

	t.Run("Invalid with non-convertible form value for slice field", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = url.Values{"ids": {"abc"}}

		err := v.UnmapOrUnmarshalValidateAndUpdate(req, &testSliceStruct{})
		assert.Error(t, err, "Expected error on unmap or unmarshal validate and update")
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

	t.Run("Valid JSON input with uuid", func(t *testing.T) {
		testStruct := &struct {
			ID       uuid.UUID  `json:"id" vld:"istuuid"`
			Optional *uuid.UUID `json:"optional,omitempty" vld:"-"`
		}{}
		id := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(`{"id":"`+id.String()+`","optional":"`+id.String()+`"}`))
		require.NoError(t, err, "Expected no error creating request")

		err = v.UnmarshalValidateAndUpdate(req, testStruct)
		assert.NoError(t, err, "Expected no error on unmarshal validate and update")
		assert.Equal(t, id, testStruct.ID, "Expected id to be updated")
		require.NotNil(t, testStruct.Optional, "Expected optional to be set")
		assert.Equal(t, id, *testStruct.Optional, "Expected optional to be updated")
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

	t.Run("Lone form value validated as list for slice field", func(t *testing.T) {
		type minOne struct {
			Tools []string `json:"tools" vld:"min1"`
		}
		type maxOne struct {
			Tools []string `json:"tools" vld:"max1"`
		}
		type maxZero struct {
			Tools []string `json:"tools" vld:"max0"`
		}
		type from struct {
			Tools []string `json:"tools" vld:"frmminer,run"`
		}
		type contains struct {
			IDs []int `json:"ids" vld:"con7"`
		}
		type maxZeroInt struct {
			IDs []int `json:"ids" vld:"max0"`
		}

		tests := []struct {
			name    string
			form    url.Values
			target  any
			wantErr bool
		}{
			{"min1 passes", url.Values{"tools": {"miner"}}, &minOne{}, false},
			{"max1 passes", url.Values{"tools": {"miner"}}, &maxOne{}, false},
			{"max0 fails", url.Values{"tools": {"miner"}}, &maxZero{}, true},
			{"frm passes", url.Values{"tools": {"miner"}}, &from{}, false},
			{"frm fails", url.Values{"tools": {"other"}}, &from{}, true},
			{"con on lone number passes", url.Values{"ids": {"7"}}, &contains{}, false},
			{"con on lone number fails", url.Values{"ids": {"8"}}, &contains{}, true},
			{"max0 on lone number fails", url.Values{"ids": {"7"}}, &maxZeroInt{}, true},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				req, err := http.NewRequest("POST", "/", bytes.NewBufferString(tt.form.Encode()))
				req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
				require.NoError(t, err, "Expected no error creating request")

				err = v.UnmapValidateAndUpdate(req, tt.target)
				if tt.wantErr {
					assert.Error(t, err, "Expected error on unmap validate and update")
				} else {
					assert.NoError(t, err, "Expected no error on unmap validate and update")
				}
			})
		}
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

	t.Run("Valid form values with uuid", func(t *testing.T) {
		testStruct := &struct {
			ID uuid.UUID `json:"id" vld:"istuuid"`
		}{}
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		require.NoError(t, err, "Expected no error creating request")
		req.Form = url.Values{"id": {"550e8400-e29b-41d4-a716-446655440000"}}

		err = v.UnmapValidateAndUpdate(req, testStruct)
		assert.NoError(t, err, "Expected no error on unmap validate and update")
		assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"), testStruct.ID, "Expected id to be updated")
	})
}

func TestUnmapOrUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	mapToUpdate := &map[string]any{}

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
	mapToUpdate := &map[string]any{}

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
	mapToUpdate := &map[string]any{}

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

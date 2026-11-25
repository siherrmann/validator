package validator

import (
	"bytes"
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWrappedValidate(t *testing.T) {
	s := testStruct{Name: "apple", Age: 2}
	err := Validate(&s)
	assert.NoError(t, err, "Expected no error on validate")
}

func TestWrappedValidateWithValidation(t *testing.T) {
	jsonInput := map[string]any{"name": "apple", "age": 2}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	result, err := ValidateWithValidation(jsonInput, validations)
	assert.NoError(t, err, "Expected no error on validate with validation")
	assert.Equal(t, "apple", result["name"], "Expected name to be 'apple'")
	assert.Equal(t, 2, result["age"], "Expected age to be 2")
}

func TestWrappedValidateAndUpdate(t *testing.T) {
	jsonInput := map[string]any{"name": "apple", "age": 2}
	var s testStruct

	err := ValidateAndUpdate(jsonInput, &s)
	assert.NoError(t, err, "Expected no error on validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedValidateAndUpdateWithValidation(t *testing.T) {
	jsonInput := map[string]any{"name": "apple", "age": 2}
	mapToUpdate := map[string]any{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := ValidateAndUpdateWithValidation(jsonInput, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, 2, mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnpackAndValidate(t *testing.T) {
	var s testStruct
	err := UnpackAndValidate(newBodyRequest(t, "application/json", `{"name":"apple","age":2}`), &s)
	assert.NoError(t, err, "Expected no error on unpack and validate")
	assert.Equal(t, testStruct{Name: "apple", Age: 2}, s)
}

func TestWrappedUnpackAndValidateWithValidation(t *testing.T) {
	validations := []model.Validation{{Key: "name", Type: model.String, Requirement: "equapple"}}
	mapToUpdate := map[string]any{}
	err := UnpackAndValidateWithValidation(newBodyRequest(t, "application/x-www-form-urlencoded", "name=apple"), &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unpack and validate with validation")
	assert.Equal(t, map[string]any{"name": "apple"}, mapToUpdate)
}

func TestWrappedUnmapAndValidate(t *testing.T) {
	values := url.Values{}
	values.Set("name", "apple")
	values.Set("age", "2")
	var s testStruct
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err, "Expected no error creating request")

	err = UnmapAndValidate(req, &s)
	assert.NoError(t, err, "Expected no error on unmap validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmarshalAndValidate(t *testing.T) {
	data := []byte(`{"name":"apple","age":2}`)
	var s testStruct
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err, "Expected no error creating request")

	err = UnmarshalAndValidate(req, &s)
	assert.NoError(t, err, "Expected no error on unmarshal validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmapAndValidateWithValidation(t *testing.T) {
	values := url.Values{}
	values.Set("name", "apple")
	values.Set("age", "2")
	mapToUpdate := map[string]any{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	req, err := http.NewRequest("POST", "/", bytes.NewBufferString(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(t, err, "Expected no error creating request")

	err = UnmapAndValidateWithValidation(req, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmap validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, float64(2), mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnmarshalAndValidateWithValidation(t *testing.T) {
	data := []byte(`{"name":"apple","age":2}`)
	mapToUpdate := map[string]any{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}
	req, err := http.NewRequest("POST", "/", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	require.NoError(t, err, "Expected no error creating request")

	err = UnmarshalAndValidateWithValidation(req, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmarshal validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, float64(2), mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnmarshalYamlAndValidate(t *testing.T) {
	data := []byte("name: apple\nage: 2\n")
	var s testStruct

	err := UnmarshalYamlAndValidate(data, &s)
	assert.NoError(t, err, "Expected no error on unmarshal yaml validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmarshalYamlAndValidateWithValidation(t *testing.T) {
	data := []byte("name: apple\nage: 2\n")
	mapToUpdate := map[string]any{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := UnmarshalYamlAndValidateWithValidation(data, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmarshal yaml validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, 2, mapToUpdate["age"], "Expected age to be 2")
}

type wrappedInjectedStruct struct {
	Owner string `json:"owner" yaml:"owner" vld:"equ${user_id}"`
}

var wrappedInjectedValidations = []model.Validation{{Key: "owner", Type: model.String, Requirement: "equ${user_id}"}}

func TestWrappedValidateContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	assert.NoError(t, ValidateContext(ctx, &wrappedInjectedStruct{Owner: "u1"}))
	assert.Error(t, ValidateContext(ctx, &wrappedInjectedStruct{Owner: "u2"}))
}

func TestWrappedValidateWithValidationContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	result, err := ValidateWithValidationContext(ctx, map[string]any{"owner": "u1"}, wrappedInjectedValidations)
	require.NoError(t, err)
	assert.Equal(t, "u1", result["owner"])
}

func TestWrappedValidateAndUpdateContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	s := &wrappedInjectedStruct{}
	require.NoError(t, ValidateAndUpdateContext(ctx, map[string]any{"owner": "u1"}, s))
	assert.Equal(t, "u1", s.Owner)
}

func TestWrappedValidateAndUpdateWithValidationContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	result := map[string]any{}
	require.NoError(t, ValidateAndUpdateWithValidationContext(ctx, map[string]any{"owner": "u1"}, &result, wrappedInjectedValidations))
	assert.Equal(t, "u1", result["owner"])
}

func TestWrappedUnmarshalYamlAndValidateContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	s := &wrappedInjectedStruct{}
	require.NoError(t, UnmarshalYamlAndValidateContext(ctx, []byte("owner: u1\n"), s))
	assert.Equal(t, "u1", s.Owner)
}

func TestWrappedUnmarshalYamlAndValidateWithValidationContext(t *testing.T) {
	ctx := Inject(context.Background(), "user_id", "u1")
	result := map[string]any{}
	require.NoError(t, UnmarshalYamlAndValidateWithValidationContext(ctx, []byte("owner: u1\n"), &result, wrappedInjectedValidations))
	assert.Equal(t, "u1", result["owner"])
}

package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestWrappedValidate(t *testing.T) {
	s := testStruct{Name: "apple", Age: 2}
	err := Validate(&s)
	assert.NoError(t, err, "Expected no error on validate")
}

func TestWrappedValidateWithValidation(t *testing.T) {
	jsonInput := model.JsonMap{"name": "apple", "age": 2}
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
	jsonInput := model.JsonMap{"name": "apple", "age": 2}
	var s testStruct

	err := ValidateAndUpdate(jsonInput, &s)
	assert.NoError(t, err, "Expected no error on validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedValidateAndUpdateWithValidation(t *testing.T) {
	jsonInput := model.JsonMap{"name": "apple", "age": 2}
	mapToUpdate := model.JsonMap{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := ValidateAndUpdateWithValidation(jsonInput, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, 2, mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnmapOrUnmarshalAndValidate(t *testing.T) {
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "2")
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var s testStruct

	err := UnmapOrUnmarshalAndValidate(req, &s)
	assert.NoError(t, err, "Expected no error on unmap or unmarshal and validate")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmapAndValidate(t *testing.T) {
	values := url.Values{}
	values.Set("name", "apple")
	values.Set("age", "2")
	var s testStruct

	err := UnmapAndValidate(values, &s)
	assert.NoError(t, err, "Expected no error on unmap and validate")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmarshalAndValidate(t *testing.T) {
	data := []byte(`{"name":"apple","age":2}`)
	var s testStruct

	err := UnmarshalAndValidate(data, &s)
	assert.NoError(t, err, "Expected no error on unmarshal and validate")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmapOrUnmarshalValidateAndUpdate(t *testing.T) {
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "2")
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	var s testStruct

	err := UnmapOrUnmarshalValidateAndUpdate(req, &s)
	assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmapValidateAndUpdate(t *testing.T) {
	values := url.Values{}
	values.Set("name", "apple")
	values.Set("age", "2")
	var s testStruct

	err := UnmapValidateAndUpdate(values, &s)
	assert.NoError(t, err, "Expected no error on unmap validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmarshalValidateAndUpdate(t *testing.T) {
	data := []byte(`{"name":"apple","age":2}`)
	var s testStruct

	err := UnmarshalValidateAndUpdate(data, &s)
	assert.NoError(t, err, "Expected no error on unmarshal validate and update")
	assert.Equal(t, "apple", s.Name, "Expected name to be 'apple'")
	assert.Equal(t, 2, s.Age, "Expected age to be 2")
}

func TestWrappedUnmapOrUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "2")
	req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	mapToUpdate := model.JsonMap{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := UnmapOrUnmarshalValidateAndUpdateWithValidation(req, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmap or unmarshal validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, float64(2), mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnmapValidateAndUpdateWithValidation(t *testing.T) {
	values := url.Values{}
	values.Set("name", "apple")
	values.Set("age", "2")
	mapToUpdate := model.JsonMap{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := UnmapValidateAndUpdateWithValidation(values, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmap validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, float64(2), mapToUpdate["age"], "Expected age to be 2")
}

func TestWrappedUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	data := []byte(`{"name":"apple","age":2}`)
	mapToUpdate := model.JsonMap{}
	validations := []model.Validation{
		{Key: "name", Requirement: "equapple"},
		{Key: "age", Requirement: "min2"},
	}

	err := UnmarshalValidateAndUpdateWithValidation(data, &mapToUpdate, validations)
	assert.NoError(t, err, "Expected no error on unmarshal validate and update with validation")
	assert.Equal(t, "apple", mapToUpdate["name"], "Expected name to be 'apple'")
	assert.Equal(t, float64(2), mapToUpdate["age"], "Expected age to be 2")
}

package validator

import (
	"bytes"
	"net/http"
	"net/url"
	"testing"

	"github.com/siherrmann/validator/model"
)

type testStruct struct {
	Name string `json:"name" vld:"equapple"`
	Age  int    `json:"age" vld:"min18"`
}

func newValidator() *Validator {
	return &Validator{}
}

func TestUnmapOrUnmarshalAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "30")
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
	req.Form = form

	err := v.UnmapOrUnmarshalAndValidate(req, ts)
	if err != nil {
		t.Errorf("UnmapOrUnmarshalAndValidate failed: %v", err)
	}
}

func TestUnmarshalAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	jsonInput := []byte(`{"name":"apple","age":30}`)

	err := v.UnmarshalAndValidate(jsonInput, ts)
	if err != nil {
		t.Errorf("UnmarshalAndValidate failed: %v", err)
	}
}

func TestUnmapAndValidate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "30")

	err := v.UnmapAndValidate(form, ts)
	if err != nil {
		t.Errorf("UnmapAndValidate failed: %v", err)
	}
}

func TestUnmapOrUnmarshalValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "30")
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
	req.Form = form

	err := v.UnmapOrUnmarshalValidateAndUpdate(req, ts)
	if err != nil {
		t.Errorf("UnmapOrUnmarshalValidateAndUpdate failed: %v", err)
	}
}

func TestUnmarshalValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	jsonInput := []byte(`{"name":"apple","age":30}`)

	err := v.UnmarshalValidateAndUpdate(jsonInput, ts)
	if err != nil {
		t.Errorf("UnmarshalValidateAndUpdate failed: %v", err)
	}
}

func TestUnmapValidateAndUpdate(t *testing.T) {
	v := newValidator()
	ts := &testStruct{}
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "30")

	err := v.UnmapValidateAndUpdate(form, ts)
	if err != nil {
		t.Errorf("UnmapValidateAndUpdate failed: %v", err)
	}
}

func TestUnmapOrUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	validations := []model.Validation{
		{
			Key:         "name",
			Type:        model.String,
			Requirement: "equapple",
		},
		{
			Key:         "age",
			Type:        model.Int,
			Requirement: "min18",
		},
	}

	t.Run("With form values", func(t *testing.T) {
		mapOut := &model.JsonMap{}
		form := url.Values{}
		form.Set("name", "apple")
		form.Set("age", "30")
		req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(nil))
		req.Form = form

		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapOut, validations)
		if err != nil {
			t.Errorf("UnmapOrUnmarshalValidateAndUpdateWithValidation failed: %v", err)
		}
	})

	t.Run("With request body", func(t *testing.T) {
		mapOut := &model.JsonMap{}
		req, _ := http.NewRequest("POST", "/", bytes.NewBufferString(`{"name":"apple","age":30}`))
		err := v.UnmapOrUnmarshalValidateAndUpdateWithValidation(req, mapOut, validations)
		if err != nil {
			t.Errorf("UnmapOrUnmarshalValidateAndUpdateWithValidation failed: %v", err)
		}
	})
}

func TestUnmarshalValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	mapOut := &model.JsonMap{}
	jsonInput := []byte(`{"name":"apple","age":30}`)
	validations := []model.Validation{}

	err := v.UnmarshalValidateAndUpdateWithValidation(jsonInput, mapOut, validations)
	if err != nil {
		t.Errorf("UnmarshalValidateAndUpdateWithValidation failed: %v", err)
	}
}

func TestUnmapValidateAndUpdateWithValidation(t *testing.T) {
	v := newValidator()
	mapOut := &model.JsonMap{}
	form := url.Values{}
	form.Set("name", "apple")
	form.Set("age", "30")
	validations := []model.Validation{}

	err := v.UnmapValidateAndUpdateWithValidation(form, mapOut, validations)
	if err != nil {
		t.Errorf("UnmapValidateAndUpdateWithValidation failed: %v", err)
	}
}

package validator

import (
	"encoding/json"
	"regexp"
	"strconv"
	"testing"
)

type TestStructNoValidation struct {
	String string   `vld:"-"`
	Int    int      `vld:"-"`
	Float  float64  `vld:"-"`
	Array  []string `vld:"-"`
}

type TestStructValidationRex struct {
	String string   `vld:"rex^[a-zA-Z0-9]+$, gr1min3"`
	Int    int      `vld:"rex^(2|3)$, gr1min3"`
	Float  float64  `vld:"rex^(2.000000|3.000000)$, gr1min3"`
	Array  []string `vld:"min3, gr1min3"`
}

type TestStructValidation struct {
	String string   `json:"string" vld:"rex^[a-zA-Z0-9]+$, gr1min3"`
	Int    int      `json:"int" vld:"equ2 || equ3, gr1min3"`
	Float  float64  `json:"float" vld:"equ2 || equ3, gr1min3"`
	Array  []string `json:"array" vld:"min3, gr1min3"`
}

func BenchmarkUnmarshalToMap(b *testing.B) {
	// unmarshal to map
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array":  ["", "", ""]}`)
	unmarshalToMap := &map[string]interface{}{}
	err := json.Unmarshal(jsonString, unmarshalToMap)
	if err != nil {
		b.Logf("error unmarshalling %s: %v", jsonString, err)
	}
}

func BenchmarkUnmarshalToStruct(b *testing.B) {
	// unmarshal to struct
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	unmarshalToStruct := &TestStructNoValidation{}
	err := json.Unmarshal(jsonString, unmarshalToStruct)
	if err != nil {
		b.Logf("error unmarshalling %s: %v", jsonString, err)
	}
}

func BenchmarkNoValidation(b *testing.B) {
	// no validation in validator
	noValidation := &TestStructNoValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	err := Validate(noValidation)
	if err != nil {
		b.Logf("error no validation %v", err)
	}
}

func BenchmarkValidationRex(b *testing.B) {
	// validation with regex in validator
	validationRex := &TestStructValidationRex{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	err := Validate(validationRex)
	if err != nil {
		b.Logf("error validation rex %v", err)
	}
}

func BenchmarkManualValidationRex(b *testing.B) {
	// manual validation with regex
	var errors int
	manualValidationRex := &TestStructValidationRex{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9]+$", manualValidationRex.String)
	if err != nil || !match {
		errors++
	}
	match, err = regexp.MatchString("^(2|3)$", strconv.Itoa(manualValidationRex.Int))
	if err != nil || !match {
		errors++
	}
	match, err = regexp.MatchString("^(2.000000|3.000000)$", strconv.FormatFloat(manualValidationRex.Float, 'f', 0, 64))
	if err != nil || !match {
		errors++
	}
	if len(manualValidationRex.Array) < 3 {
		errors++
	}
	if errors > 1 {
		b.Logf("error manual validation rex %v", err)
	}
}

func BenchmarkValidatationMinimalRex(b *testing.B) {
	// validation with minimal regex (only string) in validator
	validation := &TestStructValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	err := Validate(validation)
	if err != nil {
		b.Logf("error validation %v", err)
	}
}

func BenchmarkManualValidationMinimalRex(b *testing.B) {
	// manual validation with minimal regex (only string)
	var errors int
	manualValidation := &TestStructValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	match, err := regexp.MatchString("^[a-zA-Z0-9]+$", manualValidation.String)
	if err != nil || !match {
		errors++
	}
	if manualValidation.Int != 2 && manualValidation.Int != 3 {
		errors++
	}
	if manualValidation.Float != 2 && manualValidation.Float != 3 {
		errors++
	}
	if len(manualValidation.Array) < 3 {
		errors++
	}
	if errors > 1 {
		b.Logf("error manual validation %v", err)
	}
}

func BenchmarkValidationWithNoValidation(b *testing.B) {
	// no validation in validator
	noValidation := &TestStructNoValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	err := Validate(noValidation)
	if err != nil {
		b.Logf("error no validation %v", err)
	}
}

func BenchmarkUnmarshalAndValidate(b *testing.B) {
	// unmarshal and validate
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	unmarshalAndValidate := &TestStructValidation{}
	err := UnmarshalAndValidate(jsonString, unmarshalAndValidate)
	if err != nil {
		b.Logf("error unmarshal and validate %v", err)
	}

	b.Error("test ended")
}

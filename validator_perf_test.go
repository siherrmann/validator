package validator

import (
	"encoding/json"
	"regexp"
	"testing"
)

type TestStructNoValidation struct {
	String string   `vld:"-"`
	Int    int      `vld:"-"`
	Float  float64  `vld:"-"`
	Array  []string `vld:"-"`
}

type TestStructValidation struct {
	String string   `json:"string" vld:"rex^[a-zA-Z0-9]+$, gr1min3"`
	Int    int      `json:"int" vld:"equ2 || equ3, gr1min3"`
	Float  float64  `json:"float" vld:"equ2 || equ3, gr1min3"`
	Array  []string `json:"array" vld:"min3, gr1min3"`
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

func BenchmarkUnmarshalAndManualValidation(b *testing.B) {
	// Unmarshal JSON string to struct
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	manualValidation := &TestStructValidation{}
	err := json.Unmarshal(jsonString, manualValidation)
	if err != nil {
		b.Fatalf("error unmarshal %v", err)
	}

	// manual validation with minimal regex (only string)
	var errors int
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

func BenchmarkUnmarshalAndValidate(b *testing.B) {
	// unmarshal and validate
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	unmarshalAndValidate := &TestStructValidation{}
	err := UnmarshalAndValidate(jsonString, unmarshalAndValidate)
	if err != nil {
		b.Logf("error unmarshal and validate %v", err)
	}
}

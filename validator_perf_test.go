package validator

import (
	"encoding/json"
	"regexp"
	"strconv"
	"testing"
	"time"
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

func TestPerformanceStructValidator(t *testing.T) {
	// unmarshal to map
	jsonString := []byte(`{"string": "test", "int": 2, "float": 3.0, "array":  ["", "", ""]}`)
	unmarshalToMap := &map[string]interface{}{}
	start := time.Now()
	err := json.Unmarshal(jsonString, unmarshalToMap)
	if err != nil {
		t.Logf("error unmarshalling %s: %v", jsonString, err)
	}
	elapsed := time.Since(start)
	t.Logf("unmarshal to map took %s", elapsed)

	// unmarshal to struct
	jsonString = []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	unmarshalToStruct := &TestStructNoValidation{}
	start = time.Now()
	err = json.Unmarshal(jsonString, unmarshalToStruct)
	if err != nil {
		t.Logf("error unmarshalling %s: %v", jsonString, err)
	}
	elapsed = time.Since(start)
	t.Logf("unmarshal to struct took %s", elapsed)

	// unmarshal to map
	jsonString = []byte(`{"string": "test", "int": 2, "float": 3.0, "array":  ["", "", ""]}`)
	unmarshalToMap = &map[string]interface{}{}
	start = time.Now()
	err = json.Unmarshal(jsonString, unmarshalToMap)
	if err != nil {
		t.Logf("error unmarshalling %s: %v", jsonString, err)
	}
	elapsed = time.Since(start)
	t.Logf("unmarshal to map took %s", elapsed)

	// no validation in validator
	noValidation := &TestStructNoValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
	err = Validate(noValidation)
	if err != nil {
		t.Logf("error no validation %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("no validation took %s", elapsed)

	// validation with regex in validator
	validationRex := &TestStructValidationRex{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
	err = Validate(validationRex)
	if err != nil {
		t.Logf("error validation rex %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("validation rex took %s", elapsed)

	// manual validation with regex
	var errors int
	manualValidationRex := &TestStructValidationRex{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
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
		t.Logf("error manual validation rex %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("manual validation rex took %s", elapsed)

	// validation with minimal regex (only string) in validator
	validation := &TestStructValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
	err = Validate(validation)
	if err != nil {
		t.Logf("error validation %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("validation took %s", elapsed)

	// manual validation with minimal regex (only string)
	manualValidation := &TestStructValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
	match, err = regexp.MatchString("^[a-zA-Z0-9]+$", manualValidation.String)
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
		t.Logf("error manual validation %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("manual validation took %s", elapsed)

	// no validation in validator
	noValidation = &TestStructNoValidation{
		String: "test",
		Int:    2,
		Float:  3.0,
		Array:  []string{"", "", ""},
	}
	start = time.Now()
	err = Validate(noValidation)
	if err != nil {
		t.Logf("error no validation %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("no validation took %s", elapsed)

	// unmarshal and validate
	jsonString = []byte(`{"string": "test", "int": 2, "float": 3.0, "array": ["", "", ""]}`)
	unmarshalAndValidate := &TestStructValidation{}
	start = time.Now()
	err = UnmarshalAndValidate(jsonString, unmarshalAndValidate)
	if err != nil {
		t.Logf("error unmarshal and validate %v", err)
	}
	elapsed = time.Since(start)
	t.Logf("unmarshal and validate took %s with %v", elapsed, unmarshalAndValidate)

	t.Error("test ended")
}

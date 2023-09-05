package validator

import (
	"strings"
	"testing"
)

type TestRequestWrapper struct {
	Data         interface{}
	InvalidField string
}

type TestStructEqual struct {
	String string   `vld:"equtest"`
	Int    int      `vld:"equ3"`
	Float  float64  `vld:"equ3"`
	Array  []string `vld:"equ3"`
}

type TestStructNotEqual struct {
	String string   `vld:"neqtest"`
	Int    int      `vld:"neq3"`
	Float  float64  `vld:"neq3"`
	Array  []string `vld:"neq3"`
}

type TestStructMin struct {
	String string   `vld:"min3"`
	Int    int      `vld:"min3"`
	Float  float64  `vld:"min3"`
	Array  []string `vld:"min3"`
}

type TestStructMax struct {
	String string   `vld:"max4"`
	Int    int      `vld:"max4"`
	Float  float64  `vld:"max4"`
	Array  []string `vld:"max4"`
}

type TestStructCon struct {
	String string   `vld:"con@"`
	Int    int      `vld:"-"`
	Float  float64  `vld:"-"`
	Array  []string `vld:"con@"`
}

type TestStructRex struct {
	String string   `vld:"rex^[a-zA-Z0-9]+$"`
	Int    int      `vld:"rex^(2|3)$"`
	Float  float64  `vld:"rex^(2.000000|3.000000)$"`
	Array  []string `vld:"-"`
}

type TestStructMulti struct {
	String string   `vld:"min3 neqtest"`
	Int    int      `vld:"min3 neq4"`
	Float  float64  `vld:"min3 neq4"`
	Array  []string `vld:"min3 max4"`
}

type TestStructEmptyCondition struct {
	String string   `vld:"equ"`
	Int    int      `vld:"neq"`
	Float  float64  `vld:"min"`
	Array  []string `vld:"max"`
}

type TestStructPassword struct {
	String string `vld:"min8 max30 rex^(.*[A-Z])+(.*)$ rex^(.*[a-z])+(.*)$ rex^(.*\\d)+(.*)$ rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$"`
}

type TestStructGroup struct {
	String string   `vld:"min3, gr1min1 gr2min2"`
	Int    int      `vld:"min3, gr1min1 gr2min2"`
	Float  float64  `vld:"min3, gr1min1"`
	Array  []string `vld:"min3, gr1min1"`
}

type TestStructGroupNoGroup struct {
	String string   `vld:"min3"`
	Int    int      `vld:"min3, gr3min1"`
	Float  float64  `vld:"min3, gr3min1"`
	Array  []string `vld:"min3, gr3min1"`
}

type TestStructInvalidGroupCondition struct {
	String string   `vld:"min3, gr4min1 gr5min2"`
	Int    int      `vld:"min3, gr4min1 gr5min2"`
	Float  float64  `vld:"min3, gr4min1"`
	Array  []string `vld:"min3, gr40min1"`
}

func TestStructValidator(t *testing.T) {
	testCases := map[string]*TestRequestWrapper{
		"valid": {
			&TestStructEqual{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"notEqualString": {
			&TestStructEqual{
				String: "tes",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"String",
		},
		"notEqualInt": {
			&TestStructEqual{
				String: "test",
				Int:    -3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"Int",
		},
		"notEqualFloat": {
			&TestStructEqual{
				String: "test",
				Int:    3,
				Float:  -3.0,
				Array:  []string{"", "", ""},
			},
			"Float",
		},
		"notEqualArray": {
			&TestStructEqual{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructNotEqual{
				String: "tes",
				Int:    -3,
				Float:  -3.0,
				Array:  []string{"", ""},
			},
			"",
		},
		"equalString": {
			&TestStructNotEqual{
				String: "test",
				Int:    -3,
				Float:  -3.0,
				Array:  []string{"", ""},
			},
			"String",
		},
		"equalInt": {
			&TestStructNotEqual{
				String: "tes",
				Int:    3,
				Float:  -3.0,
				Array:  []string{"", ""},
			},
			"Int",
		},
		"equalFloat": {
			&TestStructNotEqual{
				String: "tes",
				Int:    -3,
				Float:  3.0,
				Array:  []string{"", ""},
			},
			"Float",
		},
		"equalArray": {
			&TestStructNotEqual{
				String: "tes",
				Int:    -3,
				Float:  -3.0,
				Array:  []string{"", "", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructMin{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"tooShortString": {
			&TestStructMin{
				String: "",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"String",
		},
		"tooShortInt": {
			&TestStructMin{
				String: "test",
				Int:    -3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"Int",
		},
		"tooShortFloat": {
			&TestStructMin{
				String: "test",
				Int:    3,
				Float:  -3.0,
				Array:  []string{"", "", ""},
			},
			"Float",
		},
		"tooShortArray": {
			&TestStructMin{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructMax{
				String: "test",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "", "", ""},
			},
			"",
		},
		"tooLongString": {
			&TestStructMax{
				String: "testi",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "", "", ""},
			},
			"String",
		},
		"tooLongInt": {
			&TestStructMax{
				String: "test",
				Int:    5,
				Float:  4.0,
				Array:  []string{"", "", "", ""},
			},
			"Int",
		},
		"tooLongFloat": {
			&TestStructMax{
				String: "test",
				Int:    4,
				Float:  4.1,
				Array:  []string{"", "", "", ""},
			},
			"Float",
		},
		"tooLongArray": {
			&TestStructMax{
				String: "test",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "", "", "", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructCon{
				String: "test@",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "@"},
			},
			"",
		},
		"notContainingString": {
			&TestStructCon{
				String: "test",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "@"},
			},
			"String",
		},
		"notContainingArray": {
			&TestStructCon{
				String: "test@",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructRex{
				String: "test",
				Int:    2,
				Float:  2.0,
				Array:  []string{"", "", "", ""},
			},
			"",
		},
		"notMatchingString": {
			&TestStructRex{
				String: "test@",
				Int:    2,
				Float:  2.0,
				Array:  []string{"", "", "", ""},
			},
			"String",
		},
		"notMatchingInt": {
			&TestStructRex{
				String: "test",
				Int:    -2,
				Float:  2.0,
				Array:  []string{"", "", "", ""},
			},
			"Int",
		},
		"notMatchingFloat": {
			&TestStructRex{
				String: "test",
				Int:    2,
				Float:  -2.0,
				Array:  []string{"", "", "", ""},
			},
			"Float",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructMulti{
				String: "tes",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"firstFailString": {
			&TestStructMulti{
				String: "te",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"String",
		},
		"firstFailInt": {
			&TestStructMulti{
				String: "tes",
				Int:    2,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"Int",
		},
		"firstFailFloat": {
			&TestStructMulti{
				String: "tes",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
			},
			"Float",
		},
		"firstFailArray": {
			&TestStructMulti{
				String: "tes",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", ""},
			},
			"Array",
		},
		"secondFailString": {
			&TestStructMulti{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"String",
		},
		"secondFailInt": {
			&TestStructMulti{
				String: "tes",
				Int:    4,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"Int",
		},
		"secondFailFloat": {
			&TestStructMulti{
				String: "tes",
				Int:    3,
				Float:  4.0,
				Array:  []string{"", "", ""},
			},
			"Float",
		},
		"secondFailArray": {
			&TestStructMulti{
				String: "tes",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", "", "", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"emptyConditon": {
			&TestStructEmptyCondition{
				String: "test@",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", "", "", ""},
			},
			"String",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"valid": {
			&TestStructPassword{
				String: "Password123.4",
			},
			"",
		},
		"tooShort": {
			&TestStructPassword{
				String: "Pa123.4",
			},
			"String",
		},
		"tooLong": {
			&TestStructPassword{
				String: "Password1.Password1.Password1.2",
			},
			"String",
		},
		"missingCapitalLetter": {
			&TestStructPassword{
				String: "password123.4",
			},
			"String",
		},
		"missingNonCapitalLetter": {
			&TestStructPassword{
				String: "PASSWORD123.4",
			},
			"String",
		},
		"missingDecimalLetter": {
			&TestStructPassword{
				String: "Password.",
			},
			"String",
		},
		"missingSpecialCharacter": {
			&TestStructPassword{
				String: "Password1234",
			},
			"String",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"validAll": {
			&TestStructGroup{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"validOnlyGroup2": {
			&TestStructGroup{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", ""},
			},
			"",
		},
		"noneOfGroup": {
			&TestStructGroup{
				String: "te",
				Int:    2,
				Float:  2.0,
				Array:  []string{"", ""},
			},
			"group",
		},
		"onlyGroup1": {
			&TestStructGroup{
				String: "te",
				Int:    2,
				Float:  3.0,
				Array:  []string{"", ""},
			},
			"group",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"validAll": {
			&TestStructGroupNoGroup{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"validOnlyOneOfGroup": {
			&TestStructGroupNoGroup{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", ""},
			},
			"",
		},
		"onlyGroup": {
			&TestStructGroupNoGroup{
				String: "te",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", ""},
			},
			"String",
		},
		"onlyNoGroup": {
			&TestStructGroupNoGroup{
				String: "test",
				Int:    2,
				Float:  2.0,
				Array:  []string{"", ""},
			},
			"group",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}

	testCases = map[string]*TestRequestWrapper{
		"invalidGroupConditionLast": {
			&TestStructInvalidGroupCondition{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"group",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func assertError(t testing.TB, testCase string, err error, invalidField string) {
	t.Helper()
	if len(invalidField) == 0 && err != nil {
		t.Errorf("test case: %s - wanted no invalid field, got error: %v", testCase, err)
	} else if len(invalidField) != 0 && err == nil {
		t.Errorf("test case: %s - wanted invalid field: %v, got no error", testCase, invalidField)
	} else if err != nil && !strings.Contains(err.Error(), invalidField) {
		t.Errorf("test case: %s - wanted invalid field: %v, got error: %v", testCase, invalidField, err)
	}
}

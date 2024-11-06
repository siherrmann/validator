package validator

import (
	"log"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/siherrmann/validator/model"
	"github.com/siherrmann/validator/parser"
)

type TestRequestWrapper struct {
	Data         interface{}
	InvalidField string
}

type TestRequestWrapperUpdate struct {
	Data          interface{}
	JsonMapUpdate model.JsonMap
	Error         bool
}

type TestRequestWrapperUpdateWithJson struct {
	Data       interface{}
	JsonUpdate string
	Error      bool
}

type TestRequestWrapperUpdateWithUrlValues struct {
	Data            interface{}
	UrlValuesUpdate url.Values
	Error           bool
}

func TestCaseStructIntTypes(t *testing.T) {
	type TestStructIntTypes struct {
		Int   int   `upd:"int, min3"`
		Int64 int64 `upd:"int64, min3"`
		Int32 int32 `upd:"int32, min3"`
		Int16 int16 `upd:"int16, min3"`
		Int8  int8  `upd:"int8, min3"`
	}

	testCases := map[string]*TestRequestWrapperUpdateWithJson{
		"valid": {
			&TestStructIntTypes{
				Int:   3,
				Int64: 3,
				Int32: 3,
				Int16: 3,
				Int8:  3,
			},
			`{"int": 9223372036854775807, "int64": 9223372036854775807, "int32": 2147483647, "int16": 32767, "int8": 127}`,
			false,
		},
		"invalidInt32": {
			&TestStructIntTypes{
				Int:   3,
				Int64: 3,
				Int32: 3,
				Int16: 3,
				Int8:  3,
			},
			`{"int": 9223372036854775807, "int64": 9223372036854775807, "int32": 2147483648, "int16": 32767, "int8": 127}`,
			true,
		},
		"invalidInt16": {
			&TestStructIntTypes{
				Int:   3,
				Int64: 3,
				Int32: 3,
				Int16: 3,
				Int8:  3,
			},
			`{"int": 9223372036854775807, "int64": 9223372036854775807, "int32": 2147483647, "int16": 32768, "int8": 127}`,
			true,
		},
		"invalidInt8": {
			&TestStructIntTypes{
				Int:   3,
				Int64: 3,
				Int32: 3,
				Int16: 3,
				Int8:  3,
			},
			`{"int": 9223372036854775807, "int64": 9223372036854775807, "int32": 2147483647, "int16": 32767, "int8": 128}`,
			true,
		},
	}

	for k, v := range testCases {
		err := UnmarshalValidateAndUpdate([]byte(v.JsonUpdate), v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseStructFloatTypes(t *testing.T) {
	type TestStructFloatTypes struct {
		Float64 float64 `vld:"min3"`
		Float32 float32 `vld:"min3"`
	}

	// As the json Unmarshal always unmarshals into float64
	// updating float32 will only fail with a float64 overflow.
	testCases := map[string]*TestRequestWrapperUpdateWithJson{
		"valid": {
			&TestStructFloatTypes{
				Float64: 3,
				Float32: 3,
			},
			`{"float64": 1.79769313486231570814527423731704356798070e+308, "float32": 3.40282346638528859811704183484516925440e+38}`,
			false,
		},
		"invalidFloat64": {
			&TestStructFloatTypes{
				Float64: 3,
				Float32: 3,
			},
			`{"float64": 1.79769313486231570814527423731704356798070e+309, "float32": 3.40282346638528859811704183484516925440e+38}`,
			true,
		},
		"invalidFloat32": {
			&TestStructFloatTypes{
				Float64: 3,
				Float32: 3,
			},
			`{"float64": 1.79769313486231570814527423731704356798070e+308, "float32": 3.40282346638528859811704183484516925440e+309}`,
			true,
		},
	}

	for k, v := range testCases {
		err := UnmarshalValidateAndUpdate([]byte(v.JsonUpdate), v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseStructEqual(t *testing.T) {
	type TestStructEqual struct {
		String string   `vld:"equtest"`
		Int    int      `vld:"equ3"`
		Float  float64  `vld:"equ3"`
		Array  []string `vld:"equ3"`
	}

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
}

func TestCaseStructNotEqual(t *testing.T) {
	type TestStructNotEqual struct {
		String string   `vld:"neqtest"`
		Int    int      `vld:"neq3"`
		Float  float64  `vld:"neq3"`
		Array  []string `vld:"neq3"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructMin(t *testing.T) {
	type TestStructMin struct {
		String string   `vld:"min3"`
		Int    int      `vld:"min3"`
		Float  float64  `vld:"min3"`
		Array  []string `vld:"min3"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructMax(t *testing.T) {
	type TestStructMax struct {
		String string   `vld:"max4"`
		Int    int      `vld:"max4"`
		Float  float64  `vld:"max4"`
		Array  []string `vld:"max4"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructCon(t *testing.T) {
	type TestStructCon struct {
		String string   `vld:"con@"`
		Int    int      `vld:"-"`
		Float  float64  `vld:"-"`
		Array  []string `vld:"con@"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructNotCon(t *testing.T) {
	type TestStructNotCon struct {
		String string   `vld:"nco@"`
		Int    int      `vld:"-"`
		Float  float64  `vld:"-"`
		Array  []string `vld:"nco@"`
	}

	testCases := map[string]*TestRequestWrapper{
		"valid": {
			&TestStructNotCon{
				String: "test",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", ""},
			},
			"",
		},
		"containingString": {
			&TestStructNotCon{
				String: "test@",
				Int:    4,
				Float:  4.0,
				Array:  []string{"", ""},
			},
			"String",
		},
		"containingArray": {
			&TestStructNotCon{
				String: "test",
				Int:    4,
				Float:  4.0,
				Array:  []string{"@", ""},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func TestCaseStructFrom(t *testing.T) {
	type TestStructFrom struct {
		String string   `vld:"frm@,$"`
		Int    int      `vld:"frm1,200"`
		Float  float64  `vld:"frm0.5,1"`
		Array  []string `vld:"frm@,$"`
	}

	testCases := map[string]*TestRequestWrapper{
		"valid": {
			&TestStructFrom{
				String: "@",
				Int:    1,
				Float:  0.5,
				Array:  []string{"@", "@", "@", "@"},
			},
			"",
		},
		"valid2": {
			&TestStructFrom{
				String: "$",
				Int:    200,
				Float:  1,
				Array:  []string{"$", "$", "$", "@"},
			},
			"",
		},
		"notFromString": {
			&TestStructFrom{
				String: "test",
				Int:    200,
				Float:  1,
				Array:  []string{"$", "$", "$", "@"},
			},
			"String",
		},
		"notFromInt": {
			&TestStructFrom{
				String: "@",
				Int:    4,
				Float:  1,
				Array:  []string{"$", "$", "$", "@"},
			},
			"Int",
		},
		"notFromFloat": {
			&TestStructFrom{
				String: "@",
				Int:    200,
				Float:  1.5,
				Array:  []string{"$", "$", "$", "@"},
			},
			"Float",
		},
		"notFromArray": {
			&TestStructFrom{
				String: "$",
				Int:    200,
				Float:  1,
				Array:  []string{"test", "@", "$", "@"},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func TestCaseStructNotFrom(t *testing.T) {
	type TestStructNotFrom struct {
		String string   `vld:"nfr@,$"`
		Int    int      `vld:"nfr1,2"`
		Float  float64  `vld:"nfr0.5,0.8"`
		Array  []string `vld:"nfr@,$"`
	}

	testCases := map[string]*TestRequestWrapper{
		"valid": {
			&TestStructNotFrom{
				String: "test",
				Int:    4,
				Float:  2,
				Array:  []string{"test", "test", "test", "test"},
			},
			"",
		},
		"fromString": {
			&TestStructNotFrom{
				String: "@",
				Int:    4,
				Float:  2,
				Array:  []string{"test", "test", "test", "test"},
			},
			"String",
		},
		"fromInt": {
			&TestStructNotFrom{
				String: "test",
				Int:    1,
				Float:  2,
				Array:  []string{"test", "test", "test", "test"},
			},
			"Int",
		},
		"fromFloat": {
			&TestStructNotFrom{
				String: "test",
				Int:    4,
				Float:  0.5,
				Array:  []string{"test", "test", "test", "test"},
			},
			"Float",
		},
		"fromArray": {
			&TestStructNotFrom{
				String: "test",
				Int:    4,
				Float:  2,
				Array:  []string{"@", "test", "test", "test"},
			},
			"Array",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func TestCaseStructRex(t *testing.T) {
	type TestStructRex struct {
		String string   `vld:"rex'^[a-zA-Z0-9]+$'"`
		Int    int      `vld:"rex'^(2|3)$'"`
		Float  float64  `vld:"rex'^(2|3)$'"`
		Array  []string `vld:"-"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructMulti(t *testing.T) {
	type TestStructMulti struct {
		String string   `vld:"min3 neqtest"`
		Int    int      `vld:"min3 neq4"`
		Float  float64  `vld:"min3 neq4"`
		Array  []string `vld:"min3 max4"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructMultiOr(t *testing.T) {
	type TestStructMultiOr struct {
		String string   `vld:"min4 || equte"`
		Int    int      `vld:"min3 || equ0"`
		Float  float64  `vld:"min3 || equ0"`
		Array  []string `vld:"min3 || equ0"`
	}

	testCases := map[string]*TestRequestWrapper{
		"validFirstCondition": {
			&TestStructMultiOr{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"",
		},
		"validSecondCondition": {
			&TestStructMultiOr{
				String: "te",
				Int:    0,
				Float:  0.0,
				Array:  []string{},
			},
			"",
		},
		"failString": {
			&TestStructMultiOr{
				String: "ts",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"String",
		},
		"failInt": {
			&TestStructMultiOr{
				String: "test",
				Int:    2,
				Float:  3.0,
				Array:  []string{"", "", ""},
			},
			"Int",
		},
		"failFloat": {
			&TestStructMultiOr{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
			},
			"Float",
		},
		"failArray": {
			&TestStructMultiOr{
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
}

func TestCaseTestStructEmptyCondition(t *testing.T) {
	type TestStructEmptyCondition struct {
		String string   `vld:"equ"`
		Int    int      `vld:"neq"`
		Float  float64  `vld:"min"`
		Array  []string `vld:"max"`
	}

	testCases := map[string]*TestRequestWrapper{
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
}

func TestCaseStructPassword(t *testing.T) {
	type TestStructPassword struct {
		String string `vld:"min8 max30 rex'^(.*[A-Z])+(.*)$' rex'^(.*[a-z])+(.*)$' rex'^(.*\\d)+(.*)$' rex'^(.*[\x60!@#$%^&*()_+={};/':\"|\\,.<>/?~-])+(.*)$'"`
	}

	testCases := map[string]*TestRequestWrapper{
		"valid": {
			&TestStructPassword{
				String: "Password123'4",
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
}

func TestCaseStructGroup(t *testing.T) {
	type TestStructGroup struct {
		String string            `vld:"min3, gr1min1 gr2min2"`
		Int    int               `vld:"min3, gr1min1 gr2min2"`
		Float  float64           `vld:"min3, gr1min1"`
		Array  []string          `vld:"min3, gr1min1"`
		Map    map[string]string `vld:"min3, gr1min1"`
		Struct time.Time         `vld:"min1728000000, gr1min1"`
	}

	testCases := map[string]*TestRequestWrapper{
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

}

func TestCaseStructGroupNoGroup(t *testing.T) {
	type TestStructGroupNoGroup struct {
		String string            `vld:"min3"`
		Int    int               `vld:"min3, gr3min5"`
		Float  float64           `vld:"min3, gr3min5"`
		Array  []string          `vld:"min3, gr3min5"`
		Map    map[string]string `vld:"min3, gr3min5"`
		Struct time.Time         `vld:"min1728000000, gr3min5"`
	}

	testTimeValid, _ := time.Parse("2006-01-02T15:04:05.000000", "2025-01-02T15:04:05.000000")
	testTimeInvalid, _ := time.Parse("2006-01-02T15:04:05.000000", "2024-01-02T15:04:05.000000")
	testCases := map[string]*TestRequestWrapper{
		"validAll": {
			&TestStructGroupNoGroup{
				String: "test",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
				Map:    map[string]string{"i": "j", "k": "l", "m": "n"},
				Struct: testTimeValid,
			},
			"",
		},
		"onlyGroup": {
			&TestStructGroupNoGroup{
				String: "te",
				Int:    3,
				Float:  3.0,
				Array:  []string{"", "", ""},
				Map:    map[string]string{"i": "j", "k": "l", "m": "n"},
				Struct: testTimeValid,
			},
			"String",
		},
		"onlyNoGroup": {
			&TestStructGroupNoGroup{
				String: "test",
				Int:    2,
				Float:  2.0,
				Array:  []string{"", ""},
				Map:    map[string]string{"i": "j", "k": "l"},
				Struct: testTimeInvalid,
			},
			"group",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func TestCaseStructInvalidGroupCondition(t *testing.T) {
	type TestStructInvalidGroupCondition struct {
		String string            `vld:"min3, gr4min1 gr5min2"`
		Int    int               `vld:"min3, gr4min1 gr5min2"`
		Float  float64           `vld:"min3, gr4min1"`
		Array  []string          `vld:"min3, gr4min1"`
		Map    map[string]string `vld:"min3, gr4min1"`
		Struct time.Time         `vld:"min1728000000, gr40min1"`
	}

	testCases := map[string]*TestRequestWrapper{
		"invalidGroupConditionLast": {
			&TestStructInvalidGroupCondition{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
			},
			"condition type",
		},
	}

	for k, v := range testCases {
		err := Validate(v.Data)
		assertError(t, k, err, v.InvalidField)
	}
}

func TestCaseUpdate(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, equtest"`
	}

	type TestUpdate struct {
		String string          `upd:"string, min1, gr1min7"`
		Int    int             `upd:"int, min1, gr1min7"`
		Float  float64         `upd:"float, min1, gr1min7"`
		Array  []int           `upd:"array, min1, gr1min7"`
		Date   time.Time       `upd:"date, min1, gr1min7"`
		Struct TestUpdateInner `upd:"struct, min1, gr1min7"`
		Map    model.JsonMap   `upd:"map, min1 conkey, gr1min7"`
	}

	testCasesUpdate := map[string]*TestRequestWrapperUpdate{
		"validUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Date:   time.Time{},
				Struct: TestUpdateInner{
					String: "foo",
				},
				Map: model.JsonMap{
					"key": "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []interface{}{2}, "date": "2022-01-03T15:04:05.000Z", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			false,
		},
		"invalidJsonStringUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "", "int": 2, "float": "1.2", "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidTypeStringUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": 1, "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidJsonIntUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 0, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		// TODO add strict case/case for json or url values?
		// "invalidTypeIntUpdate": {
		// 	&TestUpdate{
		// 		String: "Foo",
		// 		Int:    1,
		// 		Float:  1.1,
		// 		Array:  []int{1},
		// 		Date:   time.Time{},
		// 	},
		// 	map[string]interface{}{"string": "Bar", "int": "2", "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
		// 	true,
		// },
		"invalidJsonFloatUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 0.0, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		// TODO add strict case/case for json or url values?
		// "invalidTypeFloatUpdate": {
		// 	&TestUpdate{
		// 		String: "Foo",
		// 		Int:    1,
		// 		Float:  1.1,
		// 		Array:  []int{1},
		// 		Date:   time.Time{},
		// 	},
		// 	map[string]interface{}{"string": "Bar", "int": 2, "float": "1.2", "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
		// 	true,
		// },
		"invalidJsonArrayUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidTypeArrayUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []string{"2"}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidJsonDateUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03 15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
		// TODO add case for time.Time with timeUnix in tag
		// "invalidTypeDateUpdate": {
		// 	&TestUpdate{
		// 		String: "Foo",
		// 		Int:    1,
		// 		Float:  1.1,
		// 		Array:  []int{1},
		// 		Date:   time.Time{},
		// 	},
		// 	map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": 2024, "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
		// 	true,
		// },
		"invalidJsonStructUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "testing"}, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidTypeStructUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": 2, "map": map[string]any{"key": "test"}},
			true,
		},
		"invalidJsonMapUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
				Struct: TestUpdateInner{
					String: "foo",
				},
				Map: model.JsonMap{
					"key": "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000Z", "struct": map[string]any{"string": "test"}, "map": map[string]any{"test": "test"}},
			true,
		},
		"invalidTypeMapUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
				Struct: TestUpdateInner{
					String: "foo",
				},
				Map: model.JsonMap{
					"key": "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000Z", "struct": map[string]any{"string": "test"}, "map": map[int]any{2: "test"}},
			true,
		},
	}

	for k, v := range testCasesUpdate {
		err := ValidateAndUpdate(v.JsonMapUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseUpdatePartial(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, equtest"`
	}

	type TestUpdatePartial struct {
		String string `upd:"string, min1, gr1min2"`
		Int    int
		Float  float64
		Array  []int `upd:"array, min1, gr1min2"`
		Date   time.Time
		Struct TestUpdateInner
		Map    model.JsonMap
	}

	testCasesUpdatePartial := map[string]*TestRequestWrapperUpdate{
		"validUpdate": {
			&TestUpdatePartial{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Date:   time.Time{},
				Struct: TestUpdateInner{
					String: "foo",
				},
				Map: model.JsonMap{
					"key": "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "array": []interface{}{2}},
			false,
		},
		"validUpdateMoreValues": {
			&TestUpdatePartial{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Date:   time.Time{},
				Struct: TestUpdateInner{
					String: "foo",
				},
				Map: model.JsonMap{
					"key": "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []interface{}{2}, "date": "2022-01-03T15:04:05.000Z", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			false,
		},
		"invalidJsonStringUpdate": {
			&TestUpdatePartial{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "", "array": []int{2}},
			true,
		},
		"invalidJsonArrayUpdate": {
			&TestUpdatePartial{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "array": []int{}},
			true,
		},
	}

	for k, v := range testCasesUpdatePartial {
		err := ValidateAndUpdate(v.JsonMapUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseStructUpdateEmptyRequirement(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, -"`
	}

	type TestStructUpdateEmptyRequirement struct {
		String string          `upd:"string, -"`
		Int    int             `upd:"int, -"`
		Float  float64         `upd:"float, -"`
		Array  []string        `upd:"array, -"`
		Map    model.JsonMap   `upd:"map, -"`
		Struct TestUpdateInner `upd:"struct, -"`
	}

	testCases := map[string]*TestRequestWrapperUpdate{
		"valid": {
			&TestStructUpdateEmptyRequirement{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 3, "float": 3.2, "array": []string{"a", "b", "c"}, "struct": map[string]any{"string": "test"}, "map": map[string]any{"key1": "test", "key2": "test", "key3": "test"}},
			false,
		},
		"validOnlyString": {
			&TestStructUpdateEmptyRequirement{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			map[string]interface{}{"string": "Bar"},
			false,
		},
	}

	for k, v := range testCases {
		err := ValidateAndUpdate(v.JsonMapUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseUpdateWithJson(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, equtest"`
	}

	type TestUpdate struct {
		String string          `upd:"string, min1, gr1min7"`
		Int    int             `upd:"int, min1, gr1min7"`
		Float  float64         `upd:"float, min1, gr1min7"`
		Array  []int           `upd:"array, min1, gr1min7"`
		Date   time.Time       `upd:"date, min1, gr1min7"`
		Struct TestUpdateInner `upd:"struct, min1, gr1min7"`
		Map    model.JsonMap   `upd:"map, min1 conkey, gr1min7"`
	}

	testCasesUpdateWithJson := map[string]*TestRequestWrapperUpdateWithJson{
		"validUpdateDateUnix": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "234767127", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			false,
		},
		"validUpdateDateIso8601": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			false,
		},
		"validUpdateDateIso8601Utc": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000Z", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			false,
		},
		"validUpdateDateIso8601WithMicro": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			false,
		},
		"validUpdateDateIso8601WithMicroUtc": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000000Z", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			false,
		},
		"invalidJsonStringUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": Bar, "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidTypeStringUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": 1, "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidJsonIntUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": Bar, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		// TODO add strict case/case for json or url values?
		// "invalidTypeIntUpdate": {
		// 	&TestUpdate{
		// 		String: "Foo",
		// 		Int:    1,
		// 		Float:  1.1,
		// 		Array:  []int{1},
		// 	},
		// 	`{"string": "Bar", "int": "2", "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
		// 	true,
		// },
		"invalidJsonFloatUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": Bar, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		// TODO add strict case/case for json or url values?
		// "invalidTypeFloatUpdate": {
		// 	&TestUpdate{
		// 		String: "Foo",
		// 		Int:    1,
		// 		Float:  1.1,
		// 		Array:  []int{1},
		// 	},
		// 	`{"string": "Bar", "int": 2, "float": "1.2", "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
		// 	true,
		// },
		"invalidJsonArrayUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": Bar, "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidTypeArrayUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": ["2"], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidJsonDateUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03 15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidTypeDateUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": 2022-01-03T15:04:05.000, "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidJsonStructUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "testing"}, "map": {"key": "test"}}`,
			true,
		},
		"invalidTypeStructUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": 1, "map": {"key": "test"}}`,
			true,
		},
		"invalidJsonMapUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"test": "test"}}`,
			true,
		},
		"invalidTypeMapUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			`{"string": "Bar", "int": 2, "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {2: "test"}}`,
			true,
		},
	}

	for k, v := range testCasesUpdateWithJson {
		err := UnmarshalValidateAndUpdate([]byte(v.JsonUpdate), v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseStructUpdateWithoutKey(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, equtest"`
	}

	type TestStructUpdateWithKey struct {
		String string          `upd:"string, min3"`
		Int    int             `upd:"int, min3"`
		Float  float64         `upd:"float, min3"`
		Array  []string        `upd:"array, min3"`
		Map    model.JsonMap   `upd:"map, min3"`
		Struct TestUpdateInner `upd:"struct, min3"`
	}

	type TestStructUpdateWithoutKey struct {
		String string          `upd:"min3"`
		Int    int             `upd:"min3"`
		Float  float64         `upd:"min3"`
		Array  []string        `upd:"min3"`
		Map    model.JsonMap   `upd:"min3"`
		Struct TestUpdateInner `upd:"min3"`
	}

	testCases := map[string]*TestRequestWrapperUpdate{
		"valid": {
			&TestStructUpdateWithKey{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 3, "float": 3.2, "array": []string{"a", "b", "c"}, "struct": map[string]any{"string": "test"}, "map": map[string]any{"key1": "test", "key2": "test", "key3": "test"}},
			false,
		},
		"invalid": {
			&TestStructUpdateWithoutKey{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			map[string]interface{}{"string": "Bar", "int": 3, "float": 3.2, "array": []string{"a", "b", "c"}, "struct": map[string]any{"string": "test"}, "map": map[string]any{"key1": "test", "key2": "test", "key3": "test"}},
			true,
		},
	}

	for k, v := range testCases {
		err := ValidateAndUpdate(v.JsonMapUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseStructUrlValues(t *testing.T) {
	type TestUpdateInner struct {
		String string `upd:"string, equtest"`
	}

	type TestStructUrlValues struct {
		String string          `upd:"string, min3"`
		Int    int             `upd:"int, min3"`
		Float  float64         `upd:"float, min3"`
		Array  []string        `upd:"array, min3"`
		Map    model.JsonMap   `upd:"map, min3"`
		Struct TestUpdateInner `upd:"struct, min3"`
	}

	testCases := map[string]*TestRequestWrapperUpdateWithUrlValues{
		"validSingle": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar"}, "int": []string{"3"}, "float": []string{"3.2"}, "array": []string{`["Ah", "Eh", "Ih"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`}, "struct": []string{`{"string": "test"}`}},
			false,
		},
		"validMulti": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			false,
		},
		"invalidString": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"bu", "Bar"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidInt": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"2", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidTypeInt": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3.2", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidFloat": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"2.5", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidTypeFloat": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"Blubb", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidArray": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidTypeArray": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`[1, 2, 3]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidMap": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidTypeMap": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`["Ah", "Eh", "Ih"]`, `{"key1": "test"}`}, "struct": []string{`{"string": "test"}`, "Blubb"}},
			true,
		},
		"invalidStruct": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`{"string": "muhhh"}`, "Blubb"}},
			true,
		},
		"invalidTypeStruct": {
			&TestStructUrlValues{
				String: "test",
				Int:    3,
				Float:  2.0,
				Array:  []string{"", "", ""},
				Map: model.JsonMap{
					"key": "foo",
				},
				Struct: TestUpdateInner{
					String: "foo",
				},
			},
			url.Values{"string": []string{"Bar", "bu"}, "int": []string{"3", "2"}, "float": []string{"3.2", "3"}, "array": []string{`["Ah", "Eh", "Ih"]`, `["Oh"]`}, "map": []string{`{"key1": "test", "key2": "test", "key3": "test"}`, `{"key1": "test"}`}, "struct": []string{`["Ah", "Eh", "Ih"]`, "Blubb"}},
			true,
		},
	}

	for k, v := range testCases {
		err := UnmapValidateAndUpdate(v.UrlValuesUpdate, v.Data)
		log.Printf("after update: %v, err: %v", v.Data, err)
		assertErrorUpdate(t, k, err, v.Error)
	}
}

func TestCaseParser(t *testing.T) {
	type TestUpdate struct {
		String string `upd:"string, (min10 && max100) || equ'Bar'"`
	}

	lexer := parser.NewLexer(`(min10 && max100) || equ'Bar'`)
	p := parser.NewParser(lexer)
	r, err := p.ParseValidation()
	if err != nil {
		t.Errorf("error parsing: %v", err)
	}

	expectedValidation := "(min'10' && max'100') || equ'Bar'"
	if r.RootValue.AstGroupToString() != expectedValidation {
		t.Errorf("test case parser - wanted %s, got: %s", expectedValidation, r.RootValue.AstGroupToString())
	}

	data := &TestUpdate{
		String: "test",
	}
	err = UnmapValidateAndUpdate(url.Values{"string": []string{"Bar"}}, data)
	log.Printf("after update: %v, err: %v", data, err)
	assertErrorUpdate(t, "validSingle", err, false)

	lexer = parser.NewLexer(`(min10 && max100) || (equ'Bar')`)
	p = parser.NewParser(lexer)
	r, err = p.ParseValidation()
	if err != nil {
		t.Errorf("error parsing: %v", err)
	}

	expectedValidation = "(min'10' && max'100') || (equ'Bar')"
	if r.RootValue.AstGroupToString() != expectedValidation {
		t.Errorf("test case parser - wanted %s, got: %s", expectedValidation, r.RootValue.AstGroupToString())
	}

	data = &TestUpdate{
		String: "test",
	}
	err = UnmapValidateAndUpdate(url.Values{"string": []string{"Bar"}}, data)
	log.Printf("after update: %v, err: %v", data, err)
	assertErrorUpdate(t, "validSingle", err, false)
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

func assertErrorUpdate(t testing.TB, testCase string, err error, errorExpected bool) {
	t.Helper()
	if !errorExpected && err != nil {
		t.Errorf("test case: %s - wanted no error, got error: %v", testCase, err)
	} else if errorExpected && err == nil {
		t.Errorf("test case: %s - wanted error, got no error", testCase)
	}
}

func TestParser(t *testing.T) {
	lexer := parser.NewLexer(`min1 max2 || ((max0) && equ90 || rex')(sa jkdnf')`)
	p := parser.NewParser(lexer)
	r, err := p.ParseValidation()
	if r.RootValue == nil {
		t.Log(r.RootValue)
	} else {
		t.Logf("group: %#v", (*r.RootValue).ConditionGroup[2].ConditionGroup[2])
	}
	t.Logf("group: %#v", (*r.RootValue).AstGroupToString())
	if err != nil {
		t.Errorf("test parser - wanted no error, got error: %v", err)
	}
}

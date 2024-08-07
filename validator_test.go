package validator

import (
	"strings"
	"testing"
	"time"

	"github.com/siherrmann/validator/model"
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

type TestStructNotCon struct {
	String string   `vld:"nco@"`
	Int    int      `vld:"-"`
	Float  float64  `vld:"-"`
	Array  []string `vld:"nco@"`
}

type TestStructFrom struct {
	String string   `vld:"frm@,$"`
	Int    int      `vld:"frm1,200"`
	Float  float64  `vld:"frm0.5,1"`
	Array  []string `vld:"frm@,$"`
}

type TestStructNotFrom struct {
	String string   `vld:"nfr@,$"`
	Int    int      `vld:"nfr1,2"`
	Float  float64  `vld:"nfr0.5,0.8"`
	Array  []string `vld:"nfr@,$"`
}

type TestStructRex struct {
	String string   `vld:"rex^[a-zA-Z0-9]+$"`
	Int    int      `vld:"rex^(2|3)$"`
	Float  float64  `vld:"rex^(2.000|3.000)$"`
	Array  []string `vld:"-"`
}

type TestStructMulti struct {
	String string   `vld:"min3 neqtest"`
	Int    int      `vld:"min3 neq4"`
	Float  float64  `vld:"min3 neq4"`
	Array  []string `vld:"min3 max4"`
}

type TestStructMultiOr struct {
	String string   `vld:"min4 || equte"`
	Int    int      `vld:"min3 || equ0"`
	Float  float64  `vld:"min3 || equ0"`
	Array  []string `vld:"min3 || equ0"`
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

type TestRequestWrapperUpdate struct {
	Data       interface{}
	JsonUpdate map[string]interface{}
	Error      bool
}

type TestRequestWrapperUpdateWithJson struct {
	Data       interface{}
	JsonUpdate string
	Error      bool
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

type TestUpdatePartial struct {
	String string `upd:"string, min1, gr1min2"`
	Int    int
	Float  float64
	Array  []int `upd:"array, min1, gr1min2"`
	Date   time.Time
	Struct TestUpdateInner
	Map    model.JsonMap
}

type TestUpdateInner struct {
	String string `upd:"string, equtest"`
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

	testCases = map[string]*TestRequestWrapper{
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

	testCases = map[string]*TestRequestWrapper{
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
			map[string]interface{}{"string": "Bar", "int": 2, "float": "1.2", "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
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
		"invalidTypeIntUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": "2", "float": 1.2, "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
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
		"invalidTypeFloatUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": "1.2", "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
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
		"invalidTypeDateUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
				Date:   time.Time{},
			},
			map[string]interface{}{"string": "Bar", "int": 2, "float": 1.2, "array": []int{2}, "date": 2024, "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
			true,
		},
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
		err := ValidateAndUpdate(v.JsonUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
	}

	testCasesUpdatePartial := map[string]*TestRequestWrapperUpdate{
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
			map[string]interface{}{"string": "Bar", "int": 2, "float": "1.2", "array": []int{2}, "date": "2022-01-03T15:04:05.000", "struct": map[string]any{"string": "test"}, "map": map[string]any{"key": "test"}},
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
	}

	for k, v := range testCasesUpdatePartial {
		err := ValidateAndUpdate(v.JsonUpdate, v.Data)
		assertErrorUpdate(t, k, err, v.Error)
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
		"invalidTypeIntUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": "2", "float": 1.2, "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
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
		"invalidTypeFloatUpdate": {
			&TestUpdate{
				String: "Foo",
				Int:    1,
				Float:  1.1,
				Array:  []int{1},
			},
			`{"string": "Bar", "int": 2, "float": "1.2", "array": [2], "date": "2022-01-03T15:04:05.000", "struct": {"string": "test"}, "map": {"key": "test"}}`,
			true,
		},
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

package validators

import (
	"reflect"
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
)

func TestCheckArray(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		cond     *model.AstValue
		expected string
	}{
		{
			name:  "Invalid array type",
			value: "not an array",
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.EQUAL,
				ConditionValue: "3",
			},
			expected: "value to validate has to be a array or slice, was string",
		},
		{
			name:  "Invalid Condition Type",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  "INVALID_TYPE",
				ConditionValue: "3",
			},
			expected: "invalid condition type INVALID_TYPE",
		},
		{
			name:  "Valid NONE",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NONE,
				ConditionValue: "",
			},
			expected: "",
		},
		{
			name:  "Valid Array Length EQUAL",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.EQUAL,
				ConditionValue: "3",
			},
			expected: "",
		},
		{
			name:  "Invalid Array Length EQUAL",
			value: []string{"a", "b"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.EQUAL,
				ConditionValue: "3",
			},
			expected: "value shorter than 3",
		},
		{
			name:  "Invalid Type for length EQUAL",
			value: []string{"a", "b"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.EQUAL,
				ConditionValue: "three",
			},
			expected: "strconv.Atoi: parsing \"three\": invalid syntax",
		},
		{
			name:  "Valid Array Length NOT_EQUAL",
			value: []string{"a", "b"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_EQUAL,
				ConditionValue: "3",
			},
			expected: "",
		},
		{name: "Invalid Array Length NOT_EQUAL",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_EQUAL,
				ConditionValue: "3",
			},
			expected: "value longer than 3",
		},
		{
			name:  "Invalid Type for length NOT_EQUAL",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_EQUAL,
				ConditionValue: "three",
			},
			expected: "strconv.Atoi: parsing \"three\": invalid syntax",
		},
		{
			name:  "Valid Array Length MIN_VALUE",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MIN_VALUE,
				ConditionValue: "3",
			},
			expected: "",
		},
		{
			name:  "Invalid Array Length MIN_VALUE",
			value: []string{"a", "b"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MIN_VALUE,
				ConditionValue: "3",
			},
			expected: "value shorter than 3",
		},
		{
			name:  "Invalid Type for length MIN_VALUE",
			value: []string{"a", "b"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MIN_VALUE,
				ConditionValue: "three",
			},
			expected: "strconv.Atoi: parsing \"three\": invalid syntax",
		},
		{
			name:  "Valid Array Length MAX_VALUE",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MAX_VLAUE,
				ConditionValue: "3",
			},
			expected: "",
		},
		{
			name:  "Invalid Array Length MAX_VALUE",
			value: []string{"a", "b", "c", "d"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MAX_VLAUE,
				ConditionValue: "3",
			},
			expected: "value longer than 3",
		},
		{
			name:  "Invalid Type for length MAX_VALUE",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.MAX_VLAUE,
				ConditionValue: "three",
			},
			expected: "strconv.Atoi: parsing \"three\": invalid syntax",
		},
		{
			name:  "Valid Array CONTAINS",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.CONTAINS,
				ConditionValue: "b",
			},
			expected: "",
		},
		{
			name:  "Invalid Array CONTAINS",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.CONTAINS,
				ConditionValue: "d",
			},
			expected: "value does not contain d",
		},
		{
			name:  "Invalid contains type CONTAINS",
			value: []int{1, 2, 3},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.CONTAINS,
				ConditionValue: "d",
			},
			expected: "strconv.Atoi: parsing \"d\": invalid syntax",
		},
		{
			name:  "Valid Array NOT_CONTAINS",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_CONTAINS,
				ConditionValue: "d",
			},
			expected: "",
		},
		{
			name:  "Invalid Array NOT_CONTAINS",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_CONTAINS,
				ConditionValue: "b",
			},
			expected: "value does contain b",
		},
		{
			name:  "Invalid not contains type NOT_CONTAINS",
			value: []int{1, 2, 3},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_CONTAINS,
				ConditionValue: "d",
			},
			expected: "strconv.Atoi: parsing \"d\": invalid syntax",
		},
		{
			name:  "Valid Array FROM",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.FROM,
				ConditionValue: "a,b,c",
			},
			expected: "",
		},
		{
			name:  "Invalid Array FROM",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.FROM,
				ConditionValue: "d,e,f",
			},
			expected: "value not found in [d e f]",
		},
		{
			name:  "Valid Array NOT_FROM",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_FROM,
				ConditionValue: "d,e,f",
			},
			expected: "",
		},
		{
			name:  "Invalid Array NOT_FROM",
			value: []string{"a", "b", "c"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.NOT_FROM,
				ConditionValue: "b,c",
			},
			expected: "value found in [b c]",
		},
		{
			name:  "Valid Array REGX",
			value: []interface{}{"ab", "ac", "ade"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.REGX,
				ConditionValue: "^a.*",
			},
			expected: "",
		},
		{
			name:  "Invalid Array REGX",
			value: []interface{}{"abc", "def", "ghi"},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.REGX,
				ConditionValue: "^x.*",
			},
			expected: "value does match regex ^x.*",
		},
		{
			name:  "Invalid Array type REGX",
			value: []int{1, 2, 3},
			cond: &model.AstValue{
				Type:           model.CONDITION,
				ConditionType:  model.REGX,
				ConditionValue: "^1.*",
			},
			expected: "type slice not supported for regex",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckArray(test.value, test.cond)
			if len(test.expected) != 0 {
				assert.Error(t, err, "Expected error for test %v but got nil", test.name)
				assert.EqualError(t, err, test.expected, "Expected error message for test %v but got %v", test.name, err)
			} else {
				assert.NoError(t, err, "Expected no error for test %v but got %v", test.name, err)
			}
		})
	}
}

func TestValueContains(t *testing.T) {
	tests := []struct {
		name           string
		value          reflect.Value
		contain        string
		expectedResult string
		expectedError  string
	}{
		{
			name:          "Invalid array type",
			value:         reflect.ValueOf("not an array"),
			contain:       "2",
			expectedError: "type string not supported",
		},
		{
			name:           "Valid contains int",
			value:          reflect.ValueOf([]int{1, 2, 3}),
			contain:        "2",
			expectedResult: "2",
		},
		{
			name:           "Invalid does not contain int",
			value:          reflect.ValueOf([]int{1, 2, 3}),
			contain:        "4",
			expectedResult: "",
		},
		{
			name:          "Invalid contains type int",
			value:         reflect.ValueOf([]int{1, 2, 3}),
			contain:       "d",
			expectedError: "strconv.Atoi: parsing \"d\": invalid syntax",
		},
		{
			name:           "Valid contains float32",
			value:          reflect.ValueOf([]float32{1.1, 2.2, 3.3}),
			contain:        "2.2",
			expectedResult: "2.2",
		},
		{
			name:           "Invalid does not contain float32",
			value:          reflect.ValueOf([]float32{1.1, 2.2, 3.3}),
			contain:        "4.4",
			expectedResult: "",
		},
		{
			name:          "Invalid contains type float32",
			value:         reflect.ValueOf([]float32{1.1, 2.2, 3.3}),
			contain:       "d",
			expectedError: "strconv.ParseFloat: parsing \"d\": invalid syntax",
		},
		{
			name:           "Valid contains float64",
			value:          reflect.ValueOf([]float64{1.1, 2.2, 3.3}),
			contain:        "2.2",
			expectedResult: "2.2",
		},
		{
			name:           "Invalid does not contain float64",
			value:          reflect.ValueOf([]float64{1.1, 2.2, 3.3}),
			contain:        "4.4",
			expectedResult: "",
		},
		{
			name:          "Invalid contains type float64",
			value:         reflect.ValueOf([]float64{1.1, 2.2, 3.3}),
			contain:       "d",
			expectedError: "strconv.ParseFloat: parsing \"d\": invalid syntax",
		},
		{
			name:           "Valid contains string",
			value:          reflect.ValueOf([]string{"a", "b", "c"}),
			contain:        "b",
			expectedResult: "b",
		},
		{
			name:           "Invalid does not contain string",
			value:          reflect.ValueOf([]string{"a", "b", "c"}),
			contain:        "d",
			expectedResult: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ValueContains(test.value, test.contain)
			if len(test.expectedError) != 0 {
				assert.Error(t, err, "Expected error for test %v but got nil", test.name)
				assert.EqualError(t, err, test.expectedError, "Expected error message for test %v but got %v", test.name, err)
				assert.Empty(t, result, "Expected empty result for test %v but got %v", test.name, result)
			} else {
				assert.NoError(t, err, "Expected no error for test %v but got %v", test.name, err)
				assert.Equal(t, test.expectedResult, result, "Expected result for test %v to be %v but got %v", test.name, test.expectedResult, result)
			}
		})
	}
}

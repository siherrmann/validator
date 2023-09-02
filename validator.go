package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	EQUAL     string = "equ"
	NOT_EQUAL string = "neq"
	MIN_VALUE string = "min"
	MAX_VLAUE string = "max"
	CONTAINS  string = "con"
	REGX      string = "rex"
)

// Validate validates a given struct by vld tags.
// Validate needs a struct as input.
//
// All fields in the struct need a vld tag.
// If you want to ignore one field in the validator you can add `vld:"-"`.
// If you don't add the vld tag to every field the function will fail with an error.
//
// Conditions have different usages per variable type:
//
// equ - int/float/string == condition, len(array) == condition
//
// neq - int/float/string != condition, len(array) != condition
//
// min - int/float >= condition, len(string/array) >= condition
//
// max - int/float <= condition, len(string/array) <= condition
//
// con - strings.Contains(string, condition), contains(array, condition), int/float ignored
//
// rex - regexp.MatchString(condition, int/float/string), array ignored
//
// For con you need to put in a condition that is convertable to the underlying type of the arrary.
// Eg. for an array of int the condition must be convertable to int (bad: `vld:"conA"`, good: `vld:"con1"`).
//
// In the case of rex the int and float input will get converted to a string (strconv.Itoa(int) and fmt.Sprintf("%f", f)).
// If you want to check more complex cases you can obviously replace equ, neq, min, max and con with one regular expression.
func Validate(value any) error {
	// check if value is a struct
	if reflect.ValueOf(value).Kind() != reflect.Struct {
		return errors.New("value to validate has to be a pointer to a struct")
	}

	// get valid reflect value of struct
	structValue := reflect.ValueOf(value)

	for i := 0; i < structValue.Type().NumField(); i++ {
		tag := structValue.Type().Field(i).Tag.Get("vld")
		fieldName := structValue.Type().Field(i).Name
		if len(strings.TrimSpace(tag)) == 0 {
			return fmt.Errorf("no validate tag found for field %s, to ignore the validation put a 'vld:\"-\"' into the tag", structValue.Type().Field(i).Name)
		}

		value := structValue.Field(i)
		conditions := strings.Split(tag, ",")

		switch value.Type().Kind() {
		case reflect.Int:
			valueTemp := value.Int()
			err := checkInt(int(valueTemp), conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.Float32:
		case reflect.Float64:
			valueTemp := value.Float()
			err := checkFloat(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.String:
			valueTemp := value.String()
			err := checkString(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.Array:
		case reflect.Slice:
			valueTemp := value
			err := checkArray(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		}
	}

	return nil
}

func checkInt(i int, c []string) error {
	condition, err := getConditionByType(c, EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		equal, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if i != equal {
			return fmt.Errorf("value must be equal to %v", equal)
		}
	}

	condition, err = getConditionByType(c, NOT_EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		notEqual, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if i == notEqual {
			return fmt.Errorf("value can't be equal to %v", notEqual)
		}
	}

	condition, err = getConditionByType(c, MIN_VALUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		minValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if i < minValue {
			return fmt.Errorf("value smaller than %v", minValue)
		}
	}

	condition, err = getConditionByType(c, MAX_VLAUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		maxValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if i > maxValue {
			return fmt.Errorf("value greater than %v", maxValue)
		}
	}

	condition, err = getConditionByType(c, REGX)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		match, err := regexp.MatchString(condition, strconv.Itoa(i))
		if err != nil {
			return err
		} else if !match {
			return fmt.Errorf("value does match regex %v", condition)
		}
	}

	return nil
}

func checkFloat(f float64, c []string) error {
	condition, err := getConditionByType(c, EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		equal, err := strconv.ParseFloat(condition, 64)
		if err != nil {
			return err
		} else if f != equal {
			return fmt.Errorf("value must be equal to %v", equal)
		}
	}

	condition, err = getConditionByType(c, NOT_EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		notEqual, err := strconv.ParseFloat(condition, 64)
		if err != nil {
			return err
		} else if f == notEqual {
			return fmt.Errorf("value can't be equal to %v", notEqual)
		}
	}

	condition, err = getConditionByType(c, MIN_VALUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		minValue, err := strconv.ParseFloat(condition, 64)
		if err != nil {
			return err
		} else if f < minValue {
			return fmt.Errorf("value smaller than %v", minValue)
		}
	}

	condition, err = getConditionByType(c, MAX_VLAUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		maxValue, err := strconv.ParseFloat(condition, 64)
		if err != nil {
			return err
		} else if f > maxValue {
			return fmt.Errorf("value greater than %v", maxValue)
		}
	}

	condition, err = getConditionByType(c, REGX)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		match, err := regexp.MatchString(condition, fmt.Sprintf("%f", f))
		if err != nil {
			return err
		} else if !match {
			return fmt.Errorf("value does match regex %v", condition)
		}
	}

	return nil
}

func checkString(s string, c []string) error {
	condition, err := getConditionByType(c, EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		if s != condition {
			return fmt.Errorf("value must be equal to %v", condition)
		}
	}

	condition, err = getConditionByType(c, NOT_EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		if s == condition {
			return fmt.Errorf("value can't be equal to %v", condition)
		}
	}

	condition, err = getConditionByType(c, MIN_VALUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		minValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if len(s) < minValue {
			return fmt.Errorf("value shorter than %v", minValue)
		}
	}

	condition, err = getConditionByType(c, MAX_VLAUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		maxValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if len(s) > maxValue {
			return fmt.Errorf("value longer than %v", maxValue)
		}
	}

	condition, err = getConditionByType(c, CONTAINS)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		if !strings.Contains(s, condition) {
			return fmt.Errorf("value does not include %v", condition)
		}
	}

	condition, err = getConditionByType(c, REGX)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		match, err := regexp.MatchString(condition, s)
		if err != nil {
			return err
		} else if !match {
			return fmt.Errorf("value does match regex %v", condition)
		}
	}

	return nil
}

func checkArray(a reflect.Value, c []string) error {
	if a.Type().Kind() != reflect.Array && a.Type().Kind() != reflect.Slice {
		return fmt.Errorf("value to validate has to be a array or slice, was %v", a.Type().Kind())
	}

	condition, err := getConditionByType(c, EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		equal, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if a.Len() != equal {
			return fmt.Errorf("value shorter than %v", equal)
		}
	}

	condition, err = getConditionByType(c, NOT_EQUAL)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		notEqual, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if a.Len() == notEqual {
			return fmt.Errorf("value longer than %v", notEqual)
		}
	}

	condition, err = getConditionByType(c, MIN_VALUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		minValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if a.Len() < minValue {
			return fmt.Errorf("value shorter than %v", minValue)
		}
	}

	condition, err = getConditionByType(c, MAX_VLAUE)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		maxValue, err := strconv.Atoi(condition)
		if err != nil {
			return err
		} else if a.Len() > maxValue {
			return fmt.Errorf("value longer than %v", maxValue)
		}
	}

	condition, err = getConditionByType(c, CONTAINS)
	if err != nil {
		return err
	}
	if len(condition) != 0 {
		switch v := a.Interface().(type) {
		case []int:
			contain, err := strconv.Atoi(condition)
			if err != nil {
				return err
			} else if !contains(v, contain) {
				return fmt.Errorf("value does not contain %v", contain)
			}
		case []float32:
			contain, err := strconv.ParseFloat(condition, 32)
			if err != nil {
				return err
			} else if !contains(v, float32(contain)) {
				return fmt.Errorf("value does not contain %v", contain)
			}
		case []float64:
			contain, err := strconv.ParseFloat(condition, 64)
			if err != nil {
				return err
			} else if !contains(v, contain) {
				return fmt.Errorf("value does not contain %v", contain)
			}
		case []string:
			if !contains(v, condition) {
				return fmt.Errorf("value does not contain %v", condition)
			}
		default:
			return fmt.Errorf("type %v not supported", reflect.TypeOf(v))
		}
	}

	return nil
}

func getConditionByType(conditions []string, conditionType string) (string, error) {
	first := firstWhere(conditions, func(v string) bool { return strings.HasPrefix(v, conditionType) })
	if first == "" {
		return "", nil
	}

	condition := strings.TrimPrefix(first, conditionType)
	if len(condition) == 0 {
		return "", fmt.Errorf("empty %s value", conditionType)
	}

	return condition, nil
}

func contains[V comparable](list []V, v V) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}

func firstWhere[V comparable](list []V, where func(v V) bool) V {
	var temp V
	for i := 0; i < len(list); i++ {
		if where(list[i]) {
			temp = list[i]
			break
		}
	}
	return temp
}

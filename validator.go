package main

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
			return errors.New("no validate tag found, to ignore the validation put a 'vld:\"-\"' into the tag")
		}

		value := structValue.Field(i)
		conditions := strings.Split(tag, ",")

		switch value.Type().Kind() {
		case reflect.Float64:
			valueTemp := value.Float()
			err := checkFloat(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.Int:
			valueTemp := value.Int()
			err := checkInt(int(valueTemp), conditions)
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
	if a.Kind() != reflect.Array || a.Kind() != reflect.Slice {
		return errors.New("value to validate has to be a array or slice")
	}

	condition, err := getConditionByType(c, MIN_VALUE)
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

	// TODO fix check for different types of arrays
	// condition, err = getConditionByType(c, CONTAINS)
	// if err != nil {
	// 	return err
	// }
	// if len(condition) != 0 {
	// 	if !contains(a, condition) {
	// 		return fmt.Errorf("value does not include %v", condition)
	// 	}
	// }

	return nil
}

func getConditionByType(conditions []string, conditionType string) (string, error) {
	first := FirstWhere(conditions, func(v string) bool { return strings.HasPrefix(v, conditionType) })
	if first == "" {
		return "", nil
	}

	condition := strings.TrimPrefix(first, conditionType)
	if len(condition) == 0 {
		return "", errors.New("empty 'min' value for float")
	}

	return condition, nil
}

func Contains[V comparable](list []V, v V) bool {
	for _, s := range list {
		if v == s {
			return true
		}
	}
	return false
}

func FirstWhere[V comparable](list []V, where func(v V) bool) V {
	var temp V
	for i := 0; i < len(list); i++ {
		if where(list[i]) {
			temp = list[i]
			break
		}
	}
	return temp
}

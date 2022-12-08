package validator

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

const (
	MIN_VALUE string = "min"
	MAX_VLAUE string = "max"
	CONTAINS  string = "con"
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
			valueTemp := reflect.ValueOf(value).Float()
			err := checkFloat(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.Int:
			// TODO validate int
		case reflect.Uint:
			// TODO validate uint
		case reflect.String:
			valueTemp := reflect.ValueOf(value).String()
			err := checkString(valueTemp, conditions)
			if err != nil {
				return fmt.Errorf("field %v invalid: %v", fieldName, err.Error())
			}
		case reflect.Array:
		case reflect.Slice:
			// TODO validate array/slice
		}
	}

	return nil
}

func checkFloat(f float64, c []string) error {
	condition, err := getConditionByType(c, MIN_VALUE)
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

func checkString(s string, c []string) error {
	condition, err := getConditionByType(c, MIN_VALUE)
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

	return nil
}

func getConditionByType(conditions []string, conditionType string) (string, error) {
	if contains(conditions, conditionType) {
		condition := strings.TrimPrefix(firstWhere(conditions, func(v string) bool { return strings.HasPrefix(v, conditionType) }), conditionType)

		log.Println(condition)

		if len(condition) == 0 {
			return "", errors.New("empty 'min' value for float")
		}

		return condition, nil
	} else {
		return "", nil
	}
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

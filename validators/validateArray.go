package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckArray[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return fmt.Errorf("value to validate has to be a array or slice, was %v", rv.Kind())
	}

	switch c.ConditionType {
	case model.EQUAL:
		equal, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if rv.Len() != equal {
			return fmt.Errorf("value shorter than %v", equal)
		}
	case model.NOT_EQUAL:
		notEqual, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if rv.Len() == notEqual {
			return fmt.Errorf("value longer than %v", notEqual)
		}
	case model.MIN_VALUE:
		minValue, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if rv.Len() < minValue {
			return fmt.Errorf("value shorter than %v", minValue)
		}
	case model.MAX_VLAUE:
		maxValue, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if rv.Len() > maxValue {
			return fmt.Errorf("value longer than %v", maxValue)
		}
	case model.CONTAINS:
		contains, err := ValueContains(rv, c.ConditionValue)
		if err != nil {
			return err
		} else if len(contains) == 0 {
			return fmt.Errorf("value does not contain %v", c.ConditionValue)
		}
	case model.NOT_CONTAINS:
		contains, err := ValueContains(rv, c.ConditionValue)
		if err != nil {
			return err
		} else if len(contains) != 0 {
			return fmt.Errorf("value does contain %v", c.ConditionValue)
		}
	case model.FROM:
		fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		notFound, err := ValueFrom(rv, fromValues, true)
		if err != nil {
			return err
		} else if len(notFound) != 0 {
			return fmt.Errorf("value not found in %v", fromValues)
		}
	case model.NOT_FROM:
		notFromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		found, err := ValueFrom(rv, notFromValues, false)
		if err != nil {
			return err
		} else if len(found) != 0 {
			return fmt.Errorf("value found in %v", notFromValues)
		}
	case model.REGX:
		if vArray, ok := rv.Interface().([]any); ok {
			for _, av := range vArray {
				match, err := regexp.MatchString(c.ConditionValue, fmt.Sprint(av))
				if err != nil {
					return err
				} else if !match {
					return fmt.Errorf("value does match regex %v", c.ConditionValue)
				}
			}
		} else {
			return fmt.Errorf("type %v not supported for regex", rv.Type().Kind())
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}
	return nil
}

// ValueContains checks, if the [reflectValue] contains [contain].
func ValueContains(reflectValue reflect.Value, contain string) (string, error) {
	switch v := reflectValue.Interface().(type) {
	case []int:
		c, err := strconv.Atoi(contain)
		if err != nil {
			return "", err
		} else if slices.Contains(v, c) {
			return contain, nil
		}
		return "", nil
	case []float32:
		c, err := strconv.ParseFloat(contain, 32)
		if err != nil {
			return "", err
		} else if slices.Contains(v, float32(c)) {
			return contain, nil
		}
		return "", nil
	case []float64:
		c, err := strconv.ParseFloat(contain, 64)
		if err != nil {
			return "", err
		} else if slices.Contains(v, c) {
			return contain, nil
		}
		return "", nil
	case []string:
		if slices.Contains(v, contain) {
			return contain, nil
		}
		return "", nil
	default:
		return "", fmt.Errorf("type %v not supported", reflect.TypeOf(v))
	}
}

// ValueFrom checks, if the [reflectValue] consists of only ([shouldFind] == true)/of none ([shouldFind] == false) of the values from [from].
// If a value is not found with [shouldFind] == true the missing value is given back.
// If a value is found with [shouldFind] == false the value found is given back.
func ValueFrom(reflectValue reflect.Value, from []string, shouldFind bool) (string, error) {
	switch v := reflectValue.Interface().(type) {
	case []int:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.Atoi(fromValue)
				if err != nil {
					return false
				}
				return c == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []float32:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.ParseFloat(fromValue, 32)
				if err != nil {
					return false
				}
				return float32(c) == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []float64:
		for _, value := range v {
			contains := slices.ContainsFunc(from, func(fromValue string) bool {
				c, err := strconv.ParseFloat(fromValue, 64)
				if err != nil {
					return false
				}
				return c == value
			})
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	case []string:
		for _, value := range v {
			contains := slices.Contains(from, value)
			if contains != shouldFind {
				return fmt.Sprint(value), nil
			}
		}
	default:
		return "", fmt.Errorf("type %v not supported", reflect.TypeOf(v))
	}
	return "", nil
}

package validators

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckMap[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	rv := reflect.ValueOf(v)
	if rv.Type().Kind() != reflect.Map {
		return fmt.Errorf("value to validate has to be a map, was %v", rv.Type().Kind())
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
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = rv.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(rv))
		}
		if _, ok := jsonMap[c.ConditionValue]; !ok {
			return fmt.Errorf("value does not contain key %v", c.ConditionValue)
		}
	case model.NOT_CONTAINS:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = rv.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(rv))
		}
		if _, ok := jsonMap[c.ConditionValue]; ok {
			return fmt.Errorf("value does contain key %v", c.ConditionValue)
		}
	case model.FROM:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = rv.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(rv))
		}
		fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		for k := range jsonMap {
			if !slices.Contains(fromValues, k) {
				return fmt.Errorf("key not found in %v", fromValues)
			}
		}
	case model.NOT_FROM:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = rv.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(rv))
		}
		fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		for k := range jsonMap {
			if slices.Contains(fromValues, k) {
				return fmt.Errorf("key found in %v", fromValues)
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

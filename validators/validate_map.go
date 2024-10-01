package validators

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckMap(v reflect.Value, c *model.AstValue) error {
	if v.Type().Kind() != reflect.Map {
		return fmt.Errorf("value to validate has to be a map, was %v", v.Type().Kind())
	}

	switch c.ConditionType {
	case model.EQUAL:
		if len(c.ConditionValue) != 0 {
			equal, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if v.Len() != equal {
				return fmt.Errorf("value shorter than %v", equal)
			}
		}
	case model.NOT_EQUAL:
		if len(c.ConditionValue) != 0 {
			notEqual, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if v.Len() == notEqual {
				return fmt.Errorf("value longer than %v", notEqual)
			}
		}
	case model.MIN_VALUE:
		if len(c.ConditionValue) != 0 {
			minValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if v.Len() < minValue {
				return fmt.Errorf("value shorter than %v", minValue)
			}
		}
	case model.MAX_VLAUE:
		if len(c.ConditionValue) != 0 {
			maxValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if v.Len() > maxValue {
				return fmt.Errorf("value longer than %v", maxValue)
			}
		}
	case model.CONTAINS:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = v.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(v))
		}
		if len(c.ConditionValue) != 0 {
			if _, ok := jsonMap[c.ConditionValue]; !ok {
				return fmt.Errorf("value does not contain key %v", c.ConditionValue)
			}
		}
	case model.NOT_CONTAINS:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = v.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(v))
		}
		if len(c.ConditionValue) != 0 {
			if _, ok := jsonMap[c.ConditionValue]; ok {
				return fmt.Errorf("value does contain key %v", c.ConditionValue)
			}
		}
	case model.FROM:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = v.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(v))
		}
		if len(c.ConditionValue) != 0 {
			fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
			if err != nil {
				return err
			}
			for k := range jsonMap {
				if !slices.Contains(fromValues, k) {
					return fmt.Errorf("key not found in %v", fromValues)
				}
			}
		}
	case model.NOT_FROM:
		var jsonMap model.JsonMap
		var ok bool
		if jsonMap, ok = v.Interface().(map[string]any); !ok {
			return fmt.Errorf("value has to be of type map, was %v", reflect.TypeOf(v))
		}
		if len(c.ConditionValue) != 0 {
			fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
			if err != nil {
				return err
			}
			for k := range jsonMap {
				if slices.Contains(fromValues, k) {
					return fmt.Errorf("key found in %v", fromValues)
				}
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

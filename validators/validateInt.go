package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckInt[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	var i int
	var f float64
	var ok bool
	i, ok = any(v).(int)
	if !ok {
		if f, ok = any(v).(float64); !ok {
			return fmt.Errorf("value to validate has to be an int or float64, was %v", reflect.TypeOf(v))
		}
		i = int(f)
	}

	switch c.ConditionType {
	case model.EQUAL:
		equal, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if i != equal {
			return fmt.Errorf("value must be equal to %v", equal)
		}
	case model.NOT_EQUAL:
		notEqual, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if i == notEqual {
			return fmt.Errorf("value can't be equal to %v", notEqual)
		}
	case model.MIN_VALUE:
		minValue, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if i < minValue {
			return fmt.Errorf("value smaller than %v", minValue)
		}
	case model.MAX_VLAUE:
		maxValue, err := strconv.Atoi(c.ConditionValue)
		if err != nil {
			return err
		} else if i > maxValue {
			return fmt.Errorf("value greater than %v", maxValue)
		}
	case model.FROM:
		fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		foundInFromValues := slices.ContainsFunc(fromValues, func(fromValue string) bool {
			from, err := strconv.Atoi(fromValue)
			if err != nil {
				return false
			}
			return i == from
		})
		if !foundInFromValues {
			return fmt.Errorf("value not found in %v", fromValues)
		}
	case model.NOT_FROM:
		notFromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		foundInFromValues := slices.ContainsFunc(notFromValues, func(notFromValue string) bool {
			notFrom, err := strconv.Atoi(notFromValue)
			if err != nil {
				return false
			}
			return i == notFrom
		})
		if foundInFromValues {
			return fmt.Errorf("value found in %v", notFromValues)
		}
	case model.REGX:
		match, err := regexp.MatchString(c.ConditionValue, fmt.Sprint(i))
		if err != nil {
			return err
		} else if !match {
			return fmt.Errorf("value does match regex %v", c.ConditionValue)
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

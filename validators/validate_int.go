package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckInt(v reflect.Value, c *model.AstValue) error {
	if v.Type().Kind() != reflect.Int && v.Type().Kind() != reflect.Int64 && v.Type().Kind() != reflect.Int32 && v.Type().Kind() != reflect.Int16 && v.Type().Kind() != reflect.Int8 {
		return fmt.Errorf("value to validate has to be a int, was %v", v.Type().Kind())
	}

	i := int(v.Int())

	switch c.ConditionType {
	case model.EQUAL:
		if len(c.ConditionValue) != 0 {
			equal, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if i != equal {
				return fmt.Errorf("value must be equal to %v", equal)
			}
		}
	case model.NOT_EQUAL:
		if len(c.ConditionValue) != 0 {
			notEqual, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if i == notEqual {
				return fmt.Errorf("value can't be equal to %v", notEqual)
			}
		}
	case model.MIN_VALUE:
		if len(c.ConditionValue) != 0 {
			minValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if i < minValue {
				return fmt.Errorf("value smaller than %v", minValue)
			}
		}
	case model.MAX_VLAUE:
		if len(c.ConditionValue) != 0 {
			maxValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if i > maxValue {
				return fmt.Errorf("value greater than %v", maxValue)
			}
		}
	case model.FROM:
		if len(c.ConditionValue) != 0 {
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
		}
	case model.NOT_FROM:
		if len(c.ConditionValue) != 0 {
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
		}
	case model.REGX:
		if len(c.ConditionValue) != 0 {
			match, err := regexp.MatchString(c.ConditionValue, fmt.Sprint(i))
			if err != nil {
				return err
			} else if !match {
				return fmt.Errorf("value does match regex %v", c.ConditionValue)
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

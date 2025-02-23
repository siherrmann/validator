package validators

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/siherrmann/validator/model"
)

func CheckTime(v reflect.Value, c *model.AstValue) error {
	t, ok := v.Interface().(time.Time)
	if !ok {
		return fmt.Errorf("invalid time: %v", v.Interface())
	}

	validation := model.Validation{Type: model.Time}
	switch c.ConditionType {
	case model.EQUAL:
		if len(c.ConditionValue) != 0 {
			compareTime, err := validation.InterfaceFromString(c.ConditionValue)
			if err != nil {
				return err
			}
			if !t.Equal(compareTime.(time.Time)) {
				return fmt.Errorf("value must be equal to %v", c.ConditionValue)
			}
		}
	case model.NOT_EQUAL:
		if len(c.ConditionValue) != 0 {
			compareTime, err := validation.InterfaceFromString(c.ConditionValue)
			if err != nil {
				return err
			}
			if t.Equal(compareTime.(time.Time)) {
				return fmt.Errorf("value can't be equal to %v", c.ConditionValue)
			}
		}
	case model.MIN_VALUE:
		if len(c.ConditionValue) != 0 {
			minValue, err := validation.InterfaceFromString(c.ConditionValue)
			if err != nil {
				return err
			}
			if t.Before(minValue.(time.Time)) {
				return fmt.Errorf("value smaller than %v", minValue)
			}
		}
	case model.MAX_VLAUE:
		if len(c.ConditionValue) != 0 {
			maxValue, err := validation.InterfaceFromString(c.ConditionValue)
			if err != nil {
				return err
			}
			if t.After(maxValue.(time.Time)) {
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
				from, err := validation.InterfaceFromString(fromValue)
				if err != nil {
					return false
				}
				return t.Equal(from.(time.Time))
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
				notFrom, err := validation.InterfaceFromString(notFromValue)
				if err != nil {
					return false
				}
				return t.Equal(notFrom.(time.Time))
			})
			if foundInFromValues {
				return fmt.Errorf("value found in %v", notFromValues)
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

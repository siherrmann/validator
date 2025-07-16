package validators

import (
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/siherrmann/validator/model"
)

func CheckTime[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	var t time.Time
	var s string
	var ok bool
	t, ok = any(v).(time.Time)
	if !ok {
		if s, ok = any(v).(string); !ok {
			return fmt.Errorf("value to validate has to be a time.Time or string, was %v", reflect.TypeOf(v))
		}
		var err error
		t, err = model.UnixStringToTime(s)
		if err != nil {
			t, err = model.ISO8601StringToTime(s)
			if err != nil {
				return err
			}
		}
	}

	validation := model.Validation{Type: model.Time}
	switch c.ConditionType {
	case model.EQUAL:
		compareTime, err := validation.InterfaceFromString(c.ConditionValue)
		if err != nil {
			return err
		}
		if !t.Equal(compareTime.(time.Time)) {
			return fmt.Errorf("value must be equal to %v", c.ConditionValue)
		}
	case model.NOT_EQUAL:
		compareTime, err := validation.InterfaceFromString(c.ConditionValue)
		if err != nil {
			return err
		}
		if t.Equal(compareTime.(time.Time)) {
			return fmt.Errorf("value can't be equal to %v", c.ConditionValue)
		}
	case model.MIN_VALUE:
		minValue, err := validation.InterfaceFromString(c.ConditionValue)
		if err != nil {
			return err
		}
		if t.Before(minValue.(time.Time)) {
			return fmt.Errorf("value smaller than %v", minValue)
		}
	case model.MAX_VLAUE:
		maxValue, err := validation.InterfaceFromString(c.ConditionValue)
		if err != nil {
			return err
		}
		if t.After(maxValue.(time.Time)) {
			return fmt.Errorf("value greater than %v", maxValue)
		}
	case model.FROM:
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
	case model.NOT_FROM:
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
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

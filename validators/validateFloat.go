package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckFloat[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	f, ok := any(v).(float64)
	if !ok {
		return fmt.Errorf("value to validate has to be a float64, was %v", reflect.TypeOf(v))
	}

	switch c.ConditionType {
	case model.EQUAL:
		equal, err := strconv.ParseFloat(c.ConditionValue, 64)
		if err != nil {
			return err
		} else if f != equal {
			return fmt.Errorf("value must be equal to %v", equal)
		}
	case model.NOT_EQUAL:
		notEqual, err := strconv.ParseFloat(c.ConditionValue, 64)
		if err != nil {
			return err
		} else if f == notEqual {
			return fmt.Errorf("value can't be equal to %v", notEqual)
		}
	case model.MIN_VALUE:
		minValue, err := strconv.ParseFloat(c.ConditionValue, 64)
		if err != nil {
			return err
		} else if f < minValue {
			return fmt.Errorf("value smaller than %v", minValue)
		}
	case model.MAX_VLAUE:
		maxValue, err := strconv.ParseFloat(c.ConditionValue, 64)
		if err != nil {
			return err
		} else if f > maxValue {
			return fmt.Errorf("value greater than %v", maxValue)
		}
	case model.FROM:
		fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		foundInFromValues := false
		for _, fromValue := range fromValues {
			from, err := strconv.ParseFloat(fromValue, 64)
			if err != nil {
				return err
			}
			if f == from {
				foundInFromValues = true
				break
			}
		}
		if !foundInFromValues {
			return fmt.Errorf("value not found in %v", fromValues)
		}
	case model.NOT_FROM:
		notFromValues, err := model.GetArrayFromCondition(c.ConditionValue)
		if err != nil {
			return err
		}
		for _, notFromValue := range notFromValues {
			notFrom, err := strconv.ParseFloat(notFromValue, 64)
			if err != nil {
				return err
			}
			if f == notFrom {
				return fmt.Errorf("value found in %v", notFromValues)
			}
		}
	case model.REGX:
		match, err := regexp.MatchString(c.ConditionValue, fmt.Sprint(f))
		if err != nil {
			return err
		} else if !match {
			return fmt.Errorf("value does match regex %v", c.ConditionValue)
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid c.ConditionValue type %s", c.ConditionType)
	}

	return nil
}

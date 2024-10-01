package validators

import (
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/siherrmann/validator/model"
)

func CheckString(v reflect.Value, c *model.AstValue) error {
	if v.Type().Kind() != reflect.String {
		return fmt.Errorf("value to validate has to be a string, was %v", v.Type().Kind())
	}

	s := v.String()

	switch c.ConditionType {
	case model.EQUAL:
		if len(c.ConditionValue) != 0 {
			if s != c.ConditionValue {
				return fmt.Errorf("value must be equal to %v", c.ConditionValue)
			}
		}
	case model.NOT_EQUAL:
		if len(c.ConditionValue) != 0 {
			if s == c.ConditionValue {
				return fmt.Errorf("value can't be equal to %v", c.ConditionValue)
			}
		}
	case model.MIN_VALUE:
		if len(c.ConditionValue) != 0 {
			minValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if len(strings.TrimSpace(s)) < minValue {
				return fmt.Errorf("value shorter than %v", minValue)
			}
		}
	case model.MAX_VLAUE:
		if len(c.ConditionValue) != 0 {
			maxValue, err := strconv.Atoi(c.ConditionValue)
			if err != nil {
				return err
			} else if len(strings.TrimSpace(s)) > maxValue {
				return fmt.Errorf("value longer than %v", maxValue)
			}
		}
	case model.CONTAINS:
		if len(c.ConditionValue) != 0 {
			if !strings.Contains(s, c.ConditionValue) {
				return fmt.Errorf("value does not contain %v", c.ConditionValue)
			}
		}
	case model.NOT_CONTAINS:
		if len(c.ConditionValue) != 0 {
			if strings.Contains(s, c.ConditionValue) {
				return fmt.Errorf("value does contain %v", c.ConditionValue)
			}
		}
	case model.FROM:
		if len(c.ConditionValue) != 0 {
			fromValues, err := model.GetArrayFromCondition(c.ConditionValue)
			if err != nil {
				return err
			}
			if !slices.Contains(fromValues, s) {
				return fmt.Errorf("value not found in %v", fromValues)
			}
		}
	case model.NOT_FROM:
		if len(c.ConditionValue) != 0 {
			notFromValues, err := model.GetArrayFromCondition(c.ConditionValue)
			if err != nil {
				return err
			}
			if slices.Contains(notFromValues, s) {
				return fmt.Errorf("value found in %v", notFromValues)
			}
		}
	case model.REGX:
		if len(c.ConditionValue) != 0 {
			match, err := regexp.MatchString(c.ConditionValue, s)
			if err != nil {
				return err
			} else if !match {
				return fmt.Errorf("value does match regex %v", c.ConditionValue)
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid c.ConditionValue type %s", c.ConditionType)
	}

	return nil
}

package validators

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckBool[T comparable](v T, c *model.AstValue) error {
	if len(c.ConditionValue) == 0 {
		return nil
	}

	b, ok := any(v).(bool)
	if !ok {
		return fmt.Errorf("value to validate has to be a bool, was %v", reflect.TypeOf(v))
	}

	conditionValue, err := strconv.ParseBool(c.ConditionValue)
	if err != nil {
		return err
	}

	switch c.ConditionType {
	case model.EQUAL:
		if b != conditionValue {
			return fmt.Errorf("value must be equal to %v", conditionValue)
		}
	case model.NOT_EQUAL:
		if b == conditionValue {
			return fmt.Errorf("value can't be equal to %v", conditionValue)
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid condition type %s", c.ConditionType)
	}

	return nil
}

package validators

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func CheckBool(v reflect.Value, c *model.AstValue) error {
	if v.Type().Kind() != reflect.Bool {
		return fmt.Errorf("value to validate has to be a bool, was %v", v.Type().Kind())
	}

	b := v.Bool()

	switch c.ConditionType {
	case model.EQUAL:
		if len(c.ConditionValue) != 0 {
			equal, err := strconv.ParseBool(c.ConditionValue)
			if err != nil {
				return err
			} else if b != equal {
				return fmt.Errorf("value must be equal to %v", equal)
			}
		}
	case model.NOT_EQUAL:
		if len(c.ConditionValue) != 0 {
			notEqual, err := strconv.ParseBool(c.ConditionValue)
			if err != nil {
				return err
			} else if b == notEqual {
				return fmt.Errorf("value can't be equal to %v", notEqual)
			}
		}
	case model.NONE:
		return nil
	default:
		return fmt.Errorf("invalid c.ConditionValue type %s", c.ConditionType)
	}

	return nil
}

package validators

import (
	"fmt"

	"github.com/siherrmann/validator/model"
)

// RunFuncOnConditionGroup runs the function [f] on each condition in the [astValue].
// If the condition is a group, it recursively calls itself on the group.
// If the condition is a condition, it calls the function [f] with the input and the condition.
// If the operator is AND, it returns an error if any condition fails.
// If the operator is OR, it collects all errors and returns them if all conditions fail.
func RunFuncOnConditionGroup[T comparable](input T, astValue *model.AstValue) error {
	var errors []error
	for i, v := range astValue.ConditionGroup {
		var err error
		switch v.Type {
		case model.EMPTY:
			return nil
		case model.GROUP:
			err = RunFuncOnConditionGroup(input, v)
		case model.CONDITION:
			switch v.ConditionType {
			case model.NONE:
				continue
			case model.EQUAL:
				err = ValidateEqual(input, v)
			case model.NOT_EQUAL:
				err = ValidateNotEqual(input, v)
			case model.MIN_VALUE:
				err = ValidateMin(input, v)
			case model.MAX_VALUE:
				err = ValidateMax(input, v)
			case model.CONTAINS:
				err = ValidateContains(input, v)
			case model.NOT_CONTAINS:
				err = ValidateNotContains(input, v)
			case model.FROM:
				err = ValidateFrom(input, v)
			case model.NOT_FROM:
				err = ValidateNotFrom(input, v)
			case model.REGX:
				err = ValidateRegex(input, v)
			default:
				return fmt.Errorf("unknown condition type: %v", v.ConditionType)
			}
		}
		if err != nil {
			if (i == 0 && v.Operator == model.OR) || (i > 0 && astValue.ConditionGroup[i-1].Operator == model.OR) {
				errors = append(errors, err)
			} else {
				return err
			}
		}
	}
	if len(astValue.ConditionGroup) > 0 && len(errors) >= len(astValue.ConditionGroup) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}
	return nil
}

package validators

import (
	"fmt"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func ValidateGroup(groups map[string]string, groupSize map[string]int, groupErrors map[string][]error) error {
	if len(groups) != 0 {
		for groupName, groupCondition := range groups {
			conType := model.GetConditionType(groupCondition)

			switch conType {
			case model.MIN_VALUE:
				condition, err := model.GetConditionByType(groupCondition, model.MIN_VALUE)
				if err != nil {
					return err
				}
				if len(condition) != 0 {
					minValue, err := strconv.Atoi(condition)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) < minValue {
						return fmt.Errorf("less then %v in group %s without error, all errors: %v", minValue, groupName, groupErrors[groupName])
					}
				}
			case model.MAX_VLAUE:
				condition, err := model.GetConditionByType(groupCondition, model.MAX_VLAUE)
				if err != nil {
					return err
				}
				if len(condition) != 0 {
					maxValue, err := strconv.Atoi(condition)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) > maxValue {
						return fmt.Errorf("more then %v in group %s without error, all errors: %v", maxValue, groupName, groupErrors[groupName])
					}
				}
			default:
				return fmt.Errorf("invalid group condition type %s", conType)
			}
		}
	}
	return nil
}

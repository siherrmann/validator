package validator

import (
	"fmt"
	"strconv"
)

func validateGroup(groups map[string]string, groupSize map[string]int, groupErrors map[string][]error) error {
	if len(groups) != 0 {
		for groupName, groupCondition := range groups {
			conType := getConditionType(groupCondition)

			switch conType {
			case MIN_VALUE:
				condition, err := getConditionByType(groupCondition, MIN_VALUE)
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
			case MAX_VLAUE:
				condition, err := getConditionByType(groupCondition, MAX_VLAUE)
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

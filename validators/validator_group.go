package validators

import (
	"fmt"
	"strconv"

	"github.com/siherrmann/validator/model"
)

func ValidateGroups(groups map[string]*model.Group, groupSize map[string]int, groupErrors map[string][]error) error {
	if len(groups) != 0 {
		for groupName, group := range groups {
			switch group.ConditionType {
			case model.MIN_VALUE:
				if len(group.ConditionValue) != 0 {
					minValue, err := strconv.Atoi(group.ConditionValue)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) < minValue {
						return fmt.Errorf("less then %v in group %s without error, all errors: %v", minValue, groupName, groupErrors[groupName])
					}
				}
			case model.MAX_VLAUE:
				if len(group.ConditionValue) != 0 {
					maxValue, err := strconv.Atoi(group.ConditionValue)
					if err != nil {
						return err
					} else if (groupSize[groupName] - len(groupErrors[groupName])) > maxValue {
						return fmt.Errorf("more then %v in group %s without error, all errors: %v", maxValue, groupName, groupErrors[groupName])
					}
				}
			default:
				return fmt.Errorf("invalid group condition type %s", group.ConditionType)
			}
		}
	}
	return nil
}

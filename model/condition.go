package model

import (
	"fmt"
	"slices"
	"strings"
)

func GetConditionType(s string) (ConditionType, error) {
	var conditionType ConditionType
	if len(s) > 2 {
		conditionType = ConditionType(s[:3])
	} else {
		conditionType = ConditionType(s)
	}

	if _, ok := ValidConditionTypes[conditionType]; !ok {
		return conditionType, fmt.Errorf("invalid condition type: %s", conditionType)
	}
	return conditionType, nil
}

func GetConditionByType(conditionFull string, conditionType ConditionType) (string, error) {
	if len(conditionType) != 3 {
		return "", fmt.Errorf("length of conditionType has to be 3: %s", conditionType)
	}
	condition := strings.TrimPrefix(conditionFull, string(conditionType))
	if len(condition) == 0 {
		return "", fmt.Errorf("empty %s value", conditionType)
	}
	return condition, nil
}

func GetArrayFromCondition(condition string) ([]string, error) {
	conditionList := strings.Split(condition, ",")
	if len(conditionList) == 0 {
		return []string{}, fmt.Errorf("empty condition list %s value", condition)
	}
	return conditionList, nil
}

func GetConditionsAndOrFromString(in string) ([]string, bool) {
	or := false
	conditions := strings.Split(in, " ")
	if slices.Contains(conditions, string(OR)) {
		conditions = slices.DeleteFunc(conditions, func(v string) bool {
			return v == string(OR)
		})
		or = true
	}
	return conditions, or
}

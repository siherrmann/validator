package model

import (
	"fmt"
	"slices"
	"strings"
)

const (
	NONE         string = "-"
	EQUAL        string = "equ"
	NOT_EQUAL    string = "neq"
	MIN_VALUE    string = "min"
	MAX_VLAUE    string = "max"
	CONTAINS     string = "con"
	NOT_CONTAINS string = "nco"
	FROM         string = "frm"
	NOT_FROM     string = "nfr"
	REGX         string = "rex"
	OR           string = "||"
)

func GetConditionType(s string) string {
	if len(s) > 2 {
		return s[:3]
	}
	return s
}

func GetConditionByType(conditionFull string, conditionType string) (string, error) {
	if len(conditionType) != 3 {
		return "", fmt.Errorf("length of conditionType has to be 3: %s", conditionType)
	}
	condition := strings.TrimPrefix(conditionFull, conditionType)
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
	if slices.Contains(conditions, OR) {
		conditions = slices.DeleteFunc(conditions, func(v string) bool {
			return v == OR
		})
		or = true
	}
	return conditions, or
}

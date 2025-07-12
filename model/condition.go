package model

import (
	"fmt"
	"slices"
	"strings"
)

// ConditionType is the type for all available condition types.
type ConditionType string

// Available condition types.
const (
	NONE         ConditionType = "-"
	EQUAL        ConditionType = "equ"
	NOT_EQUAL    ConditionType = "neq"
	MIN_VALUE    ConditionType = "min"
	MAX_VLAUE    ConditionType = "max"
	CONTAINS     ConditionType = "con"
	NOT_CONTAINS ConditionType = "nco"
	FROM         ConditionType = "frm"
	NOT_FROM     ConditionType = "nfr"
	REGX         ConditionType = "rex"
)

var ValidConditionTypes = map[ConditionType]int{
	NONE:         0,
	EQUAL:        1,
	NOT_EQUAL:    2,
	MIN_VALUE:    3,
	MAX_VLAUE:    4,
	CONTAINS:     5,
	NOT_CONTAINS: 6,
	FROM:         7,
	NOT_FROM:     8,
	REGX:         9,
}

// LookupConditionType checks our validConditionType map for the scanned condition type.
// If not found, an error is returned.
func LookupConditionType(conType ConditionType) error {
	if _, ok := ValidConditionTypes[conType]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid condition type, found: %s", conType)
}

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

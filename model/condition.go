package model

import (
	"fmt"
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
	MAX_VALUE    ConditionType = "max"
	CONTAINS     ConditionType = "con"
	NOT_CONTAINS ConditionType = "nco"
	FROM         ConditionType = "frm"
	NOT_FROM     ConditionType = "nfr"
	REGX         ConditionType = "rex"
	FUNC         ConditionType = "fun"
)

var ValidConditionTypes = map[ConditionType]int{
	NONE:         0,
	EQUAL:        1,
	NOT_EQUAL:    2,
	MIN_VALUE:    3,
	MAX_VALUE:    4,
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

// GetConditionType returns the condition type from a string.
// It checks if the string starts with a valid condition type prefix.
// If the string is not valid, an error is returned.
func GetConditionType(s string) (ConditionType, error) {
	var conditionType ConditionType
	if len(s) > 3 {
		conditionType = ConditionType(s[:3])
	} else {
		conditionType = ConditionType(s)
	}

	if _, ok := ValidConditionTypes[conditionType]; !ok {
		return conditionType, fmt.Errorf("invalid condition type: %s", conditionType)
	}
	return conditionType, nil
}

// GetConditionByType extracts the condition value from a string based on the condition type.
// It trims the prefix of the condition type from the string and returns the remaining part.
// If the condition type is not valid or the value is empty, an error is returned.
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

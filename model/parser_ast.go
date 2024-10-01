package model

import (
	"fmt"
	"strings"
)

// RootNode is what starts every parsed AST.
type RootNode struct {
	RootValue *Value
}

// Value holds a Type ("Condition" or "Group") as well as a `ConditionType` and `ConditionValue`.
// The ConditionType is a [model.ConditionType] and the ConditionValue is any string (numbers are also represented as string).
type Value struct {
	Type           string // "Condition"
	ConditionType  ConditionType
	ConditionValue string
	ConditionGroup []*Value
	Operator       Operator
	Start          int
	End            int
}

func (r Value) AstGroupToString() string {
	groupConditions := []string{}
	groupString := ""
	for _, v := range r.ConditionGroup {
		if v.Type == "Group" {
			groupConditions = append(groupConditions, "("+v.AstGroupToString()+")")
		} else if v.Type == "Condition" {
			groupConditions = append(groupConditions, v.AstConditionToString())
		}
	}
	groupString = strings.Join(groupConditions, " ")
	return groupString
}

func (r Value) AstConditionToString() string {
	if len(r.Operator) > 0 {
		return fmt.Sprintf(`%v%v "%v"`, r.ConditionType, r.ConditionValue, r.Operator)
	} else {
		return fmt.Sprintf(`%v"%v"`, r.ConditionType, r.ConditionValue)
	}
}

// [ConditionType] is the type for all available condition types.
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
	AND          ConditionType = "&&"
	OR           ConditionType = "||"
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
	OR:           10,
}

// [LookupConditionType] checks our validConditionType map for the scanned condition type.
// If not found, an error is returned.
func LookupConditionType(conType ConditionType) error {
	if _, ok := ValidConditionTypes[conType]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid condition type, found: %s", conType)
}

// [Operator] is the type for all available operators.
type Operator string

// Available operators.
const (
	// Group states
	And Operator = "&&"
	Or  Operator = "||"
)

var validOperator = map[Operator]int{
	And: 0,
	Or:  1,
}

// [LookupOperator] checks our validOperator map for the scanned operator.
// If not found, an error is returned.
func LookupOperator(operator Operator) error {
	if _, ok := validOperator[operator]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid operator, found: %s", operator)
}

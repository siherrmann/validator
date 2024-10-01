package model

import (
	"fmt"
	"reflect"
	"strings"
)

// RootNode is what starts every parsed AST.
type RootNode struct {
	RootValue *AstValue
}

// AstValue holds a Type ("Condition" or "Group") as well as a `ConditionType` and `ConditionValue`.
// The ConditionType is a [model.ConditionType] and the ConditionValue is any string (numbers are also represented as string).
type AstValue struct {
	Type           AstValueType
	ConditionType  ConditionType
	ConditionValue string
	ConditionGroup ConditionGroup
	Operator       Operator
	Start          int
	End            int
}

type AstValueType string

const (
	GROUP     AstValueType = "Group"
	CONDITION AstValueType = "Condition"
)

type ConditionGroup []*AstValue

func (r AstValue) AstGroupToString() string {
	groupConditions := []string{}
	groupString := ""
	for _, v := range r.ConditionGroup {
		if v.Type == GROUP {
			if len(v.Operator) > 0 {
				groupConditions = append(groupConditions, fmt.Sprintf("(%v) %v", v.AstGroupToString(), v.Operator))
			} else {
				groupConditions = append(groupConditions, fmt.Sprintf("(%v)", v.AstGroupToString()))
			}

		} else if v.Type == CONDITION {
			groupConditions = append(groupConditions, v.AstConditionToString())
		}
	}
	groupString = strings.Join(groupConditions, " ")
	return groupString
}

func (r AstValue) AstConditionToString() string {
	if len(r.Operator) > 0 {
		return fmt.Sprintf("%v'%v' %v", r.ConditionType, r.ConditionValue, r.Operator)
	} else {
		return fmt.Sprintf("%v'%v'", r.ConditionType, r.ConditionValue)
	}
}

func (r AstValue) RunFuncOnConditionGroup(input reflect.Value, f func(reflect.Value, *AstValue) error) error {
	var errors []error
	for i, v := range r.ConditionGroup {
		var err error
		if v.Type == GROUP {
			err = v.RunFuncOnConditionGroup(input, f)
		} else if v.Type == CONDITION {
			err = f(input, v)
		}
		if err != nil && i == 0 && v.Operator == OR {
			errors = append(errors, err)
		} else if err != nil && i == 0 && v.Operator == AND {
			return err
		} else if err != nil && i > 0 && r.ConditionGroup[i-1].Operator == OR {
			errors = append(errors, err)
		} else if err != nil {
			return err
		}
	}
	if len(errors) >= len(r.ConditionGroup) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}
	return nil
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
	AND Operator = "&&"
	OR  Operator = "||"
)

var validOperator = map[Operator]int{
	AND: 0,
	OR:  1,
}

// [LookupOperator] checks our validOperator map for the scanned operator.
// If not found, an error is returned.
func LookupOperator(operator Operator) error {
	if _, ok := validOperator[operator]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid operator, found: %s", operator)
}

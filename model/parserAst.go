package model

import (
	"fmt"
	"strings"
)

// RootNode is what starts every parsed AST (=abstract syntax tree).
type RootNode struct {
	RootValue *AstValue
}

// AstValue (=abstract syntax tree value) holds a Type ("Condition" or "Group") as well as a `ConditionType` and `ConditionValue`.
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
	EMPTY     AstValueType = "Empty"
	GROUP     AstValueType = "Group"
	CONDITION AstValueType = "Condition"
)

type ConditionGroup []*AstValue

func (r AstValue) AstGroupToString() string {
	groupConditions := []string{}
	groupString := ""
	for _, v := range r.ConditionGroup {
		switch v.Type {
		case GROUP:
			if len(v.Operator) > 0 {
				groupConditions = append(groupConditions, fmt.Sprintf("(%v) %v", v.AstGroupToString(), v.Operator))
			} else {
				groupConditions = append(groupConditions, fmt.Sprintf("(%v)", v.AstGroupToString()))
			}
		case CONDITION:
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

// Operator is the type for all available operators.
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

// LookupOperator checks our validOperator map for the scanned operator.
// If not found, an error is returned.
func LookupOperator(operator Operator) error {
	if _, ok := validOperator[operator]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid operator, found: %s", operator)
}

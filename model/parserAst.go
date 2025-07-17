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

// RunFuncOnConditionGroup runs the function [f] on each condition in the [astValue].
// If the condition is a group, it recursively calls itself on the group.
// If the condition is a condition, it calls the function [f] with the input and the condition.
// If the operator is AND, it returns an error if any condition fails.
// If the operator is OR, it collects all errors and returns them if all conditions fail.
func RunFuncOnConditionGroup[T comparable](input T, astValue *AstValue, f func(T, *AstValue) error) error {
	var errors []error
	for i, v := range astValue.ConditionGroup {
		var err error
		switch v.Type {
		case EMPTY:
			return nil
		case GROUP:
			err = RunFuncOnConditionGroup(input, v, f)
		case CONDITION:
			err = f(input, v)
		}
		if err != nil && i == 0 && v.Operator == OR {
			errors = append(errors, err)
		} else if err != nil && i == 0 && v.Operator == AND {
			return err
		} else if err != nil && i > 0 && astValue.ConditionGroup[i-1].Operator == OR {
			errors = append(errors, err)
		} else if err != nil {
			return err
		}
	}
	if len(astValue.ConditionGroup) > 0 && len(errors) >= len(astValue.ConditionGroup) {
		return fmt.Errorf("no condition fulfilled, all errors: %v", errors)
	}
	return nil
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

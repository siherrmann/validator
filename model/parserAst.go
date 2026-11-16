package model

import (
	"fmt"
	"regexp"
	"strings"
)

// RootNode is what starts every parsed AST (=abstract syntax tree).
type RootNode struct {
	RootValue *AstValue
}

// ExtractDefault searches the AST for a 'DEF' condition and returns its value.
func (r *RootNode) ExtractDefault() string {
	if r.RootValue == nil {
		return ""
	}
	return r.RootValue.ExtractDefault()
}

// AstValue (=abstract syntax tree value) holds a Type ("Condition" or "Group") as well as a `ConditionType` and `ConditionValue`.
// The ConditionType is a [model.ConditionType] and the ConditionValue is any string (numbers are also represented as string).
//
// A `rel` condition holds the path of the related field in Ref and the conditions
// for the related field in ConditionGroup. A condition with a `$` value (eg. `min$min_replicas`
// or the injected value `equ${api_key}`) holds the path of the referenced field in Ref,
// its value is resolved while validating.
//
// The AST is cached and shared between goroutines, it must not be modified while validating.
type AstValue struct {
	Type           AstValueType
	ConditionType  ConditionType
	ConditionValue string
	ConditionGroup ConditionGroup
	Operator       Operator
	Start          int
	End            int
	// Ref is the path of the related (`rel`) or referenced (`$`) field.
	Ref *FieldRef
	// Regex is the compiled regular expression of a `rex` condition with a fixed value.
	Regex *regexp.Regexp
	// List is the split value of a `frm` or `nfr` condition with a fixed value.
	List []string
	// Number is the parsed value of a `min` or `max` condition, if HasNumber is set.
	Number    float64
	HasNumber bool
	// UsesFields reports if the AST contains a `rel` condition or a `$` reference.
	// It is only set on the root value by the parser.
	UsesFields bool
}

// HasFieldUsage recursively checks if the AST value contains a `rel` condition or a `$` reference.
func (r *AstValue) HasFieldUsage() bool {
	if r == nil {
		return false
	}
	for _, v := range r.ConditionGroup {
		if v.Ref != nil {
			return true
		}
		if v.Type == GROUP && v.HasFieldUsage() {
			return true
		}
	}
	return false
}

// ExtractDefault recursively searches the AST value for a 'DEF' condition and returns its value.
func (r *AstValue) ExtractDefault() string {
	for _, v := range r.ConditionGroup {
		if v.Type == CONDITION && v.ConditionType == DEF {
			return v.ConditionValue
		}
		if v.Type == GROUP {
			if def := v.ExtractDefault(); def != "" {
				return def
			}
		}
	}
	return ""
}

// HasCondition recursively checks if the AST value contains a condition for which [match] returns true.
func (r *AstValue) HasCondition(match func(ConditionType) bool) bool {
	if r == nil {
		return false
	}
	for _, v := range r.ConditionGroup {
		if v.Type == CONDITION && match(v.ConditionType) {
			return true
		}
		if v.Type == GROUP && v.HasCondition(match) {
			return true
		}
	}
	return false
}

// AstValueType is the type for all available AST value types.
type AstValueType string

const (
	EMPTY     AstValueType = "Empty"
	GROUP     AstValueType = "Group"
	CONDITION AstValueType = "Condition"
)

// ConditionGroup is a slice of AstValue pointers.
type ConditionGroup []*AstValue

// AstGroupToString converts the AstValue's ConditionGroup to a string representation.
// It iterates over each AstValue in the group and formats it based on its type.
// If the AstValue is a group, it recursively calls itself to get the string representation of the group.
// If the AstValue is a condition, it formats it as a string with its ConditionType and ConditionValue.
// The resulting string is a concatenation of all conditions and groups, separated by spaces.
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

// AstConditionToString converts the AstValue's ConditionType and ConditionValue to a string representation.
// If the AstValue has an Operator, it includes that in the string.
// The resulting string is formatted as "<ConditionType>'<ConditionValue>' <Operator>" if the Operator is present,
// or as "<ConditionType>'<ConditionValue>'" if the Operator is not present.
// A `rel` condition is formatted as "rel(<path> <conditions>)", a `$` reference as "<ConditionType>$<path>"
// and an injected value as "<ConditionType>${<name>}".
func (r AstValue) AstConditionToString() string {
	var condition string
	switch {
	case r.ConditionType == RELATION && r.Ref != nil:
		condition = fmt.Sprintf("rel(%v %v)", r.Ref.Raw, r.AstGroupToString())
	case r.Ref != nil && r.Ref.Injected:
		condition = fmt.Sprintf("%v%v", r.ConditionType, r.Ref.Raw)
	case r.Ref != nil:
		condition = fmt.Sprintf("%v$%v", r.ConditionType, r.Ref.Raw)
	default:
		condition = fmt.Sprintf("%v'%v'", r.ConditionType, r.ConditionValue)
	}
	if len(r.Operator) > 0 {
		return condition + " " + string(r.Operator)
	}
	return condition
}

// Operator is the type for all available operators.
type Operator string

// Available operators.
const (
	// Group states
	AND     Operator = "&&"
	OR      Operator = "||"
	IMPLIES Operator = "=>"
)

// validOperator is a map that holds valid operators and their corresponding integer values.
// This map is used to validate operators in the AST.
var validOperator = map[Operator]int{
	AND:     0,
	OR:      1,
	IMPLIES: 2,
}

// LookupOperator checks our validOperator map for the scanned operator.
// If not found, an error is returned.
func LookupOperator(operator Operator) error {
	if _, ok := validOperator[operator]; ok {
		return nil
	}
	return fmt.Errorf("expected a valid operator, found: %s", operator)
}

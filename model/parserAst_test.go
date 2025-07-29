package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAstGroupToString(t *testing.T) {
	tests := []struct {
		name     string
		astValue AstValue
		expected string
	}{
		{
			name:     "Valid condition",
			astValue: AstValue{ConditionGroup: ConditionGroup{&AstValue{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "1"}}},
			expected: "equ'1'",
		},
		{
			name: "Valid group",
			astValue: AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: GROUP, ConditionGroup: ConditionGroup{
						&AstValue{Type: CONDITION, ConditionType: MIN_VALUE, ConditionValue: "2", Operator: AND},
						&AstValue{Type: CONDITION, ConditionType: MAX_VALUE, ConditionValue: "10"},
					}},
				},
			},
			expected: "(min'2' && max'10')",
		},
		{
			name: "Valid groups",
			astValue: AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: GROUP, Operator: OR, ConditionGroup: ConditionGroup{
						&AstValue{Type: CONDITION, ConditionType: MIN_VALUE, ConditionValue: "2", Operator: AND},
						&AstValue{Type: CONDITION, ConditionType: MAX_VALUE, ConditionValue: "10"},
					}},
					&AstValue{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "0"},
				},
			},
			expected: "(min'2' && max'10') || equ'0'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			groupString := test.astValue.AstGroupToString()
			assert.NotEmpty(t, groupString, "Expected group string to be not empty")
			assert.Equal(t, test.expected, groupString, "Expected group string to match")
		})
	}
}

func TestLookupOperator(t *testing.T) {
	tests := []struct {
		name     string
		operator Operator
		expected error
	}{
		{
			name:     "Valid AND operator",
			operator: AND,
			expected: nil,
		},
		{
			name:     "Valid OR operator",
			operator: OR,
			expected: nil,
		},
		{
			name:     "Invalid operator",
			operator: Operator("invalid"),
			expected: fmt.Errorf("expected a valid operator, found: invalid"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			valid := LookupOperator(test.operator)
			assert.Equal(t, test.expected, valid, "Expected operator validity to match")
		})
	}
}

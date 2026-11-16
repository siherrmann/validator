package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDefault(t *testing.T) {
	tests := []struct {
		name     string
		astValue AstValue
		expected string
	}{
		{
			name: "No default",
			astValue: AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "1"},
				},
			},
			expected: "",
		},
		{
			name: "With default",
			astValue: AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: CONDITION, ConditionType: DEF, ConditionValue: "test"},
				},
			},
			expected: "test",
		},
		{
			name: "With default in group",
			astValue: AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: GROUP, ConditionGroup: ConditionGroup{
						&AstValue{Type: CONDITION, ConditionType: DEF, ConditionValue: "nested"},
					}},
				},
			},
			expected: "nested",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.astValue.ExtractDefault()
			assert.Equal(t, test.expected, result, "Expected default string to match")
		})
	}
}

func TestHasCondition(t *testing.T) {
	isMin := func(conditionType ConditionType) bool { return conditionType == MIN_VALUE }
	tests := []struct {
		name     string
		astValue *AstValue
		expected bool
	}{
		{
			name:     "Nil ast value",
			astValue: nil,
			expected: false,
		},
		{
			name: "No matching condition",
			astValue: &AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "1"},
				},
			},
			expected: false,
		},
		{
			name: "Matching condition",
			astValue: &AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "1"},
					&AstValue{Type: CONDITION, ConditionType: MIN_VALUE, ConditionValue: "1"},
				},
			},
			expected: true,
		},
		{
			name: "Matching condition in group",
			astValue: &AstValue{
				ConditionGroup: ConditionGroup{
					&AstValue{Type: GROUP, ConditionGroup: ConditionGroup{
						&AstValue{Type: CONDITION, ConditionType: MIN_VALUE, ConditionValue: "1"},
					}},
				},
			},
			expected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.astValue.HasCondition(isMin)
			assert.Equal(t, test.expected, result, "Expected HasCondition result to match")
		})
	}
}

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
			name:     "Valid IMPLIES operator",
			operator: IMPLIES,
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

func TestHasFieldUsage(t *testing.T) {
	ref := &FieldRef{Raw: "a", Path: []string{"a"}}
	assert.False(t, (*AstValue)(nil).HasFieldUsage())
	assert.False(t, (&AstValue{ConditionGroup: ConditionGroup{{Type: CONDITION, ConditionType: MIN_VALUE}}}).HasFieldUsage())
	assert.True(t, (&AstValue{ConditionGroup: ConditionGroup{{Type: CONDITION, ConditionType: MIN_VALUE, Ref: ref}}}).HasFieldUsage())
	assert.True(t, (&AstValue{ConditionGroup: ConditionGroup{{Type: GROUP, ConditionGroup: ConditionGroup{{Type: CONDITION, ConditionType: RELATION, Ref: ref}}}}}).HasFieldUsage())
}

func TestAstConditionToStringWithFields(t *testing.T) {
	ref := &FieldRef{Raw: "../type", Up: 1, Path: []string{"type"}}
	relation := AstValue{
		Type:           CONDITION,
		ConditionType:  RELATION,
		Ref:            ref,
		ConditionGroup: ConditionGroup{{Type: CONDITION, ConditionType: EQUAL, ConditionValue: "oci"}},
		Operator:       IMPLIES,
	}
	assert.Equal(t, "rel(../type equ'oci') =>", relation.AstConditionToString())

	reference := AstValue{Type: CONDITION, ConditionType: MIN_VALUE, Ref: ref}
	assert.Equal(t, "min$../type", reference.AstConditionToString())
}

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	parser := NewParser()

	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Valid input",
			input:    "equvalid",
			expected: "equ'valid'",
			wantErr:  false,
		},
		{
			name:     "Valid input with string",
			input:    "equ'valid'",
			expected: "equ'valid'",
			wantErr:  false,
		},
		{
			name:     "Valid input with newline",
			input:    "equ'valid'\n|| equ'valid2'",
			expected: "equ'valid' || equ'valid2'",
			wantErr:  false,
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "Empty condition",
			input:    "-",
			expected: "",
			wantErr:  false,
		},
		{
			name:     "Invalid condition 1",
			input:    "invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid condition 2",
			input:    "1",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty condition value",
			input:    "(equ)",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid group start",
			input:    "||",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid group end",
			input:    "(equapple || invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid second condition",
			input:    "equapple invalid",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Missing second condition in group",
			input:    "equapple && (|| equbanana)",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Invalid operator",
			input:    "equapple &| equbanana",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Condition group without operator",
			input:    "min1 max2",
			expected: "min'1' && max'2'",
			wantErr:  false,
		},
		{
			name:     "Condition group with AND operator",
			input:    "min1 && max2",
			expected: "min'1' && max'2'",
			wantErr:  false,
		},
		{
			name:     "Condition group with OR operator",
			input:    "min1 || max2",
			expected: "min'1' || max'2'",
			wantErr:  false,
		},
		{
			name:     "Condition groups with brackets without operator",
			input:    "(min1 || max2) (min3 || max4)",
			expected: "(min'1' || max'2') && (min'3' || max'4')",
			wantErr:  false,
		},
		{
			name:     "Condition groups with brackets and AND operator",
			input:    "(min1 || max2) && (min3 || max4)",
			expected: "(min'1' || max'2') && (min'3' || max'4')",
			wantErr:  false,
		},
		{
			name:     "Condition groups with brackets and OR operator",
			input:    "(min1 && max2) || (min3 && max4)",
			expected: "(min'1' && max'2') || (min'3' && max'4')",
			wantErr:  false,
		},
		{
			name:     "Condition group with nested conditions",
			input:    "(min1 && (max2 || min3)) || (min4 && max5)",
			expected: "(min'1' && (max'2' || min'3')) || (min'4' && max'5')",
			wantErr:  false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rootNode, err := parser.ParseValidation(test.input)
			if test.wantErr {
				assert.Error(t, err, "Expected an error for input: %s", test.input)
			} else {
				assert.NoError(t, err, "Did not expect an error for input: %s", test.input)
				assert.NotNil(t, rootNode.RootValue, "Expected non-nil RootValue in AST")
				assert.Len(t, parser.Errors(), 0, "Expected empty ConditionGroup in AST")
				assert.Equal(t, test.expected, rootNode.RootValue.AstGroupToString(), "Expected ConditionType to match")
			}
		})
	}
}

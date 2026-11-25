package parser

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLexer(t *testing.T) {
	lexer := NewLexer("min1")
	require.NotNil(t, lexer)
	assert.Equal(t, "min1", lexer.Input)

	token := lexer.NextToken()
	assert.Equal(t, model.LexerConditionType, token.Type)
	assert.Equal(t, "min", token.Literal)
}

func TestNextToken(t *testing.T) {
	type token struct {
		Type    model.TokenType
		Literal string
	}
	tests := []struct {
		name     string
		input    string
		expected []token
	}{
		{
			name:  "Condition",
			input: "min10",
			expected: []token{
				{model.LexerConditionType, "min"},
				{model.LexerConditionValue, "10"},
			},
		},
		{
			name:  "Quoted value",
			input: "equ'a b'",
			expected: []token{
				{model.LexerConditionType, "equ"},
				{model.LexerConditionValueString, "a b"},
			},
		},
		{
			name:  "Group and operators",
			input: "(min1 || max2) && equ3 => neq4",
			expected: []token{
				{model.LexerLeftBrace, "("},
				{model.LexerConditionType, "min"},
				{model.LexerConditionValue, "1"},
				{model.LexerOperator, "||"},
				{model.LexerConditionType, "max"},
				{model.LexerConditionValue, "2"},
				{model.LexerRightBrace, ")"},
				{model.LexerOperator, "&&"},
				{model.LexerConditionType, "equ"},
				{model.LexerConditionValue, "3"},
				{model.LexerOperator, "=>"},
				{model.LexerConditionType, "neq"},
				{model.LexerConditionValue, "4"},
			},
		},
		{
			name:  "Relation with field path",
			input: "rel(../type equoci)",
			expected: []token{
				{model.LexerConditionType, "rel"},
				{model.LexerLeftBrace, "("},
				{model.LexerFieldPath, "../type"},
				{model.LexerConditionType, "equ"},
				{model.LexerConditionValue, "oci"},
				{model.LexerRightBrace, ")"},
			},
		},
		{
			name:  "Reference and injected value",
			input: "min$other || equ${api_key}",
			expected: []token{
				{model.LexerConditionType, "min"},
				{model.LexerConditionValue, "$other"},
				{model.LexerOperator, "||"},
				{model.LexerConditionType, "equ"},
				{model.LexerConditionValue, "${api_key}"},
			},
		},
		{
			name:  "Operator does not consume the next character",
			input: "min1 &&max2 &&& equ3",
			expected: []token{
				{model.LexerConditionType, "min"},
				{model.LexerConditionValue, "1"},
				{model.LexerOperator, "&&"},
				{model.LexerConditionType, "max"},
				{model.LexerConditionValue, "2"},
				{model.LexerOperator, "&&"},
				{model.LexerOperator, "&"},
				{model.LexerConditionType, "equ"},
				{model.LexerConditionValue, "3"},
			},
		},
		{
			name:  "Braces in a value",
			input: "(rex^(a|b)+$ || rex(c))",
			expected: []token{
				{model.LexerLeftBrace, "("},
				{model.LexerConditionType, "rex"},
				{model.LexerConditionValue, "^(a|b)+$"},
				{model.LexerOperator, "||"},
				{model.LexerConditionType, "rex"},
				{model.LexerConditionValue, "(c)"},
				{model.LexerRightBrace, ")"},
			},
		},
		{
			name:  "Negative value",
			input: "min-5 && neq -",
			expected: []token{
				{model.LexerConditionType, "min"},
				{model.LexerConditionValue, "-5"},
				{model.LexerOperator, "&&"},
				{model.LexerConditionType, "neq"},
				{model.LexerEmptyRequirement, "-"},
			},
		},
		{
			name:     "Empty requirement",
			input:    "-",
			expected: []token{{model.LexerEmptyRequirement, "-"}},
		},
		{
			name:     "Illegal character",
			input:    "1",
			expected: []token{{model.LexerIllegal, "1"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			lexer := NewLexer(test.input)
			for _, expected := range test.expected {
				actual := lexer.NextToken()
				assert.Equal(t, expected, token{actual.Type, actual.Literal})
			}
			if test.expected[len(test.expected)-1].Type != model.LexerIllegal {
				assert.Equal(t, model.LexerEOF, lexer.NextToken().Type)
			}
		})
	}
}

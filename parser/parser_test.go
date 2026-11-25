package parser

import (
	"testing"

	"github.com/siherrmann/validator/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseValidation(t *testing.T) {
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
		{
			name:     "Implication",
			input:    "equoci => min1",
			expected: "equ'oci' => min'1'",
			wantErr:  false,
		},
		{
			name:     "Implication binds weaker than AND and OR and is right associative",
			input:    "min1 && max5 => equ3 || equ4 => min2",
			expected: "(min'1' && max'5') => ((equ'3' || equ'4') => min'2')",
			wantErr:  false,
		},
		{
			name:     "Implication without consequence",
			input:    "min1 =>",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Relation",
			input:    "rel(type equoci) => min1",
			expected: "rel(type equ'oci') => min'1'",
			wantErr:  false,
		},
		{
			name:     "Relation with parent path and group",
			input:    "rel(../stateful (equtrue || equ1)) => max0",
			expected: "rel(../stateful (equ'true' || equ'1')) => max'0'",
			wantErr:  false,
		},
		{
			name:     "Relation in group",
			input:    "(rel(/stateful equfalse) => max0) && unqname",
			expected: "(rel(/stateful equ'false') => max'0') && unq'name'",
			wantErr:  false,
		},
		{
			name:     "Relation without brace",
			input:    "rel type",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Relation without conditions",
			input:    "rel(type)",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Relation without closing brace",
			input:    "rel(type equoci",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Relation with invalid path",
			input:    "rel(/../type equoci)",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Reference",
			input:    "min$min_replicas && frm$../ports.port",
			expected: "min$min_replicas && frm$../ports.port",
			wantErr:  false,
		},
		{
			name:     "Quoted value is not a reference",
			input:    "equ'$literal'",
			expected: "equ'$literal'",
			wantErr:  false,
		},
		{
			name:     "Reference without field name",
			input:    "min$../",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Reference on fun",
			input:    "fun$check",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Injected value",
			input:    "equ${api_key} || frm${allowed_ids}",
			expected: "equ${api_key} || frm${allowed_ids}",
			wantErr:  false,
		},
		{
			name:     "Relation with injected value",
			input:    "rel(${role} equadmin) => min1",
			expected: "rel(${role} equ'admin') => min'1'",
			wantErr:  false,
		},
		{
			name:     "Injected value without name",
			input:    "equ${}",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Injected value without closing brace",
			input:    "equ${api_key",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Quoted injected value is a fixed value",
			input:    "equ'${api_key}'",
			expected: "equ'${api_key}'",
			wantErr:  false,
		},
		{
			name:     "Unique without value",
			input:    "unq && min1",
			expected: "unq'' && min'1'",
			wantErr:  false,
		},
		{
			name:     "Unique with key",
			input:    "unqname",
			expected: "unq'name'",
			wantErr:  false,
		},
		{
			name:     "Invalid regex",
			input:    "rex'(['",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Groups after the requirement are not part of it",
			input:    `min1`,
			expected: `min'1'`,
			wantErr:  false,
		},
		{
			name:     "Default",
			input:    `equapple || defunassigned`,
			expected: `equ'apple' || def'unassigned'`,
			wantErr:  false,
		},
		{
			name:     "Frm list",
			input:    `frma,b,c`,
			expected: `frm'a,b,c'`,
			wantErr:  false,
		},
		{
			name:     "Frm list with spaces in quotes",
			input:    `frm'a b,c'`,
			expected: `frm'a b,c'`,
			wantErr:  false,
		},
		{
			name:     "Regex with braces",
			input:    `rex^(.*[A-Z])+(.*)$`,
			expected: `rex'^(.*[A-Z])+(.*)$'`,
			wantErr:  false,
		},
		{
			name:     "Regex with braces in a group",
			input:    `(rex^(a|b)$ || max0) && min1`,
			expected: `(rex'^(a|b)$' || max'0') && min'1'`,
			wantErr:  false,
		},
		{
			name:     "Regex starting with a brace",
			input:    `rex(a|b)`,
			expected: `rex'(a|b)'`,
			wantErr:  false,
		},
		{
			name:     "Regex with special characters",
			input:    `rex^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$`,
			expected: `rex'^(.*[\x60!@#$%^&*()_+={};':\"|\\,.<>/?~-])+(.*)$'`,
			wantErr:  false,
		},
		{
			name:     "Quoted regex",
			input:    `rex'^a b$'`,
			expected: `rex'^a b$'`,
			wantErr:  false,
		},
		{
			name:     "Negative numbers",
			input:    `min-5 && max-1.5`,
			expected: `min'-5' && max'-1.5'`,
			wantErr:  false,
		},
		{
			name:     "Operator without spaces",
			input:    `min1 &&max2`,
			expected: `min'1' && max'2'`,
			wantErr:  false,
		},
		{
			name:     "Is type",
			input:    `istuuid`,
			expected: `ist'uuid'`,
			wantErr:  false,
		},
		{
			name:     "Correct min and max",
			input:    `cmn1 && cmx2`,
			expected: `cmn'1' && cmx'2'`,
			wantErr:  false,
		},
		{
			name:     "Function",
			input:    `min3 && funcheck`,
			expected: `min'3' && fun'check'`,
			wantErr:  false,
		},
		{
			name:     "Value with space",
			input:    `min 5`,
			expected: `min'5'`,
			wantErr:  false,
		},
		{
			name:     "Trailing AND",
			input:    `min1 &&`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Trailing OR",
			input:    `min1 || `,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Trailing operator in group",
			input:    `(min1 &&) || max2`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Trailing operator in relation",
			input:    `rel(a equ1 &&) => min1`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Leading operator",
			input:    `&& min1`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Two operators",
			input:    `min1 && || max2`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Three character operator",
			input:    `min1 &&& max2`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Non-numeric min",
			input:    `minabc`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Non-numeric max",
			input:    `max'abc'`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Operator without space after value",
			input:    `min1&&max2`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Non-numeric cmn",
			input:    `cmnx`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Unknown ist type",
			input:    `istuuidd`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Closing brace without opening brace",
			input:    `max1)`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Closing brace in the middle",
			input:    `min1 ) ( max2`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Opening brace without closing brace",
			input:    `(min1`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty group",
			input:    `()`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Empty requirement with conditions",
			input:    `- min1`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Unknown condition type",
			input:    `abc1`,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Unknown condition type with uppercase",
			input:    `Min1`,
			expected: "",
			wantErr:  true,
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

	t.Run("Regex is compiled", func(t *testing.T) {
		root, err := NewParser().ParseValidation("rex^[A-Z]+$")
		assert.NoError(t, err)
		condition := root.RootValue.ConditionGroup[0]
		assert.NotNil(t, condition.Regex, "Expected compiled regex")
		assert.True(t, condition.Regex.MatchString("ABC"))
	})

	t.Run("Min and max values are parsed", func(t *testing.T) {
		root, err := NewParser().ParseValidation("min1.5 && max10")
		assert.NoError(t, err)
		assert.True(t, root.RootValue.ConditionGroup[0].HasNumber)
		assert.Equal(t, 1.5, root.RootValue.ConditionGroup[0].Number)
		assert.True(t, root.RootValue.ConditionGroup[1].HasNumber)
		assert.Equal(t, 10.0, root.RootValue.ConditionGroup[1].Number)
	})

	t.Run("Reference is not parsed as number", func(t *testing.T) {
		root, err := NewParser().ParseValidation("min$other")
		assert.NoError(t, err)
		condition := root.RootValue.ConditionGroup[0]
		assert.False(t, condition.HasNumber)
		assert.Equal(t, []string{"other"}, condition.Ref.Path)
	})

	t.Run("Error message names the requirement", func(t *testing.T) {
		_, err := NewParser().ParseValidation("min1 &&")
		assert.EqualError(t, err, `invalid requirement "min1 &&": error parsing group, expected a condition after &&`)
		assert.ErrorIs(t, err, model.ErrInvalidRequirement)
	})

	t.Run("Field usage is set on the root", func(t *testing.T) {
		for input, expected := range map[string]bool{
			"min1 && max2":                  false,
			"(min1 || rel(a equ1))":         true,
			"min$a":                         true,
			"equ'$a'":                       false,
			"rel(a equ1) => (min1 || max2)": true,
		} {
			root, err := NewParser().ParseValidation(input)
			assert.NoError(t, err, input)
			assert.Equal(t, expected, root.RootValue.UsesFields, input)
		}
	})
}

func TestNewParser(t *testing.T) {
	parser := NewParser()
	require.NotNil(t, parser)
	assert.Empty(t, parser.Errors())
}

func TestErrors(t *testing.T) {
	parser := NewParser()
	_, err := parser.ParseValidation("invalid")
	require.Error(t, err)
	assert.Contains(t, parser.Errors(), "expected a valid condition type")

	_, err = parser.ParseValidation("min1")
	require.NoError(t, err)
	assert.Empty(t, parser.Errors(), "Expected the errors to be reset for a new parse")
}

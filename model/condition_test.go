package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsRequirement(t *testing.T) {
	tests := []struct {
		name     string
		input    ConditionType
		expected bool
	}{
		{name: "Min condition type", input: MIN_VALUE, expected: true},
		{name: "Correct min condition type", input: CORRECT_MIN, expected: true},
		{name: "None condition type", input: NONE, expected: false},
		{name: "Default condition type", input: DEF, expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.input.IsRequirement(), "Expected IsRequirement result to match")
		})
	}
}

func TestIsCorrect(t *testing.T) {
	tests := []struct {
		name     string
		input    ConditionType
		expected bool
	}{
		{name: "Correct min condition type", input: CORRECT_MIN, expected: true},
		{name: "Correct max condition type", input: CORRECT_MAX, expected: true},
		{name: "Min condition type", input: MIN_VALUE, expected: false},
		{name: "None condition type", input: NONE, expected: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.expected, test.input.IsCorrect(), "Expected IsCorrect result to match")
		})
	}
}

func TestLookupConditionType(t *testing.T) {
	tests := []struct {
		name    string
		input   ConditionType
		wantErr bool
	}{
		{
			name:    "Valid condition type",
			input:   EQUAL,
			wantErr: false,
		},
		{
			name:    "Valid correct min condition type",
			input:   CORRECT_MIN,
			wantErr: false,
		},
		{
			name:    "Valid correct max condition type",
			input:   CORRECT_MAX,
			wantErr: false,
		},
		{
			name:    "Valid func condition type",
			input:   FUNC,
			wantErr: false,
		},
		{
			name:    "Valid relation condition type",
			input:   RELATION,
			wantErr: false,
		},
		{
			name:    "Valid unique condition type",
			input:   UNIQUE,
			wantErr: false,
		},
		{
			name:    "Invalid condition type",
			input:   ConditionType("invalid"),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := LookupConditionType(test.input)
			if test.wantErr {
				assert.Error(t, err, "Expected an error for invalid condition type")
			} else {
				assert.NoError(t, err, "Expected no error for valid condition type")
			}
		})
	}
}

func TestGetConditionType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ConditionType
		wantErr bool
	}{
		{
			name:    "Valid condition type",
			input:   "equValue",
			want:    EQUAL,
			wantErr: false,
		},
		{
			name:    "Valid condition type without value",
			input:   "equ",
			want:    EQUAL,
			wantErr: false,
		},
		{
			name:    "Valid correct min condition type",
			input:   "cmn2",
			want:    CORRECT_MIN,
			wantErr: false,
		},
		{
			name:    "Invalid condition type",
			input:   "invalidValue",
			want:    ConditionType("inv"),
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetConditionType(test.input)
			if test.wantErr {
				assert.Error(t, err, "Expected an error for invalid condition type")
			} else {
				assert.NoError(t, err, "Expected no error for valid condition type")
				assert.Equal(t, test.want, got, "Expected the correct condition type")
			}
		})
	}
}

func TestGetConditionByType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		condType ConditionType
		want     string
		wantErr  bool
	}{
		{
			name:     "Valid condition",
			input:    "equValue",
			condType: EQUAL,
			want:     "Value",
			wantErr:  false,
		},
		{
			name:     "Valid condition without value",
			input:    "equ",
			condType: EQUAL,
			want:     "",
			wantErr:  true,
		},
		{
			name:     "Invalid condition type",
			input:    "invValue",
			condType: ConditionType("in"),
			want:     "",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := GetConditionByType(test.input, test.condType)
			if test.wantErr {
				assert.Error(t, err, "Expected an error for invalid condition type")
			} else {
				assert.NoError(t, err, "Expected no error for valid condition type")
				assert.Equal(t, test.want, got, "Expected the correct condition value")
			}
		})
	}
}

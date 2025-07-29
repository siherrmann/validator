package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

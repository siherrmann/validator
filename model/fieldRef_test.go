package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseFieldRef(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected *FieldRef
		wantErr  bool
	}{
		{name: "Same level", input: "name", expected: &FieldRef{Raw: "name", Path: []string{"name"}}},
		{name: "Nested", input: "ports.port", expected: &FieldRef{Raw: "ports.port", Path: []string{"ports", "port"}}},
		{name: "One level up", input: "../name", expected: &FieldRef{Raw: "../name", Up: 1, Path: []string{"name"}}},
		{name: "Two levels up", input: "../../a.b", expected: &FieldRef{Raw: "../../a.b", Up: 2, Path: []string{"a", "b"}}},
		{name: "Top level", input: "/ports.port", expected: &FieldRef{Raw: "/ports.port", Root: true, Path: []string{"ports", "port"}}},
		{name: "Missing name", input: "../", wantErr: true},
		{name: "Empty", input: "", wantErr: true},
		{name: "Empty segment", input: "a..b", wantErr: true},
		{name: "Top level and up", input: "/../a", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ref, err := ParseFieldRef(test.input)
			if test.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, test.expected, ref)
		})
	}
}

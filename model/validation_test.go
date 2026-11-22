package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsesFields(t *testing.T) {
	assert.True(t, (&Validation{}).UsesFields(), "A validation built by hand may use other fields")
	assert.False(t, (&Validation{NoFieldUsage: true}).UsesFields())
}

func TestInner(t *testing.T) {
	inner := []Validation{{Key: "a"}}
	assert.Equal(t, inner, (&Validation{InnerValidation: inner}).Inner())
	assert.Equal(t, inner, (&Validation{InnerValidation: []Validation{{Key: "b"}}, InnerValidationRef: &inner}).Inner(), "Expected the reference to be used")
	assert.Nil(t, (&Validation{}).Inner())
}

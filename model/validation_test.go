package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUsesFields(t *testing.T) {
	assert.True(t, (&Validation{}).UsesFields(), "A validation built by hand may use other fields")
	assert.False(t, (&Validation{NoFieldUsage: true}).UsesFields())
}

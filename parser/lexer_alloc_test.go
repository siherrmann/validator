package parser

import (
	"testing"

	"github.com/siherrmann/validator/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func drainLexer(l *Lexer) bool {
	for {
		tok := l.NextToken()
		switch tok.Type {
		case model.LexerIllegal:
			return false
		case model.LexerEOF:
			return true
		}
	}
}

func TestLexerAllocationFree(t *testing.T) {
	const tag = "(min18 && max99) || equ0"

	// An ILLEGAL token takes another path in the lexer, so check that the tag lexes cleanly before measuring.
	require.True(t, drainLexer(NewLexer(tag)))

	allocs := testing.AllocsPerRun(1000, func() {
		drainLexer(NewLexer(tag))
	})
	assert.Zero(t, allocs)
}

package parser

import (
	"testing"

	"github.com/siherrmann/validator/model"
)

// drainLexer scans l to EOF and reports whether every token it produced was
// well-formed (never ILLEGAL). It is shared between the correctness check below and
// the allocation measurement so both exercise exactly the same code path.
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

// TestLexerAllocationFree pins down the allocation budget the rest of this package's
// performance work depends on: lexing must not allocate at all. Input is a string
// (not []rune) and every readX helper returns a substring of it rather than a fresh
// string, so a full pass to EOF should cost zero allocations regardless of tag
// complexity.
func TestLexerAllocationFree(t *testing.T) {
	const tag = "(min18 && max99) || equ0"

	// Assert success once, outside the measured region. An ILLEGAL token would send
	// the lexer's default branch down a different path; measuring that by accident
	// would make this test pass for the wrong reason.
	if ok := drainLexer(NewLexer(tag)); !ok {
		t.Fatalf("expected %q to lex cleanly", tag)
	}

	allocs := testing.AllocsPerRun(1000, func() {
		drainLexer(NewLexer(tag))
	})

	if allocs != 0 {
		t.Fatalf("expected lexing %q to allocate 0 times, got %v", tag, allocs)
	}
}

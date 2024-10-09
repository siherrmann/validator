package parser

import (
	"strings"

	"github.com/siherrmann/validator/model"
)

// Lexer performs lexical analysis/scanning of the JSON
type Lexer struct {
	Input         []rune
	char          rune            // current char under examination
	lastTokenType model.TokenType // last TokenType for splitting condition type from value.
	position      int             // current position in input (points to current char)
	nextPosition  int             // current reading position in input (after current char)
	line          int             // line number for better error reporting, etc
}

// NewLexer creates and returns a pointer to the Lexer
func NewLexer(input string) *Lexer {
	l := &Lexer{Input: []rune(input)}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.nextPosition >= len(l.Input) {
		// End of input (haven't read anything yet or EOF)
		// 0 is ASCII code for "NUL" character
		l.char = 0
	} else {
		l.char = l.Input[l.nextPosition]
	}

	l.position = l.nextPosition
	l.nextPosition++
}

// NextToken switches through the lexer's current char and creates a new model.
// It then it calls readChar to advance the lexer and it returns the token.
func (l *Lexer) NextToken() model.Token {
	var t model.Token

	l.skipWhitespace()

	switch l.char {
	case '-':
		t = newToken(model.LexerEmptyRequirement, l.line, l.position, l.position+1, l.char)
	case '(':
		t = newToken(model.LexerLeftBrace, l.line, l.position, l.position+1, l.char)
	case ')':
		t = newToken(model.LexerRightBrace, l.line, l.position, l.position+1, l.char)
	case '|', '&':
		t.Literal = l.readOperator()

		t.Line = l.line
		t.Start = l.position
		t.End = l.position + 2

		t.Type = model.LexerOperator
		l.lastTokenType = model.LexerOperator
	case '\'':
		t.Literal = l.readString()

		t.Line = l.line
		t.Start = l.position
		t.End = l.position + 1

		t.Type = model.LexerConditionValue
		l.lastTokenType = model.LexerConditionValue
	case 0:
		t.Literal = ""

		t.Line = l.line

		t.Type = model.LexerEOF
		l.lastTokenType = model.LexerEOF
	default:
		if l.lastTokenType == model.LexerConditionType {
			t.Literal = l.readConditionValue()

			t.Line = l.line
			t.Start = l.position
			t.End = l.position

			t.Type = model.LexerConditionValue
			l.lastTokenType = model.LexerConditionValue
			return t
		} else if isLetter(l.char) {
			t.Literal = l.readConditionType()

			t.Line = l.line
			t.Start = l.position
			t.End = l.position

			t.Type = model.LexerConditionType
			l.lastTokenType = model.LexerConditionType
			return t
		}
		t = newToken(model.LexerIllegal, l.line, l.position, l.position, l.char)
	}

	l.readChar()

	return t
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' {
			l.line++
		}
		l.readChar()
	}
}

func newToken(tokenType model.TokenType, line, start, end int, char ...rune) model.Token {
	return model.Token{
		Type:    tokenType,
		Literal: string(char),
		Line:    line,
		Start:   start,
		End:     end,
	}
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z'
}

func isOperator(char rune) bool {
	return char == '|' || char == '&'
}

// readOperator sets a start position and reads through two characters to get a full operator
func (l *Lexer) readOperator() string {
	position := l.position
	for isOperator(l.char) && l.position < position+2 {
		l.readChar()
	}
	return string(l.Input[position:l.position])
}

// readString sets a start position and reads through characters
// until it finds a closing `'`. It stops consuming characters and
// returns the string between the start and end positions.
// The charakter`'` inside the string is escaped with a '/'
// and the '/' is removed after reading the string.
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		prevChar := l.char
		l.readChar()
		if (l.char == '\'' && prevChar != '/') || l.char == 0 {
			break
		}
	}
	// remove custom escaped `'`
	return strings.ReplaceAll(string(l.Input[position:l.position]), "/'", "'")
}

// readConditionType sets a start position and reads through 3 characters
// to get a condition type (unvalidated).
func (l *Lexer) readConditionType() string {
	position := l.position
	for isLetter(l.char) && l.position < position+3 {
		l.readChar()
	}
	return string(l.Input[position:l.position])
}

// readConditionValue sets a start position and reads through characters
// until any kind of whitespace to get a condition value (unvalidated).
func (l *Lexer) readConditionValue() string {
	position := l.position
	for l.char != ' ' && l.char != '\t' && l.char != '\n' && l.char != '\r' && l.char != ')' && l.char != 0 {
		l.readChar()
	}
	return string(l.Input[position:l.position])
}

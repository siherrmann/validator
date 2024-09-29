package model

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Start   int
	End     int
}

type TokenType string

const (
	// Token/character we don't know about
	LexerIllegal TokenType = "ILLEGAL"

	// End of file
	LexerEOF TokenType = "EOF"

	// Group
	LexerLeftBrace  TokenType = "GROUP_OPEN"
	LexerRightBrace TokenType = "GROUP_CLOSE"

	// Literals
	LexerConditionType  TokenType = "CONDITION_TYPE"
	LexerConditionValue TokenType = "CONDITION_VALUE"

	// Operators
	LexerOperator TokenType = "OPERATOR"
)

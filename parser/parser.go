package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/siherrmann/validator/model"
)

// [Parser] holds a [Lexer], errors, the [currentToken], and the [peekToken] (next token).
// Parser methods handle iterating through tokens and building and AST.
type Parser struct {
	lexer        *Lexer
	errors       []string
	currentToken model.Token
	peekToken    model.Token
}

// New takes a [Lexer], creates a [Parser] with that [Lexer], sets the current token and
// the peek token and returns the [Parser].
func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l}

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// [ParseValidation] parses tokens and creates an AST. It returns the [RootNode]
// which holds a [Value] and in it the rest of the tree.
func (p *Parser) ParseValidation() (model.RootNode, error) {
	var rootNode model.RootNode

	val := p.parseGroup()
	if val == nil {
		p.parseError(fmt.Sprintf(
			"error parsing validation, expected a value, got: %v:",
			p.currentToken.Literal,
		))
		return model.RootNode{}, errors.New(p.Errors())
	}
	rootNode.RootValue = val

	return rootNode, nil
}

// [nextToken] sets our current token to the [peekToken] and the [peekToken] to
// [p.lexer.NextToken()] which scans and returns the next token.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
	// TODO remove log
	fmt.Printf("current token: %v\n", p.peekToken)
}

func (p *Parser) currentTokenTypeIs(t model.TokenType) bool {
	return p.currentToken.Type == t
}

// [parseGroup] is called when an open left brace `(` token is found or a validation starts without a '('.
func (p *Parser) parseGroup() *model.Value {
	obj := &model.Value{Type: "Group"}
	objState := GrpStart

	for !p.currentTokenTypeIs(model.LexerEOF) {
		// TODO remove log
		fmt.Printf("current object state: %v\n", objState)
		switch objState {
		case GrpStart:
			if p.currentTokenTypeIs(model.LexerLeftBrace) {
				objState = GrpOpen
				// obj.Start = p.currentToken.Start
				p.nextToken()
			} else if p.currentTokenTypeIs(model.LexerConditionType) {
				objState = GrpOpen
				// obj.Start = p.currentToken.Start
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing validation group, expected `(` or condition, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case GrpOpen:
			if p.currentTokenTypeIs(model.LexerRightBrace) {
				p.nextToken()
				// obj.End = p.currentToken.Start
				return obj
			} else if p.currentTokenTypeIs(model.LexerLeftBrace) {
				group := p.parseGroup()
				obj.ConditionGroup = append(obj.ConditionGroup, group)
				objState = GrpOpen
			} else if p.currentTokenTypeIs(model.LexerConditionType) {
				condition := p.parseCondition()
				obj.ConditionGroup = append(obj.ConditionGroup, condition)
				objState = GrpOpen
			} else if p.currentTokenTypeIs(model.LexerOperator) {
				operator, err := p.parseOperator()
				if err != nil {
					p.parseError(fmt.Sprintf(
						"error parsing operator type %s with error: %v",
						p.currentToken.Literal,
						err.Error(),
					))
				}
				if len(obj.ConditionGroup) > 0 {
					obj.ConditionGroup[len(obj.ConditionGroup)-1].Operator = operator
				}
				objState = GrpOpen
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing group, expected RightBrace, CndType or operator token, got: %s, type: %v",
					p.currentToken.Literal,
					p.currentToken.Type,
				))
				return nil
			}
		}
	}

	// obj.End = p.currentToken.Start

	return obj
}

// [parseCondition] is used to parse a condition and setting the `conditionType`:`condition` pair.
func (p *Parser) parseCondition() *model.Value {
	condition := &model.Value{Type: "Condition"}
	conditionState := ConType

	for conditionState != ConEnd {
		switch conditionState {
		case ConType:
			if p.currentTokenTypeIs(model.LexerConditionType) {
				conditionState = ConValue
				var err error
				condition.ConditionType, err = p.parseConditionType()
				if err != nil {
					p.parseError(fmt.Sprintf(
						"error parsing condition type %s with error: %v",
						p.currentToken.Literal,
						err.Error(),
					))
				}
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing condition type, expected CndType token, got: %s",
					p.currentToken.Literal,
				))
			}
		case ConValue:
			if p.currentTokenTypeIs(model.LexerConditionValue) {
				conditionState = ConEnd
				condition.ConditionValue = p.parseConditionValue()
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing condition, expected ConValue token, got: %s",
					p.currentToken.Literal,
				))
			}
		}
	}
	return condition
}

// [parseOperator] is used to parse the condition type.
func (p *Parser) parseOperator() (model.Operator, error) {
	operator := model.Operator(p.currentToken.Literal)
	err := model.LookupOperator(operator)
	if err != nil {
		return operator, err
	}
	return operator, nil
}

// [parseConditionType] is used to parse the condition type.
func (p *Parser) parseConditionType() (model.ConditionType, error) {
	conType := model.ConditionType(p.currentToken.Literal)
	err := model.LookupConditionType(conType)
	if err != nil {
		return conType, err
	}
	return conType, nil
}

// [parseConditionValue] is used to parse the condition value (eg. 10 if min length of string is 10 (min10)).
func (p *Parser) parseConditionValue() string {
	return p.currentToken.Literal
}

// TODO implement peeking
// // expectPeekType checks the next token type against the one passed in. If it matches,
// // we call p.nextToken() to set us to the expected token and return true. If the expected
// // type does not match, we add a peek error and return false.
// func (p *Parser) expectPeekType(t model.TokenType) bool {
// 	if p.peekTokenTypeIs(t) {
// 		p.nextToken()
// 		return true
// 	}
// 	p.peekError(t)
// 	return false
// }

// func (p *Parser) peekTokenTypeIs(t model.TokenType) bool {
// 	return p.peekToken.Type == t
// }

// // peekError is a small wrapper to add a peek error to our parser's errors field.
// func (p *Parser) peekError(t model.TokenType) {
// 	msg := fmt.Sprintf(
// 		"Line: %d: Expected next token to be %s, got: %s instead",
// 		p.currentToken.Line,
// 		t,
// 		p.peekToken.Type,
// 	)
// 	p.errors = append(p.errors, msg)
// }

// parseError is very similar to `peekError`, except it simply takes a string message that
// gets appended to the parser's errors
func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors is simply a helper function that returns the parser's errors
func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}

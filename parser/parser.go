package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/siherrmann/validator/model"
)

// Parser holds a Lexer, errors, the currentToken, and the peekToken (next token).
// Parser methods handle iterating through tokens and building and AST.
type Parser struct {
	lexer        *Lexer
	errors       []string
	currentToken model.Token
	peekToken    model.Token
}

// New takes a Lexer, creates a Parser with that Lexer, sets the current token and
// the peek token and returns the Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{lexer: l}

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

// ParseValidation parses tokens and creates an AST. It returns the RootNode
// which holds a Value and in it the rest of the tree.
func (p *Parser) ParseValidation() (model.RootNode, error) {
	var rootNode model.RootNode

	val := p.parseGroup(true)
	if val == nil || len(p.Errors()) > 0 {
		p.parseError(fmt.Sprintf(
			"error parsing validation, expected a value, got: %v:",
			p.currentToken.Literal,
		))
		return model.RootNode{RootValue: &model.AstValue{}}, errors.New(p.Errors())
	}
	rootNode.RootValue = val

	return rootNode, nil
}

// nextToken sets our current token to the peekToken and the peekToken to
// p.lexer.NextToken() which scans and returns the next token.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) currentTokenTypeIs(t model.TokenType) bool {
	return p.currentToken.Type == t
}

// parseGroup is called when an open left brace `(` token is found or a requirement starts without a '('.
func (p *Parser) parseGroup(root bool) *model.AstValue {
	group := &model.AstValue{Type: model.GROUP}
	grpState := model.GrpStart

	for !p.currentTokenTypeIs(model.LexerEOF) && grpState != model.GrpEnd {
		switch grpState {
		case model.GrpStart:
			if p.currentTokenTypeIs(model.LexerLeftBrace) {
				if root {
					innerGroup := p.parseGroup(false)
					group.ConditionGroup = append(group.ConditionGroup, innerGroup)
					grpState = model.GrpOpen
				} else {
					group.Start = p.currentToken.Start
					p.nextToken()
					grpState = model.GrpOpen
				}
			} else if p.currentTokenTypeIs(model.LexerConditionType) {
				group.Start = p.currentToken.Start
				grpState = model.GrpOpen
			} else if p.currentTokenTypeIs(model.LexerEmptyRequirement) {
				group.Type = model.EMPTY
				group.Start = p.currentToken.Start
				group.End = p.currentToken.End
				grpState = model.GrpEnd
				return group
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing validation group, expected `(`, `-` or condition, got: %s",
					p.currentToken.Literal,
				))
				return nil
			}
		case model.GrpOpen:
			if p.currentTokenTypeIs(model.LexerRightBrace) {
				group.End = p.currentToken.End
				p.nextToken()
				grpState = model.GrpEnd
			} else if p.currentTokenTypeIs(model.LexerLeftBrace) {
				innerGroup := p.parseGroup(false)
				group.ConditionGroup = append(group.ConditionGroup, innerGroup)
				if len(group.ConditionGroup) > 1 && len(group.ConditionGroup[len(group.ConditionGroup)-2].Operator) == 0 {
					group.ConditionGroup[len(group.ConditionGroup)-2].Operator = model.AND
				}
			} else if p.currentTokenTypeIs(model.LexerConditionType) {
				condition := p.parseCondition()
				group.ConditionGroup = append(group.ConditionGroup, condition)
				if len(group.ConditionGroup) > 1 && len(group.ConditionGroup[len(group.ConditionGroup)-2].Operator) == 0 {
					group.ConditionGroup[len(group.ConditionGroup)-2].Operator = model.AND
				}
			} else if p.currentTokenTypeIs(model.LexerOperator) {
				operator := p.parseOperator()
				if len(group.ConditionGroup) > 0 {
					group.ConditionGroup[len(group.ConditionGroup)-1].Operator = operator
				}
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing group, expected `)`, condition or operator, got: %s, type: %v",
					p.currentToken.Literal,
					p.currentToken.Type,
				))
				return nil
			}
		}
	}

	group.End = p.currentToken.Start

	return group
}

// parseCondition is used to parse a condition and setting the `conditionType`:`condition` pair.
func (p *Parser) parseCondition() *model.AstValue {
	condition := &model.AstValue{Type: model.CONDITION}
	conditionState := model.ConType

	for conditionState != model.ConEnd {
		switch conditionState {
		case model.ConType:
			if p.currentTokenTypeIs(model.LexerConditionType) {
				condition.ConditionType = p.parseConditionType()
				p.nextToken()
				conditionState = model.ConValue
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing condition type, expected CndType token, got: %s",
					p.currentToken.Literal,
				))
				return condition
			}
		case model.ConValue:
			if p.currentTokenTypeIs(model.LexerConditionValue) || p.currentTokenTypeIs(model.LexerConditionValueString) {
				condition.ConditionValue = p.parseConditionValue()
				p.nextToken()
				conditionState = model.ConEnd
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing condition, expected ConValue token, got: %s",
					p.currentToken.Literal,
				))
				return condition
			}
		}
	}

	return condition
}

// parseOperator is used to parse the condition type.
func (p *Parser) parseOperator() model.Operator {
	operator := model.Operator(p.currentToken.Literal)
	err := model.LookupOperator(operator)
	if err != nil {
		p.parseError(fmt.Sprintf(
			"error parsing operator type %s with error: %v",
			p.currentToken.Literal,
			err.Error(),
		))
	}
	return operator
}

// parseConditionType is used to parse the condition type.
func (p *Parser) parseConditionType() model.ConditionType {
	conType := model.ConditionType(p.currentToken.Literal)
	err := model.LookupConditionType(conType)
	if err != nil {
		p.parseError(fmt.Sprintf(
			"error parsing condition type %s with error: %v",
			p.currentToken.Literal,
			err.Error(),
		))
	}
	return conType
}

// parseConditionValue is used to parse the condition value (eg. 10 if min length of string is 10 (min10)).
func (p *Parser) parseConditionValue() string {
	return p.currentToken.Literal
}

// parseError is very similar to `peekError`, except it simply takes a string message that
// gets appended to the parser's errors
func (p *Parser) parseError(msg string) {
	p.errors = append(p.errors, msg)
}

// Errors is simply a helper function that returns the parser's errors
func (p *Parser) Errors() string {
	return strings.Join(p.errors, ", ")
}

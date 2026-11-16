package parser

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
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

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseValidation parses tokens and creates an AST. It returns the RootNode
// which holds a Value and in it the rest of the tree.
// It creates a new Lexer for reset lexer state and initializes the parser's tokens.
// It returns an error if the parsing fails.
func (p *Parser) ParseValidation(validation string) (model.RootNode, error) {
	// nil rather than []string{}: append handles a nil slice fine, and a successful
	// parse (the common case) never appends, so it never pays for this allocation.
	p.errors = nil
	p.lexer = NewLexer(validation)

	// Read two tokens, so currentToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	var rootNode model.RootNode

	val := p.parseGroup(true)
	if val == nil || len(p.Errors()) > 0 {
		p.parseError(fmt.Sprintf(
			"error parsing validation, expected a value, got: %v:",
			p.currentToken.Literal,
		))
		return model.RootNode{RootValue: &model.AstValue{}}, errors.New(p.Errors())
	}
	val.UsesFields = val.HasFieldUsage()
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
	return p.parseGroupFrom(root, false)
}

// parseGroupFrom parses a group. With open, the group starts as already opened.
// This is used for the inner conditions of `rel(<path> <conditions>)`,
// where the opening brace and the path are already consumed.
func (p *Parser) parseGroupFrom(root bool, open bool) *model.AstValue {
	// Most tags hold only a handful of conditions/sub-groups; sizing the slice up
	// front avoids the repeated grow-and-copy an append-from-nil would do. A group
	// that turns out to be EMPTY or errors before appending anything just leaves this
	// capacity unused, which is cheaper than the reallocations it saves in the common case.
	group := &model.AstValue{Type: model.GROUP, ConditionGroup: make(model.ConditionGroup, 0, 4)}
	grpState := model.GrpStart
	if open {
		grpState = model.GrpOpen
		group.Start = p.currentToken.Start
	}

	for !p.currentTokenTypeIs(model.LexerEOF) && grpState != model.GrpEnd {
		if len(p.errors) > 0 {
			return nil
		}

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
					"error parsing validation group, expected left brace, `-` or condition, got: %s",
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
			} else if p.currentTokenTypeIs(model.LexerOperator) && len(group.ConditionGroup) > 0 {
				operator := p.parseOperator()
				group.ConditionGroup[len(group.ConditionGroup)-1].Operator = operator
				p.nextToken()
			} else {
				p.parseError(fmt.Sprintf(
					"error parsing group, expected right brace, condition or operator, got: %s, type: %v",
					p.currentToken.Literal,
					p.currentToken.Type,
				))
				return nil
			}
		}
	}

	if p.currentTokenTypeIs(model.LexerEOF) && grpState == model.GrpOpen && !root {
		p.parseError(fmt.Sprintf(
			"error parsing group, expected right brace, got end of line after: %s",
			p.lexer.lastTokenType,
		))
		return nil
	}

	group.End = p.currentToken.Start

	if len(p.errors) > 0 {
		return nil
	}
	if err := normalizeImplication(group); err != nil {
		p.parseError(err.Error())
		return nil
	}

	return group
}

// normalizeImplication restructures a group with an implication (`=>`) into exactly two values:
// the premise (everything before the first `=>`, carrying the `=>` operator) and the consequence
// (everything after it). `=>` binds weaker than `&&` and `||` and is right associative,
// so `a && b => c => d` is `(a && b) => (c => d)`. Doing this while parsing keeps the validation
// a single check: if the premise fails, the group is fulfilled.
func normalizeImplication(group *model.AstValue) error {
	conditions := group.ConditionGroup
	k := -1
	for i, v := range conditions {
		if v.Operator == model.IMPLIES {
			k = i
			break
		}
	}
	if k == -1 {
		return nil
	}
	if last := conditions[len(conditions)-1]; len(last.Operator) > 0 {
		return fmt.Errorf("error parsing group, expected a condition after %v", last.Operator)
	}

	premise := conditions[k]
	if k > 0 {
		inner := conditions[: k+1 : k+1]
		inner[k].Operator = ""
		premise = &model.AstValue{Type: model.GROUP, ConditionGroup: inner, Operator: model.IMPLIES, Start: inner[0].Start, End: inner[k].End}
	}

	consequence := conditions[k+1]
	if len(conditions)-(k+1) > 1 {
		inner := conditions[k+1:]
		consequence = &model.AstValue{Type: model.GROUP, ConditionGroup: inner, Start: inner[0].Start, End: inner[len(inner)-1].End}
		if err := normalizeImplication(consequence); err != nil {
			return err
		}
	}

	group.ConditionGroup = model.ConditionGroup{premise, consequence}
	return nil
}

// parseCondition is used to parse a condition and setting the `conditionType`:`condition` pair.
//
// Special cases:
//   - `rel(<path> <conditions>)` parses the path of the related field and its conditions as a group.
//   - A value starting with `$` (not quoted) is a reference to another field (eg. `min$min_replicas`)
//     or to an injected value (eg. `equ${api_key}`).
//   - `unq` can be used without a value.
//   - Fixed values are prepared once here (see precomputeConditionValue).
func (p *Parser) parseCondition() *model.AstValue {
	condition := &model.AstValue{Type: model.CONDITION}
	conditionState := model.ConType

	for conditionState != model.ConEnd {
		switch conditionState {
		case model.ConType:
			if p.currentTokenTypeIs(model.LexerConditionType) {
				condition.Start = p.currentToken.Start
				condition.ConditionType = p.parseConditionType()
				p.nextToken()
				if condition.ConditionType == model.RELATION {
					p.parseRelation(condition)
					return condition
				}
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
				if p.currentTokenTypeIs(model.LexerConditionValue) && strings.HasPrefix(condition.ConditionValue, "$") {
					p.parseReference(condition)
				} else {
					p.precomputeConditionValue(condition)
				}
				condition.End = p.currentToken.End
				p.nextToken()
				conditionState = model.ConEnd
			} else if condition.ConditionType == model.UNIQUE {
				// `unq` without a value checks the items themselves.
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

// parseRelation parses `(<path> <conditions>)` after the condition type `rel`.
func (p *Parser) parseRelation(condition *model.AstValue) {
	if !p.currentTokenTypeIs(model.LexerLeftBrace) {
		p.parseError(fmt.Sprintf("error parsing rel, expected `(`, got: %s", p.currentToken.Literal))
		return
	}
	p.nextToken()
	if !p.currentTokenTypeIs(model.LexerFieldPath) {
		p.parseError(fmt.Sprintf("error parsing rel, expected a field path, got: %s", p.currentToken.Literal))
		return
	}
	ref, err := model.ParseFieldRef(p.currentToken.Literal)
	if err != nil {
		p.parseError(fmt.Sprintf("error parsing rel: %v", err))
		return
	}
	condition.Ref = ref
	condition.ConditionValue = ref.Raw
	p.nextToken()
	if p.currentTokenTypeIs(model.LexerRightBrace) {
		p.parseError(fmt.Sprintf("error parsing rel, expected conditions for field %s", ref.Raw))
		return
	}

	inner := p.parseGroupFrom(false, true)
	if inner == nil {
		return
	}
	condition.ConditionGroup = inner.ConditionGroup
	condition.End = inner.End
}

// parseReference parses a `$<path>` or `${<name>}` condition value.
func (p *Parser) parseReference(condition *model.AstValue) {
	switch condition.ConditionType {
	case model.FUNC, model.DEF, model.IS_TYPE, model.CORRECT_MIN, model.CORRECT_MAX, model.UNIQUE:
		p.parseError(fmt.Sprintf("error parsing condition, %v can not reference a field", condition.ConditionType))
		return
	}
	// `$name` is a field (path without the `$`), `${name}` an injected value.
	raw := condition.ConditionValue
	if !strings.HasPrefix(raw, "${") {
		raw = raw[1:]
	}
	ref, err := model.ParseFieldRef(raw)
	if err != nil {
		p.parseError(fmt.Sprintf("error parsing reference: %v", err))
		return
	}
	condition.Ref = ref
}

// precomputeConditionValue prepares the fixed value of a condition once while parsing,
// so it does not have to be done on every validation:
// `rex` is compiled, `frm` and `nfr` are split and `min` and `max` are parsed to a number.
func (p *Parser) precomputeConditionValue(condition *model.AstValue) {
	switch condition.ConditionType {
	case model.REGX:
		regex, err := regexp.Compile(condition.ConditionValue)
		if err != nil {
			p.parseError(fmt.Sprintf("error parsing regex %s: %v", condition.ConditionValue, err))
			return
		}
		condition.Regex = regex
	case model.FROM, model.NOT_FROM:
		if strings.TrimSpace(condition.ConditionValue) != "" {
			condition.List = strings.Split(condition.ConditionValue, ",")
		}
	case model.MIN_VALUE, model.MAX_VALUE:
		if number, err := strconv.ParseFloat(condition.ConditionValue, 64); err == nil {
			condition.Number = number
			condition.HasNumber = true
		}
	}
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

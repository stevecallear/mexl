package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/stevecallear/mexl/ast"
	"github.com/stevecallear/mexl/ast/token"
	"github.com/stevecallear/mexl/parser/lexer"
)

type (
	Lexer interface {
		NextToken() token.Token
	}

	Parser struct {
		l            Lexer
		currentToken token.Token
		peekToken    token.Token
		errors       []string
	}

	prefixParseFn func() ast.Node
	infixParseFn  func(ast.Node) ast.Node
)

const (
	precedenceLowest int = iota + 1
	precedenceOr
	precedenceAnd
	precedenceEquals
	precedenceLessGreater
	precedenceSum
	precedenceProduct
	precedencePrefix
	precedenceStartsEndsWith
	precedenceIn
	precedenceCall
	precedenceIndex
	precedenceMember
)

var precedences = map[token.Type]int{
	token.Or:                 precedenceOr,
	token.And:                precedenceAnd,
	token.Equal:              precedenceEquals,
	token.NotEqual:           precedenceEquals,
	token.LessThan:           precedenceLessGreater,
	token.GreaterThan:        precedenceLessGreater,
	token.LessThanOrEqual:    precedenceLessGreater,
	token.GreaterThanOrEqual: precedenceLessGreater,
	token.StartsWith:         precedenceStartsEndsWith,
	token.EndsWith:           precedenceStartsEndsWith,
	token.In:                 precedenceIn,
	token.Plus:               precedenceSum,
	token.Minus:              precedenceSum,
	token.Asterisk:           precedenceProduct,
	token.Slash:              precedenceProduct,
	token.Percent:            precedenceProduct,
	token.LParen:             precedenceCall,
	token.LBracket:           precedenceIndex,
	token.Stop:               precedenceMember,
}

func New(input string) *Parser {
	return NewWithLexer(lexer.New(input))
}

func NewWithLexer(l Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken() // set peek token
	p.nextToken() // set current token
	return p
}

func (p *Parser) Parse() (ast.Node, error) {
	n := p.parseExpression(precedenceLowest)

	if p.peekToken.Type != token.EOF {
		p.error(fmt.Sprintf("multiple expressions found, next: %d", p.peekToken.Type))
	}

	if len(p.errors) > 0 {
		var b strings.Builder
		for _, err := range p.errors {
			b.WriteString(err + "\n")
		}
		return nil, errors.New(b.String())
	}

	return n, nil
}

func (p *Parser) parseExpression(precedence int) ast.Node {
	fn, ok := p.getPrefixParseFn(p.currentToken.Type)
	if !ok {
		p.error(fmt.Sprintf("no prefix parse function: %s", p.currentToken.Literal))
		return nil
	}

	e := fn()

	for !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		fn, ok := p.getInfixParseFn(p.peekToken.Type)
		if !ok {
			return e
		}

		p.nextToken()
		e = fn(e)
	}

	return e
}

func (p *Parser) getPrefixParseFn(t token.Type) (prefixParseFn, bool) {
	switch t {
	case token.Ident:
		return p.parseIdentifier, true
	case token.Int:
		return p.parseIntegerLiteral, true
	case token.Float:
		return p.parseFloatLiteral, true
	case token.String:
		return p.parseStringLiteral, true
	case token.Bang, token.Minus:
		return p.parsePrefixExpression, true
	case token.True, token.False:
		return p.parseBoolean, true
	case token.Null:
		return p.parseNull, true
	case token.LParen:
		return p.parseGroupedExpression, true
	case token.LBracket:
		return p.parseArrayLiteral, true
	default:
		return nil, false
	}
}

func (p *Parser) getInfixParseFn(t token.Type) (infixParseFn, bool) {
	switch t {
	case token.Plus, token.Minus, token.Asterisk, token.Slash, token.Percent:
		return p.parseInfixExpression, true
	case token.Equal, token.NotEqual, token.LessThan, token.LessThanOrEqual, token.GreaterThan, token.GreaterThanOrEqual:
		return p.parseInfixExpression, true
	case token.And, token.Or:
		return p.parseInfixExpression, true
	case token.StartsWith, token.EndsWith, token.In:
		return p.parseInfixExpression, true
	case token.LParen:
		return p.parseCallExpression, true
	case token.LBracket:
		return p.parseIndexExpression, true
	case token.Stop:
		return p.parseMemberExpression, true
	default:
		return nil, false
	}
}

func (p *Parser) parsePrefixExpression() ast.Node {
	e := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	e.Right = p.parseExpression(precedencePrefix)
	return e
}

func (p *Parser) parseInfixExpression(left ast.Node) ast.Node {
	e := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	pr := p.currentPrecedence()
	p.nextToken()

	e.Right = p.parseExpression(pr)
	return e
}

func (p *Parser) parseMemberExpression(left ast.Node) ast.Node {
	e := &ast.MemberExpression{Token: p.currentToken, Left: left}
	pr := p.currentPrecedence()
	p.nextToken()

	e.Member = p.parseExpression(pr)
	return e
}

func (p *Parser) parseGroupedExpression() ast.Node {
	p.nextToken()

	e := p.parseExpression(precedenceLowest)

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return e
}

func (p *Parser) parseCallExpression(fn ast.Node) ast.Node {
	return &ast.CallExpression{
		Token:     p.currentToken,
		Function:  fn,
		Arguments: p.parseExpressionList(token.RParen),
	}
}

func (p *Parser) parseExpressionList(t token.Type) []ast.Node {
	l := []ast.Node{}

	if p.peekTokenIs(t) {
		p.nextToken()
		return l
	}

	p.nextToken()
	l = append(l, p.parseExpression(precedenceLowest))

	for p.peekTokenIs(token.Comma) {
		p.nextToken()
		p.nextToken()
		l = append(l, p.parseExpression(precedenceLowest))
	}

	if !p.expectPeek(t) {
		return nil
	}

	return l
}

func (p *Parser) parseIndexExpression(left ast.Node) ast.Node {
	e := &ast.IndexExpression{Token: p.currentToken, Left: left}
	p.nextToken()

	e.Index = p.parseExpression(precedenceLowest)
	if !p.expectPeek(token.RBracket) {
		return nil
	}

	return e
}

func (p *Parser) parseIdentifier() ast.Node {
	return &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseIntegerLiteral() ast.Node {
	l := &ast.IntegerLiteral{Token: p.currentToken}

	var err error
	l.Value, err = strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		p.error(fmt.Sprintf("invalid integer literal: %s", p.currentToken.Literal))
		return nil
	}

	return l
}

func (p *Parser) parseFloatLiteral() ast.Node {
	l := &ast.FloatLiteral{Token: p.currentToken}

	var err error
	l.Value, err = strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.error(fmt.Sprintf("invalid float literal: %s", p.currentToken.Literal))
		return nil
	}

	return l
}

func (p *Parser) parseStringLiteral() ast.Node {
	return &ast.StringLiteral{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	}
}

func (p *Parser) parseBoolean() ast.Node {
	return &ast.Boolean{
		Token: p.currentToken,
		Value: p.curTokenIs(token.True),
	}
}

func (p *Parser) parseNull() ast.Node {
	return &ast.Null{
		Token: p.currentToken,
	}
}

func (p *Parser) parseArrayLiteral() ast.Node {
	return &ast.ArrayLiteral{
		Token:    p.currentToken,
		Elements: p.parseExpressionList(token.RBracket),
	}
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return precedenceLowest
}

func (p *Parser) curTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.error(fmt.Sprintf("invalid next token: %d, expected %d", p.peekToken.Type, t))
	return false
}

func (p *Parser) error(err string) {
	p.errors = append(p.errors, err)
}

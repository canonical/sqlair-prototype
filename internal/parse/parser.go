package parse

import (
	"fmt"
	"strings"
)

const (
	LOWEST = iota
	HIGHEST
)

var precedence = map[TokenType]int{
	// Token types not listed here will return LOWEST from
	// Preference() function
	RBRACKET: HIGHEST,
	PERIOD:   HIGHEST,
	RPAREN:   HIGHEST,
}

type prefixFunc func() Expression
type infixFunc func(Expression) Expression

// Parser is responsible for returning an Expression tree
// for a Sqlair DSL statement represented by a Lexer.
type Parser struct {
	lex *lexer

	// Accumulated error messages.
	errors []string

	currentToken Token
	peekToken    Token

	// Map of Prefix and Infix functions.
	// Tokens will implement at most one of them.
	prefixFn map[TokenType]prefixFunc
	infixFn  map[TokenType]infixFunc
}

// NewParser returns a reference to a Parser based on the input Lexer.
// This parser implements a top-down parsing strategy with operator
// precedence (Pratt's parser)
// https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing
func NewParser(l *lexer) *Parser {
	p := &Parser{
		lex: l,
	}

	// Fill prefix functions according to token type.
	// Token types that has prefixFunc do not care about the
	// left part of the expression.
	p.prefixFn = map[TokenType]prefixFunc{
		ASTERISK:  p.parseIdent,
		BITAND:    p.parseOutputTarget,
		DOLLAR:    p.parseInputSource,
		EQUAL:     p.parseIdent,
		IDENT:     p.parseIdent,
		NUM:       p.parseInteger,
		SEMICOLON: p.parseIdent,
		STRING:    p.parseString,
		LPAREN:    p.parseGroup,
		RPAREN:    p.parseIdent,
		COMMA:     p.parseIdent,
	}

	// Fill infix functions according to token type.
	// Tokens that has infixFunc care about the left part
	// of the expression.
	p.infixFn = map[TokenType]infixFunc{
		LBRACKET: p.parseIndex,
	}

	// Feed currentToken and peekToken
	p.nextToken()
	p.nextToken()

	return p
}

// Run returns an Expression tree using its Lexer,
// or an error for a malformed statement.
func (p *Parser) Run() (*SQLExpression, error) {
	var exp SQLExpression
	if p.currentToken.Type == EOF {
		p.errors = append(p.errors, "Empty statement")
	}
	for p.currentToken.Type != EOF {
		exp.AppendExpression(p.parseExpression(LOWEST))
		p.nextToken()
	}
	var err error
	if len(p.errors) > 0 {
		err = fmt.Errorf(strings.Join(p.errors, "\n"))
		return nil, err
	}
	return &exp, nil
}

func (p *Parser) parseExpression(prec_level int) Expression {
	prefixfunc := p.prefixFn[p.currentToken.Type]
	if prefixfunc == nil {
		// Only some tokens have a prefix function
		return nil
	}

	// Get the left part of the tree
	left := prefixfunc()

	for prec_level < p.precedence(p.peekToken) {
		infixfunc := p.infixFn[p.peekToken.Type]
		if infixfunc == nil {
			// We just need to return the left branch
			// of the expression
			return left
		}
		// Run the infix function
		p.nextToken()
		left = infixfunc(left)
	}

	return left
}

func (p *Parser) parseOutputTarget() Expression {
	marker := p.currentToken
	p.nextToken()
	name := p.parseExpression(p.precedence(p.currentToken)).(*IdentityExpression)
	// Skip period token and move to the next one
	p.nextToken()
	if p.currentToken.Type != PERIOD {
		p.errors = append(p.errors, fmt.Sprintf(
			"Line: %d, column: %d: unexpected '%s', expecting PERIOD",
			p.currentToken.Pos.Line,
			p.currentToken.Pos.Column,
			p.currentToken.Literal))
		return nil
	}
	p.nextToken()
	field := p.parseExpression(p.precedence(p.currentToken)).(*IdentityExpression)
	return NewOutputTargetExpression(marker, name, field)
}

func (p *Parser) parseInputSource() Expression {
	marker := p.currentToken
	p.nextToken()
	name := p.parseExpression(p.precedence(p.currentToken)).(*IdentityExpression)
	// Skip period token and move to the next one
	p.nextToken()
	p.nextToken()
	field := p.parseExpression(p.precedence(p.currentToken)).(*IdentityExpression)
	return NewInputSourceExpression(marker, name, field)
}

func (p *Parser) parseIdent() Expression {
	t := p.currentToken
	return &IdentityExpression{t}
}

func (p *Parser) parseString() Expression {
	return &IdentityExpression{p.currentToken}
}

func (p *Parser) parseInteger() Expression {
	return &IdentityExpression{p.currentToken}
}

func (p *Parser) parseGroup() Expression {
	// Skip left parenthesis
	p.nextToken()
	if p.currentToken.Type == RPAREN {
		// Empty group.
		p.errors = append(p.errors, fmt.Sprintf(
			"Line: %d, column: %d: unexpected ')', expecting LITERAL",
			p.currentToken.Pos.Line,
			p.currentToken.Pos.Column))
		return nil
	}
	var g GroupedColumnsExpression
	consecutiveCommas := 1
	for p.currentToken.Type != RPAREN && p.currentToken.Type != EOF {
		if p.currentToken.Type != COMMA {
			g.AppendExpression(p.parseExpression(LOWEST))
			consecutiveCommas = 0
		} else {
			consecutiveCommas++
		}
		if consecutiveCommas > 1 {
			p.errors = append(p.errors, fmt.Sprintf(
				"Line: %d, column: %d: unexpected ',', expecting LITERAL",
				p.currentToken.Pos.Line,
				p.currentToken.Pos.Column))
			return nil
		}
		p.nextToken()
	}

	if p.currentToken.Type == EOF {
		p.errors = append(p.errors, fmt.Sprintf(
			"Line: %d, column: %d: unexpected 'EOF', expecting ')'",
			p.currentToken.Pos.Line,
			p.currentToken.Pos.Column))
		return nil
	}

	return &g
}

func (p *Parser) parseIndex(left Expression) Expression {
	var pte PassThroughExpression
	pte.AppendExpression(left)
	pte.AppendExpression(&IdentityExpression{p.currentToken})
	return &pte
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) precedence(t Token) int {
	prec, ok := precedence[t.Type]
	if !ok {
		prec = LOWEST
	}
	return prec
}

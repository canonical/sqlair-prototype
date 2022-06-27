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

type PrefixFunc func() Expression
type InfixFunc func(Expression) Expression

// Parser is responsible for returning an Expression tree
// for a Sqlair DSL statement represented by a Lexer.
type Parser struct {
	lex *Lexer

	// Accumulated error messages.
	errors []string

	currentToken Token
	peekToken    Token

	// Map of Prefix and Infix functions.
	// Tokens will implement at most one of them.
	prefixfn map[TokenType]PrefixFunc
	infixfn  map[TokenType]InfixFunc
}

// NewParser returns a reference to a Parser based on the input Lexer.
// This parser implements a top-down parsing strategy with operator
// precedence (Pratt's parser)
// https://en.wikipedia.org/wiki/Operator-precedence_parser#Pratt_parsing
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		lex: l,
	}

	// Fill prefix functions according to token type.
	// Token types that has PrefixFunc do not care about the
	// left part of the expression.
	p.prefixfn = map[TokenType]PrefixFunc{
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
	// Tokens that has InfixFunc care about the left part
	// of the expression.
	p.infixfn = map[TokenType]InfixFunc{
		LBRACKET: p.parseIndex,
	}

	// Feed currentToken and peekToken
	p.NextToken()
	p.NextToken()

	return p
}

// Run returns an Expression tree using its Lexer,
// or an error for a malformed statement.
func (p *Parser) Run() (*SQLExpression, error) {
	var exp SQLExpression
	if p.currentToken.Type == EOF {
		fmt.Println("Empty statement")
	}
	for p.currentToken.Type != EOF {
		exp.Children = append(exp.Children, p.parseExpression(LOWEST))
		p.NextToken()
	}
	var err error
	if len(p.errors) > 0 {
		err = fmt.Errorf(strings.Join(p.errors, "\n"))
		return nil, err
	}
	return &exp, nil
}

func (p *Parser) parseExpression(prec_level int) Expression {
	prefixfunc := p.prefixfn[p.currentToken.Type]
	if prefixfunc == nil {
		// Only some tokens have a prefix function
		return nil
	}

	// Get the left part of the tree
	left := prefixfunc()

	for prec_level < p.Precedence(p.peekToken) {
		infixfunc := p.infixfn[p.peekToken.Type]
		if infixfunc == nil {
			// We just need to return the left branch
			// of the expression
			return left
		}
		// Run the infix function
		p.NextToken()
		left = infixfunc(left)
	}

	return left
}

func (p *Parser) parseOutputTarget() Expression {
	var ote OutputTargetExpression
	ote.Marker = p.currentToken
	p.NextToken()
	ote.Name = p.parseExpression(p.Precedence(p.currentToken))
	// Skip period token and move to the next one
	p.NextToken()
	p.NextToken()
	ote.Field = p.parseExpression(p.Precedence(p.currentToken))
	return &ote
}

func (p *Parser) parseInputSource() Expression {
	var ise InputSourceExpression
	ise.Marker = p.currentToken
	p.NextToken()
	ise.Name = p.parseExpression(p.Precedence(p.currentToken))
	// Skip period token and move to the next one
	p.NextToken()
	p.NextToken()
	ise.Field = p.parseExpression(p.Precedence(p.currentToken))
	return &ise
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
	p.NextToken()
	var g GroupedColumnsExpression
	for p.currentToken.Type != RPAREN {
		if p.currentToken.Type != COMMA {
			g.Children = append(g.Children, p.parseExpression(LOWEST))
		}
		p.NextToken()
	}
	return &g
}

func (p *Parser) parseIndex(left Expression) Expression {
	var pte PassThroughExpression
	pte.Children = append(pte.Children, left)
	pte.Children = append(pte.Children, &IdentityExpression{p.currentToken})
	return &pte
}

func (p *Parser) NextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) Precedence(t Token) int {
	prec, ok := precedence[t.Type]
	if !ok {
		prec = LOWEST
	}
	return prec
}

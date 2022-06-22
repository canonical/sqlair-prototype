package parse

import (
        "io"
        "log"
)

// Our parser has a lexer inside.
// The main function of the parser is Parse()
// that will consume tokens from the lexer and
// return a result
type Parser struct {
        l *Lexer
}

// Parser constructor function
func NewParser(s string) *Parser {
        return &Parser{l: NewLexer(s)}
}

// Returns the next expression. This can be:
//
// * A pass through expression
// * An input expression
// * An output expression
func (p *Parser) ParseNextExpression() (QueryExpression, error) {
        var token Token
	var pte PassthroughExpression

        for token = p.l.NextToken(); token.Type != BITAND && token.Type != DOLLAR && token.Type != EOF; token = p.l.NextToken() {
		// Accumulate as a pass through
		pte.Expressions = append(pte.Expressions, token)
        }

        if token.Type == EOF {
                // No more input
                return QueryExpression{}, io.EOF
        }

	// Store either & or $.
	delimiter := token

        start_pos := p.l.offset
        token = p.l.NextToken()
        if token.Type == IDENT {
                elemType := token
                token = p.l.NextToken()
		period = token
                if token.Type == PERIOD {
                        token = p.l.NextToken()
			cols := token
                        if token.Type == ASTERISK {
                                field = "all"
                        } else if token.Type == IDENT {
                                field = token.Literal
                        } else {
                                return QueryArgument{}, nil
                        }
                        return QueryArgument{
                                name:  elemType,
                                field: field,
                                from:  start_pos,
                                to:    p.l.offset,
                        }, nil
                        log.Printf("Type: %s field: %s", elemType, field)
                } else {
			pte.Expressions = append(pte.Expressions, token)
		}
        } else {
		pte.Expressions = append(pte.Expressions, token)
	}

        return QueryArgument{}, nil
}

func (p *Parser) Parse() (*QueryExpression, error) {
        var qe QueryExpression

	// Build AST
        for exp, err := p.ParseNextExpression(); err == nil; exp, err = p.ParseNextExpression() {
		qe.Expressions = append(qe.Expressions, exp)
        }

        return qe, nil
}

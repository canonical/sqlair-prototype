package parse

import (
	"fmt"
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

type EOFError struct {
}

func (e *EOFError) Error() string {
	return fmt.Sprintf("EOF")
}

// Returns the next expression. This can be:
//
// * A pass through expression
// * An input expression
// * An output expression
func (p *Parser) ParseNextExpression() (Expression, error) {
        var token Token
	var pte PassthroughExpression

	for token = p.l.NextToken(); token.Type != BITAND && token.Type != DOLLAR && token.Type != EOF; token = p.l.NextToken() {
		// Accumulate as a pass through
		pte.Elements = append(pte.Elements, Element {Token: token})
	}

        if token.Type == EOF {
                // No more input
                return &pte, &EOFError{}
        }

	// Store either & or $.
	delimiter := token

	token = p.l.NextToken()
	if token.Type == IDENT {
		elemType := token
		token = p.l.NextToken()
		period := token
		if token.Type == PERIOD {
			token = p.l.NextToken()
			cols := token
			if delimiter.Type == BITAND {
				return &OutputExpression {
					Amp:	Element {Token: delimiter},
					Name:	Element {Token: elemType},
					Period:	Element {Token: period},
					Cols:	Element {Token: cols},
				}, nil
			} else if delimiter.Type == DOLLAR {
				return &InputExpression {
					Dollar:	Element {Token: delimiter},
					Name:	Element {Token: elemType},
					Period:	Element {Token: period},
					Cols:	Element {Token: cols},
				}, nil
			}
			log.Printf("Type: %s cols: %s", elemType.Literal, cols.Literal)
		} else {
			pte.Elements = append(pte.Elements, Element {Token: token})
		}
	} else {
		pte.Elements = append(pte.Elements, Element {Token: token})
	}

        return &pte, nil
}

func (p *Parser) Parse() (*QueryExpression, error) {
        var qe QueryExpression

	// Build AST
        for exp, err := p.ParseNextExpression(); err == nil; exp, err = p.ParseNextExpression() {
		qe.Expressions = append(qe.Expressions, exp)
        }

        return &qe, nil
}

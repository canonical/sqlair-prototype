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

type QueryArgument struct {
        name  string // Name of the element (e.g. Persona)
        field string // What field we are accessing
        from  int    // start position in the string
        to    int    // end position in the string
        //TODO: do we need a reference to the original statement?
}

// Parser constructor function
func NewParser(s string) *Parser {
        return &Parser{l: NewLexer(s)}
}

func (p *Parser) NextQueryArgument() (QueryArgument, error) {
        var token Token
        for token = p.l.NextToken(); token.Type != BITAND && token.Type != EOF; token = p.l.NextToken() {
        }

        if token.Type == EOF {
                // No more input
                return QueryArgument{}, io.EOF
        }

        start_pos := p.l.offset
        token = p.l.NextToken()
        if token.Type == IDENT {
                elemType := token.Literal
                token = p.l.NextToken()
                if token.Type == PERIOD {
                        token = p.l.NextToken()
                        var field string
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
                }
        }

        return QueryArgument{}, nil
}

func (p *Parser) Parse() ([]QueryArgument, error) {
        // Get query arguments from the statement
        // until we run out of text

        //arguments := make([]QueryArgument, 1)
        var arguments []QueryArgument

        for qa, err := p.NextQueryArgument(); err == nil; qa, err = p.NextQueryArgument() {
                if qa != (QueryArgument{}) {
                        arguments = append(arguments, qa)
                }
        }

        return arguments, nil
}

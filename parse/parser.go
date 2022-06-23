package parse

// Parser is responsible for returning an Expression tree
// for a Sqlair DSL statement represented by a Lexer.
type Parser struct {
	lex *Lexer
}

// NewParser returns a reference to a Parser based on the input Lexer.
func NewParser(l *Lexer) *Parser {
	return &Parser{
		lex: l,
	}
}

// Run returns an Expression tree using its Lexer,
// or an error for a malformed statement.
func (p *Parser) Run() (Expression, error) {
	return nil, nil
}

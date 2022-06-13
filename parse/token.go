package parse

// TokenType identifies the type of a token.
type TokenType int

const (
	UNKNOWN TokenType = iota - 1
	EOF

	IDENT
	NUM // Number literal.
	STRING

	COMMA // ,

	LPAREN   // (
	RPAREN   // )
	LBRACKET // [
	RBRACKET // ]

	BITAND    // &
	PERIOD    // .
	ASTERISK  // *
	DOLLAR    // $
	EQUAL     // =
	SEMICOLON // ;
)

var knownRuneTokens = map[rune]TokenType{
	'(': LPAREN,
	')': RPAREN,
	'[': LBRACKET,
	']': RBRACKET,
	',': COMMA,
	'&': BITAND,
	'.': PERIOD,
	'*': ASTERISK,
	'$': DOLLAR,
	'=': EQUAL,
	';': SEMICOLON,
}

// Position holds the location of the token
// within the statement containing it.
type Position struct {
	// Offset is the character offset within the containing statement.
	Offset int

	// Line indicates the line on which the statement occurs.
	// It is used for reconstructing statements with
	// some resemblance to the original formatting.
	Line int
}

// Token describes the smallest part of a larger DSL statement
// that is able to reasoned about by the parser.
type Token struct {
	// Type is the type of this token.
	Type TokenType

	// Literal is the string value of the token.
	Literal string

	// Pos is the offset of this token within a statement.
	Pos Position
}

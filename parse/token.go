package parse

// TokenType identifies the type of a token.
type TokenType int

const (
	UNKNOWN TokenType = iota - 1
	EOF
	SEPARATOR

	IDENT
	INT //int literal
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

// separatorToken returns a token for a single
// whitespace at the input position.
func separatorToken(pos Position) Token {
	return Token{
		Type:    SEPARATOR,
		Literal: " ",
		Pos:     pos,
	}
}

// maybeRuneToken returns a token for the input rune if
// it from is the subset that we specifically recognise.
// If unrecognised, zero-value, false is returned.
func maybeRuneToken(char rune, pos Position) (Token, bool) {
	runeType, isRune := runeTokens[char]
	if !isRune {
		return Token{}, false
	}

	return Token{
		Type:    runeType,
		Literal: string(char),
		Pos:     pos,
	}, true
}

var runeTokens = map[rune]TokenType{
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

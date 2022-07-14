package parse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer takes a DML or SQL statement in our DSL form and
// reads it as tokens for consumption by our parser.
// It operates in a lazy fashion, requiring a call to
// `NextToken` to advance through the input.
// Offsets are zero-indexed. Lines start at 1.
type lexer struct {
	input string
	char  rune

	offset     int
	readOffset int
	line       int
	column     int
}

// NewLexer creates a new Lexer from a given input and primes it
// with the first non-whitespace character before returning.
func NewLexer(input string) *lexer {
	l := &lexer{
		input:  strings.TrimSpace(input),
		line:   1,
		column: 1,
	}
	l.nextChar()
	return l
}

// NextToken returns the next token based on
// the current offset, ignoring whitespace.
// The EOF token is returned if we have
// reached the end of the input.
func (l *lexer) NextToken() Token {
	for l.skipWhitespace() {
	}

	pos := l.position()

	if runeType, isKnown := knownRuneTokens[l.char]; isKnown {
		lit := string(l.char)
		l.nextChar()
		return Token{
			Type:    runeType,
			Literal: lit,
			Pos:     pos,
		}

	}

	return l.readComplexToken(pos)
}

func (l *lexer) position() Position {
	return Position{
		Offset: l.offset,
		Line:   l.line,
		Column: l.column - 1,
	}
}

// skipWhitespace checks if the current character is a space
// and if so, reads the next character before returning true.
func (l *lexer) skipWhitespace() bool {
	if !unicode.IsSpace(l.char) {
		return false
	}

	l.nextChar()
	return true
}

func (l *lexer) readComplexToken(pos Position) Token {
	tok := Token{Pos: pos}

	switch {
	case l.char == 0:
		tok.Type = EOF
		return tok

	case isDigit(l.char):
		tok.Type = NUM
		tok.Literal = l.readNumber()
		return tok

	case unicode.IsLetter(l.char) || l.char == '_':
		tok.Type = IDENT
		tok.Literal = l.readIdentifier()
		return tok

	case l.char == '\'':
		tok.Type = STRING
		tok.Literal = l.readString(l.char)
		return tok
	}

	tok.Type = UNKNOWN
	tok.Literal = string(l.char)
	l.nextChar()
	return tok
}

// readIdentifier calls nextChar until it detects the end of an identifier,
// then returns the range of input from when we started reading.
func (l *lexer) readIdentifier() string {
	pos := l.offset

	for unicode.IsLetter(l.char) || isDigit(l.char) || l.char == '_' {
		l.nextChar()
	}

	return l.input[pos:l.offset]
}

// readString calls nextChar until it detects the end of a quoted string,
// then returns the range of input from when we started reading.
// The return includes the quotes.
func (l *lexer) readString(r rune) string {
	pos := l.offset

	maybeCloser := true
	for {
		// Unterminated string. Will be handled downstream.
		if l.char == 0 {
			l.nextChar()
			break
		}

		if l.char == r {
			// We're looking for string terminations.
			// Each quote is regarded an opener, or potential closer.
			maybeCloser = !maybeCloser

			// If this looks like a closing quote, check if it might be an
			// escape for a following quote. If not, we're done.
			if maybeCloser && l.peek() != r {
				l.nextChar()
				break
			}
		}

		l.nextChar()
	}

	return l.input[pos:l.offset]
}

// readNumber calls nextChar until it detects the end of a number,
// then returns the range of input from when we started reading.
func (l *lexer) readNumber() string {
	pos := l.offset

	var oneDecimal bool
	for isDigit(l.char) || l.char == '.' {
		if l.char == '.' {
			oneDecimal = true
		}

		// If we've already seen a decimal point and there is one ahead,
		// just finish the token accumulation. Syntax error goes downstream.
		if oneDecimal && l.peek() == '.' {
			l.nextChar()
			break
		}

		l.nextChar()
	}

	return l.input[pos:l.offset]
}

// nextChar reads the next character from the
// input and increments the read offset.
func (l *lexer) nextChar() {
	if l.readOffset >= len(l.input) {
		l.char = 0
		l.offset = l.readOffset
		return
	}

	var size int
	l.char, size = utf8.DecodeRuneInString(l.input[l.readOffset:])
	if l.char == '\n' {
		l.line++
		l.column = 0
	}
	l.column++
	l.offset = l.readOffset
	l.readOffset += size
}

// peek returns the next rune to be read without moving the lexer forward.
func (l *lexer) peek() rune {
	if l.readOffset >= len(l.input) {
		return 0
	}

	peek, _ := utf8.DecodeRuneInString(l.input[l.readOffset:])
	return peek
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9' || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

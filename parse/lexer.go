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
type Lexer struct {
	input string
	char  rune

	offset     int
	readOffset int
	line       int
	column     int
}

// NewLexer creates a new Lexer from a given input and primes it
// with the first non-whitespace character before returning.
func NewLexer(input string) *Lexer {
	l := &Lexer{
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
func (l *Lexer) NextToken() Token {
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

func (l *Lexer) position() Position {
	return Position{
		Offset: l.offset,
		Line:   l.line,
		Column: l.column - 1,
	}
}

// skipWhitespace checks if the current character is a space
// and if so, reads the next character before returning true.
func (l *Lexer) skipWhitespace() bool {
	if !unicode.IsSpace(l.char) {
		return false
	}

	l.nextChar()
	return true
}

func (l *Lexer) readComplexToken(pos Position) Token {
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

	case isSingleQuote(l.char):
		tok.Type = STRING
		tok.Literal = l.readString(l.char)
		return tok
	}

	tok.Type = UNKNOWN
	tok.Literal = string(l.char)
	l.nextChar()
	return tok
}

// readIdentifier calls nextChar through to the end of the identifier,
// then returns the range of input from when we started reading.
func (l *Lexer) readIdentifier() string {
	pos := l.offset

	for unicode.IsLetter(l.char) || isDigit(l.char) || l.char == '_' {
		l.nextChar()
	}

	return l.input[pos:l.offset]
}

func (l *Lexer) readString(r rune) string {
	pos := l.offset

	for i := 0; true; i++ {
		// Unterminated string. Will be handled by parser.
		if l.char == 0 {
			l.nextChar()
			break
		}

		if l.char == r {
			// If this is the first time through the loop, we just append
			// the open quote. Otherwise, we're looking for termination.
			if i != 0 {
				// Check if this appears to be a closing quote.
				// If it is, return what we've accrued so far.
				// TODO (manadart 2022-06-08): This is very naive and will not
				// correctly handle some edge cases. Review for robustness.
				next := l.peek()
				if next == ' ' || next == '\n' || next == 0 || next == ';' || next == ',' || next == ')' {
					l.nextChar()
					break
				}
			}
		}

		l.nextChar()
	}

	return l.input[pos:l.offset]
}

// readNumber returns the number beginning at current offset.
func (l *Lexer) readNumber() string {
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
func (l *Lexer) nextChar() {
	if l.readOffset >= len(l.input) {
		l.char = 0
	} else {
		l.char, _ = utf8.DecodeLastRuneInString(string(l.input[l.readOffset]))
		if l.char == '\n' {
			l.line++
			l.column = 0
		}
		l.column++
	}

	l.offset = l.readOffset
	l.readOffset++
}

// peek returns the next rune to be read without moving the lexer forward.
func (l *Lexer) peek() rune {
	if l.readOffset >= len(l.input) {
		return 0
	}

	peek, _ := utf8.DecodeLastRuneInString(string(l.input[l.readOffset]))
	return peek
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9' || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

func isSingleQuote(char rune) bool {
	return char == 39
}

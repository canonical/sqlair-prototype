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
	var ret []rune

	for i := 0; true; i++ {
		switch l.char {
		case 0:
			// Unterminated string. Will be handled by parser.
			l.nextChar()
			return string(ret)

		case r:
			ret = append(ret, l.char)

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
					return string(ret)
				}
			}
		default:
			ret = append(ret, l.char)
		}
		l.nextChar()
	}

	return string(ret)
}

// readNumber returns the number beginning at current offset.
func (l *Lexer) readNumber() string {
	var ret []rune

	ret = append(ret, l.char)
	l.nextChar()

	for isDigit(l.char) || l.char == '.' {
		if l.char == '.' {
			if l.peek() == '.' {
				return string(ret)
			}
		}

		ret = append(ret, l.char)
		l.nextChar()
	}
	return string(ret)
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

// peek attempts to read the next rune if available.
func (l *Lexer) peek() rune {
	return l.peekN(0)
}

// peekN attempts to read the next indicated by
// the input offset from the current offset.
func (l *Lexer) peekN(n int) rune {
	if l.readOffset+n >= len(l.input) {
		return 0
	}
	peek, _ := utf8.DecodeLastRuneInString(string(l.input[l.readOffset+n]))
	return peek
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9' || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

func isSingleQuote(char rune) bool {
	return char == 39
}

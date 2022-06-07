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
// Offsets and lines are zero-indexed.
type Lexer struct {
	input []rune
	char  rune

	offset     int
	readOffset int
	line       int
}

// NewLexer creates a new Lexer from a given input and primes it
// with the first non-whitespace character before returning.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: []rune(strings.TrimSpace(input))}
	l.nextChar()
	return l
}

// NextToken attempts to grab the next token available.
// Multiple whitespaces are skipped over and conflated as a single separator.
func (l *Lexer) NextToken() Token {
	pos := l.position()

	var skipped bool
	for l.skipWhitespace() {
		skipped = true
	}
	if skipped {
		return separatorToken(pos)
	}

	if token, isRune := maybeRuneToken(l.char, pos); isRune {
		l.nextChar()
		return token
	}

	return l.readToken(pos)
}

func (l *Lexer) position() Position {
	return Position{
		Offset: l.offset,
		Line:   l.line,
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

func (l *Lexer) readToken(pos Position) Token {
	tok := Token{Pos: pos}

	switch {
	case l.char == 0:
		tok.Type = EOF
		return tok

	case isDigit(l.char):
		tok.Type = INT
		tok.Literal = l.readNumber()
		return tok

	case isLetter(l.char):
		tok.Type = IDENT
		tok.Literal = l.readIdentifier()
		return tok

	case isSingleQuote(l.char):
		tok.Type = STRING
		tok.Literal = l.readString(l.char)
		return tok
	}

	l.nextChar()

	tok.Type = UNKNOWN
	tok.Literal = string(l.char)
	return tok
}

// readIdentifier calls nextChar through to the end of the identifier,
// then returns the range of input from when we started reading.
func (l *Lexer) readIdentifier() string {
	pos := l.offset

	for isLetter(l.char) || isDigit(l.char) || l.char == '-' {
		l.nextChar()
	}

	return string(l.input[pos:l.offset])
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
			if i == 0 {
				// First time through the loop. Just append the quote.
				ret = append(ret, l.char)
			} else {
				ret = append(ret, l.char)

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
		l.char = l.input[l.readOffset]
		if l.char == '\n' {
			l.line++
		}
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
	return l.input[l.readOffset+n]
}

func isLetter(char rune) bool {
	return 'a' <= char && char <= 'z' ||
		'A' <= char && char <= 'Z' ||
		char == '_' ||
		char >= utf8.RuneSelf && unicode.IsLetter(char)
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9' || char >= utf8.RuneSelf && unicode.IsDigit(char)
}

func isSingleQuote(char rune) bool {
	return char == 39
}

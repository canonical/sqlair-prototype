package parse

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeToken(t TokenType, literal string, offset, line, column int) Token {
	return Token{
		t,
		literal,
		Position{
			offset,
			line,
			column,
		},
	}
}

// Use go test -v for printing the AST
func printAST(t *testing.T, e Expression, indentation int) {
	var ind string
	if indentation == 0 {
		t.Log("Dumping AST:")
	}
	for i := 0; i < indentation; i++ {
		ind = ind + " "
	}
	t.Log(ind, e.String(), "[", reflect.TypeOf(e).String(), "]")
	for _, c := range e.Expressions() {
		printAST(t, c, indentation+4)
	}
}

// Check we handle spaces and line breaks properly
func TestParserCarriageReturnAndSpaces(t *testing.T) {
	stmt := `SELECT      a
	AS myalias    FROM
	person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	// We get an AST.
	assert.NotEqual(t, r, nil)
	// We parse properly ignoring blanks and new lines.
	assert.Equal(t, r.String(), "SELECT a AS myalias FROM person")
	// Top of the AST is a SQL Expression.
	assert.Equal(t, reflect.TypeOf(r), reflect.TypeOf(&SQLExpression{}))
	// The rest of the expressions are Indent. Nothing to be reflected
	// in this statement.
	for _, c := range r.Expressions() {
		t.Log(reflect.TypeOf(c), reflect.TypeOf(&IdentityExpression{}))
		assert.Equal(t, reflect.TypeOf(c), reflect.TypeOf(&IdentityExpression{}))
	}
}

// Check we parse column groups properly
func TestParserSimpleGroup(t *testing.T) {
	stmt := `SELECT (a, b) AS &Person.* FROM person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)

	var expected SQLExpression
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "SELECT", 0, 1, 1),
		})
	var gce GroupedColumnsExpression
	gce.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "a", 8, 1, 9),
		})
	gce.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "b", 11, 1, 12),
		})
	expected.AppendExpression(&gce)
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "AS", 14, 1, 15),
		})
	expected.AppendExpression(
		NewOutputTargetExpression(
			makeToken(BITAND, "&", 17, 1, 18),
			&IdentityExpression{
				makeToken(IDENT, "Person", 18, 1, 19),
			},
			&IdentityExpression{
				makeToken(ASTERISK, "*", 25, 1, 26),
			},
		),
	)
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "FROM", 27, 1, 28),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "person", 32, 1, 33),
		})

	assert.NotEqual(t, r, nil)
	assert.Equal(t, &expected, r)
}

// Check that we return an error for empty groups
func TestErrorEmptyGroupedColumn(t *testing.T) {
	stmt := `SELECT () from person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, e := p.Run()

	assert.NotEqual(t, nil, e)
	t.Log(e)
	assert.NotEqual(t, nil, r)
}

// Check that we return an error for empty groups
func TestErrorBadFormedGroupedColumn(t *testing.T) {
	stmt := `SELECT (a, b,, c) from person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, e := p.Run()

	assert.NotEqual(t, nil, e)
	t.Log(e)
	assert.NotEqual(t, nil, r)
}

// Check that we return an error for empty groups
func TestErrorUnfinishedGroup(t *testing.T) {
	stmt := `SELECT (a, b, from person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, e := p.Run()

	assert.NotEqual(t, nil, e)
	t.Log(e)
	assert.NotEqual(t, nil, r)
}

func TestErrorEmptyStatement(t *testing.T) {
	stmt := ``

	l := NewLexer(stmt)
	p := NewParser(l)

	r, e := p.Run()

	assert.NotEqual(t, nil, e)
	t.Log(e)
	assert.NotEqual(t, nil, r)
}

// Check parsing and detection of output expressions
func TestParserSimpleOutputTarget(t *testing.T) {
	stmt := `SELECT &Person.* FROM person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	var expected SQLExpression
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "SELECT", 0, 1, 1),
		},
	)
	expected.AppendExpression(
		NewOutputTargetExpression(
			makeToken(BITAND, "&", 7, 1, 8),
			&IdentityExpression{
				makeToken(IDENT, "Person", 8, 1, 9),
			},
			&IdentityExpression{
				makeToken(ASTERISK, "*", 15, 1, 16),
			},
		),
	)
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "FROM", 17, 1, 18),
		},
	)
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "person", 22, 1, 23),
		},
	)
	assert.NotEqual(t, r, nil)
	assert.Equal(t, &expected, r)
}

func TestErrorMissingPeriodOutputTarget(t *testing.T) {
	stmt := `SELECT * as &Person* from person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, e := p.Run()

	assert.NotEqual(t, nil, e)
	t.Log(e)
	assert.NotEqual(t, nil, r)
}

// Check parsing and detection of input source expressions
func TestParserSimpleInputSource(t *testing.T) {
	stmt := `UPDATE person SET surname='Hitchens' WHERE id=$Person.id;`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	var expected SQLExpression
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "UPDATE", 0, 1, 1),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "person", 7, 1, 8),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "SET", 14, 1, 15),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "surname", 18, 1, 19),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(EQUAL, "=", 25, 1, 26),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(STRING, "'Hitchens'", 26, 1, 27),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "WHERE", 37, 1, 38),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(IDENT, "id", 43, 1, 44),
		})
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(EQUAL, "=", 45, 1, 46),
		})
	expected.AppendExpression(
		NewInputSourceExpression(
			makeToken(DOLLAR, "$", 46, 1, 47),
			&IdentityExpression{
				makeToken(IDENT, "Person", 47, 1, 48),
			},
			&IdentityExpression{
				makeToken(IDENT, "id", 54, 1, 55),
			},
		),
	)
	expected.AppendExpression(
		&IdentityExpression{
			makeToken(SEMICOLON, ";", 56, 1, 57),
		},
	)
	assert.NotEqual(t, r, nil)
	assert.Equal(t, &expected, r)
}
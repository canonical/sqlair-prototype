package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Use go test -v for printing the AST
func printAST(t *testing.T, e Expression, indentation int) {
	var ind string
	if indentation == 0 {
		t.Log("Dumping AST:")
	}
	for i := 0; i < indentation; i++ {
		ind = ind + " "
	}
	t.Log(ind, e.String(), "[", e.Type(), "]")
	for _, c := range e.Expressions() {
		printAST(t, c, indentation+4)
	}
}

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
	assert.Equal(t, r.Type(), SQL)
	// The rest of the expressions are Indent. Nothing to be reflected
	// in this statement.
	for _, c := range r.Expressions() {
		assert.Equal(t, c.Type(), Identity)
	}
}

func TestParserSimpleGroup(t *testing.T) {
	stmt := `SELECT (a, b) AS &Person.* FROM person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	expected := &SQLExpression{
		Children: []Expression{
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "SELECT",
					Pos: Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
			},
			&GroupedColumnsExpression{
				Children: []Expression{
					&IdentityExpression{
						Token{Type: IDENT,
							Literal: "a",
							Pos: Position{
								Offset: 8,
								Line:   1,
								Column: 9,
							},
						},
					},
					&IdentityExpression{
						Token{Type: IDENT,
							Literal: "b",
							Pos: Position{
								Offset: 11,
								Line:   1,
								Column: 12,
							},
						},
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "AS",
					Pos: Position{
						Offset: 14,
						Line:   1,
						Column: 15,
					},
				},
			},
			&OutputTargetExpression{
				Marker: Token{Type: BITAND,
					Literal: "&",
					Pos: Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
				Name: &IdentityExpression{
					Token{Type: IDENT,
						Literal: "Person",
						Pos: Position{
							Offset: 18,
							Line:   1,
							Column: 19,
						},
					},
				},
				Field: &IdentityExpression{
					Token{Type: ASTERISK,
						Literal: "*",
						Pos: Position{
							Offset: 25,
							Line:   1,
							Column: 26,
						},
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "FROM",
					Pos: Position{
						Offset: 27,
						Line:   1,
						Column: 28,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "person",
					Pos: Position{
						Offset: 32,
						Line:   1,
						Column: 33,
					},
				},
			},
		},
	}

	assert.NotEqual(t, r, nil)
	assert.Equal(t, expected, r)
}

func TestParserSimpleOutputTarget(t *testing.T) {
	stmt := `SELECT &Person.* FROM person`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	expected := &SQLExpression{
		Children: []Expression{
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "SELECT",
					Pos: Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
			},
			&OutputTargetExpression{
				Marker: Token{Type: BITAND,
					Literal: "&",
					Pos: Position{
						Offset: 7,
						Line:   1,
						Column: 8,
					},
				},
				Name: &IdentityExpression{
					Token{Type: IDENT,
						Literal: "Person",
						Pos: Position{
							Offset: 8,
							Line:   1,
							Column: 9,
						},
					},
				},
				Field: &IdentityExpression{
					Token{Type: ASTERISK,
						Literal: "*",
						Pos: Position{
							Offset: 15,
							Line:   1,
							Column: 16,
						},
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "FROM",
					Pos: Position{
						Offset: 17,
						Line:   1,
						Column: 18,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "person",
					Pos: Position{
						Offset: 22,
						Line:   1,
						Column: 23,
					},
				},
			},
		},
	}
	assert.NotEqual(t, r, nil)
	assert.Equal(t, expected, r)
}

func TestParserSimpleInputSource(t *testing.T) {
	stmt := `UPDATE person SET surname='Hitchens' WHERE id=$Person.id;`

	l := NewLexer(stmt)
	p := NewParser(l)

	r, _ := p.Run()

	printAST(t, r, 0)
	expected := &SQLExpression{
		Children: []Expression{
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "UPDATE",
					Pos: Position{
						Offset: 0,
						Line:   1,
						Column: 1,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "person",
					Pos: Position{
						Offset: 7,
						Line:   1,
						Column: 8,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "SET",
					Pos: Position{
						Offset: 14,
						Line:   1,
						Column: 15,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "surname",
					Pos: Position{
						Offset: 18,
						Line:   1,
						Column: 19,
					},
				},
			},
			&IdentityExpression{
				Token{Type: EQUAL,
					Literal: "=",
					Pos: Position{
						Offset: 25,
						Line:   1,
						Column: 26,
					},
				},
			},
			&IdentityExpression{
				Token{Type: STRING,
					Literal: "'Hitchens'",
					Pos: Position{
						Offset: 26,
						Line:   1,
						Column: 27,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "WHERE",
					Pos: Position{
						Offset: 37,
						Line:   1,
						Column: 38,
					},
				},
			},
			&IdentityExpression{
				Token{Type: IDENT,
					Literal: "id",
					Pos: Position{
						Offset: 43,
						Line:   1,
						Column: 44,
					},
				},
			},
			&IdentityExpression{
				Token{Type: EQUAL,
					Literal: "=",
					Pos: Position{
						Offset: 45,
						Line:   1,
						Column: 46,
					},
				},
			},
			&InputSourceExpression{
				Marker: Token{Type: DOLLAR,
					Literal: "$",
					Pos: Position{
						Offset: 46,
						Line:   1,
						Column: 47,
					},
				},
				Name: &IdentityExpression{
					Token{Type: IDENT,
						Literal: "Person",
						Pos: Position{
							Offset: 47,
							Line:   1,
							Column: 48,
						},
					},
				},
				Field: &IdentityExpression{
					Token{Type: IDENT,
						Literal: "id",
						Pos: Position{
							Offset: 54,
							Line:   1,
							Column: 55,
						},
					},
				},
			},
			&IdentityExpression{
				Token{Type: SEMICOLON,
					Literal: ";",
					Pos: Position{
						Offset: 56,
						Line:   1,
						Column: 57,
					},
				},
			},
		},
	}
	assert.NotEqual(t, r, nil)
	assert.Equal(t, expected, r)
}

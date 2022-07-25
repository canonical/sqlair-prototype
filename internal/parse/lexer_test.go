package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexerSimpleCorrectLiterals(t *testing.T) {
	// Note multiple spaces and carriage returns.
	stmt := `
  SELECT *  AS &Person.* 
FROM   person
WHERE address_id = $Address.id;
`

	expected := []string{
		"SELECT", "*", "AS", "&", "Person", ".", "*", "FROM", "person",
		"WHERE", "address_id", "=", "$", "Address", ".", "id", ";",
	}

	assert.Equal(t, expected, stringsFromTokens(tokensForStatement(stmt)))
}

func TestLexerSimpleCorrectNumbers(t *testing.T) {
	stmt := `
SELECT * AS &Person.* 
FROM   person
WHERE  salary = 100000.5;`

	expected := []string{
		"SELECT", "*", "AS", "&", "Person", ".", "*", "FROM", "person",
		"WHERE", "salary", "=", "100000.5", ";",
	}

	assert.Equal(t, expected, stringsFromTokens(tokensForStatement(stmt)))
}

func TestLexerBadNumber(t *testing.T) {
	stmt := `
SELECT * AS &Person.* 
FROM   person
WHERE  salary = 100000.5.6;`

	expected := []string{
		"SELECT", "*", "AS", "&", "Person", ".", "*", "FROM", "person",
		"WHERE", "salary", "=", "100000.5", ".", "6", ";",
	}

	assert.Equal(t, expected, stringsFromTokens(tokensForStatement(stmt)))
}

func TestLexerSimpleCorrectQuotedString(t *testing.T) {
	stmt := `
SELECT * AS &Person.* 
FROM   person
WHERE  name IN ('Lorn', 'Onos T''oolan', '', ''' ''');`

	expected := []string{
		"SELECT", "*", "AS", "&", "Person", ".", "*", "FROM", "person",
		"WHERE", "name", "IN", "(", "'Lorn'", ",", "'Onos T''oolan'", ",", "''", ",", "''' '''", ")", ";",
	}

	assert.Equal(t, expected, stringsFromTokens(tokensForStatement(stmt)))
}

func TestLexerCorrectPositions(t *testing.T) {
	stmt := `
SELECT *
FROM person`

	tokens := tokensForStatement(stmt)

	positions := []Position{
		{Offset: 0, Line: 1, Column: 1},
		{Offset: 7, Line: 1, Column: 8},
		{Offset: 9, Line: 2, Column: 1},
		{Offset: 14, Line: 2, Column: 6},
	}

	for i, token := range tokens {
		assert.Equal(t, positions[i], token.Pos)
	}
}

func TestLexerUnterminatedString(t *testing.T) {
	stmt := `select 's`

	tokens := tokensForStatement(stmt)

	assert.Equal(t, STRING, tokens[1].Type)
	assert.Equal(t, "'s", tokens[1].Literal)
	assert.Equal(t, 8, tokens[1].Pos.Column)
}

func TestLexerUnknownToken(t *testing.T) {
	stmt := `SELECT #a AS badtoken FROM t`

	expected := []string{
		"SELECT", "#", "a", "AS", "badtoken", "FROM", "t",
	}

	tokens := tokensForStatement(stmt)

	assert.Equal(t, expected, stringsFromTokens(tokens))
	assert.Equal(t, UNKNOWN, tokens[1].Type)
}

func tokensForStatement(stmt string) []Token {
	lex := NewLexer(stmt)

	var tokens []Token
	for token := lex.NextToken(); token.Type != EOF; token = lex.NextToken() {
		tokens = append(tokens, token)
	}

	return tokens
}

func stringsFromTokens(tokens []Token) []string {
	str := make([]string, len(tokens))
	for i, t := range tokens {
		str[i] = t.Literal
	}
	return str
}

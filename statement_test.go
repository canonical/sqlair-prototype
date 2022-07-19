package sqlair

import (
	"testing"

	"github.com/canonical/sqlair/internal/parse"
	sqlairtesting "github.com/canonical/sqlair/internal/testing"
	"github.com/stretchr/testify/assert"
)

func TestTypesForStatementUniqueNames(t *testing.T) {
	type DifferentPerson struct{}

	argTypes, err := typesForStatement([]any{DifferentPerson{}, sqlairtesting.Person{}})
	assert.Nil(t, err)

	assert.Len(t, argTypes, 2)
	assert.Equal(t, argTypes["DifferentPerson"].Name(), "DifferentPerson")
	assert.Equal(t, argTypes["Person"].Name(), "Person")
}

func TestTypesForStatementDuplicateNamesError(t *testing.T) {
	type Person struct{}

	_, err := typesForStatement([]any{Person{}, sqlairtesting.Person{}})
	assert.Error(t, err, NewErrTypeNameNotUnique("Person"))
}

// TODO (manadart 2022-07-15): The tests below are for verification during
// an intermediate stage. They will be subject to deletion shortly.

func TestInterpretFullCoverage(t *testing.T) {
	type address struct{}

	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}, address{}})
	assert.Nil(t, err)

	err = interpret(getExpression(), argTypes)
	assert.Nil(t, err)
}

func TestInterpretMissingTypesError(t *testing.T) {
	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}})
	assert.Nil(t, err)

	err = interpret(getExpression(), argTypes)
	assert.Error(t, err, NewErrTypeInfoNotPresent("address"))
}

func TestInterpretSuperfluousTypesError(t *testing.T) {
	type address struct{}
	type notUsed struct{}

	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}, address{}, notUsed{}})
	assert.Nil(t, err)

	err = interpret(getExpression(), argTypes)
	assert.Error(t, err, NewErrSuperfluousType("notUsed"))
}

func getExpression() parse.Expression {
	exp := &parse.SQLExpression{}

	tokens := tokensForStatement("&Person.*")
	exp.AppendExpression(parse.NewOutputTargetExpression(
		tokens[0], parse.NewIdentityExpression(tokens[1]), parse.NewIdentityExpression(tokens[3])))

	tokens = tokensForStatement("$address.id")
	exp.AppendExpression(parse.NewInputSourceExpression(
		tokens[0], parse.NewIdentityExpression(tokens[1]), parse.NewIdentityExpression(tokens[3])))

	return exp
}

func tokensForStatement(stmt string) []parse.Token {
	lex := parse.NewLexer(stmt)

	var tokens []parse.Token
	for token := lex.NextToken(); token.Type != parse.EOF; token = lex.NextToken() {
		tokens = append(tokens, token)
	}

	return tokens
}

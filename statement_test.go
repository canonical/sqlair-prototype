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

func TestValidateExpressionTypesFullCoverage(t *testing.T) {
	type address struct{}

	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}, address{}})
	assert.Nil(t, err)

	err = validateExpressionTypes(getExpression(), argTypes)
	assert.Nil(t, err)
}

func TestValidateExpressionTypesMissingTypesError(t *testing.T) {
	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}})
	assert.Nil(t, err)

	err = validateExpressionTypes(getExpression(), argTypes)
	assert.Error(t, err, NewErrTypeInfoNotPresent("address"))
}

func TestValidateExpressionTypesSuperfluousTypesError(t *testing.T) {
	type address struct{}
	type notUsed struct{}

	argTypes, err := typesForStatement([]any{sqlairtesting.Person{}, address{}, notUsed{}})
	assert.Nil(t, err)

	err = validateExpressionTypes(getExpression(), argTypes)
	assert.Error(t, err, NewErrSuperfluousType("notUsed"))
}

func getExpression() parse.Expression {
	return &sqlairtesting.SimpleExpression{
		T: parse.SQL,
		E: []*sqlairtesting.SimpleExpression{
			{
				T: parse.OutputTarget,
				E: []*sqlairtesting.SimpleExpression{
					{},
					{
						T: parse.Identity,
						S: "Person",
					},
				},
			},
			{
				T: parse.InputSource,
				E: []*sqlairtesting.SimpleExpression{
					{},
					{
						T: parse.Identity,
						S: "address",
					},
				},
			},
		},
	}
}

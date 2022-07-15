// Package parse_test is used to avoid an import loop.
// The testing package imports parse.
package parse_test

import (
	"errors"
	"testing"

	"github.com/canonical/sqlair/internal/parse"
	sqlairtesting "github.com/canonical/sqlair/internal/testing"
	"github.com/stretchr/testify/assert"
)

var _ parse.Expression = (*parse.IdentityExpression)(nil)

func TestIdentityExpression(t *testing.T) {
	literal := "identity"
	exp := parse.NewIdentityExpression(tokensForStatement(literal)[0])

	assert.Equal(t, literal, exp.String())
	assert.Nil(t, exp.Expressions())
}

var _ parse.TypeMappingExpression = (*parse.OutputTargetExpression)(nil)

func TestOutputTargetExpression(t *testing.T) {
	literal := "&Person.*"
	tokens := tokensForStatement(literal)
	exp := parse.NewOutputTargetExpression(
		tokens[0], parse.NewIdentityExpression(tokens[1]), parse.NewIdentityExpression(tokens[3]))

	assert.Equal(t, literal, exp.String())
	assert.Equal(t, "Person", exp.TypeName().String())
}

var _ parse.TypeMappingExpression = (*parse.InputSourceExpression)(nil)

func TestInputSourceExpression(t *testing.T) {
	literal := "$Address.id"
	tokens := tokensForStatement(literal)
	exp := parse.NewInputSourceExpression(
		tokens[0], parse.NewIdentityExpression(tokens[1]), parse.NewIdentityExpression(tokens[3]))

	assert.Equal(t, literal, exp.String())
	assert.Equal(t, "Address", exp.TypeName().String())
}

func TestWalk(t *testing.T) {
	expr := &sqlairtesting.SimpleExpression{
		T: parse.SQL,
		E: []*sqlairtesting.SimpleExpression{
			{
				T: parse.InputSource,
				E: []*sqlairtesting.SimpleExpression{
					{
						T: parse.PassThrough,
					},
				},
			},
			{
				T: parse.Identity,
			},
			{
				T: parse.GroupedColumns,
			},
		},
	}

	var types []parse.ExpressionType
	visit := func(e parse.Expression) error {
		if e.Type() == parse.Identity {
			return errors.New("stop")
		}
		types = append(types, e.Type())
		return nil
	}

	err := parse.Walk(expr, visit)

	// We expect to descend depth first into the expression tree,
	// and stop at the `Identity` expression.
	assert.NotNil(t, err)
	assert.Equal(t, []parse.ExpressionType{parse.SQL, parse.InputSource, parse.PassThrough}, types)
}

func tokensForStatement(stmt string) []parse.Token {
	lex := parse.NewLexer(stmt)

	var tokens []parse.Token
	for token := lex.NextToken(); token.Type != parse.EOF; token = lex.NextToken() {
		tokens = append(tokens, token)
	}

	return tokens
}

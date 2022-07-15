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

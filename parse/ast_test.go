// Package parse_test is used to avoid an import loop.
// The testing package imports parse.
package parse_test

import (
	"testing"

	"github.com/canonical/sqlair/parse"
	sqlairtesting "github.com/canonical/sqlair/testing"
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
	visit := func(e parse.Expression) bool {
		if e.Type() == parse.Identity {
			return false
		}
		types = append(types, e.Type())
		return true
	}

	finished := parse.Walk(expr, visit)

	// We expect to descend depth first into the expression tree,
	// and stop at the `Identity` expression.
	assert.False(t, finished)
	assert.Equal(t, []parse.ExpressionType{parse.SQL, parse.InputSource, parse.PassThrough}, types)
}

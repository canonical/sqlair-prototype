package parse

import (
	"testing"

	sqlairtesting "github.com/canonical/sqlair/testing"
	"github.com/stretchr/testify/assert"
)

func TestWalk(t *testing.T) {
	expr := &sqlairtesting.SimpleExpression{
		T: SQL,
		E: []*sqlairtesting.SimpleExpression{
			{
				T: InputSource,
				E: []*sqlairtesting.SimpleExpression{
					{
						T: PassThrough,
					},
				},
			},
			{
				T: Identity,
			},
			{
				T: GroupedColumns,
			},
		},
	}

	var types []ExpressionType
	visit := func(e Expression) bool {
		if e.Type() == Identity {
			return false
		}
		types = append(types, e.Type())
		return true
	}

	finished := Walk(expr, visit)

	// We expect to descend depth first into the expression tree,
	// and stop at the `Identity` expression.
	assert.False(t, finished)
	assert.Equal(t, []ExpressionType{SQL, InputSource, PassThrough}, types)
}

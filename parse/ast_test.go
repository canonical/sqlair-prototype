package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// testExp is a minimal implementation of Expression.
type testExp struct {
	t ExpressionType
	e []*testExp
}

func (e *testExp) Type() ExpressionType { return e.t }
func (e *testExp) Begin() Position       { return Position{} }
func (e *testExp) End() Position       { return Position{} }
func (e *testExp) String() string       { return "" }

func (e *testExp) Expressions() []Expression {
	res := make([]Expression, len(e.e))
	for i, exp := range e.e {
		res[i] = exp
	}
	return res
}

func TestWalk(t *testing.T) {
	expr := &testExp{
		t: SQL,
		e: []*testExp{
			{
				t: InputSource,
				e: []*testExp{
					{
						t: PassThrough,
					},
				},
			},
			{
				t: Identity,
			},
			{
				t: GroupedColumns,
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

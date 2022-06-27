package testing

import "github.com/canonical/sqlair/parse"

// SimpleExpression is a minimal implementation of parse.Expression.
type SimpleExpression struct {
	T parse.ExpressionType
	E []*SimpleExpression
	S string
}

func (e *SimpleExpression) Type() parse.ExpressionType { return e.T }
func (e *SimpleExpression) Begin() parse.Position      { return parse.Position{} }
func (e *SimpleExpression) End() parse.Position        { return parse.Position{} }
func (e *SimpleExpression) String() string             { return e.S }

func (e *SimpleExpression) Expressions() []parse.Expression {
	res := make([]parse.Expression, len(e.E))
	for i, exp := range e.E {
		res[i] = exp
	}
	return res
}

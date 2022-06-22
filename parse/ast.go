package parser

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Expression defines a type of AST node for outlining an expression.
type Expression interface {
	Pos() Position
	End() Position

	String() string
}

// QueryExpression represents a query full of expressions
type QueryExpression struct {
	Expressions []Expression
}

// Pos returns the first position of the query expression.
func (e *QueryExpression) Pos() Position {
	if len(e.Expressions) > 0 {
		return e.Expressions[0].Pos()
	}
	return Position{}
}

// End returns the last position of the query expression.
func (e *QueryExpression) End() Position {
	if num := len(e.Expressions); num > 0 {
		return e.Expressions[num-1].End()
	}
	return Position{}
}

func (e *QueryExpression) String() string {
	var out bytes.Buffer

	for _, s := range e.Expressions {
		out.WriteString(s.String())
	}

	return out.String()
}

// InputExpression.
// Represents things like $Person.name
type InputExpression struct {
	Dollar 	Expression
	Name	Expression
	Period	Expression
	Cols	Expression
}

func (ie *InputExpression) Pos() Position {
	return ie.Dollar.Pos()
		
}

func (ie *InputExpression) End() Position {
	return Cols.End()
}

func (ie *InputExpression) String() string {
	return Dollar.String() + Name.String() + Period.String() + Cols.String()
}

// OutputExpression.
// Represents things like &Person.*
type OutputExpression struct {
	Amp 	Expression
	Name	Expression
	Period	Expression
	Cols	Expression
}

func (oe *OutputExpression) Pos() Position {
	return oe.Amp.Pos()	
}

func (oe *OutputExpression) End() Position {
	return Cols.End()
}

func (oe *OutputExpression) String() string {
	return Amp.String() + Name.String() + Period.String() + Cols.String()
}

// PassthroughExpression.
// Represents all the things that are neither Input nor Output Expressions.
// Things we will not modify and we are not interested in. For instance:
// ---> SELECT A + 7 AS myalias... <---
type PassthroughExpression struct {
	Expressions []Expression
}

func (pe *PassthroughExpression) Pos() Position {
	if len(pe.Expressions) > 0 {
		return pe.Expressions[0].Pos()
	}
	return Position{}
}

func (pe *PassthroughExpression) End() Positions {
	if num := len(pe.Expressions); num > 0 {
		return pe.Expressions[num - 1].End()
	}
	return Position{}
}

// PassthroughElement
// Represents a single element inside a PassthroughExpression.
// This can be tokens of any type (identifiers, operators, etc...)
type PassthroughElement {
	Token Token
}

func (e *PassthroughElement) Pos() Position {
	return e.Token.Pos()
}

func (e *PassthroughElement) End() Position {
	return e.Token.End()
}

func (e *PassthroughElement) String() string {
	return e.Token.Literal
}

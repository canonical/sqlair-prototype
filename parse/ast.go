package parse

import (
	"bytes"
	"unicode/utf8"
)

// Definition of the AST for parsing the DSL of sqlair.
// Nodes are all expressions. We define three types of
// expressions:
//
//  * InputExpression:		for query arguments we want to use in the query.
//  * OutputExpression:		for placeholders we want to fill with the
//				results of a query.
//  * PassthroughExpression:	for elements that will be passed verbatim to the
//				underlying sql layer.

// Expression defines a type of AST node for outlining an expression.
// It is printable and holds references to the starting and ending position.
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
		return e.Expressions[num - 1].End()
	}
	return Position{}
}

// Printable method
func (e *QueryExpression) String() string {
	var out bytes.Buffer

	for _, s := range e.Expressions {
		out.WriteString(s.String())
	}

	return out.String()
}

// InputExpression.
// Represents something like $Person.name
type InputExpression struct {
	Dollar 	Element
	Name	Element
	Period	Element
	Cols	Element
}

// Pos returns the first position of the input expression.
func (ie *InputExpression) Pos() Position {
	return ie.Dollar.Pos()
}

// End returns the last position of the input expression.
func (ie *InputExpression) End() Position {
	return ie.Cols.End()
}

// Printable method
func (ie *InputExpression) String() string {
	return ie.Dollar.String() + ie.Name.String() + ie.Period.String() + ie.Cols.String()
}

// OutputExpression.
// Represents things like &Person.*
type OutputExpression struct {
	Amp 	Element
	Name	Element
	Period	Element
	Cols	Element
}

// Pos returns the first position of the output expression.
func (oe *OutputExpression) Pos() Position {
	return oe.Amp.Pos()	
}

// End returns the last position of the output expression.
func (oe *OutputExpression) End() Position {
	return oe.Cols.End()
}

// Printable method
func (oe *OutputExpression) String() string {
	return oe.Amp.String() + oe.Name.String() + oe.Period.String() + oe.Cols.String()
}

// PassthroughExpression.
// Represents all the things that are neither Input nor Output Expressions.
// Things we will not modify and we are not interested in. For instance:
// ---> SELECT A + 7 AS myalias... <---
type PassthroughExpression struct {
	Elements []Element
}

// Pos returns the first position of the passthrough expression.
func (pe *PassthroughExpression) Pos() Position {
	if len(pe.Elements) > 0 {
		return pe.Elements[0].Pos()
	}
	return Position{}
}

// End returns the last position of the passthrough expression.
func (pe *PassthroughExpression) End() Position {
	if num := len(pe.Elements); num > 0 {
		return pe.Elements[num - 1].End()
	}
	return Position{}
}

// Printable method
func (pe *PassthroughExpression) String() string {
	var out bytes.Buffer

	for _, s := range pe.Elements {
		out.WriteString(s.String())
	}

	return out.String()
}

// Element
// Represents a single element inside a PassthroughExpression.
// This can be tokens of any type (identifiers, operators, etc...)
type Element struct {
	Token Token
}

// Pos returns the first position of the passthrough element.
func (e *Element) Pos() Position {
	return e.Token.Pos
}

// End returns the last position of the passthrough element.
func (e *Element) End() Position {
	length := utf8.RuneCountInString(e.Token.Literal)
	return Position {
		Line: e.Token.Pos.Line,
		Column: e.Token.Pos.Column + length,
	}
}

// Printable method
func (e *Element) String() string {
	return e.Token.Literal
}

package parse

import (
	"bytes"
	"strings"
)

type ExpressionType int64

const (
	// SQL is a parent expression representing a
	// full structured query language query.
	SQL ExpressionType = iota

	// DML is a parent expression representing a data modification
	// language statement, i.e insert, update or delete.
	DML

	// DDL is a parent expression representing a data definition
	// language statement such as a table creation.
	DDL

	// GroupedColumns is an expression representing a list of columns to be
	// selected into a struct reference, as found within a query.
	// Example:
	// "(id, name)" in "SELECT (id, name) AS &Person.* FROM person;"
	GroupedColumns

	OutputTarget

	InputSource

	Identity

	// PassThrough is an expression representing a chunk of SQL, DML or SQL
	// that Sqlair will effectively ignore and pass to the DB as is.
	PassThrough
)

// Expression describes a token or tokens in a Sqlair DSL statement
// that represent a coherent, discrete subset of the DSL grammar.
type Expression interface {
	// Type indicates the type of this expression.
	Type() ExpressionType

	// Expressions returns the child expressions
	// that constitute this parent expression.
	Expressions() []Expression

	// Begin returns the starting position of the expression.
	Begin() Position

	// End returns the end position of the expression.
	End() Position

	// String returns the string that constitutes the expression.
	String() string
}

// TypeMappingExpression describes an expression that
// is for mapping inputs or outputs to Go types.
type TypeMappingExpression interface {
	Expression

	// TypeName returns the type name used in this expression,
	// such as "Person" in "&Person.*" or "$Person.id".
	TypeName() Expression
}

type SQLExpression struct {
	Children []Expression
}

func (sql *SQLExpression) Type() ExpressionType {
	return SQL
}

func (sql *SQLExpression) Expressions() []Expression {
	return sql.Children
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (sql *SQLExpression) Begin() Position {
	return beginChildren(sql.Children)
}

func (sql *SQLExpression) End() Position {
	return endChildren(sql.Children)
}

func (sql *SQLExpression) String() string {
	var sb strings.Builder
	for i, exp := range sql.Children {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(exp.String())
	}
	return sb.String()
}

type DMLExpression struct {
	Children []Expression
}

func (dml *DMLExpression) Type() ExpressionType {
	return DML
}

func (dml *DMLExpression) Expressions() []Expression {
	return dml.Children
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (dml *DMLExpression) Begin() Position {
	return beginChildren(dml.Children)
}

func (dml *DMLExpression) End() Position {
	return endChildren(dml.Children)
}

func (dml *DMLExpression) String() string {
	var sb bytes.Buffer
	for _, exp := range dml.Children {
		sb.WriteString(exp.String())
	}
	return sb.String()
}

type DDLExpression struct {
	Children []Expression
}

func (ddl *DDLExpression) Type() ExpressionType {
	return DDL
}

func (ddl *DDLExpression) Expressions() []Expression {
	return ddl.Children
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (ddl *DDLExpression) Begin() Position {
	return beginChildren(ddl.Children)
}

func (ddl *DDLExpression) End() Position {
	return endChildren(ddl.Children)
}

func (ddl *DDLExpression) String() string {
	var sb bytes.Buffer
	for _, exp := range ddl.Children {
		sb.WriteString(exp.String())
	}
	return sb.String()
}

type GroupedColumnsExpression struct {
	Children []Expression
}

func (gce *GroupedColumnsExpression) Type() ExpressionType {
	return GroupedColumns
}

func (gce *GroupedColumnsExpression) Expressions() []Expression {
	return gce.Children
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (gce *GroupedColumnsExpression) Begin() Position {
	return beginChildren(gce.Children)
}

func (gce *GroupedColumnsExpression) End() Position {
	return endChildren(gce.Children)
}

func (gce *GroupedColumnsExpression) String() string {
	var sb strings.Builder
	sb.WriteByte('(')
	for i, exp := range gce.Children {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(exp.String())
	}
	sb.WriteByte(')')
	return sb.String()
}

// OutputTargetExpression is an expression representing a type
// into which the output of a SQL query is to be mapped.
// Example:
// "&Person.*" in "SELECT &Person.* FROM person;"
type OutputTargetExpression struct {
	marker Token
	name   *IdentityExpression
	field  *IdentityExpression
}

// NewOutputTargetExpression returns a reference to a new
// OutputTargetExpression based on the input arguments.
func NewOutputTargetExpression(
	marker Token, name *IdentityExpression, field *IdentityExpression,
) *OutputTargetExpression {
	return &OutputTargetExpression{
		marker: marker,
		name:   name,
		field:  field,
	}
}

func (ote *OutputTargetExpression) Type() ExpressionType {
	return OutputTarget
}

// Expressions implements Expression by returning the child Expressions.
func (ote *OutputTargetExpression) Expressions() []Expression {
	return []Expression{ote.name, ote.field}
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (ote *OutputTargetExpression) Begin() Position {
	return ote.marker.Pos
}

func (ote *OutputTargetExpression) End() Position {
	return ote.field.End()
}

func (ote *OutputTargetExpression) String() string {
	return strings.Join([]string{ote.marker.Literal, ote.name.String(), ".", ote.field.String()}, "")
}

func (ote *OutputTargetExpression) TypeName() Expression {
	return ote.name
}

// InputSourceExpression is an expression representing a type
// from which parameters of a statement are to be sourced.
// Example:
// "$Person.id" in "UPDATE person SET surname='Hitchens' WHERE id=$Person.id;"
type InputSourceExpression struct {
	marker Token
	name   Expression
	field  Expression
}

// NewInputSourceExpression returns a reference to a new
// InputSourceExpression based on the input arguments.
func NewInputSourceExpression(
	marker Token, name *IdentityExpression, field *IdentityExpression,
) *InputSourceExpression {
	return &InputSourceExpression{
		marker: marker,
		name:   name,
		field:  field,
	}
}

func (ise *InputSourceExpression) Type() ExpressionType {
	return InputSource
}

// Expressions implements Expression by returning the child Expressions.
func (ise *InputSourceExpression) Expressions() []Expression {
	return []Expression{ise.name, ise.field}
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (ise *InputSourceExpression) Begin() Position {
	return ise.marker.Pos
}

func (ise *InputSourceExpression) End() Position {
	return ise.field.End()
}

func (ise *InputSourceExpression) String() string {
	return strings.Join([]string{ise.marker.Literal, ise.name.String(), ".", ise.field.String()}, "")
}

func (ise *InputSourceExpression) TypeName() Expression {
	return ise.name
}

// IdentityExpression is an expression that identifies a single entity.
type IdentityExpression struct {
	token Token
}

// NewIdentityExpression returns a reference to a new
// IdentityExpression based on the input Token.
func NewIdentityExpression(token Token) *IdentityExpression {
	return &IdentityExpression{token: token}
}

func (ie *IdentityExpression) Type() ExpressionType {
	return Identity
}

// Expressions implements Expression by returning the child Expressions.
func (ie *IdentityExpression) Expressions() []Expression {
	return nil
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (ie *IdentityExpression) Begin() Position {
	return ie.token.Pos
}

func (ie *IdentityExpression) End() Position {
	return Position{
		Offset: ie.token.Pos.Offset + len(ie.token.Literal),
	}
}

func (ie *IdentityExpression) String() string {
	return ie.token.Literal
}

type PassThroughExpression struct {
	Children []Expression
}

func (pt *PassThroughExpression) Type() ExpressionType {
	return PassThrough
}

// Expressions implements Expression by returning the child Expressions.
func (pt *PassThroughExpression) Expressions() []Expression {
	return pt.Children
}

// Begin implements Expression by returning the
// Position of this Expression's first Token.
func (pt *PassThroughExpression) Begin() Position {
	return beginChildren(pt.Children)
}

func (pt *PassThroughExpression) End() Position {
	return endChildren(pt.Children)
}

func (pt *PassThroughExpression) String() string {
	var sb strings.Builder
	for _, exp := range pt.Children {
		sb.WriteString(exp.String())
	}
	return sb.String()
}

// Walk recursively iterates depth-first over the input expression tree,
// calling the input function for each visited expression.
// If it returns an error, the iteration terminates.
func Walk(parent Expression, visit func(Expression) error) error {
	if err := visit(parent); err != nil {
		return err
	}
	for _, child := range parent.Expressions() {
		if err := Walk(child, visit); err != nil {
			return err
		}
	}
	return nil
}

// beginChildren is a helper method that returns the position of the
// first expression in the children array or an empty position if the
// array is empty
func beginChildren(children []Expression) Position {
	if len(children) > 0 {
		return children[0].Begin()
	}
	return Position{}
}

// endChildren is a helper method that returns the position of the
// last expression in the children array or an empty position if the
// array is empty
func endChildren(children []Expression) Position {
	if l := len(children); l > 0 {
		return children[l-1].End()
	}
	return Position{}
}


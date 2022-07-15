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

	// OutputTarget is an expression representing a type into
	// which the output of a SQL query is to be mapped.
	// Example:
	// "&Person.*" in "SELECT &Person.* FROM person;"
	OutputTarget

	// InputSource is an expression representing a type from
	// which parameters of a statement are to be sourced.
	// Example:
	// "$Person.id" in "UPDATE person SET surname='Hitchens' WHERE id=$Person.id;"
	InputSource

	// Identity is an expression that identifies a single entity.
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

// beginChildren is a helper method that returns the position of the
// first expression in the children array or an empty position if the
// array is empty
func beginChildren(children []Expression) Position {
	var p Position
	if len(children) > 0 {
		p = children[0].Begin()
	}
	return p
}

// endChildren is a helper method that returns the position of the
// last expression in the children array or an empty position if the
// array is empty
func endChildren(children []Expression) Position {
	var p Position
	if l := len(children); l > 0 {
		p = children[l-1].End()
	}
	return p
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

type OutputTargetExpression struct {
	Marker Token
	Name   Expression
	Field  Expression
}

func (ote *OutputTargetExpression) Type() ExpressionType {
	return OutputTarget
}

func (ote *OutputTargetExpression) Expressions() []Expression {
	return []Expression{ote.Name, ote.Field}
}

func (ote *OutputTargetExpression) Begin() Position {
	return ote.Marker.Pos
}

func (ote *OutputTargetExpression) End() Position {
	return ote.Field.End()
}

func (ote *OutputTargetExpression) String() string {
	var sb strings.Builder
	sb.WriteString(ote.Marker.Literal)
	sb.WriteString(ote.Name.String())
	sb.WriteByte('.')
	sb.WriteString(ote.Field.String())
	return sb.String()
}

type InputSourceExpression struct {
	Marker Token
	Name   Expression
	Field  Expression
}

func (ise *InputSourceExpression) Type() ExpressionType {
	return InputSource
}

func (ise *InputSourceExpression) Expressions() []Expression {
	return []Expression{ise.Name, ise.Field}
}

func (ise *InputSourceExpression) Begin() Position {
	return ise.Marker.Pos
}

func (ise *InputSourceExpression) End() Position {
	return ise.Field.End()
}

func (ise *InputSourceExpression) String() string {
	var sb strings.Builder
	sb.WriteString(ise.Marker.Literal)
	sb.WriteString(ise.Name.String())
	sb.WriteByte('.')
	sb.WriteString(ise.Field.String())
	return sb.String()
}

type IdentityExpression struct {
	Token Token
}

func (ie *IdentityExpression) Type() ExpressionType {
	return Identity
}

func (ie *IdentityExpression) Expressions() []Expression {
	return nil
}

func (ie *IdentityExpression) Begin() Position {
	return ie.Token.Pos
}

func (ie *IdentityExpression) End() Position {
	return Position{
		Offset: ie.Token.Pos.Offset + len(ie.Token.Literal),
	}
}

func (ie *IdentityExpression) String() string {
	return ie.Token.Literal
}

type PassThroughExpression struct {
	Children []Expression
}

func (pt *PassThroughExpression) Type() ExpressionType {
	return PassThrough
}

func (pt *PassThroughExpression) Expressions() []Expression {
	return pt.Children
}

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

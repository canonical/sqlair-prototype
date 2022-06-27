package parse

import (
	"bytes"
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

	// Returns the starting position of the expression.
	Begin()	Position

	// Returns the end position of the expression.
	End() Position

	// String returns the string that constitutes the expression.
	String() string
}

/*
	SQL Expression
*/
type SQLExpression struct {
	children []Expression
}

func (self *SQLExpression) Type() ExpressionType {
	return SQL
}

func (self *SQLExpression) Expressions() []Expression {
	return self.children
}

func (self *SQLExpression) Begin() Position {
	var p Position
	if len(self.children) > 0 {
		p = self.children[0].Begin()
	}
	return p
}

func (self *SQLExpression) End() Position {
	var p Position
	if l:=len(self.children); l > 0 {
		p = self.children[l - 1].End()
	}
	return p
}

func (self *SQLExpression) String() string {
	var os bytes.Buffer
	for _, exp:= range self.children {
		os.WriteString(exp.String())
	}
	return os.String()
}

/*
	DML Expression
*/
type DMLExpression struct {
	children []Expression
}

func (self *DMLExpression) Type() ExpressionType {
	return DML
}

func (self *DMLExpression) Expressions () []Expression {
	return self.children
}

func (self *DMLExpression) Begin() Position {
	var p Position
	if len(self.children) > 0 {
		p = self.children[0].Begin()
	}
	return p
}

func (self *DMLExpression) End() Position {
	var p Position
	if l:=len(self.children); l > 0 {
		p = self.children[l - 1].End()
	}
	return p
}

func (self *DMLExpression) String() string {
	var os bytes.Buffer
	for _, exp:= range self.children {
		os.WriteString(exp.String())
	}
	return os.String()
}

/*
	DDL Expression
*/
type DDLExpression struct {
	children []Expression
}

func (self *DDLExpression) Type() ExpressionType {
	return DDL
}

func (self *DDLExpression) Expressions () []Expression {
	return self.children
}

func (self *DDLExpression) Begin() Position {
	var p Position
	if len(self.children) > 0 {
		p = self.children[0].Begin()
	}
	return p
}

func (self *DDLExpression) End() Position {
	var p Position
	if l:=len(self.children); l > 0 {
		p = self.children[l - 1].End()
	}
	return p
}

func (self *DDLExpression) String() string {
	var os bytes.Buffer
	for _, exp:= range self.children {
		os.WriteString(exp.String())
	}
	return os.String()
}


/*
	GroupedColumns Expression
*/
type GroupedColumnsExpression struct {
	children []Expression
}

func (self *GroupedColumnsExpression) Type() ExpressionType {
	return GroupedColumns
}

func (self *GroupedColumnsExpression) Expressions () []Expression {
	return self.children
}

func (self *GroupedColumnsExpression) Begin() Position {
	var p Position
	if len(self.children) > 0 {
		p = self.children[0].Begin()
	}
	return p
}

func (self *GroupedColumnsExpression) End() Position {
	var p Position
	if l:=len(self.children); l > 0 {
		p = self.children[l - 1].End()
	}
	return p
}

func (self *GroupedColumnsExpression) String() string {
	var os bytes.Buffer
	os.WriteString("(")
	for _, exp:= range self.children {
		os.WriteString(exp.String())
		os.WriteString(",")
	}
	os.WriteString(")")
	return os.String()
}


/*
	OutputTarget Expression
*/
type OutputTargetExpression struct {
	Marker	Token
	Name 	Expression
	Period	Token
	Field	Expression
}

func (self *OutputTargetExpression) Type() ExpressionType {
	return OutputTarget
}

func (self *OutputTargetExpression) Expressions () []Expression {
	return []Expression {self.Name, self.Field}
}

func (self *OutputTargetExpression) Begin() Position {
	return self.Marker.Pos
}

func (self *OutputTargetExpression) End() Position {
	return self.Field.End()
}

func (self *OutputTargetExpression) String() string {
	var os bytes.Buffer
	os.WriteString(self.Marker.Literal)
	os.WriteString(self.Name.String())
	os.WriteString(self.Period.Literal)
	os.WriteString(self.Field.String())
	return os.String()
}

/*
	InputSource Expression
*/

type InputSourceExpression struct {
	Marker	Token
	Name 	Expression
	Period	Token
	Field	Expression
}

func (self *InputSourceExpression) Type() ExpressionType {
	return InputSource
}

func (self *InputSourceExpression) Expressions () []Expression {
	return []Expression {self.Name, self.Field}
}

func (self *InputSourceExpression) String() string {
	var os bytes.Buffer
	os.WriteString(self.Marker.Literal)
	os.WriteString(self.Name.String())
	os.WriteString(self.Period.Literal)
	os.WriteString(self.Field.String())
	return os.String()
}

/*
	Identity Expression
*/

type IdentityExpression struct {
	Token Token
}

func (self *IdentityExpression) Type() ExpressionType {
	return Identity
}

func (self *IdentityExpression) Expressions () []Expression {
	return nil
}

func (self *IdentityExpression) String() string {
	return self.Token.Literal
}


/*
	PassThrough Expression
*/

type PassThroughExpression struct {
	children []Expression
}

func (self *PassThroughExpression) Type() ExpressionType {
	return PassThrough
}

func (self *PassThroughExpression) Expressions () []Expression {
	return self.children
}

func (self *PassThroughExpression) String() string {
	var os bytes.Buffer
	for _, exp:= range self.children {
		os.WriteString(exp.String())
	}
	return os.String()
}

// Walk recursively iterates depth-first over the input expression tree,
// calling the input function. If it returns false, the iteration terminates.
func Walk(parent Expression, visit func(Expression) bool) bool {
	if !visit(parent) {
		return false
	}
	for _, child := range parent.Expressions() {
		if !Walk(child, visit) {
			return false
		}
	}
	return true
}

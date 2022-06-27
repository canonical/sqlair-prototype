package parse

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

	// String returns the string that constitutes the expression.
	String() string
}

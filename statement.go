package sqlair

import (
	"github.com/canonical/sqlair/parse"
	sqlairreflect "github.com/canonical/sqlair/reflect"
)

// types is a convenience type alias for reflection
// information indexed by type name.
type types = map[string]sqlairreflect.Info

// Statement represents a prepared Sqlair DSL statement
// that can be executed by the database.
type Statement struct {
	// expression is the parsed expression tree for this statement.
	expression parse.Expression

	// argTypes holds the reflection info for types used in this statement.
	argTypes types
}

// Prepare accepts a raw DSL string and optionally,
// objects from which to infer type information.
// - The string is parsed into an expression tree.
// - Any input objects have their reflection information retrieved/generated.
// - The reflection information is matched with the parser output to generate
//   a Statement that can be passed to the database for execution.
func Prepare(stmt string, args ...any) (*Statement, error) {
	lex := parse.NewLexer(stmt)
	parser := parse.NewParser(lex)

	exp, err := parser.Run()
	if err != nil {
		return nil, err
	}

	argTypes, err := typesForStatement(args)
	if err != nil {
		return nil, err
	}

	if err := validateExpressionTypes(exp, argTypes); err != nil {
		return nil, err
	}

	return &Statement{
		expression: exp,
		argTypes:   argTypes,
	}, nil
}

// typesForStatement returns reflection information for the input arguments.
// The reflected type name of each argument must be unique in the list,
// which means declaring new local types to avoid ambiguity.
//
// Example:
//
//     type Person struct{}
//     type Manager Person
//
//     stmt, err := sqlair.Prepare(`
//     SELECT p.* AS &Person.*,
//            m.* AS &Manager.*
//       FROM person AS p
//       JOIN person AS m
//         ON p.manager_id = m.id
//      WHERE p.name = 'Fred'`, Person{}, Manager{})
//
func typesForStatement(args []any) (types, error) {
	c := sqlairreflect.Cache()
	argTypes := make(types)

	for _, arg := range args {
		reflected, err := c.Reflect(arg)
		if err != nil {
			return nil, err
		}

		name := reflected.Name()
		if _, ok := argTypes[name]; ok {
			return nil, NewErrTypeNameNotUnique(name)
		}

		argTypes[name] = reflected
	}

	return argTypes, nil
}

// validateExpressionTypes walks the input expression tree to ensure:
// - Each input/output target in expression has type information in argTypes.
// - All type information is actually required by the input/output targets.
func validateExpressionTypes(statementExp parse.Expression, argTypes types) error {
	var err error
	seen := make(map[string]bool)

	visit := func(exp parse.Expression) bool {
		if t := exp.Type(); t != parse.OutputTarget && t != parse.InputSource {
			return true
		}

		// Select the first identity, such as "Person"
		// in the case of "$Person.id".
		// Ensure that there is type information for it.
		typeName := exp.Expressions()[1].String()
		if _, ok := argTypes[typeName]; !ok {
			err = NewErrTypeInfoNotPresent(typeName)
			return false
		}

		seen[typeName] = true
		return true
	}

	// If we did not complete the walk through the tree,
	// return the error that we encountered.
	if !parse.Walk(statementExp, visit) {
		return err
	}

	// Now compare the type names that we saw against what we have information
	// for. If unused types were supplied, it is an error condition.
	for name := range argTypes {
		if _, ok := seen[name]; !ok {
			return NewErrSuperfluousType(name)
		}
	}

	return nil
}

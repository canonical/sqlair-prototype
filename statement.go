package sqlair

import (
	"github.com/canonical/sqlair/internal/parse"
	sqlairreflect "github.com/canonical/sqlair/internal/reflect"
)

// typeMap is a convenience type alias for reflection
// information indexed by type name.
type typeMap = map[string]sqlairreflect.Info

// Statement represents a prepared Sqlair DSL statement
// that can be executed by the database.
type Statement struct {
	// expression is the parsed expression tree for this statement.
	expression parse.Expression

	// argTypes holds the reflection info for types used in this statement.
	argTypes typeMap
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

	if err := interpret(exp, argTypes); err != nil {
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
func typesForStatement(args []any) (typeMap, error) {
	c := sqlairreflect.Cache()
	argTypes := make(typeMap)

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

// interpret walks the input expression tree to ensure:
// - Each input/output target in expression has type information in argTypes.
// - All type information is actually required by the input/output targets.
// - TODO (manadart 2022-07-15): Add further interpreter behaviour.
func interpret(statementExp parse.Expression, argTypes typeMap) error {
	var err error
	seen := make(map[string]bool)

	visit := func(exp parse.Expression) error {
		switch e := exp.(type) {
		case *parse.OutputTargetExpression, *parse.InputSourceExpression:
			if seen, err = validateExpressionType(e.(parse.TypeMappingExpression), argTypes, seen); err != nil {
				return err
			}
		}

		return nil
	}

	if err := parse.Walk(statementExp, visit); err != nil {
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

// validateExpressionType ensures that the type name identity from the input
// expression is present in the input type information. If it is not, an error
// is returned. The list of seen types is updated and returned.
func validateExpressionType(
	exp parse.TypeMappingExpression, argTypes typeMap, seen map[string]bool,
) (map[string]bool, error) {
	typeName := exp.TypeName().String()
	if _, ok := argTypes[typeName]; !ok {
		return seen, NewErrTypeInfoNotPresent(typeName)
	}

	seen[typeName] = true
	return seen, nil
}

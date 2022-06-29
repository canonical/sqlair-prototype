package sqlair

import (
	"github.com/canonical/sqlair/parse"
	sqlairreflect "github.com/canonical/sqlair/reflect"
)

// Statement represents a prepared Sqlair DSL statement
// that can be executed by the database.
type Statement struct {
}

// Prepare accepts a raw DSL string and optionally,
// objects from which to infer type information.
// - The string is parsed.
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

	reflected := make([]sqlairreflect.Kind, len(args))
	c := sqlairreflect.Cache()
	for i, arg := range args {
		if reflected[i], err = c.Reflect(arg); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

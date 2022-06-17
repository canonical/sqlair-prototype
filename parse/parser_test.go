package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
  Check that we can correctly parse a single argument
  in the statement
*/
func TestParserSimple(t *testing.T) {
	stmt := `
		SELECT * AS &Person.* 
		FROM person
		WHERE address_id = $Address.id;
	`

	p := NewParser(stmt)

	ret, err := p.Parse()

	assert.Equal(t, 1, len(ret))
	expected := QueryArgument {
		name:"Person",
		field:"all",
		from:13,
		to:21,
	}
	assert.Equal(t, expected, ret[0])
	assert.Nil(t, nil, err)
}

/*
  Check that we can correctly parse multiple arguments
  in the statement
*/
func TestParserMultiple(t *testing.T) {
	stmt := `
		SELECT p.* AS &Person.*, 
		(a.district, a.street) AS &Address.*
		FROM person AS p 
		JOIN address AS a
		ON p.address_id = a.id
		WHERE p.name = 'Fred'
	`

	p := NewParser(stmt)

	ret, err := p.Parse()

	assert.Nil(t, nil, err)
	assert.Equal(t, 2, len(ret))

	expectedPerson := QueryArgument {
		name:"Person",
		field:"all",
		from:15,
		to:23,
	}
	expectedAddress := QueryArgument {
		name:"Address",
		field:"all",
		from:55,
		to:64,
	}

	assert.Equal(t, expectedPerson, ret[0])
	assert.Equal(t, expectedAddress, ret[1])
}

/*
  Check that we dont validate something that is
  not a real placeholder
*/
func TestParserIgnoreAmpIdent(t *testing.T) {
	/*
	  Note that we are not checking the (incorrect) SQL syntax.
	  We are only interested in checking if we mistakingly
	  thing &a is a placeholder as understood by our DSL.
	*/
	stmt := `
		SELECT &Person.*, &a
		FROM person
		WHERE p.name = 'Fred'
	`

	p := NewParser(stmt)

	ret, err := p.Parse()

	assert.Nil(t, nil, err)
	assert.Equal(t, 1, len(ret))

	expectedPerson := QueryArgument {
		name:"Person",
		field:"all",
		from:8,
		to:16,
	}

	assert.Equal(t, expectedPerson, ret[0])
}

/*
  Check that we dont validate something that is
  not a real placeholder
*/
//func TestParserIgnoreAmpIdentPeriod(t *testing.T) {
//	/*
//	  Note that we are not checking the (incorrect) SQL syntax.
//	  We are only interested in checking if we mistakingly
//	  thing &a is a placeholder as understood by our DSL.
//	*/
//	stmt := `
//		SELECT &Person.*, &a.
//		FROM person
//		WHERE p.name = 'Fred'
//	`
//
//	p := NewParser(stmt)
//
//	ret, err := p.Parse()
//
//	log.Println(ret[0])
//	log.Println(ret[1])
//
//	assert.Nil(t, nil, err)
//	assert.Equal(t, 1, len(ret))
//
//	expectedPerson := QueryArgument{name:"Person", field:"all", from:8, to:16}
//	assert.Equal(t, expectedPerson, ret[0])
//}

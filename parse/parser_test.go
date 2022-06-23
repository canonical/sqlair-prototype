package parse

import (
	//"fmt"
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

	assert.Equal(t, 2, len(ret.Expressions))
	//fmt.Print(ret.Expressions)
	assert.Nil(t, nil, err)
}


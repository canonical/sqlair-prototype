package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMaybeRuneTokenTrue(t *testing.T) {
	pos := Position{Offset: 21, Line: 3}
	token, isRune := maybeRuneToken(')', pos)

	assert.True(t, isRune)
	assert.Equal(t, Token{
		Type:    RPAREN,
		Literal: ")",
		Pos:     pos,
	}, token)
}

func TestMaybeRuneTokenFalse(t *testing.T) {
	_, isRune := maybeRuneToken('6', Position{})
	assert.False(t, isRune)
}

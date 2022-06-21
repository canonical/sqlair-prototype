package reflect

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestReflectSimple(t *testing.T) {
	var i int64

	info, err := NewCache().Reflect(i)
	assert.Nil(t, err)

	assert.Equal(t, reflect.Int64, info.Kind())

	_, ok := info.(Value)
	assert.True(t, ok)
}

func TestReflectStruct(t *testing.T) {
	type something struct {
		ID      int64  `db:"id"`
		Name    string `db:"name,omitempty"`
		NotInDB string
	}

	s := something{
		ID:      99,
		Name:    "Chainheart Machine",
		NotInDB: "doesn't matter",
	}

	info, err := NewCache().Reflect(s)
	assert.Nil(t, err)

	assert.Equal(t, reflect.Struct, info.Kind())

	st, ok := info.(Struct)
	assert.True(t, ok)

	assert.Len(t, st.Fields, 2)

	id, ok := st.Fields["id"]
	assert.True(t, ok)
	assert.Equal(t, "ID", id.Name)
	assert.False(t, id.OmitEmpty)

	name, ok := st.Fields["name"]
	assert.True(t, ok)
	assert.Equal(t, "Name", name.Name)
	assert.True(t, name.OmitEmpty)
}

func TestReflectBadTagError(t *testing.T) {
	type something struct {
		ID int64 `db:"id,bad-juju"`
	}

	s := something{ID: 99}

	_, err := NewCache().Reflect(s)
	assert.Error(t, errors.New(`unexpected tag value "bad-juju"`), err)
}

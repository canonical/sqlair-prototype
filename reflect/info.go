package reflect

import (
	"reflect"
)

// Kind describes the ability to return the reflect.Kind of the receiver.
type Kind interface {
	Kind() reflect.Kind
}

// Value represents reflection information for a simple type.
// It wraps a reflect.Value in order to implement Kind.
type Value struct {
	value reflect.Value
}

// Kind returns the Value's reflect.Kind.
func (r Value) Kind() reflect.Kind {
	return r.value.Kind()
}

// Field represents a single field from a struct type.
type Field struct {
	// Name is the name of the struct field.
	Name string

	// OmitEmpty is true when "omitempty" is
	// a property of the field's "db" tag.
	OmitEmpty bool
	value     reflect.Value
}

// Struct represents reflected information about a struct type.
type Struct struct {
	// Name is the name of the struct's type.
	Name string

	// Fields maps "db" tags to struct fields.
	// Sqlair does not care about fields without a "db" tag.
	Fields map[string]Field
	value  reflect.Value
}

// Kind returns the Struct's reflect.Kind.
func (r Struct) Kind() reflect.Kind {
	return r.value.Kind()
}

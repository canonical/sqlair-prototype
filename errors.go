package sqlair

import "fmt"

// ErrTypeNameNotUnique is an error indicating that the objects
// passed as arguments to statement preparation do not constitute
// a list of unique type names.
type ErrTypeNameNotUnique struct {
	name string
}

// NewErrTypeNameNotUnique returns a new error
// for the input non-unique type name.
func NewErrTypeNameNotUnique(name string) error {
	return &ErrTypeNameNotUnique{name: name}
}

// Error implements error, returning a message indicating the non-unique type.
func (e *ErrTypeNameNotUnique) Error() string {
	return fmt.Sprintf("names for supplied types are not unique; %q is ambiguous", e.name)
}

// ErrTypeInfoNotPresent is an error indicating that the one of the
// names used in a DSL statement, for input or output type mapping,
// does not have associated type information.
type ErrTypeInfoNotPresent struct {
	name string
}

// NewErrTypeInfoNotPresent returns a new error
// for the input type without associated info.
func NewErrTypeInfoNotPresent(name string) error {
	return &ErrTypeInfoNotPresent{name: name}
}

// Error implements error, returning a message
// indicating the unrepresented identity.
func (e *ErrTypeInfoNotPresent) Error() string {
	return fmt.Sprintf("identity %q has no associated object from which to derive type information", e.name)
}

// ErrSuperfluousType is an error indicating that an object was supplied as an
// argument to Prepare, but its type information is not reflected in the
// parsed expression, making it redundant.
type ErrSuperfluousType struct {
	name string
}

// NewErrSuperfluousType returns a new error
// for the input redundant type name.
func NewErrSuperfluousType(name string) error {
	return &ErrSuperfluousType{name: name}
}

// Error implements error, returning a message
// indicating the redundant identity.
func (e *ErrSuperfluousType) Error() string {
	return fmt.Sprintf("type with name %q was supplied, but is not used in the statement", e.name)
}

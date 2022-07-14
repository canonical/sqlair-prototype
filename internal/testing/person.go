package testing

// Person is a simple test subject for inputs to sqlair.Prepare.
type Person struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}
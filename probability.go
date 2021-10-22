// Collections where returned values are based on probability.
package probability

// Selector is a type that can randomly select a value based on the value's
// probability.
type Selector interface {
	// Select gets a random value.
	Select() interface{}
}

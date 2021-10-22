// Markov contains types to handle Markov chains
package markov

// Chain is created from a feed of tokens and can create a generator that
// uses probability to return tokens.
type Chain interface {
	// MakeGenerator creates a generator function that acts as an interator
	// to create tokens.
	MakeGenerator() Generator
	// Seed seeds the random value.
	Seed(seed int64)
	// Feed feeds data into the Markov chain and returns an error if the
	// feed fails.
	Feed(interface{}) error
}

// Generator generates tokens based on a Markov chain and returns a boolean
// to determine if the generation has completed (stopped).
type Generator interface {
	// Next returns the next token, and a boolean if there are not more tokens
	// to be returned.
	Next() (token interface{}, stop bool)
	// HasNext returns if there is a next token to return. This is a convenience
	// that should return the opposite of Next's stop value.
	HasNext() bool
}

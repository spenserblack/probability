package probability

import "math/rand"

// ByCount defines probability by the count of inserted values.
type ByCount struct {
	m     map[interface{}]int
	total int
	// TODO Set random value with seed
	r *rand.Rand
}

// NewByCount initializes a new ByCount struct.
func NewByCount() *ByCount {
	m := make(map[interface{}]int)
	return &ByCount{m, 0, rand.New(rand.NewSource(1))}
}

// Seed sets the seed for for the randomization. This can be useful,
// as ByCount is always initialized with the same seed, making it deterministic.
func (bc *ByCount) Seed(seed int64) {
	bc.r.Seed(seed)
}

// Insert increments the count of the matching value, and returns the count.
func (bc *ByCount) Insert(v interface{}) (count int) {
	bc.total += 1
	bc.m[v] += 1
	return bc.m[v]
}

// Keys returns all of the keys that have been counted.
func (bc ByCount) Keys() []interface{} {
	keys := make([]interface{}, 0, len(bc.m))
	for key := range bc.m {
		keys = append(keys, key)
	}
	return keys
}

// Get gets the count of the provided value.
func (bc ByCount) Get(v interface{}) int {
	return bc.m[v]
}

// Total gets the total count of all values.
func (bc ByCount) Total() int {
	return bc.total
}

// Probability returns the probability of a single value being picked from
// all of the values.
func (bc ByCount) Probability(v interface{}) float32 {
	return float32(bc.Get(v)) / float32(bc.Total())
}

// Probabilities returns a map of all values and their probability.
func (bc ByCount) Probabilities() (probabilities map[interface{}]float32) {
	probabilities = make(map[interface{}]float32)
	for v := range bc.m {
		probabilities[v] = bc.Probability(v)
	}
	return
}

// Select gets a random value based on the count of that value.
// It can be thought of as pulling a value out of a bag of values.
func (bc ByCount) Select() interface{} {
	if bc.total == 0 {
		return nil
	}
	index := bc.r.Intn(bc.total)
	bottom := 0
	for value, count := range bc.m {
		top := bottom + count
		if bottom <= index && index < top {
			return value
		}
		bottom = top
	}
	panic("Unreachable")
	return nil
}

package probability

// ByCount defines probability by the count of inserted values.
type ByCount struct {
	m     map[interface{}]int
	total int
}

// NewByCount initializes a new ByCount struct.
func NewByCount() *ByCount {
	m := make(map[interface{}]int)
	return &ByCount{m, 0}
}

// Insert increments the count of the matching value, and returns the count.
func (bc *ByCount) Insert(v interface{}) (count int) {
	bc.total += 1
	bc.m[v] += 1
	return bc.m[v]
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
	for v := range bc.m {
		probabilities[v] = bc.Probability(v)
	}
	return
}

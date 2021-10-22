package probability_test

import (
	"fmt"

	"github.com/spenserblack/probability"
)

// Probabilities are set based on the number of times each value has been
// inserted. You can think of it as a bag of marbles, or a shuffled deck.
func ExampleByCount() {
	const (
		Red int = iota
		Green
		Blue
	)

	byCount := probability.NewByCount()
	for _, marble := range []int{Red, Green, Red, Blue} {
		byCount.Insert(marble)
	}

	fmt.Printf("Probability of getting a red marble: %.2f\n", byCount.Probability(Red))
	fmt.Printf("Probability of getting a green marble: %.2f\n", byCount.Probability(Green))
	fmt.Printf("Probability of getting a blue marble: %.2f\n", byCount.Probability(Blue))
	// Output:
	// Probability of getting a red marble: 0.50
	// Probability of getting a green marble: 0.25
	// Probability of getting a blue marble: 0.25
}

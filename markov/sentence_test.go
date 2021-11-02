package markov

import (
	"testing"

	"github.com/spenserblack/probability"
)

// TestSentenceFeed checks that the proper values are created in the chain.
func TestSentenceFeed(t *testing.T) {
	chain := NewSentenceChain(2)

	chain.Feed([]string{"a", "b", "a", "c", "a", "b", "a", "d"})

	expectedCounts := make(sentenceCountMap)
	expectedCounts["a\000b"] = probability.NewByCount()
	expectedCounts["b\000a"] = probability.NewByCount()
	expectedCounts["a\000c"] = probability.NewByCount()
	expectedCounts["c\000a"] = probability.NewByCount()
	expectedCounts["a\000d"] = probability.NewByCount()
	expectedCounts["a\000b"].Insert("a")
	expectedCounts["a\000b"].Insert("a")
	expectedCounts["b\000a"].Insert("c")
	expectedCounts["b\000a"].Insert("d")
	expectedCounts["a\000c"].Insert("a")
	expectedCounts["c\000a"].Insert("b")
	expectedCounts["a\000d"].Insert("")

	t.Logf(`Sentence chain: %#v`, chain.chain)
	for prefix, suffixes := range chain.chain {
		t.Logf(`Suffix counts for prefix %q: %#v`, prefix, *suffixes)
	}

	if probability := chain.initialPrefixes.Probability("a\000b"); probability != 1.0 {
		t.Errorf(`Probability of sentence-starter %q = %f, want 1.0`, "a\000b", probability)
	}

	for expectedPrefix := range expectedCounts {
		if chainCounts, ok := chain.chain[expectedPrefix]; ok {
			// There are counts of suffixes in the actual value

			for _, countKey := range expectedCounts[expectedPrefix].Keys() {
				expectedSuffixCount := expectedCounts[expectedPrefix].Get(countKey)
				actualSuffixCount := chainCounts.Get(countKey)
				t.Logf(
					`Count for suffix %q for prefix %q = %d`,
					countKey,
					expectedPrefix,
					actualSuffixCount,
				)

				if actualSuffixCount != expectedSuffixCount {
					t.Errorf(
						`Count of suffix %q for prefix %q = %d, want %d`,
						countKey,
						expectedPrefix,
						actualSuffixCount,
						expectedSuffixCount,
					)
				}
			}
		} else {
			t.Errorf(`Want prefix %q`, expectedPrefix)
		}
	}
}

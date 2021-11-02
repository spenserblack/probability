package markov

import (
	"testing"

	"github.com/spenserblack/probability"
)

// TestWordFeed checks that the proper values are created in the chain.
func TestWordFeed(t *testing.T) {
	chain := NewWordChain(2)

	chain.Feed([]rune{'a', 'b', 'a', 'c', 'a', 'b', 'a', 'd'})

	expectedCounts := make(runeCountMap)
	expectedCounts["ab"] = probability.NewByCount()
	expectedCounts["ba"] = probability.NewByCount()
	expectedCounts["ac"] = probability.NewByCount()
	expectedCounts["ca"] = probability.NewByCount()
	expectedCounts["ad"] = probability.NewByCount()
	expectedCounts["ab"].Insert('a')
	expectedCounts["ba"].Insert('c')
	expectedCounts["ac"].Insert('a')
	expectedCounts["ca"].Insert('b')
	expectedCounts["ab"].Insert('a')
	expectedCounts["ad"].Insert(nullRune)

	t.Logf(`Word chain: %#v`, chain.chain)
	for prefix, suffixes := range chain.chain {
		t.Logf(`Suffix counts for prefix %q: %#v`, prefix, *suffixes)
	}

	if probability := chain.initialPrefixes.Probability("ab"); probability != 1.0 {
		t.Errorf(`Probability of sentence-starter "ab" = %f, want 1.0`, probability)
	}

	for expectedPrefix := range expectedCounts {
		if chainCounts, ok := chain.chain[expectedPrefix]; ok {
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
						`Count of suffix %q for prefix  %q = %d, want %d`,
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

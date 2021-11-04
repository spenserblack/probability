package markov

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/spenserblack/probability"
)

// WordJoiner is used to join words into a single string so that they can be
// used as a map key.
const wordJoiner string = "\000"

// SentenceCountMap maps prefixes to a count of possible suffixes.
type sentenceCountMap = map[string]*probability.ByCount

// SentenceChain is fed sentence(s) to create a chain of words in a sentence.
type SentenceChain struct {
	chain           sentenceCountMap
	initialPrefixes *probability.ByCount
	prefix          int
	r               *rand.Rand
}

// SentenceGenerator is created from a SentenceChain and generates a series of
// words (tokens).
type SentenceGenerator struct {
	chain  SentenceChain
	prefix []string
}

// NewSentenceChain creates a new sentence chain. It will still need to be fed.
// prefix is the number of tokens to be used as the prefix in the chain.
func NewSentenceChain(prefix int) *SentenceChain {
	return &SentenceChain{
		make(sentenceCountMap),
		probability.NewByCount(),
		prefix,
		rand.New(rand.NewSource(1)),
	}
}

// Feed takes a sentence as an array of strings, where each string is a word, or
// token. The first positional tokens are the sentence-starting prefix, and the
// last positional token ends the sentence.
func (c *SentenceChain) Feed(tokens []string) error {
	if len(tokens) < c.prefix {
		return fmt.Errorf("Got %d tokens, want >=%d", len(tokens), c.prefix)
	}

	prefix := make([]string, 0, c.prefix)
	prefix = append(prefix, tokens[:c.prefix]...)

	// NOTE Defines the sentence-starter
	c.initialPrefixes.Insert(strings.Join(prefix, wordJoiner))

	for _, suffix := range tokens[c.prefix:] {
		key := strings.Join(prefix, wordJoiner)
		c.Insert(key, suffix)
		oldPrefix := prefix[1:]
		prefix = make([]string, 0, c.prefix)
		prefix = append(prefix, oldPrefix...)
		prefix = append(prefix, suffix)
	}
	c.Insert(strings.Join(prefix, wordJoiner), "")
	return nil
}

// Seed will seed the random number generator used by the SentenceChain. Each
// SentenceChain chain has its own random number generator, but each
// SentenceGenerator created from a SentenceChain will share the same random
// number generator as the SentenceChain.
func (c *SentenceChain) Seed(seed int64) {
	c.r.Seed(seed)
}

// Insert inserts the prefix/suffix pair, and creates the prefix key if necessary.
//
// suffix should be an empty string if the prefix is intended to be able to end
// the sentence.
func (c *SentenceChain) Insert(prefix string, suffix string) (created bool) {
	byCount, exists := c.chain[prefix]
	if !exists {
		byCount = probability.NewByCount()
		byCount.SetRandom(c.r)
		c.chain[prefix] = byCount
		created = true
	}
	byCount.Insert(suffix)
	return
}

// MakeGenerator creates a generator from a chain for generating words (tokens)
// to build out a sentence.
func (c SentenceChain) MakeGenerator() *SentenceGenerator {
	prefix := strings.Split(c.initialPrefixes.Select().(string), wordJoiner)
	return &SentenceGenerator{c, prefix}
}

// Next returns the next token in the sequence. Stop determines if that is the
// last token.
func (g *SentenceGenerator) Next() (token string, stop bool) {
	token = g.prefix[0]
	nextWordSelector, ok := g.chain.chain[strings.Join(g.prefix, wordJoiner)]
	if token == "" {
		return token, true
	}

	var nextWord string
	if !ok {
		nextWord = ""
	} else {
		nextWord = nextWordSelector.Select().(string)
	}

	g.prefix = g.prefix[1:]
	g.prefix = append(g.prefix, nextWord)
	return
}

// HasNext returns true if another word can be returned.
func (g SentenceGenerator) HasNext() bool {
	return g.prefix[0] != ""
}

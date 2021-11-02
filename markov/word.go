package markov

import (
	"fmt"
	"math/rand"

	"github.com/spenserblack/probability"
)

// NullRune is used to detect when a chain should end.
const nullRune rune = '\000'

// RuneCountMap maps prefixes to a count of possible suffixes.
type runeCountMap = map[string]*probability.ByCount

// WordChain is fed word(s) to create a chain of characters in a word.
//
// It does not have to be a *literal* word. For example, "Hello, World!" can
// also be considered a "word". What's important is that the input and output
// is a chain of characters. This can help not only to generate a word, but to
// generate a whole sentence which may have newly generated words.
type WordChain struct {
	chain           runeCountMap
	initialPrefixes *probability.ByCount
	prefix          int
	r               *rand.Rand
}

// WordGenerator is created from a WordChain and generates a series of
// characters (tokens).
type WordGenerator struct {
	chain  WordChain
	prefix []rune
}

// NewWordChain creates a new word chain. It will still need to be fed.
// prefix is the number of tokens to be used as the prefix in the chain.
func NewWordChain(prefix int) *WordChain {
	return &WordChain{
		make(runeCountMap),
		probability.NewByCount(),
		prefix,
		rand.New(rand.NewSource(1)),
	}
}

// Feed takes a word as a slice of runes. The first positional tokens in the
// slice are the word-starting prefix, and the last positional token ends the
// word.
//
// With a small prefix and a lot of words feed into the chain, an extremely
// long "word" is possible, so this should be used with care.
func (c *WordChain) Feed(tokens []rune) error {
	if len(tokens) < c.prefix {
		return fmt.Errorf("Got %d tokens, want >=%d", len(tokens), c.prefix)
	}

	prefix := make([]rune, 0, c.prefix)
	prefix = append(prefix, tokens[:c.prefix]...)

	// NOTE Defines the word-starter
	c.initialPrefixes.Insert(string(prefix))

	for _, suffix := range tokens[c.prefix:] {
		key := string(prefix)
		c.Insert(key, suffix)
		oldPrefix := prefix[1:]
		prefix = make([]rune, 0, c.prefix)
		prefix = append(prefix, oldPrefix...)
		prefix = append(prefix, suffix)
	}
	c.Insert(string(prefix), nullRune)
	return nil
}

// Seed will seed the random number generator used by the WordChain. Each
// WordChain chain has its own random number generator, but each WordGenerator
// created from a WordChain will share the same random number generator as the
// WordChain.
func (c *WordChain) Seed(seed int64) {
	c.r.Seed(seed)
}

// Insert inserts the prefix/suffix pair, and creates the prefix key if
// necessary.
//
// Prefix should be the chain of runes, as a string.
//
// Suffix should be NullRune if the prefix is intended to be able to end the
// word.
func (c *WordChain) Insert(prefix string, suffix rune) (created bool) {
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

// MakeGenerator creates a geneartor from a chain for generating characters
// (tokens) to build out a sentence.
func (c WordChain) MakeGenerator() *WordGenerator {
	prefix := []rune(c.initialPrefixes.Select().(string))
	return &WordGenerator{c, prefix}
}

// Next returns the next token in the sequence. Stop determines if that is the
// last token.
func (g *WordGenerator) Next() (token rune, stop bool) {
	token = g.prefix[0]
	nextRuneSelector, ok := g.chain.chain[string(g.prefix)]
	if token == nullRune {
		return token, true
	}
	var nextRune rune
	if !ok {
		nextRune = nullRune
	} else {
		nextRune = nextRuneSelector.Select().(rune)
	}

	g.prefix = g.prefix[1:]
	g.prefix = append(g.prefix, nextRune)
	return
}

// HasNext returns true if another token can be returned.
func (g WordGenerator) HasNext() bool {
	return g.prefix[0] != nullRune
}

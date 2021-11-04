// Command-line interface for markov chains.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/spf13/cobra"

	"github.com/spenserblack/probability/markov"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// RootCmd is the main command.
var rootCmd = &cobra.Command{
	Use:   "markov",
	Short: "CLI tool to use Markov chains",
	Long:  `A CLI tool to feed a Markov chain and generate tokens.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.New("No valid subcommand or args provided")
	},
}

// SentenceCmd is used to manage a Markov generator for sentences.
var SentenceCmd = &cobra.Command{
	Use:   "sentence",
	Args:  cobra.MinimumNArgs(1),
	Short: "Manage a sentence chain/generator",
	Long: heredoc.Doc(`
		Manage a sentence chain/generator, where each token is a word.
		Each argument will be a "sentence" to feed the chain.
	`),
	Example: `sentence "my sentence" "my other sentence"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chain := markov.NewSentenceChain(PrefixFlag)
		chain.Seed(time.Now().UnixNano())
		for _, sentence := range args {
			words := strings.Fields(sentence)
			if err := chain.Feed(words); err != nil {
				return err
			}
		}
		generator := chain.MakeGenerator()
		for generator.HasNext() {
			word, _ := generator.Next()
			fmt.Printf("%s ", word)
		}
		fmt.Println()
		return nil
	},
}

// WordCmd is used to manage a Markov generator for words.
var WordCmd = &cobra.Command{
	Use:   "word",
	Args:  cobra.MinimumNArgs(1),
	Short: "Manage a word chain/generator",
	Long: heredoc.Doc(`
		Manage a word chain/generator, where each token is a letter.
		Each argument will be a "sentence" to feed the chain.

		The definition of "word" is used loosely in this context. Any string of
		characters is allowed, including spacing and punctuation. Because of this,
		this subcommand can be used to not only generate a random word, but be
		used as different, more random algorithm for generating a sentence.
	`),
	Example: `word "foo" "bar"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chain := markov.NewWordChain(PrefixFlag)
		chain.Seed(time.Now().UnixNano())
		for _, word := range args {
			if err := chain.Feed([]rune(word)); err != nil {
				return err
			}
		}
		generator := chain.MakeGenerator()
		for generator.HasNext() {
			letter, _ := generator.Next()
			fmt.Printf("%c", letter)
		}
		fmt.Println()
		return nil
	},
}

// PrefixFlag is the number of tokens to use as the prefix when generating
// tokens.
var PrefixFlag int

func init() {
	for _, subcommand := range []*cobra.Command{SentenceCmd, WordCmd} {
		rootCmd.AddCommand(subcommand)
		subcommand.Flags().IntVar(
			&PrefixFlag,
			"prefix",
			1,
			"The number of tokens to use as a prefix",
		)
	}
}

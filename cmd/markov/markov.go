// Command-line interface for markov chains.
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
	Long: `Manage a sentence chain/generator, where each token is a word.
		Each argument will be a "sentence" to feed the chain.`,
	Example: `sentence "my sentence" "my other sentence"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		chain := markov.NewSentenceChain(PrefixFlag)
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

// PrefixFlag is the number of tokens to use as the prefix when generating
// tokens.
var PrefixFlag int

func init() {
	rootCmd.AddCommand(SentenceCmd)
	SentenceCmd.Flags().IntVar(
		&PrefixFlag,
		"prefix",
		1,
		"The number of tokens to use as a prefix",
	)
}

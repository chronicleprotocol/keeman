package main

import (
	"os"

	"github.com/chronicleprotocol/keeman/pkg/cobra"
)

func main() {
	opts, cmd := cobra.Command()
	cmd.PersistentFlags().StringVarP(
		&opts.InputFile,
		"input",
		"i",
		"",
		"input file path",
	)
	cmd.PersistentFlags().BoolVarP(
		&opts.Verbose,
		"verbose",
		"v",
		false,
		"verbose logging",
	)
	cmd.AddCommand(
		cobra.NewDerive(opts),
		cobra.GenerateSeed(opts),
		cobra.NewList(opts),
	)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

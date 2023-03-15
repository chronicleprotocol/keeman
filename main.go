package main

import (
	"os"

	"github.com/chronicleprotocol/keeman/cobra"
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
	// cmd.PersistentFlags().StringVarP(
	// 	&opts.OutputFile,
	// 	"output",
	// 	"o",
	// 	"",
	// 	"output file path",
	// )
	cmd.PersistentFlags().BoolVarP(
		&opts.Verbose,
		"verbose",
		"v",
		false,
		"verbose logging",
	)
	cmd.AddCommand(
		cobra.NewDerive(opts),
		cobra.NewDeriveTf(),
		cobra.GenerateSeed(opts),
		cobra.NewList(opts),
	)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
	os.Exit(0)
}

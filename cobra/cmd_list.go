package cobra

import (
	"github.com/spf13/cobra"
)

func NewList(opts *Options) *cobra.Command {
	var all bool
	var index int
	cmd := &cobra.Command{
		Use:     "list [--all]",
		Aliases: []string{"l"},
		Short:   "List word count and first word from the input, omitting the comments.",
		RunE: func(_ *cobra.Command, args []string) error {
			if all {
				lines, err := linesFromFile(opts.InputFile)
				if err != nil {
					return err
				}
				for _, l := range lines {
					printLine(l)
				}
				return nil
			}
			l, err := lineFromFile(opts.InputFile, index)
			if err != nil {
				return err
			}
			printLine(l)
			return nil
		},
	}
	cmd.Flags().IntVar(
		&index,
		"index",
		0,
		"data index",
	)
	cmd.Flags().BoolVarP(
		&all,
		"all",
		"a",
		false,
		"all data",
	)
	return cmd
}

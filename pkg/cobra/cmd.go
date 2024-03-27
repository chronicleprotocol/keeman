package cobra

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chronicleprotocol/keeman/pkg/txt"
)

type Options struct {
	InputFile string
	Verbose   bool
}

func Command() (*Options, *cobra.Command) {
	return &Options{}, &cobra.Command{
		Use: "keeman",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}
}

func lineFromFile(filename string, idx int) (string, error) {
	lines, err := linesFromFile(filename)
	if err != nil {
		return "", err
	}
	return selectLine(lines, idx)
}

func selectLine(lines []string, lineIdx int) (string, error) {
	if len(lines) <= lineIdx {
		return "", fmt.Errorf("data needs %d line(s)", lineIdx+1)
	}
	return lines[lineIdx], nil
}

func linesFromFile(filename string) ([]string, error) {
	file, fileClose, err := inputFileOrStdin(filename)
	if err != nil {
		return nil, err
	}
	defer func() { err = fileClose() }()
	return txt.ReadNonEmptyLines(file, 0, false)
}

func inputFileOrStdin(inputFilePath string) (*os.File, func() error, error) {
	if inputFilePath != "" {
		f, err := os.Open(inputFilePath)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to open file: %w", err)
		}
		return f, f.Close, nil
	}
	stdin, err := NonEmptyStdIn()
	return stdin, func() error { return nil }, err
}

func NonEmptyStdIn() (*os.File, error) {
	if fi, err := os.Stdin.Stat(); err != nil {
		return nil, fmt.Errorf("unable to stat stdin: %w", err)
	} else if fi.Size() <= 0 && fi.Mode()&os.ModeNamedPipe == 0 {
		return nil, errors.New("stdin is empty")
	}
	return os.Stdin, nil
}

package cobra

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/chronicleprotocol/keeman/txt"
)

type Options struct {
	InputFile  string
	OutputFile string
	Verbose    bool
}

func Command() (*Options, *cobra.Command) {
	return &Options{}, &cobra.Command{
		Use: "keeman",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
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

func linesFromFile(filename string) ([]string, error) {
	file, fileClose, err := inputFileOrStdin(filename)
	if err != nil {
		return nil, err
	}
	defer func() { err = fileClose() }()
	return txt.ReadNonEmptyLines(file, 0, false)
}

func selectLine(lines []string, lineIdx int) (string, error) {
	if len(lines) <= lineIdx {
		return "", fmt.Errorf("data needs %d line(s)", lineIdx+1)
	}
	return lines[lineIdx], nil
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
func printLine(l string) {
	split := strings.Split(l, " ")
	fmt.Println(len(split), split[0])
}

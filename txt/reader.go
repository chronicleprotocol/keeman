package txt

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

func ReadNonEmptyLines(r io.Reader, limit int, withComments bool) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if !withComments {
			text = stripFromFirstChar(text, "#;")
		}

		text = strings.Trim(text, "\t \n")
		if len(text) > 0 {
			lines = append(lines, text)
		}

		if 0 < limit && limit <= len(lines) {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func stripFromFirstChar(s, chars string) string {
	if cut := strings.IndexAny(s, chars); cut >= 0 {
		return strings.TrimRightFunc(s[:cut], unicode.IsSpace)
	}
	return s
}

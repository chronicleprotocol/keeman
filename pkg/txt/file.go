package txt

import (
	"fmt"
	"log"
	"os"
)

func ReadLines(filename string, limit int, withComment bool) ([]string, error) {
	file, fileClose, err := FileOpen(filename)
	if err != nil {
		return nil, err
	}
	defer fileClose()
	return FileReadLines(file, limit, withComment)
}

func FileReadLines(file *os.File, limit int, withComment bool) ([]string, error) {
	if fileIsEmpty(file) {
		return nil, fmt.Errorf("file is empty: %s", file.Name())
	}
	return ReadNonEmptyLines(file, limit, withComment)
}

func FileOpen(filename string) (*os.File, func(), error) {
	file, err := os.Open(filename)
	return file, func() {
		if err := file.Close(); err != nil {
			log.Fatal(file.Name(), err)
		}
	}, err
}

func fileIsEmpty(file *os.File) bool {
	info, err := file.Stat()
	return err != nil || info.Size() == 0 && info.Mode()&os.ModeNamedPipe == 0
}

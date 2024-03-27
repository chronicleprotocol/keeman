package wl

import (
	_ "embed"
	"fmt"
	"hash/crc32"
	"strings"
)

func init() {
	// Ensure word list is correct
	// $ wget -o hdwallet/wl/polish.txt https://raw.githubusercontent.com/p2w34/blog/master/2019-03-25-polish-word-list-bip0039/polish.txt
	// $ crc32 hdwallet/wl/polish.txt
	// d574dc8c
	checksum := crc32.ChecksumIEEE([]byte(polish))
	if fmt.Sprintf("%x", checksum) != "d574dc8c" {
		panic("polish checksum invalid")
	}
}

// Polish is a slice of mnemonic words taken from:
// https://github.com/p2w34/blog/blob/master/2019-03-25-polish-word-list-bip0039/polish-word-list-bip0039.md
var Polish = strings.Split(strings.TrimSpace(polish), "\n")

//go:embed polish.txt
var polish string

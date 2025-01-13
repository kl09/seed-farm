package mnemonic

import (
	_ "embed"
	"fmt"
	"hash/crc32"
	"strconv"
	"strings"
)

//go:embed vocabulary.txt
var words string

// Words a slice of mnemonic words taken from the bip39 specification
// https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt
var (
	Words = strings.Split(strings.TrimSpace(words), "\n")
	// WordsShortly is a map of mnemonic words to their index in the Words slice.
	WordsShortly = map[string]int{}
)

func init() {
	var expected uint32 = 0xc1dbd296
	// Ensure word list is correct
	// $ wget https://raw.githubusercontent.com/bitcoin/bips/master/bip-0039/english.txt
	// $ crc32 english.txt
	// OUTPUT: c1dbd296
	checksum := crc32.ChecksumIEEE([]byte(words))
	if checksum != expected {
		panic(fmt.Errorf("wordlist checksum mismatch: expected %x, got %x", expected, checksum))
	}

	for i, w := range Words {
		WordsShortly[w] = i
	}
}

// FromMnemonicToShort converts a mnemonic to a short representation.
func FromMnemonicToShort(mnemonic string) string {
	var short string
	for _, s := range strings.Split(mnemonic, " ") {
		i, ok := WordsShortly[s]
		if ok {
			short += strconv.Itoa(i) + "."
		} else {
			short += s + "."
		}
	}

	return short
}

// FromShortToMnemonic converts a short representation to a mnemonic.
func FromShortToMnemonic(short string) (string, error) {
	var mnemonic string
	for _, sh := range strings.Split(short, ".") {
		if sh == "" {
			continue
		}

		i, err := strconv.Atoi(sh)
		if err != nil {
			return "", fmt.Errorf("from short to mnemonic: %w", err)
		}

		for word, v := range WordsShortly {
			if v == i {
				mnemonic += word + " "
			}
		}
	}

	return strings.TrimRight(mnemonic, " "), nil
}

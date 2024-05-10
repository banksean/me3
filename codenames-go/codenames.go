package main

import (
	"fmt"
	"strings"

	_ "embed"
)

var (
	//go:embed wordlist.txt
	wordlistFile string
)

func main() {
	fmt.Printf("wordlist:\n%v\n", strings.Split(wordlistFile, "\n"))
}

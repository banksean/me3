package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/banksean/me3/blaim"
)

func processLogFile(in io.Reader) error {
	b, err := io.ReadAll(in)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		parsed, err := blaim.ParseLogLine(line)
		if err != nil {
			return fmt.Errorf("error: %v", err)
		}
		fmt.Printf("parsed: %#v\n", parsed)
	}
	return nil
}

func main() {
	processLogFile(os.Stdin)
}

package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/banksean/me3/blaim"
	"github.com/sourcegraph/go-diff/diff"
)

const acceptLogEnvVar = "ACCEPT_LOG"

func processLogFile(in io.Reader) (map[string][]*blaim.AcceptLogLine, error) {
	ret := map[string][]*blaim.AcceptLogLine{}

	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		parsed, err := blaim.ParseLogLine(line)
		if err != nil {
			return nil, fmt.Errorf("error: %v", err)
		}
		if parsed == nil {
			continue
		}
		if _, ok := ret[parsed.FileName]; !ok {
			ret[parsed.FileName] = []*blaim.AcceptLogLine{}
		}
		ret[parsed.FileName] = append(ret[parsed.FileName], parsed)
	}
	return ret, nil
}

func getAdditions(body string) (string, int) {
	ret := []string{}
	start := -1
	for i, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "+") {
			if len(ret) == 0 {
				start = i
			}
			ret = append(ret, line[1:])
		}
	}
	return strings.Join(ret, "\n"), start
}

func main() {
	diffReader := diff.NewMultiFileDiffReader(os.Stdin)
	acceptLogPath := os.Getenv(acceptLogEnvVar)
	if acceptLogPath == "" {
		fmt.Printf("%s environment variable is not set", acceptLogEnvVar)
		os.Exit(1)
	}
	logFile, err := os.Open(acceptLogPath)
	if err != nil {
		fmt.Printf("error opening accept log at %s: %v\n", acceptLogPath, err)
		os.Exit(1)
	}
	acceptsForFile, err := processLogFile(logFile)
	if err != nil {
		fmt.Printf("error processing accept log: %v\n", err)
		os.Exit(1)
	}
	for i := 0; ; i++ {
		fdiff, err := diffReader.ReadFile()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("err reading diff: %s", err)
		}
		origName := fdiff.OrigName[2:]
		newName := fdiff.NewName[2:]
		accepts := acceptsForFile[origName]
		if newName != origName {
			accepts = append(accepts, acceptsForFile[newName]...)
		}
		for _, hunk := range fdiff.Hunks {
			body, addsStart := getAdditions(string(hunk.Body))
			for _, accept := range accepts {
				idx := strings.Index(body, accept.Text)
				if idx != -1 {
					fmt.Printf("found a matching accept log entry for %q starting at position %d on line %d of %s:\n%v\n",
						accept.Text, idx, int32(addsStart)+hunk.NewStartLine, newName, body)
				}
			}
		}
	}
}

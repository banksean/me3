package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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

// Diff hunks contain both additions and deletions, but we only
// care about the additions here. Returns the text of just the added
// lines, if any, and the offset for the line within the hunk where
// the additions start.
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

type BlaimLine struct {
	Filename        string                `json:"filename"`
	Position        blaim.Position        `json:"position"`
	Text            string                `json:"text"`
	InferenceConfig blaim.InferenceConfig `json:"inference_config"`
}

// Parses the accept logs, compares their contents to the current git diff results
// and produces a json-formatted array of BlaimLine objects, one for each git diff hunk
// that contains text that appears in the accept logs.
func generateBlaimFile() {
	diffReader := diff.NewMultiFileDiffReader(os.Stdin)
	acceptLogPath := os.Getenv(acceptLogEnvVar)
	if acceptLogPath == "" {
		fmt.Printf("%s environment variable is not set\n", acceptLogEnvVar)
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

	// Read the git diff output and check for blaim entries for each file mentioned
	// in the diff.
	for {
		fdiff, err := diffReader.ReadFile()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("err reading diff: %s", err)
		}
		// Strip the "a/" and "b/" prefixes from the diff file names.
		origName := fdiff.OrigName[2:]
		newName := fdiff.NewName[2:]
		accepts := acceptsForFile[origName]
		// If the filename changed in this diff, group the accept logs for the
		// old name and the new name together.
		if newName != origName {
			accepts = append(accepts, acceptsForFile[newName]...)
		}
		blaimLines := []BlaimLine{}

		// Now check each "hunk" in the diff'd file to see if there are any
		// entries in the .blaim file about it.
		for _, hunk := range fdiff.Hunks {
			addedInDiffHunk, addsStart := getAdditions(string(hunk.Body))
			for _, accept := range accepts {
				idx := strings.Index(addedInDiffHunk, accept.Text)
				if idx < 0 {
					continue
				}
				lineOffset := len(strings.Split(addedInDiffHunk[:idx], "\n"))

				blaimLine := BlaimLine{
					Filename: accept.FileName,
					Position: blaim.Position{
						Line:      lineOffset + addsStart + int(hunk.NewStartLine),
						Character: idx,
					},
					Text:            accept.Text,
					InferenceConfig: accept.InferenceConfig,
				}
				blaimLines = append(blaimLines, blaimLine)
			}
		}
		if len(blaimLines) == 0 {
			continue
		}
		jsonBytes, err := json.MarshalIndent(blaimLines, "", "\t")
		if err != nil {
			fmt.Printf("error marshaling blaimLines: %v", err)
			os.Exit(1)
		}
		if true {
			fmt.Printf("%s\n", string(jsonBytes))

		}
	}
}

// Represents the set of blaim lines and ranges for a particular source file.
type BlaimRangeSet struct {
	blaimLines []*BlaimLine
}

func (s *BlaimRangeSet) ForSourceLine(lineNumber int) []*BlaimLine {
	ret := []*BlaimLine{}
	for _, blaimLine := range s.blaimLines {
		textLines := strings.Split(blaimLine.Text, "\n")
		if lineNumber >= blaimLine.Position.Line &&
			lineNumber <= blaimLine.Position.Line+len(textLines) {
			ret = append(ret, blaimLine)
		}
	}
	return ret
}

// parses a json-formatted list of BlaimLine objects from stdin,
// and produces a line-by-line annotation of AI-generated code for
// each file mentioned in the BlaimLine input list.
func annotate() error {
	reader := os.Stdin
	um := json.NewDecoder(reader)
	blaimLines := []*BlaimLine{}
	for {
		err := um.Decode(&blaimLines)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error decoding BlaimLines: %v", err)
			os.Exit(1)
		}
	}
	blaimLinesByFile := map[string][]*BlaimLine{}

	// Group the blaim lines by the source file path they refer to.
	for _, blaimLine := range blaimLines {
		if _, ok := blaimLinesByFile[blaimLine.Filename]; !ok {
			blaimLinesByFile[blaimLine.Filename] = []*BlaimLine{}
		}
		blaimLinesByFile[blaimLine.Filename] = append(blaimLinesByFile[blaimLine.Filename], blaimLine)
	}

	for fileName, fileBlaimLines := range blaimLinesByFile {
		blaimRangeSet := &BlaimRangeSet{
			blaimLines: fileBlaimLines,
		}
		fileBytes, err := os.ReadFile(filepath.Join(baseDir, fileName))
		if err != nil {
			return err
		}

		fileLines := strings.Split(string(fileBytes), "\n")

		// TODO: this doesn't handle multi-line code suggestions well.
		// For instance, if an accepted suggestion spans multiple lines,
		// (conains \n characters) then this will only annotate the *first*
		// line containing the generated code suggestion.
		for lineNumber, lineText := range fileLines {
			linePrefix := "  "
			blaimLineMatches := blaimRangeSet.ForSourceLine(lineNumber)
			if len(blaimLineMatches) > 0 {
				linePrefix = "* "
			}
			fmt.Printf("%s%s\n", linePrefix, lineText)
		}
	}
	return nil
}

var (
	baseDir = "/Users/seanmccullough/code/me3"
)

func main() {
	cmd := "generate"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}
	if cmd == "generate" {
		generateBlaimFile()
		return
	}
	if cmd == "annotate" {
		if err := annotate(); err != nil {
			fmt.Printf("annotate error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	fmt.Printf("unrecognized command: %q\n", os.Args[1])
	os.Exit(1)
}

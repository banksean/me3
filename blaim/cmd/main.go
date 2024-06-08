// package main is the main cli entry point for blaim commands
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
	"github.com/urfave/cli/v2"
)

// processAcceptedSuggestionsLog parses the contents of a "accepted.suggestions.log" file
// which the VS Code extension has been writing entries to as the user has edited
// code and accepted AI-generated suggestions.
func processAcceptedSuggestionsLog(in io.Reader) (map[string][]*blaim.AcceptLogLine, error) {
	ret := map[string][]*blaim.AcceptLogLine{}

	b, err := io.ReadAll(in)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		parsed, err := blaim.ParseAcceptLogLine(line)
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

func searchMatchingRanges(diff, aiInsert string) []blaim.Range {
	idx := strings.Index(diff, aiInsert)
	if idx == -1 {
		return nil
	}
	ret := []blaim.Range{}
	return ret
}

// Parses the accept logs, compares their contents to the current git diff results
// and produces a json-formatted array of BlaimLine objects, one for each git diff hunk
// that contains text that appears in the accept logs.
func generate(diffStream, logReader io.Reader, out io.Writer) error {
	diffReader := diff.NewMultiFileDiffReader(diffStream)

	acceptsForFile, err := processAcceptedSuggestionsLog(logReader)
	if err != nil {
		return fmt.Errorf("error processing accept log: %v", err)
	}

	// Read the git diff output and check for blaim entries for each file mentioned
	// in the diff.
	for {
		fdiff, err := diffReader.ReadFile()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("err reading diff: %s", err)
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
		blaimLines := []blaim.BlaimLine{}

		// Now check each "hunk" in the diff'd file to see if there are any
		// entries in the .blaim file about it.
		for _, hunk := range fdiff.Hunks {
			addedInDiffHunk, addsStart := getAdditions(string(hunk.Body))
			for _, accept := range accepts {
				// Check for exact matches:
				idx := strings.Index(addedInDiffHunk, accept.Text)
				lineOffset := -1
				if idx > 0 {
					lineOffset = len(strings.Split(addedInDiffHunk[:idx], "\n"))
				} else {
					// TODO: Implement the suffix/prefix heuristic used by the vscocde-extension
				}
				if lineOffset == -1 {
					continue
				}
				blaimLine := blaim.BlaimLine{
					FileName: accept.FileName,
					Range: blaim.Range{
						Start: blaim.Position{
							Line:      lineOffset + addsStart + int(hunk.NewStartLine),
							Character: idx,
						},
						// TODO: figure out the End position.
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
		m := json.NewEncoder(out)
		m.SetIndent("", "\t")
		err = m.Encode(blaimLines)
		if err != nil {
			return fmt.Errorf("error marshaling blaimLines: %v", err)
		}
	}
	return nil
}

// Represents the set of blaim lines and ranges for a particular source file.
type BlaimRangeSet struct {
	blaimLines []*blaim.BlaimLine
}

// ForSourceLine returns the BlaimLines that cover that line of the source file.
func (s *BlaimRangeSet) ForSourceLine(lineNumber int) []*blaim.BlaimLine {
	ret := []*blaim.BlaimLine{}
	// This problem is straight out of a coding interview question:
	//   "Given a list of ranges (start, end), write a function that returns true
	//   if a particular value falls within any of the ranges."
	// And reader, this O(number of ranges) solution is not what the interviewer
	// wants to see, but it's sufficient for this PoC:
	for _, blaimLine := range s.blaimLines {
		textLines := strings.Split(blaimLine.Text, "\n")
		if lineNumber >= blaimLine.Range.Start.Line &&
			lineNumber <= blaimLine.Range.Start.Line+len(textLines) {
			ret = append(ret, blaimLine)
		}
	}
	return ret
}

func formatAnnotationLinePrefix(line *blaim.BlaimLine) string {
	return fmt.Sprintf("[%s, temp: %.1f] ", line.InferenceConfig.ModelName, line.InferenceConfig.Temperature)
}

// parses a json-formatted list of BlaimLine objects from stdin,
// and produces a line-by-line annotation of AI-generated code for
// each file mentioned in the BlaimLine input list.
func annotate() error {
	reader := os.Stdin
	um := json.NewDecoder(reader)
	blaimLines := []*blaim.BlaimLine{}
	for {
		err := um.Decode(&blaimLines)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error decoding BlaimLines: %v", err)
			os.Exit(1)
		}
	}
	blaimLinesByFile := map[string][]*blaim.BlaimLine{}

	longestLinePrefixLen := -1

	// Group the blaim lines by the source file path they refer to.
	for _, blaimLine := range blaimLines {
		if _, ok := blaimLinesByFile[blaimLine.FileName]; !ok {
			blaimLinesByFile[blaimLine.FileName] = []*blaim.BlaimLine{}
		}
		blaimLinesByFile[blaimLine.FileName] = append(blaimLinesByFile[blaimLine.FileName], blaimLine)
		linePrefix := formatAnnotationLinePrefix(blaimLine)
		if len(linePrefix) > longestLinePrefixLen {
			longestLinePrefixLen = len(linePrefix)
		}
	}

	defaultLinePrefix := strings.Repeat(" ", longestLinePrefixLen)

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
			linePrefix := defaultLinePrefix
			blaimLineMatches := blaimRangeSet.ForSourceLine(lineNumber)
			if len(blaimLineMatches) > 0 {
				// Just use the first blaim entry, if there is more than one for
				// this line of the diff.
				linePrefix = formatAnnotationLinePrefix(blaimLineMatches[0])
			}
			fmt.Printf("%s%s\n", linePrefix, lineText)
		}
	}
	return nil
}

var (
	baseDir                    string
	acceptedSuggestionsLogPath string
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "root",
				Value:       ".",
				Usage:       "path to the root of the git checkout",
				Destination: &baseDir,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "generate",
				Aliases: []string{"g"},
				Usage:   "generate a .blaim file from git diff output at stdin, and the contents of accepted.suggestions.log",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "accept-log",
						Value:       "",
						Usage:       "path to the accepted.suggestions.log file",
						Destination: &acceptedSuggestionsLogPath,
					},
				},
				Action: func(cCtx *cli.Context) error {
					logFile, err := os.Open(acceptedSuggestionsLogPath)
					if err != nil {
						return fmt.Errorf("error opening accept log at %s: %v", acceptedSuggestionsLogPath, err)
					}

					return generate(os.Stdin, logFile, os.Stdout)
				},
			},
			{
				Name:    "annotate",
				Aliases: []string{"a"},
				Usage:   "produce a line-by-line annotation of source files that contain machine-generated code changes",
				Action: func(cCtx *cli.Context) error {
					return annotate()
				},
			},
		},
		Name:  "blaim",
		Usage: "manage the attributrion of machine-generated code changes",
		Action: func(*cli.Context) error {
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

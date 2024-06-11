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

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/sourcegraph/go-diff/diff"
	"github.com/urfave/cli/v2"
	"gopkg.in/vmarkovtsev/go-lcss.v1"
)

const (
	minEditDistanceSimilarity = 0.8
	minLCS                    = 20
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
func getAdditions(body string) string {
	ret := []string{}
	for _, line := range strings.Split(body, "\n") {
		if strings.HasPrefix(line, "+") {
			ret = append(ret, line[1:])
		}
	}
	return strings.Join(ret, "\n")
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
			addedInDiffHunk := getAdditions(string(hunk.Body))
			// Now find any acceptLog entriesd that match the added text.
			matchingBlaimLines := getMatchingAcceptLogsForHunk(accepts, addedInDiffHunk)
			// Offset the line numbers in the blaim entries by the start line of the hunk
			// so they line up with the full file contents.
			for _, match := range matchingBlaimLines {
				match.Range.Start.Line += int(hunk.NewStartLine) + 1
				match.Range.End.Line += int(hunk.NewStartLine) + 1
				blaimLines = append(blaimLines, match)
			}
		}
		if len(blaimLines) == 0 {
			continue
		}
		m := json.NewEncoder(out)
		m.SetIndent("", "  ")
		err = m.Encode(blaimLines)
		if err != nil {
			return fmt.Errorf("error marshaling blaimLines: %v", err)
		}
	}
	return nil
}

func indexToPos(s string, i int) blaim.Position {
	prefixLines := strings.Split(s[:i], "\n")
	startLine := len(prefixLines) + 1
	startChar := i - len(strings.Join(prefixLines[:len(prefixLines)-1], "\n"))
	return blaim.Position{
		Line:      startLine,
		Character: startChar,
	}
}

// Things to watch out for:
// - The accept log text may not exactly match the diff text, so we need to do some fuzzy matching.
// - The accept log text may span multiple lines, so we need to handle that.
// - The line numbers in the accept log may not match the line numbers in the diff
// - The user may have accepted suggstions in a different order than they appear in the diff
func getMatchingAcceptLogsForHunk(accepts []*blaim.AcceptLogLine, addedInDiffHunk string) []blaim.BlaimLine {
	blaimLines := []blaim.BlaimLine{}
	for _, accept := range accepts {
		targetString := accept.Text
		startIdx := strings.Index(addedInDiffHunk, targetString)
		endIdx := -1
		if startIdx >= 0 { // exact match
			// the accepted text starts at lineOffset within the diff hunk, so count the newlines preceding the accepted text
			endIdx = startIdx + len(targetString)
		} else { // check for a fuzzy match
			// This edit distance check skips too many true positives, so we're not using it for now.
			if false {
				ld := fuzzy.LevenshteinDistance(addedInDiffHunk, targetString)
				similarity := float32(len(addedInDiffHunk)-ld) / float32(len(addedInDiffHunk))

				fmt.Printf("similarity: %f\n", similarity)
				if similarity < minEditDistanceSimilarity {
					continue
				}
			}

			// Find the longest common substring between the diff hunk and the accept log text
			targetString = string(lcss.LongestCommonSubstring([]byte(addedInDiffHunk), []byte(targetString)))
			if len(targetString) < minLCS {
				continue
			}
			startIdx = strings.Index(addedInDiffHunk, targetString)
			endIdx = startIdx + len(targetString)
		}
		if endIdx == -1 { // no match
			continue
		}

		startPos, endPos := indexToPos(addedInDiffHunk, startIdx), indexToPos(addedInDiffHunk, endIdx)
		blaimLine := blaim.BlaimLine{
			FileName: accept.FileName,
			Range: blaim.Range{
				Start: startPos,
				End:   endPos,
			},
			Text:            accept.Text,
			InferenceConfig: accept.InferenceConfig,
		}
		blaimLines = append(blaimLines, blaimLine)
	}
	return blaimLines
}

// Represents the set of blaim lines and ranges for a particular source file.
type BlaimRangeSet struct {
	blaimLines []*blaim.BlaimLine
}

// ForSourceLine returns the BlaimLines that cover that line of the source file.
func (s *BlaimRangeSet) ForSourceLine(lineNumber int) []*blaim.BlaimLine {
	ret := []*blaim.BlaimLine{}
	for _, blaimLine := range s.blaimLines {
		if lineNumber >= blaimLine.Range.Start.Line &&
			lineNumber <= blaimLine.Range.End.Line {
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
func annotate(blaimReader io.Reader, out io.Writer) error {
	// Group the blaim lines by the source file path they refer to.
	blaimLinesByFile, err := readBlaimFile(blaimReader)
	if err != nil {
		return err
	}

	// Read the contents of each file in the diff
	for fileName, fileBlaimLines := range blaimLinesByFile {
		blaimRangeSet := &BlaimRangeSet{
			blaimLines: fileBlaimLines,
		}
		fileBytes, err := os.ReadFile(filepath.Join(baseDir, fileName))
		if err != nil {
			return err
		}

		annotateLines(fileBytes, blaimRangeSet, out)
	}
	return nil
}

func annotateLines(fileBytes []byte, blaimRangeSet *BlaimRangeSet, out io.Writer) {
	fileLines := strings.Split(string(fileBytes), "\n")
	prefixLines := []string{}

	longestLinePrefixLen := 0

	for lineNumber := range fileLines {
		blaimLineMatches := blaimRangeSet.ForSourceLine(lineNumber + 1)
		if len(blaimLineMatches) > 0 {
			linePrefix := formatAnnotationLinePrefix(blaimLineMatches[0])
			prefixLines = append(prefixLines, linePrefix)
			if len(linePrefix) > longestLinePrefixLen {
				longestLinePrefixLen = len(linePrefix)
			}
		} else {
			prefixLines = append(prefixLines, "")
		}
	}

	defaultPrefix := strings.Repeat(" ", longestLinePrefixLen)
	for lineNumber, lineText := range fileLines {
		linePrefix := prefixLines[lineNumber]
		if linePrefix == "" {
			linePrefix = defaultPrefix
		}
		fmt.Fprintf(out, "%s%s\n", linePrefix, lineText)
	}
}

func readBlaimFile(blaimReader io.Reader) (map[string][]*blaim.BlaimLine, error) {
	um := json.NewDecoder(blaimReader)
	blaimLines := []*blaim.BlaimLine{}
	for {
		err := um.Decode(&blaimLines)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error decoding BlaimLines: %v", err)
			return nil, err
		}
	}
	blaimLinesByFile := map[string][]*blaim.BlaimLine{}

	for _, blaimLine := range blaimLines {
		if _, ok := blaimLinesByFile[blaimLine.FileName]; !ok {
			blaimLinesByFile[blaimLine.FileName] = []*blaim.BlaimLine{}
		}
		blaimLinesByFile[blaimLine.FileName] = append(blaimLinesByFile[blaimLine.FileName], blaimLine)
	}
	return blaimLinesByFile, nil
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
					return annotate(os.Stdin, os.Stdout)
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

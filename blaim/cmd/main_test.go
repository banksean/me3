package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sourcegraph/go-diff/diff"

	_ "embed"
)

var (
	//go:embed testdata/diff.txt
	diffText string

	//go:embed testdata/accepted.suggestions.log
	acceptedSuggestionsLogText string

	//go:embed testdata/expected_blaim.json
	expectedBlaimText string

	//go:embed testdata/playground.js
	playgroundJS string

	//go:embed testdata/expected_annotate.txt
	expectedAnnotateText string
)

// TestGenerate needs to read in some example data:
// A git diff output stream
// An accepted.suggestions.log stream
// and make sure that generate produces the correct condensted blaim list.
func TestGenerate(t *testing.T) {

	out := &bytes.Buffer{}
	generate(strings.NewReader(diffText), strings.NewReader(acceptedSuggestionsLogText), out)
	got := out.String()
	diff := cmp.Diff(expectedBlaimText, got)
	if diff != "" {
		fmt.Printf("expected: %s\n", expectedBlaimText)
		fmt.Printf("got: %s\n", got)
		t.Errorf("diff: %s", diff)
	}
}

func TestProcessAcceptedSuggestionsLog(t *testing.T) {
	acceptsForFile, err := processAcceptedSuggestionsLog(strings.NewReader(acceptedSuggestionsLogText))
	if err != nil {
		t.Errorf("error processing accept log: %v", err)
	}
	if len(acceptsForFile["playground.js"]) != 3 {
		t.Errorf("expected 3 files, got %d", len(acceptsForFile["playground.js"]))
	}
}

func TestGetAdditions(t *testing.T) {
	diffReader := diff.NewMultiFileDiffReader(strings.NewReader(diffText))
	fdiff, err := diffReader.ReadFile()
	if err != nil {
		t.Errorf("error reading diff: %v", err)
	}
	if len(fdiff.Hunks) != 1 {
		t.Errorf("expected 1 hunk, got %d", len(fdiff.Hunks))
	}
	for _, hunk := range fdiff.Hunks {
		additions := getAdditions(string(hunk.Body))
		if len(additions) == 0 {
			t.Errorf("expected some additions, got none")
		}
	}
}

func TestFoo(t *testing.T) {

}

func TestReadBlaimFile(t *testing.T) {
	blaimLinesByFile, err := readBlaimFile(strings.NewReader(expectedBlaimText))
	if err != nil {
		t.Errorf("error reading blaim file: %v", err)
	}
	if len(blaimLinesByFile) != 1 {
		t.Errorf("expected 1 file1, got %d", len(blaimLinesByFile))
	}
	if len(blaimLinesByFile["playground.js"]) != 3 {
		t.Errorf("expected 3 blaim lines for playground.js, got %d", len(blaimLinesByFile["playground.js"]))
	}
}

func TestAnnotateLines(t *testing.T) {
	blaimLinesByFile, err := readBlaimFile(strings.NewReader(expectedBlaimText))
	if err != nil {
		t.Errorf("error reading blaim file: %v", err)
	}
	blaimRangeSet := BlaimRangeSet{
		blaimLinesByFile["playground.js"],
	}
	lineNumber := 4
	blaimLineMatches := blaimRangeSet.ForSourceLine(lineNumber)
	if len(blaimLineMatches) != 1 {
		t.Errorf("expected 1 blaim line match, got %d", len(blaimLineMatches))
		return
	}
	blaimLine := blaimLineMatches[0]
	expected := "[codegemma, temp: 0.2] "
	if formatAnnotationLinePrefix(blaimLine) != expected {
		t.Errorf("expected %s, got %s", expected, formatAnnotationLinePrefix(blaimLine))
	}

	out := &bytes.Buffer{}
	annotateLines([]byte(playgroundJS), &blaimRangeSet, out)
	got := out.String()
	diff := cmp.Diff(expectedAnnotateText, got)
	if diff != "" {
		t.Errorf("diff: %s", diff)
	}
}

package main

import (
	"bytes"
	"strings"
	"testing"

	_ "embed"
)

var (
	//go:embed testdata/diff.txt
	diffText string

	//go:embed testdata/accepted.suggestions.log
	acceptedSuggestionsLogText string

	//go:embed testdata/expected_blaim.json
	expectedBlaimText string
)

// TestGenerate needs to read in some example data:
// A git diff output stream
// An accepted.suggestions.log stream
// and make sure that generate produces the correct condensted blaim list.
func TestGenerate(t *testing.T) {

	out := &bytes.Buffer{}
	generate(strings.NewReader(diffText), strings.NewReader(acceptedSuggestionsLogText), out)
	got := out.String()
	if got != expectedBlaimText {
		t.Errorf("got:\n%s\nexpected:\n%s\n", got, expectedBlaimText)
	}
}

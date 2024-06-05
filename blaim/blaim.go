package blaim

import (
	"encoding/json"
	"strings"
)

// BlaimLine represents an entry in a .blaim file.
// Each entry describes a single range of text in
// a named source code file.  The description of the
// text range includes information about how the text
// range was generated (e.g. name of the model, inference
// request parameters etc.)
type BlaimLine struct {
	// Filename is the path of a file that contains an AI-generated code suggestion.
	Filename string `json:"filename"`
	// Range specifies the start and end position of the inserted text.
	Range Range `json:"range"`
	// Text is the raw text of the AI-generated code suggestion.
	Text string `json:"text"`
	// InferenceConfig describes the request sent to the code-generating model,
	// (e.g. the name of the model, temperature etc).
	InferenceConfig InferenceConfig `json:"inference_config"`
}

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

type GitCommit struct {
	Type   int    `json:"type"`
	Name   string `json:"name"`
	Commit string `json:"commit"`
	Ahead  int    `json:"ahead"`
	Behind int    `json:"behind"`
}

type InferenceConfig struct {
	Endpoint    string  `json:"endpoint"`
	MaxLines    int     `json:"maxLines"`
	MaxTokens   int     `json:"maxTokens"`
	Temperature float32 `json:"temperature"`
	ModelName   string  `json:"modelName"`
	ModelFormat string  `json:"modelFormat"`
	Delay       int     `json:"delay"`
}

type AcceptLogLine struct {
	FileName        string
	Position        Position
	Text            string
	HeadGitCommit   GitCommit
	InferenceConfig InferenceConfig
}

func ParseAcceptLogLine(logLine string) (*AcceptLogLine, error) {
	jsonStart := strings.Index(logLine, "] ")
	if jsonStart == -1 {
		// Not an error condition, since the vs code extension logs can contain any arbitrary string.
		// We just ignore anything that doesn't parse like we expect.
		return nil, nil
	}
	jsonText := logLine[jsonStart+2:]
	line := &AcceptLogLine{}
	err := json.Unmarshal([]byte(jsonText), &line)
	return line, err
}

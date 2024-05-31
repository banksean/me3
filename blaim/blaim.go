package blaim

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
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

func ParseLogLine(logLine string) (*AcceptLogLine, error) {
	jsonStart := strings.Index(logLine, "] ")
	if jsonStart == -1 {
		return nil, fmt.Errorf("couldn't find the start of json data in %q", logLine)
	}
	jsonText := logLine[jsonStart+2:]
	line := &AcceptLogLine{}
	err := json.Unmarshal([]byte(jsonText), &line)
	return line, err
}

package blaim

import (
	"reflect"
	"testing"
)

func TestParseLine(t *testing.T) {

	for _, test := range []struct {
		logLine       string
		expected      *AcceptLogLine
		expectedError error
	}{
		{
			logLine: `2024-05-31 14:14:17.804 [info] {"fileName":"inline-completions/playground.js","position":{"line":20,"character":9},"text":"foo(){\n  return \"bar\";\n}","headGitCommit":{"type":0,"name":"logaccepts","commit":"f0d3f3eea79cff732255067ba85588a2bbc4d7c3","ahead":0,"behind":0}
		,"inferenceConfig":{"endpoint":"http://127.0.0.1:11434","bearerToken":"","maxLines":16,"maxTokens":256,"temperature":0.2,"modelName":"stable-code:3b-code-q4_0","modelFormat":"stable-code","delay":250}}`,
			expected: &AcceptLogLine{
				FileName: "inline-completions/playground.js",
				Position: Position{20, 9},
				Text:     "foo(){\n  return \"bar\";\n}",
				HeadGitCommit: GitCommit{
					Type:   0,
					Name:   "logaccepts",
					Commit: "f0d3f3eea79cff732255067ba85588a2bbc4d7c3",
					Ahead:  0,
					Behind: 0,
				},
				InferenceConfig: InferenceConfig{
					Endpoint:    "http://127.0.0.1:11434",
					MaxLines:    16,
					MaxTokens:   256,
					Temperature: 0.2,
					ModelName:   "stable-code:3b-code-q4_0",
					ModelFormat: "stable-code",
					Delay:       250,
				},
			},
		},
	} {
		got, err := ParseLogLine(test.logLine)
		if !reflect.DeepEqual(test.expected, got) {
			t.Errorf("expected %v, got %v\n", test.expected, got)
		}
		if err != test.expectedError {
			t.Errorf("expected %v, got %v\n", test.expectedError, err)
		}
	}
}

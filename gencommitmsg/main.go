package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sashabaranov/go-openai"
)

type generator struct {
	client *openai.Client
}

func (g *generator) prompt() string {
	locale := "en"
	titleChars := 72
	tmpl := `Generate a concise git commit message written in present tense for the following code diff with the given specifications below:
Message language: %s,
Commit message must be a maximum of %d characters.
Exclude anything unnecessary such as translation. Your entire response will be passed directly into git commit.
`
	return fmt.Sprintf(tmpl, locale, titleChars)
}

func (g *generator) commitMessage(diff string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: g.prompt(),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: diff,
				},
			},
		},
	)
	buf, err := json.MarshalIndent(resp, "", "  ")

	return string(buf), err
}

func getDiff(rootDir string) (string, error) {
	cmd := exec.Command("git", "diff", rootDir, ":(exclude)go.*", ":(exclude)*repositories.bzl")
	cmd.Dir = rootDir
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	fmt.Printf("cmd: %+v\n", cmd)
	if err := cmd.Start(); err != nil {
		return "", err
	}

	ret, err := io.ReadAll(stdout)
	if err != nil {
		return "", err
	}

	return string(ret), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <path to git repository>\n", os.Args[0])
		return
	}
	rootDir := os.Args[1]

	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if len(openaiAPIKey) == 0 {
		fmt.Printf("OPENAI_API_KEY environment variable is not set\n")
		return
	}

	g := &generator{
		client: openai.NewClient(openaiAPIKey),
	}
	diff, err := getDiff(rootDir)
	if err != nil {
		fmt.Printf("getDiff error: %v\n", err)
		return
	}
	fmt.Printf("diff: \n%s\n", diff)

	msg, err := g.commitMessage(diff)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	fmt.Printf("response: \n%s\n", msg)
}

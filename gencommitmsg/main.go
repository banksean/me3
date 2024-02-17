package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sashabaranov/go-openai"
)

// This text is taken verbatim from "The seven rules of a great Git commit message": https://cbea.ms/git-commit/
const prompt = `Generate a concise git commit message written in present tense for the following code diff with the given specifications below:
Separate subject from body with a blank line
Limit the subject line to 50 characters
Capitalize the subject line
Do not end the subject line with a period
Use the imperative mood in the subject line
Wrap the body at 72 characters
Use the body to explain what and why vs. how
`

var (
	help        = flag.Bool("h", false, "prtint this help message and exit")
	temperature = flag.Float64("t", 1.0, "temperature for the GPT-3.5-turbo model")
)

type generator struct {
	client *openai.Client
}

func (g *generator) commitMessage(diff string) (string, error) {
	resp, err := g.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       openai.GPT3Dot5Turbo,
			Temperature: float32(*temperature),
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: prompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: diff,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	res := ""
	err = fmt.Errorf("no assistant response")

	for _, choice := range resp.Choices {
		if choice.Message.Role == openai.ChatMessageRoleAssistant {
			res = choice.Message.Content
			err = nil
		}
	}

	return res, err
}

func getDiff(rootDir string) (string, error) {
	cmd := exec.Command("git", "diff", rootDir, ":(exclude)go.mod", ":(exclude)go.sum", ":(exclude)*repositories.bzl")
	cmd.Dir = rootDir
	cmd.Env = os.Environ()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
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
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s <path to git repository>\nOr alternatively: %s $(pwd) > .gitmessage && git commit\n", os.Args[0], os.Args[0])
		os.Exit(1)
	}
	rootDir := os.Args[len(os.Args)-1]

	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if len(openaiAPIKey) == 0 {
		fmt.Fprintf(os.Stderr, "OPENAI_API_KEY environment variable is not set\n")
		os.Exit(1)
	}

	g := &generator{
		client: openai.NewClient(openaiAPIKey),
	}
	diff, err := getDiff(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getDiff error: %v\n", err)
		os.Exit(1)
	}

	msg, err := g.commitMessage(diff)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ChatCompletion error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", msg)
}

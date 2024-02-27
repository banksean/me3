package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"time"

	"github.com/invopop/jsonschema"

	ollama "github.com/jmorganca/ollama/api"
	"github.com/sashabaranov/go-openai"
)

type CommitMessageResponse struct {
	Message string `json:"commit-msg" jsonschema_description:"contents of the commit message"`
}

// This text is taken verbatim from "The seven rules of a great Git commit message": https://cbea.ms/git-commit/
const prompt2 = `Generate a concise git commit message written in present tense for the following code diff with the given specifications below:
Separate subject from body with a blank line
Limit the subject line to 50 characters
Capitalize the subject line
Do not end the subject line with a period
Use the imperative mood in the subject line
Wrap the body at 72 characters
Use the body to explain what and why vs. how
The output should be JSON, and use this schema:
{"commit-msg": "contents of the commit message"}
\n
`
const prompt = `Only use the following information to answer the question. 
- Do not use anything else
- Do not use your own knowledge.
- Do not use your own opinion.
- Do not use anything that is not in the diff.
- Don not use the character "` + "`" + `" or "'" in your answer.
- Be as concise as possible.
- Be as specific as possible.
- Be as accurate as possible.
Task: Write a git commit message for the following diff
`

const (
	generatorFlagOLlama = "ollama"
	generatorFlagOpenAI = "openai"
)

var (
	help              = flag.Bool("h", false, "prtint this help message and exit")
	model             = flag.String("model", "codellama:7b", "name of the ollama model to use")
	generator         = flag.String("generator", generatorFlagOLlama, "generator type")
	temperature       = flag.Float64("t", 1.0, "temperature for the GPT-3.5-turbo model")
	commitMsgFilename = flag.String("commit-msg-file", "", "file to write the commit message to")
	commitSrc         = flag.String("commit-source", "", "source of the commit message")
	commitSHA1        = flag.String("sha1", "", "SHA1 of the commit")
)

type Generator interface {
	GenerateCommitMessage(ctx context.Context, diff string) (string, error)
}

type ollamaGenerator struct {
	model  string
	client *ollama.Client
}

var _ Generator = &ollamaGenerator{}

func NewOLlamaGenerator() (*ollamaGenerator, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	return &ollamaGenerator{
		//model: "llama2:latest",
		// model: "stable-code:3b-code-q4_0",
		//model:  "codellama:7b",
		model:  *model,
		client: client,
	}, nil
}

func (g *ollamaGenerator) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	r := &jsonschema.Reflector{
		AllowAdditionalProperties: false,
		Anonymous:                 true,
		DoNotReference:            true,
	}
	s := r.ReflectFromType(reflect.TypeOf(&CommitMessageResponse{}))
	schemaStr, err := json.MarshalIndent(s.Properties, "", "  ")
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(os.Stderr, "schema: %s\n", string(schemaStr))
	streaming := false
	request := ollama.GenerateRequest{
		Model:   g.model,
		Prompt:  diff,
		Context: []int{},
		//Format:  "json",
		Template: `[INST] <<SYS>>{{ .System }}<</SYS>>

		{{ .Prompt }} [/INST]`,
		Stream: &streaming,
		System: prompt, // + string(schemaStr),
	}
	ret := ""
	fn := func(response ollama.GenerateResponse) error {
		ret += response.Response
		return nil
	}

	ctx, done := context.WithTimeout(ctx, 10*time.Second)
	defer done()
	fmt.Fprintf(os.Stderr, "request: %#v\n", request)
	if err := g.client.Generate(ctx, &request, fn); err != nil {
		if errors.Is(err, context.Canceled) {
			return err.Error(), err
		}
		return err.Error(), err
	}
	fmt.Fprintf(os.Stderr, "raw response text: %s\n", ret)
	resp := &CommitMessageResponse{}
	if err := json.Unmarshal([]byte(ret), resp); err != nil {
		fmt.Fprintf(os.Stderr, "could not un-marshal response json:\n%s\n", ret)
		return "", err
	}
	return resp.Message, nil
}

type openAIGenerator struct {
	client *openai.Client
}

var _ Generator = &openAIGenerator{}

func NewOpenAIGenerator() (*openAIGenerator, error) {
	openaiAPIKey := os.Getenv("OPENAI_API_KEY")
	if len(openaiAPIKey) == 0 {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	g := &openAIGenerator{
		client: openai.NewClient(openaiAPIKey),
	}
	return g, nil
}

func (g *openAIGenerator) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
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

	diff, err := getDiff(rootDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getDiff error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	var g Generator

	if *generator == generatorFlagOLlama {
		g, err = NewOLlamaGenerator()
		if err != nil {
			fmt.Fprintf(os.Stderr, "NewOLlamaGenerator error: %v\n", err)
			os.Exit(1)
		}
	} else if *generator == generatorFlagOpenAI {
		g, err = NewOpenAIGenerator()
		if err != nil {
			fmt.Fprintf(os.Stderr, "NewOpenAIGenerator error: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Fprintf(os.Stderr, "unrecognized generator type: %v\n", *generator)
		os.Exit(1)
	}

	msg, err := g.GenerateCommitMessage(ctx, diff)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GenerateCommitMessage error: %v\n", err)
		os.Exit(1)
	}

	if *commitMsgFilename != "" {
		err := os.WriteFile(filepath.Join(rootDir, *commitMsgFilename), []byte(msg), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WriteFile error: %v\n", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s\n", msg)
}

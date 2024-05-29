package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"text/template"

	"bitbucket.org/creachadair/stringset"
	ollama "github.com/jmorganca/ollama/api"
)

var (
	spyMasterPromptTmpl = template.Must(template.New("spymaster").Parse(
		`Your task is to provide me with a single word clue to help me identify one of the words in the following list:
	{{range .OurWords }}{{. | printf "%q"}} {{end}}
Your clue cannot be any of the words in that list.
Your clue cannot be a slight variation of any of the words in that list.
Your clue must NOT be associated with any of the words in the following list:
	{{range .NotOurWords }}{{. | printf "%q"}} {{end}}
In particular, DO NOT offer a clue that might suggest the word {{ .AssassinWord | printf "%q" }}, because you will cause us to lose the game.
Respond only with the single word clue.  
Do not provide any explanation for why you chose the single word clue.
> `))

	fieldAgentPromptTmpl = template.Must(template.New("fieldagent").Parse(
		`Based on the following clue: {{.Clue | printf "%q"}},
	Your task is to identify one of the words in the following list:
	{{range .Words }}{{. | printf "%q"}} {{end}}
	Your guess MUST BE one and only one word from the above list.
	Do not guess a word that is not in that list.
	Your guess MUST NOT BE the word {{.Clue | printf "%q"}}.
	Respond only with the single word, lowercase, with no punctuation.
	Do NOT respond with any text OTHER THAN THAT ONE WORD.
	> `))
)

type OLlamaSpyMasterTurn struct {
	team       string
	model      string
	promptTmpl *template.Template
	client     *ollama.Client
}

var _ Player = &OLlamaSpyMasterTurn{}

func NewOLllamaSpyMaster(team, model string) (*OLlamaSpyMasterTurn, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	ret := &OLlamaSpyMasterTurn{
		team:       team,
		client:     client,
		model:      model,
		promptTmpl: spyMasterPromptTmpl,
	}

	return ret, nil
}

func (s *OLlamaSpyMasterTurn) Team() string {
	return s.team
}

func (s *OLlamaSpyMasterTurn) Move(game *gameBoard) error {
	game.WriteTable(os.Stdout, true)

	ourRemainingWords := game.cards[s.team].Clone()
	ourRemainingWords.Remove(*game.guessedWords)

	notOurWords := stringset.New()

	for team, teamCards := range game.cards {
		if team == s.team {
			continue
		}
		teamCards := teamCards.Clone()
		teamCards.Remove(*game.guessedWords)
		notOurWords.Add(teamCards.Elements()...)
	}
	assassinWord := game.cards[ASSASSIN].Elements()[0]
	promptData := map[string]any{
		"OurWords":     ourRemainingWords.Elements(),
		"NotOurWords":  notOurWords.Elements(),
		"AssassinWord": assassinWord,
	}
	prompt := &strings.Builder{}
	if err := s.promptTmpl.Execute(prompt, promptData); err != nil {
		return err
	}
	streaming := false
	request := ollama.GenerateRequest{
		Model:   s.model,
		Prompt:  prompt.String(),
		Context: []int{},
		Stream:  &streaming,
	}
	input := ""
	fn := func(response ollama.GenerateResponse) error {
		input += response.Response
		return nil
	}
	ctx := context.Background()
	err := s.client.Generate(ctx, &request, fn)
	if err != nil {
		return err
	}

	input = strings.ToLower(input)
	game.state = game.transitions[s.team+"CLUE"]
	fmt.Printf("%s team (%s) SpyMaster's clue: %q\n", s.team, s.model, input)
	game.clue[s.team] = input
	return nil
}

type OLlamaFieldAgentTurn struct {
	team       string
	model      string
	promptTmpl *template.Template
	client     *ollama.Client
}

var _ Player = &OLlamaFieldAgentTurn{}

func NewOLlamaFieldAgent(team, model string) (*OLlamaFieldAgentTurn, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	ret := &OLlamaFieldAgentTurn{
		team:       team,
		client:     client,
		promptTmpl: fieldAgentPromptTmpl,
		model:      model,
	}

	return ret, nil
}

func (s *OLlamaFieldAgentTurn) Team() string {
	return s.team
}

func (s *OLlamaFieldAgentTurn) Move(game *gameBoard) error {
	var err error
	for errRetries := 3; errRetries > 0; errRetries-- {
		if errRetries < 3 {
			fmt.Printf("\tRetrying guess (%d retries remaining), due to previous error\n", errRetries)
		}
		remainingWords := stringset.New()
		for _, cards := range game.cards {
			remainingWords.Add(cards.Elements()...)
		}
		remainingWords.Remove(*game.guessedWords)
		clue := game.clue[s.team]
		promptData := map[string]any{
			"Clue":  clue,
			"Words": remainingWords.Elements(),
		}
		prompt := &strings.Builder{}
		if err := s.promptTmpl.Execute(prompt, promptData); err != nil {
			return err
		}
		streaming := false
		request := ollama.GenerateRequest{
			Model:   s.model,
			Prompt:  prompt.String(),
			Context: []int{},
			Stream:  &streaming,
		}
		input := ""
		fn := func(response ollama.GenerateResponse) error {
			input += response.Response
			return nil
		}
		ctx := context.Background()
		err := s.client.Generate(ctx, &request, fn)
		if err != nil {
			return err
		}
		fmt.Printf("%s team (%s) FieldAgent guessed %q based on the clue %q\n", s.team, s.model, input, clue)
		team := ""
		team, err = game.guess(input)
		if err != nil {
			fmt.Printf("\tError during guess: %v\n", err)
			continue
		}

		if team == s.team {
			fmt.Printf("\tCORRECT: ")
		} else {
			fmt.Printf("\tINCORRECT: ")
		}
		fmt.Printf("\t%q belongs to team %s\n", input, team)
		game.state = game.transitions[s.team+"GUESS"]
		break
	}
	return err
}

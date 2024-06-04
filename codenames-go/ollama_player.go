package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"text/template"

	"bitbucket.org/creachadair/stringset"
	ollama "github.com/jmorganca/ollama/api"
)

const (
	moveRetryCount  = 3
	preserveContext = false
)

var (
	spyMasterPromptTmpl = template.Must(template.New("spymaster").Parse(
		`Your task is to provide me with a single word clue to help me identify one of the target words in the following OUR_CARDS list:
	{{range .OurWords }}{{. | printf "%q"}} {{end}}
Your clue cannot be any of the words in OUR_CARDS.
Your clue cannot be a slight variation of any of the words in OUR_CARDS.
Your clue should make me think of one of the words in OUR_CARDS.
Your clue must NOT be associated with any of the words in the following THEIR_CARDS list:
	{{range .NotOurWords }}{{. | printf "%q"}} {{end}}
Your clue must not be any of the words in the THEIR_CARDS list.
In particular, DO NOT offer a clue that might suggest the word {{ .AssassinWord | printf "%q" }}, because you will cause us to lose the game.
Respond with a json object like this example:
{
	"clue": "shoe",
	"target": "sock, fit",
	"explanation": "The word 'shoe' is related to the words 'sock' and 'fit' from our word list because they both have to do with feet and none of the  words in THIER_CARDS are associated with shoes."
}
> `))

	fieldAgentPromptTmpl = template.Must(template.New("fieldagent").Parse(
		`Based on the following clue: {{.Clue | printf "%q"}},
	Your task is to identify one of the words in the following list:
	{{range .Words }}{{. | printf "%q"}} {{end}}
	Your guess MUST BE one and only one word from the above list.
	Do not guess a word that is not in that list.
	Your guess MUST NOT BE any of these words: 
	{{.Clue | printf "%q"}}{{range .PreviousGuessErrors }} {{. | printf "%q"}}{{end}}.
	Respond with a json object like this example:
	{
		"guess": "sock",
		"explanation": "Based on the clue word 'shoe', the word 'sock' seems like the best match because those two objects are frequently used together and none of the other cards are associated with the word 'sock'."
	}
	> `))
)

type OLlamaSpyMaster struct {
	team            string
	model           string
	promptTmpl      *template.Template
	client          *ollama.Client
	modelContext    []int
	preserveContext bool
}

var _ Player = &OLlamaSpyMaster{}

func NewOLllamaSpyMaster(team, model string) (*OLlamaSpyMaster, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	ret := &OLlamaSpyMaster{
		team:            team,
		client:          client,
		model:           model,
		promptTmpl:      spyMasterPromptTmpl,
		preserveContext: preserveContext,
	}

	return ret, nil
}

func (s *OLlamaSpyMaster) Team() string {
	return s.team
}

type SpyMasterResponse struct {
	Clue        string `json:"clue"`
	Target      string `json:"target"`
	Explanation string `json:"explanation"`
}

func (s *OLlamaSpyMaster) Move(game *gameBoard) error {
	return retry(moveRetryCount, func() error {
		fmt.Printf("%s team (%s) SpyMaster is thinking of a clue...\n", s.team, s.model)
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
			Format:  "json",
			Context: s.modelContext,
			Stream:  &streaming,
		}
		input := ""
		fn := func(response ollama.GenerateResponse) error {
			input += response.Response
			if s.preserveContext {
				s.modelContext = response.Context
			}
			return nil
		}
		ctx := context.Background()
		err := s.client.Generate(ctx, &request, fn)
		if err != nil {
			return err
		}

		spyMasterResponse := &SpyMasterResponse{}
		err = json.Unmarshal([]byte(input), spyMasterResponse)
		if err != nil {
			return err
		}
		input = strings.ToLower(input)
		game.state = game.transitions[s.team+"CLUE"]
		clue := spyMasterResponse.Clue
		explanation := spyMasterResponse.Explanation

		fmt.Printf("%s team (%s) SpyMaster's clue: %q\n", s.team, s.model, clue)
		game.clue[s.team] = clue
		game.explanation[s.team] = fmt.Sprintf("[%s] %s", spyMasterResponse.Target, explanation)
		return nil
	})
}

type OLlamaFieldAgent struct {
	team            string
	model           string
	promptTmpl      *template.Template
	modelContext    []int
	client          *ollama.Client
	preserveContext bool
}

var _ Player = &OLlamaFieldAgent{}

func NewOLlamaFieldAgent(team, model string) (*OLlamaFieldAgent, error) {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		return nil, err
	}

	ret := &OLlamaFieldAgent{
		team:            team,
		client:          client,
		promptTmpl:      fieldAgentPromptTmpl,
		model:           model,
		preserveContext: preserveContext,
	}

	return ret, nil
}

func (s *OLlamaFieldAgent) Team() string {
	return s.team
}

type FieldAgentResponse struct {
	Guess       string `json:"guess"`
	Explanation string `json:"explanation"`
}

func (s *OLlamaFieldAgent) Move(game *gameBoard) error {
	previousGuessErrors := []string{}

	return retry(moveRetryCount, func() error {
		remainingWords := stringset.New()
		for _, cards := range game.cards {
			remainingWords.Add(cards.Elements()...)
		}
		remainingWords.Remove(*game.guessedWords)
		clue := game.clue[s.team]
		promptData := map[string]any{
			"Clue":                clue,
			"Words":               remainingWords.Elements(),
			"PreviousGuessErrors": previousGuessErrors,
		}
		prompt := &strings.Builder{}
		if err := s.promptTmpl.Execute(prompt, promptData); err != nil {
			return err
		}
		streaming := false
		request := ollama.GenerateRequest{
			Model:   s.model,
			Prompt:  prompt.String(),
			Format:  "json",
			Context: s.modelContext,
			Stream:  &streaming,
		}
		input := ""
		fn := func(response ollama.GenerateResponse) error {
			input += response.Response
			if s.preserveContext {
				s.modelContext = response.Context
			}
			return nil
		}
		ctx := context.Background()
		err := s.client.Generate(ctx, &request, fn)
		if err != nil {
			return err
		}

		fieldAgentResponse := &FieldAgentResponse{}
		err = json.Unmarshal([]byte(input), fieldAgentResponse)
		if err != nil {
			return err
		}

		fmt.Printf("%s team (%s) FieldAgent guessed %q based on the clue %q\n", s.team, s.model, fieldAgentResponse.Guess, clue)
		team := ""
		team, err = game.guess(fieldAgentResponse.Guess)
		if err != nil {
			previousGuessErrors = append(previousGuessErrors, fieldAgentResponse.Guess)
			return err
		}

		if team == s.team {
			fmt.Printf("\tCORRECT: ")
		} else {
			fmt.Printf("\tINCORRECT: ")
		}
		fmt.Printf("\t%q belongs to team %s\n", fieldAgentResponse.Guess, team)
		game.state = game.transitions[s.team+"GUESS"]
		game.explanation[s.team] = fieldAgentResponse.Explanation
		return nil
	})
}

package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"bitbucket.org/creachadair/stringset"
	ollama "github.com/jmorganca/ollama/api"
)

type OLlamaSpyMasterTurn struct {
	team   string
	model  string
	client *ollama.Client
}

var _ Player = &OLlamaSpyMasterTurn{}

func NewOLllamaSpyMaster(team, model string) *OLlamaSpyMasterTurn {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := &OLlamaSpyMasterTurn{
		team:   team,
		client: client,
		model:  model,
	}

	return ret
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

	prompt := fmt.Sprintf(`Your task is to provide me with a single word clue to help me identify one of the words in the following list:
	[%s]
Your clue cannot be any of the words in that list.
Your clue cannot be a slight variation of any of the words in that list.
Your clue must NOT be associated with any of the words in the following list:
	[%s]
In particular, DO NOT offer a clue that might suggest the word %q, because you will cause us to lose the game.
Respond only with the single word clue.  
Do not provide any explanation for why you chose the single word clue.
> `,
		strings.Join(ourRemainingWords.Elements(), ", "), strings.Join(notOurWords.Elements(), ", "), assassinWord)

	streaming := false
	request := ollama.GenerateRequest{
		Model:   s.model,
		Prompt:  prompt,
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

	game.state = game.transitions[s.team+"CLUE"]
	fmt.Printf("%s team SpyMaster's clue: %q\n", s.team, input)
	game.clue[s.team] = input
	return nil
}

type OLlamaFieldAgentTurn struct {
	team   string
	model  string
	client *ollama.Client
}

var _ Player = &OLlamaFieldAgentTurn{}

func NewOLllamaOLlamaFieldAgent(team, model string) *OLlamaFieldAgentTurn {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := &OLlamaFieldAgentTurn{
		team:   team,
		client: client,
		model:  model,
	}

	return ret
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
		prompt := fmt.Sprintf(`Based on the following clue: %s,
Your task is to identify one of the words in the following list:
	[%s]
Your guess MUST BE one and only one word from the above list.
Do not guess a word that is not in that list.
Your guess MUST NOT BE the word %q.
Respond only with the single word, lowercase, with no punctuation.
Do NOT respond with any text OTHER THAN THAT ONE WORD.
> `,
			clue, strings.Join(remainingWords.Elements(), ", "), clue)

		streaming := false
		request := ollama.GenerateRequest{
			Model:   s.model,
			Prompt:  prompt,
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
		fmt.Printf("%s team FieldAgent guessed %q based on the clue %q\n", s.team, input, clue)
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

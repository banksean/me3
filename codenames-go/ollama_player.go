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
	client *ollama.Client
}

var _ Player = &OLlamaSpyMasterTurn{}

func NewOLllamaSpyMaster(team string) *OLlamaSpyMasterTurn {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := &OLlamaSpyMasterTurn{
		team:   team,
		client: client,
	}

	return ret
}

func (s *OLlamaSpyMasterTurn) Move(game *gameBoard) error {
	input, err := s.PromptInput(game)
	if err != nil {
		return err
	}
	return s.ProcessInput(game, input)
}

func (s *OLlamaSpyMasterTurn) PromptInput(game *gameBoard) (string, error) {
	ourRemainingWords := game.cards[s.team].Clone()
	ourRemainingWords.Remove(*game.guessedWords)

	notOurWords := stringset.New()

	for team, teamCards := range game.cards {
		if team == s.team {
			continue
		}
		teamCards := teamCards.Clone()
		teamCards.Remove(*game.guessedWords)
		notOurWords.Union(teamCards)
	}

	prompt := fmt.Sprintf(`%s team spymaster
Your task is to provide me with a single word clue to help me identify one of the words in the following list:
	%s
Your clue cannot be any of the words in that list.
Your clue cannot be a slight variation of any of the words in that list.
Your clue must not be associated with any of the words in the following list:
	%s
Respond only with the single word clue.  
Do not provide any explanation for why you chose the single word clue.
> `,
		s.team, strings.Join(ourRemainingWords.Elements(), ", "), strings.Join(notOurWords.Elements(), ", "))

	streaming := false
	request := ollama.GenerateRequest{
		Model:   "llama3",
		Prompt:  prompt,
		Context: []int{},
		Stream:  &streaming,
	}
	ret := ""
	fn := func(response ollama.GenerateResponse) error {
		ret += response.Response
		return nil
	}
	ctx := context.Background()
	if err := s.client.Generate(ctx, &request, fn); err != nil {
		return err.Error(), err
	}

	return ret, nil
}

func (s *OLlamaSpyMasterTurn) ProcessInput(game *gameBoard, input string) error {
	game.state = game.transitions[s.team+"CLUE"]

	game.clue[s.team] = input
	return nil
}

type OLlamaFieldAgentTurn struct {
	team   string
	client *ollama.Client
}

var _ Player = &OLlamaFieldAgentTurn{}

func NewOLllamaOLlamaFieldAgentTurn(team string) *OLlamaFieldAgentTurn {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := &OLlamaFieldAgentTurn{
		team:   team,
		client: client,
	}

	return ret
}

func (s *OLlamaFieldAgentTurn) Move(game *gameBoard) error {
	input, err := s.PromptInput(game)
	if err != nil {
		return err
	}
	return s.ProcessInput(game, input)
}

func (s *OLlamaFieldAgentTurn) PromptInput(game *gameBoard) (string, error) {
	ourRemainingWords := game.cards[s.team].Clone()
	ourRemainingWords.Remove(*game.guessedWords)

	notOurWords := stringset.New()

	for team, teamCards := range game.cards {
		if team == s.team {
			continue
		}
		teamCards := teamCards.Clone()
		teamCards.Remove(*game.guessedWords)
		notOurWords.Union(teamCards)
	}

	prompt := fmt.Sprintf(`%s team spymaster
Your task is to provide me with a single word clue to help me identify one of the words in the following list:
	%s
Your clue cannot be any of the words in that list.
Your clue cannot be a slight variation of any of the words in that list.
Your clue must not be associated with any of the words in the following list:
	%s
Respond only with the single word clue.  
Do not provide any explanation for why you chose the single word clue.
> `,
		s.team, strings.Join(ourRemainingWords.Elements(), ", "), strings.Join(notOurWords.Elements(), ", "))

	streaming := false
	request := ollama.GenerateRequest{
		Model:   "llama3",
		Prompt:  prompt,
		Context: []int{},
		Stream:  &streaming,
	}
	ret := ""
	fn := func(response ollama.GenerateResponse) error {
		ret += response.Response
		return nil
	}
	ctx := context.Background()
	if err := s.client.Generate(ctx, &request, fn); err != nil {
		return err.Error(), err
	}

	return ret, nil
}

func (s *OLlamaFieldAgentTurn) ProcessInput(game *gameBoard, input string) error {
	game.state = game.transitions[s.team+"CLUE"]
	game.clue[s.team] = input
	return nil
}

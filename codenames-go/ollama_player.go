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
	game   *gameBoard
	team   string
	client *ollama.Client
}

func NewOLllamaSpyMaster(game *gameBoard, team string) *OLlamaSpyMasterTurn {
	client, err := ollama.ClientFromEnvironment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	ret := &OLlamaSpyMasterTurn{
		team:   team,
		game:   game,
		client: client,
	}

	return ret
}

func (s *OLlamaSpyMasterTurn) PromptInput() (string, error) {
	ourRemainingWords := s.game.cards[s.team].Clone()
	ourRemainingWords.Remove(*s.game.guessedWords)

	notOurWords := stringset.New()

	for team, teamCards := range s.game.cards {
		if team == s.team {
			continue
		}
		teamCards := teamCards.Clone()
		teamCards.Remove(*s.game.guessedWords)
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

func (s *OLlamaSpyMasterTurn) ProcessInput(input string) error {
	s.game.state = s.game.transitions[s.team+"CLUE"]

	fat := s.game.state.(*HumanFieldAgentTurn)
	fat.clue = input
	return nil
}

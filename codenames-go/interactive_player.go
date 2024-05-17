package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"

	"bitbucket.org/creachadair/stringset"
)

type HumanSpyMasterTurn struct {
	team string
	rl   *readline.Instance
}

var _ Player = &HumanSpyMasterTurn{}

func (s *HumanSpyMasterTurn) Move(game *gameBoard) error {
	input, err := s.PromptInput(game)
	if err != nil {
		return err
	}
	return s.ProcessInput(game, input)
}

func (s *HumanSpyMasterTurn) PromptInput(game *gameBoard) (string, error) {
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
		notOurWords.Union(teamCards)
	}

	fmt.Printf(`%s team spymaster
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

	input, err := s.rl.Readline()

	return input, err
}

func (s *HumanSpyMasterTurn) ProcessInput(game *gameBoard, input string) error {
	game.state = game.transitions[s.team+"CLUE"]
	game.clue[s.team] = input
	return nil
}

type HumanFieldAgentTurn struct {
	team string
	rl   *readline.Instance
}

var _ Player = &HumanFieldAgentTurn{}

func (s *HumanFieldAgentTurn) Move(game *gameBoard) error {
	input, err := s.PromptInput(game)
	if err != nil {
		return err
	}
	return s.ProcessInput(game, input)
}

func (s *HumanFieldAgentTurn) PromptInput(game *gameBoard) (string, error) {
	game.WriteTable(os.Stdout, false)
	clue := game.clue[s.team]
	fmt.Printf("%s team field agent: make a guess for %q\n> ", s.team, clue)
	input, err := s.rl.Readline()
	return input, err
}

func (s *HumanFieldAgentTurn) ProcessInput(game *gameBoard, input string) error {
	team, err := game.guess(input)
	if err != nil {
		return err
	}
	if team == s.team {
		fmt.Printf("\nCORRECT ")
	} else {
		fmt.Printf("\nINCORRECT ")
	}
	fmt.Printf("%q belongs to team %s\n\n", input, team)
	game.state = game.transitions[s.team+"GUESS"]

	return nil
}

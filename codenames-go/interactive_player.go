package main

import (
	"fmt"
	"os"
	"strings"
)

type HumanSpyMasterTurn struct {
	game *gameBoard
	team string
}

func (s *HumanSpyMasterTurn) PromptInput() (string, error) {
	s.game.WriteTable(os.Stdout, true)

	ourCards := s.game.cards[s.team]

	ourWords := []string{}
	notOurWords := []string{}

	for w := range ourCards {
		ourWords = append(ourWords, w)
	}
	for team, teamCards := range s.game.cards {
		if team == s.team {
			continue
		}

		for w := range teamCards {
			notOurWords = append(notOurWords, w)
		}
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
		s.team, strings.Join(ourWords, ", "), strings.Join(notOurWords, ", "))

	input, err := s.game.rl.Readline()

	return input, err
}

func (s *HumanSpyMasterTurn) ProcessInput(input string) error {
	s.game.state = s.game.transitions[s.team+"CLUE"]
	fat := s.game.state.(*HumanFieldAgentTurn)
	fat.clue = input
	return nil
}

type HumanFieldAgentTurn struct {
	game *gameBoard
	team string
	clue string
}

func (s *HumanFieldAgentTurn) PromptInput() (string, error) {
	s.game.WriteTable(os.Stdout, false)

	fmt.Printf("%s team field agent: make a guess for %q\n> ", s.team, s.clue)
	input, err := s.game.rl.Readline()
	return input, err
}

func (s *HumanFieldAgentTurn) ProcessInput(input string) error {
	team, err := s.game.guess(input)
	if err != nil {
		return err
	}
	if team == s.team {
		fmt.Printf("\nCORRECT ")
	} else {
		fmt.Printf("\nINCORRECT ")
	}
	fmt.Printf("%q belongs to team %s\n\n", input, team)
	s.game.state = s.game.transitions[s.team+"GUESS"]

	return nil
}

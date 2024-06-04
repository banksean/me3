package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

type HumanSpyMasterTurn struct {
	team string
	rl   *readline.Instance
}

var _ Player = &HumanSpyMasterTurn{}

func (s *HumanSpyMasterTurn) Team() string {
	return s.team
}

func (s *HumanSpyMasterTurn) Move(game *gameBoard) error {
	game.WriteTable(os.Stdout, true)

	fmt.Printf(`%s team spymaster, provide a clue for one of your words:
> `,
		s.team)

	input, err := s.rl.Readline()
	if err != nil {
		return err
	}
	if strings.HasPrefix(input, "/explain") {
		team := s.team
		if len(input) > len("/explain") {
			team = strings.TrimSpace(input[len("/explain"):])
		}
		fmt.Printf("explanation: %s\n", game.explanation[team])
		return nil
	}

	game.state = game.transitions[s.team+"CLUE"]
	game.clue[s.team] = input
	return nil
}

type HumanFieldAgentTurn struct {
	team string
	rl   *readline.Instance
}

var _ Player = &HumanFieldAgentTurn{}

func (s *HumanFieldAgentTurn) Team() string {
	return s.team
}

func (s *HumanFieldAgentTurn) Move(game *gameBoard) error {
	fmt.Printf("\n")
	game.WriteTable(os.Stdout, false)
	fmt.Printf("\n")
	clue := game.clue[s.team]
	fmt.Printf("%s team field agent: make a guess for %q\n> ", s.team, clue)
	input, err := s.rl.Readline()

	if err != nil {
		return err
	}

	if strings.HasPrefix(input, "/explain") {
		team := s.team
		if len(input) > len("/explain") {
			team = strings.TrimSpace(input[len("/explain"):])
		}
		fmt.Printf("explanation: %s\n", game.explanation[team])
		return nil
	}

	team, err := game.guess(input)
	if err != nil {
		fmt.Printf("error: %s", err)
		return nil
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

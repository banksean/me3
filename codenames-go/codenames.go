package main

import (
	"fmt"
	"math/rand"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/jedib0t/go-pretty/v6/table"

	_ "embed"
)

var (
	//go:embed wordlist.txt
	wordlistFile string
)

type gameBoard struct {
	cards        map[string]map[string]any
	guessedWords map[string]any
	teamForCard  map[string]string
	state        gameState
}

/*
States:
- RedClue
- RedGuess
- BlueClue
- BlueGuess
*/

// "State Machine" pattern.
type gameState interface {
	PromptInput() (string, error)
	Guess(string) error
	Clue(string, int) error
	Pass() error
}

type SpyMasterTurn struct {
	game gameBoard
	team string
}

func (s *SpyMasterTurn) PromptInput() (string, error) {
	fmt.Printf("%s team spymaster: offer a clue\n> ", s.team)
	return "", nil
}

func (s *SpyMasterTurn) Guess(string) error {
	return fmt.Errorf("spymaster cannot make guesses")
}

func (s *SpyMasterTurn) Clue(clue string, num int) error {
	return nil
}

func (s *SpyMasterTurn) Pass() error {
	return fmt.Errorf("spymaster cannot pass")
}

type FieldAgentTurn struct {
	game     gameBoard
	team     string
	clue     string
	numWords int
}

func (s *FieldAgentTurn) PromptInput() (string, error) {
	fmt.Printf("%s team field agent: make a guess for %q, %d\n> ", s.team, s.clue, s.numWords)
	return "", nil
}

func (s *FieldAgentTurn) Guess(string) error {
	return nil
}

func (s *FieldAgentTurn) Clue(clue string, num int) error {
	return fmt.Errorf("field agent cannot offer clues")
}

func (s *FieldAgentTurn) Pass() error {
	return fmt.Errorf("field agent cannot pass")
}

func (g *gameBoard) String() string {
	tw := table.NewWriter()
	allCards := []string{}
	for _, cards := range g.cards {
		for c := range cards {
			allCards = append(allCards, c)
		}
	}
	for i, c := range allCards {
		if _, ok := g.guessedWords[c]; ok {
			team := g.teamForCard[c]
			allCards[i] = fmt.Sprintf("%s [%s]", c, team)
		}
	}
	sort.Strings(allCards)

	for row := 0; row < 5; row++ {
		tableRow := table.Row{}
		rowCards := allCards[row*5 : row*5+5]
		for _, c := range rowCards {
			tableRow = append(tableRow, c)
		}
		tw.AppendRow(tableRow)
	}
	return tw.Render()
}

func (g *gameBoard) currentScore() map[string]int {
	ret := map[string]int{}
	for w := range g.guessedWords {
		for team, teamCards := range g.cards {
			if _, ok := teamCards[w]; ok {
				ret[team]++
			}
		}
	}
	return ret
}

func (g *gameBoard) guess(word string) (string, error) {
	if _, ok := g.guessedWords[word]; ok {
		return "", fmt.Errorf("already guessed %q", word)
	}
	g.guessedWords[word] = nil
	for t, c := range g.cards {
		if _, ok := c[word]; ok {
			return t, nil
		}
	}
	return "", fmt.Errorf("%q is not one of the cards currently in play", word)
}

type deck []string

func draw(d deck, n int) (map[string]any, deck) {
	ret := map[string]any{}
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(d))
		ret[d[idx]] = nil
		d = append(d[:idx], d[idx+1:]...)
	}

	return ret, d
}

func main() {
	words := strings.Split(wordlistFile, "\n")
	d := deck(words)

	game := &gameBoard{
		cards:        make(map[string]map[string]any),
		guessedWords: make(map[string]any),
		teamForCard:  make(map[string]string),
	}
	game.cards["RED"], d = draw(d, 8)
	game.cards["BLUE"], d = draw(d, 9)
	game.cards["BYSTANDER"], d = draw(d, 7)
	game.cards["ASSASSIN"], _ = draw(d, 1)
	for team, cards := range game.cards {
		for card := range cards {
			game.teamForCard[card] = team
		}
	}

	rl, err := readline.New("> ")
	if err != nil {
		panic(err)
	}
	defer rl.Close()

	for {
		fmt.Printf("game:\n%+v\n", game.String())
		fmt.Printf("score:\n%+v\n", game.currentScore())
		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}

		team, err := game.guess(line)
		if err != nil {
			fmt.Printf("error: %v\n", err)
		} else {
			fmt.Printf("%q belongs to %q\n", line, team)
		}
	}
}

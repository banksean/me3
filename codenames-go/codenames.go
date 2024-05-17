package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"

	_ "embed"
)

const (
	RED       = "RED"
	BLUE      = "BLUE"
	BYSTANDER = "BYSTANDER"
	ASSASSIN  = "ASSASSIN"
)

var (
	//go:embed wordlist.txt
	wordlistFile string
)

type gameBoard struct {
	cards        map[string]*stringset.Set
	guessedWords *stringset.Set
	teamForCard  map[string]string
	state        Player
	transitions  map[string]Player
	clue         map[string]string
}

type Team struct {
	Name                  string
	SpyMaster, FieldAgent Player
}

func (t *Team) Turn(g *gameBoard) error {
	return nil
}

// "State Machine" pattern.
type Player interface {
	Move(*gameBoard) error
}

const GUESSED = "guessed-"

var (
	defaultColor = tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor, tablewriter.BgBlackColor}
	teamColor    = map[string]tablewriter.Colors{
		RED:                 {tablewriter.Normal, tablewriter.FgRedColor, tablewriter.BgBlackColor},
		BLUE:                {tablewriter.Normal, tablewriter.FgBlueColor, tablewriter.BgBlackColor},
		BYSTANDER:           {tablewriter.Normal, tablewriter.FgWhiteColor, tablewriter.BgBlackColor},
		ASSASSIN:            {tablewriter.Normal, tablewriter.FgYellowColor, tablewriter.BgBlackColor},
		RED + GUESSED:       {tablewriter.Normal, tablewriter.FgBlackColor, tablewriter.BgRedColor},
		BLUE + GUESSED:      {tablewriter.Normal, tablewriter.FgBlackColor, tablewriter.BgBlueColor},
		BYSTANDER + GUESSED: {tablewriter.Normal, tablewriter.FgBlackColor, tablewriter.BgWhiteColor},
		ASSASSIN + GUESSED:  {tablewriter.Normal, tablewriter.FgBlackColor, tablewriter.BgYellowColor},
	}
)

func (g *gameBoard) WriteTable(w io.Writer, spyMasterView bool) {
	tw := tablewriter.NewWriter(w)
	tw.SetBorder(true)
	tw.SetRowLine(true)
	tw.SetAlignment(tablewriter.ALIGN_CENTER)
	allCards := []string{}
	for _, cards := range g.cards {
		for c := range *cards {
			allCards = append(allCards, c)
		}
	}
	sort.Strings(allCards)
	allColors := []tablewriter.Colors{}
	for i, c := range allCards {
		guessed := g.guessedWords.Contains(c)
		if guessed || spyMasterView {
			team := g.teamForCard[c]
			color := teamColor[team]

			if guessed {
				allCards[i] = strings.ToUpper(c)
				color = teamColor[team+GUESSED]
			}
			allColors = append(allColors, color)

		} else {
			allColors = append(allColors, defaultColor)
		}
	}

	for row := 0; row < 5; row++ {
		rowCards := allCards[row*5 : row*5+5]
		rowColors := allColors[row*5 : row*5+5]

		tw.Rich(rowCards, rowColors)
	}
	tw.Render()
}

func (g *gameBoard) currentScore() map[string]int {
	ret := map[string]int{}
	for w := range *g.guessedWords {
		for team, teamCards := range g.cards {
			if teamCards.Contains(w) {
				ret[team]++
			}
		}
	}
	return ret
}

func (g *gameBoard) guess(word string) (string, error) {
	if g.guessedWords.Contains(word) {
		return "", fmt.Errorf("already guessed %q", word)
	}
	g.guessedWords.Add(word)
	for t, c := range g.cards {
		if c.Contains(word) {
			return t, nil
		}
	}
	return "", fmt.Errorf("%q is not one of the cards currently in play", word)
}

func draw(d []string, n int) (*stringset.Set, []string) {
	ret := stringset.New()
	for i := 0; i < n; i++ {
		idx := rand.Intn(len(d))
		ret.Add(d[idx])
		d = append(d[:idx], d[idx+1:]...)
	}

	return &ret, d
}

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func NewGameBoard(d []string) *gameBoard {
	guessed := stringset.New()
	game := &gameBoard{
		cards:        make(map[string]*stringset.Set),
		guessedWords: &guessed,
		teamForCard:  make(map[string]string),
		clue:         map[string]string{},
	}
	game.cards[RED], d = draw(d, 8)
	game.cards[BLUE], d = draw(d, 9)
	game.cards[BYSTANDER], d = draw(d, 7)
	game.cards[ASSASSIN], _ = draw(d, 1)
	for team, cards := range game.cards {
		for _, card := range cards.Elements() {
			game.teamForCard[card] = team
		}
	}

	return game
}

func main() {
	rl, err := readline.NewEx(
		&readline.Config{
			InterruptPrompt:     "^C",
			Prompt:              "> ",
			FuncFilterInputRune: filterInput,
		},
	)

	if err != nil {
		panic(err)
	}
	defer rl.Close()

	words := strings.Split(wordlistFile, "\n")
	game := NewGameBoard(words)

	//ollamaSpyMasterRed := NewOLllamaSpyMaster(RED)
	ollamaSpyMasterBlue := NewOLllamaSpyMaster(BLUE)

	ollamaFieldAgentRed := NewOLllamaOLlamaFieldAgent(RED)

	var redFieldAgent Player = &HumanFieldAgentTurn{
		team: RED,
		rl:   rl,
	}
	var redSpyMaster Player = &HumanSpyMasterTurn{
		team: RED,
		rl:   rl,
	}
	var blueFieldAgent Player = &HumanFieldAgentTurn{
		team: BLUE,
		rl:   rl,
	}
	var blueSpyMaster Player = &HumanSpyMasterTurn{
		team: BLUE,
		rl:   rl,
	}

	//redSpyMaster = ollamaSpyMasterRed
	redFieldAgent = ollamaFieldAgentRed
	blueSpyMaster = ollamaSpyMasterBlue

	game.transitions = map[string]Player{
		RED + "CLUE":   redFieldAgent,
		RED + "GUESS":  blueSpyMaster,
		BLUE + "CLUE":  blueFieldAgent,
		BLUE + "GUESS": redSpyMaster,
	}
	game.state = redSpyMaster

	for {
		err = game.state.Move(game)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}

		// check for win conditions
		s := game.currentScore()
		winner := ""
		for team, score := range s {
			if len(*game.cards[team]) == score {
				winner = team
			}
		}

		if winner != "" {
			fmt.Printf("Winner is %s\nFinal scores:\n%v\n", winner, s)
			os.Exit(0)
		}
	}
}

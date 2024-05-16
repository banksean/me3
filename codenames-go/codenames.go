package main

import (
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strings"

	"github.com/chzyer/readline"
	"github.com/olekukonko/tablewriter"

	_ "embed"
)

const (
	RED       = "red"
	BLUE      = "blue"
	BYSTANDER = "bystander"
	ASSASSIN  = "assassin"
)

var (
	//go:embed wordlist.txt
	wordlistFile string
)

type gameBoard struct {
	cards        map[string]map[string]any
	guessedWords map[string]any
	teamForCard  map[string]string
	rl           *readline.Instance
	state        gameState
	transitions  map[string]gameState
}

// "State Machine" pattern.
type gameState interface {
	PromptInput() (string, error)
	ProcessInput(string) error
}

var (
	defaultColor = tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor, tablewriter.BgBlackColor}
	teamColor    = map[string]tablewriter.Colors{
		RED:       {tablewriter.Bold, tablewriter.FgRedColor, tablewriter.BgBlackColor},
		BLUE:      {tablewriter.Bold, tablewriter.FgBlueColor, tablewriter.BgBlackColor},
		BYSTANDER: {tablewriter.Bold, tablewriter.FgWhiteColor, tablewriter.BgBlackColor},
		ASSASSIN:  {tablewriter.Bold, tablewriter.FgYellowColor, tablewriter.BgBlackColor},
	}
)

func (g *gameBoard) WriteTable(w io.Writer, spyMasterView bool) {
	tw := tablewriter.NewWriter(w)
	tw.SetBorder(false)
	tw.SetBorders(tablewriter.Border{})
	tw.SetCenterSeparator(" ")
	tw.SetColumnSeparator(" ")
	tw.SetRowSeparator(" ")
	tw.SetRowLine(true)
	allCards := []string{}
	for _, cards := range g.cards {
		for c := range cards {
			allCards = append(allCards, c)
		}
	}
	sort.Strings(allCards)
	allColors := []tablewriter.Colors{}
	for i, c := range allCards {
		if _, guessed := g.guessedWords[c]; guessed || spyMasterView {
			team := g.teamForCard[c]
			allColors = append(allColors, teamColor[team])
			if guessed {
				allCards[i] = strings.ToUpper(c)
			}
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

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return r, false
	}
	return r, true
}

func NewGameBoard(d []string) *gameBoard {
	game := &gameBoard{
		cards:        make(map[string]map[string]any),
		guessedWords: make(map[string]any),
		teamForCard:  make(map[string]string),
	}
	game.cards[RED], d = draw(d, 8)
	game.cards[BLUE], d = draw(d, 9)
	game.cards[BYSTANDER], d = draw(d, 7)
	game.cards[ASSASSIN], _ = draw(d, 1)
	//	pcItems := []readline.PrefixCompleterInterface{}
	for team, cards := range game.cards {
		for card := range cards {
			game.teamForCard[card] = team
			//		pcItems = append(pcItems, readline.PcItem(card))
		}
	}

	ollamaSpyMasterRed := NewOLllamaSpyMaster(game, RED)
	ollamaSpyMasterBlue := NewOLllamaSpyMaster(game, BLUE)

	var redFieldAgent gameState = &HumanFieldAgentTurn{
		team: RED,
		game: game,
	}
	var redSpyMaster gameState = &HumanSpyMasterTurn{
		team: RED,
		game: game,
	}
	var blueFieldAgent gameState = &HumanFieldAgentTurn{
		team: BLUE,
		game: game,
	}
	var blueSpyMaster gameState = &HumanSpyMasterTurn{
		team: BLUE,
		game: game,
	}

	redSpyMaster = ollamaSpyMasterRed
	blueSpyMaster = ollamaSpyMasterBlue

	game.transitions = map[string]gameState{
		"REDCLUE":   redFieldAgent,
		"REDGUESS":  blueSpyMaster,
		"BLUECLUE":  blueFieldAgent,
		"BLUEGUESS": redSpyMaster,
	}
	game.state = redSpyMaster
	return game
}

func main() {
	words := strings.Split(wordlistFile, "\n")
	d := deck(words)
	game := NewGameBoard(d)

	rl, err := readline.NewEx(
		&readline.Config{
			InterruptPrompt:     "^C",
			Prompt:              "> ",
			FuncFilterInputRune: filterInput,
			//AutoComplete:        readline.NewPrefixCompleter(pcItems...),
		},
	)

	if err != nil {
		panic(err)
	}
	defer rl.Close()
	game.rl = rl

	for {
		fmt.Printf("score: %v\n", game.currentScore())
		line, err := game.state.PromptInput()
		if err != nil { // io.EOF
			break
		}
		err = game.state.ProcessInput(strings.TrimSpace(line))
		if err != nil {
			fmt.Printf("error: %v\n", err)
		}
	}
}

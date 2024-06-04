package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"strings"

	"bitbucket.org/creachadair/stringset"
	"github.com/chzyer/readline"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	_ "embed"
)

const (
	RED       = "RED"
	BLUE      = "BLUE"
	BYSTANDER = "BYSTANDER"
	ASSASSIN  = "ASSASSIN"

	INTERACTIVE_PLAYER = "interactive"
)

var (
	//go:embed wordlist.txt
	wordlistFile string

	blueSpyMasterType  = flag.String("blue-spy-master", INTERACTIVE_PLAYER, "type of player for the blue spy master role")
	blueFieldAgentType = flag.String("blue-field-agent", INTERACTIVE_PLAYER, "type of player for the blue field agent role")
	redSpyMasterType   = flag.String("red-spy-master", INTERACTIVE_PLAYER, "type of player for the red spy master role")
	redFieldAgentType  = flag.String("red-field-agent", INTERACTIVE_PLAYER, "type of player for the red field agent role")
)

type gameBoard struct {
	cards        map[string]*stringset.Set
	guessedWords *stringset.Set
	teamForCard  map[string]string
	state        Player
	transitions  map[string]Player
	clue         map[string]string
	explanation  map[string]string
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
	Team() string
}

const GUESSED = "guessed-"

var (
	defaultColor = text.Colors{text.FgWhite, text.BgBlack}
	teamColor    = map[string]text.Colors{
		RED:                 {text.FgRed, text.BgBlack},
		BLUE:                {text.FgBlue, text.BgBlack},
		BYSTANDER:           {text.FgWhite, text.BgBlack},
		ASSASSIN:            {text.FgYellow, text.BgBlack},
		RED + GUESSED:       {text.FgBlack, text.BgRed},
		BLUE + GUESSED:      {text.FgBlack, text.BgBlue},
		BYSTANDER + GUESSED: {text.FgBlack, text.BgWhite},
		ASSASSIN + GUESSED:  {text.FgBlack, text.BgYellow},
	}
)

func (g *gameBoard) WriteTable(w io.Writer, spyMasterView bool) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	allCards := []string{}
	for _, cards := range g.cards {
		for c := range *cards {
			allCards = append(allCards, c)
		}
	}
	sort.Strings(allCards)
	for i, c := range allCards {
		guessed := g.guessedWords.Contains(c)
		if guessed || spyMasterView {
			if guessed {
				allCards[i] = strings.ToUpper(c)
			}
		}
	}
	cardTransformer := text.Transformer(func(c interface{}) string {
		card := c.(string)
		team := g.teamForCard[card]
		guessed := g.guessedWords.Contains(card)
		if guessed || spyMasterView {
			return teamColor[team].Sprint(card)
		}
		return defaultColor.Sprint(card)
	})

	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateHeader = true
	t.Style().Options.SeparateRows = true
	t.Style().Format.Header = text.FormatDefault

	header := table.Row{}
	columnConfigs := []table.ColumnConfig{}
	for i := 0; i < 5; i++ {
		card := allCards[i]
		header = append(header, card)
		columnConfigs = append(columnConfigs, table.ColumnConfig{
			Name:              card,
			Transformer:       cardTransformer,
			TransformerHeader: cardTransformer,
			Align:             text.AlignCenter,
			AlignHeader:       text.AlignCenter,
		})
	}
	t.AppendHeader(header)
	t.SetColumnConfigs(columnConfigs)

	for row := 1; row < 5; row++ {
		rowCards := allCards[row*5 : row*5+5]
		row := table.Row{}
		for _, card := range rowCards {
			row = append(row, card)
		}
		t.AppendRow(row)
	}
	t.Render()
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
		explanation:  map[string]string{},
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
	flag.Parse()

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

	var redFieldAgent, redSpyMaster, blueFieldAgent, blueSpyMaster Player

	if *redFieldAgentType == INTERACTIVE_PLAYER {
		redFieldAgent = &HumanFieldAgentTurn{
			team: RED,
			rl:   rl,
		}
	} else {
		redFieldAgent, err = NewOLlamaFieldAgent(RED, *redFieldAgentType)
		if err != nil {
			panic(err.Error())
		}
	}

	if *redSpyMasterType == INTERACTIVE_PLAYER {
		redSpyMaster = &HumanSpyMasterTurn{
			team: RED,
			rl:   rl,
		}
	} else {
		redSpyMaster, err = NewOLllamaSpyMaster(RED, *redSpyMasterType)
		if err != nil {
			panic(err.Error())
		}
	}

	if *blueFieldAgentType == INTERACTIVE_PLAYER {
		blueFieldAgent = &HumanFieldAgentTurn{
			team: BLUE,
			rl:   rl,
		}
	} else {
		blueFieldAgent, err = NewOLlamaFieldAgent(BLUE, *blueFieldAgentType)
		if err != nil {
			panic(err.Error())
		}
	}

	if *blueSpyMasterType == INTERACTIVE_PLAYER {
		blueSpyMaster = &HumanSpyMasterTurn{
			team: BLUE,
			rl:   rl,
		}
	} else {
		blueSpyMaster, err = NewOLllamaSpyMaster(BLUE, *blueSpyMasterType)
		if err != nil {
			panic(err.Error())
		}
	}

	game.transitions = map[string]Player{
		RED + "CLUE":   redFieldAgent,
		RED + "GUESS":  blueSpyMaster,
		BLUE + "CLUE":  blueFieldAgent,
		BLUE + "GUESS": redSpyMaster,
	}
	game.state = redSpyMaster

	for {
		// check for win conditions
		s := game.currentScore()
		winner := ""
		for team, score := range s {
			if len(*game.cards[team]) == score {
				winner = team
			}
		}

		if winner == ASSASSIN {
			// The turn will have moved on to the team who DIDN'T pick the assassin word,
			// so they would be the winner.
			winner = game.state.Team()
		}
		if winner == RED || winner == BLUE {
			fmt.Printf("Winner is %s\nFinal scores:\n%v\n", winner, s)
			os.Exit(0)
		}

		err = game.state.Move(game)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			break
		}
	}
}

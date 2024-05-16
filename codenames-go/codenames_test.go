package main

import (
	"fmt"
	"testing"

	"bitbucket.org/creachadair/stringset"
)

func newTestGameBoard() *gameBoard {
	deck := []string{}
	for i := 0; i < 25; i++ {
		deck = append(deck, fmt.Sprintf("word-%d", i))
	}
	guessed := stringset.New()
	game := &gameBoard{
		cards: map[string]*stringset.Set{
			RED:       {},
			BLUE:      {},
			BYSTANDER: {},
			ASSASSIN:  {},
		},
		guessedWords: &guessed,
		teamForCard:  make(map[string]string),
	}

	for _, card := range deck[0:8] {
		game.cards[RED].Add(card)
	}

	for _, card := range deck[8:17] {
		game.cards[BLUE].Add(card)
	}

	for _, card := range deck[17:24] {
		game.cards[BYSTANDER].Add(card)
	}

	for _, card := range deck[24:] {
		game.cards[ASSASSIN].Add(card)
	}

	return game
}

func TestGameBoard(t *testing.T) {
	g := newTestGameBoard()

	for i := 0; i < 25; i++ {
		team, err := g.guess(fmt.Sprintf("word-%d", i))
		if err != nil {
			t.Errorf("expected nil error, got: %s", err)
		}
		if i < 8 {
			if team != RED {
				t.Errorf("%d: expected RED, got %q", i, team)
			}
		} else if i < 17 {
			if team != BLUE {
				t.Errorf("%d: expected BLUE, got %q", i, team)
			}
		} else if i < 24 {
			if team != BYSTANDER {
				t.Errorf("%d: expected BYSTANDER, got %q", i, team)
			}
		} else {
			if team != ASSASSIN {
				t.Errorf("%d: expected ASSASSIN, got %q", i, team)
			}
		}
	}

	s := g.currentScore()
	for k, v := range map[string]int{
		ASSASSIN:  1,
		BYSTANDER: 7,
		RED:       8,
		BLUE:      9,
	} {
		if s[k] != v {
			t.Errorf("expected score for %q to be %d, but was %d instead", k, v, s[k])
		}
	}
}

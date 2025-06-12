package solitaire

import "fmt"

// Versions - Games supported by this engine.
type Versions int

//nolint:revive // Ignore the need to comment on each version.
const (
	Klondike Versions = iota
)

func (solitaireGames Versions) String() string {
	return [...]string{
		"Klondike",
	}[solitaireGames]
}

// Engine of the application.

// This co-ordinates the resources needed for a game, and manages the rules that
// apply for the given game.

// Game - the game being played.
type Game struct {
	Deck        *Deck
	Foundations []Foundation
	Tableau     []*Tableau
}

// NewGame - Create a new Game.
func NewGame(decks, foundations, tableaus int, foundationBase, tableauBase Rank) *Game {
	deck := CreateDecks(decks)
	deck.Shuffle()

	return &Game{
		Deck: deck,
		Foundations: CreateFoundations(
			foundations,
			foundationBase,
			func(
				foundation Foundation,
				card SuitedCard,
			) bool {
				// Card must be in the same suit.
				if card.Suit != foundation.Base.Suit {
					return false
				}

				// If no cards have been added to the foundation then the card must
				// be the same as the base card.
				if foundation.Len() == 0 && card.Rank != foundation.Base.Rank {
					return false
				}

				// Handle solitaire.Ace being put on top of a solitaire.King - can
				// only happen when more than one deck is in the Deck.
				topCard, err := foundation.Top()
				if err != nil {
					return false
				}

				if (card.Rank == Ace) && (topCard.Rank == King) {
					return true
				}

				// Card being added must be rank 1 higher than the latest on the
				// foundation stack.
				if card.Rank == (topCard.Rank + 1) {
					return true
				}

				return false
			}),
		Tableau: CreateTableaus(
			tableaus,
			func(
				tableau *Tableau,
				card SuitedCard,
			) bool {
				topCard, err := tableau.Top()
				if err != nil {
					fmt.Println("Error getting tableau top", err.Error())
					return false
				}

				// Card must be in the opposite coloured suit.
				if card.Suit == topCard.Suit {
					return false
				}

				// If the suit of the card plus the topcard of the tableau add
				// up to an even number, then they are both the same colour.
				// Even plus even == even
				// Odd Plus odd == even
				// Even + Odd == odd.
				if (card.Suit+topCard.Suit)%2 == 0 {
					return false
				}

				// If the tableau is empty, only the tableau base can be put on it.
				if tableau.Len() == 0 && card.Rank != tableauBase {
					return false
				}

				// Handle solitaire.King being put on top of a solitaire.Ace - can
				// only happen when more than one deck is in the Deck.
				if (card.Rank == King) && (topCard.Rank == Ace) {
					return true
				}

				// Card being added must be rank 1 lower than the latest on the
				// foundation stack.
				if card.Rank == (topCard.Rank - 1) {
					return true
				}

				return false
			},
		),
	}
}

// Start - start the game.
func (game Game) Start(name Versions) {
	if name == Klondike {
		game.StartKlondike()
	}
}

// StartKlondike - start a game of Klondike Solitaire.
func (game Game) StartKlondike() {
	// Klondike has 7 stacks in its tableau.
	const numStacks = 7
	game.Tableau = make([]*Tableau, 0, numStacks)

	counter := 0
	for idx := 0; idx <= numStacks; idx++ {
		// The form a triangle with each successive stack holding one more card
		// than the last.
		if idx < counter {
			continue
		}

		visible := false

		// The final card to be added is visible.
		if counter == idx {
			visible = true
		}

		game.Tableau[idx].Add(game.Deck.Deal(), visible)

		counter++
	}
}

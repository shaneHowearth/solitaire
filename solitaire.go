package solitaire

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
	Tableau     Tableau
}

// NewGame - Create a new Game.
func NewGame(decks, foundations int, foundationBase Card) *Game {
	deck := CreateDecks(decks)
	deck.Shuffle()

	return &Game{
		Deck:        deck,
		Foundations: CreateFoundations(foundations, foundationBase),
	}
}

// Start - start the game.
func (game Game) Start(name Versions) {
	if name == Klondike {
		game.StartKlondike()
	}
}

// StartKlondike - start a game of Klondije Solitaire.
func (game Game) StartKlondike() {
	// Klondike has 7 stacks in its tableau.
	const numStacks = 7
	game.Tableau = make(Tableau, 0, numStacks)

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

package game

import "github.com/shanehowearth/solitaire/state"

// Klondike - https://en.wikipedia.org/wiki/Klondike_(solitaire)
type Klondike struct{}

// Ensure that Klondike implements game.Variant
var _ Variant = (*Klondike)(nil)

// Name - name of the variant.
func (*Klondike) Name() string {
	return "Klondike"
}

// TableauGridSize - The size of the grid required by klondike.
func (*Klondike) TableauGridSize() (int, int) {
	const height = 1

	return height, numKlondikeTableau
}

// Decks - How many decks of cards are required to play klondike.
func (*Klondike) Decks() int {
	return 1
}

const numKlondikeTableau = 7

// Tableau - how the tableau are defined.
func (*Klondike) Tableau() (
	number int,
	basecard state.Rank,
	addRule func(state.Tableau, state.SuitedCard) bool,
) {
	return numKlondikeTableau, state.King, MinusOneRule
}

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (*Klondike) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*Klondike) Foundations() (
	number int,
	basecard state.Rank,
	addRule func(state.Foundation, state.SuitedCard) bool,
) {
	const foundationCount = 4
	return foundationCount, state.Ace, PlusOneRule
}

// SetupDeal - deal all the cards into the right piles ready to start a game.
func (*Klondike) SetupDeal() {
	panic("Not Implemented")
}

package game

import "github.com/shanehowearth/solitaire"

// Klondike - https://en.wikipedia.org/wiki/Klondike_(solitaire)
type Klondike struct{}

// Name - name of the variant.
func (Klondike) Name() string {
	return "Klondike"
}

// GridSize - The size of the grid required by klondike.
func (Klondike) GridSize() (int, int) {
	const height = 1

	const width = 7

	return height, width
}

// Decks - How many decks of cards are required to play klondike.
func (Klondike) Decks() int {
	return 1
}

const numKlondikeTableau = 7

// Tableau - how the tableau are defined.
func (Klondike) Tableau() (
	number int,
	basecard solitaire.Rank,
	addRule func(solitaire.Tableau, solitaire.SuitedCard) bool,
) {
	return numKlondikeTableau, solitaire.King, MinusOneRule
}

// Layout - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (Klondike) Layout(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (Klondike) Foundations() (
	number int,
	basecard solitaire.Rank,
	addRule func(solitaire.Foundation, solitaire.SuitedCard) bool,
) {
	const foundationCount = 4
	return foundationCount, solitaire.Ace, PlusOneRule
}

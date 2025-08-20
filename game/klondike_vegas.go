package game

import "github.com/shanehowearth/solitaire/state"

// KlondikeVegas - https://en.wikipedia.org/wiki/Klondike2_(solitaire)
type KlondikeVegas struct{}

// Ensure that KlondikeVegas implements game.Variant.
var _ Variant = (*KlondikeVegas)(nil)

// Name - name of the variant.
func (*KlondikeVegas) Name() string {
	return "Klondike (Vegas style)"
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of KlondikeVegas.
func (*KlondikeVegas) Reserves() (
	number int,
	cardCount [][2]int,
	addRule func(*state.Reserve, state.SuitedCard) bool,
) {
	const reserveCount = 0
	return reserveCount, [][2]int{}, func(*state.Reserve, state.SuitedCard) bool { return false }
}

// TableauGridSize - The size of the grid required by klondikeVegas.
func (*KlondikeVegas) TableauGridSize() (int, int) {
	const height = 1

	return height, numKlondikeVegasTableau
}

// Decks - How many decks of cards are required to play klondikeVegas.
func (*KlondikeVegas) Decks() int {
	return 1
}

const numKlondikeVegasTableau = 7

// Tableau - how the tableau are defined.
func (*KlondikeVegas) Tableau() (
	number int,
	basecard state.Rank,
	addRule func(*state.Tableau, state.SuitedCard) bool,
) {
	return numKlondikeVegasTableau, state.King, MinusOneRule
}

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (*KlondikeVegas) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*KlondikeVegas) Foundations() (
	number int,
	basecard state.Rank,
	addRule func(state.Foundation, state.SuitedCard) bool,
) {
	const foundationCount = 4
	return foundationCount, state.Ace, PlusOneRule
}

// SetupTableauCardCounts - Should return a list of ints, the first int will be
// the number of cards going into the first tableau, the second will be how many
// cards are visible in that tableau. The third and fourth ints will apply to
// the second tableau, etc.
func (*KlondikeVegas) SetupTableauCardCounts() [][2]int {
	//nolint:revive // Ignore the constant complaint.
	return [][2]int{{1, 1}, {2, 1}, {3, 1}, {4, 1}, {5, 1}, {6, 1}, {7, 1}}
}

// HowToPlay - Tell the player how to play the game.
func (*KlondikeVegas) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The four foundations (rectangles in the upper right of the board) are built up by suit from Ace (low in this game) to King, and the tableau piles can be built down by alternate colors.
`,
		`Every face-up card in a partial pile, or a complete pile, can be moved, as a unit, to another tableau pile on the basis of what can be moved onto that other pile. Any empty piles can be filled with a King, or a pile of cards with a King.
`,
		`The aim of the game is to build up four stacks of cards starting with Ace and ending with King all of the same suit, on one of the four foundations, at which time the player would have won.
`,
		`There are no redeals available in this style of Klondike.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*KlondikeVegas) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for idx := range foundations {
		if foundations[idx].Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - This is the maximum number of times that the stock can be
// refreshed from the waste.
func (*KlondikeVegas) MaxRedeals() int {
	// Never allowed to redeal.
	return 0
}

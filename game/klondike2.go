package game

import "github.com/shanehowearth/solitaire/state"

// Klondike2 - https://en.wikipedia.org/wiki/Klondike2_(solitaire)
type Klondike2 struct{}

// Ensure that Klondike2 implements game.Variant.
var _ Variant = (*Klondike2)(nil)

// Name - name of the variant.
func (*Klondike2) Name() string {
	return "Klondike2"
}

// TableauGridSize - The size of the grid required by klondike2.
func (*Klondike2) TableauGridSize() (int, int) {
	const height = 1

	return height, numKlondike2Tableau
}

// Decks - How many decks of cards are required to play klondike2.
func (*Klondike2) Decks() int {
	return 1
}

const numKlondike2Tableau = 7

// Tableau - how the tableau are defined.
func (*Klondike2) Tableau() (
	number int,
	basecard state.Rank,
	addRule func(*state.Tableau, state.SuitedCard) bool,
) {
	return numKlondike2Tableau, state.King, MinusOneRule
}

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (*Klondike2) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*Klondike2) Foundations() (
	number int,
	basecard state.Rank,
	addRule func(state.Foundation, state.SuitedCard) bool,
) {
	const foundationCount = 4
	return foundationCount, state.Ace, PlusOneRule
}

// SetupDealCardCounts - Should return a list of ints, the first int will be the
// number of cards going into the first tableau, the second will be how many
// cards are visible in that tableau. The third and fourth ints will apply
// to the second tableau, etc.
func (*Klondike2) SetupDealCardCounts() []int {
	//nolint:revive // Ignore the constant complaint.
	return []int{1, 1, 2, 1, 3, 1, 4, 1, 5, 1, 6, 1, 7, 1}
}

// HowToPlay - Tell the player how to play the game.
func (*Klondike2) HowToPlay() []string {
	lines := []string{
		`The four foundations (rectangles in the upper right of the board) are
		built up by suit from Ace (low in this game) to King, and the tableau
		piles can be built down by alternate colors.  Every face-up card in a
		partial pile, or a complete pile, can be moved, as a unit, to another
		tableau pile on the basis of its highest card.  Any empty piles can be
		filled with a King, or a pile of cards with a King.  The aim of the game
		is to build up four stacks of cards starting with Ace and ending with
		King all of the same suit, on one of the four foundations, at which time
		the player would have won.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Klondike2) HasWon(_ []*state.Tableau, _ []*state.Foundation) bool {
	return true
}

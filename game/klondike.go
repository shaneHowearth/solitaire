package game

import "github.com/shanehowearth/solitaire/state"

// Klondike - https://en.wikipedia.org/wiki/Klondike_(solitaire)
type Klondike struct{}

// Ensure that Klondike implements game.Variant.
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
	addRule func(*state.Tableau, state.SuitedCard) bool,
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

// SetupDealCardCounts - Should return a list of ints, the first int will be the
// number of cards going into the first tableau, the second will be how many
// cards are visible in that tableau. The third and fourth ints will apply
// to the second tableau, etc.
func (*Klondike) SetupDealCardCounts() []int {
	//nolint:revive // Ignore the constant complaint.
	return []int{1, 1, 2, 1, 3, 1, 4, 1, 5, 1, 6, 1, 7, 1}
}

// HowToPlay - Tell the player how to play the game.
func (*Klondike) HowToPlay() []string {
	lines := []string{
		`The four foundations (rectangles in the upper right of the board) are built up by suit from Ace (low in this game) to King, and the tableau piles can be built down by alternate colors.
`,
		`Every face-up card in a partial pile, or a complete pile, can be moved, as a unit, to another tableau pile on the basis of what can be moved onto that other pile. Any empty piles can be filled with a King, or a pile of cards with a King.
`,
		`The aim of the game is to build up four stacks of cards starting with Ace and ending with King all of the same suit, on one of the four foundations, at which time the player would have won.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Klondike) HasWon(tableaus []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}
	return true
}

// MaxRedeals - how many redeals are allowed
func (*Klondike) MaxRedeals() int {
	// Allow an unlimited number of redeals.
	return -1
}

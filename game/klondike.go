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
	const numKlondikeTableau = 7

	return height, numKlondikeTableau
}

// Decks - How many decks of cards are required to play klondike.
func (*Klondike) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of Klondike.
func (*Klondike) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (*Klondike) Tableau() []state.StackSpec {
	// return numKlondikeTableau, state.King, MinusOneRule
	return []state.StackSpec{
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{2, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{3, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{4, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{5, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{6, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{7, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
	}
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
func (*Klondike) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades},
			AddRule:  PlusOneRule,
		},
	}
}

// HowToPlay - Tell the player how to play the game.
func (*Klondike) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The four foundations (rectangles in the upper right of the board) are built up by suit from Ace (low in this game) to King, and the tableau piles can be built down by alternate colors.
`,
		`Every face-up card in a partial pile, or a complete pile, can be moved, as a unit, to another tableau pile on the basis of what can be moved onto that other pile. Any empty piles can be filled with a King, or a pile of cards with a King.
`,
		`The aim of the game is to build up four stacks of cards starting with Ace and ending with King all of the same suit, on one of the four foundations, at which time the player would have won.`,
		`There are unlimited redeals available in this style of Klondike.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Klondike) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*Klondike) MaxRedeals() int {
	// Allow an unlimited number of redeals.
	return -1
}

// Move -
func (*Klondike) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination)
}

// Compact
func (*Klondike) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Klondike) Talon() bool {
	return true
}

// Redeal
func (*Klondike) Redeal(_ *state.Talon, _ []*state.Tableau) {}

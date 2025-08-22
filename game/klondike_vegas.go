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
func (*KlondikeVegas) Reserves() []state.StackSpec {
	return []state.StackSpec{}
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
func (*KlondikeVegas) Tableau() []state.StackSpec {
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
func (*KlondikeVegas) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*KlondikeVegas) Foundations() []state.StackSpec {
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

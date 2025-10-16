package game

import (
	"fmt"

	"github.com/shanehowearth/solitaire/state"
)

// Yukon - https://en.wikipedia.org/wiki/Yukon_(solitaire)
type Yukon struct{}

// Ensure that Yukon implements game.Variant.
var _ Variant = (*Yukon)(nil)

// Name - name of the variant.
func (*Yukon) Name() string {
	return "Yukon"
}

// TableauGridSize - The size of the grid required by klondike.
func (*Yukon) TableauGridSize() (int, int) {
	const height = 1
	const numYukonTableau = 7

	return height, numYukonTableau
}

// Decks - How many decks of cards are required to play klondike.
func (*Yukon) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of Yukon.
func (*Yukon) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (*Yukon) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{6, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{7, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{8, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{9, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{10, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   MinusOneRule,
			CardCount: [2]int{11, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
	}
}

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (*Yukon) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*Yukon) Foundations() []state.StackSpec {
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
func (*Yukon) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`There are 7 rows of tableau stacks, and 4 foundations that build up in suit. Their lengths are respectively: 1, 6, 7, 8, 9, 10, 11. The topmost 5 cards of each stack are face-up, while the rest of the cards underneath them are face-down. (Second stack onwards is an ascending pyramid of face-down cards, each with 5 face-up cards appended on top).
`,
		`Groups of cards can be moved; the cards below the one to be moved do not need to be in any order, except that the starting and target cards must be built in sequence and in alternate color. For example, a group starting with a Red 3 can be moved on top a Black 4, and the cards below the Red 3 can differ.
`,
		`There are no redeals available in this style of Yukon.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Yukon) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*Yukon) MaxRedeals() int {
	// Allow an zero redeals.
	return 0
}

// Move -
func (*Yukon) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination, false)
}

// Compact
func (*Yukon) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Yukon) Talon() bool {
	return false
}

// Redeal
func (*Yukon) Redeal(_ *state.Talon, _ []*state.Tableau) {}

// FoundationBase
func (*Yukon) FoundationBase() bool {
	return false
}

func (yukon *Yukon) AvailableMoves(tableau []state.Tableau,
	foundations []state.Foundation,
	_ []state.Talon,
	_ []state.Reserve,
) []string {
	hints := []string{}

	// Check tableau to foundation moves
	for foundationIdx := range foundations {
		for sourceIdx := range tableau {
			if hint := checkMove(
				tableau[sourceIdx].Stack,
				foundations[foundationIdx].Stack,
				fmt.Sprintf("Tableau %d", sourceIdx+1),
				fmt.Sprintf("Foundation %d", foundationIdx+1),
				false,
				true,
			); hint != "" {
				hints = append(hints, hint)
			}
		}
	}

	// Check tableau to tableau moves
	for destIdx := range tableau {
		for sourceIdx := range tableau {
			if destIdx == sourceIdx {
				continue
			}
			if hint := checkMove(
				tableau[sourceIdx].Stack,
				tableau[destIdx].Stack,
				fmt.Sprintf("Tableau %d", sourceIdx+1),
				fmt.Sprintf("Tableau %d", destIdx+1),
				false,
				true,
			); hint != "" {
				hints = append(hints, hint)
			}
		}
	}

	return hints
}

package game

import (
	"fmt"

	"github.com/shanehowearth/solitaire/state"
)

// Russian - https://en.wikipedia.org/wiki/Russian_(solitaire)
type Russian struct{}

// Ensure that Russian implements game.Variant.
var _ Variant = (*Russian)(nil)

// Name - name of the variant.
func (*Russian) Name() string {
	return "Russian"
}

// TableauGridSize - The size of the grid required by klondike.
func (*Russian) TableauGridSize() (int, int) {
	const height = 1
	const numRussianTableau = 7

	return height, numRussianTableau
}

// Decks - How many decks of cards are required to play klondike.
func (*Russian) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of Russian.
func (*Russian) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (*Russian) Tableau() []state.StackSpec {
	var RussianMinusOneRule = func(tableau *state.Stack, card state.SuitedCard) bool {
		// Handle when the tableau is empty.
		if (*tableau).Len() == 0 {
			if card.Rank == tableau.Base.Rank {
				return true
			}
		}

		// Get the card currently at the top of the tableau.
		topCard, err := (*tableau).Top()
		if err != nil {
			return false
		}

		// If the card is the opposite colour, and is one down in rank
		// then it can go onto the tableau.
		if (card.Suit == topCard.Suit) && (topCard.Rank-card.Rank) == 1 {
			return true
		}

		// All other cases the card should not be added to the tableau.
		return false
	}

	return []state.StackSpec{
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{6, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{7, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{8, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{9, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{10, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   RussianMinusOneRule,
			CardCount: [2]int{11, 5},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
	}
}

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
// Tableau and Grid are 0 indexed.
func (*Russian) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*Russian) Foundations() []state.StackSpec {
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
func (*Russian) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`There are 7 rows of tableau stacks, and 4 foundations that build up in suit. Their lengths are respectively: 1, 6, 7, 8, 9, 10, 11. The topmost 5 cards of each stack are face-up, while the rest of the cards underneath them are face-down. (Second stack onwards is an ascending pyramid of face-down cards, each with 5 face-up cards appended on top).
`,
		`Groups of cards can be moved; the cards below the one to be moved do not need to be in any order, except that the starting and target cards must be built in sequence and in alternate color. For example, a group starting with a 3♥ can be moved on top a 4♥, and the cards below the 3♥ can differ.
`,
		`There are no redeals available in this style of Russian.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Russian) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*Russian) MaxRedeals() int {
	// Allow an zero redeals.
	return 0
}

// Move -
func (*Russian) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination, false)
}

// Compact
func (*Russian) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Russian) Talon() bool {
	return false
}

// Redeal
func (*Russian) Redeal(_ *state.Talon, _ []*state.Tableau) {}

// FoundationBase
func (*Russian) FoundationBase() bool {
	return false
}

// AvailableMoves - return a list of the available moves.
func (*Russian) AvailableMoves(
	tableau []state.Tableau,
	foundations []state.Foundation,
	talon []state.Talon,
	reserves []state.Reserve,
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
			); hint != "" {
				hints = append(hints, hint)
			}
		}
	}

	return hints
}

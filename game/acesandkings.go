package game

import "github.com/shanehowearth/solitaire/state"

// AcesAndKings - https://en.wikipedia.org/wiki/Aces_and_Kings
type AcesAndKings struct{}

// Ensure that AcesAndKings implements game.Variant.
var _ Variant = (*AcesAndKings)(nil)

// Name - name of the variant.
func (*AcesAndKings) Name() string {
	return "Aces and Kings"
}

// TableauGridSize - The size of the grid required by Aces and Kings.
func (*AcesAndKings) TableauGridSize() (int, int) {
	const height = 1
	const numAcesAndKingsTableau = 4

	return height, numAcesAndKingsTableau
}

// Decks - How many decks of cards are required to play Aces and Kings.
func (*AcesAndKings) Decks() int {
	return 2
}

// Reserves - how the reserves are defined.
// Note that there are two reserves required in a game of Aces and Kings.
func (*AcesAndKings) Reserves() []state.StackSpec {
	return []state.StackSpec{
		{
			AddRule: func(*state.Stack, state.SuitedCard) bool {
				// Nothing can be added to a reserve.
				return false
			},
			CardCount: [2]int{13, 13},
		},
		{
			AddRule: func(*state.Stack, state.SuitedCard) bool {
				// Nothing can be added to a reserve.
				return false
			},
			CardCount: [2]int{13, 13},
		},
	}
}

// Tableau - how the tableau are defined.
func (*AcesAndKings) Tableau() []state.StackSpec {

	return []state.StackSpec{
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
	}
}

// Foundations - how the foundations are defined.
func (*AcesAndKings) Foundations() []state.StackSpec {
	const foundationCount = 8
	foundations := make([]state.StackSpec, 0, foundationCount)
	// Count, state.Ace, PlusOneRule
	for idx := 0; idx < foundationCount; idx++ {
		// Aces and Kings has eight foundations in total. Four foundations start with an Ace and build up regardless of suit, e.g. A♥, 2♠, 3♦, 4♦.

		// The other four start with a King and build down regardless of suit, e.g. K♣, Q♥, J♠, 10♣
		foundationSpec := state.StackSpec{}
		if idx < 4 {
			foundationSpec.BaseCard = state.SuitedCard{Rank: state.Ace, Suit: state.Undefined}
			foundationSpec.AddRule = func(stack *state.Stack, card state.SuitedCard) bool {
				if stack.Len() == 0 && card.Rank == state.Ace {
					return true
				}
				// Once a stack has ace through King or King through Ace, it's
				// full.
				if stack.Len() == 13 {
					return false
				}
				// Any suit can be added to a foundation, as long as the card is
				// one higher than what's already there.
				topcard, err := stack.Top()
				if err != nil {
					return false
				}

				// Card must be one HIGHER than existing top card.
				if card.Rank-topcard.Rank == 1 {
					return true
				}

				// All other possibilities are excluded.
				return false
			}
		} else {
			foundationSpec.BaseCard = state.SuitedCard{Rank: state.King, Suit: state.Undefined}
			foundationSpec.AddRule = func(stack *state.Stack, card state.SuitedCard) bool {
				if stack.Len() == 0 && card.Rank == state.King {
					return true
				}
				// Any suit can be added to a foundation, as long as the card is
				// one higher than what's already there.
				topcard, err := stack.Top()
				if err != nil {
					return false
				}

				// Card must be one LOWER than existing top card.
				if topcard.Rank-card.Rank == 1 {
					return true
				}

				// All other possibilities are excluded.
				return false
			}
		}

		foundations = append(foundations, foundationSpec)
	}

	return foundations
}

// HowToPlay - Tell the player how to play the game.
func (*AcesAndKings) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`Aces and Kings has eight foundations in total. Four foundations start with an Ace and build up regardless of suit, e.g. A♥, 2♠, 3♦, 4♦.

The other four start with a King and build down regardless of suit, e.g. K♣, Q♥, J♠, 10♣

The games include two reserve piles of thirteen cards each that can only be played onto the foundations. At the bottom left, there is a stock pile and a waste pile. At the bottom right, there are the four tableau piles of one face-up card each.
`,
		`Cards in the tableau can only be moved to the foundations; empty tableau spaces should be immediately filled with the top card from the stock pile. If desired, the top card on a foundation pile can be transferred onto another one, as long as it is of compatible rank. When the stock pile is empty, the game is over: if all of the foundations can be completed, the game has been won.
`,
		`There are no redeals available in this style of Solitaire.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*AcesAndKings) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*AcesAndKings) MaxRedeals() int {
	// Allow an unlimited number of redeals.
	return 0
}

// Move -
func (*AcesAndKings) Move(source, destination *state.Stack, maxRedeals int) bool {
	return Move(source, destination, maxRedeals)
}

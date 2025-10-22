package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// https://en.wikipedia.org/wiki/Agnes_(card_game)
type Agnes struct{}

// Ensure that Agnes implements game.Variant.
var _ Variant = (*Agnes)(nil)

// Name - name of the variant.
func (*Agnes) Name() string {
	return "Agnes"
}

// TableauGridSize - The size of the grid required by Agnes.
func (*Agnes) TableauGridSize() (int, int) {
	const height = 1
	const numAgnesTableau = 7

	return height, numAgnesTableau
}

// Decks - How many decks of cards are required to play Agnes.
func (*Agnes) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of Agnes.
func (*Agnes) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (*Agnes) Tableau() []state.StackSpec {
	var agnesRule = func(tableau *state.Stack, card state.SuitedCard) bool {
		// Handle when the tableau is empty.
		// Nothing can be added to an empty tableau except when dealing.
		if tableau.Len() == 0 {
			return false
		}

		// Get the card currently at the top of the tableau.
		topCard, err := tableau.Top()
		if err != nil {
			return false
		}

		// Allow Kings to go onto Aces (of the same colour).
		if ((card.Suit+topCard.Suit)%2 == 0) && (topCard.Rank == state.Ace && card.Rank == state.King) {
			return true
		}

		// If the card is the same colour, and is one down in rank
		// then it can go onto the tableau.
		if ((card.Suit+topCard.Suit)%2 == 0) && (topCard.Rank-card.Rank) == 1 {
			return true
		}

		// All other cases the card should not be added to the tableau.
		return false
	}

	return []state.StackSpec{
		{
			AddRule:   agnesRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{2, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{3, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{4, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{5, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{6, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   agnesRule,
			CardCount: [2]int{7, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
	}
}

// Foundations - how the foundations are defined.
func (*Agnes) Foundations() []state.StackSpec {
	var AgnesPlusOneRule = func(foundation *state.Stack, card state.SuitedCard) bool {
		// Handle when the foundation is empty.
		if foundation.Len() == 0 {
			if card.Suit == foundation.Base.Suit && card.Rank == foundation.Base.Rank {
				return true
			}
		}

		// Get the card currently at the top of the foundation.
		topCard, err := foundation.Top()
		if err != nil {
			return false
		}

		// Allow Aces to go onto Kings (of the same suit).
		if (card.Suit == topCard.Suit) && (topCard.Rank == state.King && card.Rank == state.Ace) {
			return true
		}

		// If the card is the same suit, and is one up in rank
		// then it can go onto the foundation.
		if card.Suit == foundation.Base.Suit && (card.Rank-topCard.Rank) == 1 {
			return true
		}

		// All other cases the card should not be added to the foundation.
		return false
	}

	return []state.StackSpec{
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts},
			AddRule:  AgnesPlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds},
			AddRule:  AgnesPlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs},
			AddRule:  AgnesPlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades},
			AddRule:  AgnesPlusOneRule,
		},
	}
}

// HasWon - How to tell if the game has been won.
func (*Agnes) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*Agnes) MaxRedeals() int {
	// There are 3 redeals allowed
	return 3
}

// Move -
func (*Agnes) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination, true)
}

// Compact
func (*Agnes) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Agnes) Talon() bool {
	return false
}

// HowToPlay - Tell the player how to play the game.
func (*Agnes) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`After shuffling, twenty-eight cards are dealt, face up, in seven rows and aligned left.[d] The first row has seven cards, the second six cards and so on. The tableau thus forms a right-angled triangle. The twenty-ninth card is turned up and placed above the tableau as the first "master card" (foundation card). The remaining three master cards are those of the same rank as the first and are placed in a row to its right when they become available.[1]
`,
		`The aim is to build on the master cards in suit and ascending sequence, turning the corner with the Aces if necessary. The lowest card in each column is "exposed" i.e. available to be moved to a foundation pile or onto the bottom card of the same colour in another column and in descending sequence. Two or more cards may be moved from one column to another as a packet if they are of the same suit as well as the same colour.[e] So the ♠8 can be moved onto the ♣9, but they cannot be moved on together. In addition, any exposed card can be moved into a vacancy in the tableau, but such spaces need not be occupied.[1][f]
`,
		`When no more moves are possible on the initial layout, seven more cards are dealt to the bottom of the seven columns, after which any further moves may be carried out. Once all desired moves have been made, another seven cards are dealt to the columns and so on until the pack is exhausted. After the third deal, there will be two cards left in hand which may be looked at before the third deal is played. When all moves have been made, they are dealt to the first two columns.[g] Only one deal is permitted and if any cards remain in the tableau, the patience has failed.[1]
`,
	}

	return lines
}

// Redeal
func (*Agnes) Redeal(talon *state.Talon, tableau []*state.Tableau) {
	// Place a card on each of the tableau.
	for idx := range tableau {
		top, err := talon.Stock.Top()
		if err != nil {
			// No more cards.
			return
		}

		// Temporarily set the tableau rule to allow any cards to be added.
		holdRule := tableau[idx].Stack.Rule
		tableau[idx].Stack.Rule = func(state.SuitedCard) bool { return true }

		tableau[idx].Stack.Add(top, true)

		_, _ = talon.Stock.Deal()
		tableau[idx].Stack.Rule = holdRule
	}

}

// FoundationBase
func (*Agnes) FoundationBase() bool {
	return true
}

// AvailableMoves - return a list of the available moves.
func (*Agnes) AvailableMoves(
	tableau []state.Tableau,
	foundations []state.Foundation,
	talon []state.Talon,
	reserves []state.Reserve,
) []state.Move {
	moves := []state.Move{}

	// Check tableau to foundation moves
	for foundationIdx := range foundations {
		for sourceIdx := range tableau {
			if move := checkMove(
				tableau[sourceIdx].Stack,
				foundations[foundationIdx].Stack,
				true,
				true,
			); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
	}

	// Check tableau to tableau moves
	for destIdx := range tableau {
		for sourceIdx := range tableau {
			if destIdx == sourceIdx {
				continue
			}
			if move := checkMove(
				tableau[sourceIdx].Stack,
				tableau[destIdx].Stack,
				true,
				true,
			); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
	}

	return moves
}

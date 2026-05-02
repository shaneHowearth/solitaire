package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type WestcliffAmerican struct{}

var _ Variant = (*WestcliffAmerican)(nil)

func (*WestcliffAmerican) Name() string { return "Westcliff (American)" }

func (*WestcliffAmerican) Category() Category {
	return CatKlondike
}

func (*WestcliffAmerican) Description() string {
	return "A broader version of Klondike with ten tableau columns, offering more opportunities for movement and sequence building."
}

func (*WestcliffAmerican) Decks() int                  { return 1 }
func (*WestcliffAmerican) TableauGridSize() (int, int) { return 1, 10 }
func (*WestcliffAmerican) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (w *WestcliffAmerican) Tableau() []state.StackSpec {
	specs := make([]state.StackSpec, 10)
	for i := 0; i < 10; i++ {
		specs[i] = state.StackSpec{AddRule: w.buildRule, CardCount: [2]int{2, 1}}
	}
	return specs
}

func (*WestcliffAmerican) Fanned() bool { return true }

func (*WestcliffAmerican) buildRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return true
	}
	top, _ := s.Top()
	isTopRed := top.Suit == state.Hearts || top.Suit == state.Diamonds
	isCRed := c.Suit == state.Hearts || c.Suit == state.Diamonds
	return isTopRed != isCRed && int(top.Rank)-int(c.Rank) == 1
}

func (*WestcliffAmerican) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (w *WestcliffAmerican) Move(s, d *state.Stack, _ []*state.Tableau) bool {
	return Move(s, d, true)
}

func (w *WestcliffAmerican) AvailableMoves(tableau []state.Tableau, foundations []state.Foundation, talon []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}

	// 1. Check Waste
	if card, err := talon[0].Waste.Top(); err == nil {
		// To Foundations
		for fIdx := range foundations {
			if foundations[fIdx].Stack.Rule(card) {
				moves = append(moves, state.Move{Source: *talon[0].Waste, Destination: *foundations[fIdx].Stack, NumberMoving: 1, SourceCardTop: card})
			}
		}
		// To Tableau
		for tIdx := range tableau {
			if w.buildRule(tableau[tIdx].Stack, card) {
				moves = append(moves, state.Move{Source: *talon[0].Waste, Destination: *tableau[tIdx].Stack, NumberMoving: 1, SourceCardTop: card})
			}
		}
	}

	// 2. Check Tableau
	for srcIdx := range tableau {
		src := tableau[srcIdx].Stack
		if src.Len() == 0 {
			continue
		}

		// Move whole sequences or single cards
		clone := src.Clone()
		var sequence []state.SuitedCard
		for clone.Len() > 0 {
			card, _ := clone.Top()
			if len(sequence) > 0 {
				last := sequence[len(sequence)-1]
				// Check if this card can be part of an alternating color sequence
				isRed := card.Suit == state.Hearts || card.Suit == state.Diamonds
				isLastRed := last.Suit == state.Hearts || last.Suit == state.Diamonds
				if isRed == isLastRed || int(card.Rank)-int(last.Rank) != 1 {
					break // Sequence broken
				}
			}
			sequence = append(sequence, card)
			_, _ = clone.Deal()

			// Check if this sub-sequence can move to another tableau
			for destIdx := range tableau {
				if srcIdx == destIdx {
					continue
				}
				dest := tableau[destIdx].Stack
				if w.buildRule(dest, card) {
					moves = append(moves, state.Move{Source: *src, Destination: *dest, NumberMoving: len(sequence), SourceCardTop: card})
				}
			}
		}

		// Top card to Foundation
		if card, err := src.Top(); err == nil {
			for fIdx := range foundations {
				if foundations[fIdx].Stack.Rule(card) {
					moves = append(moves, state.Move{Source: *src, Destination: *foundations[fIdx].Stack, NumberMoving: 1, SourceCardTop: card})
				}
			}
		}
	}
	return moves
}

func (*WestcliffAmerican) HowToPlay() []string {
	return []string{"10 columns of 3. Build down by alternating color. Any card to empty space."}
}

func (*WestcliffAmerican) HasWon(_ []*state.Tableau, f []*state.Foundation) bool {
	for _, foundation := range f {
		if foundation.Stack.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*WestcliffAmerican) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*WestcliffAmerican) Talon() bool                                   { return true }
func (*WestcliffAmerican) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*WestcliffAmerican) FoundationBase() bool                          { return false }
func (*WestcliffAmerican) MaxRedeals() int                               { return 0 }

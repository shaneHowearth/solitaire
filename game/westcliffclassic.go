package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type WestcliffClassic struct{}

var _ Variant = (*WestcliffClassic)(nil)

func (*WestcliffClassic) Name() string { return "Westcliff (Classic)" }

func (*WestcliffClassic) Category() Category {
	return CatKlondike
}

func (*WestcliffClassic) Description() string {
	return "The traditional Westcliff ruleset. Deal from the stock one by one to navigate through the deck and clear the tableau."
}

func (*WestcliffClassic) Decks() int                  { return 1 }
func (*WestcliffClassic) TableauGridSize() (int, int) { return 1, 7 }
func (*WestcliffClassic) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (w *WestcliffClassic) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
		{AddRule: w.buildRule, CardCount: [2]int{2, 1}},
	}
}

func (*WestcliffClassic) Fanned() bool { return true }

func (*WestcliffClassic) buildRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return true
	}
	top, _ := s.Top()
	isTopRed := top.Suit == state.Hearts || top.Suit == state.Diamonds
	isCRed := c.Suit == state.Hearts || c.Suit == state.Diamonds
	return isTopRed != isCRed && int(top.Rank)-int(c.Rank) == 1
}

func (*WestcliffClassic) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (w *WestcliffClassic) Move(s, d *state.Stack, _ []*state.Tableau, _ []*state.Reserve) bool {
	return Move(s, d, true)
}

func (w *WestcliffClassic) AvailableMoves(tableau []*state.Tableau, foundations []*state.Foundation, talon []*state.Talon, _ []*state.Reserve) []state.Move {
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

func (*WestcliffClassic) HowToPlay() []string {
	return []string{"Build down by alternating color. Spaces can be filled by any card."}
}

func (*WestcliffClassic) HasWon(_ []*state.Tableau, f []*state.Foundation) bool {
	for _, foundation := range f {
		if foundation.Stack.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*WestcliffClassic) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*WestcliffClassic) Talon() bool                                   { return true }
func (*WestcliffClassic) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*WestcliffClassic) FoundationBase() bool                          { return false }
func (*WestcliffClassic) MaxRedeals() int                               { return 0 }

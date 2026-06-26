package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Easthaven struct{}

var _ Variant = (*Easthaven)(nil)

func (*Easthaven) Name() string { return "Easthaven" }

func (*Easthaven) Category() Category {
	return CatKlondike
}

func (*Easthaven) Description() string {
	return "A blend of Klondike and Spider. Dealing from the stock adds one card to each tableau pile, creating a unique challenge in uncovering hidden cards."
}

func (*Easthaven) Decks() int                  { return 1 }
func (*Easthaven) TableauGridSize() (int, int) { return 1, 7 }
func (*Easthaven) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (e *Easthaven) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
		{AddRule: e.easthavenRule, CardCount: [2]int{2, 1}},
	}
}

func (*Easthaven) Fanned() bool { return true }

func (e *Easthaven) easthavenRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return c.Rank == state.King
	}
	top, _ := s.Top()
	isTopRed := top.Suit == state.Hearts || top.Suit == state.Diamonds
	isCRed := c.Suit == state.Hearts || c.Suit == state.Diamonds
	return isTopRed != isCRed && int(top.Rank)-int(c.Rank) == 1
}

func (*Easthaven) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (e *Easthaven) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination, false)
}

func (e *Easthaven) AvailableMoves(tableau []state.Tableau, foundations []state.Foundation, talon []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}

	// 1. Check Waste
	if card, err := talon[0].Waste.Top(); err == nil {
		for fIdx := range foundations {
			if foundations[fIdx].Stack.Rule(card) {
				moves = append(moves, state.Move{Source: *talon[0].Waste, Destination: *foundations[fIdx].Stack, NumberMoving: 1, SourceCardTop: card})
			}
		}
		for tIdx := range tableau {
			if e.easthavenRule(tableau[tIdx].Stack, card) {
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

		clone := src.Clone()
		var sequence []state.SuitedCard
		for clone.Len() > 0 {
			card, _ := clone.Top()
			if len(sequence) > 0 {
				last := sequence[len(sequence)-1]
				isRed := card.Suit == state.Hearts || card.Suit == state.Diamonds
				isLastRed := last.Suit == state.Hearts || last.Suit == state.Diamonds
				if isRed == isLastRed || int(card.Rank)-int(last.Rank) != 1 {
					break
				}
			}
			sequence = append(sequence, card)
			_, _ = clone.Deal()

			for destIdx := range tableau {
				if srcIdx == destIdx {
					continue
				}
				dest := tableau[destIdx].Stack
				// Use e.easthavenRule to ensure empty spaces only get Kings
				if e.easthavenRule(dest, card) {
					moves = append(moves, state.Move{Source: *src, Destination: *dest, NumberMoving: len(sequence), SourceCardTop: card})
				}
			}
		}

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

func (e *Easthaven) Redeal(talon *state.Talon, tableaus []*state.Tableau) {
	e.Move(talon.Waste, talon.Stock, tableaus)
}

func (*Easthaven) HowToPlay() []string {
	return []string{"Build down by alternating color. Only Kings can fill empty spaces."}
}

func (*Easthaven) HasWon(_ []*state.Tableau, f []*state.Foundation) bool {
	for _, foundation := range f {
		if foundation.Stack.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*Easthaven) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Easthaven) Talon() bool                                   { return true }
func (*Easthaven) FoundationBase() bool                          { return false }
func (*Easthaven) MaxRedeals() int                               { return -1 }

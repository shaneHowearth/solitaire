package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Whitehead - https://en.wikipedia.org/wiki/Whitehead_(solitaire).
type Whitehead struct{}

var _ Variant = (*Whitehead)(nil)

func (*Whitehead) Name() string {
	return "Whitehead"
}

func (*Whitehead) Category() Category {
	return CatKlondike
}

func (*Whitehead) Description() string {
	return "A Klondike variant where all cards are dealt face up. You build down by suit color (Red on Red, Black on Black) and can move sequences of the same suit."
}

func (*Whitehead) TableauGridSize() (int, int) {
	return 1, 7
}

func (*Whitehead) Decks() int {
	return 1
}

func (*Whitehead) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

func (w *Whitehead) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 1}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 2}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 3}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 4}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 5}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 6}},
		{AddRule: w.sameColorRule, CardCount: [2]int{0, 7}},
	}
}

func (*Whitehead) sameColorRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return true
	}
	top, _ := s.Top()
	isTopRed := top.Suit == state.Hearts || top.Suit == state.Diamonds
	isCRed := c.Suit == state.Hearts || c.Suit == state.Diamonds
	if isTopRed != isCRed {
		return false
	}
	return int(top.Rank)-int(c.Rank) == 1
}

func (*Whitehead) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (*Whitehead) HowToPlay() []string {
	return []string{
		`Tableau: All cards face up. Build down by same color.`,
		`Foundations: Build up by suit from Ace to King.`,
		`Moving: Only sequences of the SAME SUIT can be moved as a unit.`,
		`Empty Spaces: Can be filled by any card or valid sequence.`,
	}
}

func (*Whitehead) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*Whitehead) MaxRedeals() int {
	return -1
}

func (w *Whitehead) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	if destination == nil || source == nil || source.Len() == 0 {
		return false
	}

	// 1. SYSTEM MOVES: Stock to Waste OR Waste to Stock (Redeal)
	// These moves ignore rank/color rules and just transfer the cards.
	if destination.Type == state.StackWaste || destination.Type == state.StackTalon {
		// When dealing from Stock to Waste, move 1 card.
		// When redealing from Waste to Stock, move all cards.
		count := 1
		if destination.Type == state.StackTalon {
			count = source.Len()
		}
		return w.performManualMove(source, destination, count)
	}

	// 2. FOUNDATION: Single card, must follow suit/rank rules.
	if destination.Type == state.StackFoundation {
		top, _ := source.Top()
		if destination.Rule(top) {
			return w.performManualMove(source, destination, 1)
		}
		return false
	}

	// 3. TABLEAU: Whitehead "Build by Color, Move by Suit" Logic.
	clone := source.Clone()
	var currentPile []state.SuitedCard

	for clone.Len() > 0 {
		card, _ := clone.Top()

		if len(currentPile) > 0 {
			above := currentPile[len(currentPile)-1]
			// Only same-suit sequences can be picked up as a pile.
			if card.Suit != above.Suit || int(card.Rank)-int(above.Rank) != 1 {
				break
			}
		}

		currentPile = append(currentPile, card)
		_, _ = clone.Deal()

		// Does this specific sub-pile's bottom card fit the destination?
		if w.sameColorRule(destination, card) {
			return w.performManualMove(source, destination, len(currentPile))
		}
	}

	return false
}

func (w *Whitehead) performManualMove(source, destination *state.Stack, count int) bool {
	temp := state.NewStack(count, state.SuitedCard{}, nil, state.StackUndefined)
	for i := 0; i < count; i++ {
		c, _ := source.Top()
		// If we are moving TO the Talon, make it invisible.
		// Otherwise, keep it visible.
		temp.Add(c, true)
		_, _ = source.Deal()
	}

	for temp.Len() > 0 {
		c, _ := temp.Top()

		isVisible := true
		if destination.Type == state.StackTalon {
			isVisible = false
		}

		destination.Add(c, isVisible)
		_, _ = temp.Deal()
	}

	// Only auto-flip the source if it's a Tableau/Reserve.
	// Stock/Waste handle their own visibility.
	if source.Len() > 0 && (source.Type == state.StackTableau || source.Type == state.StackReserve) {
		t, _ := source.Top()
		_, _ = source.Deal()
		source.Add(t, true)
	}

	return true
}

func (w *Whitehead) AvailableMoves(
	tableau []state.Tableau,
	foundations []state.Foundation,
	talon []state.Talon,
	_ []state.Reserve,
) []state.Move {
	moves := []state.Move{}

	// Waste
	if card, err := talon[0].Waste.Top(); err == nil {
		for fIdx := range foundations {
			if foundations[fIdx].Stack.Rule(card) {
				moves = append(moves, state.Move{Source: *talon[0].Waste, Destination: *foundations[fIdx].Stack, NumberMoving: 1, SourceCardTop: card})
			}
		}
		for tIdx := range tableau {
			if w.sameColorRule(tableau[tIdx].Stack, card) {
				moves = append(moves, state.Move{Source: *talon[0].Waste, Destination: *tableau[tIdx].Stack, NumberMoving: 1, SourceCardTop: card})
			}
		}
	}

	// Tableau
	for srcIdx := range tableau {
		src := tableau[srcIdx].Stack
		if src.Len() == 0 {
			continue
		}

		clone := src.Clone()
		var seq []state.SuitedCard
		for clone.Len() > 0 {
			c, _ := clone.Top()
			if len(seq) > 0 {
				prev := seq[len(seq)-1]
				if c.Suit != prev.Suit || int(c.Rank)-int(prev.Rank) != 1 {
					break
				}
			}
			seq = append(seq, c)
			_, _ = clone.Deal()

			for destIdx := range tableau {
				if srcIdx == destIdx {
					continue
				}
				dest := tableau[destIdx].Stack
				if w.sameColorRule(dest, c) {
					moves = append(moves, state.Move{
						Source:        *src,
						Destination:   *dest,
						NumberMoving:  len(seq),
						SourceCardTop: c,
					})
				}
			}
		}

		// To Foundation
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

func (*Whitehead) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Whitehead) Talon() bool                                   { return true }
func (*Whitehead) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*Whitehead) FoundationBase() bool                          { return false }

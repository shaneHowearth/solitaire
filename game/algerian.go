package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Algerian - https://en.wikipedia.org/wiki/Algerian_(card_game).
type Algerian struct{}

var _ Variant = (*Algerian)(nil)

func (*Algerian) Name() string {
	return "Algerian"
}

func (*Algerian) TableauGridSize() (int, int) {
	return 1, 8
}

func (*Algerian) Decks() int {
	return 2
}

func (*Algerian) Reserves() []state.StackSpec {
	reserves := make([]state.StackSpec, 6)
	for i := range reserves {
		reserves[i] = state.StackSpec{
			CardCount: [2]int{4, 4},
			AddRule:   func(_ *state.Stack, _ state.SuitedCard) bool { return false },
		}
	}
	return reserves
}

func (a *Algerian) Tableau() []state.StackSpec {
	tableaus := make([]state.StackSpec, 8)
	for i := range tableaus {
		tableaus[i] = state.StackSpec{
			CardCount: [2]int{1, 1},
			AddRule:   a.tableauRule,
		}
	}
	return tableaus
}

func (*Algerian) tableauRule(tableau *state.Stack, card state.SuitedCard) bool {
	if tableau.Len() == 0 {
		return true
	}
	top, _ := tableau.Top()
	if top.Suit != card.Suit {
		return false
	}

	diff := int(card.Rank) - int(top.Rank)
	// Supports bidirectional building and wrapping (K-A or A-K)
	// Ace is 0, King is 12.
	// diff 1 or -1 is adjacent.
	// diff 12 is K (12) - A (0). diff -12 is A (0) - K (12).
	return diff == 1 || diff == -1 || diff == 12 || diff == -12
}

func (a *Algerian) Foundations() []state.StackSpec {
	foundations := make([]state.StackSpec, 8)
	for i := 0; i < 8; i++ {
		spec := state.StackSpec{}
		if i < 4 {
			spec.BaseCard = state.SuitedCard{Rank: state.Ace, Suit: state.Undefined}
			spec.AddRule = a.foundationUpRule
		} else {
			spec.BaseCard = state.SuitedCard{Rank: state.King, Suit: state.Undefined}
			spec.AddRule = a.foundationDownRule
		}
		foundations[i] = spec
	}
	return foundations
}

func (*Algerian) foundationUpRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return c.Rank == state.Ace
	}
	top, _ := s.Top()
	return c.Suit == top.Suit && int(c.Rank)-int(top.Rank) == 1
}

func (*Algerian) foundationDownRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return c.Rank == state.King
	}
	top, _ := s.Top()
	return c.Suit == top.Suit && int(top.Rank)-int(c.Rank) == 1
}

func (*Algerian) HowToPlay() []string {
	return []string{
		`Eight foundations: 4 Ace-up, 4 King-down. Build in suit.`,
		`Tableau builds up or down in suit, wrapping allowed.`,
		`Cards can transfer between foundations of the same suit if ranks meet.`,
		`Only one card can be moved at a time.`,
	}
}

func (a *Algerian) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	card, err := source.Top()
	if err != nil {
		return false
	}

	// Logic for Foundation to Foundation transfer
	if source.Type == state.StackFoundation && destination.Type == state.StackFoundation {
		topDest, errDest := destination.Top()
		if errDest == nil && card.Suit == topDest.Suit {
			diff := int(card.Rank) - int(topDest.Rank)
			if diff == 1 || diff == -1 {
				c, _ := source.Deal()
				destination.Add(c, true)
				return true
			}
		}
	}

	// Use the destination's Rule (the bound func from your NewStack pattern)
	if destination.Rule(card) {
		c, _ := source.Deal()
		destination.Add(c, true)
		return true
	}
	return false
}

func (*Algerian) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Algerian) Talon() bool                                   { return true }
func (*Algerian) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*Algerian) FoundationBase() bool                          { return false }
func (*Algerian) MaxRedeals() int                               { return 0 }

func (*Algerian) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 104
}

// AvailableMoves - return a list of the available moves.
func (a *Algerian) AvailableMoves(
	tableaus []state.Tableau,
	foundations []state.Foundation,
	_ []state.Talon,
	reserves []state.Reserve,
) []state.Move {
	var moves []state.Move

	// Helper to check and record a move
	checkMove := func(srcStack *state.Stack, destStack *state.Stack) {
		card, err := srcStack.Top()
		if err != nil {
			return
		}

		// 1. Check Foundation-to-Foundation transfer (special Algerian rule)
		if srcStack.Type == state.StackFoundation && destStack.Type == state.StackFoundation {
			topDest, errDest := destStack.Top()
			if errDest == nil && card.Suit == topDest.Suit {
				diff := int(card.Rank) - int(topDest.Rank)
				if diff == 1 || diff == -1 {
					moves = append(moves, state.Move{
						Source:        *srcStack,
						Destination:   *destStack,
						NumberMoving:  1,
						SourceCardTop: card,
					})
					return
				}
			}
		}

		// 2. Check standard rules via the destination's bound Rule
		if destStack.Rule(card) {
			moves = append(moves, state.Move{
				Source:        *srcStack,
				Destination:   *destStack,
				NumberMoving:  1,
				SourceCardTop: card,
			})
		}
	}

	// Iterate through sources and targets
	for i := range tableaus {
		// Tableau -> Foundations
		for j := range foundations {
			checkMove(tableaus[i].Stack, foundations[j].Stack)
		}
		// Tableau -> Tableau
		for j := range tableaus {
			if i == j {
				continue
			}
			checkMove(tableaus[i].Stack, tableaus[j].Stack)
		}
	}

	for i := range reserves {
		// Reserves -> Foundations
		for j := range foundations {
			checkMove(reserves[i].Stack, foundations[j].Stack)
		}
		// Reserves -> Tableau
		for j := range tableaus {
			checkMove(reserves[i].Stack, tableaus[j].Stack)
		}
	}

	// Foundation -> Foundation
	for i := range foundations {
		for j := range foundations {
			if i == j {
				continue
			}
			checkMove(foundations[i].Stack, foundations[j].Stack)
		}
	}

	return moves
}

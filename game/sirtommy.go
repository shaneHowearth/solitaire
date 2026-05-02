package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Sir Tommy - https://en.wikipedia.org/wiki/Sir_Tommy
type SirTommy struct{}

var _ Variant = (*SirTommy)(nil)

func (*SirTommy) Name() string {
	return "Sir Tommy"
}

func (*SirTommy) Category() Category {
	return CatFoundation
}

func (*SirTommy) Description() string {
	return "The ancestor of many modern games. There is no building on the tableau; you must use the four discard piles to organize cards for the foundations."
}

func (*SirTommy) TableauGridSize() (int, int) {
	// 4 wastepiles and 4 foundations
	return 1, 4
}

func (*SirTommy) Decks() int {
	return 1
}

func (*SirTommy) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

func (s *SirTommy) Tableau() []state.StackSpec {
	// Sir Tommy uses 4 wastepiles as its tableau.
	// Wikipedia: "Once placed, it cannot be moved [to another wastepile]"
	wastepiles := make([]state.StackSpec, 4)
	for i := 0; i < 4; i++ {
		wastepiles[i] = state.StackSpec{
			CardCount: [2]int{0, 0},
			// Cards can always be placed on wastepiles, but logic in Move
			// will ensure they only come from the Talon/Stock.
			AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true },
		}
	}
	return wastepiles
}

func (*SirTommy) Fanned() bool { return false }

func (s *SirTommy) Foundations() []state.StackSpec {
	foundations := make([]state.StackSpec, 4)
	for i := 0; i < 4; i++ {
		foundations[i] = state.StackSpec{
			// Foundations start with Aces.
			BaseCard: state.SuitedCard{Rank: state.Ace},
			AddRule:  s.foundationRule,
		}
	}
	return foundations
}

func (*SirTommy) foundationRule(s *state.Stack, c state.SuitedCard) bool {
	// If empty, only an Ace can start.
	if s.Len() == 0 {
		return c.Rank == state.Ace
	}

	// Build up regardless of suit.
	top, _ := s.Top()
	return int(c.Rank)-int(top.Rank) == 1
}

func (*SirTommy) HowToPlay() []string {
	return []string{
		`Cards are dealt one at a time from the stock.`,
		`Foundations build up from Ace to King regardless of suit.`,
		`Cards from the stock that cannot be played to foundations are placed on one of four wastepiles.`,
		`Once a card is placed on a wastepile, it cannot be moved to another wastepile.`,
		`Only the top card of each wastepile is available for the foundations.`,
	}
}

func (s *SirTommy) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	card, err := source.Top()
	if err != nil {
		return false
	}

	// Core Rule: You can't move cards between wastepiles (Tableau to Tableau).
	if source.Type == state.StackTableau && destination.Type == state.StackTableau {
		return false
	}

	if destination.Rule(card) {
		c, _ := source.Deal()
		destination.Add(c, true)
		return true
	}
	return false
}

func (*SirTommy) Compact(_ *state.Stack, _ *state.Stack, _ []*state.Tableau) {}

func (*SirTommy) Talon() bool { return true }

func (*SirTommy) MaxRedeals() int {
	return 0 // Only one pass through the deck.
}

func (*SirTommy) Redeal(_ *state.Talon, _ []*state.Tableau) {}

func (*SirTommy) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 52
}

func (*SirTommy) FoundationBase() bool { return false }

func (s *SirTommy) AvailableMoves(
	tableaus []state.Tableau,
	foundations []state.Foundation,
	talons []state.Talon,
	_ []state.Reserve,
) []state.Move {
	var moves []state.Move

	check := func(src *state.Stack, dest *state.Stack) {
		if card, err := src.Top(); err == nil {
			// Validate movement using the Move logic (prevents Tableau->Tableau)
			if dest.Rule(card) {
				if src.Type == state.StackTableau && dest.Type == state.StackTableau {
					return
				}
				moves = append(moves, state.Move{
					Source: *src, Destination: *dest, NumberMoving: 1, SourceCardTop: card,
				})
			}
		}
	}

	// 1. Waste (Tableau) to Foundations
	for i := range tableaus {
		for j := range foundations {
			check(tableaus[i].Stack, foundations[j].Stack)
		}
	}

	// 2. Stock (Talon Waste) to Foundations or Waste (Tableau)
	for i := range talons {
		for j := range foundations {
			check(talons[i].Waste, foundations[j].Stack)
		}
		for j := range tableaus {
			check(talons[i].Waste, tableaus[j].Stack)
		}
	}

	return moves
}

package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Flower Garden - https://en.wikipedia.org/wiki/Flower_Garden_(solitaire)
type FlowerGarden struct{}

var _ Variant = (*FlowerGarden)(nil)

func (*FlowerGarden) Name() string {
	return "Flower Garden"
}

func (*FlowerGarden) Category() Category {
	return CatSpecialty
}

func (*FlowerGarden) Description() string {
	return "An open-information game where all cards are visible from the start. Use the 16-card 'bouquet' (reserve) to build long sequences in the garden."
}

func (*FlowerGarden) TableauGridSize() (int, int) {
	// 6 columns (flower beds)
	return 1, 6
}

func (*FlowerGarden) Decks() int {
	return 1
}

func (f *FlowerGarden) Reserves() []state.StackSpec {
	// The Bouquet: 16 cards, all are available for play.
	// In many engines, this is represented as a single pile where
	// only the top is accessible, but the player can choose any.
	// To simplify for this engine, we treat it as one stack where
	// the player pulls the top card.
	return []state.StackSpec{
		{
			CardCount: [2]int{16, 16},
			// Cannot build onto the bouquet.
			AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return false },
		},
	}
}

func (f *FlowerGarden) Tableau() []state.StackSpec {
	// The Garden: 6 flower beds with 6 cards each.
	flowerBeds := make([]state.StackSpec, 6)
	for i := 0; i < 6; i++ {
		flowerBeds[i] = state.StackSpec{
			CardCount: [2]int{6, 6},
			AddRule:   f.tableauRule,
		}
	}
	return flowerBeds
}

func (*FlowerGarden) Fanned() bool { return true }

func (*FlowerGarden) tableauRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		// Empty flower beds can be filled with any card.
		return true
	}

	top, _ := s.Top()
	// Build DOWN regardless of suit.
	// top (e.g., 5/4) - card (e.g., 4/3) == 1
	return int(top.Rank)-int(c.Rank) == 1
}

func (f *FlowerGarden) Foundations() []state.StackSpec {
	foundations := make([]state.StackSpec, 4)
	for i := 0; i < 4; i++ {
		foundations[i] = state.StackSpec{
			BaseCard: state.SuitedCard{Rank: state.Ace},
			AddRule:  f.foundationRule,
		}
	}
	return foundations
}

func (*FlowerGarden) foundationRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return c.Rank == state.Ace
	}

	top, _ := s.Top()
	// Build UP in SUIT.
	if c.Suit != top.Suit {
		return false
	}
	return int(c.Rank)-int(top.Rank) == 1
}

func (*FlowerGarden) HowToPlay() []string {
	return []string{
		`The Garden: 6 columns of 6 cards each.`,
		`The Bouquet (Reserve): 16 cards available to play at any time.`,
		`Foundations: Build up in suit from Ace to King.`,
		`Tableau (Flower Beds): Build down regardless of suit.`,
		`Empty flower beds can be filled with any card.`,
		`Cards are moved one at a time.`,
	}
}

func (f *FlowerGarden) Move(source, destination *state.Stack, _ []*state.Tableau, _ []*state.Reserve) bool {
	card, err := source.Top()
	if err != nil {
		return false
	}

	if destination.Rule(card) {
		c, _ := source.Deal()
		destination.Add(c, true)
		return true
	}
	return false
}

func (*FlowerGarden) Compact(_ *state.Stack, _ *state.Stack, _ []*state.Tableau) {}

func (*FlowerGarden) Talon() bool { return false }

func (*FlowerGarden) MaxRedeals() int { return 0 }

func (*FlowerGarden) Redeal(_ *state.Talon, _ []*state.Tableau) {}

func (*FlowerGarden) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 52
}

func (*FlowerGarden) FoundationBase() bool { return false }

func (f *FlowerGarden) AvailableMoves(
	tableaus []*state.Tableau,
	foundations []*state.Foundation,
	talons []*state.Talon,
	reserves []*state.Reserve,
) []state.Move {
	var moves []state.Move

	check := func(src *state.Stack, dest *state.Stack) {
		if card, err := src.Top(); err == nil {
			if dest.Rule(card) {
				moves = append(moves, state.Move{
					Source: *src, Destination: *dest, NumberMoving: 1, SourceCardTop: card,
				})
			}
		}
	}

	// 1. Sources: Tableau (Garden)
	for i := range tableaus {
		// To Foundation
		for j := range foundations {
			check(tableaus[i].Stack, foundations[j].Stack)
		}
		// To other Tableau columns
		for j := range tableaus {
			if i != j {
				check(tableaus[i].Stack, tableaus[j].Stack)
			}
		}
	}

	// 2. Sources: Reserve (Bouquet)
	for i := range reserves {
		// To Foundation
		for j := range foundations {
			check(reserves[i].Stack, foundations[j].Stack)
		}
		// To Tableau
		for j := range tableaus {
			check(reserves[i].Stack, tableaus[j].Stack)
		}
	}

	return moves
}

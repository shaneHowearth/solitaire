package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// American Toad - https://en.wikipedia.org/wiki/American_Toad_(card_game).
type AmericanToad struct {
	baseRank state.Rank
}

var _ Variant = (*AmericanToad)(nil)

func (*AmericanToad) Name() string {
	return "American Toad"
}

func (*AmericanToad) Category() Category {
	return CatKlondike
}

func (*AmericanToad) Description() string {
	return "A challenging Canfield variant using two decks. It features a large reserve pile and foundations that must start with a randomly dealt rank."
}

func (*AmericanToad) TableauGridSize() (int, int) {
	return 1, 8
}

func (*AmericanToad) Decks() int {
	return 2
}

func (a *AmericanToad) Reserves() []state.StackSpec {
	// One reserve pile with 20 cards.
	return []state.StackSpec{
		{
			CardCount: [2]int{20, 20},
			// Cannot build onto the reserve.
			AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return false },
		},
	}
}

func (a *AmericanToad) Tableau() []state.StackSpec {
	tableaus := make([]state.StackSpec, 8)
	for i := 0; i < 8; i++ {
		tableaus[i] = state.StackSpec{
			CardCount: [2]int{1, 1},
			AddRule:   a.tableauRule,
		}
	}
	return tableaus
}

func (a *AmericanToad) tableauRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		// Technically anything can land here, but the Move/AvailableMoves
		// logic will prevent Tableau->Tableau moves for empty spaces.
		return true
	}

	top, _ := s.Top()
	if c.Suit != top.Suit {
		return false
	}

	// Tableau builds DOWN in SUIT with wrapping.
	if top.Rank == state.Ace && c.Rank == state.King {
		return true
	}

	return int(top.Rank)-int(c.Rank) == 1
}

func (a *AmericanToad) Foundations() []state.StackSpec {
	foundations := make([]state.StackSpec, 8)
	for i := 0; i < 8; i++ {
		foundations[i] = state.StackSpec{
			// The engine will set the Stack.Base of the first foundation
			// when it deals the starter card.
			AddRule: a.foundationRule,
		}
	}
	return foundations
}

func (a *AmericanToad) foundationRule(s *state.Stack, c state.SuitedCard) bool {
	// If the slot is empty
	if s.Len() == 0 {
		// If this is the very first foundation (position 0), the engine
		// often places the first card there automatically.
		// For the other 7 piles, they must match the rank of the first one.

		// If your engine populates s.Base when FoundationBase() is true:
		return c.Rank == s.Base.Rank
	}

	// Building UP in SUIT with wrapping
	top, _ := s.Top()
	if c.Suit != top.Suit {
		return false
	}

	// Wrap: Ace onto King
	if top.Rank == state.King && c.Rank == state.Ace {
		return true
	}

	// Build up: card (e.g., Two/1) - top (e.g., Ace/0) == 1
	return int(c.Rank)-int(top.Rank) == 1
}

func (*AmericanToad) HowToPlay() []string {
	return []string{
		`Foundations build up in suit, wrapping King to Ace.`,
		`Tableau builds down in suit, wrapping Ace to King.`,
		`Reserve of 20 cards automatically fills empty tableau spaces.`,
		`Move only the top card or the entire tableau pile.`,
		`Once reserve is empty, only cards from the pack can fill tableau spaces.`,
		`Two passes through the pack allowed.`,
	}
}

func (a *AmericanToad) Move(source, destination *state.Stack, tableaus []*state.Tableau) bool {
	card, err := source.Top()
	if err != nil {
		return false
	}

	if destination.Rule(card) {
		// Rule for filling empty spaces once reserve is empty:
		// "spaces in the tableau can be filled with a card from the pack, but NOT from another tableau pile."
		if destination.Type == state.StackTableau && destination.Len() == 0 && source.Type == state.StackTableau {
			return false
		}

		c, _ := source.Deal()
		destination.Add(c, true)
		return true
	}
	return false
}

// Compact handles the automatic filling of empty tableau spaces from the reserve.
func (*AmericanToad) Compact(_ *state.Stack, reserve *state.Stack, tableaus []*state.Tableau) {
	if reserve == nil || reserve.Len() == 0 {
		return
	}

	for i := range tableaus {
		if tableaus[i].Len() == 0 {
			card, err := reserve.Deal()
			if err == nil {
				tableaus[i].Stack.Add(card, true)
			}
			// One fill per compaction cycle is safer, but we can fill all.
			if reserve.Len() == 0 {
				break
			}
		}
	}
}

func (*AmericanToad) Talon() bool { return true }

func (*AmericanToad) MaxRedeals() int {
	return 1 // Two passes total = 1 redeal.
}

func (a *AmericanToad) Redeal(talon *state.Talon, _ []*state.Tableau) {
	// standard recycling logic from waste to stock.
}

func (*AmericanToad) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 104
}

func (*AmericanToad) FoundationBase() bool { return true }

func (a *AmericanToad) AvailableMoves(
	tableaus []state.Tableau,
	foundations []state.Foundation,
	talons []state.Talon,
	reserves []state.Reserve,
) []state.Move {
	var moves []state.Move

	check := func(src *state.Stack, dest *state.Stack) {
		if card, err := src.Top(); err == nil {
			if dest.Rule(card) {
				// CRITICAL RULE: Empty spaces cannot be filled from another tableau pile.
				if dest.Type == state.StackTableau && dest.Len() == 0 && src.Type == state.StackTableau {
					return
				}
				moves = append(moves, state.Move{
					Source: *src, Destination: *dest, NumberMoving: 1, SourceCardTop: card,
				})
			}
		}
	}

	// 1. Sources: Tableau
	for i := range tableaus {
		for j := range foundations {
			check(tableaus[i].Stack, foundations[j].Stack)
		}
		for j := range tableaus {
			if i != j {
				check(tableaus[i].Stack, tableaus[j].Stack)
			}
		}
	}

	// 2. Sources: Reserve
	for i := range reserves {
		for j := range foundations {
			check(reserves[i].Stack, foundations[j].Stack)
		}
		for j := range tableaus {
			check(reserves[i].Stack, tableaus[j].Stack)
		}
	}

	// 3. Sources: Waste (from Talon)
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

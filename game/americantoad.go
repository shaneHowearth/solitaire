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

func (*AmericanToad) Fanned() bool { return true }

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

	// Foundation is complete once it holds a full 13-rank cycle.
	if s.Len() >= int(state.RankCount) {
		return false
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

func (a *AmericanToad) Move(source, destination *state.Stack, tableaus []*state.Tableau, reserves []*state.Reserve) bool {
	// Let checkMove validate and determine how many cards. This guarantees 1:1 behavior with the engine.
	move := checkMove(source, destination, true, true)
	if move.NumberMoving == 0 {
		return false
	}

	// Rule for filling empty spaces once reserve is empty:
	// "spaces in the tableau can be filled with a card from the pack, but NOT from another tableau pile."
	if destination.Type == state.StackTableau && destination.Len() == 0 && source.Type == state.StackTableau {
		return false
	}

	// Rule: Move only the top card or the entire tableau pile.
	if source.Type == state.StackTableau && move.NumberMoving != 1 && move.NumberMoving != source.Len() {
		return false
	}

	canMove := true
	if source.Type == state.StackWaste && destination.Type == state.StackTalon {
		if destination.Len() != 0 {
			canMove = false
		} else {
			canMove = destination.CanReceiveMore()
		}
	}

	// Execute the physical move
	if canMove {
		return Move(source, destination, true)
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

func (a *AmericanToad) Redeal(talon *state.Talon, tableaus []*state.Tableau) {
	a.Move(talon.Waste, talon.Stock, tableaus, nil)
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
	tableaus []*state.Tableau,
	foundations []*state.Foundation,
	talons []*state.Talon,
	reserves []*state.Reserve,
) []state.Move {
	var moves []state.Move

	// 1. Sources: Tableau
	for i := range tableaus {
		// Tableau to Foundations
		for j := range foundations {
			if move := checkMove(tableaus[i].Stack, foundations[j].Stack, false, true); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}

		// Tableau to Tableau
		for j := range tableaus {
			if i != j {
				if move := checkMove(tableaus[i].Stack, tableaus[j].Stack, false, true); move.NumberMoving > 0 {
					// CRITICAL RULE: Empty spaces cannot be filled from another tableau pile.
					if tableaus[j].Stack.Len() == 0 {
						continue
					}

					// CRITICAL RULE: Move only the top card or the entire tableau pile.
					if move.NumberMoving != 1 && move.NumberMoving != tableaus[i].Stack.Len() {
						continue
					}

					moves = append(moves, move)
				}
			}
		}
	}

	// 2. Sources: Reserve
	for i := range reserves {
		top, err := reserves[i].Stack.Top()
		if err != nil {
			continue // reserve empty
		}
		for j := range foundations {
			if foundations[j].Stack.Rule(top) {
				moves = append(moves, state.Move{
					Source:        *reserves[i].Stack,
					Destination:   *foundations[j].Stack,
					NumberMoving:  1,
					SourceCardTop: top,
				})
			}
		}
		for j := range tableaus {
			if tableaus[j].Stack.Rule(top) {
				moves = append(moves, state.Move{
					Source:        *reserves[i].Stack,
					Destination:   *tableaus[j].Stack,
					NumberMoving:  1,
					SourceCardTop: top,
				})
			}
		}
	}

	return moves
}

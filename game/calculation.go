package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Calculation struct{}

var _ Variant = (*Calculation)(nil)

func (*Calculation) Name() string       { return "Calculation" }
func (*Calculation) Category() Category { return CatFoundation }
func (*Calculation) Description() string {
	return "A classic foundation game where foundations build by different intervals (1, 2, 3, 4)."
}

func (*Calculation) TableauGridSize() (int, int) { return 1, 4 }
func (*Calculation) Decks() int                  { return 1 }
func (*Calculation) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Calculation) Tableau() []state.StackSpec {
	return []state.StackSpec{{AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true }},
		{AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true }},
		{AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true }},
		{AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true }}}
}

func (*Calculation) Fanned() bool                          { return false }
func (*Calculation) TableauPosition(i int) (int, int, int) { return 0, i, 0 }

func (*Calculation) Foundations() []state.StackSpec {
	// Foundations: A, 2, 3, 4. Rules: +1, +2, +3, +4.
	// Rank cycling (A=1, K=13).
	rule := func(inc int) func(*state.Stack, state.SuitedCard) bool {

		return func(foundation *state.Stack, c state.SuitedCard) bool {
			if foundation.Len() == 0 {
				return int(c.Rank)+1 == inc
			}
			top, _ := foundation.Top()
			// Foundations stop being built up at the King
			if top.Rank == state.King {
				return false
			}
			expected := ((int(top.Rank)+inc-1)%13 + 1) % 13
			return int(c.Rank) == expected
		}
	}
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Undefined}, AddRule: rule(1)},
		{BaseCard: state.SuitedCard{Rank: state.Two, Suit: state.Undefined}, AddRule: rule(2)},
		{BaseCard: state.SuitedCard{Rank: state.Three, Suit: state.Undefined}, AddRule: rule(3)},
		{BaseCard: state.SuitedCard{Rank: state.Four, Suit: state.Undefined}, AddRule: rule(4)},
	}
}

func (*Calculation) HowToPlay() []string {
	return []string{"Build foundations by 1, 2, 3, and 4.", "Tableau piles hold any card.", "Play cards from stock to tableau or foundation."}
}

func (*Calculation) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 52
}
func (*Calculation) MaxRedeals() int { return 0 }
func (*Calculation) Move(src, dest *state.Stack, _ []*state.Tableau) bool {
	// Once a card is placed on a tableau, it stays there or goes on a
	// foundation.
	if src.Type == state.StackTableau && dest.Type != state.StackFoundation {
		return false
	}

	// Nothing gets moved back onto the talon
	if dest.Type == state.StackTalon {
		return false
	}

	card, err := src.Top()
	if err != nil {
		return false
	}
	if dest.Rule(card) {
		c, _ := src.Deal()
		dest.Add(c, true)
		return true
	}
	return false
}
func (*Calculation) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Calculation) Talon() bool                                   { return true }
func (*Calculation) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*Calculation) FoundationBase() bool                          { return false }
func (*Calculation) AvailableMoves(t []state.Tableau, f []state.Foundation, talons []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}
	for _, talon := range talons {
		for i := range t {
			if m := checkMove(talon.Waste, t[i].Stack, false, false); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
		for i := range f {
			if m := checkMove(talon.Waste, f[i].Stack, false, false); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
	}
	for i := range t {
		for j := range f {
			if m := checkMove(t[i].Stack, f[j].Stack, false, false); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
	}
	return moves
}

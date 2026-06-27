package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Appreciate struct{}

var _ Variant = (*Appreciate)(nil)

func (*Appreciate) Name() string       { return "Appreciate" }
func (*Appreciate) Category() Category { return CatFoundation }
func (*Appreciate) Description() string {
	return "A Calculation-style game with a full 8x6 tableau and open information."
}

func (*Appreciate) TableauGridSize() (int, int) { return 6, 8 }
func (*Appreciate) Decks() int                  { return 1 }
func (*Appreciate) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Appreciate) Tableau() []state.StackSpec {
	spec := state.StackSpec{AddRule: func(_ *state.Stack, _ state.SuitedCard) bool { return true }}
	res := make([]state.StackSpec, 48)
	for i := range res {
		res[i] = spec
	}
	return res
}

func (*Appreciate) Fanned() bool                          { return false }
func (*Appreciate) TableauPosition(i int) (int, int, int) { return i / 8, i % 8, 0 }

func (*Appreciate) Foundations() []state.StackSpec {
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

func (*Appreciate) HowToPlay() []string {
	return []string{"Tableau is an 8x6 grid.", "Foundations build by 1, 2, 3, 4.", "Information is open."}
}

func (*Appreciate) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	total := 0
	for _, f := range foundations {
		total += f.Len()
	}
	return total == 52
}

func (*Appreciate) MaxRedeals() int { return 0 }
func (*Appreciate) Move(src, dest *state.Stack, _ []*state.Tableau) bool {
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
func (*Appreciate) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Appreciate) Talon() bool                                   { return true }
func (*Appreciate) Redeal(_ *state.Talon, _ []*state.Tableau)     {}
func (*Appreciate) FoundationBase() bool                          { return false }

func (*Appreciate) AvailableMoves(t []state.Tableau, f []state.Foundation, talons []state.Talon, _ []state.Reserve) []state.Move {
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

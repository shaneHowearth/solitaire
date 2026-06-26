package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Tasmanian struct{}

var _ Variant = (*Tasmanian)(nil)

func (*Tasmanian) Name() string { return "Tasmanian Solitaire" }

func (*Tasmanian) Category() Category { return CatSpider }

func (*Tasmanian) Description() string {
	return "A variant of Australian Patience with unlimited redeals."
}

func (*Tasmanian) TableauGridSize() (int, int) { return 1, 7 }

func (*Tasmanian) Decks() int { return 1 }

func (*Tasmanian) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Tasmanian) Tableau() []state.StackSpec {
	var MinusOneSuitRule = func(tableau *state.Stack, card state.SuitedCard) bool {
		if (*tableau).Len() == 0 {
			return card.Rank == state.King
		}
		topCard, err := (*tableau).Top()
		if err != nil {
			return false
		}
		return (card.Suit == topCard.Suit) && (topCard.Rank-card.Rank == 1)
	}
	spec := state.StackSpec{
		AddRule:   MinusOneSuitRule,
		CardCount: [2]int{4, 4},
		BaseCard:  state.SuitedCard{Rank: state.King},
	}
	return []state.StackSpec{spec, spec, spec, spec, spec, spec, spec}
}

func (*Tasmanian) Fanned() bool { return true }

func (*Tasmanian) TableauPosition(i int) (int, int, int) { return 0, i, 0 }

func (*Tasmanian) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (*Tasmanian) HowToPlay() []string {
	return []string{
		"Similar to Australian Patience, build down in suit on the 7 piles of 4 cards.",
		"Any face-up card can be moved with its pile.",
		"Empty spaces can only be filled by a King.",
		"This version allows unlimited redeals of the stock.",
	}
}

func (*Tasmanian) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, f := range foundations {
		if f.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*Tasmanian) MaxRedeals() int { return -1 }

func (*Tasmanian) Move(s, d *state.Stack, _ []*state.Tableau) bool { return Move(s, d, true) }

func (*Tasmanian) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

func (*Tasmanian) Talon() bool { return true }

func (t *Tasmanian) Redeal(talon *state.Talon, tableaus []*state.Tableau) {
	t.Move(talon.Waste, talon.Stock, tableaus)
}

func (*Tasmanian) FoundationBase() bool { return false }

func (*Tasmanian) AvailableMoves(t []state.Tableau, f []state.Foundation, _ []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}
	for fIdx := range f {
		for tIdx := range t {
			if m := checkMove(t[tIdx].Stack, f[fIdx].Stack, true, true); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
	}
	for dIdx := range t {
		for sIdx := range t {
			if dIdx == sIdx {
				continue
			}
			if m := checkMove(t[sIdx].Stack, t[dIdx].Stack, true, true); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
	}
	return moves
}

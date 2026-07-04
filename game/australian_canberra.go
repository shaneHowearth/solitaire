package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Canberra struct{}

var _ Variant = (*Canberra)(nil)

func (*Canberra) Name() string { return "Canberra" }

func (*Canberra) Category() Category { return CatSpider }

func (*Canberra) Description() string {
	return "A variant of Australian Patience allowing one redeal of the stock."
}

func (*Canberra) TableauGridSize() (int, int) {
	return 1, 7
}

func (*Canberra) Decks() int { return 1 }

func (*Canberra) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Canberra) Tableau() []state.StackSpec {
	// Build down in suit
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

func (*Canberra) Fanned() bool { return true }

func (*Canberra) TableauPosition(i int) (int, int, int) { return 0, i, 0 }

func (*Canberra) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (*Canberra) HowToPlay() []string {
	return []string{
		"The Tableau is filled with 7 piles of 4 cards each. Build down in suit.",
		"Like Yukon, any face-up card can be moved, and all cards on top move with it.",
		"Only a King can be moved to an empty space.",
		"One pass through the deck is allowed, plus one redeal.",
	}
}

func (*Canberra) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, f := range foundations {
		if f.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*Canberra) MaxRedeals() int { return 1 }

func (*Canberra) Move(s, d *state.Stack, _ []*state.Tableau, _ []*state.Reserve) bool {
	return Move(s, d, true)
}

func (*Canberra) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

func (*Canberra) Talon() bool { return true }

func (c *Canberra) Redeal(talon *state.Talon, tableaus []*state.Tableau) {
	c.Move(talon.Waste, talon.Stock, tableaus, nil)
}

func (*Canberra) FoundationBase() bool { return false }

func (*Canberra) AvailableMoves(t []*state.Tableau, f []*state.Foundation, _ []*state.Talon, _ []*state.Reserve) []state.Move {
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

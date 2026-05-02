package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Kuipers struct{}

var _ Variant = (*Kuipers)(nil)

func (*Kuipers) Name() string       { return "Kuipers" }
func (*Kuipers) Category() Category { return CatKlondike }
func (*Kuipers) Description() string {
	return "Eight columns instead of seven, with unlimited single-card draws from the deck."
}

func (*Kuipers) TableauGridSize() (int, int) { return 1, 8 }
func (*Kuipers) Decks() int                  { return 1 }
func (*Kuipers) Fanned() bool                { return true }
func (*Kuipers) Talon() bool                 { return true }
func (*Kuipers) MaxRedeals() int             { return -1 }
func (*Kuipers) FoundationBase() bool        { return false }

func (*Kuipers) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Kuipers) Tableau() []state.StackSpec {
	specs := make([]state.StackSpec, 8)
	for i := 0; i < 8; i++ {
		specs[i] = state.StackSpec{
			AddRule:   MinusOneRule,
			CardCount: [2]int{i + 1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		}
	}
	return specs
}

func (*Kuipers) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs},
			AddRule:  PlusOneRule,
		},
		{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades},
			AddRule:  PlusOneRule,
		},
	}
}

func (*Kuipers) HowToPlay() []string {
	return []string{
		"Eight tableau columns are dealt (1 card in the first, increasing to 8 in the last).",
		"Foundations are built up by suit from Ace to King.",
		"Tableau piles are built down by alternating colors.",
		"The stock is dealt one card at a time to the waste pile.",
		"Unlimited passes through the deck are allowed.",
	}
}

func (k *Kuipers) Move(s, d *state.Stack, t []*state.Tableau) bool { return Move(s, d, true) }

func (*Kuipers) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

func (k *Kuipers) AvailableMoves(t []state.Tableau, f []state.Foundation, tal []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}
	for fIdx := range f {
		for tIdx := range t {
			if move := checkMove(t[tIdx].Stack, f[fIdx].Stack, false, true); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
		if len(tal) > 0 {
			if move := checkMove(tal[0].Waste, f[fIdx].Stack, false, true); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
	}
	for dIdx := range t {
		for sIdx := range t {
			if dIdx == sIdx {
				continue
			}
			if move := checkMove(t[sIdx].Stack, t[dIdx].Stack, false, true); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
		if len(tal) > 0 {
			if move := checkMove(tal[0].Waste, t[dIdx].Stack, false, true); move.NumberMoving > 0 {
				moves = append(moves, move)
			}
		}
	}
	return moves
}

func (*Kuipers) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (k *Kuipers) Redeal(talon *state.Talon, _ []*state.Tableau) {
	if talon.Stock.Len() > 0 {
		card, _ := talon.Stock.Deal()
		talon.Waste.Add(card, true)
	} else {
		for talon.Waste.Len() > 0 {
			card, _ := talon.Waste.Deal()
			talon.Stock.Add(card, false)
		}
	}
}

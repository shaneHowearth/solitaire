package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Somerset struct{}

var _ Variant = (*Somerset)(nil)

func (*Somerset) Name() string       { return "Somerset" }
func (*Somerset) Category() Category { return CatKlondike }
func (*Somerset) Description() string {
	return "A Klondike variant where all cards are dealt into the tableau in a 10-9-8...3 pattern."
}

func (*Somerset) TableauGridSize() (int, int) { return 1, 8 }
func (*Somerset) Decks() int                  { return 1 }
func (*Somerset) Fanned() bool                { return true }
func (*Somerset) Talon() bool                 { return false }
func (*Somerset) MaxRedeals() int             { return 0 }
func (*Somerset) FoundationBase() bool        { return false }

func (*Somerset) Reserves() []state.StackSpec { return []state.StackSpec{} }

func (*Somerset) Tableau() []state.StackSpec {
	counts := []int{10, 9, 8, 7, 6, 5, 4, 3}
	specs := make([]state.StackSpec, len(counts))
	for i, count := range counts {
		specs[i] = state.StackSpec{
			AddRule:   MinusOneRule,
			CardCount: [2]int{count, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		}
	}
	return specs
}

func (*Somerset) Foundations() []state.StackSpec {
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

func (*Somerset) HowToPlay() []string {
	return []string{
		"The entire deck is dealt into 8 tableau columns in a 10, 9, 8, 7, 6, 5, 4, 3 pattern.",
		"Build foundations up by suit from Ace to King.",
		"Tableau piles are built down by alternating colors.",
		"Empty tableau spaces can be filled with any King.",
	}
}

func (s *Somerset) Move(src, dst *state.Stack, t []*state.Tableau) bool { return Move(src, dst, true) }

func (*Somerset) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

func (s *Somerset) AvailableMoves(t []state.Tableau, f []state.Foundation, _ []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}
	for fIdx := range f {
		for tIdx := range t {
			if move := checkMove(t[tIdx].Stack, f[fIdx].Stack, false, true); move.NumberMoving > 0 {
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
	}
	return moves
}

func (*Somerset) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*Somerset) Redeal(_ *state.Talon, _ []*state.Tableau) {}

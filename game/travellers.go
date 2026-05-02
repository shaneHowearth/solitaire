package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type Travellers struct{}

var _ Variant = (*Travellers)(nil)

func (*Travellers) Name() string                { return "Travellers" }
func (*Travellers) Category() Category          { return CatSpecialty }
func (*Travellers) Reserves() []state.StackSpec { return []state.StackSpec{} }
func (*Travellers) Description() string {
	return "A rhythmic shuttling game. Move the top card to the bottom of the pile matching its rank. The game ends when the fourth King appears."
}

func (*Travellers) Decks() int { return 1 }

func (*Travellers) TableauGridSize() (int, int) {
	// 13 piles fit well in a 3x5 or 2x7 layout
	return 3, 5
}

func (*Travellers) Tableau() []state.StackSpec {
	specs := make([]state.StackSpec, 13)
	for i := 0; i < 13; i++ {
		targetRank := state.Rank(i)

		// Piles 0-11 (A-Q) start fully face-down (4, 0)
		count := [2]int{4, 0}

		// Pile 12 (Kings) starts with one card visible to begin the chain
		if i == 12 {
			count = [2]int{4, 1}
		}

		specs[i] = state.StackSpec{
			CardCount: count,
			AddRule: func(s *state.Stack, c state.SuitedCard) bool {
				return c.Rank == targetRank
			},
		}
	}
	return specs
}

func (*Travellers) Fanned() bool { return true }

// Foundations are technically empty because Kings have their own Tableau pile
func (*Travellers) Foundations() []state.StackSpec { return []state.StackSpec{} }

func (*Travellers) HowToPlay() []string {
	return []string{
		"13 piles represent Ace through King.",
		"Only the King pile starts with a face-up card.",
		"Move a card to the bottom of its matching rank pile.",
		"The card uncovered at the top of that pile becomes your next move.",
		"The game is won if all 13 piles are sorted into four-of-a-kind.",
	}
}

func (t *Travellers) Move(source, destination *state.Stack, allTableaus []*state.Tableau) bool {
	if source.Len() == 0 || destination.Type != state.StackTableau {
		return false
	}

	card, _ := source.Top()

	// Find destination index
	destIdx := -1
	for i, tab := range allTableaus {
		if tab.Stack == destination {
			destIdx = i
			break
		}
	}

	if destIdx != -1 && card.Rank == state.Rank(destIdx) {
		if source == destination && source.Len() <= 1 {
			return false
		}

		// Perform the Shuttle
		c, _ := source.Deal()
		destination.AddBottom(c, true)

		// 1. ALWAYS reveal the card at the top of the destination.
		// This keeps the "rhythmic shuttle" moving.
		t.revealNext(destination)

		// 2. ONLY reveal the next card from the King pile if we just
		// "parked" a King there. This ensures we only start a new
		// chain when the previous one hits a dead end.
		if source.TableauPosition == 12 && card.Rank == state.King {
			t.revealNext(source)
		}

		return true
	}

	return false
}

func (t *Travellers) revealNext(s *state.Stack) {
	if s.Len() > 0 {
		top, err := s.Top()
		if err == nil && !top.Visible {
			// Engine flip: remove and re-add face-up
			_, _ = s.Deal()
			s.Add(top, true)
		}
	}
}

func (*Travellers) HasWon(tableau []*state.Tableau, _ []*state.Foundation) bool {
	for i, t := range tableau {
		if t.Stack.Len() != 4 {
			return false
		}
		for _, c := range t.Stack.GetCards() {
			if c.Rank != state.Rank(i) {
				return false
			}
		}
	}
	return true
}

func (*Travellers) MaxRedeals() int                                     { return 0 }
func (*Travellers) Talon() bool                                         { return false }
func (*Travellers) Redeal(talon *state.Talon, tableau []*state.Tableau) {}
func (*Travellers) Compact(_, _ *state.Stack, _ []*state.Tableau)       {}
func (*Travellers) FoundationBase() bool                                { return false }

func (t *Travellers) AvailableMoves(tableau []state.Tableau, _ []state.Foundation, _ []state.Talon, _ []state.Reserve) []state.Move {
	moves := []state.Move{}

	// Priority: Only the "Intruder" (card on the wrong pile) is moveable.
	// This enforces the single-thread rhythmic movement.
	for i, tab := range tableau {
		if tab.Stack.Len() == 0 {
			continue
		}
		top, _ := tab.Stack.Top()

		// If it's face-up and doesn't belong here, it's the active card
		if top.Visible && top.Rank != state.Rank(i) {
			targetIdx := int(top.Rank)
			moves = append(moves, state.Move{
				Source:        *tab.Stack,
				Destination:   *tableau[targetIdx].Stack,
				NumberMoving:  1,
				SourceCardTop: top,
			})
			return moves // Strictly one active card at a time
		}
	}

	return moves
}

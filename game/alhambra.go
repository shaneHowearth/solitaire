package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Alhambra implements the Alhambra solitaire variant.
type Alhambra struct{}

// Ensure Alhambra implements the Variant interface.
var _ Variant = (*Alhambra)(nil)

// Name returns the name of the variant.
func (*Alhambra) Name() string { return "Alhambra" }

// Category returns the category classification.
func (*Alhambra) Category() Category { return CatKlondike }

// Description provides a brief overview of the game rules.
func (*Alhambra) Description() string {
	return "A unique two-deck game where cards from eight reserves can be built onto foundations or back onto the waste pile with bidirectional wrapping sequences."
}

// Decks - Alhambra requires two full decks of playing cards.
func (*Alhambra) Decks() int { return 2 }

// Fanned - Reserve piles display cards face-up.
func (*Alhambra) Fanned() bool { return true }

// TableauGridSize - Alhambra does not use a traditional building tableau.
func (*Alhambra) TableauGridSize() (int, int) { return 0, 0 }

// Tableau - Returns empty as cards are dealt exclusively into reserves.
func (*Alhambra) Tableau() []state.StackSpec { return []state.StackSpec{} }

// Reserves - 8 reserve piles initialized with 4 face-up cards each.
func (*Alhambra) Reserves() []state.StackSpec {
	specs := make([]state.StackSpec, 8)
	for i := 0; i < 8; i++ {
		specs[i] = state.StackSpec{
			CardCount: [2]int{4, 4},
			AddRule: func(s *state.Stack, _ state.SuitedCard) bool {
				// Building onto or between reserve piles is strictly forbidden
				return false
			},
		}
	}
	return specs
}

// Foundations defines 4 ascending stacks (Ace->King) and 4 descending stacks (King->Ace).
func (a *Alhambra) Foundations() []state.StackSpec {
	suits := []state.Suit{state.Hearts, state.Diamonds, state.Clubs, state.Spades}
	specs := make([]state.StackSpec, 8)

	// First 4 foundations: Ascending (Ace to King)
	for i, suit := range suits {
		specs[i] = state.StackSpec{
			BaseCard: state.SuitedCard{Rank: state.Ace, Suit: suit},
			AddRule:  PlusOneRule,
		}
	}

	// Last 4 foundations: Descending (King to Ace)
	for i, suit := range suits {
		specs[i+4] = state.StackSpec{
			BaseCard: state.SuitedCard{Rank: state.King, Suit: suit},
			AddRule:  a.descendingFoundationRule,
		}
	}

	return specs
}

// Custom rule for descending foundations (King downwards to Ace).
func (*Alhambra) descendingFoundationRule(foundation *state.Stack, card state.SuitedCard) bool {
	if foundation.Len() == 0 {
		return card.Rank == foundation.Base.Rank && card.Suit == foundation.Base.Suit
	}
	topCard, err := foundation.Top()
	if err != nil {
		return false
	}
	return card.Suit == foundation.Base.Suit && int(topCard.Rank)-int(card.Rank) == 1
}

// Checks if a card can be played onto the Waste pile from a reserve pile.
func (*Alhambra) canBuildOnWaste(waste *state.Stack, card state.SuitedCard) bool {
	topCard, err := waste.Top()
	if err != nil {
		// If waste is somehow completely empty, any card can land
		return true
	}

	if card.Suit != topCard.Suit {
		return false
	}

	// Allowed ranks are ±1 with wrap-around protection modulo 13
	diffUp := (int(card.Rank) - int(topCard.Rank) + 13) % 13
	diffDown := (int(topCard.Rank) - int(card.Rank) + 13) % 13

	return diffUp == 1 || diffDown == 1
}

// Move executes card routing while preserving engine immutability state boundaries.
func (a *Alhambra) Move(source, destination *state.Stack, _ []*state.Tableau, _ []*state.Reserve) bool {
	card, err := source.Top()
	if err != nil {
		return false
	}

	// Route 1: Target is a Foundation
	if destination.Type == state.StackFoundation {
		if !destination.Rule(card) {
			return false
		}
		_, _ = source.Deal()
		destination.Add(card, true)
		return true
	}

	// Route 2: Target is the Waste Pile (Only valid from Reserves)
	if destination.Type == state.StackWaste && source.Type == state.StackReserve {
		if !a.canBuildOnWaste(destination, card) {
			return false
		}
		_, _ = source.Deal()
		destination.Add(card, true)
		return true
	}

	// Route 3: Standard Deal simulation pass from Stock to Waste
	if destination.Type == state.StackWaste && source.Type == state.StackTalon {
		_, _ = source.Deal()
		destination.Add(card, true)
		return true
	}

	// Route 4: Moving from Waste to stock (redeal).
	if destination.Type == state.StackTalon && source.Type == state.StackWaste && destination.Len() == 0 {
		if !destination.CanReceiveMore() {
			return false
		}
		temp := state.NewStack(
			source.Len(),
			state.SuitedCard{},
			func(*state.Stack) func(state.SuitedCard) bool {
				return func(state.SuitedCard) bool {
					return true
				}
			},
			state.StackUndefined,
		)

		// Collect exactly numCards from the source stack.
		for source.Len() > 0 {
			card, err := source.Top()
			if err != nil {
				// This shouldn't happen since CanMove validated it.
				break
			}

			temp.Add(card, true)

			_, _ = source.Deal()
		}

		temp.Reverse()

		// Move all cards from temp to destination.
		for temp.Len() > 0 {
			card, err := temp.Top()
			if err != nil {
				break
			}

			destination.Add(card, false)
			_, _ = temp.Deal()
		}
	}

	return false
}

// AvailableMoves scans game components to compute all legal execution paths.
func (a *Alhambra) AvailableMoves(
	_ []*state.Tableau,
	foundations []*state.Foundation,
	talon []*state.Talon,
	reserves []*state.Reserve,
) []state.Move {
	var moves []state.Move
	if len(talon) == 0 {
		return moves
	}
	waste := talon[0].Waste

	// 1. Scan Top of Waste Pile for plays onto Foundations
	if wasteCard, err := waste.Top(); err == nil {
		for fIdx := range foundations {
			if foundations[fIdx].Stack.Rule(wasteCard) {
				moves = append(moves, state.Move{
					Source:        *waste,
					Destination:   *foundations[fIdx].Stack,
					NumberMoving:  1,
					SourceCardTop: wasteCard,
				})
			}
		}
	}

	// 2. Scan Reserve Piles for plays onto Foundations or Waste
	for rIdx := range reserves {
		resCard, err := reserves[rIdx].Stack.Top()
		if err != nil {
			continue // This reserve column is exhausted
		}

		// Can play to Foundations?
		for fIdx := range foundations {
			if foundations[fIdx].Stack.Rule(resCard) {
				moves = append(moves, state.Move{
					Source:        *reserves[rIdx].Stack,
					Destination:   *foundations[fIdx].Stack,
					NumberMoving:  1,
					SourceCardTop: resCard,
				})
			}
		}

		// Can play to Waste?
		if a.canBuildOnWaste(waste, resCard) {
			moves = append(moves, state.Move{
				Source:        *reserves[rIdx].Stack,
				Destination:   *waste,
				NumberMoving:  1,
				SourceCardTop: resCard,
			})
		}
	}

	return moves
}

// Redeal flips the waste pile over cleanly to form a new stock. Max 2 redeals.
func (a *Alhambra) Redeal(talon *state.Talon, tableaus []*state.Tableau) {
	a.Move(talon.Waste, talon.Stock, tableaus, nil)
}

// HasWon confirms success when all 8 foundation components reach maximum size capacity.
func (*Alhambra) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, f := range foundations {
		if f.Stack.Len() != state.RankCount {
			return false
		}
	}
	return true
}

// MaxRedeals allows a maximum of 2 flipped passes after the initial run.
func (*Alhambra) MaxRedeals() int { return 2 }

func (*Alhambra) HowToPlay() []string {
	return []string{
		"Foundations are split: 4 build up from Ace to King, 4 build down from King to Ace by matching suit.",
		"Cards from the 8 reserves can be played to foundations, or down onto the waste pile if matching suit and +/- 1 rank (wrapping allowed).",
		"Empty reserve spaces cannot be refilled. No building between reserve piles is allowed.",
		"The waste pile can be flipped over to form a new stock up to two times.",
	}
}

func (*Alhambra) Compact(_, _ *state.Stack, _ []*state.Tableau) {}
func (*Alhambra) FoundationBase() bool                          { return false }
func (*Alhambra) Talon() bool                                   { return true }

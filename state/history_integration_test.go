package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func TestUndoRestoresExactPreviousState(t *testing.T) {
	deck := state.CreateDecks(1)
	s := &state.State{
		Deck: deck,
		Talon: state.NewTalon(1, 1, func(s *state.Stack) func(state.SuitedCard) bool {
			return func(c state.SuitedCard) bool { return true }
		}),
	}

	// Create some history
	for i := 0; i < 5; i++ {
		s.Talon.Stock.Add(s.Deck.Deal(), false)
	}

	history := state.History{}

	// Initial snapshot (State A)
	history.Update(*s)

	// Perform a move (State B)
	s.Talon.Deal()
	history.Update(*s)

	assert.Equal(t, 4, s.Talon.Stock.Len(), "Card should have been dealt")
	assert.Equal(t, 1, s.Talon.Waste.Len(), "Waste should have 1 card")

	// Perform Undo (Back to State A)
	history.Undo(s)

	// Final Assertions
	assert.Equal(t, 5, s.Talon.Stock.Len(), "Stock count should be restored")
	assert.Equal(t, 0, s.Talon.Waste.Len(), "Waste should be empty after undo")

	// Deep Isolation Check
	// Modify the history's stored state (if we could access it) to ensure
	// the live state isn't just a pointer to history.
	s.Talon.Stock.Deal()

	// If the history and live state shared pointers, a future undo would be corrupted.
	history.Undo(s)
	assert.Equal(t, 4, s.Talon.Stock.Len(), "Undo should still restore original count despite interim live changes")
}

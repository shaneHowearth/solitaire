package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func TestUndoRestoresExactPreviousState(t *testing.T) {
	deck := state.CreateDecks(1)
	testState := &state.State{
		Deck: deck,
		Talon: state.NewTalon(1, 1, func(testState *state.Stack) func(state.SuitedCard) bool {
			return func(c state.SuitedCard) bool { return true }
		}),
	}

	// Create some history.
	for i := 0; i < 5; i++ {
		testState.Talon.Stock.Add(testState.Deck.Deal(), false)
	}

	history := state.History{}

	// Initial snapshot (State A).
	history.Update(*testState)

	// Perform a move (State B).
	_ = testState.Talon.Deal()
	history.Update(*testState)

	assert.Equal(t, 4, testState.Talon.Stock.Len(), "Card should have been dealt")
	assert.Equal(t, 1, testState.Talon.Waste.Len(), "Waste should have 1 card")

	// Perform Undo (Back to State A).
	history.Undo(testState)

	// Final Assertions.
	assert.Equal(t, 5, testState.Talon.Stock.Len(), "Stock count should be restored")
	assert.Equal(t, 0, testState.Talon.Waste.Len(), "Waste should be empty after undo")

	// Deep Isolation Check.
	// Modify the history's stored state (if we could access it) to ensure.
	// the live state isn't just a pointer to history.
	_, _ = testState.Talon.Stock.Deal()

	// If the history and live state shared pointers, a future undo would be corrupted.
	history.Undo(testState)
	assert.Equal(t, 4, testState.Talon.Stock.Len(),
		"Undo should still restore original count despite interim live changes",
	)
}

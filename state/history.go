package state

type History struct {
	states  []*State
	current int // where the current state is in the list.
}

func (history *History) Update(gameState State) {
	clonedState := gameState.CloneGameState()

	// Handle the "Redo" branch: If the user made a new move after an Undo,.
	// we must truncate and overwrite any future "redo" states.
	if history.current < len(history.states) {
		history.states = history.states[:history.current]
	}

	// Append the newly isolated state and advance the current index.
	history.states = append(history.states, clonedState)
	history.current = len(history.states)
}

func (history *History) Undo(gameState *State) {
	if history.current <= 1 {
		// Nothing to do.
		return
	}

	// Get the pointer to the historical snapshot.
	historicalSnapshot := history.states[history.current-2]

	// Create a fresh, isolated DEEP CLONE of the historical snapshot.
	// This ensures that the state we are about to restore is NOT sharing any.
	// pointers with the history slice.
	restoredState := historicalSnapshot.CloneGameState()

	// Overwrite the *value* of the live game state with the.
	// *value* of the newly restored and isolated clone.
	*gameState = *restoredState

	history.current--
}

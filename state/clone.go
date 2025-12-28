package state

// CloneGameState -.
func (gameState *State) CloneGameState() *State {
	if gameState == nil {
		return nil
	}

	clonedState := &State{}

	clonedState.Talon = gameState.Talon.Clone()
	clonedState.Deck = gameState.Deck.Clone()

	clonedState.Tableau = make([]*Tableau, len(gameState.Tableau))
	for idx, tableau := range gameState.Tableau {
		// clonedTableau = append(clonedTableau, gameState.Tableau[idx].Clone()).
		clonedState.Tableau[idx] = tableau.Clone()
	}

	clonedState.Foundations = make([]*Foundation, len(gameState.Foundations))
	for idx, foundation := range gameState.Foundations {
		// foundation := gameState.Foundations[idx].Clone().
		// clonedFoundations = append(clonedFoundations, &foundation).
		// fClone := foundation.Clone().
		clonedState.Foundations[idx] = foundation.Clone()
	}

	clonedState.Reserves = make([]*Reserve, len(gameState.Reserves))
	for idx, reserve := range gameState.Reserves {
		// clonedReserves = append(clonedReserves, gameState.Reserves[idx].Clone()).
		clonedState.Reserves[idx] = reserve.Clone()
	}

	return clonedState
	// return &State{.
	// 	Deck:        clonedDeck,.
	// 	Tableau:     clonedTableau,.
	// 	Talon:       clonedTalon,.
	// 	Foundations: clonedFoundations,.
	// 	Reserves:    clonedReserves,.
	// }.
}

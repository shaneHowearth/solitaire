package solitaire

import (
	"github.com/shanehowearth/solitaire/state"
)

func (instance *Instance) setupGameState() {
	// Get the foundation information for the game.
	numFoundations, foundationBase, foundationRule := instance.Game.Foundations()
	// Get the tableau information for the game.
	numTableau, tableauBase, tableauRule := instance.Game.Tableau()

	// Create the state/model for the game.
	gameState := state.New(
		instance.Game.Decks(),
		numFoundations,
		foundationBase,
		foundationRule,
		numTableau,
		tableauBase,
		tableauRule,
		1,
		1,
		// Talon rule is to allow everything to be added to its stacks.
		func(state.SuitedCard) bool {
			return true
		},
	)

	// Copy the game state instantiated model into the current instance.
	instance.Foundations = gameState.Foundations
	instance.Tableau = gameState.Tableau
	instance.Talon = gameState.Talon
	instance.Deck = gameState.Deck

	instance.dealCards()
}

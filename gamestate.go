package solitaire

import (
	"github.com/shanehowearth/solitaire/state"
)

func (instance *Instance) setupGameState() {
	// Get the foundation information for the game.
	numFoundations, foundationBase, foundationRule := instance.Game.Foundations()
	// Get the tableau information for the game.
	numTableau, tableauBase, tableauRule := instance.Game.Tableau()
	// Get any reserves information for the game.
	numReserves, _, reserveRule := instance.Game.Reserves()

	// Create the state/model for the game.
	gameState := state.New(
		// Number of decks.
		instance.Game.Decks(),
		// Foundation Setup.
		numFoundations,
		foundationBase,
		foundationRule,
		// Tableau Setup.
		numTableau,
		tableauBase,
		tableauRule,
		// Reserve Setup.
		numReserves,
		reserveRule,
		// Talon Setup.
		1,
		1,
		// Talon rule is to allow everything to be added to its stacks.
		func(state.SuitedCard) bool {
			return true
		},
	)

	// Copy the game state instantiated model into the current instance.
	instance.Reserves = gameState.Reserves
	instance.Foundations = gameState.Foundations
	instance.Tableau = gameState.Tableau
	instance.Talon = gameState.Talon
	instance.Deck = gameState.Deck

	instance.dealCards()
}

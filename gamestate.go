package solitaire

import (
	"github.com/shanehowearth/solitaire/state"
)

func (instance *Instance) setupGameState() {
	// Get the foundation information for the game.
	foundationSpec := instance.Game.Foundations()
	// Get the tableau information for the game.
	tableauSpec := instance.Game.Tableau()
	// Get any reserves information for the game.
	reserveSpec := instance.Game.Reserves()

	// Create the state/model for the game.
	gameState := state.New(
		// Number of decks.
		instance.Game.Decks(),
		// Foundation Setup.
		foundationSpec,
		// Tableau Setup.
		tableauSpec,
		// Reserve Setup.
		reserveSpec,
		// Talon Setup.
		instance.Game.MaxRedeals(),
		1, // how many cards per deal.
		// Talon rule is to allow everything to be added to its stacks.
		func(*state.Stack) func(state.SuitedCard) bool {
			return func(state.SuitedCard) bool {
				return true
			}
		},
	)

	// Copy the game state instantiated model into the current instance.
	instance.State.Reserves = gameState.Reserves
	instance.State.Foundations = gameState.Foundations
	instance.State.Tableau = gameState.Tableau
	instance.State.Talon = gameState.Talon
	instance.State.Deck = gameState.Deck

	// Create a new history for this instance.
	instance.history = state.History{}

	instance.dealCards()
}

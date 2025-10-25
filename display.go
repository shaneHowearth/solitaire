package solitaire

import (
	"fmt"
)

// createGamePage - create the game board page dynamically.
func (instance *Instance) createGamePage() {
	// Create the board that will be displayed.
	tableauHeight, tableauWidth := instance.Game.TableauGridSize()
	foundationSpec := instance.Game.Foundations()
	reserveSpec := instance.Game.Reserves()
	howTo := instance.Game.HowToPlay()

	// Create the board layout.
	instance.Display.CreateBoard(
		instance.Game.Name(),
		tableauHeight,
		tableauWidth,
		len(reserveSpec),
		len(foundationSpec),
		howTo,
	)

	// Update the display with current game state.
	instance.updateDisplay()
}

// updateDisplay - update the display with current game state.
func (instance *Instance) updateDisplay() {
	// Tell the board what to display in each box.
	for idx := range instance.State.Foundations {
		// Set the foundation title.
		instance.Display.FoundationTitle(idx,
			fmt.Sprintf(" %s %s ",
				instance.State.Foundations[idx].Base.Rank.String(),
				instance.State.Foundations[idx].Base.Suit.String(),
			),
		)

		// Tell the foundation what cards it is holding.
		instance.Display.FoundationPrint(idx,
			instance.State.Foundations[idx].Stack.Cards(),
		)
	}

	// Tell each Tableau what cards it is holding.
	for idx := range instance.State.Tableau {
		instance.Display.TableauPrint(idx,
			instance.State.Tableau[idx].Stack.Cards(),
			instance.State.Tableau[idx].Stack.ShowCount,
		)
	}

	// Tell each Reserve what cards it is holding.
	for idx := range instance.State.Reserves {
		instance.Display.ReservePrint(idx,
			instance.State.Reserves[idx].Stack.Cards(),
		)
	}

	// Display the Talon.
	instance.Display.TalonPrint(instance.State.Talon.Stock.Cards())

	// Display the Waste.
	instance.Display.WastePrint(instance.State.Talon.Waste.Cards())
}

func (instance *Instance) redeal() {
	instance.Game.Redeal(instance.State.Talon, instance.State.Tableau)
	instance.updateDisplay()
}

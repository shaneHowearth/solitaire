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
	for idx := range instance.Foundations {
		// Set the foundation title.
		instance.Display.FoundationTitle(idx,
			fmt.Sprintf(" %s %s ",
				instance.Foundations[idx].Base.Rank.String(),
				instance.Foundations[idx].Base.Suit.String(),
			),
		)

		// Tell the foundation what cards it is holding.
		instance.Display.FoundationPrint(idx,
			instance.Foundations[idx].Stack.Cards(),
		)
	}

	// Tell each Tableau what cards it is holding.
	for idx := range instance.Tableau {
		instance.Display.TableauPrint(idx,
			instance.Tableau[idx].Stack.Cards(),
		)
	}

	// Tell each Reserve what cards it is holding.
	for idx := range instance.Reserves {
		instance.Display.ReservePrint(idx,
			instance.Reserves[idx].Stack.Cards(),
		)
	}

	// Display the Talon.
	instance.Display.TalonPrint(instance.Talon.Stock.Cards())

	// Display the Waste.
	instance.Display.WastePrint(instance.Talon.Waste.Cards())
}

func (instance *Instance) redeal() {
	instance.Game.Redeal(instance.Talon, instance.Tableau)
	instance.updateDisplay()
}

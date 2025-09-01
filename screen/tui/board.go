//nolint:mnd // Small numeric literals for UI positioning are clearer inline
package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/state"
)

// CreateBoard - Create the board that the game will use.
func (display *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) {
	gamePage := display.createGamePage(name, tableauHeight, tableauWidth, reserveCount, foundationCount, howTo)

	display.screens[name] = gamePage
}

func (display *Display) createGamePage(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) tview.Primitive {
	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)

	title := tview.NewTextView().
		SetText(strings.Join(append([]string{fmt.Sprintf("Playing: %s\n", name)}, howTo...), "\n")).
		SetWordWrap(true).
		SetTextColor(tview.Styles.PrimaryTextColor)
	title.SetBorder(true)

	foundationsRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	// #########
	// # TALON #
	// #########
	display.stack = make([]*tview.TextView, 1)
	talonIndex := 0 // There's typically only one talon.
	talon := tview.NewTextView().SetDynamicColors(true)
	talon.SetWordWrap(true).SetBorder(true).SetTitle(" Stock ")
	talon.SetBackgroundColor(display.defaultBgColor)

	talon.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick && talon.HasFocus() {
			display.selectComponent(state.StackTalon, talonIndex)
			return tview.MouseConsumed, nil
		}

		return action, event
	})

	display.stack[talonIndex] = talon

	foundationsRow.AddItem(
		talon, 0, 1, true,
	)

	// #########
	// # WASTE #
	// #########
	display.waste = make([]*tview.TextView, 1)
	waste := tview.NewTextView().SetDynamicColors(true)

	wasteIndex := 0 // There's typically only one waste pile.

	waste.SetBorder(true).SetTitle(" Waste ")
	waste.SetBackgroundColor(display.defaultBgColor)
	// Add waste selection capability.

	waste.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick && waste.HasFocus() {
			display.selectComponent(state.StackWaste, wasteIndex)
			return tview.MouseConsumed, nil
		}

		return action, event
	})

	display.waste[wasteIndex] = waste
	foundationsRow.AddItem(waste, 0, 1, true)

	// Add a box for each foundation.
	// ###############
	// # FOUNDATIONS #
	// ###############
	display.foundations = make([]*tview.TextView, foundationCount)

	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView().SetDynamicColors(true)

		// Add some decorations to the box.
		foundationIdx := idx

		foundation.SetBorder(true)
		foundation.SetBackgroundColor(display.defaultBgColor)

		foundation.SetMouseCapture(
			func(action tview.MouseAction, event *tcell.EventMouse) (
				tview.MouseAction, *tcell.EventMouse) {
				if action == tview.MouseLeftClick && foundation.HasFocus() {
					display.selectComponent(state.StackFoundation, foundationIdx)
					// Return nil, nil to completely consume the event.
					return tview.MouseConsumed, nil
				}

				return action, event
			})

		display.foundations[foundationIdx] = foundation

		foundationsRow.AddItem(
			foundation, 0, 1, true,
		)
	}

	tableauArea := tview.NewFlex().SetDirection(tview.FlexRow)

	// ############
	// # TABLEAUS #
	// ############
	display.tableau = make([]*tview.TextView, tableauHeight*tableauWidth)

	for idx := 0; idx < tableauHeight; idx++ {
		// Create a row holding flex.
		tableauRow := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Only do this in the first tableauRow
		if idx == 0 {
			// ############
			// # RESERVES #
			// ############

			display.reserves = make([]*tview.TextView, reserveCount)
			for reserveIdx := 0; reserveIdx < reserveCount; reserveIdx++ {
				reserve := tview.NewTextView().SetDynamicColors(true)

				// Add some decorations to the box.
				reserve.SetBorder(true).SetTitle(" Reserve ")
				reserve.SetBackgroundColor(display.defaultBgColor)

				reserve.SetMouseCapture(
					func(action tview.MouseAction, event *tcell.EventMouse) (
						tview.MouseAction, *tcell.EventMouse) {
						if action == tview.MouseLeftClick && reserve.HasFocus() {
							display.selectComponent(state.StackReserve, reserveIdx)
							// Return nil, nil to completely consume the event.
							return tview.MouseConsumed, nil
						}

						return action, event
					})

				display.reserves[reserveIdx] = reserve

				tableauRow.AddItem(
					reserve, 0, 1, true,
				)
			}
		}

		for widthIdx := 0; widthIdx < tableauWidth; widthIdx++ {
			tableau := tview.NewTextView().SetDynamicColors(true)

			tableau.SetBorder(true)
			tableau.SetBackgroundColor(display.defaultBgColor)

			tableauIdx := idx*tableauWidth + widthIdx
			display.tableau[tableauIdx] = tableau

			tableau.SetMouseCapture(
				func(action tview.MouseAction, event *tcell.EventMouse) (
					tview.MouseAction, *tcell.EventMouse) {
					if action == tview.MouseLeftClick && tableau.HasFocus() {
						display.selectComponent(state.StackTableau, tableauIdx)
						return action, nil
					}

					return action, event
				})
			tableauRow.AddItem(tableau, 0, 1, true)

		}

		// Add the row to the tableau.
		tableauArea.AddItem(tableauRow, 0, 1, true)
	}

	// Controls/Help.
	controls := tview.NewTextView().
		SetText("Press 'n' to start a new game, 'm' to return to game menu, Ctrl+C or 'q' to quit").
		SetTextAlign(tview.AlignCenter)
	controls.SetBorder(true).SetTitle("Controls")

	display.App.EnableMouse(true)

	// Add the rows to the main rows container.
	mainRows.
		AddItem(title, 0, 1, false).
		AddItem(foundationsRow, 8, 0, true).
		AddItem(tableauArea, 0, 1, true).
		AddItem(controls, 3, 0, false)

	// Add the main rows to the window container.
	mainRows.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'm':
			display.Show("Games")
			return nil
		case 'n':
			display.onGameSelected(display.Selected)
			return nil
		case 'q':
			display.App.Stop()
			return nil
		}

		return event
	})

	display.screens[name] = mainRows

	return mainRows
}

// TalonPrint -
func (display *Display) TalonPrint(value []string) {
	if len(value) > 0 {
		display.stack[0].SetText(
			value[len(value)-1],
		)
	} else {
		display.stack[0].SetText(
			emptyStack,
		)
	}
}

// WastePrint -
func (display *Display) WastePrint(value []string) {
	if len(value) > 0 {
		textColor := "[-]"
		if strings.Contains(value[len(value)-1], state.Hearts.String()) ||
			strings.Contains(value[len(value)-1], state.Diamonds.String()) {
			textColor = "[red]"
		}
		value[len(value)-1] = fmt.Sprintf("%s%s%s", textColor, value[len(value)-1], "[-]")

		display.waste[0].SetText(
			value[len(value)-1],
		)
	} else {
		display.waste[0].SetText(
			emptyStack,
		)
	}
}

// FoundationTitle -
func (display *Display) FoundationTitle(num int, value string) {
	display.foundations[num].SetTitle(value)
}

// FoundationPrint -
func (display *Display) FoundationPrint(num int, value []string) {
	if len(value) > 0 {
		// Only print the top card.
		textColor := "[-]"
		if strings.Contains(value[len(value)-1], state.Hearts.String()) ||
			strings.Contains(value[len(value)-1], state.Diamonds.String()) {
			textColor = "[red]"
		}
		value[len(value)-1] = fmt.Sprintf("%s%s%s", textColor, value[len(value)-1], "[-]")

		display.foundations[num].SetText(
			value[len(value)-1],
		)
	}
}

const emptyStack = ""

// ReservePrint -
func (display *Display) ReservePrint(idx int, value []string) {
	if len(value) > 0 {
		for i := range value {
			textColor := "[-]"
			if strings.Contains(value[i], state.Hearts.String()) ||
				strings.Contains(value[i], state.Diamonds.String()) {
				textColor = "[red]"
			}
			value[i] = fmt.Sprintf("%s%s%s", textColor, value[i], "[-]")
		}

		display.reserves[idx].SetText(
			strings.Join(value, "\n"),
		)
	} else {
		display.reserves[idx].SetText(
			emptyStack,
		)
	}
}

// TableauPrint -
func (display *Display) TableauPrint(idx int, value []string) {
	if len(value) > 0 {
		for i := range value {
			textColor := "[-]"
			if strings.Contains(value[i], state.Hearts.String()) ||
				strings.Contains(value[i], state.Diamonds.String()) {
				textColor = "[red]"
			}
			value[i] = fmt.Sprintf("%s%s%s", textColor, value[i], "[-]")
		}

		display.tableau[idx].SetText(
			strings.Join(value, "\n"),
		)
	} else {
		display.tableau[idx].SetText(
			emptyStack,
		)
	}
}

// Update the selectComponent method to use the callback:
func (display *Display) selectComponent(componentType state.StackType, index int) {
	if display.processingClick {
		return
	}

	display.processingClick = true

	defer func() {
		display.processingClick = false
	}()

	// Validate the selection first.
	var component *tview.TextView

	switch componentType {
	case state.StackFoundation:
		if index < 0 || index >= len(display.foundations) || display.foundations[index] == nil {
			return
		}

		component = display.foundations[index]
	case state.StackReserve:
		if index < 0 || index >= len(display.reserves) || display.reserves[index] == nil {
			return
		}

		component = display.reserves[index]
	case state.StackTableau:
		if index < 0 || index >= len(display.tableau) || display.tableau[index] == nil {
			return
		}

		component = display.tableau[index]
	case state.StackTalon:
		if index < 0 || index >= len(display.stack) || display.stack[index] == nil {
			return
		}

		component = display.stack[index]
	case state.StackWaste:
		if index < 0 || index >= len(display.waste) || display.waste[index] == nil {
			return
		}

		component = display.waste[index]
	default:
		return
	}

	// Clear previous selection.
	if display.selectedIndex != -1 {
		// Tell the controller.
		display.componentSelectedCallback(display.selectedComponentType, display.selectedIndex, componentType, index)

		display.clearCurrentSelection()
	} else {
		// Set the new selection.
		display.selectedComponentType = componentType
		display.selectedIndex = index
		component.SetBackgroundColor(display.selectedBgColor)

		// Update the display (use a different goroutine to prevent a lockup).
		go func() {
			display.App.Draw() // Call from a goroutine to avoid blocking.
		}()
	}
}

// clearCurrentSelection - helper to clear the current selection.
func (display *Display) clearCurrentSelection() {
	if display.selectedIndex < 0 {
		return
	}

	var component *tview.TextView

	switch display.selectedComponentType {
	case state.StackFoundation:
		if display.selectedIndex < len(display.foundations) && display.foundations[display.selectedIndex] != nil {
			component = display.foundations[display.selectedIndex]
		}
	case state.StackReserve:
		if display.selectedIndex < len(display.reserves) && display.reserves[display.selectedIndex] != nil {
			component = display.reserves[display.selectedIndex]
		}
	case state.StackTableau:
		if display.selectedIndex < len(display.tableau) && display.tableau[display.selectedIndex] != nil {
			component = display.tableau[display.selectedIndex]
		}
	case state.StackTalon:
		if display.selectedIndex < len(display.stack) && display.stack[display.selectedIndex] != nil {
			component = display.stack[display.selectedIndex]
		}
	case state.StackWaste:
		if display.selectedIndex < len(display.waste) && display.waste[display.selectedIndex] != nil {
			component = display.waste[display.selectedIndex]
		}
	default:
		// Shouldn't be here.
		panic(fmt.Sprintf("Bad component type received %d", display.selectedComponentType))
	}

	if component != nil {
		component.SetBackgroundColor(display.defaultBgColor)
		display.selectedComponentType = -1
		display.selectedIndex = -1

		go func() {
			display.App.Draw() // Call from a goroutine to avoid blocking.
		}()
	}
}

// GetSelectedComponent - get the currently selected component type and index.
func (display *Display) GetSelectedComponent() (state.StackType, int) {
	return display.selectedComponentType, display.selectedIndex
}

// ClearSelection - clear the current selection.
func (display *Display) ClearSelection() {
	display.App.QueueUpdate(func() {
		display.clearCurrentSelection()
		display.selectedIndex = -1
	})
}

// HasSelection - check if there's currently a selection.
func (display *Display) HasSelection() bool {
	return display.selectedIndex >= 0
}

package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/state"
)

// CreateBoard - Create the board that the game will use.
func (display *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, foundationCount int,
	foundationBase state.Rank,
) {
	gamePage := display.createGamePage(name, tableauHeight, tableauWidth, foundationCount, foundationBase)

	display.screens[name] = gamePage
}

func (display *Display) createGamePage(
	name string,
	tableauHeight, tableauWidth, foundationCount int,
	foundationBase state.Rank,
) tview.Primitive {
	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)

	title := tview.NewTextView().
		SetText(fmt.Sprintf("Playing: %s", name)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor)
	title.SetBorder(true)

	foundationsRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	/////////////
	/// TALON ///
	/////////////
	talonIndex := 0 // There's typically only one talon
	talon := tview.NewTextView().SetDynamicColors(true)
	talon.SetWordWrap(true).SetBorder(true).SetTitle("Talon")
	talon.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick && talon.HasFocus() {
			display.selectComponent(screen.ComponentTalon, talonIndex)
			return tview.MouseConsumed, nil
		}
		return action, event
	})

	talon.SetBackgroundColor(display.defaultBgColor)
	display.stack = append(display.stack, talon)

	foundationsRow.AddItem(
		talon, 0, 1, true,
	)

	/////////////
	/// WASTE ///
	/////////////
	waste := tview.NewTextView()
	display.waste = append(display.waste, waste)
	waste.SetBorder(true).SetTitle("Waste")
	waste.SetBackgroundColor(display.defaultBgColor)
	// Add waste selection capability
	wasteIndex := 0 // There's typically only one waste pile

	waste.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		if action == tview.MouseLeftClick && waste.HasFocus() {
			display.selectComponent(screen.ComponentWaste, wasteIndex)
			return tview.MouseConsumed, nil
		}
		return action, event
	})

	foundationsRow.AddItem(waste, 0, 1, true)

	// Add a box for each foundation.
	///////////////////
	/// FOUNDATIONS ///
	///////////////////
	display.foundations = make([]*tview.TextView, foundationCount)

	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView()

		// Add some decorations to the box.
		foundationIdx := idx
		foundation.Box.SetBorder(true).SetTitle(foundationBase.String())
		foundation.SetBackgroundColor(display.defaultBgColor)

		// For foundations - replace your current foundation mouse handler:
		foundation.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick && foundation.HasFocus() {
				display.selectComponent(screen.ComponentFoundation, foundationIdx)
				// Return nil, nil to completely consume the event
				return tview.MouseConsumed, nil
			}
			return action, event
		})

		display.foundations[foundationIdx] = foundation

		foundationsRow.AddItem(
			foundation, 0, 1, true,
		)
	}

	// The tableau.
	////////////////
	/// TABLEAUS ///
	////////////////
	tableauArea := tview.NewFlex().SetDirection(tview.FlexColumn)
	display.tableau = make([]*tview.TextView, tableauHeight*tableauWidth)

	for idx := 0; idx < tableauHeight*tableauWidth; idx++ {
		tableau := tview.NewTextView()

		tableau.SetBorder(true).SetTitle(fmt.Sprintf(""))
		tableau.SetBackgroundColor(display.defaultBgColor)
		tableauIdx := idx
		display.tableau[tableauIdx] = tableau

		tableau.SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
			if action == tview.MouseLeftClick && tableau.HasFocus() {
				display.selectComponent(screen.ComponentTableau, tableauIdx)
				return action, nil
			}
			return action, event
		})

		// Add the row to the tableau.
		tableauArea.AddItem(tableau, 0, 1, true)
	}

	// Controls/Help
	controls := tview.NewTextView().
		SetText("Press 'q' to return to game selection, Ctrl+C to quit").
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
		case 'q':
			display.Show("Games")
			return nil
		}
		return event
	})

	display.screens[name] = mainRows

	return mainRows
}

// FoundationTitle -
func (display *Display) FoundationTitle(num int, value string) {
	display.foundations[num].SetTitle(value)
}

// FoundationPrint -
func (display *Display) FoundationPrint(num int, value []string) {
	if len(value) > 0 {
		display.foundations[num].SetText(
			value[len(value)-1],
		)
	}
}

// TableauPrint -
func (display *Display) TableauPrint(idx int, value []string) {
	if len(value) > 0 {
		display.tableau[idx].SetText(
			strings.Join(value, "\n"),
		)
	}
}

// Update the selectComponent method to use the callback:
func (display *Display) selectComponent(componentType screen.ComponentType, index int) {
	if display.processingClick {
		return
	}
	display.processingClick = true
	defer func() {
		display.processingClick = false
	}()

	// Validate the selection first
	var component *tview.TextView
	switch componentType {
	case screen.ComponentFoundation:
		if index < 0 || index >= len(display.foundations) || display.foundations[index] == nil {
			return
		}
		component = display.foundations[index]
	case screen.ComponentTableau:
		if index < 0 || index >= len(display.tableau) || display.tableau[index] == nil {
			return
		}
		component = display.tableau[index]
	case screen.ComponentTalon:
		if index < 0 || index >= len(display.stack) || display.stack[index] == nil {
			return
		}
		component = display.stack[index]
	case screen.ComponentWaste:
		if index < 0 || index >= len(display.waste) || display.waste[index] == nil {
			return
		}
		component = display.waste[index]
	default:
		return
	}

	// Clear previous selection
	if display.selectedIndex > -1 {
		display.clearCurrentSelection()
	} else {
		// Set the new selection
		display.selectedComponentType = componentType
		display.selectedIndex = index
		component.SetBackgroundColor(display.selectedBgColor)
	}
}

// clearCurrentSelection - helper to clear the current selection
func (display *Display) clearCurrentSelection() {
	if display.selectedIndex < 0 {
		return
	}

	var component *tview.TextView
	switch display.selectedComponentType {
	case screen.ComponentFoundation:
		if display.selectedIndex < len(display.foundations) && display.foundations[display.selectedIndex] != nil {
			component = display.foundations[display.selectedIndex]
		}
	case screen.ComponentTableau:
		if display.selectedIndex < len(display.tableau) && display.tableau[display.selectedIndex] != nil {
			component = display.tableau[display.selectedIndex]
		}
	case screen.ComponentTalon:
		if display.selectedIndex < len(display.stack) && display.stack[display.selectedIndex] != nil {
			component = display.stack[display.selectedIndex]
		}
	case screen.ComponentWaste:
		if display.selectedIndex < len(display.waste) && display.waste[display.selectedIndex] != nil {
			component = display.waste[display.selectedIndex]
		}
	}

	if component != nil {
		component.SetBackgroundColor(display.defaultBgColor)
		display.selectedIndex = -1
	}
}

// getComponentName - helper to get component name for display
func (display *Display) getComponentName(componentType screen.ComponentType) string {
	switch componentType {
	case screen.ComponentFoundation:
		return "foundation"
	case screen.ComponentTableau:
		return "tableau"
	case screen.ComponentTalon:
		return "talon"
	case screen.ComponentWaste:
		return "waste"
	default:
		return "component"
	}
}

// GetSelectedComponent - get the currently selected component type and index
func (display *Display) GetSelectedComponent() (screen.ComponentType, int) {
	return display.selectedComponentType, display.selectedIndex
}

// ClearSelection - clear the current selection
func (display *Display) ClearSelection() {
	display.App.QueueUpdate(func() {
		display.clearCurrentSelection()
		display.selectedIndex = -1
	})
}

// HasSelection - check if there's currently a selection
func (display *Display) HasSelection() bool {
	return display.selectedIndex >= 0
}

package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/state"
)

// CreateBoard - Create the board that the game will use.
func (display *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, foundationCount int,
	foundationBase state.Rank,
) {
	display.foundations = make([]*tview.TextView, 0, foundationCount)
	display.tableau = make([]*tview.TextView, 0, tableauHeight*tableauWidth)

	// There are two rows, the top one with the Talon and Foundations, and
	// the bottom one with the Tableaus.

	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)
	// The top row.

	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	display.app.EnableMouse(true)
	// toggle := &struct{ Set bool }{Set: false}
	// Add a box for the stack.
	stack := tview.NewTextView()
	display.stack = append(display.stack, stack)
	stack.SetWordWrap(true).SetBorder(true).SetTitle("Talon").SetFocusFunc(
		func() {
			// if !toggle.Set {
			// stack.SetBackgroundColor(tcell.ColorRed)
			display.stack[0].SetText("Talon")
			display.app.QueueUpdateDraw(func() {})
			// fmt.Println("Talon")
			// } else {
			// 	stack.SetBackgroundColor(tcell.ColorDefault)
			// 	stack.SetText("Talon Unset")
			// }
			// toggle.Set = !toggle.Set
		},
	).SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		// stack.GetFocu
		return action, event
	},
	)

	topRow.AddItem(
		stack, 0, 1, false,
	)

	// Add a box for the waste.
	waste := tview.NewTextView()
	display.waste = append(display.waste, waste)
	waste.SetBorder(true).SetTitle("Waste").SetFocusFunc(
		func() {
			// if !toggle.Set {
			// waste.SetBackgroundColor(tcell.ColorRed)
			// waste.SetText("Waste")
			// } else {
			// 	waste.SetBackgroundColor(tcell.ColorDefault)
			// 	waste.SetText("Waste Unset")
			// }
			// toggle.Set = !toggle.Set
		},
		// ).SetMouseCapture(func(action tview.MouseAction, event *tcell.EventMouse) (tview.MouseAction, *tcell.EventMouse) {
		// 	return action, event
		// }
	)

	topRow.AddItem(
		waste, 0, 1, false,
	)

	// Add a box for each foundation.
	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView()

		display.foundations = append(display.foundations, foundation)

		// Add some decorations to the box.
		foundation.Box.SetBorder(true).SetTitle(foundationBase.String())

		topRow.AddItem(
			foundation, 0, 1, false,
		)
	}

	// Add the top row to the main rows container.
	mainRows.AddItem(topRow, 0, 1, false)

	// The tableau.
	tableau := tview.NewFlex().SetDirection(tview.FlexRow)

	for idx := 0; idx < tableauHeight; idx++ {
		// Create a new row for the tableau.
		tableauRow := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Add columns to the row.
		for colIdx := 0; colIdx < tableauWidth; colIdx++ {
			tableauCell := tview.NewTextView()

			display.tableau = append(display.tableau, tableauCell)
			tableauCell.SetBorder(true)

			tableauRow.AddItem(
				tableauCell, 0, 1, false,
			)
		}

		// Add the row to the tableau.
		tableau.AddItem(tableauRow, 0, 1, false)
	}

	// Add the tableau to the main rows container.
	mainRows.AddItem(tableau, 0, tableauHeight, false)

	mainWindow := tview.NewFlex()

	// Add the main rows to the window container.
	mainWindow.AddItem(mainRows, 0, 2, true)

	display.screens[name] = mainWindow
}

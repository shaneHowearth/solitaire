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
	display.stack = make([]*tview.TextView, 0, 1)
	display.waste = make([]*tview.TextView, 0, 1)
	display.foundations = make([]*tview.TextView, 0, foundationCount)
	display.tableau = make([]*tview.TextView, 0, tableauHeight*tableauWidth)

	// There are two rows, the top one with the Talon and Foundations, and
	// the bottom one with the Tableaus.
	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)

	// The top row.
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	display.App.EnableMouse(true)

	// Add a box for the talon.
	/////////////
	/// TALON ///
	/////////////
	talon := tview.NewTextView().SetDynamicColors(true)
	display.stack = append(display.stack, talon)
	talon.SetWordWrap(true).SetBorder(true).SetTitle("Talon").
		SetFocusFunc(func() {
			switch talon.GetBackgroundColor() {
			case tcell.ColorRed:
				talon.SetBackgroundColor(tcell.ColorDefault)
			default:
				talon.SetBackgroundColor(tcell.ColorRed)
			}
		},
		)

	topRow.AddItem(
		talon, 0, 1, true,
	)

	// Add a box for the waste.
	/////////////
	/// WASTE ///
	/////////////
	waste := tview.NewTextView()
	display.waste = append(display.waste, waste)
	waste.SetBorder(true).SetTitle("Waste").
		SetFocusFunc(func() {
			switch waste.GetBackgroundColor() {
			case tcell.ColorRed:
				waste.SetBackgroundColor(tcell.ColorDefault)
			default:
				waste.SetBackgroundColor(tcell.ColorRed)
			}
		},
		)

	topRow.AddItem(
		waste, 0, 1, true,
	)

	// Add a box for each foundation.
	///////////////////
	/// FOUNDATIONS ///
	///////////////////
	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView()

		display.foundations = append(display.foundations, foundation)

		// Add some decorations to the box.
		foundation.Box.SetBorder(true).SetTitle(foundationBase.String()).
			SetFocusFunc(func() {
				switch foundation.GetBackgroundColor() {
				case tcell.ColorRed:
					foundation.SetBackgroundColor(tcell.ColorDefault)
				default:
					foundation.SetBackgroundColor(tcell.ColorRed)
				}
			},
			)

		topRow.AddItem(
			foundation, 0, 1, true,
		)
	}

	// Add the top row to the main rows container.
	mainRows.AddItem(topRow, 0, 1, false)

	// The tableau.
	////////////////
	/// TABLEAUS ///
	////////////////
	tableau := tview.NewFlex().SetDirection(tview.FlexRow)

	for idx := 0; idx < tableauHeight; idx++ {
		// Create a new row for the tableau.
		tableauRow := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Add columns to the row.
		for colIdx := 0; colIdx < tableauWidth; colIdx++ {
			tableauCell := tview.NewTextView()

			display.tableau = append(display.tableau, tableauCell)
			tableauCell.SetBorder(true)
			tableauCell.SetFocusFunc(func() {
				switch tableauCell.GetBackgroundColor() {
				case tcell.ColorRed:
					tableauCell.SetBackgroundColor(tcell.ColorDefault)
				default:
					tableauCell.SetBackgroundColor(tcell.ColorRed)
				}
			},
			)

			tableauRow.AddItem(
				tableauCell, 0, 1, true,
			)
		}

		// Add the row to the tableau.
		tableau.AddItem(tableauRow, 0, 1, false)
	}

	// Add the tableau to the main rows container.
	mainRows.AddItem(tableau, 0, tableauHeight, false)

	mainWindow := tview.NewFlex()

	// Add the main rows to the window container.
	mainWindow.AddItem(mainRows, 0, 2, false)

	display.screens[name] = mainWindow
}

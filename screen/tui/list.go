package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

type Display struct {
	app *tview.Application
}

func (display *Display) Splash(games []game.Variant) {
	display.app = tview.NewApplication()
	list := tview.NewList()
	for idx, game := range games {
		list.AddItem(game.Name(), "", []rune(fmt.Sprintf("%d", idx+1))[0], func() { display.Board(game) })
	}

	if err := display.app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

func (display *Display) Board(variant game.Variant) {
	// How big is the tableau.
	// How big is the foundation.
	tableauHeight, tableauWidth := variant.TableauGridSize()
	foundationCount, baseRank, _ := variant.Foundations()
	// variant.TableauPosition

	// There are two rows, the top one with the Talon and Foundations, and
	// the bottom one with the Tableaus.

	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)
	// The top row.

	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add a box for the stack.
	topRow.AddItem(
		tview.NewBox().SetBorder(true).SetTitle("Talon"), 0, 1, false,
	)

	// Add a box for the waste.
	topRow.AddItem(
		tview.NewBox().SetBorder(true).SetTitle("Waste"), 0, 1, false,
	)

	// Add a box for each foundation.
	for idx := 0; idx < foundationCount; idx++ {
		topRow.AddItem(
			tview.NewBox().SetBorder(true).SetTitle(baseRank.String()), 0, 1, false,
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
			tableauRow.AddItem(
				tview.NewBox().SetBorder(true).SetTitle("Box"), 0, 1, false,
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

	// Display the window container.
	if err := display.app.SetRoot(mainWindow, true).SetFocus(mainWindow).Run(); err != nil {
		panic(err)
	}
}

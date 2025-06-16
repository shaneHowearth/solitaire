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
	// tableauHeight, tableauWidth := variant.TableauGridSize
	foundationCount, _, _ := variant.Foundations()
	// variant.TableauPosition

	flex := tview.NewFlex()
	flex.AddItem(
		// There are two rows, the top one with the Talon and Foundations, and
		// the bottom one with the Tableaus.
		tview.NewFlex().SetDirection(tview.FlexRow).
			// The top row.
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(
					// This is two wide.
					tview.NewBox().SetBorder(true).SetTitle("Talon"), 0, 2, false,
				).
				AddItem(
					// This is foundationCount wide.
					tview.NewBox().SetBorder(true).SetTitle("Foundations"), 0, foundationCount, false,
				),
				// The flexbox that holds the foundation and talon is one high.
				0, 1, false,
			).
			// The tableau.

			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				).
				AddItem(
					tview.NewBox().SetBorder(true).SetTitle("Tableau"), 0, 1, false,
				),
				0, 3, false,
			),
		// End Tableau.
		0, 2, false,
	)
	// AddItem(tview.NewBox().SetBorder(true).SetTitle("Right (20 cols)"), 20, 1, false)

	if err := display.app.SetRoot(flex, true).SetFocus(flex).Run(); err != nil {
		panic(err)
	}
}

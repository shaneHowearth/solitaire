package tui

// This is the set of methods that are concerned with the list of games that the
// user is able to choose from.

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// CreateGameListPage - Lists all variants of games that the application knows about,
// allowing the user to select which variant to play.
func (display *Display) CreateGameListPage(games []game.Variant) *tview.List {
	list := tview.NewList()

	for idx, game := range games {
		list.AddItem(game.Name(),
			"",
			[]rune(fmt.Sprintf("%d", idx+1))[0],
			func() {
				display.Selected = game
				tableauHeight, tableauWidth := game.TableauGridSize()
				foundationCount, foundationBase, _ := game.Foundations()
				display.CreateBoard(
					game.Name(),
					tableauHeight,
					tableauWidth,
					foundationCount,
					foundationBase,
				)
				display.Show(game.Name())
			},
		)
	}

	// Add a quit option.
	list.AddItem("Quit", "", 'q', func() { display.App.Stop() })

	return list
}

// GetSelected - Get the game that was selected by the user.
func (display *Display) GetSelected() game.Variant {
	return display.Selected
}

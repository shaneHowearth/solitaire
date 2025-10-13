package tui

import (
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// CreateGameListPage - Lists all variants of games that the application knows about,
// allowing the user to select which variant to play.
func (display *Display) createGameListPage(games []game.Variant) *tview.List {
	list := tview.NewList()
	list.SetTitle("Select a Solitaire Game (↑/↓ to navigate, Enter to select)").SetBorder(true)

	for _, g := range games {
		// Create closure to capture the correct game instance
		game := g

		list.AddItem(
			game.Name(),
			"", // No secondary text needed
			0,  // No shortcut key
			func() {
				display.onGameSelected(game)
			},
		)
	}

	// Add a quit option.
	list.AddItem("Quit", "Exit the application", 'q', func() {
		display.App.Stop()
	})

	return list
}

// GetSelected - Get the game that was selected by the user.
func (display *Display) GetSelected() game.Variant {
	return display.Selected
}

// onGameSelected - handle game selection.
func (display *Display) onGameSelected(selectedGame game.Variant) {
	display.Selected = selectedGame
	// Call the callback if set.
	if display.gameSelectedCallback != nil {
		display.gameSelectedCallback(selectedGame)
	}
}

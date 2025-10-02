package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// createGameListPage creates a searchable game list with filtering.
func (display *Display) createGameListPage(games []game.Variant) tview.Primitive {
	// Store all games for filtering
	allGames := games

	// Create the list
	list := tview.NewList()
	list.SetTitle("Select a Solitaire Game").SetBorder(true)
	list.ShowSecondaryText(false)

	// Create search input field
	searchField := tview.NewInputField().
		SetLabel("Filter: ").
		SetFieldWidth(30)

	// Function to populate/update the list based on search
	updateList := func(searchTerm string) {
		list.Clear()
		searchLower := strings.ToLower(searchTerm)

		// Add matching games
		for _, g := range allGames {
			game := g // Capture for closure
			gameName := game.Name()

			// Filter by search term
			if searchTerm == "" || strings.Contains(strings.ToLower(gameName), searchLower) {
				list.AddItem(
					gameName,
					"",
					0,
					func() {
						display.onGameSelected(game)
					},
				)
			}
		}

		// Always add quit option
		list.AddItem("Quit", "Exit the application", 'q', func() {
			display.App.Stop()
		})

		// Update title with count
		if searchTerm == "" {
			list.SetTitle("Select a Solitaire Game (↑/↓ to navigate, Enter to select)")
		} else {
			count := list.GetItemCount() - 1 // Subtract quit option
			list.SetTitle("Select a Solitaire Game (" + searchTerm + " - " +
				fmt.Sprintf("%d matches", count) + ")")
		}
	}

	// Initialize with all games
	updateList("")

	// Update list as user types
	searchField.SetChangedFunc(func(text string) {
		updateList(text)
	})

	// Handle special keys in search field
	searchField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyDown:
			// Move focus to list
			display.App.SetFocus(list)
			return nil
		case tcell.KeyEsc:
			// Clear search
			searchField.SetText("")
			updateList("")
			return nil
		}
		return event
	})

	// Handle keys in list
	list.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			// Any typing moves focus back to search
			display.App.SetFocus(searchField)
			// Don't consume the event, let it go to search field
			return event
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			// Backspace in list moves to search
			display.App.SetFocus(searchField)
			return event
		}
		return event
	})

	// Create flex layout
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(searchField, 1, 0, true).
		AddItem(list, 0, 1, false)

	return flex
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

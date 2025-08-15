// than as named constants.
//
//nolint:mnd // Small numeric literals for UI positioning are clearer inline
package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// ShowWinnerModal displays a custom styled win modal.
func (display *Display) ShowWinnerModal(gameName string, score int) {
	// Create a custom modal using Flex containers.
	winText := tview.NewTextView()
	winText.
		SetText(fmt.Sprintf(
			"🏆 VICTORY! 🏆\n\n"+
				"Game: %s\n"+
				"Final Score: %d\n\n"+
				"🎉 Amazing work! 🎉",
			gameName, score,
		)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tcell.ColorYellow).
		SetBackgroundColor(tcell.ColorDarkGreen)

	winText.SetDynamicColors(true). // Enable dynamic colours for emojis.
					SetWrap(true)

	winText.SetBorder(true).
		SetBorderColor(tcell.ColorGold).
		SetTitle(" 🎊 CONGRATULATIONS 🎊 ")

		// Handle key presses directly in the modal.
	winText.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'n', 'N':
			display.StartNewGame()
			return nil
		case 'm', 'M':
			display.ShowGameMenu()
			return nil
		case 'q', 'Q':
			display.App.Stop()
			return nil
		case 27: // ESC key.
			display.pages.RemovePage("winner")
			return nil
		}

		return event
	})

	// Create buttons manually using a flex.
	buttonFlex := tview.NewFlex().SetDirection(tview.FlexColumn)

	newGameBtn := tview.NewButton("N\u0332ew Game")
	newGameBtn.
		SetSelectedFunc(func() {
			display.StartNewGame()
		}).
		SetBackgroundColor(tcell.ColorDarkBlue)

	menuBtn := tview.NewButton("M\u0332ain Menu")
	menuBtn.
		SetSelectedFunc(func() {
			display.ShowGameMenu()
		}).
		SetBackgroundColor(tcell.ColorDarkBlue)

	quitBtn := tview.NewButton("Q\u0332uit")
	quitBtn.
		SetSelectedFunc(func() {
			display.App.Stop()
		}).
		SetBackgroundColor(tcell.ColorDarkRed)

	// Style buttons.
	newGameBtn.SetLabelColor(tcell.ColorWhite)
	menuBtn.SetLabelColor(tcell.ColorWhite)
	quitBtn.SetLabelColor(tcell.ColorWhite)

	buttonFlex.
		AddItem(tview.NewBox(), 0, 1, false). // Spacer.
		AddItem(newGameBtn, 12, 0, false).
		AddItem(tview.NewBox(), 0, 1, false). // Spacer.
		AddItem(menuBtn, 12, 0, false).
		AddItem(tview.NewBox(), 0, 1, false). // Spacer.
		AddItem(quitBtn, 8, 0, false).
		AddItem(tview.NewBox(), 0, 1, false) // Spacer.

		// Position in top-right corner - most of screen remains visible!
	modalContent := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(winText, 0, 1, true).         // Give more space to text.
		AddItem(tview.NewBox(), 1, 0, false). // Small spacer.
		AddItem(buttonFlex, 3, 0, false)

	modalContainer := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 0, 3, false).AddItem(nil, 0, 3, false). // 75% empty space on left.
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).          // Top spacer.
				AddItem(modalContent, 15, 0, true). // Modal content with fixed height.
				AddItem(nil, 0, 1, false),          // Bottom spacer.
						40, 0, true). // Fixed width for modal.
		AddItem(nil, 0, 1, false) // Right spacer.

	display.pages.AddPage("winner", modalContainer, true, true)
}

// StartNewGame - Start a new game.
func (display *Display) StartNewGame() {
	display.onGameSelected(display.Selected)
	display.pages.RemovePage("winner")
}

// ShowGameMenu - Show the menu of games.
func (display *Display) ShowGameMenu() {
	display.Show("Games")
	display.pages.RemovePage("winner")
}

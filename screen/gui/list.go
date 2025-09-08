package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (display *DisplayGUI) createGameListScreen() fyne.CanvasObject {
	title := widget.NewLabel("Select a Solitaire Game")
	title.Alignment = fyne.TextAlignCenter

	gameButtons := container.NewVBox()
	for _, game := range display.games {
		gameCopy := game // Important: capture the loop variable
		button := widget.NewButton(game.Name(), func() {
			display.Selected = gameCopy.Name()
			if display.gameSelectedCallback != nil {
				display.gameSelectedCallback(gameCopy)
			}
		})
		gameButtons.Add(button)
	}

	content := container.NewBorder(title, nil, nil, nil, gameButtons)
	return content
}

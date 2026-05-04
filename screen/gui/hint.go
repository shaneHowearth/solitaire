package gui

import (
	"fmt"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (d *Display) showHintModal() {
	if d.gameHint == nil {
		return
	}

	moves := d.gameHint()
	if len(moves) == 0 {
		dialog.ShowInformation("Hints", "No moves available! Try drawing from the Stock.", d.Window)
		return
	}

	// Create a container to hold our hint rows
	hintList := container.NewVBox()

	for i, m := range moves {
		// Create a descriptive string for the move
		// e.g., "1) A♠ from Waste → Foundation 1"
		src := d.formatLocation(m.Source)
		dst := d.formatLocation(m.Destination)

		cardStr := m.SourceCardTop.String()
		if m.NumberMoving > 1 {
			cardStr = fmt.Sprintf("%s (+%d cards)", cardStr, m.NumberMoving-1)
		}

		hintLabel := widget.NewLabel(fmt.Sprintf("%d) %s: %s → %s", i+1, cardStr, src, dst))
		hintList.Add(hintLabel)
	}

	styledContent := d.styledModalContent(hintList, 500, 350)

	dialog.ShowCustom("Available Hints", "Close", styledContent, d.Window)
}

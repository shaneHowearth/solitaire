package gui

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (d *Display) showHowToModal(gameName string, howTo []string) {
	fullText := strings.Join(howTo, "\n\n")
	// Create a text label with the instructions
	content := widget.NewLabel(fullText)
	// content := widget.NewRichTextFromMarkdown(strings.Join(howTo, "\n\n"))
	content.Wrapping = fyne.TextWrapWord

	styledContent := d.styledModalContent(content, 600, 450)

	// Show it as a modal dialog
	title := fmt.Sprintf("How to Play: %s", gameName)
	dialog.ShowCustom(title, "Back to Game", styledContent, d.Window)
}

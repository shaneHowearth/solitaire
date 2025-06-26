package tui

import "strings"

// FoundationTitle -
func (display *Display) FoundationTitle(num int, value string) {
	display.foundations[num].SetTitle(value)
}

// FoundationPrint -
func (display *Display) FoundationPrint(num int, value []string) {
	if len(value) > 0 {
		display.foundations[num].SetText(
			value[len(value)-1],
		)
	}
}

// TableauPrint -
func (display *Display) TableauPrint(idx int, value []string) {
	if len(value) > 0 {
		display.tableau[idx].SetText(
			strings.Join(value, "\n"),
		)
	}
}

package tui

import (
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// Display -
type Display struct {
	app         *tview.Application
	stack       []*tview.TextView
	waste       []*tview.TextView
	foundations []*tview.TextView
	tableau     []*tview.TextView
	Selected    game.Variant
	games       []game.Variant
	screens     map[string]tview.Primitive
}

// New - create a new display.
func New(games []game.Variant) *Display {
	app := tview.NewApplication()

	display := &Display{
		app:     app,
		games:   games,
		screens: make(map[string]tview.Primitive),
	}
	display.screens["Games"] = display.CreateGameListPage(games)
	display.app.SetRoot(display.screens["Games"], true).EnableMouse(true)

	return display
}

// Show - show the named screen.
func (display *Display) Show(name string) {
	if err := display.app.SetRoot(display.screens[name], true).
		SetFocus(display.screens[name]).
		EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

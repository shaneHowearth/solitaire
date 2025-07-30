package tui

import (
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// Display -
type Display struct {
	App         *tview.Application
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
		App:     app,
		games:   games,
		screens: make(map[string]tview.Primitive),
	}
	display.screens["Games"] = display.CreateGameListPage(games)
	display.App.SetRoot(display.screens["Games"], true).EnableMouse(true)

	return display
}

// Show - show the named screen.
func (display *Display) Show(name string) {
	display.App.SetRoot(display.screens[name], true)

	display.App.ForceDraw()
}

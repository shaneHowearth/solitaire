package tui

import (
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

// Display -
type Display struct {
	App                  *tview.Application
	stack                []*tview.TextView
	waste                []*tview.TextView
	foundations          []*tview.TextView
	tableau              []*tview.TextView
	Selected             game.Variant
	games                []game.Variant
	screens              map[string]tview.Primitive
	gameSelectedCallback func(game.Variant)
}

// New - create a new display.
func New(games []game.Variant) *Display {
	app := tview.NewApplication()

	display := &Display{
		App:     app,
		games:   games,
		screens: make(map[string]tview.Primitive),
	}

	display.screens["Games"] = display.createGameListPage(games)
	// display.App.SetRoot(display.screens["Games"], true).EnableMouse(true)

	return display
}

func (display *Display) Run() error {
	return display.App.Run()
}

// Show - show the named screen.
func (display *Display) Show(name string) {
	display.App.SetRoot(display.screens[name], true).EnableMouse(true)
}

// SetGameSelectedCallback - set the callback for when a game is selected
func (display *Display) SetGameSelectedCallback(callback func(game.Variant)) {
	display.gameSelectedCallback = callback
}

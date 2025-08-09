package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen"
)

// Display -
type Display struct {
	App                       *tview.Application
	stack                     []*tview.TextView
	waste                     []*tview.TextView
	foundations               []*tview.TextView
	tableau                   []*tview.TextView
	Selected                  game.Variant
	games                     []game.Variant
	screens                   map[string]tview.Primitive
	gameSelectedCallback      func(game.Variant)
	componentSelectedCallback func(screen.ComponentType, int)

	selectedComponentType screen.ComponentType
	selectedIndex         int
	defaultBgColor        tcell.Color
	selectedBgColor       tcell.Color

	processingClick bool
}

// New - create a new display.
func New(games []game.Variant) *Display {
	app := tview.NewApplication()

	display := &Display{
		App:                   app,
		games:                 games,
		screens:               make(map[string]tview.Primitive),
		selectedComponentType: screen.ComponentFoundation,
		selectedIndex:         -1, // No selection initially
		defaultBgColor:        tcell.ColorDefault,
		selectedBgColor:       tcell.ColorRed,
	}

	display.screens["Games"] = display.createGameListPage(games)
	// display.App.SetRoot(display.screens["Games"], true).EnableMouse(true)

	return display
}

// Add this method to the TUI Display:
func (display *Display) SetComponentSelectedCallback(callback func(screen.ComponentType, int)) {
	display.componentSelectedCallback = callback
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

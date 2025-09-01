package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// Display - the struct to hold all of the state that the display needs.
type Display struct {
	App                       *tview.Application
	pages                     *tview.Pages
	stack                     []*tview.TextView
	waste                     []*tview.TextView
	foundations               []*tview.TextView
	reserves                  []*tview.TextView
	tableau                   []*tview.TextView
	Selected                  game.Variant
	games                     []game.Variant
	screens                   map[string]tview.Primitive
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
	componentSelectedCallback func(state.StackType, int, state.StackType, int)

	selectedComponentType state.StackType
	selectedIndex         int
	defaultBgColor        tcell.Color
	selectedBgColor       tcell.Color

	processingClick bool
}

// Initialise pages system.
func (display *Display) initializePages() {
	display.pages = tview.NewPages()
	display.App.SetRoot(display.pages, true)
}

// New - create a new display.
func New(games []game.Variant) *Display {
	app := tview.NewApplication()

	display := &Display{
		App:                   app,
		games:                 games,
		screens:               make(map[string]tview.Primitive),
		selectedComponentType: state.StackFoundation,
		selectedIndex:         -1, // No selection initially.
		defaultBgColor:        tcell.ColorDefault,
		selectedBgColor:       tcell.ColorRed,
	}

	display.initializePages()

	display.screens["Games"] = display.createGameListPage(games)

	return display
}

// SetComponentSelectedCallback - sets the callback to be used when a component
// is selected.
func (display *Display) SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int)) {
	display.componentSelectedCallback = callback
}

// Run - Run the application.
func (display *Display) Run() error {
	return display.App.Run()
}

// Show - show the named screen.
func (display *Display) Show(name string) {
	if screen, exists := display.screens[name]; exists {
		// Add the screen as a page if it doesn't exist.
		if display.pages.HasPage(name) {
			display.pages.RemovePage(name)
		}

		display.pages.AddPage(name, screen, true, true)

		// Switch to the page.
		display.pages.SwitchToPage(name)
		display.App.EnableMouse(true)
	}
}

// SetGameSelectedCallback - set the callback for when a game is selected.
func (display *Display) SetGameSelectedCallback(callback func(game.Variant)) {
	display.gameSelectedCallback = callback
}

func (display *Display) SetGameRedealCallback(callback func()) {
	display.gameRedealCallback = callback
}

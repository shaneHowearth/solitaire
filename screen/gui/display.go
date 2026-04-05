package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

const (
	greenR       = 46
	greenG       = 125
	greenB       = 50
	greenA       = 255
	windowWidth  = 1200
	windowHeight = 900
	// cardWidth    = 100
	// cardHeight = 145
)

// Display - the struct to hold all of the state that the display needs.
type Display struct {
	App    fyne.App
	Window fyne.Window

	// Structural containers for the GUI
	// Tabletop holds the background + the CardLayer
	Tabletop  *fyne.Container
	CardLayer *fyne.Container

	// These replace the slices of *tview.TextView
	// We don't need slices of widgets because we draw dynamically
	// into the CardLayer, but we keep these to track layout metadata.
	foundations []*fyne.Container
	reserves    []*fyne.Container
	tableau     []*fyne.Container
	stack       []*fyne.Container // Talon/Stock
	waste       []*fyne.Container

	Selected game.Variant
	games    []game.Variant

	// screens replaces tview.Pages/screens map
	screens map[string]fyne.CanvasObject

	// Callbacks matching the interface
	gameHint                  func() []state.Move
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
	gameUndoCallback          func()
	componentSelectedCallback func(state.StackType, int, state.StackType, int)

	// Selection state for move logic
	selectedComponentType state.StackType
	selectedIndex         int

	// Standardized colors for the GUI
	defaultBgColor  color.Color
	selectedBgColor color.Color

	processingClick bool
}

// New - create a new Fyne display.
// New - create a new Fyne display.
func New(variants []game.Variant) *Display {
	myApp := app.New()
	window := myApp.NewWindow("Irate Sol")

	// 1. Define the Green Felt color
	greenFelt := color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA}

	d := &Display{
		App:             myApp,
		Window:          window,
		games:           variants,
		screens:         make(map[string]fyne.CanvasObject),
		defaultBgColor:  greenFelt,
		selectedBgColor: color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red for selection
	}

	// 2. Setup the "Initializing" screen with the green background
	bg := canvas.NewRectangle(d.defaultBgColor)
	label := widget.NewLabel("Initializing...")

	// Stack the label on top of the green background
	initContent := container.NewStack(bg, container.NewCenter(label))

	window.SetContent(initContent)

	// 3. Set a taller height to handle the tableau
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	return d
}

// SetComponentSelectedCallback - sets the callback used when a component is selected.
func (display *Display) SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int)) {
	display.componentSelectedCallback = callback
}

// Run - Start the application.
func (display *Display) Run() error {
	display.Window.ShowAndRun()
	return nil
}

// Show - Swaps the root content of the window.
// Replaces display.pages.SwitchToPage(name).
func (d *Display) Show(name string) {
	fmt.Printf("Show called for: %q\n", name)

	// If the engine asks for the menu, but we want to jump to a game:
	if name == "Games" {
		if d.gameSelectedCallback != nil && len(d.games) > 0 {
			fmt.Println("Bypassing menu: selecting first game...")
			// This triggers the engine to call CreateBoard()
			d.gameSelectedCallback(d.games[0])
			return
		}
	}

	screen, exists := d.screens[name]
	if exists {
		fmt.Printf("Displaying screen: %q\n", name)
		d.Window.SetContent(screen)
	} else {
		fmt.Printf("Screen %q still doesn't exist after callback.\n", name)
	}

	fmt.Println("Exiting Show")
}

// SetGameSelectedCallback - set the callback for when a game is selected.
func (display *Display) SetGameSelectedCallback(callback func(game.Variant)) {
	display.gameSelectedCallback = callback
}

// SetGameRedealCallback - sets the redeal logic.
func (display *Display) SetGameRedealCallback(callback func()) {
	display.gameRedealCallback = callback
}

// SetGameUndoCallback - sets the undo logic.
func (display *Display) SetGameUndoCallback(callback func()) {
	display.gameUndoCallback = callback
}

// SetHintsCallback - sets the hint logic.
func (display *Display) SetHintsCallback(callback func() []state.Move) {
	display.gameHint = callback
}

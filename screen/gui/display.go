package gui

import (
	"fmt"
	"image/color"
	"strings"

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
	cardWidth    = 100
	cardHeight   = 145
	verticalFan  = 30
)

type Display struct {
	App       fyne.App
	Window    fyne.Window
	Tabletop  *fyne.Container
	CardLayer *fyne.Container

	// Containers for specific game areas
	talonBox      *fyne.Container
	wasteBox      *fyne.Container
	foundationBox *fyne.Container
	tableauBox    *fyne.Container
	reserveBox    *fyne.Container // Added for Reserve support

	Selected game.Variant
	games    []game.Variant
	screens  map[string]fyne.CanvasObject

	// Store dimensions for index mapping
	tableauWidth  int
	tableauHeight int

	gameHint                  func() []state.Move
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
	gameUndoCallback          func()
	componentSelectedCallback func(state.StackType, int, state.StackType, int)

	selectedComponentType state.StackType
	selectedIndex         int

	defaultBgColor  color.Color
	selectedBgColor color.Color
	processingClick bool

	foundationHints map[int]string
}

func New(variants []game.Variant) *Display {
	myApp := app.New()
	window := myApp.NewWindow("Irate Sol")

	greenFelt := color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA}

	d := &Display{
		App:             myApp,
		Window:          window,
		games:           variants,
		screens:         make(map[string]fyne.CanvasObject),
		defaultBgColor:  greenFelt,
		selectedBgColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		selectedIndex:   -1,
	}

	bg := canvas.NewRectangle(d.defaultBgColor)
	label := widget.NewLabel("Initializing...")
	initContent := container.NewStack(bg, container.NewCenter(label))

	window.SetContent(initContent)
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	return d
}

func (d *Display) Run() error {
	d.Window.ShowAndRun()
	return nil
}

func (d *Display) Show(name string) {
	if name == "Games" && d.gameSelectedCallback != nil && len(d.games) > 0 {
		d.gameSelectedCallback(d.games[0])
		return
	}
	if screen, exists := d.screens[name]; exists {
		d.Window.SetContent(screen)
		d.Window.Content().Refresh()
	}
}

func (d *Display) selectComponent(sType state.StackType, index int) {
	if d.processingClick {
		return
	}
	d.processingClick = true
	defer func() { d.processingClick = false }()

	if d.selectedIndex != -1 {
		// Second click: handle move
		d.componentSelectedCallback(d.selectedComponentType, d.selectedIndex, sType, index)
		d.selectedIndex = -1
		d.selectedComponentType = -1
	} else {
		// First click: select
		d.selectedIndex = index
		d.selectedComponentType = sType
	}

	// Force the specific boxes to refresh their children (the cards)
	d.talonBox.Refresh()
	d.wasteBox.Refresh()
	d.foundationBox.Refresh()
	d.tableauBox.Refresh()
}

func (d *Display) getCardFilename(card string) string {
	if card == "" || card == "--" || card == "🂠" {
		return "assets/cards/PySolFC/back131.gif"
	}
	card = strings.TrimSpace(card)

	rank := "00"
	if strings.Contains(card, "Ace") {
		rank = "01"
	} else if strings.Contains(card, "Jack") {
		rank = "11"
	} else if strings.Contains(card, "Queen") {
		rank = "12"
	} else if strings.Contains(card, "King") {
		rank = "13"
	} else {
		numStr := ""
		for _, r := range card {
			if r >= '0' && r <= '9' {
				numStr += string(r)
			} else if numStr != "" {
				break
			}
		}
		if len(numStr) == 1 {
			rank = "0" + numStr
		} else if len(numStr) == 2 {
			rank = numStr
		}
	}

	suit := "z"
	if strings.Contains(card, "♠") {
		suit = "s"
	} else if strings.Contains(card, "♥") {
		suit = "h"
	} else if strings.Contains(card, "♦") {
		suit = "d"
	} else if strings.Contains(card, "♣") {
		suit = "c"
	}

	return fmt.Sprintf("assets/cards/PySolFC/%s%s.gif", rank, suit)
}

// Callback Setters
func (d *Display) SetComponentSelectedCallback(cb func(state.StackType, int, state.StackType, int)) {
	d.componentSelectedCallback = cb
}
func (d *Display) SetGameSelectedCallback(cb func(game.Variant)) { d.gameSelectedCallback = cb }
func (d *Display) SetGameRedealCallback(cb func())               { d.gameRedealCallback = cb }
func (d *Display) SetGameUndoCallback(cb func())                 { d.gameUndoCallback = cb }
func (d *Display) SetHintsCallback(cb func() []state.Move)       { d.gameHint = cb }

// State Getters
func (d *Display) HasSelection() bool { return d.selectedIndex >= 0 }
func (d *Display) ClearSelection()    { d.selectedIndex = -1; d.selectedComponentType = -1 }
func (d *Display) GetSelectedComponent() (state.StackType, int) {
	return d.selectedComponentType, d.selectedIndex
}

func (d *Display) ShowWinnerModal(gameName string, score int) {
	fmt.Printf("Winner: %s Score: %d\n", gameName, score)
}

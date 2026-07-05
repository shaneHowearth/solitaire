package gui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire"
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

	Selected   game.Variant
	games      []game.Variant
	screens    map[string]fyne.CanvasObject
	titleLabel *widget.Label

	// Store dimensions for index mapping
	tableauWidth  int
	tableauHeight int

	gameHint                   func() []state.Move
	gameSelectedCallback       func(game.Variant)
	gameRedealCallback         func()
	gameUndoCallback           func()
	componentSelectedCallback  func(state.StackType, int, state.StackType, int)
	gameRefreshVisualsCallback func()

	selectedComponentType state.StackType
	selectedIndex         int

	defaultBgColor  color.Color
	selectedBgColor color.Color
	processingClick bool

	foundationHints map[int]string

	redealBtn *widget.Button

	// Zoom properties
	zoomLevel   float32
	baseWidth   float32
	baseHeight  float32
	cardWidth   float32
	cardHeight  float32
	verticalFan float32
}

func New(app fyne.App, variants []game.Variant) *Display {
	// myApp := app.New()
	window := app.NewWindow("Irate Sol")

	greenFelt := color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA}

	d := &Display{
		App:             app,
		Window:          window,
		games:           variants,
		screens:         make(map[string]fyne.CanvasObject),
		defaultBgColor:  greenFelt,
		selectedBgColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		selectedIndex:   -1,
		zoomLevel:       1.0,
		baseWidth:       100.0,
		baseHeight:      145.0,
		cardWidth:       100.0,
		cardHeight:      145.0,
		verticalFan:     30.0,
	}

	if len(variants) > 0 {
		d.Selected = variants[0]
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

	defer d.RefreshAll()

	if d.Selected == nil {
		return
	}

	if sType == state.StackTalon {
		// If Talon() is false (Agnes), we trigger Redeal.
		if !d.Selected.Talon() {
			if d.gameRedealCallback != nil {
				d.gameRedealCallback()
			}
			d.ClearSelection()
			return // EXIT HERE so it doesn't try to "select" the talon
		}
	}

	if d.selectedIndex != -1 {
		if d.componentSelectedCallback != nil {
			d.componentSelectedCallback(d.selectedComponentType, d.selectedIndex, sType, index)
		}
		d.selectedIndex = -1
		d.selectedComponentType = -1
	} else {
		d.selectedIndex = index
		d.selectedComponentType = sType
	}
}

func (d *Display) RefreshAll() {
	if d.talonBox != nil {
		d.talonBox.Refresh()
	}
	if d.wasteBox != nil {
		d.wasteBox.Refresh()
	}
	if d.foundationBox != nil {
		d.foundationBox.Refresh()
	}
	if d.tableauBox != nil {
		d.tableauBox.Refresh()
	}
	if d.reserveBox != nil {
		d.reserveBox.Refresh()
	}
}
func (d *Display) getCardResource(card string) fyne.Resource {
	// 1. Get the relative path using your existing logic
	path := d.getCardFilename(card)

	// 2. Read the bytes from the embedded filesystem
	// Replace 'solitaire' with the actual name of your root package if different
	data, err := solitaire.ResourceFS.ReadFile(path)
	if err != nil {
		// Log the error or return a fallback resource (like the card back)
		// so the app doesn't crash or show empty space
		fallback, _ := solitaire.ResourceFS.ReadFile("assets/cards/PySolFC/back131.gif")
		return fyne.NewStaticResource("fallback.gif", fallback)
	}

	// 3. Return as a Fyne resource
	return fyne.NewStaticResource(path, data)
}

// Keep this as a private helper for the mapping logic
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
	d.LaunchFireworks()
}

func (d *Display) showGamePicker() {
	var pickerDialog dialog.Dialog

	// 1. Load Recently Played from Preferences
	prefs := d.App.Preferences()
	recentStr := prefs.StringWithFallback("recent_games", "")
	recentBox := container.NewHBox()

	if recentStr != "" {
		for _, rName := range strings.Split(recentStr, ",") {
			name := rName
			btn := widget.NewButtonWithIcon(name, theme.HistoryIcon(), func() {
				d.switchToGame(name)
				pickerDialog.Hide()
			})
			recentBox.Add(btn)
		}
	} else {
		recentBox.Add(widget.NewLabel("No recent games played yet."))
	}

	// 2. Build Categorized Accordion
	catGroups := make(map[game.Category][]game.Variant)
	for _, g := range d.games {
		catGroups[g.Category()] = append(catGroups[g.Category()], g)
	}

	accordion := widget.NewAccordion()
	displayOrder := []game.Category{
		game.CatKlondike,
		game.CatSpider,
		game.CatFoundation,
		game.CatPairing,
		game.CatSpecialty,
	}

	for _, cat := range displayOrder {
		games, exists := catGroups[cat]
		if !exists {
			continue
		}

		catContent := container.NewVBox()
		for _, g := range games {
			selectedGame := g // Capture for closure

			// UI Elements for the row
			nameLabel := widget.NewLabelWithStyle(selectedGame.Name(), fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			descLabel := widget.NewLabel(selectedGame.Description())
			descLabel.Wrapping = fyne.TextWrapWord

			// Invisible-ish button to act as the click surface
			selectBtn := widget.NewButton("", func() {
				d.switchToGame(selectedGame.Name())
				pickerDialog.Hide()
			})

			// Stack the labels ON TOP of the button
			// The button fills the stack, making the whole area clickable
			clickableRow := container.NewStack(
				selectBtn,
				container.NewPadded(container.NewVBox(nameLabel, descLabel)),
			)

			catContent.Add(clickableRow)
			catContent.Add(widget.NewSeparator())
		}
		accordion.Append(widget.NewAccordionItem(cat.String(), catContent))
	}

	// 3. Assemble the Dashboard
	scrollArea := container.NewVScroll(accordion)
	scrollArea.SetMinSize(fyne.NewSize(600, 450))

	dashboard := container.NewVBox(
		widget.NewLabelWithStyle("QUICK FIND", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		container.NewHScroll(recentBox),
		widget.NewSeparator(),
		widget.NewLabelWithStyle("GAME LIBRARY", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
		scrollArea,
	)

	// 4. Create and Show the Dialog
	pickerDialog = dialog.NewCustom("Select Your Game", "Cancel", dashboard, d.Window)
	pickerDialog.Show()
}

func (d *Display) formatLocation(stack state.Stack) string {
	switch stack.Type {
	case state.StackTableau:
		// Map flat index back to Row/Col for the user
		row := (stack.TableauPosition / d.tableauWidth) + 1
		col := (stack.TableauPosition % d.tableauWidth) + 1
		rowStr := ""
		if d.tableauHeight > 1 {
			rowStr = fmt.Sprintf("Row %d, Col ", row)
		}
		return fmt.Sprintf("Tableau %s%d", rowStr, col)
	case state.StackFoundation:
		return fmt.Sprintf("Foundation %d", stack.FoundationPosition+1)
	case state.StackReserve:
		return fmt.Sprintf("Reserve")
	case state.StackWaste:
		return "Waste"
	default:
		return "Stock"
	}

	return ""
}

func (d *Display) SetGameRefreshVisualsCallback(cb func()) {
	d.gameRefreshVisualsCallback = cb
}

func (d *Display) SetZoom(level float32) {
	if level < 0.5 {
		level = 0.5
	} else if level > 2.0 {
		level = 2.0
	}

	d.zoomLevel = level
	d.cardWidth = d.baseWidth * d.zoomLevel
	d.cardHeight = d.baseHeight * d.zoomLevel
	d.verticalFan = 30.0 * d.zoomLevel

	// 1. Update the layout of existing cards
	if d.tableauBox != nil {
		d.RecomputeTableauPiles()
		d.tableauBox.Refresh()
	}

	// 2. Trigger a full refresh if needed
	if d.gameRefreshVisualsCallback != nil {
		d.gameRefreshVisualsCallback()
	}

	// 3. Force Fyne to re-layout the window
	if d.Window != nil && d.Window.Content() != nil {
		d.Window.Content().Refresh()
	}
}

func (d *Display) RecomputeTableauPiles() {
	// Walk the tableauBox to find all fanned containers
	for _, row := range d.tableauBox.Objects {
		if rowContainer, ok := row.(*fyne.Container); ok {
			for _, pile := range rowContainer.Objects {
				// Dig into the VBox to find the Stack (where cards live)
				if vbox, ok := pile.(*fyne.Container); ok {
					for _, obj := range vbox.Objects {
						if stack, ok := obj.(*fyne.Container); ok {
							// This is the container with NoLayout
							d.repositionCardsInStack(stack)
						}
					}
				}
			}
		}
	}
}

func (d *Display) repositionCardsInStack(stack *fyne.Container) {
	for i, obj := range stack.Objects {
		if card, ok := obj.(*CardWidget); ok {
			// 1. Resize to new zoom
			card.Resize(fyne.NewSize(d.cardWidth, d.cardHeight))
			// 2. Re-calculate fan position based on NEW d.verticalFan
			yPos := float32(i) * d.verticalFan
			card.Move(fyne.NewPos(0, yPos))
		}
	}
}

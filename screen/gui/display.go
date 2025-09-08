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

// DisplayGUI represents the main game display system
type DisplayGUI struct {
	App     fyne.App
	Window  fyne.Window
	games   []game.Variant
	screens map[string]fyne.CanvasObject

	// Game components
	stock       *PileWidget
	waste       *PileWidget
	foundations []*PileWidget
	reserves    []*PileWidget
	tableau     []*PileWidget

	// Selection state
	selectedComponentType state.StackType
	selectedIndex         int
	processingClick       bool
	Selected              string

	// Callbacks
	componentSelectedCallback func(state.StackType, int, state.StackType, int)
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
}

// New creates a new Fyne-based game display
func New(games []game.Variant) *DisplayGUI {
	myApp := app.New()
	window := myApp.NewWindow("Solitaire Game")

	display := &DisplayGUI{
		App:                   myApp,
		Window:                window,
		games:                 games,
		screens:               make(map[string]fyne.CanvasObject),
		selectedComponentType: -1,
		selectedIndex:         -1,
	}

	display.setupWindow()
	return display
}

// setupKeyboardShortcuts sets up window-level keyboard shortcuts
func (display *DisplayGUI) setupKeyboardShortcuts() {
	display.Window.Canvas().SetOnTypedKey(func(key *fyne.KeyEvent) {
		switch key.Name {
		case fyne.KeyN, fyne.KeyR:
			if display.gameRedealCallback != nil {
				display.gameRedealCallback()
			}
		case fyne.KeyM:
			display.Show("Games")
		case fyne.KeyQ:
			display.App.Quit()
		}
	})
}

// setupWindow configures the main window
func (display *DisplayGUI) setupWindow() {
	display.Window.Resize(fyne.NewSize(1200, 800))
	display.Window.CenterOnScreen()

	// Set up keyboard shortcuts at window level
	display.setupKeyboardShortcuts()

	// Create initial game selection screen
	gameListScreen := display.createGameListScreen()
	display.screens["Games"] = gameListScreen
	display.Show("Games")
}

// CreateBoard creates the game board layout
func (display *DisplayGUI) CreateBoard(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) {
	gamePage := display.createGamePage(name, tableauHeight, tableauWidth, reserveCount, foundationCount, howTo)
	display.screens[name] = gamePage
	display.Show(name)
}

func (display *DisplayGUI) createGamePage(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) fyne.CanvasObject {

	// Create game components
	display.stock = NewPileWidget("Stock", state.StackTalon, 0, display)
	display.waste = NewPileWidget("Waste", state.StackWaste, 0, display)

	// Foundations
	display.foundations = make([]*PileWidget, foundationCount)
	for i := 0; i < foundationCount; i++ {
		display.foundations[i] = NewPileWidget(fmt.Sprintf("Foundation %d", i+1), state.StackFoundation, i, display)
	}

	// Reserves
	if reserveCount > 0 {
		display.reserves = make([]*PileWidget, reserveCount)
		for i := 0; i < reserveCount; i++ {
			display.reserves[i] = NewPileWidget(fmt.Sprintf("Reserve %d", i+1), state.StackReserve, i, display)
		}
	}

	// Tableau
	display.tableau = make([]*PileWidget, tableauWidth)
	for i := 0; i < tableauWidth; i++ {
		display.tableau[i] = NewPileWidget(fmt.Sprintf("Column %d", i+1), state.StackTableau, i, display)
	}

	// Layout
	topRow := container.NewHBox(display.stock, display.waste)

	// Add spacer to separate stock/waste from foundations
	spacer := widget.NewLabel("")
	spacer.Resize(fyne.NewSize(50, 10)) // Fixed width spacer
	topRow.Add(spacer)

	// Add foundations
	for _, foundation := range display.foundations {
		topRow.Add(foundation)
	}

	// Middle section with reserves
	var middleRow *fyne.Container
	if reserveCount > 0 {
		middleRow = container.NewHBox()
		for _, reserve := range display.reserves {
			middleRow.Add(reserve)
		}
	}

	// Tableau row
	tableauRow := container.NewHBox()
	for _, col := range display.tableau {
		tableauRow.Add(col)
	}

	// Controls
	controlsRow := container.NewHBox(
		widget.NewButton("New Game", func() {
			if display.gameRedealCallback != nil {
				display.gameRedealCallback()
			}
		}),
		widget.NewButton("Back to Games", func() {
			display.Show("Games")
		}),
	)

	// Instructions
	instructions := widget.NewLabel(strings.Join(howTo, " "))
	instructions.Wrapping = fyne.TextWrapWord

	// Main layout - use NewVBox and keep the layout manager
	var content *fyne.Container
	if middleRow != nil {
		content = container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Playing: %s", name)),
			instructions,
			topRow,
			middleRow,
			tableauRow,
			controlsRow,
		)
	} else {
		content = container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Playing: %s", name)),
			instructions,
			topRow,
			tableauRow,
			controlsRow,
		)
	}

	// Create a green background rectangle
	greenBackground := canvas.NewRectangle(color.RGBA{0, 128, 0, 255}) // Solitaire table green

	// Layer the background behind the content
	gamePageWithBackground := container.NewBorder(nil, nil, nil, nil, greenBackground, content)

	return gamePageWithBackground
}

// Update methods for different stack types
func (display *DisplayGUI) TalonPrint(value []string) {
	if display.stock != nil {
		if len(value) > 0 {
			display.stock.SetCards([]string{value[len(value)-1]})
		} else {
			display.stock.Clear()
		}
	}
}

func (display *DisplayGUI) WastePrint(value []string) {
	if display.waste != nil {
		if len(value) > 0 {
			text := value[len(value)-1]
			if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
				text = fmt.Sprintf("[red]%s[-]", text)
			}
			display.waste.SetCards([]string{text})
		} else {
			display.waste.Clear()
		}
	}
}

func (display *DisplayGUI) FoundationTitle(num int, value string) {
	if num < len(display.foundations) && display.foundations[num] != nil {
		display.foundations[num].title = value
		display.foundations[num].Refresh()
	}
}

func (display *DisplayGUI) FoundationPrint(num int, value []string) {
	if num < len(display.foundations) && display.foundations[num] != nil {
		if len(value) > 0 {
			texts := make([]string, len(value))
			for i, v := range value {
				text := v
				if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
					text = fmt.Sprintf("[red]%s[-]", text)
				}
				texts[i] = text
			}
			display.foundations[num].SetCards(texts)
		} else {
			display.foundations[num].Clear()
		}
	}
}

func (display *DisplayGUI) ReservePrint(idx int, value []string) {
	if idx < len(display.reserves) && display.reserves[idx] != nil {
		if len(value) > 0 {
			texts := make([]string, len(value))
			for i, v := range value {
				text := v
				if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
					text = fmt.Sprintf("[red]%s[-]", text)
				}
				texts[i] = text
			}
			display.reserves[idx].SetCards(texts)
		} else {
			display.reserves[idx].Clear()
		}
	}
}

func (display *DisplayGUI) TableauPrint(idx int, value []string, _ int) {
	if idx < len(display.tableau) && display.tableau[idx] != nil {
		if len(value) > 0 {
			texts := make([]string, len(value))
			for i, v := range value {
				text := v
				if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
					text = fmt.Sprintf("[red]%s[-]", text)
				}
				texts[i] = text
			}
			display.tableau[idx].SetCards(texts)
		} else {
			display.tableau[idx].Clear()
		}
	}
}

// Selection handling
func (display *DisplayGUI) selectComponent(componentType state.StackType, index int) {
	if display.processingClick {
		return
	}

	display.processingClick = true
	defer func() {
		display.processingClick = false
	}()

	// If something is already selected, try to make a move
	if display.selectedIndex != -1 {
		if display.componentSelectedCallback != nil {
			display.componentSelectedCallback(
				display.selectedComponentType, display.selectedIndex,
				componentType, index,
			)
		}
		display.clearCurrentSelection()
	} else {
		// Make new selection
		display.selectedComponentType = componentType
		display.selectedIndex = index
		display.highlightSelection()
	}
}

func (display *DisplayGUI) highlightSelection() {
	// This would highlight the selected component
	// Implementation depends on how you want to show selection
}

func (display *DisplayGUI) clearCurrentSelection() {
	display.selectedComponentType = -1
	display.selectedIndex = -1
	// Clear any visual selection indicators
}

// Interface methods
func (display *DisplayGUI) SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int)) {
	display.componentSelectedCallback = callback
}

func (display *DisplayGUI) SetGameSelectedCallback(callback func(game.Variant)) {
	display.gameSelectedCallback = callback
}

func (display *DisplayGUI) SetGameRedealCallback(callback func()) {
	display.gameRedealCallback = callback
}

func (display *DisplayGUI) Run() error {
	display.Window.ShowAndRun()
	return nil
}

func (display *DisplayGUI) Show(name string) {
	if screen, exists := display.screens[name]; exists {
		display.Window.SetContent(screen)
	}
}

func (display *DisplayGUI) GetSelectedComponent() (state.StackType, int) {
	return display.selectedComponentType, display.selectedIndex
}

func (display *DisplayGUI) ClearSelection() {
	display.clearCurrentSelection()
}

func (display *DisplayGUI) HasSelection() bool {
	return display.selectedIndex >= 0
}

func (display *DisplayGUI) ShowWinnerModal(winner string, score int) {
	// Show a modal dialog for game completion
	modal := widget.NewModalPopUp(
		container.NewVBox(
			widget.NewLabel(fmt.Sprintf("Congratulations! %s", winner)),
			widget.NewLabel(fmt.Sprintf("Score: %d", score)),
			widget.NewButton("New Game", func() {
				if display.gameRedealCallback != nil {
					display.gameRedealCallback()
				}
			}),
			widget.NewButton("Back to Games", func() {
				display.Show("Games")
			}),
		),
		display.Window.Canvas(),
	)
	modal.Show()
}

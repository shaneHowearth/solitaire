package gui

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	// Import your game package
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// Display represents the main game display system
type Display struct {
	App                   fyne.App
	Window                fyne.Window
	games                 []game.Variant // Replace with your actual game.Variant type
	screens               map[string]fyne.CanvasObject
	container             *container.AppTabs // For managing different screens
	stack                 []*CardArea
	waste                 []*CardArea
	foundations           []*CardArea
	reserves              []*CardArea
	tableau               []*CardArea
	selectedComponentType state.StackType // Replace with your actual state.StackType
	selectedIndex         int
	defaultBgColor        color.Color
	selectedBgColor       color.Color

	processingClick bool
	Selected        string
	// Callbacks
	componentSelectedCallback func(state.StackType, int, state.StackType, int)
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
}

// New creates a new Fyne-based game display
func New(games []game.Variant) *Display {
	myApp := app.New()
	window := myApp.NewWindow("Solitaire Game")

	display := &Display{
		App:                   myApp,
		Window:                window,
		games:                 games,
		screens:               make(map[string]fyne.CanvasObject),
		selectedComponentType: state.StackFoundation,
		selectedIndex:         -1,
		defaultBgColor:        color.RGBA{R: 0, G: 100, B: 0, A: 255}, // Green table
		selectedBgColor:       color.RGBA{R: 255, G: 0, B: 0, A: 255}, // Red selection
	}

	display.initializeScreens()
	display.setupWindow()

	return display
}

// initializeScreens sets up the different screens/pages
func (display *Display) initializeScreens() {
	// Create container for managing different screens
	display.container = container.NewAppTabs()

	// Create the main game screen
	display.screens["Game"] = display.createGameScreen()

	// Create the game selection screen
	display.screens["Games"] = display.createGameListPage(display.games)

	// Create settings screen
	display.screens["Settings"] = display.createSettingsScreen()
}

// createGameScreen creates the main game playing area
func (display *Display) createGameScreen() fyne.CanvasObject {
	// Game table background
	gameTable := canvas.NewRectangle(display.defaultBgColor)

	// Create areas for different card stacks
	foundationArea := display.createFoundationArea()
	tableauArea := display.createTableauArea()
	stockArea := display.createStockArea()

	// Menu bar for game controls
	menuBar := container.NewHBox(
		widget.NewButton("New Game", func() {
			if display.gameRedealCallback != nil {
				display.gameRedealCallback()
			}
		}),
		widget.NewButton("Undo", func() {
			// Implement undo functionality
		}),
		widget.NewButton("Hint", func() {
			// Implement hint functionality
		}),
		widget.NewButton("Back to Games", func() {
			display.Show("Games")
		}),
	)

	// Layout the game screen
	gameContent := container.NewBorder(
		menuBar,        // Top
		nil,            // Bottom
		stockArea,      // Left
		foundationArea, // Right
		tableauArea,    // Center
	)

	// Overlay the content on the game table
	return container.NewMax(gameTable, gameContent)
}

// createFoundationArea creates the foundation piles area
func (display *Display) createFoundationArea() fyne.CanvasObject {
	foundations := container.NewVBox()

	// Create 4 foundation piles (typical for solitaire)
	for i := 0; i < 4; i++ {
		pile := display.createCardPile(fmt.Sprintf("Foundation %d", i), state.StackFoundation, i)
		foundations.Add(pile)
	}

	return foundations
}

// createTableauArea creates the main playing area
func (display *Display) createTableauArea() fyne.CanvasObject {
	tableau := container.NewHBox()

	// Create 7 tableau columns (typical for Klondike)
	for i := 0; i < 7; i++ {
		column := display.createCardPile(fmt.Sprintf("Column %d", i), state.StackType(i+1), i)
		tableau.Add(column)
	}

	return tableau
}

// createStockArea creates the stock and waste piles
func (display *Display) createStockArea() fyne.CanvasObject {
	stockPile := display.createCardPile("Stock", state.StackType(99), 0)
	wastePile := display.createCardPile("Waste", state.StackType(98), 0)

	return container.NewVBox(stockPile, wastePile)
}

// createCardPile creates a clickable card pile area
func (display *Display) createCardPile(name string, stackType state.StackType, index int) fyne.CanvasObject {
	// Create a tappable area for the card pile
	pile := canvas.NewRectangle(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	pile.Resize(fyne.NewSize(80, 120)) // Standard card size

	// Add tap handling
	tappable := container.NewWithoutLayout(pile)

	// You'll need to implement proper tap handling here
	// This is a simplified version
	return tappable
}

// createGameListPage creates the game selection screen
func (display *Display) createGameListPage(games []game.Variant) fyne.CanvasObject {
	var gameNames []string
	for _, game := range games {
		gameNames = append(gameNames, game.Name())
	}

	gameList := widget.NewList(
		func() int { return len(games) },
		func() fyne.CanvasObject { return widget.NewLabel("Game Name") },
		func(id widget.ListItemID, object fyne.CanvasObject) {
			object.(*widget.Label).SetText(games[id].Name())
		},
	)

	gameList.OnSelected = func(id widget.ListItemID) {
		if display.gameSelectedCallback != nil {
			display.gameSelectedCallback(games[id])
		}
		display.Show("Game")
	}

	return container.NewBorder(
		widget.NewLabel("Select a Game"),
		nil,
		nil,
		nil,
		gameList,
	)
}

// createSettingsScreen creates a settings/options screen
func (display *Display) createSettingsScreen() fyne.CanvasObject {
	return container.NewVBox(
		widget.NewLabel("Game Settings"),
		widget.NewCheck("Auto-move cards", func(checked bool) {
			// Handle auto-move setting
		}),
		widget.NewCheck("Show hints", func(checked bool) {
			// Handle hints setting
		}),
		widget.NewButton("Back", func() {
			display.Show("Games")
		}),
	)
}

// setupWindow configures the main window
func (display *Display) setupWindow() {
	display.Window.Resize(fyne.NewSize(1200, 800))
	display.Window.CenterOnScreen()

	// Start with the games list
	display.Show("Games")
}

// SetComponentSelectedCallback sets the callback for component selection
func (display *Display) SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int)) {
	display.componentSelectedCallback = callback
}

// SetGameSelectedCallback sets the callback for game selection
func (display *Display) SetGameSelectedCallback(callback func(game.Variant)) {
	display.gameSelectedCallback = callback
}

// SetGameRedealCallback sets the callback for redealing cards
func (display *Display) SetGameRedealCallback(callback func()) {
	display.gameRedealCallback = callback
}

// Run starts the application
func (display *Display) Run() {
	display.Window.ShowAndRun()
}

// Show displays the named screen
func (display *Display) Show(name string) {
	if screen, exists := display.screens[name]; exists {
		display.Window.SetContent(screen)
	}
}

// Additional helper methods for game state management
func (display *Display) UpdateGameState() {
	// Refresh the game display based on current state
	// This would update card positions, selections, etc.
}

func (display *Display) HighlightComponent(stackType state.StackType, index int) {
	// Highlight a specific component (card pile, etc.)
	// You'd implement visual feedback here
}

// onGameSelected handles game selection
func (display *Display) onGameSelected(gameName string) {
	// Find the game by name and call the callback
	for _, game := range display.games {
		if game.Name() == gameName {
			if display.gameSelectedCallback != nil {
				display.gameSelectedCallback(game)
			}
			break
		}
	}
}

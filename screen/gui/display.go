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

	// Card being dragged.
	draggedCard *CardWidget

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

	// Tableau - create total tableau widgets (width * height)
	totalTableauColumns := tableauWidth * tableauHeight
	display.tableau = make([]*PileWidget, totalTableauColumns)
	for i := 0; i < totalTableauColumns; i++ {
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

	// Tableau rows - distribute columns across multiple rows
	var tableauRows []*fyne.Container

	colIndex := 0
	for row := 0; row < tableauHeight && colIndex < totalTableauColumns; row++ {
		tableauRow := container.NewHBox()

		// Add 'tableauWidth' columns to each row
		for col := 0; col < tableauWidth && colIndex < totalTableauColumns; col++ {
			tableauRow.Add(display.tableau[colIndex])
			colIndex++
		}

		tableauRows = append(tableauRows, tableauRow)
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
	contentItems := []fyne.CanvasObject{
		widget.NewLabel(fmt.Sprintf("Playing: %s", name)),
		instructions,
		topRow,
	}

	// Add middle row if it exists
	if middleRow != nil {
		contentItems = append(contentItems, middleRow)
	}

	// Add all tableau rows
	for _, tableauRow := range tableauRows {
		contentItems = append(contentItems, tableauRow)
	}

	// Add controls at the end
	contentItems = append(contentItems, controlsRow)

	// Create the main content container
	content := container.NewVBox(contentItems...)

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
			display.foundations[num].SetCards(value)
		} else {
			display.foundations[num].Clear()
		}
	}
}

func (display *DisplayGUI) ReservePrint(idx int, value []string) {
	if idx < len(display.reserves) && display.reserves[idx] != nil {
		if len(value) > 0 {
			display.reserves[idx].SetCards(value)
		} else {
			display.reserves[idx].Clear()
		}
	}
}

func (display *DisplayGUI) TableauPrint(idx int, value []string, _ int) {
	if idx < len(display.tableau) && display.tableau[idx] != nil {
		if len(value) > 0 {
			display.tableau[idx].SetCards(value)
		} else {
			display.tableau[idx].Clear()
		}
	}
}

// Selection handling
func (display *DisplayGUI) selectComponent(componentType state.StackType, index int) {
	fmt.Printf("DEBUG: selectComponent called with type=%d, index=%d\n", componentType, index)
	fmt.Printf("DEBUG: Current selection: type=%d, index=%d\n", display.selectedComponentType, display.selectedIndex)

	if display.processingClick {
		fmt.Printf("DEBUG: Already processing click, ignoring\n")
		return
	}

	display.processingClick = true
	defer func() {
		display.processingClick = false
		fmt.Printf("DEBUG: Finished processing click\n")
	}()

	// If something is already selected, try to make a move
	if display.selectedIndex != -1 {
		fmt.Printf("DEBUG: Making move from (%d,%d) to (%d,%d)\n",
			display.selectedComponentType, display.selectedIndex, componentType, index)

		if display.componentSelectedCallback != nil {
			fmt.Printf("DEBUG: Calling componentSelectedCallback\n")
			display.componentSelectedCallback(
				display.selectedComponentType, display.selectedIndex,
				componentType, index,
			)
		} else {
			fmt.Printf("DEBUG: componentSelectedCallback is nil!\n")
		}
		display.clearCurrentSelection()
	} else {
		// Make new selection
		fmt.Printf("DEBUG: Making new selection: type=%d, index=%d\n", componentType, index)
		display.selectedComponentType = componentType
		display.selectedIndex = index
		display.highlightSelection()
	}
}

func (display *DisplayGUI) highlightSelection() {
	fmt.Printf("DEBUG: Highlighting selection: type=%d, index=%d\n",
		display.selectedComponentType, display.selectedIndex)

	// Clear any existing selection first
	display.clearVisualSelection()

	// Find and highlight the selected component
	switch display.selectedComponentType {
	case state.StackTalon:
		if display.stock != nil {
			display.highlightPile(display.stock)
		}
	case state.StackWaste:
		if display.waste != nil {
			display.highlightPile(display.waste)
		}
	case state.StackFoundation:
		if display.selectedIndex < len(display.foundations) && display.foundations[display.selectedIndex] != nil {
			display.highlightPile(display.foundations[display.selectedIndex])
		}
	case state.StackReserve:
		if display.selectedIndex < len(display.reserves) && display.reserves[display.selectedIndex] != nil {
			display.highlightPile(display.reserves[display.selectedIndex])
		}
	case state.StackTableau:
		if display.selectedIndex < len(display.tableau) && display.tableau[display.selectedIndex] != nil {
			display.highlightPile(display.tableau[display.selectedIndex])
		}
	}
}

func (display *DisplayGUI) clearCurrentSelection() {
	fmt.Printf("DEBUG: Clearing selection (was type=%d, index=%d)\n",
		display.selectedComponentType, display.selectedIndex)

	// Clear visual selection from previously selected component
	display.clearVisualSelection()

	display.selectedComponentType = -1
	display.selectedIndex = -1
}

func (display *DisplayGUI) highlightPile(pile *PileWidget) {
	fmt.Printf("DEBUG: Highlighting pile: %s\n", pile.title)

	// Highlight the top card if there are cards
	if len(pile.cards) > 0 {
		topCard := pile.cards[len(pile.cards)-1]
		topCard.SetSelected(true)
	}
}

// Helper method to clear visual selection
func (display *DisplayGUI) clearVisualSelection() {
	fmt.Printf("DEBUG: Clearing visual selection\n")

	// Clear selection from all piles
	allPiles := []*PileWidget{}
	if display.stock != nil {
		allPiles = append(allPiles, display.stock)
	}
	if display.waste != nil {
		allPiles = append(allPiles, display.waste)
	}
	allPiles = append(allPiles, display.foundations...)
	allPiles = append(allPiles, display.tableau...)
	allPiles = append(allPiles, display.reserves...)

	for _, pile := range allPiles {
		if pile == nil {
			continue
		}
		for _, card := range pile.cards {
			card.SetSelected(false)
		}
	}
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

func (display *DisplayGUI) handleCardDrop(card *CardWidget, target fyne.CanvasObject) {
	fmt.Printf("DEBUG: handleCardDrop called\n")

	if pile, ok := target.(*PileWidget); ok {
		fmt.Printf("DEBUG: Drop target is pile: %s\n", pile.title)

		// Use the game's callback to handle the move
		if display.componentSelectedCallback != nil {
			display.componentSelectedCallback(
				card.stackType, card.index,
				pile.stackType, pile.index,
			)
		} else {
			fmt.Printf("DEBUG: No callback set for handling moves\n")
		}
	} else {
		fmt.Printf("DEBUG: Drop target is not a pile: %T\n", target)
	}
}

func (display *DisplayGUI) canMoveCardToPile(card *CardWidget, pile *PileWidget) bool {
	// Implement your solitaire rules here
	switch pile.stackType {
	case state.StackFoundation:
		return display.isValidFoundationMove(card, pile)
	case state.StackTableau:
		return display.isValidTableauMove(card, pile)
	default:
		return false
	}
}

func (display *DisplayGUI) isValidFoundationMove(card *CardWidget, pile *PileWidget) bool {
	// Foundation rules: Same suit, ascending order
	if len(pile.cards) == 0 {
		return strings.HasPrefix(card.cardName, "Ace")
	}
	// Add more validation logic here
	return false
}

func (display *DisplayGUI) isValidTableauMove(card *CardWidget, pile *PileWidget) bool {
	// Tableau rules: Alternating colors, descending order
	if len(pile.cards) == 0 {
		return strings.HasPrefix(card.cardName, "King")
	}
	// Add more validation logic here
	return false
}

func (display *DisplayGUI) moveCardToPile(card *CardWidget, targetPile *PileWidget) {
	// This would typically trigger a callback to your game logic
	if display.componentSelectedCallback != nil {
		display.componentSelectedCallback(
			card.stackType, card.index,
			targetPile.stackType, targetPile.index,
		)
	}
}

func (display *DisplayGUI) findObjectAtPosition(pos fyne.Position) fyne.CanvasObject {
	fmt.Printf("DEBUG: Finding object at position: %v\n", pos)

	// Check all piles to see if position is within their bounds
	allPiles := []*PileWidget{}
	if display.stock != nil {
		allPiles = append(allPiles, display.stock)
	}
	if display.waste != nil {
		allPiles = append(allPiles, display.waste)
	}
	allPiles = append(allPiles, display.foundations...)
	allPiles = append(allPiles, display.tableau...)
	allPiles = append(allPiles, display.reserves...)

	for _, pile := range allPiles {
		if pile == nil {
			continue
		}
		pilePos := pile.Position()
		pileSize := pile.Size()

		if pos.X >= pilePos.X && pos.X <= pilePos.X+pileSize.Width &&
			pos.Y >= pilePos.Y && pos.Y <= pilePos.Y+pileSize.Height {
			fmt.Printf("DEBUG: Found pile at position: %s\n", pile.title)
			return pile
		}
	}

	fmt.Printf("DEBUG: No pile found at position\n")
	return nil
}

// Bring card to front during drag to solve z-order issues
func (display *DisplayGUI) bringCardToFront(card *CardWidget) {
	fmt.Printf("DEBUG: Bringing card to front: %s\n", card.cardName)

	// Store reference to dragged card
	display.draggedCard = card

	// In Fyne, we can't easily change z-order of widgets in containers
	// But we can work around this by ensuring the card is rendered last
	// The card's position will be updated and it should appear on top

	// Force refresh of the entire display to ensure proper rendering order
	if display.Window != nil && display.Window.Content() != nil {
		display.Window.Content().Refresh()
	}
}

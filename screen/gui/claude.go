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
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/state"
)

var _ screen.Display = (*DisplayGUI)(nil)

// CardWidget represents a single card that can be tapped
type CardWidget struct {
	widget.BaseWidget
	rect      *canvas.Rectangle
	text      *canvas.Text
	stackType state.StackType
	index     int
	display   *DisplayGUI
	selected  bool
}

func NewCardWidget(stackType state.StackType, index int, display *DisplayGUI) *CardWidget {
	c := &CardWidget{
		rect:      canvas.NewRectangle(color.RGBA{0, 100, 0, 255}), // Green background
		text:      canvas.NewText("", color.Black),
		stackType: stackType,
		index:     index,
		display:   display,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &cardRenderer{
		card:    c,
		objects: []fyne.CanvasObject{c.rect, c.text},
	}
}

func (c *CardWidget) Tapped(*fyne.PointEvent) {
	c.display.selectComponent(c.stackType, c.index)
}

func (c *CardWidget) SetText(text string) {
	// Handle color formatting - convert [red] tags to actual red color
	if strings.Contains(text, "[red]") {
		text = strings.ReplaceAll(text, "[red]", "")
		text = strings.ReplaceAll(text, "[-]", "")
		c.text.Color = color.RGBA{200, 0, 0, 255} // Red for hearts/diamonds
	} else {
		c.text.Color = color.Black
	}
	c.text.Text = text
	c.text.Refresh()
}

func (c *CardWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		c.rect.FillColor = color.RGBA{255, 255, 0, 255} // Yellow for selection
	} else {
		c.rect.FillColor = color.RGBA{0, 100, 0, 255} // Green default
	}
	c.rect.Refresh()
}

func (c *CardWidget) SetEmpty() {
	c.text.Text = ""
	c.rect.FillColor = color.RGBA{50, 50, 50, 255} // Dark gray for empty
	c.rect.Refresh()
	c.text.Refresh()
}

type cardRenderer struct {
	card    *CardWidget
	objects []fyne.CanvasObject
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.card.rect.Resize(size)
	r.card.text.Move(fyne.NewPos(5, 5))
	r.card.text.Resize(fyne.NewSize(size.Width-10, size.Height-10))
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 80)
}

func (r *cardRenderer) Refresh() {
	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cardRenderer) Destroy() {}

// PileWidget represents a pile of cards
type PileWidget struct {
	widget.BaseWidget
	cards     []*CardWidget
	stackType state.StackType
	index     int
	title     string
	display   *DisplayGUI
}

func NewPileWidget(title string, stackType state.StackType, index int, display *DisplayGUI) *PileWidget {
	p := &PileWidget{
		title:     title,
		stackType: stackType,
		index:     index,
		display:   display,
		cards:     make([]*CardWidget, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *PileWidget) CreateRenderer() fyne.WidgetRenderer {
	titleText := canvas.NewText(p.title, color.Black)
	titleText.Alignment = fyne.TextAlignCenter

	return &pileRenderer{
		pile:    p,
		title:   titleText,
		objects: []fyne.CanvasObject{titleText},
	}
}

func (p *PileWidget) AddCard(text string) {
	card := NewCardWidget(p.stackType, p.index, p.display)
	card.SetText(text)
	p.cards = append(p.cards, card)
	p.Refresh()
}

func (p *PileWidget) SetCards(texts []string) {
	p.cards = make([]*CardWidget, 0, len(texts))
	for _, text := range texts {
		card := NewCardWidget(p.stackType, p.index, p.display)
		card.SetText(text)
		p.cards = append(p.cards, card)
	}
	p.Refresh()
}

func (p *PileWidget) Clear() {
	p.cards = make([]*CardWidget, 0)
	p.Refresh()
}

func (p *PileWidget) Tapped(*fyne.PointEvent) {
	p.display.selectComponent(p.stackType, p.index)
}

type pileRenderer struct {
	pile    *PileWidget
	title   *canvas.Text
	objects []fyne.CanvasObject
}

func (r *pileRenderer) Layout(size fyne.Size) {
	r.title.Move(fyne.NewPos(0, 0))
	r.title.Resize(fyne.NewSize(size.Width, 20))

	// Layout cards vertically with slight offset
	cardHeight := float32(25)
	startY := float32(25)

	for i, card := range r.pile.cards {
		card.Move(fyne.NewPos(0, startY+float32(i)*cardHeight))
		card.Resize(fyne.NewSize(size.Width, 80))
	}
}

func (r *pileRenderer) MinSize() fyne.Size {
	baseHeight := float32(105) // Title + one card
	if len(r.pile.cards) > 1 {
		baseHeight += float32(len(r.pile.cards)-1) * 25 // Overlap
	}
	return fyne.NewSize(80, baseHeight)
}

func (r *pileRenderer) Refresh() {
	// Update objects list with current cards
	r.objects = []fyne.CanvasObject{r.title}
	for _, card := range r.pile.cards {
		r.objects = append(r.objects, card)
	}

	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *pileRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *pileRenderer) Destroy() {}

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

// setupWindow configures the main window
func (display *DisplayGUI) setupWindow() {
	display.Window.Resize(fyne.NewSize(1200, 800))
	display.Window.CenterOnScreen()

	// Create initial game selection screen
	gameListScreen := display.createGameListScreen()
	display.screens["Games"] = gameListScreen
	display.Show("Games")
}

func (display *DisplayGUI) createGameListScreen() fyne.CanvasObject {
	title := widget.NewLabel("Select a Solitaire Game")
	title.Alignment = fyne.TextAlignCenter

	gameButtons := container.NewVBox()
	for _, game := range display.games {
		gameCopy := game // Important: capture the loop variable
		button := widget.NewButton(game.Name(), func() {
			display.Selected = gameCopy.Name()
			if display.gameSelectedCallback != nil {
				display.gameSelectedCallback(gameCopy)
			}
		})
		gameButtons.Add(button)
	}

	content := container.NewBorder(title, nil, nil, nil, gameButtons)
	return content
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
	// Add spacer
	topRow.Add(widget.NewLabel(""))

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

	// Main layout
	content := container.NewVBox(
		widget.NewLabel(fmt.Sprintf("Playing: %s", name)),
		instructions,
		topRow,
	)

	if middleRow != nil {
		content.Add(middleRow)
	}

	content.Add(tableauRow)
	content.Add(controlsRow)

	// Handle keyboard shortcuts
	content = container.NewWithoutLayout(content)

	return content
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

package gui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
	// Import your state package
	// "github.com/shanehowearth/solitaire/state"
)

// CardArea represents a clickable card area in Fyne
type CardArea struct {
	*widget.Card
	background *canvas.Rectangle
	label      *widget.RichText
	stackType  state.StackType
	index      int
	display    *Display
}

// NewCardArea creates a new clickable card area
func NewCardArea(title string, stackType state.StackType, index int, display *Display) *CardArea {
	background := canvas.NewRectangle(display.defaultBgColor)

	label := widget.NewRichText()
	label.Wrapping = fyne.TextWrapWord

	content := container.NewMax(background, label)
	card := widget.NewCard(title, "", content)

	cardArea := &CardArea{
		Card:       card,
		background: background,
		label:      label,
		stackType:  stackType,
		index:      index,
		display:    display,
	}

	// Add tap handling
	// card.OnTapped = func() {
	// 	cardArea.display.selectComponent(stackType, index)
	// }

	return cardArea
}

// SetText updates the card content with color formatting
func (ca *CardArea) SetText(text string) {
	if text == "" {
		ca.label.ParseMarkdown("")
		return
	}

	// Convert color formatting from tview to Fyne rich text
	richText := ca.convertColorFormat(text)
	ca.label.ParseMarkdown(richText)
}

// SetBackgroundColor changes the background color
func (ca *CardArea) SetBackgroundColor(c color.Color) {
	ca.background.FillColor = c
	ca.background.Refresh()
}

// convertColorFormat converts tview color tags to Fyne rich text
func (ca *CardArea) convertColorFormat(text string) string {
	// Convert [red]text[-] to red markdown
	if strings.Contains(text, "[red]") {
		text = strings.ReplaceAll(text, "[red]", "**")
		text = strings.ReplaceAll(text, "[-]", "**")
	}
	return text
}

// CreateBoard creates the game board layout
func (display *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) {
	gamePage := display.createGamePage(name, tableauHeight, tableauWidth, reserveCount, foundationCount, howTo)
	display.screens[name] = gamePage
}

func (display *Display) createGamePage(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) fyne.CanvasObject {

	// Title and instructions
	titleText := fmt.Sprintf("Playing: %s\n\n%s", name, strings.Join(howTo, "\n"))
	title := widget.NewCard("", "", widget.NewRichText(&widget.TextSegment{
		Text:  titleText,
		Style: widget.RichTextStyle{},
	}))

	// Top row: Stock, Waste, Foundations
	topRow := container.NewHBox()

	// STOCK/TALON
	display.stack = make([]*CardArea, 1)
	talon := NewCardArea("Stock", state.StackTalon, 0, display)
	display.stack[0] = talon
	topRow.Add(talon)

	// WASTE
	display.waste = make([]*CardArea, 1)
	waste := NewCardArea("Waste", state.StackWaste, 0, display)
	display.waste[0] = waste
	topRow.Add(waste)

	// Add spacer
	topRow.Add(widget.NewLabel(""))

	// FOUNDATIONS
	display.foundations = make([]*CardArea, foundationCount)
	for idx := 0; idx < foundationCount; idx++ {
		foundation := NewCardArea(fmt.Sprintf("Foundation %d", idx), state.StackFoundation, idx, display)
		display.foundations[idx] = foundation
		topRow.Add(foundation)
	}

	// Middle area: Reserves and Tableau
	middleArea := container.NewVBox()

	// RESERVES (if any)
	if reserveCount > 0 {
		reserveRow := container.NewHBox()
		display.reserves = make([]*CardArea, reserveCount)

		for idx := 0; idx < reserveCount; idx++ {
			reserve := NewCardArea("Reserve", state.StackReserve, idx, display)
			display.reserves[idx] = reserve
			reserveRow.Add(reserve)
		}

		// Add spacers to align with tableau
		for i := reserveCount; i < tableauWidth; i++ {
			reserveRow.Add(widget.NewLabel(""))
		}

		middleArea.Add(reserveRow)
	}

	// TABLEAU
	display.tableau = make([]*CardArea, tableauHeight*tableauWidth)

	for rowIdx := 0; rowIdx < tableauHeight; rowIdx++ {
		tableauRow := container.NewHBox()

		for colIdx := 0; colIdx < tableauWidth; colIdx++ {
			tableauIdx := rowIdx*tableauWidth + colIdx
			tableau := NewCardArea(fmt.Sprintf("Col %d", colIdx), state.StackTableau, tableauIdx, display)
			display.tableau[tableauIdx] = tableau
			tableauRow.Add(tableau)
		}

		middleArea.Add(tableauRow)
	}

	// Controls
	controls := widget.NewCard("Controls", "",
		widget.NewRichText(&widget.TextSegment{
			Text:  "Click to select cards. Press 'N' for new game, 'M' for menu, 'R' to redeal",
			Style: widget.RichTextStyle{},
		}))

	// Main layout
	mainContent := container.NewBorder(
		title,      // Top
		controls,   // Bottom
		nil,        // Left
		nil,        // Right
		middleArea, // Center
	)

	// Wrap in a container that handles keyboard input
	gameContainer := &GameContainer{
		Container: container.NewMax(mainContent),
		display:   display,
	}

	return gameContainer
}

// GameContainer wraps the game content and handles keyboard input
type GameContainer struct {
	*fyne.Container
	display *Display
}

// TypedKey handles keyboard shortcuts
func (gc *GameContainer) TypedKey(key *fyne.KeyEvent) {
	switch key.Name {
	case fyne.KeyN:
		gc.display.onGameSelected(gc.display.Selected)
	case fyne.KeyM:
		gc.display.Show("Games")
	case fyne.KeyQ:
		gc.display.App.Quit()
	case fyne.KeyR:
		if gc.display.gameRedealCallback != nil {
			gc.display.gameRedealCallback()
		}
	}
}

// Update methods for different stack types

// TalonPrint updates the stock pile display
func (display *Display) TalonPrint(value []string) {
	if len(display.stack) > 0 && display.stack[0] != nil {
		if len(value) > 0 {
			display.stack[0].SetText(value[len(value)-1])
		} else {
			display.stack[0].SetText("")
		}
	}
}

// WastePrint updates the waste pile display
func (display *Display) WastePrint(value []string) {
	if len(display.waste) > 0 && display.waste[0] != nil {
		if len(value) > 0 {
			text := value[len(value)-1]
			// Add red color for hearts and diamonds
			if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
				text = fmt.Sprintf("[red]%s[-]", text)
			}
			display.waste[0].SetText(text)
		} else {
			display.waste[0].SetText("")
		}
	}
}

// FoundationTitle updates a foundation's title
func (display *Display) FoundationTitle(num int, value string) {
	if num < len(display.foundations) && display.foundations[num] != nil {
		display.foundations[num].SetTitle(value)
	}
}

// FoundationPrint updates a foundation pile display
func (display *Display) FoundationPrint(num int, value []string) {
	if num < len(display.foundations) && display.foundations[num] != nil {
		if len(value) > 0 {
			text := value[len(value)-1]
			// Add red color for hearts and diamonds
			if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
				text = fmt.Sprintf("[red]%s[-]", text)
			}
			display.foundations[num].SetText(text)
		} else {
			display.foundations[num].SetText("")
		}
	}
}

// ReservePrint updates a reserve pile display
func (display *Display) ReservePrint(idx int, value []string) {
	if idx < len(display.reserves) && display.reserves[idx] != nil {
		if len(value) > 0 {
			var coloredValues []string
			for _, v := range value {
				text := v
				if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
					text = fmt.Sprintf("[red]%s[-]", text)
				}
				coloredValues = append(coloredValues, text)
			}
			display.reserves[idx].SetText(strings.Join(coloredValues, "\n"))
		} else {
			display.reserves[idx].SetText("")
		}
	}
}

// TableauPrint updates a tableau pile display
func (display *Display) TableauPrint(idx int, value []string) {
	if idx < len(display.tableau) && display.tableau[idx] != nil {
		if len(value) > 0 {
			var coloredValues []string
			for _, v := range value {
				text := v
				if strings.Contains(text, "♥") || strings.Contains(text, "♦") {
					text = fmt.Sprintf("[red]%s[-]", text)
				}
				coloredValues = append(coloredValues, text)
			}
			display.tableau[idx].SetText(strings.Join(coloredValues, "\n"))
		} else {
			display.tableau[idx].SetText("")
		}
	}
}

// Selection handling methods

func (display *Display) selectComponent(componentType state.StackType, index int) {
	if display.processingClick {
		return
	}

	display.processingClick = true
	defer func() {
		display.processingClick = false
	}()

	// Get the component
	var component *CardArea

	switch componentType {
	case state.StackFoundation:
		if index >= 0 && index < len(display.foundations) {
			component = display.foundations[index]
		}
	case state.StackReserve:
		if index >= 0 && index < len(display.reserves) {
			component = display.reserves[index]
		}
	case state.StackTableau:
		if index >= 0 && index < len(display.tableau) {
			component = display.tableau[index]
		}
	case state.StackTalon:
		if index >= 0 && index < len(display.stack) {
			component = display.stack[index]
		}
	case state.StackWaste:
		if index >= 0 && index < len(display.waste) {
			component = display.waste[index]
		}
	}

	if component == nil {
		return
	}

	// Handle selection
	if display.selectedIndex != -1 {
		// There's already a selection, try to make a move
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
		component.SetBackgroundColor(display.selectedBgColor)
	}
}

func (display *Display) clearCurrentSelection() {
	if display.selectedIndex < 0 {
		return
	}

	var component *CardArea

	switch display.selectedComponentType {
	case state.StackFoundation:
		if display.selectedIndex < len(display.foundations) {
			component = display.foundations[display.selectedIndex]
		}
	case state.StackReserve:
		if display.selectedIndex < len(display.reserves) {
			component = display.reserves[display.selectedIndex]
		}
	case state.StackTableau:
		if display.selectedIndex < len(display.tableau) {
			component = display.tableau[display.selectedIndex]
		}
	case state.StackTalon:
		if display.selectedIndex < len(display.stack) {
			component = display.stack[display.selectedIndex]
		}
	case state.StackWaste:
		if display.selectedIndex < len(display.waste) {
			component = display.waste[display.selectedIndex]
		}
	}

	if component != nil {
		component.SetBackgroundColor(display.defaultBgColor)
	}

	display.selectedComponentType = -1
	display.selectedIndex = -1
}

// GetSelectedComponent returns the currently selected component
func (display *Display) GetSelectedComponent() (state.StackType, int) {
	return display.selectedComponentType, display.selectedIndex
}

// ClearSelection clears the current selection
func (display *Display) ClearSelection() {
	display.clearCurrentSelection()
}

// HasSelection checks if there's a current selection
func (display *Display) HasSelection() bool {
	return display.selectedIndex >= 0
}

// Add these fields to your Display struct:
// stack       []*CardArea
// waste       []*CardArea
// foundations []*CardArea
// reserves    []*CardArea
// tableau     []*CardArea
// processingClick bool
// selectedComponentType state.StackType
// selectedIndex int

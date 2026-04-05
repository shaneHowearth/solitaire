package gui

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

const (
	cardWidth    = 100
	cardHeight   = 145
	cardPadding  = 20
	verticalFan  = 30
	headerHeight = 120
)

// CreateBoard initializes the GUI components for a specific game variant.
func (d *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) {
	// 1. Create the card layer
	d.CardLayer = container.NewWithoutLayout()

	// 2. Setup the Header (replaces BOX 1 & 2 from TUI)
	nameLabel := widget.NewLabel("Playing: " + name)

	btnBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("New", func() { d.gameRedealCallback() }),
		widget.NewButton("Undo", func() { d.gameUndoCallback() }),
		widget.NewButton("Quit", func() { d.App.Quit() }),
	)

	topRow := container.NewBorder(nil, nil, nameLabel, btnBar)

	// 3. Setup Instructions (replaces BOX 3 from TUI)
	instructionBox := widget.NewLabel(strings.Join(howTo, "\n"))
	instructionBox.Wrapping = fyne.TextWrapWord

	headerContainer := container.NewVBox(topRow, instructionBox)

	// 4. Combine into the Master Layout
	bg := canvas.NewRectangle(d.defaultBgColor)
	gameArea := container.NewStack(bg, d.CardLayer)

	// This mirrors your mainRows.AddItem logic from TUI
	mainLayout := container.NewBorder(headerContainer, nil, nil, nil, gameArea)

	// 5. Store in the screens map
	d.screens[name] = mainLayout

	// 6. Force the window to switch to this new board immediately
	d.Window.SetContent(mainLayout)
	fmt.Printf("Board %q created and set to window content\n", name)
}

// --- Print Methods (Implementing the Interface) ---

func (d *Display) TalonPrint(value []string) {
	d.renderPile(value, 50, 20, state.StackTalon, 0, 1)
}

func (d *Display) WastePrint(value []string) {
	d.renderPile(value, 50+cardWidth+cardPadding, 20, state.StackWaste, 0, 1)
}

func (d *Display) FoundationTitle(num int, value string) {
	// In GUI, we typically use the FoundationPrint to update the visual.
}

func (d *Display) FoundationPrint(num int, value []string) {
	// Foundations usually start after Talon/Waste
	x := float32(350 + (num * (cardWidth + cardPadding)))
	d.renderPile(value, x, 20, state.StackFoundation, num, 1)
}

func (d *Display) ReservePrint(idx int, value []string) {
	x := float32(50 + (idx * (cardWidth + cardPadding)))
	d.renderPile(value, x, 180, state.StackReserve, idx, 1)
}

func (d *Display) TableauPrint(idx int, value []string, showCount int) {
	x := float32(50 + (idx * (cardWidth + cardPadding)))
	y := float32(180)
	// If reserves exist, shift tableau down
	d.renderPile(value, x, y, state.StackTableau, idx, showCount)
}

// --- Internal Rendering Logic ---

func (d *Display) renderPile(cards []string, x, y float32, sType state.StackType, idx int, showCount int) {
	// In a real implementation, we would manage specific objects.
	// For this port, we add objects to the CardLayer.

	if len(cards) == 0 {
		slot := d.createEmptySlot(x, y, sType, idx)
		d.CardLayer.Add(slot)
		return
	}

	for i, cardStr := range cards {
		// Only offset if it's the Tableau
		yOffset := float32(0)
		if sType == state.StackTableau {
			yOffset = float32(i * verticalFan)
		}

		// Determine if face up
		displayStr := "🂠"
		if i >= len(cards)-showCount || showCount == 0 {
			displayStr = cardStr
		}

		// Check if this component is selected to apply highlight
		isSel := (d.selectedComponentType == sType && d.selectedIndex == idx && i == len(cards)-1)

		c := d.createCard(displayStr, x, y+yOffset, isSel, func() {
			d.selectComponent(sType, idx)
		})
		d.CardLayer.Add(c)
	}
	d.CardLayer.Refresh()
}

func (d *Display) createCard(val string, x, y float32, selected bool, tapped func()) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.White)
	if selected {
		rect.StrokeColor = d.selectedBgColor
		rect.StrokeWidth = 3
	} else {
		rect.StrokeColor = color.Black
		rect.StrokeWidth = 1
	}
	rect.Resize(fyne.NewSize(cardWidth, cardHeight))

	text := canvas.NewText(val, color.Black)
	if strings.Contains(val, "♥") || strings.Contains(val, "♦") {
		text.Color = color.RGBA{200, 0, 0, 255}
	}
	text.Alignment = fyne.TextAlignCenter

	// Wrap in a button-like tappable container
	btn := widget.NewButton("", tapped)

	card := container.NewStack(rect, text, btn)
	card.Move(fyne.NewPos(x, y))
	card.Resize(fyne.NewSize(cardWidth, cardHeight))
	return card
}

func (d *Display) createEmptySlot(x, y float32, sType state.StackType, idx int) fyne.CanvasObject {
	rect := canvas.NewRectangle(color.Transparent)
	rect.StrokeColor = color.RGBA{255, 255, 255, 50}
	rect.StrokeWidth = 1
	rect.Resize(fyne.NewSize(cardWidth, cardHeight))

	btn := widget.NewButton("", func() { d.selectComponent(sType, idx) })

	slot := container.NewStack(rect, btn)
	slot.Move(fyne.NewPos(x, y))
	slot.Resize(fyne.NewSize(cardWidth, cardHeight))
	return slot
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
	// Note: Refreshing happens in the next updateDisplay cycle called by Instance
}

func (d *Display) HasSelection() bool { return d.selectedIndex >= 0 }
func (d *Display) ClearSelection()    { d.selectedIndex = -1; d.selectedComponentType = -1 }
func (d *Display) GetSelectedComponent() (state.StackType, int) {
	return d.selectedComponentType, d.selectedIndex
}
func (d *Display) ShowWinnerModal(gameName string, score int) {
	// Implementation using dialog.ShowInformation...
}

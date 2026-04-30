package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

func (d *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, reserveCount, foundationCount int,
	howTo []string,
) {
	// Store dimensions for TableauPrint calculations
	d.tableauWidth = tableauWidth
	d.tableauHeight = tableauHeight

	d.talonBox = container.NewHBox()
	d.wasteBox = container.NewHBox()
	d.foundationBox = container.NewHBox()
	d.reserveBox = container.NewHBox()
	d.tableauBox = container.NewVBox()

	// 1. Initialize Talon and Waste
	d.talonBox.Add(d.buildPile(nil, state.StackTalon, 0, 0, ""))
	d.wasteBox.Add(d.buildPile(nil, state.StackWaste, 0, 0, ""))

	// 2. Initialize Foundations
	for i := 0; i < foundationCount; i++ {
		d.foundationBox.Add(d.buildPile(nil, state.StackFoundation, i, 0, "A"))
	}

	// 3. Initialize Reserves
	for i := 0; i < reserveCount; i++ {
		d.reserveBox.Add(d.buildPile(nil, state.StackReserve, i, 0, ""))
	}

	// 4. Build Tableau Rows
	for h := 0; h < tableauHeight; h++ {
		row := container.NewHBox()
		for w := 0; w < tableauWidth; w++ {
			tableauIdx := h*tableauWidth + w
			row.Add(d.buildPile(nil, state.StackTableau, tableauIdx, 0, ""))
		}
		d.tableauBox.Add(row)
	}

	// 5. UI Component Setup (Select, Buttons)
	var gameNames []string
	for _, g := range d.games {
		gameNames = append(gameNames, g.Name())
	}

	gameSelector := widget.NewSelect(gameNames, func(selected string) {
		if selected == name {
			return
		}
		for _, g := range d.games {
			if g.Name() == selected {
				d.Selected = g
				d.gameSelectedCallback(g)
				return
			}
		}
	})
	gameSelector.SetSelected(name)

	d.redealBtn = widget.NewButton("Redeal", func() {
		if d.gameRedealCallback != nil {
			d.gameRedealCallback()
			d.RefreshAll()
		}
	})
	d.redealBtn.Disable()

	btnBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("How to Play", func() { d.showHowToModal(name, howTo) }),
		widget.NewButton("Hint", func() { d.showHintModal() }),
		widget.NewButton("New", func() {
			d.gameRedealCallback()
			d.ClearBoard()
			for _, g := range d.games {
				if g.Name() == name {
					d.gameSelectedCallback(g)
					break
				}
			}
		}),
		d.redealBtn,
		widget.NewButton("Undo", func() { d.gameUndoCallback() }),
		widget.NewButton("Quit", func() { d.App.Quit() }),
	)

	// --- FIX: Header Declaration ---
	header := container.NewVBox(
		container.NewBorder(nil, nil, gameSelector, btnBar),
	)

	d.foundationHints = make(map[int]string)

	// 6. Layout Top Area
	topArea := container.NewBorder(nil, nil, container.NewHBox(d.talonBox, d.wasteBox), nil, d.foundationBox)

	// 7. Delineation Logic for Reserves
	var mainBoard fyne.CanvasObject
	if reserveCount > 0 {
		reserveSection := container.NewVBox(
			widget.NewLabelWithStyle("RESERVES", fyne.TextAlignLeading, fyne.TextStyle{Italic: true}),
			d.reserveBox,
			canvas.NewLine(color.RGBA{R: 255, G: 255, B: 255, A: 50}),
		)
		mainBoard = container.NewVBox(container.NewPadded(reserveSection), d.tableauBox)
	} else {
		mainBoard = d.tableauBox
	}

	tableauScroll := container.NewScroll(mainBoard)

	// Construct Tabletop
	d.Tabletop = container.NewBorder(topArea, nil, nil, nil, tableauScroll)

	bg := canvas.NewRectangle(d.defaultBgColor)
	d.CardLayer = container.NewStack(bg, d.Tabletop)

	// Set the Window Content
	d.Window.SetContent(container.NewBorder(header, nil, nil, nil, d.CardLayer))
}

func (d *Display) buildPile(cards []string, sType state.StackType, idx int, showCount int, baseCard string) fyne.CanvasObject {
	// 1. Handle Empty Piles (Placeholders)
	if len(cards) == 0 {
		c := NewCardWidget(baseCard, sType, idx, d)
		c.IsPlaceholder = true
		return container.NewVBox(c, layout.NewSpacer())
	}

	// 2. Handle Multi-Row Tableau (Foundation-style: Single card, no fan)
	if d.tableauHeight > 1 && sType == state.StackTableau {
		topFace := "--"
		// If showCount allows it, show the actual face of the last card
		if showCount == 0 || len(cards) > 0 {
			topFace = cards[len(cards)-1]
		}

		cWidget := NewCardWidget(topFace, sType, idx, d)
		cWidget.Resize(fyne.NewSize(float32(cardWidth), float32(cardHeight)))
		// No Move() call needed, it defaults to (0,0)

		return container.NewVBox(cWidget, layout.NewSpacer())
	}

	// 3. Handle Standard Fanned Tableau (Klondike-style) or other stacks (Waste/Talon)
	cardContainer := container.NewWithoutLayout()
	currentHeight := float32(cardHeight)

	for i, face := range cards {
		cardFace := "--"
		if showCount == 0 || i >= len(cards)-showCount {
			cardFace = face
		}

		cWidget := NewCardWidget(cardFace, sType, idx, d)
		cWidget.Resize(fyne.NewSize(float32(cardWidth), float32(cardHeight)))

		yPos := float32(0)
		if sType == state.StackTableau {
			yPos = float32(i * verticalFan)
		}
		cWidget.Move(fyne.NewPos(0, yPos))

		if yPos+float32(cardHeight) > currentHeight {
			currentHeight = yPos + float32(cardHeight)
		}
		cardContainer.Add(cWidget)
	}

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(float32(cardWidth), currentHeight))
	return container.NewVBox(container.NewStack(spacer, cardContainer), layout.NewSpacer())
}

func (d *Display) TableauPrint(idx int, value []string, showCount int) {
	newPile := d.buildPile(value, state.StackTableau, idx, showCount, "")

	rowIdx := idx / d.tableauWidth
	colInRowIdx := idx % d.tableauWidth

	// Simpler logic: rowIdx now directly matches the VBox index
	if rowIdx < len(d.tableauBox.Objects) {
		if row, ok := d.tableauBox.Objects[rowIdx].(*fyne.Container); ok {
			if colInRowIdx < len(row.Objects) {
				row.Objects[colInRowIdx] = newPile
				row.Refresh()
			}
		}
	}
}

// FoundationPrint - Handles the foundations in the top-right
func (d *Display) FoundationPrint(num int, value []string) {
	// Use the hint from the engine if available, otherwise default to "A"
	target := "A"
	if hint, ok := d.foundationHints[num]; ok {
		target = "A" + hint
	}

	newPile := d.buildPile(value, state.StackFoundation, num, 1, target)

	if num < len(d.foundationBox.Objects) {
		d.foundationBox.Objects[num] = newPile
		d.foundationBox.Refresh()
	} else {
		d.foundationBox.Add(newPile)
	}
}

func (d *Display) TalonPrint(value []string) {
	newPile := d.buildPile(value, state.StackTalon, 0, 1, "")
	if len(d.talonBox.Objects) > 0 {
		d.talonBox.Objects[0] = newPile
	} else {
		d.talonBox.Add(newPile)
	}
	d.talonBox.Refresh()
}

// SetRedealStatus allows the controller to toggle the button state
func (d *Display) SetRedealStatus(available bool) {
	if d.redealBtn == nil {
		return
	}
	if available {
		d.redealBtn.Enable()
	} else {
		d.redealBtn.Disable()
	}
}

func (d *Display) WastePrint(value []string) {
	newPile := d.buildPile(value, state.StackWaste, 0, 1, "")
	if len(d.wasteBox.Objects) > 0 {
		d.wasteBox.Objects[0] = newPile
	} else {
		d.wasteBox.Add(newPile)
	}
	d.wasteBox.Refresh()
}

func (d *Display) ClearBoard() {
	d.talonBox.Objects = nil
	d.wasteBox.Objects = nil
	d.foundationBox.Objects = nil
	d.tableauBox.Objects = nil

	d.talonBox.Refresh()
	d.wasteBox.Refresh()
	d.foundationBox.Refresh()
	d.tableauBox.Refresh()
}

// ReservePrint - New function to handle the Reserve stacks you added to the layout
func (d *Display) ReservePrint(idx int, value []string) {
	// Reserves usually show all cards or just the top one depending on the game
	newPile := d.buildPile(value, state.StackReserve, idx, 0, "")

	if idx < len(d.reserveBox.Objects) {
		d.reserveBox.Objects[idx] = newPile
		d.reserveBox.Refresh()
	} else {
		d.reserveBox.Add(newPile)
	}
}

func (d *Display) FoundationTitle(num int, value string) {
	if d.foundationHints == nil {
		d.foundationHints = make(map[int]string)
	}
	d.foundationHints[num] = value

	// Re-print the foundation to show the new hint immediately
	d.FoundationPrint(num, nil)
}

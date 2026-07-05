package gui

import (
	"image/color"
	"strings"

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

	// --- 5. UI Component Setup (Game Library & Buttons) ---

	gameLibraryBtn := widget.NewButton("Change Game", func() {
		d.showGamePicker()
	})

	d.redealBtn = widget.NewButton("Redeal", func() {
		if d.gameRedealCallback != nil {
			d.gameRedealCallback()
			d.RefreshAll()
		}
	})
	d.redealBtn.Disable()

	btnBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("-", func() { d.SetZoom(d.zoomLevel - 0.1) }), // Zoom Out
		widget.NewButton("+", func() { d.SetZoom(d.zoomLevel + 0.1) }), // Zoom In
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

	d.titleLabel = widget.NewLabelWithStyle("Playing "+name, fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Italic: true})

	header := container.NewVBox(
		container.NewBorder(nil, nil, gameLibraryBtn, btnBar, d.titleLabel),
		canvas.NewLine(color.RGBA{R: 255, G: 255, B: 255, A: 50}),
	)

	d.foundationHints = make(map[int]string)

	// 6. Layout Top Area
	// Calculate the gap between Talon/Waste and Foundations (1/8 card width)
	gapSize := d.cardWidth / 8
	gapSpacer := canvas.NewRectangle(color.Transparent)
	gapSpacer.SetMinSize(fyne.NewSize(gapSize, 1))

	// Group Talon/Waste with the gap
	leftTopSection := container.NewHBox(d.talonBox, d.wasteBox, gapSpacer)

	// Foundations are set to Center, allowing the gap to push them right
	topArea := container.NewBorder(nil, nil, leftTopSection, nil, d.foundationBox)

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

	// --- Global Padding (Left and Top) ---
	padSize := d.cardWidth / 4

	leftSpacer := canvas.NewRectangle(color.Transparent)
	leftSpacer.SetMinSize(fyne.NewSize(padSize, 1))

	topSpacer := canvas.NewRectangle(color.Transparent)
	topSpacer.SetMinSize(fyne.NewSize(1, padSize))

	// Assemble the entire board content
	gameBoard := container.NewBorder(topArea, nil, nil, nil, tableauScroll)

	// Wrap the entire game board in a border with both global spacers
	d.Tabletop = container.NewBorder(topSpacer, nil, leftSpacer, nil, gameBoard)
	// ------------------------------------------

	bg := canvas.NewRectangle(d.defaultBgColor)
	d.CardLayer = container.NewStack(bg, d.Tabletop)

	// Set the Window Content
	d.Window.SetContent(container.NewBorder(header, nil, nil, nil, d.CardLayer))
}

func (d *Display) switchToGame(name string) {
	for _, g := range d.games {
		if g.Name() == name {
			// 1. Update the metadata and persistence
			d.recordGamePlayed(name)
			d.Selected = g

			if d.titleLabel != nil {
				d.titleLabel.SetText("Playing " + name)
			}

			// 2. CRITICAL: Clear the existing UI objects
			d.ClearBoard()

			// 3. TRIGGER THE ENGINE
			// This is the bridge to your main.go or controller.
			// It calls the variant's Setup() to generate the new deck/piles.
			if d.gameSelectedCallback != nil {
				d.gameSelectedCallback(g)
			}

			// 4. Refresh the visuals
			d.RefreshAll()
			return
		}
	}
}

func (d *Display) recordGamePlayed(name string) {
	p := d.App.Preferences()
	recent := p.StringWithFallback("recent_games", "")

	var games []string
	if recent != "" {
		games = strings.Split(recent, ",")
	}

	// Move current game to the front
	newRecent := []string{name}
	for _, g := range games {
		// Only add if it's not the current game AND not an empty string
		if g != name && g != "" && len(newRecent) < 5 {
			newRecent = append(newRecent, g)
		}
	}

	p.SetString("recent_games", strings.Join(newRecent, ","))
}

func (d *Display) buildPile(cards []string, sType state.StackType, idx int, showCount int, baseCard string) fyne.CanvasObject {
	// 1. Handle Empty Piles
	if len(cards) == 0 {
		c := NewCardWidget(baseCard, sType, idx, d)
		c.IsPlaceholder = true
		return container.NewVBox(c, layout.NewSpacer())
	}

	// 2. Determine if this specific stack should fan
	// Foundations and Talon almost never fan in this engine.
	// Tableau only fans if the Game Variant says so.
	shouldFan := sType == state.StackTableau && d.Selected.Fanned()

	if !shouldFan {
		topFace := cards[len(cards)-1]
		cWidget := NewCardWidget(topFace, sType, idx, d)
		cWidget.Resize(fyne.NewSize(float32(d.cardWidth), float32(d.cardHeight)))
		return container.NewVBox(cWidget, layout.NewSpacer())
	}

	// 3. Handle Fanned Logic
	cardContainer := container.NewWithoutLayout()
	currentHeight := float32(d.cardHeight)

	for i, face := range cards {
		cardFace := "--"
		if showCount == 0 || i >= len(cards)-showCount {
			cardFace = face
		}

		cWidget := NewCardWidget(cardFace, sType, idx, d)
		cWidget.Resize(fyne.NewSize(float32(d.cardWidth), float32(d.cardHeight)))

		// Calculate vertical offset
		yPos := float32(i) * d.verticalFan
		cWidget.Move(fyne.NewPos(0, yPos))

		if yPos+float32(d.cardHeight) > currentHeight {
			currentHeight = yPos + float32(d.cardHeight)
		}
		cardContainer.Add(cWidget)
	}

	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(float32(d.cardWidth), currentHeight))
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

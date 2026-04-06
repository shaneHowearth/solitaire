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
	d.talonBox = container.NewHBox()
	d.wasteBox = container.NewHBox()
	d.foundationBox = container.NewHBox()
	d.tableauBox = container.NewHBox()

	// Updated calls to buildPile with the 5th argument (empty string for placeholders)
	d.talonBox.Add(d.buildPile(nil, state.StackTalon, 0, 0, ""))
	d.wasteBox.Add(d.buildPile(nil, state.StackWaste, 0, 0, ""))

	for i := 0; i < foundationCount; i++ {
		// Foundations usually start expecting an Ace
		d.foundationBox.Add(d.buildPile(nil, state.StackFoundation, i, 0, "A"))
	}
	for i := 0; i < tableauWidth; i++ {
		d.tableauBox.Add(d.buildPile(nil, state.StackTableau, i, 0, ""))
	}

	// Dropdown Menu Setup
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
				d.gameSelectedCallback(g)
				return
			}
		}
	})
	gameSelector.SetSelected(name)

	btnBar := container.NewHBox(
		layout.NewSpacer(),
		widget.NewButton("New", func() {
			d.gameRedealCallback()
			d.ClearBoard()
			// Re-trigger the game selection to force a fresh draw
			for _, g := range d.games {
				if g.Name() == name {
					d.gameSelectedCallback(g)
					break
				}
			}
		}),
		widget.NewButton("Undo", func() { d.gameUndoCallback() }),
		widget.NewButton("Quit", func() { d.App.Quit() }),
	)

	instructionBox := widget.NewLabel(strings.Join(howTo, "\n"))
	instructionBox.Wrapping = fyne.TextWrapWord

	header := container.NewVBox(
		container.NewBorder(nil, nil, gameSelector, btnBar),
		instructionBox,
	)

	topArea := container.NewBorder(nil, nil, container.NewHBox(d.talonBox, d.wasteBox), nil, d.foundationBox)
	tableauContainer := container.NewHScroll(d.tableauBox)

	d.Tabletop = container.NewVBox(
		topArea,
		container.NewPadded(tableauContainer),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(d.defaultBgColor)
	d.CardLayer = container.NewStack(bg, d.Tabletop)

	d.Window.SetContent(container.NewBorder(header, nil, nil, nil, d.CardLayer))
}

func (d *Display) buildPile(cards []string, sType state.StackType, idx int, showCount int, baseCard string) fyne.CanvasObject {
	if len(cards) == 0 {
		c := NewCardWidget(baseCard, sType, idx, d)
		c.IsPlaceholder = true // Now valid because we updated the struct
		return container.NewVBox(c, layout.NewSpacer())
	}

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

	pileStack := container.NewStack(spacer, cardContainer)
	return container.NewVBox(pileStack, layout.NewSpacer())
}

func (d *Display) TableauPrint(idx int, value []string, showCount int) {
	newPile := d.buildPile(value, state.StackTableau, idx, showCount, "")
	if idx < len(d.tableauBox.Objects) {
		d.tableauBox.Objects[idx] = newPile
	} else {
		d.tableauBox.Add(newPile)
	}
	d.tableauBox.Refresh()
}

func (d *Display) FoundationPrint(num int, value []string) {
	// Map the foundation index to a suit for the hint
	suits := []string{"♠", "♥", "♣", "♦"}
	target := "A"
	if num < len(suits) {
		target = "A" + suits[num]
	}

	newPile := d.buildPile(value, state.StackFoundation, num, 1, target)
	if num < len(d.foundationBox.Objects) {
		d.foundationBox.Objects[num] = newPile
	} else {
		d.foundationBox.Add(newPile)
	}
	d.foundationBox.Refresh()
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

func (d *Display) ReservePrint(idx int, value []string)  {}
func (d *Display) FoundationTitle(num int, value string) {}

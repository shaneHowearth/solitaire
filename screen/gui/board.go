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

	// Initial empty state
	d.talonBox.Add(d.buildPile(nil, state.StackTalon, 0, 0))
	d.wasteBox.Add(d.buildPile(nil, state.StackWaste, 0, 0))

	for i := 0; i < foundationCount; i++ {
		d.foundationBox.Add(d.buildPile(nil, state.StackFoundation, i, 0))
	}
	for i := 0; i < tableauWidth; i++ {
		d.tableauBox.Add(d.buildPile(nil, state.StackTableau, i, 0))
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
		widget.NewButton("New", func() { d.gameRedealCallback() }),
		widget.NewButton("Undo", func() { d.gameUndoCallback() }),
		widget.NewButton("Quit", func() { d.App.Quit() }),
	)

	instructionBox := widget.NewLabel(strings.Join(howTo, "\n"))
	instructionBox.Wrapping = fyne.TextWrapWord

	header := container.NewVBox(
		container.NewBorder(nil, nil, gameSelector, btnBar),
		instructionBox,
	)

	// Layout: Talon/Waste Left, Foundations Right
	topArea := container.NewBorder(nil, nil, container.NewHBox(d.talonBox, d.wasteBox), nil, d.foundationBox)

	// The tableauContainer allows horizontal scrolling if the board is wide.
	tableauContainer := container.NewHScroll(d.tableauBox)

	// Tabletop VBox with a Spacer at the end to pin everything to the top.
	d.Tabletop = container.NewVBox(
		topArea,
		container.NewPadded(tableauContainer),
		layout.NewSpacer(),
	)

	bg := canvas.NewRectangle(d.defaultBgColor)
	d.CardLayer = container.NewStack(bg, d.Tabletop)

	d.Window.SetContent(container.NewBorder(header, nil, nil, nil, d.CardLayer))
}

func (d *Display) buildPile(cards []string, sType state.StackType, idx int, showCount int) fyne.CanvasObject {
	// 1. Handle the Empty State
	if len(cards) == 0 {
		c := NewCardWidget("", sType, idx, d)
		// Wrapping in a VBox with a Spacer ensures the empty card slot
		// stays at its MinSize height and doesn't stretch down.
		return container.NewVBox(c, layout.NewSpacer())
	}

	// 2. Handle the populated state
	cardContainer := container.NewWithoutLayout()
	currentHeight := float32(cardHeight)

	for i, face := range cards {
		cardFace := "--"
		// Logic check: ensure showCount is respected
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

	// Create a transparent rectangle to define the exact hit-box/size of the stack
	spacer := canvas.NewRectangle(color.Transparent)
	spacer.SetMinSize(fyne.NewSize(float32(cardWidth), currentHeight))

	// Stack the transparent base with the card container
	pileStack := container.NewStack(spacer, cardContainer)

	// 3. Final Wrap: NewVBox with a Spacer
	// This is the critical change: it forces the stack to be only as tall
	// as 'currentHeight', pushing all remaining window space below it.
	return container.NewVBox(pileStack, layout.NewSpacer())
}

func (d *Display) TableauPrint(idx int, value []string, showCount int) {
	if idx < len(d.tableauBox.Objects) {
		d.tableauBox.Objects[idx] = d.buildPile(value, state.StackTableau, idx, showCount)
		d.tableauBox.Refresh()
	}
}

func (d *Display) FoundationPrint(num int, value []string) {
	if num < len(d.foundationBox.Objects) {
		d.foundationBox.Objects[num] = d.buildPile(value, state.StackFoundation, num, 1)
		d.foundationBox.Refresh()
	}
}

func (d *Display) TalonPrint(value []string) {
	if len(d.talonBox.Objects) > 0 {
		d.talonBox.Objects[0] = d.buildPile(value, state.StackTalon, 0, 1)
		d.talonBox.Refresh()
	}
}

func (d *Display) WastePrint(value []string) {
	if len(d.wasteBox.Objects) > 0 {
		d.wasteBox.Objects[0] = d.buildPile(value, state.StackWaste, 0, 1)
		d.wasteBox.Refresh()
	}
}

func (d *Display) ClearBoard()                           {}
func (d *Display) ReservePrint(idx int, value []string)  {}
func (d *Display) FoundationTitle(num int, value string) {}

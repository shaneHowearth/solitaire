package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

const (
	greenR       = 46
	greenG       = 125
	greenB       = 50
	greenA       = 255
	windowWidth  = 1200
	windowHeight = 600
	cardWidth    = 100
	cardHeight   = 145
)

func main() {
	myApp := app.New()

	window := myApp.NewWindow("Irate Sol")

	// Green Rectangle.
	tabletop := canvas.NewRectangle(color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA})

	// Card
	cardBase := canvas.NewRectangle(color.White)
	cardBase.StrokeColor = color.Black
	cardBase.StrokeWidth = 2
	cardBase.Resize(fyne.NewSize(cardWidth, cardHeight))

	// 3. Add a Label for the Card Value
	cardText := canvas.NewText("A ♠", color.Black)
	cardText.TextSize = 20
	cardText.Alignment = fyne.TextAlignCenter

	// 4. Group them in a Container
	// We use container.NewMax so the text sits on top of the rectangle
	card := container.NewMax(cardBase, cardText)

	// Position the card manually for now
	card.Move(fyne.NewPos(50, 50))

	// 5. The Table (Using a Stacked Container)
	// container.NewStack puts the first item at the bottom (the green felt)
	// and subsequent items on top (the cards)
	content := container.NewStack(tabletop, container.NewWithoutLayout(card))

	window.SetContent(content)
	window.Resize(fyne.NewSize(windowWidth, windowHeight))
	window.ShowAndRun()
}

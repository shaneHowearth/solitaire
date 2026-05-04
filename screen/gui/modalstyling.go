package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

func (d *Display) styledModalContent(content fyne.CanvasObject, width, height float32) fyne.CanvasObject {
	// 1. Consistent Background: Deep Grey
	bg := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})

	// 2. Consistent Border: Gold
	bg.StrokeColor = color.RGBA{R: 255, G: 215, B: 0, A: 255}
	bg.StrokeWidth = 2
	bg.CornerRadius = 4

	// 3. Wrap content in Scroll + Padding
	scroll := container.NewVScroll(container.NewPadded(content))
	scroll.SetMinSize(fyne.NewSize(width, height))

	// 4. Stack them
	return container.NewStack(bg, scroll)
}

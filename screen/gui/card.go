package gui

import (
	"image/color"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

type CardWidget struct {
	widget.BaseWidget
	Face    string
	Stack   state.StackType
	Index   int
	Display *Display
}

func NewCardWidget(face string, sType state.StackType, idx int, d *Display) *CardWidget {
	c := &CardWidget{
		Face:    face,
		Stack:   sType,
		Index:   idx,
		Display: d,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *CardWidget) Tapped(_ *fyne.PointEvent) {
	c.Display.selectComponent(c.Stack, c.Index)
}

func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	// Base layer: A dark rectangle so the card always has a "body"
	bg := canvas.NewRectangle(color.RGBA{0, 0, 0, 180})
	bg.StrokeColor = color.RGBA{255, 255, 255, 50}
	bg.StrokeWidth = 1

	filename := c.Display.getCardFilename(c.Face)

	var mainContent fyne.CanvasObject

	// Check if file exists
	if _, err := os.Stat(filename); err == nil {
		img := canvas.NewImageFromFile(filename)
		img.FillMode = canvas.ImageFillStretch // Force it to fill the card dimensions
		mainContent = img
	} else {
		// If file is missing, show a clear white text label on red
		bg.FillColor = color.RGBA{150, 0, 0, 255}
		txt := canvas.NewText(c.Face, color.White)
		txt.TextSize = 10
		txt.Alignment = fyne.TextAlignCenter
		mainContent = txt
	}

	// Combine them
	stack := container.NewStack(bg, mainContent)

	return &cardRenderer{
		stack:   stack,
		objects: []fyne.CanvasObject{stack},
	}
}

// Custom renderer to ensure Resize is handled correctly
type cardRenderer struct {
	stack   *fyne.Container
	objects []fyne.CanvasObject
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(cardWidth), float32(cardHeight))
}

func (r *cardRenderer) Refresh() {
	canvas.Refresh(r.stack)
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cardRenderer) Destroy() {}

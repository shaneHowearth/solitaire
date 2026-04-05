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
	// Base background
	bg := canvas.NewRectangle(color.RGBA{0, 0, 0, 180})

	// Selection Overlay (Initially transparent)
	overlay := canvas.NewRectangle(color.Transparent)

	filename := c.Display.getCardFilename(c.Face)
	var mainContent fyne.CanvasObject

	if _, err := os.Stat(filename); err == nil {
		img := canvas.NewImageFromFile(filename)
		img.FillMode = canvas.ImageFillStretch
		mainContent = img
	} else {
		bg.FillColor = color.RGBA{150, 0, 0, 255}
		txt := canvas.NewText(c.Face, color.White)
		txt.TextSize = 10
		txt.Alignment = fyne.TextAlignCenter
		mainContent = txt
	}

	// Layer order: Background -> Card Image -> Selection Tint
	stack := container.NewStack(bg, mainContent, overlay)

	r := &cardRenderer{
		overlay: overlay,
		stack:   stack,
		objects: []fyne.CanvasObject{stack},
		card:    c,
	}
	r.Refresh()
	return r
}

type cardRenderer struct {
	overlay *canvas.Rectangle
	stack   *fyne.Container
	objects []fyne.CanvasObject
	card    *CardWidget
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.stack.Resize(size)
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(float32(cardWidth), float32(cardHeight))
}

func (r *cardRenderer) Refresh() {
	selType, selIdx := r.card.Display.GetSelectedComponent()

	if r.card.Display.HasSelection() && r.card.Stack == selType && r.card.Index == selIdx {
		// Apply a semi-transparent blue tint to "shade" the card
		r.overlay.FillColor = color.RGBA{R: 0, G: 120, B: 215, A: 100}
	} else {
		// Return to transparent
		r.overlay.FillColor = color.Transparent
	}

	r.overlay.Refresh()
	canvas.Refresh(r.card)
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cardRenderer) Destroy() {}

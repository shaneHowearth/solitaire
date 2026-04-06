package gui

import (
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

type CardWidget struct {
	widget.BaseWidget
	Face          string
	Stack         state.StackType
	Index         int
	Display       *Display
	IsPlaceholder bool
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

func (c *CardWidget) Tapped(event *fyne.PointEvent) {
	c.Display.selectComponent(c.Stack, c.Index)
}

func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(color.RGBA{R: 20, G: 20, B: 20, A: 255}) // Slightly darker for contrast
	bg.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 40}
	bg.StrokeWidth = 1

	overlay := canvas.NewRectangle(color.Transparent)
	var mainContent fyne.CanvasObject

	if c.IsPlaceholder {
		placeholder := canvas.NewRectangle(color.RGBA{R: 100, G: 100, B: 100, A: 40})
		placeholder.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 60}
		placeholder.StrokeWidth = 2

		if c.Face != "" {
			// Determine color based on suit
			textColor := color.RGBA{R: 220, G: 220, B: 220, A: 180} // Off-white for Spades/Clubs
			if strings.Contains(c.Face, "♥") || strings.Contains(c.Face, "♦") {
				textColor = color.RGBA{R: 230, G: 50, B: 50, A: 180} // Muted Red
			}

			txt := canvas.NewText(c.Face, textColor)
			txt.TextSize = 36 // Larger for better visibility
			txt.Alignment = fyne.TextAlignCenter
			txt.TextStyle.Bold = true

			mainContent = container.NewStack(placeholder, container.NewCenter(txt))
		} else {
			mainContent = placeholder
		}
	} else {
		filename := c.Display.getCardFilename(c.Face)
		if _, err := os.Stat(filename); err == nil {
			img := canvas.NewImageFromFile(filename)
			img.FillMode = canvas.ImageFillStretch
			mainContent = img
		} else {
			txt := canvas.NewText(c.Face, color.White)
			txt.Alignment = fyne.TextAlignCenter
			mainContent = txt
		}
	}

	content := container.NewStack(bg, mainContent, overlay)

	r := &cardRenderer{
		overlay: overlay,
		stack:   content,
		objects: []fyne.CanvasObject{content},
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
	if r.card.Display.selectedComponentType == r.card.Stack && r.card.Display.selectedIndex == r.card.Index {
		r.overlay.FillColor = color.RGBA{R: 0, G: 0, B: 0, A: 90}
	} else {
		r.overlay.FillColor = color.Transparent
	}
	r.overlay.Refresh()
	canvas.Refresh(r.card)
}

func (r *cardRenderer) Destroy() {}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

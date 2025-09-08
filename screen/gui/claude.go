package gui

import (
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/state"
)

var _ screen.Display = (*DisplayGUI)(nil)

// CardWidget represents a single card that can be tapped
type CardWidget struct {
	widget.BaseWidget
	rect      *canvas.Rectangle
	text      *canvas.Text
	stackType state.StackType
	index     int
	display   *DisplayGUI
	selected  bool
}

func NewCardWidget(stackType state.StackType, index int, display *DisplayGUI) *CardWidget {
	c := &CardWidget{
		rect:      canvas.NewRectangle(color.RGBA{0, 100, 0, 255}), // Green background
		text:      canvas.NewText("", color.Black),
		stackType: stackType,
		index:     index,
		display:   display,
	}
	c.ExtendBaseWidget(c)
	return c
}

func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &cardRenderer{
		card:    c,
		objects: []fyne.CanvasObject{c.rect, c.text},
	}
}

func (c *CardWidget) Tapped(*fyne.PointEvent) {
	c.display.selectComponent(c.stackType, c.index)
}

func (c *CardWidget) SetText(text string) {
	// Handle color formatting - convert [red] tags to actual red color
	if strings.Contains(text, "[red]") {
		text = strings.ReplaceAll(text, "[red]", "")
		text = strings.ReplaceAll(text, "[-]", "")
		c.text.Color = color.RGBA{200, 0, 0, 255} // Red for hearts/diamonds
	} else {
		c.text.Color = color.Black
	}
	c.text.Text = text
	c.text.Refresh()
}

func (c *CardWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		c.rect.FillColor = color.RGBA{255, 255, 0, 255} // Yellow for selection
	} else {
		c.rect.FillColor = color.RGBA{0, 100, 0, 255} // Green default
	}
	c.rect.Refresh()
}

func (c *CardWidget) SetEmpty() {
	c.text.Text = ""
	c.rect.FillColor = color.RGBA{50, 50, 50, 255} // Dark gray for empty
	c.rect.Refresh()
	c.text.Refresh()
}

type cardRenderer struct {
	card    *CardWidget
	objects []fyne.CanvasObject
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.card.rect.Resize(size)
	r.card.text.Move(fyne.NewPos(5, 5))
	r.card.text.Resize(fyne.NewSize(size.Width-10, size.Height-10))
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 80)
}

func (r *cardRenderer) Refresh() {
	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cardRenderer) Destroy() {}

// PileWidget represents a pile of cards
type PileWidget struct {
	widget.BaseWidget
	cards     []*CardWidget
	stackType state.StackType
	index     int
	title     string
	display   *DisplayGUI
}

func NewPileWidget(title string, stackType state.StackType, index int, display *DisplayGUI) *PileWidget {
	p := &PileWidget{
		title:     title,
		stackType: stackType,
		index:     index,
		display:   display,
		cards:     make([]*CardWidget, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}

func (p *PileWidget) CreateRenderer() fyne.WidgetRenderer {
	titleText := canvas.NewText(p.title, color.Black)
	titleText.Alignment = fyne.TextAlignCenter

	return &pileRenderer{
		pile:    p,
		title:   titleText,
		objects: []fyne.CanvasObject{titleText},
	}
}

func (p *PileWidget) AddCard(text string) {
	card := NewCardWidget(p.stackType, p.index, p.display)
	card.SetText(text)
	p.cards = append(p.cards, card)
	p.Refresh()
}

func (p *PileWidget) SetCards(texts []string) {
	p.cards = make([]*CardWidget, 0, len(texts))
	for _, text := range texts {
		card := NewCardWidget(p.stackType, p.index, p.display)
		card.SetText(text)
		p.cards = append(p.cards, card)
	}
	p.Refresh()
}

func (p *PileWidget) Clear() {
	p.cards = make([]*CardWidget, 0)
	p.Refresh()
}

func (p *PileWidget) Tapped(*fyne.PointEvent) {
	p.display.selectComponent(p.stackType, p.index)
}

type pileRenderer struct {
	pile    *PileWidget
	title   *canvas.Text
	objects []fyne.CanvasObject
}

func (r *pileRenderer) Layout(size fyne.Size) {
	r.title.Move(fyne.NewPos(0, 0))
	r.title.Resize(fyne.NewSize(size.Width, 20))

	// Layout cards vertically with slight offset
	cardHeight := float32(25)
	startY := float32(25)

	for i, card := range r.pile.cards {
		card.Move(fyne.NewPos(0, startY+float32(i)*cardHeight))
		card.Resize(fyne.NewSize(size.Width, 80))
	}
}

func (r *pileRenderer) MinSize() fyne.Size {
	baseHeight := float32(105) // Title + one card
	if len(r.pile.cards) > 1 {
		baseHeight += float32(len(r.pile.cards)-1) * 25 // Overlap
	}
	return fyne.NewSize(80, baseHeight)
}

func (r *pileRenderer) Refresh() {
	// Update objects list with current cards
	r.objects = []fyne.CanvasObject{r.title}
	for _, card := range r.pile.cards {
		r.objects = append(r.objects, card)
	}

	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *pileRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *pileRenderer) Destroy() {}

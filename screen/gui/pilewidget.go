package gui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

// PileWidget - also make it a drop target
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
	p.Refresh() // Refresh to update the display
}

// You'll also need CreateRenderer for PileWidget
func (p *PileWidget) CreateRenderer() fyne.WidgetRenderer {
	titleText := canvas.NewText(p.title, color.Black)
	titleText.Alignment = fyne.TextAlignCenter

	return &pileRenderer{
		pile: p,
		// title: titleText,
	}
}

// Define the pileRenderer struct
type pileRenderer struct {
	pile        *PileWidget
	placeholder *canvas.Rectangle
	objects     []fyne.CanvasObject
}

func (r *pileRenderer) Layout(size fyne.Size) {
	// Layout cards in the pile
	for i, card := range r.pile.cards {
		cardSize := fyne.NewSize(80, 120) // Standard card size

		if r.pile.stackType == state.StackTableau {
			// Tableau cards cascade down
			card.Resize(cardSize)
			card.Move(fyne.NewPos(0, float32(i*25))) // 25 pixel offset per card
		} else {
			// Other piles stack on top of each other
			card.Resize(cardSize)
			card.Move(fyne.NewPos(0, 0))
		}
	}

	// Layout placeholder if no cards
	if len(r.pile.cards) == 0 && r.placeholder != nil {
		r.placeholder.Resize(fyne.NewSize(80, 120))
		r.placeholder.Move(fyne.NewPos(0, 0))
	}
}

func (r *pileRenderer) MinSize() fyne.Size {
	if r.pile.stackType == state.StackTableau && len(r.pile.cards) > 0 {
		// Tableau needs extra height for cascading cards
		extraHeight := float32((len(r.pile.cards) - 1) * 25)
		return fyne.NewSize(80, 120+extraHeight)
	}
	return fyne.NewSize(80, 120)
}

func (r *pileRenderer) Refresh() {
	// Refresh the title
	// r.title.Refresh()

	// Refresh all cards
	for _, card := range r.pile.cards {
		card.Refresh()
	}
}

func (r *pileRenderer) Objects() []fyne.CanvasObject {
	// Return objects dynamically
	objects := []fyne.CanvasObject{}
	for _, card := range r.pile.cards {
		objects = append(objects, card)
	}
	return objects
}

func (r *pileRenderer) Destroy() {
	// Cleanup if needed
}

func (p *PileWidget) AcceptsDrop() bool {
	return true
}

func (p *PileWidget) DropAccepted(e *fyne.DragEvent) {
	// This will be handled by the card's DragEnd method
}

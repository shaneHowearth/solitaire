package gui

import (
	"fmt"
	"image/color"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
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
	image     *canvas.Image // Add this field
	cardName  string        // Add this field
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
	objects := []fyne.CanvasObject{c.rect, c.text}
	if c.image != nil {
		objects = append(objects, c.image)
	}
	return &cardRenderer{
		card:    c,
		objects: objects,
	}
}

type cardRenderer struct {
	card    *CardWidget
	objects []fyne.CanvasObject
}

func (r *cardRenderer) Layout(size fyne.Size) {
	r.card.rect.Resize(size)

	if r.card.image != nil {
		r.card.image.Move(fyne.NewPos(0, 0))
		r.card.image.Resize(size)
	}

	r.card.text.Move(fyne.NewPos(5, 5))
	r.card.text.Resize(fyne.NewSize(size.Width-10, size.Height-10))
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(60, 80)
}

func (r *cardRenderer) Refresh() {
	// Update objects list to include image if it exists
	r.objects = []fyne.CanvasObject{r.card.rect, r.card.text}
	if r.card.image != nil {
		r.objects = append(r.objects, r.card.image)
	}

	for _, obj := range r.objects {
		obj.Refresh()
	}
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *cardRenderer) Destroy() {}

func (c *CardWidget) Tapped(*fyne.PointEvent) {
	c.display.selectComponent(c.stackType, c.index)
}

func (c *CardWidget) SetText(text string) {
	c.cardName = text

	// If we have a valid card name, try to load the image
	if text != "" && text != " " {
		c.SetCardImage(text)
	} else {
		c.SetEmpty()
	}
}

func (c *CardWidget) SetCardImage(cardName string) {
	c.cardName = cardName

	// Get the image path
	imagePath := c.getImagePath(cardName)

	// Load the image
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("ERROR: Card image file does not exist: %s\n", imagePath)
	}
	resource := storage.NewFileURI(imagePath)
	c.image = canvas.NewImageFromURI(resource)
	c.image.FillMode = canvas.ImageFillContain

	// Hide text when showing image
	c.text.Text = ""
	c.rect.FillColor = color.Transparent

	c.Refresh()
}

func (c *CardWidget) SetCardBack() {
	c.cardName = "BACK"

	// Load card back image
	resource := storage.NewFileURI("./cmd/fyne/cards/back01.png")
	c.image = canvas.NewImageFromURI(resource)
	c.image.FillMode = canvas.ImageFillContain

	// Hide text and make background transparent
	c.text.Text = ""
	c.rect.FillColor = color.Transparent
	c.Refresh()
}

func (c *CardWidget) SetEmpty() {
	c.cardName = ""
	c.text.Text = ""
	c.image = nil                                  // Remove image for empty cards
	c.rect.FillColor = color.RGBA{50, 50, 50, 255} // Dark gray for empty
	c.Refresh()
}

func (c *CardWidget) SetSelected(selected bool) {
	c.selected = selected
	if selected {
		// For image cards, add a border or overlay instead of changing background
		if c.image != nil {
			c.rect.FillColor = color.RGBA{255, 255, 0, 100} // Semi-transparent yellow
		} else {
			c.rect.FillColor = color.RGBA{255, 255, 0, 255} // Yellow for selection
		}
	} else {
		// For image cards, make background transparent
		if c.image != nil {
			c.rect.FillColor = color.Transparent
		} else {
			c.rect.FillColor = color.RGBA{0, 100, 0, 255} // Green default
		}
	}
	c.rect.Refresh()
}

func (c *CardWidget) getImagePath(cardName string) string {
	if cardName == "BACK" {
		return "./cmd/fyne/cards/back01.png"
	}

	var rank, suit string

	// Parse the card name with better error handling
	if len(cardName) < 2 {
		return "./cmd/fyne/cards/back01.png" // Fallback to card back
	}

	parts := strings.Split(cardName, " ")
	if len(parts) != 2 {
		return "./cmd/fyne/cards/back01.png" // Invalid rank, show card back
	}
	suitPart := parts[1]
	rankPart := parts[0]

	// Convert rank
	switch rankPart {
	case "Ace":
		rank = "01"
	case "Jack":
		rank = "11"
	case "Queen":
		rank = "12"
	case "King":
		rank = "13"
	case "10":
		rank = "10"
	default:
		if len(rankPart) == 1 && rankPart >= "2" && rankPart <= "9" {
			rank = fmt.Sprintf("0%s", rankPart)
		} else {
			return "./cmd/fyne/cards/back01.png" // Invalid rank, show card back
		}
	}

	// Convert suit
	switch suitPart {
	case "♠":
		suit = "s" // Spades
	case "♥":
		suit = "h" // Hearts
	case "♦":
		suit = "d" // Diamonds
	case "♣":
		suit = "c" // Clubs
	default:
		return "./cmd/fyne/cards/back01.png" // Invalid suit, show card back
	}

	// Construct file path - adjust this to match your file structure
	return fmt.Sprintf("./cmd/fyne/cards/%s%s.png", rank, suit)
}

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

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
	"github.com/shanehowearth/solitaire/state"
)

// CardWidget represents a single card that can be tapped and dragged
type CardWidget struct {
	widget.BaseWidget
	rect      *canvas.Rectangle
	text      *canvas.Text
	image     *canvas.Image
	cardName  string
	stackType state.StackType
	index     int
	display   *DisplayGUI
	selected  bool

	// Drag state
	dragging    bool
	dragStart   fyne.Position
	originalPos fyne.Position
}

func NewCardWidgetFromString(cardStr string, display *DisplayGUI) *CardWidget {
	c := &CardWidget{
		rect:     canvas.NewRectangle(color.RGBA{255, 255, 255, 255}), // White background
		text:     canvas.NewText(cardStr, color.Black),
		cardName: cardStr,
		display:  display,
		index:    0,
	}

	// Set text properties
	c.text.Alignment = fyne.TextAlignCenter
	c.text.TextStyle = fyne.TextStyle{Bold: true}

	c.ExtendBaseWidget(c)
	return c
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

func (c *CardWidget) SetText(text string) {
	c.cardName = text

	// If we have a valid card name, try to load the image
	if text != "" && text != " " {
		c.SetCardImage(text)
	} else {
		c.SetEmpty()
	}
}

func (c *CardWidget) SetEmpty() {
	c.cardName = ""
	c.text.Text = ""
	c.image = nil                                  // Remove image for empty cards
	c.rect.FillColor = color.RGBA{50, 50, 50, 255} // Dark gray for empty
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

func (c *CardWidget) SetCardImage(cardName string) {
	c.cardName = cardName

	// Get the image path
	imagePath := c.getImagePath(cardName)

	// Load the image
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		fmt.Printf("ERROR: Card image file does not exist: %s\n", imagePath)
		// Instead of failing silently, let's show text as fallback
		c.image = nil
		c.text.Text = cardName
		c.rect.FillColor = color.RGBA{255, 255, 255, 255} // White background
		c.Refresh()
		return
	}

	resource := storage.NewFileURI(imagePath)
	c.image = canvas.NewImageFromURI(resource)
	c.image.FillMode = canvas.ImageFillContain

	// Hide text when showing image
	c.text.Text = ""
	c.rect.FillColor = color.Transparent

	c.Refresh()
}

// Implement fyne.Draggable interface
func (c *CardWidget) Dragged(e *fyne.DragEvent) {
	if !c.dragging {
		return
	}

	// Calculate new position based on drag offset
	newPos := fyne.NewPos(
		c.originalPos.X+e.Position.X-c.dragStart.X,
		c.originalPos.Y+e.Position.Y-c.dragStart.Y,
	)

	// Move the card
	c.Move(newPos)

	// Optional: Show drag feedback
	if c.image != nil {
		c.rect.FillColor = color.RGBA{255, 255, 100, 100} // Semi-transparent yellow
	} else {
		c.rect.FillColor = color.RGBA{255, 255, 100, 255} // Yellow
	}
	c.rect.Refresh()
}

func (c *CardWidget) DragEnd() {
	if !c.dragging {
		return
	}

	c.dragging = false

	// Find what's under the cursor
	dropTarget := c.findDropTarget()

	if dropTarget != nil {
		// Valid drop - let the display handle the move
		c.display.handleCardDrop(c, dropTarget)
	} else {
		// Invalid drop - snap back to original position
		c.Move(c.originalPos)
	}

	// Restore normal appearance
	c.SetSelected(c.selected)
}

func (c *CardWidget) findDropTarget() fyne.CanvasObject {
	// Get current position
	currentPos := c.Position()
	cardSize := c.Size()
	centerPos := fyne.NewPos(
		currentPos.X+cardSize.Width/2,
		currentPos.Y+cardSize.Height/2,
	)

	// Find what object is at the center of the card
	return c.display.findObjectAtPosition(centerPos)
}

// Implement standard widget methods
func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &cardRenderer{card: c}
}

func (c *CardWidget) Tapped(e *fyne.PointEvent) {
	// Start drag operation
	c.dragging = true
	c.dragStart = e.Position
	c.originalPos = c.Position()

	// Also handle selection
	c.display.selectComponent(c.stackType, c.index)
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

// cardRenderer implementation
type cardRenderer struct {
	card *CardWidget
}

func (r *cardRenderer) MinSize() fyne.Size {
	return fyne.NewSize(80, 120) // Standard card size
}

func (r *cardRenderer) Objects() []fyne.CanvasObject {
	// Return objects dynamically, not from a stored slice
	objects := []fyne.CanvasObject{}

	// Always include the rectangle (for background/border)
	if r.card.rect != nil {
		objects = append(objects, r.card.rect)
	}

	// Include text if it has content
	if r.card.text != nil && r.card.text.Text != "" {
		objects = append(objects, r.card.text)
	}

	// Include image if it exists
	if r.card.image != nil {
		objects = append(objects, r.card.image)
	}

	return objects
}

func (r *cardRenderer) Destroy() {
	// Cleanup if needed
}

func (r *cardRenderer) Layout(size fyne.Size) {
	// Layout rectangle (background/border)
	if r.card.rect != nil {
		r.card.rect.Resize(size)
		r.card.rect.Move(fyne.NewPos(0, 0))
	}

	// Layout image if it exists
	if r.card.image != nil {
		r.card.image.Move(fyne.NewPos(0, 0))
		r.card.image.Resize(size)
	}

	// Layout text if it has content
	if r.card.text != nil && r.card.text.Text != "" {
		// Center the text
		textSize := r.card.text.MinSize()
		r.card.text.Move(fyne.NewPos(
			(size.Width-textSize.Width)/2,
			(size.Height-textSize.Height)/2,
		))
		r.card.text.Resize(textSize)
	}
}

func (r *cardRenderer) Refresh() {
	// Refresh the individual objects
	if r.card.rect != nil {
		r.card.rect.Refresh()
	}
	if r.card.text != nil {
		r.card.text.Refresh()
	}
	if r.card.image != nil {
		r.card.image.Refresh()
	}
}

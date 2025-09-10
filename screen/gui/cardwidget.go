package gui

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/state"
)

const (
	// Card dimensions
	CardWidth  = 80
	CardHeight = 120

	// Drag sensitivity threshold
	DragThreshold = 3.0
)

var (
	// Color constants
	WhiteColor       = color.RGBA{255, 255, 255, 255}
	BlackColor       = color.RGBA{0, 0, 0, 255}
	GreenColor       = color.RGBA{0, 100, 0, 255}
	DarkGrayColor    = color.RGBA{50, 50, 50, 255}
	YellowColor      = color.RGBA{255, 255, 0, 255}
	SemiYellowColor  = color.RGBA{255, 255, 100, 150}
	TransparentColor = color.Transparent

	// URI cache to avoid reloading the same images
	uriCache    = make(map[string]fyne.URI)
	uriCacheMux sync.RWMutex
)

// CardState represents the visual state of a card
type CardState int

const (
	CardStateNormal CardState = iota
	CardStateSelected
	CardStateDragging
	CardStateEmpty
	CardStateBack
)

// CardWidget represents a single card that can be tapped and dragged
type CardWidget struct {
	widget.BaseWidget

	// Visual components
	rect  *canvas.Rectangle
	text  *canvas.Text
	image *canvas.Image

	// Card properties
	cardName  string
	stackType state.StackType
	index     int
	state     CardState

	// References
	display *DisplayGUI

	// Drag state
	isDragging     bool
	dragStartPos   fyne.Position
	mouseStartPos  fyne.Position
	originalParent fyne.CanvasObject

	// Cached values
	minSize fyne.Size
}

// NewCardWidgetFromString creates a card widget from a card string representation
func NewCardWidgetFromString(cardStr string, display *DisplayGUI) *CardWidget {
	c := &CardWidget{
		rect:     canvas.NewRectangle(WhiteColor),
		text:     canvas.NewText(cardStr, BlackColor),
		cardName: cardStr,
		display:  display,
		index:    0,
		state:    CardStateNormal,
		minSize:  fyne.NewSize(CardWidth, CardHeight),
	}

	c.setupText()
	c.ExtendBaseWidget(c)
	return c
}

// NewCardWidget creates an empty card widget for a specific stack position
func NewCardWidget(stackType state.StackType, index int, display *DisplayGUI) *CardWidget {
	c := &CardWidget{
		rect:      canvas.NewRectangle(GreenColor),
		text:      canvas.NewText("", BlackColor),
		stackType: stackType,
		index:     index,
		display:   display,
		state:     CardStateEmpty,
		minSize:   fyne.NewSize(CardWidth, CardHeight),
	}

	c.setupText()
	c.ExtendBaseWidget(c)
	return c
}

// setupText configures the text properties
func (c *CardWidget) setupText() {
	if c.text != nil {
		c.text.Alignment = fyne.TextAlignCenter
		c.text.TextStyle = fyne.TextStyle{Bold: true}
	}
}

// SetText sets the card text and updates the visual representation
func (c *CardWidget) SetText(text string) {
	c.cardName = text

	if text != "" && text != " " {
		if err := c.setCardImage(text); err != nil {
			log.Printf("Failed to set card image for %s: %v", text, err)
			// Fallback to text display
			c.setTextDisplay(text)
		}
	} else {
		c.SetEmpty()
	}
}

// SetEmpty configures the card as empty
func (c *CardWidget) SetEmpty() {
	c.cardName = ""
	c.state = CardStateEmpty
	c.text.Text = ""
	c.image = nil
	c.updateAppearance()
}

// SetCardBack configures the card to show its back
func (c *CardWidget) SetCardBack() {
	c.cardName = "BACK"
	c.state = CardStateBack

	if err := c.loadCardBackImage(); err != nil {
		log.Printf("Failed to load card back image: %v", err)
		c.setTextDisplay("BACK")
	}
}

// loadCardBackImage loads the card back image
func (c *CardWidget) loadCardBackImage() error {
	imagePath := "./cmd/fyne/cards/back01.png"
	image, err := c.getOrCreateImage(imagePath)
	if err != nil {
		return err
	}

	c.image = image
	c.text.Text = ""
	c.updateAppearance()
	return nil
}

// setTextDisplay sets up the card for text-only display
func (c *CardWidget) setTextDisplay(text string) {
	c.text.Text = text
	c.image = nil
	c.updateAppearance()
}

// getImagePath generates the file path for a card image
func (c *CardWidget) getImagePath(cardName string) (string, error) {
	if cardName == "BACK" {
		return "./cmd/fyne/cards/back01.png", nil
	}

	if len(cardName) < 2 {
		return "", fmt.Errorf("invalid card name: too short")
	}

	parts := strings.Split(cardName, " ")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid card name format: %s", cardName)
	}

	rank, err := c.parseRank(parts[0])
	if err != nil {
		return "", fmt.Errorf("invalid rank: %w", err)
	}

	suit, err := c.parseSuit(parts[1])
	if err != nil {
		return "", fmt.Errorf("invalid suit: %w", err)
	}

	return filepath.Join("./cmd/fyne/cards", fmt.Sprintf("%s%s.png", rank, suit)), nil
}

// parseRank converts a rank string to its file representation
func (c *CardWidget) parseRank(rankStr string) (string, error) {
	switch rankStr {
	case "Ace":
		return "01", nil
	case "Jack":
		return "11", nil
	case "Queen":
		return "12", nil
	case "King":
		return "13", nil
	case "10":
		return "10", nil
	default:
		if len(rankStr) == 1 && rankStr >= "2" && rankStr <= "9" {
			return fmt.Sprintf("0%s", rankStr), nil
		}
		return "", fmt.Errorf("unknown rank: %s", rankStr)
	}
}

// parseSuit converts a suit symbol to its file representation
func (c *CardWidget) parseSuit(suitStr string) (string, error) {
	switch suitStr {
	case "♠":
		return "s", nil // Spades
	case "♥":
		return "h", nil // Hearts
	case "♦":
		return "d", nil // Diamonds
	case "♣":
		return "c", nil // Clubs
	default:
		return "", fmt.Errorf("unknown suit: %s", suitStr)
	}
}

// getOrCreateImage returns a cached image or creates a new one
func (c *CardWidget) getOrCreateImage(imagePath string) (*canvas.Image, error) {
	uriCacheMux.RLock()
	if cachedURI, exists := uriCache[imagePath]; exists {
		uriCacheMux.RUnlock()
		// Create a new image instance from the cached URI
		newImage := canvas.NewImageFromURI(cachedURI)
		newImage.FillMode = canvas.ImageFillContain
		return newImage, nil
	}
	uriCacheMux.RUnlock()

	// Check if file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("image file does not exist: %s", imagePath)
	}

	resource := storage.NewFileURI(imagePath)
	image := canvas.NewImageFromURI(resource)
	image.FillMode = canvas.ImageFillContain

	// Cache the URI
	uriCacheMux.Lock()
	uriCache[imagePath] = resource
	uriCacheMux.Unlock()

	return image, nil
}

// setCardImage sets the card image from card name
func (c *CardWidget) setCardImage(cardName string) error {
	c.cardName = cardName

	imagePath, err := c.getImagePath(cardName)
	if err != nil {
		return err
	}

	image, err := c.getOrCreateImage(imagePath)
	if err != nil {
		return err
	}

	c.image = image
	c.text.Text = ""
	c.state = CardStateNormal
	c.updateAppearance()
	return nil
}

// updateAppearance updates the visual appearance based on current state
func (c *CardWidget) updateAppearance() {
	switch c.state {
	case CardStateEmpty:
		c.rect.FillColor = DarkGrayColor
	case CardStateSelected:
		if c.image != nil {
			c.rect.FillColor = SemiYellowColor
		} else {
			c.rect.FillColor = YellowColor
		}
	case CardStateDragging:
		if c.image != nil {
			c.rect.FillColor = SemiYellowColor
		} else {
			c.rect.FillColor = YellowColor
		}
	case CardStateBack, CardStateNormal:
		if c.image != nil {
			c.rect.FillColor = TransparentColor
		} else {
			c.rect.FillColor = WhiteColor
		}
	}
	c.Refresh()
}

// MouseDown handles mouse press events
func (c *CardWidget) MouseDown(e *fyne.PointEvent) {
	c.mouseStartPos = e.Position
	c.dragStartPos = c.Position()
}

// Dragged handles drag events
func (c *CardWidget) Dragged(e *fyne.DragEvent) {
	if !c.isDragging {
		// Check if we've moved far enough to start dragging
		distance := c.calculateDistance(e.Position, c.mouseStartPos)
		if distance < DragThreshold {
			return
		}

		c.startDragging()
	}

	newPos := c.calculateDragPosition(e.Position)
	c.Move(newPos)
}

// calculateDistance calculates the distance between two points
func (c *CardWidget) calculateDistance(p1, p2 fyne.Position) float64 {
	dx := float64(p1.X - p2.X)
	dy := float64(p1.Y - p2.Y)
	return dx*dx + dy*dy // Using squared distance to avoid sqrt
}

// startDragging initializes the drag operation
func (c *CardWidget) startDragging() {
	c.isDragging = true
	c.state = CardStateDragging
	c.display.bringCardToFront(c)
	c.updateAppearance()
}

// calculateDragPosition calculates the new position during drag
func (c *CardWidget) calculateDragPosition(dragPos fyne.Position) fyne.Position {
	return fyne.NewPos(
		c.dragStartPos.X+dragPos.X-c.mouseStartPos.X,
		c.dragStartPos.Y+dragPos.Y-c.mouseStartPos.Y,
	)
}

// DragEnd handles the end of drag operations
func (c *CardWidget) DragEnd() {
	if !c.isDragging {
		return
	}

	c.isDragging = false
	c.state = CardStateNormal

	if dropTarget := c.findDropTarget(); dropTarget != nil {
		c.handleSuccessfulDrop(dropTarget)
	} else {
		c.snapBackToOriginalPosition()
	}

	c.updateAppearance()
}

// findDropTarget finds a valid drop target at the current position
func (c *CardWidget) findDropTarget() *PileWidget {
	cardPos := c.Position()
	cardSize := c.Size()

	for _, pile := range c.getAllPiles() {
		if pile == nil || c.isSamePile(pile) {
			continue
		}

		pilePos := pile.Position()
		pileSize := pile.Size()

		// Check if rectangles overlap
		if c.rectanglesOverlap(cardPos, cardSize, pilePos, pileSize) {
			return pile
		}
	}

	return nil
}

// isSamePile checks if the pile is the same as the card's origin
func (c *CardWidget) isSamePile(pile *PileWidget) bool {
	return pile.stackType == c.stackType && pile.index == c.index
}

// rectanglesOverlap checks if two rectangles overlap
func (c *CardWidget) rectanglesOverlap(pos1 fyne.Position, size1 fyne.Size, pos2 fyne.Position, size2 fyne.Size) bool {
	return pos1.X < pos2.X+size2.Width && pos1.X+size1.Width > pos2.X &&
		pos1.Y < pos2.Y+size2.Height && pos1.Y+size1.Height > pos2.Y
}

// getAllPiles returns all available pile widgets
func (c *CardWidget) getAllPiles() []*PileWidget {
	var allPiles []*PileWidget

	if c.display.stock != nil {
		allPiles = append(allPiles, c.display.stock)
	}
	if c.display.waste != nil {
		allPiles = append(allPiles, c.display.waste)
	}

	allPiles = append(allPiles, c.display.foundations...)
	allPiles = append(allPiles, c.display.tableau...)
	allPiles = append(allPiles, c.display.reserves...)

	return allPiles
}

// handleSuccessfulDrop processes a successful drop operation
func (c *CardWidget) handleSuccessfulDrop(targetPile *PileWidget) {
	if c.display.componentSelectedCallback != nil {
		c.display.componentSelectedCallback(
			c.stackType, c.index,
			targetPile.stackType, targetPile.index,
		)
	}
}

// snapBackToOriginalPosition moves the card back to its starting position
func (c *CardWidget) snapBackToOriginalPosition() {
	c.Move(c.dragStartPos)
}

// Tapped handles tap events
func (c *CardWidget) Tapped(e *fyne.PointEvent) {
	if !c.isDragging {
		c.display.selectComponent(c.stackType, c.index)
	}
}

// SetSelected updates the selection state
func (c *CardWidget) SetSelected(selected bool) {
	if selected {
		c.state = CardStateSelected
	} else {
		c.state = CardStateNormal
	}
	c.updateAppearance()
}

// CreateRenderer creates the widget renderer
func (c *CardWidget) CreateRenderer() fyne.WidgetRenderer {
	return &cardRenderer{card: c}
}

// MinSize returns the minimum size of the widget
func (c *CardWidget) MinSize() fyne.Size {
	return c.minSize
}

// cardRenderer implements the widget renderer
type cardRenderer struct {
	card *CardWidget
}

// MinSize returns the minimum size for the renderer
func (r *cardRenderer) MinSize() fyne.Size {
	return r.card.minSize
}

// Objects returns the objects to be rendered
func (r *cardRenderer) Objects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, 0, 3)

	if r.card.rect != nil {
		objects = append(objects, r.card.rect)
	}

	if r.card.text != nil && r.card.text.Text != "" {
		objects = append(objects, r.card.text)
	}

	if r.card.image != nil {
		objects = append(objects, r.card.image)
	}

	return objects
}

// Layout arranges the objects within the widget
func (r *cardRenderer) Layout(size fyne.Size) {
	if r.card.rect != nil {
		r.card.rect.Resize(size)
		r.card.rect.Move(fyne.NewPos(0, 0))
	}

	if r.card.image != nil {
		r.card.image.Move(fyne.NewPos(0, 0))
		r.card.image.Resize(size)
	}

	if r.card.text != nil && r.card.text.Text != "" {
		textSize := r.card.text.MinSize()
		r.card.text.Move(fyne.NewPos(
			(size.Width-textSize.Width)/2,
			(size.Height-textSize.Height)/2,
		))
		r.card.text.Resize(textSize)
	}
}

// Refresh refreshes the visual components
func (r *cardRenderer) Refresh() {
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

// Destroy cleans up resources
func (r *cardRenderer) Destroy() {
	// Cleanup if needed - URIs are cached globally
}

// ClearURICache clears the global URI cache (useful for cleanup or memory management)
func ClearURICache() {
	uriCacheMux.Lock()
	defer uriCacheMux.Unlock()

	for key := range uriCache {
		delete(uriCache, key)
	}
}

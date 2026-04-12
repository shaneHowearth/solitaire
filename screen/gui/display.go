package gui

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

const (
	greenR       = 46
	greenG       = 125
	greenB       = 50
	greenA       = 255
	windowWidth  = 1200
	windowHeight = 900
	cardWidth    = 100
	cardHeight   = 145
	verticalFan  = 30
)

type Display struct {
	App       fyne.App
	Window    fyne.Window
	Tabletop  *fyne.Container
	CardLayer *fyne.Container

	// Containers for specific game areas
	talonBox      *fyne.Container
	wasteBox      *fyne.Container
	foundationBox *fyne.Container
	tableauBox    *fyne.Container
	reserveBox    *fyne.Container // Added for Reserve support

	Selected game.Variant
	games    []game.Variant
	screens  map[string]fyne.CanvasObject

	// Store dimensions for index mapping
	tableauWidth  int
	tableauHeight int

	gameHint                  func() []state.Move
	gameSelectedCallback      func(game.Variant)
	gameRedealCallback        func()
	gameUndoCallback          func()
	componentSelectedCallback func(state.StackType, int, state.StackType, int)

	selectedComponentType state.StackType
	selectedIndex         int

	defaultBgColor  color.Color
	selectedBgColor color.Color
	processingClick bool

	foundationHints map[int]string
}

func New(variants []game.Variant) *Display {
	myApp := app.New()
	window := myApp.NewWindow("Irate Sol")

	greenFelt := color.RGBA{R: greenR, G: greenG, B: greenB, A: greenA}

	d := &Display{
		App:             myApp,
		Window:          window,
		games:           variants,
		screens:         make(map[string]fyne.CanvasObject),
		defaultBgColor:  greenFelt,
		selectedBgColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		selectedIndex:   -1,
	}

	if len(variants) > 0 {
		d.Selected = variants[0]
	}

	bg := canvas.NewRectangle(d.defaultBgColor)
	label := widget.NewLabel("Initializing...")
	initContent := container.NewStack(bg, container.NewCenter(label))

	window.SetContent(initContent)
	window.Resize(fyne.NewSize(windowWidth, windowHeight))

	return d
}

func (d *Display) Run() error {
	d.Window.ShowAndRun()
	return nil
}

func (d *Display) Show(name string) {
	if name == "Games" && d.gameSelectedCallback != nil && len(d.games) > 0 {
		d.gameSelectedCallback(d.games[0])
		return
	}
	if screen, exists := d.screens[name]; exists {
		d.Window.SetContent(screen)
		d.Window.Content().Refresh()
	}
}

func (d *Display) selectComponent(sType state.StackType, index int) {
	if d.processingClick {
		return
	}

	d.processingClick = true
	defer func() { d.processingClick = false }()

	defer d.RefreshAll()

	if d.Selected == nil {
		return
	}

	if sType == state.StackTalon {
		// If Talon() is false (Agnes), we trigger Redeal.
		if !d.Selected.Talon() {
			if d.gameRedealCallback != nil {
				d.gameRedealCallback()
			}
			d.ClearSelection()
			return // EXIT HERE so it doesn't try to "select" the talon
		}
	}

	if d.selectedIndex != -1 {
		if d.componentSelectedCallback != nil {
			d.componentSelectedCallback(d.selectedComponentType, d.selectedIndex, sType, index)
		}
		d.selectedIndex = -1
		d.selectedComponentType = -1
	} else {
		d.selectedIndex = index
		d.selectedComponentType = sType
	}
}

func (d *Display) RefreshAll() {
	if d.talonBox != nil {
		d.talonBox.Refresh()
	}
	if d.wasteBox != nil {
		d.wasteBox.Refresh()
	}
	if d.foundationBox != nil {
		d.foundationBox.Refresh()
	}
	if d.tableauBox != nil {
		d.tableauBox.Refresh()
	}
	if d.reserveBox != nil {
		d.reserveBox.Refresh()
	}
}

func (d *Display) getCardFilename(card string) string {
	if card == "" || card == "--" || card == "🂠" {
		return "assets/cards/PySolFC/back131.gif"
	}
	card = strings.TrimSpace(card)

	rank := "00"
	if strings.Contains(card, "Ace") {
		rank = "01"
	} else if strings.Contains(card, "Jack") {
		rank = "11"
	} else if strings.Contains(card, "Queen") {
		rank = "12"
	} else if strings.Contains(card, "King") {
		rank = "13"
	} else {
		numStr := ""
		for _, r := range card {
			if r >= '0' && r <= '9' {
				numStr += string(r)
			} else if numStr != "" {
				break
			}
		}
		if len(numStr) == 1 {
			rank = "0" + numStr
		} else if len(numStr) == 2 {
			rank = numStr
		}
	}

	suit := "z"
	if strings.Contains(card, "♠") {
		suit = "s"
	} else if strings.Contains(card, "♥") {
		suit = "h"
	} else if strings.Contains(card, "♦") {
		suit = "d"
	} else if strings.Contains(card, "♣") {
		suit = "c"
	}

	return fmt.Sprintf("assets/cards/PySolFC/%s%s.gif", rank, suit)
}

// Callback Setters
func (d *Display) SetComponentSelectedCallback(cb func(state.StackType, int, state.StackType, int)) {
	d.componentSelectedCallback = cb
}
func (d *Display) SetGameSelectedCallback(cb func(game.Variant)) { d.gameSelectedCallback = cb }
func (d *Display) SetGameRedealCallback(cb func())               { d.gameRedealCallback = cb }
func (d *Display) SetGameUndoCallback(cb func())                 { d.gameUndoCallback = cb }
func (d *Display) SetHintsCallback(cb func() []state.Move)       { d.gameHint = cb }

// State Getters
func (d *Display) HasSelection() bool { return d.selectedIndex >= 0 }
func (d *Display) ClearSelection()    { d.selectedIndex = -1; d.selectedComponentType = -1 }
func (d *Display) GetSelectedComponent() (state.StackType, int) {
	return d.selectedComponentType, d.selectedIndex
}

func (d *Display) ShowWinnerModal(gameName string, score int) {
	d.LaunchFireworks()
}

type particle struct {
	shape *canvas.Circle
	velX  float32
	velY  float32
	clr   color.RGBA
}

func (d *Display) LaunchFireworks() {
	// Particle tracker to manage movement and color per spark
	type particle struct {
		shape *canvas.Circle
		velX  float32
		velY  float32
		clr   color.RGBA
	}

	// Create a dedicated overlay layer for fireworks so they don't interfere with card widgets
	fireworkLayer := container.NewWithoutLayout()
	d.CardLayer.Add(fireworkLayer)

	// Launch multiple bursts at staggered intervals for a better show
	for i := 0; i < 6; i++ {
		go func(burstIdx int) {
			// Delay each burst slightly
			time.Sleep(time.Duration(burstIdx) * 500 * time.Millisecond)

			// Randomize the origin point within the window bounds (with padding)
			size := d.Window.Canvas().Size()
			origin := fyne.NewPos(
				150+rand.Float32()*(size.Width-300),
				150+rand.Float32()*(size.Height-300),
			)

			count := 45 // Number of sparks per burst
			burst := make([]particle, count)

			for j := 0; j < count; j++ {
				// Generate a bright, high-saturation random color
				sparkColor := color.RGBA{
					R: uint8(rand.IntN(155) + 100),
					G: uint8(rand.IntN(155) + 100),
					B: uint8(rand.IntN(155) + 100),
					A: 255,
				}

				p := canvas.NewCircle(sparkColor)
				p.Resize(fyne.NewSize(4, 4))
				p.Move(origin)

				// Physics: Random angle and initial speed
				angle := rand.Float64() * 2 * math.Pi
				speed := 2.0 + rand.Float64()*5.0

				burst[j] = particle{
					shape: p,
					velX:  float32(math.Cos(angle) * speed),
					velY:  float32(math.Sin(angle) * speed),
					clr:   sparkColor,
				}
				fireworkLayer.Add(p)
			}

			// Define the 2-second animation sequence
			anim := fyne.NewAnimation(time.Second*2, func(v float32) {
				for j := range burst {
					p := &burst[j]

					// 1. Update Position based on velocity
					p.shape.Move(p.shape.Position().Add(fyne.NewPos(p.velX, p.velY)))

					// 2. Apply Gravity (pulls velY downward over time)
					p.velY += 0.12

					// 3. Apply Air Friction (optional, slows them down slightly)
					p.velX *= 0.99
					p.velY *= 0.99

					// 4. Fade Alpha based on animation progress (v goes from 0.0 to 1.0)
					newAlpha := uint8(255 * (1.0 - v))
					fadeColor := p.clr
					fadeColor.A = newAlpha

					p.shape.FillColor = fadeColor
					p.shape.Refresh()
				}

				// 5. Cleanup: Remove objects from the canvas tree once animation completes
				if v == 1.0 {
					for _, p := range burst {
						fireworkLayer.Remove(p.shape)
					}
					// Optional: if all bursts are done, you could remove fireworkLayer itself
				}
			})

			anim.Start()
		}(i)
	}
}

func (d *Display) showHintModal() {
	if d.gameHint == nil {
		return
	}

	moves := d.gameHint()
	if len(moves) == 0 {
		dialog.ShowInformation("Hints", "No moves available! Try drawing from the Stock.", d.Window)
		return
	}

	// Create a container to hold our hint rows
	hintList := container.NewVBox()

	for i, m := range moves {
		// Create a descriptive string for the move
		// e.g., "1) A♠ from Waste → Foundation 1"
		src := d.formatLocation(m.Source)
		dst := d.formatLocation(m.Destination)

		cardStr := m.SourceCardTop.String()
		if m.NumberMoving > 1 {
			cardStr = fmt.Sprintf("%s (+%d cards)", cardStr, m.NumberMoving-1)
		}

		hintLabel := widget.NewLabel(fmt.Sprintf("%d) %s: %s → %s", i+1, cardStr, src, dst))
		hintList.Add(hintLabel)
	}

	styledContent := d.styledModalContent(hintList, 500, 350)

	dialog.ShowCustom("Available Hints", "Close", styledContent, d.Window)
}

func (d *Display) formatLocation(stack state.Stack) string {
	switch stack.Type {
	case state.StackTableau:
		// Map flat index back to Row/Col for the user
		row := (stack.TableauPosition / d.tableauWidth) + 1
		col := (stack.TableauPosition % d.tableauWidth) + 1
		rowStr := ""
		if d.tableauHeight > 1 {
			rowStr = fmt.Sprintf("Row %d, Col ", row)
		}
		return fmt.Sprintf("Tableau %s%d", rowStr, col)
	case state.StackFoundation:
		return fmt.Sprintf("Foundation %d", stack.FoundationPosition+1)
	case state.StackReserve:
		return fmt.Sprintf("Reserve")
	case state.StackWaste:
		return "Waste"
	default:
		return "Stock"
	}

	return ""
}

func (d *Display) showHowToModal(gameName string, howTo []string) {
	fullText := strings.Join(howTo, "\n\n")
	// Create a text label with the instructions
	content := widget.NewLabel(fullText)
	// content := widget.NewRichTextFromMarkdown(strings.Join(howTo, "\n\n"))
	content.Wrapping = fyne.TextWrapWord

	styledContent := d.styledModalContent(content, 600, 450)

	// Show it as a modal dialog
	title := fmt.Sprintf("How to Play: %s", gameName)
	dialog.ShowCustom(title, "Back to Game", styledContent, d.Window)
}

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

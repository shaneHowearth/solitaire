package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/state"
)

// CreateBoard - Create the board that the game will use.
func (display *Display) CreateBoard(
	name string,
	tableauHeight, tableauWidth, foundationCount int,
	foundationBase state.Rank,
) {
	gamePage := display.createGamePage(name, tableauHeight, tableauWidth, foundationCount, foundationBase)

	display.screens[name] = gamePage
}

func (display *Display) createGamePage(
	name string,
	tableauHeight, tableauWidth, foundationCount int,
	foundationBase state.Rank,
) tview.Primitive {
	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)

	title := tview.NewTextView().
		SetText(fmt.Sprintf("Playing: %s", name)).
		SetTextAlign(tview.AlignCenter).
		SetTextColor(tview.Styles.PrimaryTextColor)
	title.SetBorder(true)

	// Add a box for each foundation.
	///////////////////
	/// FOUNDATIONS ///
	///////////////////
	foundationsRow := tview.NewFlex().SetDirection(tview.FlexColumn)
	display.foundations = make([]*tview.TextView, foundationCount)
	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView()

		// Add some decorations to the box.
		foundation.Box.SetBorder(true).SetTitle(foundationBase.String())
		display.foundations[idx] = foundation

		foundationsRow.AddItem(
			foundation, 0, 1, false,
		)
	}

	// The tableau.
	////////////////
	/// TABLEAUS ///
	////////////////
	tableauArea := tview.NewFlex().SetDirection(tview.FlexColumn)
	display.tableau = make([]*tview.TextView, tableauHeight*tableauWidth)

	for idx := 0; idx < tableauHeight*tableauWidth; idx++ {
		tableau := tview.NewTextView()

		tableau.SetBorder(true).SetTitle(fmt.Sprintf(""))
		display.tableau[idx] = tableau

		// Add the row to the tableau.
		tableauArea.AddItem(tableau, 0, 1, false)
	}

	// Controls/Help
	controls := tview.NewTextView().
		SetText("Press 'q' to return to game selection, Ctrl+C to quit").
		SetTextAlign(tview.AlignCenter)
	controls.SetBorder(true).SetTitle("Controls")

	display.stack = make([]*tview.TextView, 0, 1)
	display.waste = make([]*tview.TextView, 0, 1)

	// There are two rows, the top one with the Talon and Foundations, and
	// the bottom one with the Tableaus.
	// The top row.
	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	display.App.EnableMouse(true)

	// Add a box for the talon.
	/////////////
	/// TALON ///
	/////////////
	talon := tview.NewTextView().SetDynamicColors(true)
	display.stack = append(display.stack, talon)
	talon.SetWordWrap(true).SetBorder(true).SetTitle("Talon").
		SetFocusFunc(func() {
			switch talon.GetBackgroundColor() {
			case tcell.ColorRed:
				talon.SetBackgroundColor(tcell.ColorDefault)
			default:
				talon.SetBackgroundColor(tcell.ColorRed)
			}
		},
		)

	topRow.AddItem(
		talon, 0, 1, true,
	)

	// Add a box for the waste.
	/////////////
	/// WASTE ///
	/////////////
	waste := tview.NewTextView()
	display.waste = append(display.waste, waste)
	waste.SetBorder(true).SetTitle("Waste")

	// Add the top row to the main rows container.
	mainRows.
		AddItem(title, 0, 1, false).
		AddItem(foundationsRow, 8, 0, false).
		AddItem(tableauArea, 0, 1, true).
		AddItem(controls, 3, 0, false)

	// Add the main rows to the window container.
	mainRows.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 'q':
			display.Show("Games")
			return nil
		}
		return event
	})

	display.screens[name] = mainRows

	return mainRows
}

// FoundationTitle -
func (display *Display) FoundationTitle(num int, value string) {
	display.foundations[num].SetTitle(value)
	go func() {
		time.Sleep(1 * time.Millisecond)
		display.App.ForceDraw()
	}()
}

// FoundationPrint -
func (display *Display) FoundationPrint(num int, value []string) {
	if len(value) > 0 {
		display.foundations[num].SetText(
			value[len(value)-1],
		)
	}
}

// TableauPrint -
func (display *Display) TableauPrint(idx int, value []string) {
	if len(value) > 0 {
		display.tableau[idx].SetText(
			strings.Join(value, "\n"),
		)
	}
}

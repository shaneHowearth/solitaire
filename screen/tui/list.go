package tui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// Display -
type Display struct {
	app         *tview.Application
	root        *tview.Flex
	stack       []*tview.TextView
	waste       []*tview.TextView
	foundations []*tview.TextView
	tableau     [][]*tview.TextView // Row, Column
	Selected    game.Variant
	games       []game.Variant
	screens     map[string]tview.Primitive
}

func New(games []game.Variant) *Display {
	app := tview.NewApplication()

	display := &Display{
		app:     app,
		games:   games,
		screens: make(map[string]tview.Primitive),
	}
	display.screens["Games"] = display.CreateGameListPage(games)
	display.app.SetRoot(display.screens["Games"], true)

	return display
}

// Show - show the named screen.
func (display *Display) Show(name string) {
	display.app.SetRoot(display.screens[name], true).SetFocus(display.screens[name]).Run()
}

func (display *Display) Switch(page string) {
	fmt.Println("Called")
}

// CreateGameListPage - Lists all variants of games that the application knows about,
// allowing the user to select which variant to play.
func (display *Display) CreateGameListPage(games []game.Variant) *tview.List {
	list := tview.NewList()

	for idx, game := range games {
		list.AddItem(game.Name(),
			"",
			[]rune(fmt.Sprintf("%d", idx+1))[0],
			func() {
				display.Selected = game
				display.app.Stop()
			},
		)
	}

	// Add a quit option.
	list.AddItem("Quit", "", 'q', func() { display.app.Stop() })

	return list
}

// GetSelected -
func (display *Display) GetSelected() game.Variant {
	return display.Selected
}

// Board - Create the board that the game will use.
func (display *Display) CreateBoard(name string, tableauHeight, tableauWidth, foundationCount int, foundationBase state.Rank) {

	display.foundations = make([]*tview.TextView, 0, foundationCount)
	display.tableau = make([][]*tview.TextView, tableauHeight)

	for idx := range display.tableau {
		display.tableau[idx] = make([]*tview.TextView, 0, tableauWidth)
	}

	// There are two rows, the top one with the Talon and Foundations, and
	// the bottom one with the Tableaus.
	stackFunc := func(
		screen tcell.Screen,
		x, y, width, height int,
	) (
		int, int, int, int,
	) {
		tview.Print(screen, "I am the stack", x, height/2-1, width, tview.AlignCenter, tcell.ColorDefault)
		return x, y, width, height
	}

	mainRows := tview.NewFlex().SetDirection(tview.FlexRow)
	// The top row.

	topRow := tview.NewFlex().SetDirection(tview.FlexColumn)

	// Add a box for the stack.
	stack := tview.NewTextView()
	display.stack = append(display.stack, stack)

	topRow.AddItem(
		stack.SetWordWrap(true).SetDrawFunc(stackFunc).SetBorder(true).SetTitle("Talon"), 0, 1, false,
	)

	// Add a box for the waste.
	waste := tview.NewTextView()
	display.waste = append(display.waste, waste)

	topRow.AddItem(
		waste.SetBorder(true).SetTitle("Waste"), 0, 1, false,
	)

	// Add a box for each foundation.
	for idx := 0; idx < foundationCount; idx++ {
		foundation := tview.NewTextView()

		display.foundations = append(display.foundations, foundation)

		topRow.AddItem(
			foundation.Box.SetBorder(true).SetTitle(foundationBase.String()), 0, 1, false,
		)
	}

	// Add the top row to the main rows container.
	mainRows.AddItem(topRow, 0, 1, false)

	// The tableau.
	tableau := tview.NewFlex().SetDirection(tview.FlexRow)

	for idx := 0; idx < tableauHeight; idx++ {
		// Create a new row for the tableau.
		tableauRow := tview.NewFlex().SetDirection(tview.FlexColumn)

		// Add columns to the row.
		for colIdx := 0; colIdx < tableauWidth; colIdx++ {
			tableauCell := tview.NewTextView()

			display.tableau[idx] = append(display.tableau[idx], tableauCell)

			tableauRow.AddItem(
				tableauCell.SetBorder(true).SetTitle("Box"), 0, 1, false,
			)
		}

		// Add the row to the tableau.
		tableau.AddItem(tableauRow, 0, 1, false)
	}

	// Add the tableau to the main rows container.
	mainRows.AddItem(tableau, 0, tableauHeight, false)

	mainWindow := tview.NewFlex()

	// Add the main rows to the window container.
	mainWindow.AddItem(mainRows, 0, 2, true)

	display.screens[name] = mainWindow

}

// FirstDeal -
func (*Display) FirstDeal(variant game.Variant) {
	// TODO Marry the position of the stacks with cells in the display.
	variant.SetupDeal()
}

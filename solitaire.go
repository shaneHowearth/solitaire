package solitaire

import (
	"fmt"

	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/screen/tui"
	"github.com/shanehowearth/solitaire/state"
)

// This is the controller, it acts as the intermediary between the view and
// the model.
// Actions of the user are captured by the view, passed to the controller which
// then instructs the model on what to do. The change to the model is then
// relayed to the view to be displayed to the user.

// Instance - holder of the information required for an instance of the game
// TODO - give this baby a proper name :)
type Instance struct {
	Game game.Variant
	// View.
	Display screen.Display
	// Model.
	Foundations []state.Foundation
	Tableau     []*state.Tableau
	Talon       *state.Talon
	Deck        *state.Deck

	// Track first selection for move operations
	firstSelection     bool
	firstComponentType screen.ComponentType
	firstIndex         int
}

// New - create a new instance.
func New() *Instance {
	return &Instance{
		firstSelection: false,
		firstIndex:     -1,
	}
}

// Start - start the game.
func (instance *Instance) Start() error {
	// Available games.
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})
	variants = append(variants, &game.Klondike2{})

	instance.Display = tui.New(variants)

	instance.Display.SetGameSelectedCallback(instance.onGameSelected)
	// Set up component selection callback
	if tuiDisplay, ok := instance.Display.(*tui.Display); ok {
		tuiDisplay.SetComponentSelectedCallback(instance.onComponentSelected)
	}

	// Show the list of games available to play.
	instance.Display.Show("Games")

	return instance.Display.Run()
}

func (instance *Instance) onGameSelected(selectedGame game.Variant) {
	instance.Game = selectedGame

	// Create and set up the game state
	instance.setupGameState()

	// Create the game board page dynamically
	instance.createGamePage()

	// Switch to the game page
	instance.Display.Show(instance.Game.Name())
}

// onComponentSelected - handle component selection events
func (instance *Instance) onComponentSelected(componentType screen.ComponentType, index int) {
	if !instance.firstSelection {
		// First selection - just mark the source
		instance.firstSelection = true
		instance.firstComponentType = componentType
		instance.firstIndex = index
		return
	}

	// Second selection - attempt to move from first to second
	defer func() {
		// Always clear selection state after attempting a move
		instance.firstSelection = false
		instance.firstIndex = -1
		instance.Display.ClearSelection()
	}()
}

func (instance *Instance) setupGameState() {
	// Get the foundation information for the game.
	numFoundations, foundationBase, foundationRule := instance.Game.Foundations()
	// Get the tableau information for the game.
	numTableau, _, tableauRule := instance.Game.Tableau()

	// Create the state/model for the game.
	gameState := state.New(
		instance.Game.Decks(),
		numFoundations,
		foundationBase,
		foundationRule,
		numTableau,
		tableauRule,
		1,
		1,
		// Talon rule is to allow everything to be added to its stacks.
		func(state.SuitedCard) bool { return true },
	)

	// Copy the game state instantiated model into the current instance.
	instance.Foundations = gameState.Foundations
	instance.Tableau = gameState.Tableau
	instance.Talon = gameState.Talon
	instance.Deck = gameState.Deck

	instance.dealCards()
}

// dealCards - deal cards to tableau
func (instance *Instance) dealCards() {
	// Shuffle the cards.
	instance.Deck.Shuffle()

	numTableau, _, _ := instance.Game.Tableau()
	counts := instance.Game.SetupDealCardCounts()

	// Deal the cards out onto the different stacks (talon, tableau).
	for idx := 0; idx < numTableau; idx++ {
		// Grab a copy of the existing rule on the stack and replace it with
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being
		// dealt most definitely do not adhere to it (the rule).
		rule := instance.Tableau[idx].Stack.Rule
		instance.Tableau[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
		countIdx := idx * 2
		numCards := counts[countIdx]
		numOpen := counts[countIdx+1]

		for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
			card := instance.Deck.Deal()
			instance.Tableau[idx].Add(card, false)
		}

		for openIdx := 0; openIdx < numOpen; openIdx++ {
			card := instance.Deck.Deal()
			instance.Tableau[idx].Add(card, true)
		}

		// Return the rule to its correct state.
		instance.Tableau[idx].Stack.Rule = rule
	}

}

// createGamePage - create the game board page dynamically
func (instance *Instance) createGamePage() {
	// Create the board that will be displayed.
	// instance.CreateBoard(instance.Game)
	tableauHeight, tableauWidth := instance.Game.TableauGridSize()
	foundationCount, foundationBase, _ := instance.Game.Foundations()

	// Create the board layout
	instance.Display.CreateBoard(
		instance.Game.Name(),
		tableauHeight,
		tableauWidth,
		foundationCount,
		foundationBase,
	)

	// Update the display with current game state
	instance.updateDisplay()
}

// updateDisplay - update the display with current game state
func (instance *Instance) updateDisplay() {

	// Tell the board what to display in each box.
	for idx := range instance.Foundations {
		instance.Display.FoundationTitle(idx,
			fmt.Sprintf("%s %s",
				instance.Foundations[idx].Base.Rank.String(),
				instance.Foundations[idx].Base.Suit.String(),
			),
		)
		instance.Display.FoundationPrint(idx,
			instance.Foundations[idx].Stack.Cards(),
		)
	}

	// Display each tableau.
	for idx := range instance.Tableau {
		instance.Display.TableauPrint(idx,
			instance.Tableau[idx].Stack.Cards(),
		)
	}
}

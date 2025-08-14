package solitaire

import (
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
	Foundations []*state.Foundation
	Tableau     []*state.Tableau
	Talon       *state.Talon
	Deck        *state.Deck

	// Track first selection for move operations
	firstSelection     bool
	firstComponentType state.StackType
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
// Moving cards from one stack to another.
func (instance *Instance) onComponentSelected(
	fromComponentType state.StackType, fromIndex int,
	toComponentType state.StackType, toIndex int,
) {

	if fromComponentType == toComponentType && fromIndex == toIndex {
		// Nothing to do.
		return
	}

	var fromStack *state.Stack
	switch fromComponentType {
	case state.StackFoundation:
		fromStack = instance.Foundations[fromIndex].Stack
	case state.StackTableau:
		fromStack = instance.Tableau[fromIndex].Stack
	case state.StackTalon:
		fromStack = instance.Talon.Stock
	case state.StackWaste:
		fromStack = instance.Talon.Waste
	}

	var toStack *state.Stack
	switch toComponentType {
	case state.StackFoundation:
		toStack = instance.Foundations[toIndex].Stack
	case state.StackTableau:
		toStack = instance.Tableau[toIndex].Stack
	case state.StackWaste:
		toStack = instance.Talon.Waste
	}

	fromStack.Move(toStack)

	instance.updateDisplay()

	if instance.Game.HasWon(instance.Tableau, instance.Foundations) {
		// TODO: Add a score to display.
		instance.Display.ShowWinnerModal(instance.Game.Name(), 100)
	}

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

	// Put one card onto the Waste.
	card := instance.Deck.Deal()
	instance.Talon.Waste.Add(card, true)

	// Put the rest of the cards onto the talon.
	for {
		if instance.Deck.Len() == 0 {
			break
		}
		card := instance.Deck.Deal()
		instance.Talon.Stock.Add(card, false)
	}

}

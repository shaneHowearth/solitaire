package solitaire

import (
	"fmt"

	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/state"
)

// This is the controller, it acts as the intermediary between the view and.
// the model.
// Actions of the user are captured by the view, passed to the controller which.
// then instructs the model on what to do. The change to the model is then.
// relayed to the view to be displayed to the user.

// Instance - holder of the information required for an instance of the game.
// TODO - give this baby a proper name :).
type Instance struct {
	Game game.Variant
	// View.
	Display screen.Display
	// State of the current game.
	State   state.State
	History state.History

	// Track first selection for move operations.
	firstSelection bool
	firstIndex     int
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
	// Call the onGameSelected function when a game is selected.
	instance.Display.SetGameSelectedCallback(instance.onGameSelected)

	// Set up component selection callback.
	instance.Display.SetComponentSelectedCallback(instance.onComponentSelected)

	// Set up Hint function callback.
	instance.Display.SetHintsCallback(instance.onHints)

	// Show the list of games available to play.
	instance.Display.Show("Games")

	return instance.Display.Run()
}

func (instance *Instance) onGameSelected(selectedGame game.Variant) {
	instance.Game = selectedGame

	// Create and set up the game state.
	instance.setupGameState()

	// Create the game board page dynamically.
	instance.createGamePage()

	// Set the redeal callback for this game.
	instance.Display.SetGameRedealCallback(instance.redeal)

	// Set the undo callback for this game.
	instance.Display.SetGameUndoCallback(instance.undo)

	// Switch to the game page.
	instance.Display.Show(instance.Game.Name())
}

// onComponentSelected - handle component selection events.
// Moving cards from one stack to another.
func (instance *Instance) onComponentSelected(
	fromComponentType state.StackType, fromIndex int,
	toComponentType state.StackType, toIndex int,
) {
	var fromStack *state.Stack

	switch fromComponentType {
	case state.StackFoundation:
		if fromIndex < 0 || fromIndex >= len(instance.State.Foundations) {
			return // Out of bounds filler foundation, drop event smoothly
		}
		fromStack = instance.State.Foundations[fromIndex].Stack
	case state.StackTableau:
		if fromIndex < 0 || fromIndex >= len(instance.State.Tableau) {
			return // Out of bounds filler grid slot, drop event smoothly
		}
		fromStack = instance.State.Tableau[fromIndex].Stack
	case state.StackTalon:
		fromStack = instance.State.Talon.Stock
		if !instance.Game.Talon() {
			return
		}
	case state.StackWaste:
		fromStack = instance.State.Talon.Waste
	case state.StackReserve:
		if fromIndex < 0 || fromIndex >= len(instance.State.Reserves) {
			return // Out of bounds filler reserve slot, drop event smoothly
		}
		fromStack = instance.State.Reserves[fromIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'fromComponentType' %d", fromComponentType))
	}

	var toStack *state.Stack

	switch toComponentType {
	case state.StackFoundation:
		if toIndex < 0 || toIndex >= len(instance.State.Foundations) {
			return // Out of bounds filler foundation target, drop event smoothly
		}
		toStack = instance.State.Foundations[toIndex].Stack
	case state.StackTableau:
		if toIndex < 0 || toIndex >= len(instance.State.Tableau) {
			return // Out of bounds filler grid target, drop event smoothly
		}
		toStack = instance.State.Tableau[toIndex].Stack
	case state.StackTalon:
		toStack = instance.State.Talon.Stock
	case state.StackWaste:
		toStack = instance.State.Talon.Waste
	case state.StackReserve:
		if toIndex < 0 || toIndex >= len(instance.State.Reserves) {
			return // Out of bounds filler reserve target, drop event smoothly
		}
		toStack = instance.State.Reserves[toIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'toComponentType' %d", toComponentType))
	}

	change := instance.Game.Move(fromStack, toStack, instance.State.Tableau, instance.State.Reserves)

	instance.Game.Compact(instance.State.Talon.Stock, instance.State.Talon.Waste, instance.State.Tableau)

	if change {
		instance.History.Update(instance.State)
	}

	instance.updateDisplay()

	if instance.Game.HasWon(instance.State.Tableau, instance.State.Foundations) {
		// TODO: Add a score to display.
		const score = 100
		instance.Display.ShowWinnerModal(instance.Game.Name(), score)
	}
}

func (instance *Instance) undo() {
	instance.History.Undo(&instance.State)

	instance.updateDisplay()
}

// dealCards - deal cards to tableau.
func (instance *Instance) dealCards() {
	// Shuffle the cards.
	instance.State.Deck.Shuffle()

	tableauSpec := instance.Game.Tableau()
	reserveSpec := instance.Game.Reserves()

	// Deal the cards out onto the tableau.
	for idx := 0; idx < len(tableauSpec); idx++ {
		// Grab a copy of the existing rule on the stack and replace it with.
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being.
		// dealt most definitely do not adhere to it (the rule).
		rule := instance.State.Tableau[idx].Stack.Rule
		instance.State.Tableau[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
		numCards := tableauSpec[idx].CardCount[0]
		numOpen := tableauSpec[idx].CardCount[1]

		for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
			card := instance.State.Deck.Deal()
			instance.State.Tableau[idx].Stack.Add(card, false)
		}

		for openIdx := 0; openIdx < numOpen; openIdx++ {
			card := instance.State.Deck.Deal()

			if _, ok := instance.State.Tableau[idx].Stack.SkipCards[card]; ok {
				continue
			}

			instance.State.Tableau[idx].Stack.Add(card, true)
		}

		instance.State.Tableau[idx].Stack.ShowCount = tableauSpec[idx].ShowCount

		// Return the rule to its correct state.
		instance.State.Tableau[idx].Stack.Rule = rule
	}

	// Deal cards to any reserves.
	for idx := 0; idx < len(reserveSpec); idx++ {
		// Grab a copy of the existing rule on the stack and replace it with.
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being.
		// dealt most definitely do not adhere to it (the rule).
		rule := instance.State.Reserves[idx].Stack.Rule
		instance.State.Reserves[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
		numCards := reserveSpec[idx].CardCount[0]
		numOpen := reserveSpec[idx].CardCount[1]

		for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
			card := instance.State.Deck.Deal()
			instance.State.Reserves[idx].Stack.Add(card, false)
		}

		for openIdx := 0; openIdx < numOpen; openIdx++ {
			card := instance.State.Deck.Deal()
			instance.State.Reserves[idx].Stack.Add(card, true)
		}

		// Return the rule to its correct state.
		instance.State.Reserves[idx].Stack.Rule = rule
	}

	// Put one card onto the Waste.
	if instance.Game.Talon() {
		card := instance.State.Deck.Deal()
		instance.State.Talon.Waste.Add(card, true)
	}

	if instance.Game.FoundationBase() {
		card := instance.State.Deck.Deal()

		// Determine which pile the first card actually lands on based on suit
		toChange := int(card.Suit)

		// Loop through ALL foundations (handles 4 for Agnes, 8 for American Toad)
		for i := 0; i < len(instance.State.Foundations); i++ {
			// Calculate suit for this specific foundation (0=H, 1=D, 2=C, 3=S, 4=H...)
			suit := state.Suit(i % 4)
			suitedBase := state.SuitedCard{Rank: card.Rank, Suit: suit}

			// Set the base on both the Foundation and the underlying Stack
			instance.State.Foundations[i].Base = suitedBase
			instance.State.Foundations[i].Stack.Base = suitedBase
		}

		// Add the starter card to the correct pile
		backUpRule := instance.State.Foundations[toChange].Stack.Rule
		instance.State.Foundations[toChange].Stack.Rule = func(state.SuitedCard) bool { return true }
		instance.State.Foundations[toChange].Stack.Add(card, true)
		instance.State.Foundations[toChange].Stack.Rule = backUpRule
	}

	// Put the rest of the cards onto the talon.
	for instance.State.Deck.Len() != 0 {
		card := instance.State.Deck.Deal()
		instance.State.Talon.Stock.Add(card, false)
	}

	// Create the first history.
	instance.History.Update(instance.State)
}

func (instance *Instance) onHints() []state.Move {
	// TODO - fix this hack.
	hints := instance.Game.AvailableMoves(
		func() []*state.Tableau {
			values := make([]*state.Tableau, len(instance.State.Tableau))
			for i, p := range instance.State.Tableau {
				values[i] = p
			}

			return values
		}(),
		func() []*state.Foundation {
			values := make([]*state.Foundation, len(instance.State.Foundations))
			for i, p := range instance.State.Foundations {
				values[i] = p
			}

			return values
		}(),
		func() []*state.Talon {
			return []*state.Talon{instance.State.Talon}
		}(),
		func() []*state.Reserve {
			values := make([]*state.Reserve, len(instance.State.Reserves))
			for i, p := range instance.State.Reserves {
				values[i] = p
			}

			return values
		}(),
	)

	return hints
}

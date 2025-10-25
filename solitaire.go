package solitaire

import (
	"fmt"
	"log"

	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen"
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
	// State of the current game.
	State state.State

	// Track first selection for move operations.
	firstSelection bool
	firstIndex     int

	// History of this instance.
	history state.History
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

	// Switch to the game page.
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
		fromStack = instance.State.Foundations[fromIndex].Stack
	case state.StackTableau:
		fromStack = instance.State.Tableau[fromIndex].Stack
	case state.StackTalon:
		fromStack = instance.State.Talon.Stock
		if !instance.Game.Talon() {
			return
		}
	case state.StackWaste:
		fromStack = instance.State.Talon.Waste
	case state.StackReserve:
		fromStack = instance.State.Reserves[fromIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'fromComponentType' %d", fromComponentType))
	}

	var toStack *state.Stack
	switch toComponentType {
	case state.StackFoundation:
		toStack = instance.State.Foundations[toIndex].Stack
	case state.StackTableau:
		toStack = instance.State.Tableau[toIndex].Stack
	case state.StackTalon:
		toStack = instance.State.Talon.Stock
	case state.StackWaste:
		toStack = instance.State.Talon.Waste
	case state.StackReserve:
		toStack = instance.State.Reserves[toIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'toComponentType' %d", toComponentType))
	}

	instance.Game.Move(fromStack, toStack, instance.State.Tableau)

	instance.Game.Compact(instance.State.Talon.Stock, instance.State.Talon.Waste, instance.State.Tableau)

	instance.history.Update()

	instance.updateDisplay()

	if instance.Game.HasWon(instance.State.Tableau, instance.State.Foundations) {
		// TODO: Add a score to display.
		const score = 100
		instance.Display.ShowWinnerModal(instance.Game.Name(), score)
	}
}

// dealCards - deal cards to tableau.
func (instance *Instance) dealCards() {
	// Shuffle the cards.
	instance.State.Deck.Shuffle()

	tableauSpec := instance.Game.Tableau()
	reserveSpec := instance.Game.Reserves()

	// Deal the cards out onto the tableau.
	for idx := 0; idx < len(tableauSpec); idx++ {
		// Grab a copy of the existing rule on the stack and replace it with
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being
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
		// Grab a copy of the existing rule on the stack and replace it with
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being
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
		toChange := 99
		switch card.Suit {
		case state.Hearts:
			toChange = 0
		case state.Diamonds:
			toChange = 1
		case state.Clubs:
			toChange = 2
		case state.Spades:
			toChange = 3
		}

		backUpRule := instance.State.Foundations[toChange].Stack.Rule
		instance.State.Foundations[toChange].Stack.Rule = func(state.SuitedCard) bool { return true }
		instance.State.Foundations[toChange].Stack.Add(card, true)
		instance.State.Foundations[0].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Hearts}
		instance.State.Foundations[1].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Diamonds}
		instance.State.Foundations[2].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Clubs}
		instance.State.Foundations[3].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Spades}
		instance.State.Foundations[0].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Hearts}
		instance.State.Foundations[1].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Diamonds}
		instance.State.Foundations[2].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Clubs}
		instance.State.Foundations[3].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Spades}
		instance.State.Foundations[toChange].Stack.Rule = backUpRule
	}

	// Put the rest of the cards onto the talon.
	for instance.State.Deck.Len() != 0 {
		card := instance.State.Deck.Deal()
		instance.State.Talon.Stock.Add(card, false)
	}
}

func (instance *Instance) onHints() {
	//TODO - fix this hack
	hints := instance.Game.AvailableMoves(
		func() []state.Tableau {
			values := make([]state.Tableau, len(instance.State.Tableau))
			for i, p := range instance.State.Tableau {
				values[i] = *p
			}
			return values
		}(),
		func() []state.Foundation {
			values := make([]state.Foundation, len(instance.State.Foundations))
			for i, p := range instance.State.Foundations {
				values[i] = *p
			}
			return values
		}(),
		func() []state.Talon {
			return []state.Talon{*instance.State.Talon}
		}(),
		func() []state.Reserve {
			values := make([]state.Reserve, len(instance.State.Reserves))
			for i, p := range instance.State.Reserves {
				values[i] = *p
			}
			return values
		}(),
	)

	log.Printf("HINTS %#v", hints)
}

package solitaire

import (
	"fmt"

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
	// Model.
	Foundations []*state.Foundation
	Tableau     []*state.Tableau
	Reserves    []*state.Reserve
	Talon       *state.Talon // This has both Stock and Waste stacks.
	Deck        *state.Deck

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

	instance.Display.SetGameSelectedCallback(instance.onGameSelected)

	// Set up component selection callback.
	instance.Display.SetComponentSelectedCallback(instance.onComponentSelected)

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
		fromStack = instance.Foundations[fromIndex].Stack
	case state.StackTableau:
		fromStack = instance.Tableau[fromIndex].Stack
	case state.StackTalon:
		fromStack = instance.Talon.Stock
	case state.StackWaste:
		fromStack = instance.Talon.Waste
	case state.StackReserve:
		fromStack = instance.Reserves[fromIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'fromComponentType' %d", fromComponentType))
	}

	var toStack *state.Stack
	switch toComponentType {
	case state.StackFoundation:
		toStack = instance.Foundations[toIndex].Stack
	case state.StackTableau:
		toStack = instance.Tableau[toIndex].Stack
	case state.StackTalon:
		toStack = instance.Talon.Stock
	case state.StackWaste:
		toStack = instance.Talon.Waste
	case state.StackReserve:
		toStack = instance.Reserves[toIndex].Stack
	default:
		panic(fmt.Sprintf("Got impossible 'toComponentType' %d", toComponentType))
	}

	instance.Game.Move(fromStack, toStack, instance.Tableau)

	instance.Game.Compact(instance.Talon.Stock, instance.Talon.Waste, instance.Tableau)

	instance.updateDisplay()

	if instance.Game.HasWon(instance.Tableau, instance.Foundations) {
		// TODO: Add a score to display.
		const score = 100
		instance.Display.ShowWinnerModal(instance.Game.Name(), score)
	}
}

// dealCards - deal cards to tableau.
func (instance *Instance) dealCards() {
	// Shuffle the cards.
	instance.Deck.Shuffle()

	tableauSpec := instance.Game.Tableau()
	reserveSpec := instance.Game.Reserves()

	// Deal the cards out onto the tableau.
	for idx := 0; idx < len(tableauSpec); idx++ {
		// Grab a copy of the existing rule on the stack and replace it with
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being
		// dealt most definitely do not adhere to it (the rule).
		rule := instance.Tableau[idx].Stack.Rule
		instance.Tableau[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
		numCards := tableauSpec[idx].CardCount[0]
		numOpen := tableauSpec[idx].CardCount[1]

		for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
			card := instance.Deck.Deal()
			instance.Tableau[idx].Stack.Add(card, false)
		}

		for openIdx := 0; openIdx < numOpen; openIdx++ {
			card := instance.Deck.Deal()

			if _, ok := instance.Tableau[idx].Stack.SkipCards[card]; ok {
				continue
			}

			instance.Tableau[idx].Stack.Add(card, true)
		}

		// Return the rule to its correct state.
		instance.Tableau[idx].Stack.Rule = rule
	}

	// Deal cards to any reserves.
	for idx := 0; idx < len(reserveSpec); idx++ {
		// Grab a copy of the existing rule on the stack and replace it with
		// one that will allow us to deal anything.
		// FTR the existing rule prevents a deal because the cards being
		// dealt most definitely do not adhere to it (the rule).
		rule := instance.Reserves[idx].Stack.Rule
		instance.Reserves[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
		numCards := reserveSpec[idx].CardCount[0]
		numOpen := reserveSpec[idx].CardCount[1]

		for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
			card := instance.Deck.Deal()
			instance.Reserves[idx].Stack.Add(card, false)
		}

		for openIdx := 0; openIdx < numOpen; openIdx++ {
			card := instance.Deck.Deal()
			instance.Reserves[idx].Stack.Add(card, true)
		}

		// Return the rule to its correct state.
		instance.Reserves[idx].Stack.Rule = rule
	}

	// Put one card onto the Waste.
	if instance.Game.Talon() {
		card := instance.Deck.Deal()
		instance.Talon.Waste.Add(card, true)
	}

	if instance.Game.FoundationBase() {
		card := instance.Deck.Deal()
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

		backUpRule := instance.Foundations[toChange].Stack.Rule
		instance.Foundations[toChange].Stack.Rule = func(state.SuitedCard) bool { return true }
		instance.Foundations[toChange].Stack.Add(card, true)
		instance.Foundations[0].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Hearts}
		instance.Foundations[1].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Diamonds}
		instance.Foundations[2].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Clubs}
		instance.Foundations[3].Base = state.SuitedCard{Rank: card.Rank, Suit: state.Spades}
		instance.Foundations[0].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Hearts}
		instance.Foundations[1].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Diamonds}
		instance.Foundations[2].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Clubs}
		instance.Foundations[3].Stack.Base = state.SuitedCard{Rank: card.Rank, Suit: state.Spades}
		instance.Foundations[toChange].Stack.Rule = backUpRule
	}

	// Put the rest of the cards onto the talon.
	for instance.Deck.Len() != 0 {
		card := instance.Deck.Deal()
		instance.Talon.Stock.Add(card, false)
	}
}

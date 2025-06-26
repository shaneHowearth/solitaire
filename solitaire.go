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
	Display     screen.Display
	Game        game.Variant
	Foundations []state.Foundation
	Tableau     []*state.Tableau
	Talon       *state.Talon
}

// New - create a new instance.
func New(display screen.Display) *Instance {
	return &Instance{
		Display: display,
	}
}

// Start - start the game.
func (instance *Instance) Start() {
	instance.Display.Show("Games")
	instance.ChooseGame()

	if instance.Game != nil {
		numFoundations, foundationBase, foundationRule := instance.Game.Foundations()
		numTableau, _, tableauRule := instance.Game.Tableau()
		gameState := state.New(
			instance.Game.Decks(),
			numFoundations,
			foundationBase,
			foundationRule,
			numTableau,
			tableauRule,
			1,
			1,
			// Talon rule is to allow everythign to be added to its stacks.
			func(state.SuitedCard) bool { return true },
		)

		instance.Foundations = gameState.Foundations
		instance.Tableau = gameState.Tableau
		instance.Talon = gameState.Talon
		counts := instance.Game.SetupDeal()

		gameState.Deck.Shuffle()

		for idx := 0; idx < numTableau; idx++ {
			// Grab a copy of the existing rule on the stack and rplace it with
			// one that will allow us to deal.
			rule := instance.Tableau[idx].Stack.Rule
			instance.Tableau[idx].Stack.Rule = func(state.SuitedCard) bool { return true }
			countIdx := idx * 2
			numCards := counts[countIdx]
			numOpen := counts[countIdx+1]

			for dealIdx := 0; dealIdx < numCards-numOpen; dealIdx++ {
				card := gameState.Deck.Deal()
				instance.Tableau[idx].Add(card, false)
			}

			for openIdx := 0; openIdx < numOpen; openIdx++ {
				card := gameState.Deck.Deal()
				instance.Tableau[idx].Add(card, true)
			}

			// Return the rule to its correct state.
			instance.Tableau[idx].Stack.Rule = rule
		}

		instance.CreateBoard(instance.Game)

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

		// TODO - need this to be row/col?
		for idx := range instance.Tableau {
			instance.Display.TableauPrint(idx,
				instance.Tableau[idx].Stack.Cards(),
			)
		}

		instance.Display.Show(instance.Game.Name())
	}
}

// ChooseGame - Get the game choice from the user.
func (instance *Instance) ChooseGame() {
	instance.Game = instance.Display.GetSelected()
}

// CreateBoard - game display
// Create tableaus, foundations, and talons.
// Get the cards into the right places to begin.
func (instance *Instance) CreateBoard(game game.Variant) {
	tableauHeight, tableauWidth := game.TableauGridSize()
	foundationCount, foundationBase, _ := game.Foundations()
	instance.Display.CreateBoard(
		game.Name(),
		tableauHeight,
		tableauWidth,
		foundationCount,
		foundationBase,
	)
}

// React to each move that the user makes.
// Splash when winning, followed by asking if the user wants to quit, play the
// same variant again, or a new one.
// One day, ask if the user wants to save the current state on exit.
// TODO: Points?

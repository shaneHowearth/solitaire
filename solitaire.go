package solitaire

import (
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
	Tableau     state.Tableau
	Talon       state.Talon
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
		instance.CreateBoard(instance.Game)
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

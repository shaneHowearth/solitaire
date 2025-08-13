package screen

import (
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// Display - defines what a display of the game needs to do.
type Display interface {
	Run() error
	Show(name string)

	// Game selection.
	SetGameSelectedCallback(callback func(game.Variant))

	// Component selection callback
	SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int))

	// Board creation.
	CreateBoard(name string, tableauHeight, tableauWidth, foundationCount int, foundationBase state.Rank, howToPlay []string)

	// Display updates.
	FoundationTitle(num int, value string)
	FoundationPrint(num int, value []string)
	TableauPrint(idx int, value []string)
	TalonPrint(value []string)
	WastePrint(value []string)

	// Component selection.
	GetSelectedComponent() (state.StackType, int)
	ClearSelection()
	HasSelection() bool
}

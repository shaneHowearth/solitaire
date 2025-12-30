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

	// Component selection callback.
	SetComponentSelectedCallback(callback func(state.StackType, int, state.StackType, int))

	// Game Redeal callback.
	SetGameRedealCallback(callback func())

	// Game Undo callback.
	SetGameUndoCallback(callback func())

	// Board creation.
	CreateBoard(
		name string,
		tableauHeight, tableauWidth, reserveCount, foundationCount int,
		howToPlay []string,
	)

	// Display updates.
	FoundationTitle(num int, value string)
	FoundationPrint(num int, value []string)
	ReservePrint(idx int, value []string)
	TableauPrint(idx int, value []string, showCount int)
	TalonPrint(value []string)
	WastePrint(value []string)

	// Component selection.
	GetSelectedComponent() (state.StackType, int)
	ClearSelection()
	HasSelection() bool

	// Game Hint.
	SetHintsCallback(callback func() []state.Move)

	// Winner modal.
	ShowWinnerModal(string, int)
}

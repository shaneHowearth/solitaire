package screen

import (
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// ComponentType represents the type of game component
type ComponentType int

const (
	ComponentFoundation ComponentType = iota
	ComponentTableau
	ComponentTalon
	ComponentWaste
)

// Display - defines what a display of the game needs to do.
type Display interface {
	Run() error
	Show(name string)

	// Game selection.
	SetGameSelectedCallback(callback func(game.Variant))

	// Component selection callback
	SetComponentSelectedCallback(callback func(ComponentType, int))

	// Board creation.
	CreateBoard(name string, tableauHeight, tableauWidth, foundationCount int, foundationBase state.Rank)

	// Display updates.
	FoundationTitle(num int, value string)
	FoundationPrint(num int, value []string)
	TableauPrint(idx int, value []string)

	// Component selection.
	GetSelectedComponent() (ComponentType, int)
	ClearSelection()
	HasSelection() bool
}

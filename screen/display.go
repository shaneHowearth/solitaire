package screen

import (
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
)

// Display - defines what a display of the game needs to do.
type Display interface {
	Show(string)
	GetSelected() game.Variant
	CreateBoard(name string, tableauHeight, tableauWidth, foundationCount int, foundationBase state.Rank)
}

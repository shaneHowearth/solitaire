package screen

import "github.com/shanehowearth/solitaire/game"

// Display - defines what a display of the game needs to do.
type Display interface {
	// Splash screen - opportunity to show user all of the variants available.
	Splash([]game.Variant)
}

package main

import (
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/tui"
)

func main() {
	// Which variants of solitaire are available to play.
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})

	// select which display is going to be used.
	display := tui.Display{}

	// Splash screen - has options for user to choose which variant of solitaire
	// they wish to play.
	display.Splash(variants)

}

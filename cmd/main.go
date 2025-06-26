package main

import (
	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/tui"
)

func main() {
	// Which variants of state are available to play.
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})
	variants = append(variants, &game.Klondike2{})

	instance := solitaire.New(tui.New(variants))

	instance.Start()
}

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

	instance := &solitaire.Instance{}

	instance.Display = tui.New(variants)
	// select which display is going to be used.

	instance.Display.Show("Games")
	instance.ChooseGame(variants)
	if instance.Game != nil {
		instance.CreateBoard(instance.Game)
		instance.Display.Show(instance.Game.Name())
	}
}

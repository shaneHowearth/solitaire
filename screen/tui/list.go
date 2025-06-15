package tui

import (
	"fmt"

	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/game"
)

type Display struct{}

func (*Display) Splash(games []game.Variant) {
	app := tview.NewApplication()
	list := tview.NewList()
	for idx, game := range games {
		list.AddItem(game.Name(), "", []rune(fmt.Sprintf("%d", idx+1))[0], game.SetupDeal)
	}

	if err := app.SetRoot(list, true).SetFocus(list).Run(); err != nil {
		panic(err)
	}
}

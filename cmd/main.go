package main

import (
	"fmt"

	"github.com/shanehowearth/solitaire"
)

func main() {
	for i := 1; i < solitaire.SuitCount; i++ {
		deck := solitaire.CreateDecks(i)
		fmt.Println(deck)
		fmt.Println("-----------------------------------------")
		deck.Shuffle()
		fmt.Println(deck)
		fmt.Println("-----------------------------------------")
	}
}

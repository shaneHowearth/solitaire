package solitaire

import (
	"slices"
)

// Deck - all the cards for the game, will always be a multiple of 52.
type Deck []SuitedCard

var deck = []SuitedCard{
	{Ace, Spade}, {Ace, Heart}, {Ace, Club}, {Ace, Diamond},
	{Two, Spade}, {Two, Heart}, {Two, Club}, {Two, Diamond},
	{Three, Spade}, {Three, Heart}, {Three, Club}, {Three, Diamond},
	{Four, Spade}, {Four, Heart}, {Four, Club}, {Four, Diamond},
	{Five, Spade}, {Five, Heart}, {Five, Club}, {Five, Diamond},
	{Six, Spade}, {Six, Heart}, {Six, Club}, {Six, Diamond},
	{Seven, Spade}, {Seven, Heart}, {Seven, Club}, {Seven, Diamond},
	{Eight, Spade}, {Eight, Heart}, {Eight, Club}, {Eight, Diamond},
	{Nine, Spade}, {Nine, Heart}, {Nine, Club}, {Nine, Diamond},
	{Ten, Spade}, {Ten, Heart}, {Ten, Club}, {Ten, Diamond},
	{Jack, Spade}, {Jack, Heart}, {Jack, Club}, {Jack, Diamond},
	{Queen, Spade}, {Queen, Heart}, {Queen, Club}, {Queen, Diamond},
	{King, Spade}, {King, Heart}, {King, Club}, {King, Diamond},
}

// CreateDecks - Create a deck of cards.
func CreateDecks(n int) Deck {
	if n < 1 {
		panic("cannot create zero decks")
	}

	return slices.Repeat(deck, n)
}

package solitaire

import (
	"math/rand/v2"
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
func CreateDecks(n int) *Deck {
	if n < 1 {
		panic("cannot create zero decks")
	}

	d := Deck(slices.Repeat(deck, n))
	return &d
}

// Len - Length of the Deck.
func (usedDeck *Deck) Len() int {
	if usedDeck == nil {
		return 0
	}

	return len((*usedDeck))
}

// Shuffle - Shuffles deck in place.
func (usedDeck *Deck) Shuffle() {
	rand.Shuffle(usedDeck.Len(), func(i, j int) {
		(*usedDeck)[i], (*usedDeck)[j] = (*usedDeck)[j], (*usedDeck)[i]
	})
}

// Deal returns and removes the final card in the deck.
func (usedDeck *Deck) Deal() SuitedCard {
	if usedDeck.Len() == 0 {
		// Empty, shouldn't be here.
	}

	card := (*usedDeck)[(usedDeck.Len() - 1)]

	*usedDeck = slices.Delete(*usedDeck, usedDeck.Len()-1, usedDeck.Len())

	return card
}

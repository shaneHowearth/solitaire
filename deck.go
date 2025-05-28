package solitaire

import (
	"math/rand/v2"
	"slices"
)

// Deck - all the cards for the game, will always be a multiple of 52.
type Deck []SuitedCard

var deck = []SuitedCard{
	{Card: Ace, Suit: Spades}, {Card: Ace, Suit: Hearts}, {Card: Ace, Suit: Clubs}, {Card: Ace, Suit: Diamonds},
	{Card: Two, Suit: Spades}, {Card: Two, Suit: Hearts}, {Card: Two, Suit: Clubs}, {Card: Two, Suit: Diamonds},
	{Card: Three, Suit: Spades}, {Card: Three, Suit: Hearts}, {Card: Three, Suit: Clubs}, {Card: Three, Suit: Diamonds},
	{Card: Four, Suit: Spades}, {Card: Four, Suit: Hearts}, {Card: Four, Suit: Clubs}, {Card: Four, Suit: Diamonds},
	{Card: Five, Suit: Spades}, {Card: Five, Suit: Hearts}, {Card: Five, Suit: Clubs}, {Card: Five, Suit: Diamonds},
	{Card: Six, Suit: Spades}, {Card: Six, Suit: Hearts}, {Card: Six, Suit: Clubs}, {Card: Six, Suit: Diamonds},
	{Card: Seven, Suit: Spades}, {Card: Seven, Suit: Hearts}, {Card: Seven, Suit: Clubs}, {Card: Seven, Suit: Diamonds},
	{Card: Eight, Suit: Spades}, {Card: Eight, Suit: Hearts}, {Card: Eight, Suit: Clubs}, {Card: Eight, Suit: Diamonds},
	{Card: Nine, Suit: Spades}, {Card: Nine, Suit: Hearts}, {Card: Nine, Suit: Clubs}, {Card: Nine, Suit: Diamonds},
	{Card: Ten, Suit: Spades}, {Card: Ten, Suit: Hearts}, {Card: Ten, Suit: Clubs}, {Card: Ten, Suit: Diamonds},
	{Card: Jack, Suit: Spades}, {Card: Jack, Suit: Hearts}, {Card: Jack, Suit: Clubs}, {Card: Jack, Suit: Diamonds},
	{Card: Queen, Suit: Spades}, {Card: Queen, Suit: Hearts}, {Card: Queen, Suit: Clubs}, {Card: Queen, Suit: Diamonds},
	{Card: King, Suit: Spades}, {Card: King, Suit: Hearts}, {Card: King, Suit: Clubs}, {Card: King, Suit: Diamonds},
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
		// TODO: better to send an error back.
		return SuitedCard{}
	}

	card := (*usedDeck)[(usedDeck.Len() - 1)]

	*usedDeck = slices.Delete(*usedDeck, usedDeck.Len()-1, usedDeck.Len())

	return card
}

// Top - Length of the Deck.
func (usedDeck *Deck) Top() SuitedCard {
	return (*usedDeck)[usedDeck.Len()-1]
}

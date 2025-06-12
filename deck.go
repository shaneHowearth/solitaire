package solitaire

import (
	"math/rand/v2"
	"slices"
)

// Deck - all the cards for the game, will always be a multiple of 52.
type Deck []SuitedCard

var deck = []SuitedCard{
	{Rank: Ace, Suit: Spades}, {Rank: Ace, Suit: Hearts}, {Rank: Ace, Suit: Clubs}, {Rank: Ace, Suit: Diamonds},
	{Rank: Two, Suit: Spades}, {Rank: Two, Suit: Hearts}, {Rank: Two, Suit: Clubs}, {Rank: Two, Suit: Diamonds},
	{Rank: Three, Suit: Spades}, {Rank: Three, Suit: Hearts}, {Rank: Three, Suit: Clubs}, {Rank: Three, Suit: Diamonds},
	{Rank: Four, Suit: Spades}, {Rank: Four, Suit: Hearts}, {Rank: Four, Suit: Clubs}, {Rank: Four, Suit: Diamonds},
	{Rank: Five, Suit: Spades}, {Rank: Five, Suit: Hearts}, {Rank: Five, Suit: Clubs}, {Rank: Five, Suit: Diamonds},
	{Rank: Six, Suit: Spades}, {Rank: Six, Suit: Hearts}, {Rank: Six, Suit: Clubs}, {Rank: Six, Suit: Diamonds},
	{Rank: Seven, Suit: Spades}, {Rank: Seven, Suit: Hearts}, {Rank: Seven, Suit: Clubs}, {Rank: Seven, Suit: Diamonds},
	{Rank: Eight, Suit: Spades}, {Rank: Eight, Suit: Hearts}, {Rank: Eight, Suit: Clubs}, {Rank: Eight, Suit: Diamonds},
	{Rank: Nine, Suit: Spades}, {Rank: Nine, Suit: Hearts}, {Rank: Nine, Suit: Clubs}, {Rank: Nine, Suit: Diamonds},
	{Rank: Ten, Suit: Spades}, {Rank: Ten, Suit: Hearts}, {Rank: Ten, Suit: Clubs}, {Rank: Ten, Suit: Diamonds},
	{Rank: Jack, Suit: Spades}, {Rank: Jack, Suit: Hearts}, {Rank: Jack, Suit: Clubs}, {Rank: Jack, Suit: Diamonds},
	{Rank: Queen, Suit: Spades}, {Rank: Queen, Suit: Hearts}, {Rank: Queen, Suit: Clubs}, {Rank: Queen, Suit: Diamonds},
	{Rank: King, Suit: Spades}, {Rank: King, Suit: Hearts}, {Rank: King, Suit: Clubs}, {Rank: King, Suit: Diamonds},
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

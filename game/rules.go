package game

import "github.com/shanehowearth/solitaire"

// PlusOneRule - A rule for foundations where the card being added must be the
// same suit, and be of a rank one higher than the existing top card.
var PlusOneRule = func(foundation solitaire.Foundation, card solitaire.SuitedCard) bool {
	// Handle when the foundation is empty.
	if foundation.Len() == 0 {
		if card.Suit == foundation.Base.Suit && card.Rank == foundation.Base.Rank {
			return true
		}
	}

	// Get the card currently at the top of the foundation.
	topCard, err := foundation.Top()
	if err != nil {
		return false
	}

	// If the card is the same suit, and is one up in rank
	// then it can go onto the foundation.
	if card.Suit == foundation.Base.Suit && (card.Rank-topCard.Rank) == 1 {
		return true
	}

	// All other cases the card should not be added to the foundation.
	return false
}

// MinusOneRule - A rule for tableau where the card being added must be the
// opposite colour, and a rank of one less, than the existing top card on that
// tableau.
var MinusOneRule = func(tableau solitaire.Tableau, card solitaire.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		if card.Rank == tableau.Base {
			return true
		}
	}

	// Get the card currently at the top of the tableau.
	topCard, err := tableau.Top()
	if err != nil {
		return false
	}

	// If the card is the opposite colour, and is one down in rank
	// then it can go onto the tableau.
	if ((card.Suit+topCard.Suit)%2 == 1) && (topCard.Rank-card.Rank) == 1 {
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

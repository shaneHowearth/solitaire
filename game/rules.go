package game

import "github.com/shanehowearth/solitaire/state"

// PlusOneRule - A rule for stacks where the card being added must be the.
// same suit, and be of a rank one higher than the existing top card.
var PlusOneRule = func(foundation *state.Stack, card state.SuitedCard) bool {
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

	// If the card is the same suit, and is one up in rank.
	// then it can go onto the foundation.
	if card.Suit == foundation.Base.Suit && (card.Rank-topCard.Rank) == 1 {
		return true
	}

	// All other cases the card should not be added to the foundation.
	return false
}

// MinusOneRule - A rule for tableau where the card being added must be the.
// opposite colour, and a rank of one less, than the existing top card on that.
// tableau.
var MinusOneRule = func(tableau *state.Stack, card state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if (*tableau).Len() == 0 {
		if card.Rank == tableau.Base.Rank {
			return true
		}
	}

	// Get the card currently at the top of the tableau.
	topCard, err := (*tableau).Top()
	if err != nil {
		return false
	}

	// If the card is the opposite colour, and is one down in rank.
	// then it can go onto the tableau.
	if ((card.Suit+topCard.Suit)%2 == 1) && (topCard.Rank-card.Rank) == 1 {
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// PlusOneWraparoundRule handles sequences like Q-K-A-2.
func PlusOneWraparoundRule(foundation *state.Stack, card state.SuitedCard) bool {
	if foundation.Len() == 0 {
		return true
	}
	top, _ := foundation.Top()

	// Ensure suits match
	if top.Suit != card.Suit {
		return false
	}

	// Logic: (TopRank + 1) % 13 == CardRank % 13
	// Using the next rank in the sequence with wraparound
	nextRank := state.Rank((int(top.Rank) + 1) % 13)

	return card.Rank == nextRank
}

func isRed(s state.Suit) bool {
	return s == state.Hearts || s == state.Diamonds
}

func MinusOneWraparoundRule(s *state.Stack, c state.SuitedCard) bool {
	if s.Len() == 0 {
		return true
	}
	top, _ := s.Top()

	// Alternating colors check
	if isRed(top.Suit) == isRed(c.Suit) {
		return false
	}

	// Logic: (TopRank - 1) % 13.
	targetRank := state.Rank((int(top.Rank) + 12) % 13)
	return c.Rank == targetRank
}

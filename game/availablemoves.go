package game

// Helpers for checking if a move is available, or not.

import (
	"fmt"

	"github.com/shanehowearth/solitaire/state"
)

// checkMove checks if a move is possible and returns a formatted hint string
func checkMove(source, destination *state.Stack, sourceName, destName string) string {
	canMove, numCards := CanMove(source, destination, false)

	if !canMove || numCards == 0 {
		return ""
	}

	bottomCard, err := getBottomMovingCard(source, numCards)
	if err != nil {
		return ""
	}

	cardStr := fmt.Sprintf("%s %s", bottomCard.Rank.String(), bottomCard.Suit.String())

	if numCards == 1 {
		return fmt.Sprintf("Move %s from %s to %s", cardStr, sourceName, destName)
	}

	return fmt.Sprintf("Move %s (+ %d card(s)) from %s to %s",
		cardStr, numCards-1, sourceName, destName)
}

// getBottomMovingCard returns the bottom card of the sequence that would be moved
func getBottomMovingCard(source *state.Stack, numCards int) (state.SuitedCard, error) {
	sourceClone := source.Clone()

	var bottomCard state.SuitedCard
	for i := 0; i < numCards; i++ {
		card, err := sourceClone.Top()
		if err != nil {
			return state.SuitedCard{}, err
		}
		bottomCard = card
		sourceClone.Deal()
	}

	return bottomCard, nil
}

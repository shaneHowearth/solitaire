package game

// Helpers for checking if a move is available, or not.

import (
	"errors"

	"github.com/shanehowearth/solitaire/state"
)

// checkMove checks if a move is possible and returns a state.Move.
func checkMove(
	source, destination *state.Stack,
	keepSequence, noKings bool,
) state.Move {
	canMove, numCards := CanMove(source, destination, keepSequence)

	if !canMove || numCards == 0 {
		return state.Move{}
	}

	// Ignore moving piles to an empty stack, and leaving an empty stack.
	// eg. King... to an empty stack.
	if noKings && destination.Len() == 0 && numCards == source.Len() {
		return state.Move{}
	}

	destinationTop, err := destination.Top()
	if err != nil {
		if !errors.Is(err, state.ErrEmpty) { //nolint:staticcheck // This no-op is here for documentation purposes.
			// No-op
			// If the destination is not allowed to be empty that's handled
			// above.
		} else {
			return state.Move{}
		}
	}

	sourceBottomCard, err := getBottomMovingCard(source, numCards)
	if err != nil {
		return state.Move{}
	}

	return state.Move{
		Source:                *source,
		Destination:           *destination,
		NumberMoving:          numCards,
		SourceCardTop:         sourceBottomCard,
		DestinationCardBottom: destinationTop,
	}
}

// getBottomMovingCard returns the bottom card of the sequence that would be moved.
func getBottomMovingCard(source *state.Stack, numCards int) (state.SuitedCard, error) {
	sourceClone := source.Clone()

	var bottomCard state.SuitedCard

	for i := 0; i < numCards; i++ {
		card, err := sourceClone.Top()
		if err != nil {
			return state.SuitedCard{}, err
		}

		bottomCard = card

		_, _ = sourceClone.Deal()
	}

	return bottomCard, nil
}

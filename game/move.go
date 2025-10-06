package game

import (
	"log"

	"github.com/shanehowearth/solitaire/state"
)

// CanMove checks if cards can be moved from source to destination without modifying the original stacks.
// Returns whether the move is valid and the number of cards that would be moved.
func CanMove(source, destination *state.Stack, keepSequence bool) (bool, int) {
	if destination == nil {
		return false, 0
	}

	if destination.Type == state.StackReserve {
		// No cards can be moved ONTO a reserve.
		return false, 0
	}

	if destination.Type == state.StackWaste && source.Type != state.StackTalon {
		// Waste can only receive cards from the stock.
		return false, 0
	}

	// Clone both stacks to avoid modifying originals
	sourceClone := source.Clone()
	destClone := destination.Clone()

	// Temporary stack that will hold cards that will be moved.
	temp := state.NewStack(
		sourceClone.Len(),
		state.SuitedCard{},
		func(*state.Stack) func(state.SuitedCard) bool {
			return func(state.SuitedCard) bool {
				return true
			}
		},
		state.StackUndefined,
	)

	count := 0
	canMove := false

	if destination.Type == state.StackTalon {
		// Stock can only be given cards when it is empty.
		if destination.Len() != 0 {
			return false, 0
		}

		// Stock can only be given cards from the Deck (first deal) or the
		// Waste. As the Deck isn't typed, we'll exclude the foundation and
		// tableau from being able to add to the stock.
		if source.Type == state.StackFoundation || source.Type == state.StackTableau {
			return false, 0
		}
	}

	// This loop attempts to find the group of cards that may be moved.
	for {
		sourceTop, err := sourceClone.Top()
		if err != nil {
			break
		}

		if !sourceTop.Visible && (sourceClone.Type != state.StackTalon) {
			break
		}

		count++

		temp.Add(sourceTop, true)

		_, err = sourceClone.Deal()
		if err != nil {
			log.Printf("Stack Deal err %v", err)
		}

		if destination.Rule(sourceTop) && destination.Type != state.StackTalon {
			canMove = true
			break
		}

		// Only tableau or talon can have more than 1 cards moved at once.
		if destination.Type != state.StackTableau && destination.Type != state.StackTalon {
			break
		}

		// The waste can only move multiple cards to the Talon.
		if source.Type == state.StackWaste && destination.Type != state.StackTalon {
			break
		}

		// The Talon can never move multiple cards.
		if source.Type == state.StackTalon {
			break
		}
	}

	if source.Type == state.StackWaste && destination.Type == state.StackTalon {
		if destination.CanReceiveMore() {
			canMove = true
		} else {
			canMove = false
		}
	}

	// Check that all the cards on the temp stack can be moved.
	// This validates the sequence if keepSequence is true
	temp2 := state.NewStack(
		sourceClone.Len(),
		state.SuitedCard{},
		func(*state.Stack) func(state.SuitedCard) bool {
			return func(state.SuitedCard) bool {
				return true
			}
		},
		state.StackUndefined,
	)

	if canMove && temp.Len() > 0 {
		// Clone the destination for testing
		testDest := destClone.Clone()

		// Try placing cards directly on the cloned destination
		sequenceValid := true
		if keepSequence {
			for i := temp.Len() - 1; i >= 0; i-- {
				card, _ := temp.Deal()
				temp2.Add(card, true)

				if !testDest.Rule(card) {
					sequenceValid = false
					break
				}

				testDest.Add(card, true) // Add to test next card
			}
		}

		// Save remaining cards from temp if we broke early
		for temp.Len() > 0 {
			card, _ := temp.Deal()
			temp2.Add(card, true)
		}
		canMove = sequenceValid
	}

	if canMove {
		return true, count
	}

	return false, 0
}

// Move - Move card(s) from one stack to another.
func Move(source, destination *state.Stack, keepSequence bool) bool {
	// First check if the move is valid
	canMove, numCards := CanMove(source, destination, keepSequence)

	if !canMove || numCards == 0 {
		return false
	}

	// Temporary stack to collect cards from source
	temp := state.NewStack(
		numCards,
		state.SuitedCard{},
		func(*state.Stack) func(state.SuitedCard) bool {
			return func(state.SuitedCard) bool {
				return true
			}
		},
		state.StackUndefined,
	)

	// Collect exactly numCards from the source stack
	for i := 0; i < numCards; i++ {
		card, err := source.Top()
		if err != nil {
			// This shouldn't happen since CanMove validated it
			log.Printf("Error getting card from source: %v", err)
			break
		}
		temp.Add(card, true)
		_, err = source.Deal()
		if err != nil {
			log.Printf("Error dealing card from source: %v", err)
		}
	}

	// Handle reversal for waste->talon moves
	if source.Type == state.StackWaste && destination.Type == state.StackTalon {
		temp.Reverse()
	}

	// Move all cards from temp to destination
	for temp.Len() > 0 {
		card, err := temp.Top()
		if err != nil {
			break
		}

		if destination.Type == state.StackTalon {
			destination.Add(card, false)
		} else {
			destination.Add(card, true)
		}

		_, _ = temp.Deal()
	}

	// Make the top card visible on source stack if needed
	if source.Type == state.StackTableau || source.Type == state.StackReserve {
		newTop, err := source.Top()
		if err != nil {
			// Source is empty, which is fine
			return true
		}

		// Remove and re-add to make it visible
		_, _ = source.Deal()
		source.Add(newTop, true)
	}

	return true
}

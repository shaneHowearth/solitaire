package game

import (
	"log"

	"github.com/shanehowearth/solitaire/state"
)

// Move - Move card(s) from one stack to another.
func Move(source, destination *state.Stack) bool {
	if destination == nil {
		return false
	}

	if destination.Type == state.StackReserve {
		// No cards can be moved ONTO a reserve.
		return false
	}

	// Can we move multiple cards?
	// Temporary stack that will hold cards that will be moved.
	temp := state.NewStack(
		source.Len(),
		state.SuitedCard{},
		func(state.SuitedCard) bool { return true },
		state.StackUndefined,
	)

	count := 0
	canMove := false

	if destination.Type == state.StackTalon {
		// Stock can only be given cards when it is empty.
		if destination.Len() != 0 {
			return false
		}

		// Stock can only be given cards from the Deck (first deal) or the
		// Waste. As the Deck isn't typed, we'll exclude the foundation and
		// tableau from being able to add to the stock.
		if source.Type == state.StackFoundation || source.Type == state.StackTableau {
			return false
		}
	}

	// This loop attempts to find the group of cards that may be moved.
	for {
		sourceTop, err := source.Top()
		if err != nil {
			break
		}

		if !sourceTop.Visible && (source.Type != state.StackTalon) {
			break
		}

		count++

		temp.Add(sourceTop, true)

		_, err = source.Deal()
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

			temp.Reverse()
		} else {
			canMove = false
		}
	}

	if !canMove {
		// Put the cards back and finish.
		for {
			savedRule := source.Rule
			source.Rule = func(state.SuitedCard) bool { return true }
			top, err := temp.Top()
			if err != nil {
				source.Rule = savedRule
				break
			}

			stackTop, _ := source.Top()
			if stackTop.Rank == top.Rank && stackTop.Suit == top.Suit {
				log.Printf("Got same card %v %v", stackTop, top)

				source.Rule = savedRule

				break
			}

			if source.Type == state.StackTalon {
				source.Add(top, false)
			} else {
				source.Add(top, true)
			}

			_, _ = temp.Deal()
		}
	} else {
		// If the card can be added to the destination add it, and drop it from the
		// source.
		for {
			top, err := temp.Top()
			if err != nil {
				break
			}

			if destination.Type == state.StackTalon {
				destination.Add(top, false)
			} else {
				destination.Add(top, true)
			}
			// Pull the top card off the source stack.
			_, _ = temp.Deal()
		}
		// Make the top card visible.
		if source.Type == state.StackTableau || source.Type == state.StackReserve {
			newTop, err := source.Top()
			if err != nil {
				return true
			}

			_, _ = source.Deal()
			source.Add(newTop, true)
		}
	}

	//

	return true
}

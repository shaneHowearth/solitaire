package state

import (
	"errors"
	"log"
	"slices"
)

// ErrEmpty - Error emitted when the stack is empty.
var ErrEmpty = errors.New("Empty")

// Add - add a suited card to the stack.
func (stack *Stack) Add(card SuitedCard, visible bool) {
	card.Visible = visible

	*stack.cards = append(*stack.cards, card)
}

// Deal removes and returns the final card in the deck.
func (stack *Stack) Deal() (SuitedCard, error) {
	if stack.Len() == 0 {
		// Empty, shouldn't be here.
		return SuitedCard{}, ErrEmpty
	}

	card := (*stack.cards)[stack.Len()-1]

	cards := make([]SuitedCard, stack.Len()-1)
	for idx := 0; idx < stack.Len()-1; idx++ {
		cards[idx] = (*stack.cards)[idx]
	}

	stack.cards = &cards

	return card, nil
}

// Move - Move card(s) from one stack to another.
func (stack *Stack) Move(destination *Stack, maxRedeals int) bool {
	if destination == nil {
		return false
	}

	// Nothing to do.
	if stack.cards == destination.cards {
		return false
	}

	// Can we move multiple cards?
	// Temporary stack that will hold cards that will be moved.
	temp := NewStack(
		15,
		func(SuitedCard) bool { return true },
		StackUndefined,
	)

	count := 0
	canMove := false

	if destination.Type == StackTalon {
		// Stock can only be given cards when it is empty.
		if destination.Len() != 0 {
			return false
		}

		// Stock can only be given cards from the Deck (first deal) or the
		// Stock. As the Deck isn't typed, we'll exclude the foundation and
		// tableau from being able to add to the stock.
		if stack.Type == StackFoundation || stack.Type == StackTableau {
			return false
		}
	}

	for {
		top, err := stack.Top()
		if err != nil {
			break
		}

		if !top.Visible && (stack.Type != StackTalon) {
			break
		}

		count++

		temp.Add(top, true)

		_, err = stack.Deal()
		if err != nil {
			log.Printf("Stack Deal err %v", err)
		}

		if destination.Rule(top) && destination.Type != StackTalon {
			canMove = true
			break
		}

		// Only tableau or talon can have more than 1 cards moved at once.
		if destination.Type != StackTableau && destination.Type != StackTalon {
			break
		}
	}

	if stack.Type == StackWaste && destination.Type == StackTalon {
		if destination.Received < maxRedeals || maxRedeals == -1 {
			canMove = true

			slices.Reverse(*temp.cards)

			destination.Received++
		} else {
			canMove = false
		}
	}

	if !canMove {
		savedRule := stack.Rule
		stack.Rule = func(SuitedCard) bool { return true }
		// Put the cards back and finish.
		for {
			top, err := temp.Top()
			if err != nil {
				stack.Rule = savedRule
				break
			}

			stackTop, _ := stack.Top()
			if stackTop.Rank == top.Rank && stackTop.Suit == top.Suit {
				log.Printf("Got same card %v %v", stackTop, top)

				stack.Rule = savedRule

				break
			}

			stack.Add(top, true)

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

			if destination.Type == StackTalon {
				destination.Add(top, false)
			} else {
				destination.Add(top, true)
			}
			// Pull the top card off the source stack.
			_, _ = temp.Deal()
		}
		// Make the top card visible.
		if stack.Type == StackTableau {
			newTop, err := stack.Top()
			if err != nil {
				return true
			}

			_, _ = stack.Deal()
			stack.Add(newTop, true)
		}
	}

	//

	return true
}

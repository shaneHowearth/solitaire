package state

import (
	"errors"
	"log"
)

// ErrEmpty - Error emitted when the stack is empty.
var ErrEmpty = errors.New("Empty")

// Add - add a suited card to the stack.
func (stack *Stack) Add(card SuitedCard, visible bool) {
	card.Visible = visible

	*stack.cards = append(*stack.cards, card)
}

// Deal returns and removes the final card in the deck.
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

// Move - Move card(s) on the tableau stack to the nominated pile (foundation or
// tableau).
func (stack *Stack) Move(destination *Stack) bool {
	if destination == nil {
		return false
	}

	if stack.cards == destination.cards {
		return false
	}

	// Can we move multiple cards?
	// Should I put the cards onto a temporary stack, then unwind if there's no
	// cards able to be moved?
	temp := NewStack(
		15,
		func(SuitedCard) bool { return true },
		StackUndefined,
	)

	count := 0
	canMove := false

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

		if destination.Rule(top) {
			canMove = true
			break
		}

		// Only tableau can have more than 1 cards moved at once.
		if stack.Type != StackTableau {
			break
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
			destination.Add(top, true)
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

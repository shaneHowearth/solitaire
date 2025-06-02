package solitaire

import (
	"errors"
	"slices"
)

// Stack (or Pile)
// A stack of cards one on top of the other and squared such that only the topmost card,
// whether face up or face down is visible.[5]
type Stack []SuitedCard

// ErrEmpty - Error emitted when the stack is empty.
var ErrEmpty = errors.New("Empty")

// There are two types of stack, the one that holds the reserve of cards that
// are yet to be played, and the ones that are on the tableau.

// Add - add a suited card to the stack.
func (stack *Stack) Add(card SuitedCard, visible bool) {
	card.Visible = visible

	*stack = append(*stack, card)
}

// Len - return the length of the stack.
func (stack *Stack) Len() int {
	return len((*stack))
}

// Top - the card that can be accessed immediately.
func (stack *Stack) Top() (SuitedCard, error) {
	if stack.Len() == 0 {
		return SuitedCard{}, ErrEmpty
	}

	return (*stack)[stack.Len()-1], nil
}

// Deal returns and removes the final card in the deck.
func (stack *Stack) Deal() (SuitedCard, error) {
	if stack.Len() == 0 {
		// Empty, shouldn't be here.
		// TODO: better to send an error back.
		return SuitedCard{}, ErrEmpty
	}

	card := (*stack)[(stack.Len() - 1)]

	*stack = slices.Delete(*stack, stack.Len()-1, stack.Len())

	return card, nil
}

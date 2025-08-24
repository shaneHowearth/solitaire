package state

import (
	"errors"
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

func (stack *Stack) Reverse() {
	slices.Reverse(*stack.cards)
}

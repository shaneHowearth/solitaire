package state

import (
	"errors"
	"fmt"
)

// Stack - A stack of cards one on top of the other and squared such that only
// the topmost card, whether face up or face down is visible.[5]
type Stack struct {
	cards *[]SuitedCard
	Rule  func(SuitedCard) bool
	Type  StackType
}

// StackType represents the type of game component
type StackType int

const (
	StackUndefined StackType = iota
	StackFoundation
	StackTableau
	StackTalon
	StackWaste
)

// ErrEmpty - Error emitted when the stack is empty.
var ErrEmpty = errors.New("Empty")

// There are two types of stack, the one that holds the reserve of cards that
// are yet to be played, and the ones that are on the tableau.

// NewStack - Create a new stack with an empty slice of SuitedCards that has a
// capacity of n.
func NewStack(number int, rule func(SuitedCard) bool, componentType StackType) *Stack {
	cards := make([]SuitedCard, 0, number)

	return &Stack{
		cards: &cards,
		Rule:  rule,
		Type:  componentType,
	}
}

// Add - add a suited card to the stack.
func (stack *Stack) Add(card SuitedCard, visible bool) {
	card.Visible = visible

	*stack.cards = append(*stack.cards, card)
}

// Len - return the length of the stack.
func (stack *Stack) Len() int {
	return len(*stack.cards)
}

// Top - the card that can be accessed immediately.
func (stack *Stack) Top() (SuitedCard, error) {
	if stack.Len() == 0 {
		return SuitedCard{}, ErrEmpty
	}

	return (*stack.cards)[stack.Len()-1], nil
}

// Deal returns and removes the final card in the deck.
func (stack *Stack) Deal() (SuitedCard, error) {
	if stack.Len() == 0 {
		// Empty, shouldn't be here.
		// TODO: better to send an error back.
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

// Move - Move n cards on the tableau stack to the nominated pile (foundation or
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
	temp := NewStack(15, func(SuitedCard) bool { return true }, StackUndefined)
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

		newTop, err := stack.Top()
		if err != nil {
			return true
		}
		_, _ = stack.Deal()
		stack.Add(newTop, true)
	}

	return true
}

const blankCard = "--"

// Cards - a string representation of the ards in the stack.
func (stack *Stack) Cards() []string {
	cardPile := []string{}

	for _, card := range *stack.cards {
		cardStr := ""
		if card.Visible {
			cardStr = fmt.Sprintf("%s %s",
				card.Rank.String(),
				card.Suit.String(),
			)
		} else {
			cardStr = blankCard
		}

		cardPile = append(cardPile, cardStr)
	}

	return cardPile
}

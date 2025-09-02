package state

import (
	"fmt"
	"math/rand/v2"
)

// Stack - A stack of cards one on top of the other and squared such that only
// the topmost card, whether face up or face down is visible.[5]
type Stack struct {
	cards           *[]SuitedCard
	Base            SuitedCard
	Rule            func(SuitedCard) bool
	Type            StackType
	Received        int // count how many times this stack has received cards.
	MaxRedeals      int
	TableauPosition int // Simple count to determine where in the Tableau this is.
	SkipCards       map[SuitedCard]struct{}
}

// StackType represents the type of game component.
type StackType int

const (
	// StackUndefined - Undefined Stock type.
	StackUndefined StackType = iota
	// StackFoundation - Foundations.
	StackFoundation
	// StackTableau - Tableau.
	StackTableau
	// StackTalon - Talon, or Stock.
	StackTalon
	// StackWaste - Waste.
	StackWaste
	// StackReserve - Reserve.
	StackReserve
)

// There are two types of stack, the one that holds the reserve of cards that
// are yet to be played, and the ones that are on the tableau.

// NewStack - Create a new stack with an empty slice of SuitedCards that has a
// capacity of n.
func NewStack(number int, base SuitedCard, rule func(SuitedCard) bool, componentType StackType) *Stack {
	cards := make([]SuitedCard, 0, number)

	return &Stack{
		Base:  base,
		cards: &cards,
		Rule:  rule,
		Type:  componentType,
	}
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

const blankCard = "--"

// Cards - a string representation of the cards in the stack.
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

// CanReceiveMore - returns whether the stack is allowed to receive another
// bunch. Used to check if there has been more redeals than the game specifies
// are allowed.
func (stack *Stack) CanReceiveMore() bool {
	if stack.Received < stack.MaxRedeals || stack.MaxRedeals == -1 {
		stack.Received++

		return true
	}

	return false
}

func (stack *Stack) Shuffle() {
	rand.Shuffle(stack.Len(), func(i, j int) {
		(*stack.cards)[i], (*stack.cards)[j] = (*stack.cards)[j], (*stack.cards)[i]
	})
}

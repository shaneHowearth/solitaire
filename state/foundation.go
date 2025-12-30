package state

import "fmt"

// BaseCard - the card that the Foundation starts at.
type BaseCard Rank

// Foundation - The final place for cards. Cards are built up in piles from the.
// (nominated) base card.
// There are at least 4 Foundations per game, sometimes more (there will always.
// be a multiple of 4 foundations though).
type Foundation struct {
	Stack *Stack
	Base  SuitedCard
}

// CreateFoundations - Create the foundations that will host the cards.
func CreateFoundations(foundationSpec []StackSpec) []*Foundation {
	if len(foundationSpec)%SuitCount != 0 {
		panic("Number of foundations must be a multiple of the number of suits in a deck")
	}

	foundations := make([]*Foundation, 0, len(foundationSpec))

	for idx := 0; idx < len(foundationSpec); idx++ {
		if foundationSpec[idx].AddRule == nil {
			panic(fmt.Sprintf("Cannot create foundation %d without a rule.", idx))
		}

		foundation := &Foundation{
			Base: foundationSpec[idx].BaseCard,
		}

		stack := NewStack(RankCount,
			foundationSpec[idx].BaseCard,
			func(foundationStack *Stack) func(SuitedCard) bool {
				return func(card SuitedCard) bool {
					return foundationSpec[idx].AddRule(foundationStack, card)
				}
			},
			StackFoundation,
		)

		stack.FoundationPosition = idx
		foundation.Stack = stack

		foundations = append(foundations, foundation)
	}

	return foundations
}

// Len - the length of the stack inside the foundation.
func (foundation Foundation) Len() int {
	return foundation.Stack.Len()
}

// Top - the top most card on the stack inside the foundation.
func (foundation Foundation) Top() (SuitedCard, error) {
	return foundation.Stack.Top()
}

func (foundation Foundation) Clone() *Foundation {
	return &Foundation{
		Stack: foundation.Stack.Clone(),
		Base:  foundation.Base,
	}
}

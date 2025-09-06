package state

import "fmt"

// BaseCard - the card that the Foundation starts at.
type BaseCard Rank

// Foundation - The final place for cards. Cards are built up in piles from the
// (nominated) base card.
// There are at least 4 Foundations per game, sometimes more (there will always
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

	for i := 0; i < len(foundationSpec); i++ {
		if foundationSpec[i].AddRule == nil {
			panic(fmt.Sprintf("Cannot create foundation %d without a rule.", i))
		}

		foundation := &Foundation{
			Base: foundationSpec[i].BaseCard,
		}

		stack := NewStack(RankCount,
			foundationSpec[i].BaseCard,
			func(foundationStack *Stack) func(SuitedCard) bool {
				return func(card SuitedCard) bool {
					return foundationSpec[i].AddRule(foundationStack, card)
				}
			},
			StackFoundation,
		)

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

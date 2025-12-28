package state

import "fmt"

type Reserve struct {
	Stack *Stack
}

// CreateReserves - Create the reserve that will host the cards.
func CreateReserves(reserveSpec []StackSpec) []*Reserve {
	if len(reserveSpec) < 1 {
		return []*Reserve{}
	}

	reserves := make([]*Reserve, 0, len(reserveSpec))

	for idx := 0; idx < len(reserveSpec); idx++ {
		if reserveSpec[idx].AddRule == nil {
			panic(fmt.Sprintf("Cannot create reserve %d without a rule.", idx))
		}

		reserve := &Reserve{}

		stack := NewStack(RankCount,
			reserveSpec[idx].BaseCard,
			func(targetStack *Stack) func(SuitedCard) bool {
				return func(card SuitedCard) bool {
					return reserveSpec[idx].AddRule(targetStack, card)
				}
			},
			StackReserve,
		)

		reserve.Stack = stack

		reserves = append(reserves, reserve)
	}

	return reserves
}

// Len - the length of the stack inside the foundation.
func (reserve *Reserve) Len() int {
	return reserve.Stack.Len()
}

// Top - the top most card on the stack inside the foundation.
func (reserve *Reserve) Top() (SuitedCard, error) {
	return reserve.Stack.Top()
}

func (reserve *Reserve) Clone() *Reserve {
	return &Reserve{
		Stack: reserve.Stack.Clone(),
	}
}

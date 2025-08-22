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

	for i := 0; i < len(reserveSpec); i++ {
		if reserveSpec[i].AddRule == nil {
			panic(fmt.Sprintf("Cannot create reserve %d without a rule.", i))
		}

		reserve := &Reserve{}

		stack := NewStack(RankCount,
			reserveSpec[i].BaseCard,
			func(card SuitedCard) bool {
				return reserveSpec[i].AddRule(reserve.Stack, card)
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

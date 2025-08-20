package state

type Reserve struct {
	Stack *Stack
}

// CreateReserves - Create the reserve that will host the cards.
func CreateReserves(number int, rule func(*Reserve, SuitedCard) bool) []*Reserve {
	if number < 1 {
		return []*Reserve{}
	}

	if rule == nil {
		panic("Cannot create depots without a rule.")
	}

	reserves := make([]*Reserve, 0, number)

	for i := 0; i < number; i++ {
		reserve := &Reserve{}

		stack := NewStack(RankCount,
			func(card SuitedCard) bool {
				return rule(reserve, card)
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

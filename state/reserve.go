package state

type Reserve struct {
	Stack *Stack
}

// CreateDepot - Create the reserve that will host the cards.
func CreateDepot(number int, rule func(Reserve, SuitedCard) bool) []*Reserve {
	if number < 1 {
		panic("Cannot have zero iwdepots")
	}

	if rule == nil {
		panic("Cannot create depots without a rule.")
	}

	depots := make([]*Reserve, 0, number)

	for i := 0; i < number; i++ {
		reserve := &Reserve{}

		stack := NewStack(RankCount,
			func(card SuitedCard) bool {
				return rule(*reserve, card)
			},
			StackReserve,
		)

		reserve.Stack = stack

		depots = append(depots, reserve)
	}

	return depots
}

// Len - the length of the stack inside the foundation.
func (reserve Reserve) Len() int {
	return reserve.Stack.Len()
}

// Top - the top most card on the stack inside the foundation.
func (reserve Reserve) Top() (SuitedCard, error) {
	return reserve.Stack.Top()
}

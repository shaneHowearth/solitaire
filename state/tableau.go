package state

// Tableau - An arrangement of cards on the table, typically comprising several
// depots i.e. places where columns of overlapping cards may be formed.
type Tableau struct {
	Stack *Stack
	Base  Rank
}

// CreateTableaus - Create the tableaus that will host the cards.
func CreateTableaus(number int, base Rank, rule func(*Tableau, SuitedCard) bool) []*Tableau {
	if number < 1 {
		panic("Cannot have zero tableaus")
	}

	if rule == nil {
		panic("Cannot create tableaus without a rule.")
	}

	tableaus := make([]*Tableau, 0, number)

	for i := 0; i < number; i++ {
		tableau := Tableau{}
		stack := NewStack(RankCount,
			func(card SuitedCard) bool {
				return rule(&tableau, card)
			},
			StackTableau,
		)
		tableau.Stack = stack
		tableau.Base = base

		tableaus = append(tableaus, &tableau)
	}

	return tableaus
}

// // Add - Add a card to the tableau.
// func (tableau *Tableau) Add(card SuitedCard, visible bool) bool {
// 	if tableau.Stack.Rule(card) {
// 		tableau.Stack.Add(card, visible)
// 		return true
// 	}

// 	return false
// }

// Len - the length of the stack inside the tableau.
func (tableau *Tableau) Len() int {
	return tableau.Stack.Len()
}

// Top - the top most card on the stack inside the tableau.
func (tableau *Tableau) Top() (SuitedCard, error) {
	return tableau.Stack.Top()
}

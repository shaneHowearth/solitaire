package state

import "fmt"

// Tableau - An arrangement of cards on the table, typically comprising several
// depots i.e. places where columns of overlapping cards may be formed.
type Tableau struct {
	Stack *Stack
	Base  Rank
}

// CreateTableaus - Create the tableaus that will host the cards.
func CreateTableaus(tableauSpec []StackSpec) []*Tableau {
	if len(tableauSpec) < 1 {
		panic("Cannot have zero tableaus")
	}

	tableaus := make([]*Tableau, 0, len(tableauSpec))

	for i := 0; i < len(tableauSpec); i++ {
		if tableauSpec[i].AddRule == nil {
			panic(fmt.Sprintf("Cannot create tableau %d without a rule.", i))
		}

		tableau := Tableau{
			Base:  tableauSpec[i].BaseCard.Rank,
			Count: i,
		}

		stack := NewStack(RankCount,
			tableauSpec[i].BaseCard,
			func(card SuitedCard) bool {
				return tableauSpec[i].AddRule(tableau.Stack, card)
			},
			StackTableau,
		)
		tableau.Stack = stack

		tableaus = append(tableaus, &tableau)
	}

	return tableaus
}

// Len - the length of the stack inside the tableau.
func (tableau *Tableau) Len() int {
	return tableau.Stack.Len()
}

// Top - the top most card on the stack inside the tableau.
func (tableau *Tableau) Top() (SuitedCard, error) {
	return tableau.Stack.Top()
}

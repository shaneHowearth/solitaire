package solitaire

import (
	"errors"
	"log/slog"
)

// Tableau - An arrangement of cards on the table, typically comprising several
// depots i.e. places where columns of overlapping cards may be formed.
type Tableau struct {
	Stack *Stack
	Base  Rank
}

// CreateTableaus - Create the tableaus that will host the cards.
func CreateTableaus(number int, rule func(*Tableau, SuitedCard) bool) []*Tableau {
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
		)
		tableau.Stack = stack

		tableaus = append(tableaus, &tableau)
	}

	return tableaus
}

// Empty - the tableau is empty.
func (tableau *Tableau) Empty() bool {
	return tableau.Len() == 0
}

// Add - Add a card to the tableau.
func (tableau *Tableau) Add(card SuitedCard, visible bool) bool {
	if tableau.Stack.Rule(card) {
		tableau.Stack.Add(card, visible)
		return true
	}

	return false
}

// Len - the length of the stack inside the tableau.
func (tableau *Tableau) Len() int {
	return tableau.Stack.Len()
}

// Top - the top most card on the stack inside the tableau.
func (tableau *Tableau) Top() (SuitedCard, error) {
	return tableau.Stack.Top()
}

// MultiMove - ONLY to be used when moving cards from one tableau stack to
// another tableau stack.
// This will move multiple cards at once.
func (tableau *Tableau) MultiMove(destination *Tableau) bool {
	// 1. Find where in the donor stack can move from.
	movable := 0
	canMove := false

	for movable := tableau.Stack.Len() - 1; movable >= 0; movable-- {
		card := (*tableau.Stack.cards)[movable]

		// Only move visible cards.
		if !card.Visible {
			break
		}

		// Stop when the card that can be moved to the destination is found.
		if destination.Stack.Rule(card) {
			canMove = true
			break
		}
	}

	// 2. Move those cards to the tmp Stack that will allow any cards to be
	// added.
	if canMove {
		temporaryStack := NewStack(
			(movable - tableau.Stack.Len()),
			func(SuitedCard) bool { return true },
		)

		for donorIdx := tableau.Stack.Len() - 1; donorIdx >= movable; donorIdx-- {
			card, err := tableau.Top()
			if err != nil {
				// shouldn't be able to get here.
				return false
			}

			temporaryStack.Add(card, true)
		}

		// Move the cards from the temporary stack onto the destination.
		for {
			card, err := temporaryStack.Deal()
			if errors.Is(err, ErrEmpty) {
				break
			}

			if !destination.Add(card, true) {
				// Should not happen.
				slog.Warn("")
				return false
			}
		}

		return true
	}

	// 3. Move the cards from the temp stack to the destination.
	return false
}

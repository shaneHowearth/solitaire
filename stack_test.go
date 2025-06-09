package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	shuffledDeck := solitaire.CreateDecks(1)
	shuffledDeck.Shuffle()

	testcases := map[string]struct {
		Number  int
		Visible bool
	}{
		"Add one":         {Number: 1},
		"Add ten":         {Number: 10},
		"Add ten visible": {Number: 10, Visible: true},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			stack := solitaire.NewStack(testcase.Number, func(solitaire.SuitedCard) bool { return true })

			for x := 1; x <= testcase.Number; x++ {
				card := shuffledDeck.Deal()
				card.Visible = testcase.Visible
				stack.Add(card, testcase.Visible)
				assert.Equalf(t, x, stack.Len(), "Stack length error got %d want %d", stack.Len(), x)
				topCard, err := stack.Top()
				assert.Nilf(t, err, "Getting the top card errored %v", err)

				assert.EqualExportedValuesf(t,
					topCard, card,
					"Got different cards got %#v want %#v",
					topCard, card)
			}
		})
	}
}

func Test_Move(t *testing.T) {
	source := solitaire.NewStack(4, func(solitaire.SuitedCard) bool { return true })
	standardDeck := solitaire.CreateDecks(1)

	for idx := 0; idx < 10; idx++ {
		source.Add(standardDeck.Deal(), true)
	}

	testcases := map[string]struct {
		number           int
		sourceStack      *solitaire.Stack
		destinationStack *solitaire.Stack
		output           bool
	}{
		// Cannot move card to the same stack as it came from.
		"Should not be able to move a card to the same stack that it came from": {
			number:           1,
			sourceStack:      source,
			destinationStack: source,
			output:           false,
		},
		"Move to an empty stack where the rule allows the move": {
			number:           1,
			sourceStack:      source,
			destinationStack: solitaire.NewStack(0, func(solitaire.SuitedCard) bool { return true }),
			output:           true,
		},
		"Move to an empty stack where the rule denies the move": {
			number:           1,
			sourceStack:      source,
			destinationStack: solitaire.NewStack(0, func(solitaire.SuitedCard) bool { return false }),
			output:           false,
		},
		// Waste to Empty Foundation.
		// Waste to partially filled Foundation.
		// Waste to nominated Tableau.
		// Waste to empty Tableau.
		// Tableau to Foundation.
		// Foundation to nominated tableau.
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			output := testcase.sourceStack.Move(testcase.number, testcase.destinationStack)
			assert.Equalf(t, testcase.output, output, "got %t want %t", output, testcase.output)
		})
	}
}

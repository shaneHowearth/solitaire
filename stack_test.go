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
	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			stack := solitaire.NewStack(testCase.Number, func(solitaire.SuitedCard) bool { return true })

			for x := 1; x <= testCase.Number; x++ {
				card := shuffledDeck.Deal()
				card.Visible = testCase.Visible
				stack.Add(card, testCase.Visible)
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

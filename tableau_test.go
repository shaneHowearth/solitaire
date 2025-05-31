package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTableaus(t *testing.T) {
	testcases := map[string]struct {
		WillPanic    bool
		PanicMessage string
		Number       int
		Rule         func(solitaire.Tableau, solitaire.SuitedCard) bool
	}{
		"Zero tableaus": {
			WillPanic:    true,
			PanicMessage: "Cannot have zero tableaus",
			Number:       0,
			Rule: func(
				solitaire.Tableau,
				solitaire.SuitedCard,
			) bool {
				// Allow everything to be added.
				return true
			},
		},
		"No rule": {
			WillPanic:    true,
			PanicMessage: "Cannot create tableaus without a rule.",
			Number:       1,
		},
		"Seven tableaus (klondike)": {
			Number: 7,
			Rule: func(
				solitaire.Tableau,
				solitaire.SuitedCard,
			) bool {
				// Allow everything to be added.
				return true
			},
		},
	}
	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			if testCase.WillPanic {
				assert.PanicsWithValue(t, testCase.PanicMessage,
					func() {
						solitaire.CreateTableaus(
							testCase.Number,
							testCase.Rule,
						)
					},
				)
			}

			if !testCase.WillPanic {
				tableau := solitaire.CreateTableaus(
					testCase.Number,
					testCase.Rule,
				)

				// tableau has the correct number of elements.
				assert.Equalf(t, testCase.Number, len(tableau),
					"tableau has incorrect number of elements, want: %d, got: %d", testCase.Number, len(tableau))
			}
		})
	}
}

func Test_Empty(t *testing.T) {
	testcases := map[string]struct {
		Number   int
		Expected bool
	}{
		"Not empty": {Number: 6, Expected: false},
		"Empty":     {Number: 0, Expected: true},
	}

	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			tableaus := solitaire.CreateTableaus(
				2,
				func(
					solitaire.Tableau,
					solitaire.SuitedCard,
				) bool {
					// Allow everything to be added.
					return true
				},
			)

			standardDeck := solitaire.CreateDecks(1)

			for idx := 0; idx < testCase.Number; idx++ {
				card := standardDeck.Deal()
				tableaus[0].Add(card, true)
			}

			assert.Equalf(
				t,
				testCase.Expected,
				tableaus[0].Empty(),
				"tableaus emptiness got %t want %t",
				tableaus[0].Empty(),
				testCase.Expected,
			)
		})
	}
}

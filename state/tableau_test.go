package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTableaus(t *testing.T) {
	testcases := map[string]struct {
		WillPanic    bool
		PanicMessage string
		Number       int
		Rule         func(*state.Tableau, state.SuitedCard) bool
	}{
		"Zero tableaus": {
			WillPanic:    true,
			PanicMessage: "Cannot have zero tableaus",
			Number:       0,
			Rule: func(
				*state.Tableau,
				state.SuitedCard,
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
				*state.Tableau,
				state.SuitedCard,
			) bool {
				// Allow everything to be added.
				return true
			},
		},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			if testcase.WillPanic {
				assert.PanicsWithValue(t, testcase.PanicMessage,
					func() {
						state.CreateTableaus(
							testcase.Number,
							testcase.Rule,
						)
					},
				)
			}

			if !testcase.WillPanic {
				tableau := state.CreateTableaus(
					testcase.Number,
					testcase.Rule,
				)

				// tableau has the correct number of elements.
				assert.Equalf(t, testcase.Number, len(tableau),
					"tableau has incorrect number of elements, want: %d, got: %d", testcase.Number, len(tableau))
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

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			tableaus := state.CreateTableaus(
				2,
				func(
					*state.Tableau,
					state.SuitedCard,
				) bool {
					// Allow everything to be added.
					return true
				},
			)

			standardDeck := state.CreateDecks(1)

			for idx := 0; idx < testcase.Number; idx++ {
				card := standardDeck.Deal()
				tableaus[0].Add(card, true)
			}

			assert.Equalf(
				t,
				testcase.Expected,
				tableaus[0].Empty(),
				"tableaus emptiness got %t want %t",
				tableaus[0].Empty(),
				testcase.Expected,
			)
		})
	}
}

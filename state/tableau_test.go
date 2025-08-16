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
		Rank         state.Rank
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
							testcase.Rank,
							testcase.Rule,
						)
					},
				)
			}

			if !testcase.WillPanic {
				tableau := state.CreateTableaus(
					testcase.Number,
					testcase.Rank,
					testcase.Rule,
				)

				// tableau has the correct number of elements.
				assert.Equalf(t, testcase.Number, len(tableau),
					"tableau has incorrect number of elements, want: %d, got: %d", testcase.Number, len(tableau))
			}
		})
	}
}

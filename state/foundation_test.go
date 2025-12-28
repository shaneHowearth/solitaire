package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func Test_CreateFoundations(t *testing.T) {
	BaseRank := state.Ace

	testcases := map[string]struct {
		Number       int
		WillPanic    bool
		PanicMessage string
	}{
		"GIVEN that the number of foundations requested is Zero " +
			"WHEN Create Foundations is called " +
			"THEN the code will panic.": {
			Number: 0, WillPanic: true,
			PanicMessage: "Cannot have zero foundations",
		},
		"GIVEN that the number of foundation requested is Two " +
			"WHEN Create foundation is called " +
			"THEN the code will panic.": {
			Number: 2, WillPanic: true,
			PanicMessage: "Number of foundations must be a multiple of the number of suits in a deck",
		},
		"GIVEN that the number of foundation requested is Four " +
			"WHEN Create foundation is called " +
			"THEN Four foundation will be created.": {Number: 4},
		"GIVEN that the number of foundation requested is Four " +
			"WHEN Create foundation is called " +
			"THEN Four foundations will be created.": {Number: 8},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			// 1. Dynamically build the specs slice based on testcase.Number
			specs := make([]state.StackSpec, testcase.Number)
			for i := 0; i < testcase.Number; i++ {
				specs[i] = state.StackSpec{
					BaseCard: state.SuitedCard{Rank: BaseRank, Suit: state.Suit(i % 4)},
					AddRule:  game.MinusOneRule,
				}
			}

			if testcase.WillPanic {
				assert.PanicsWithValue(t, testcase.PanicMessage, func() {
					state.CreateFoundations(specs)
				})
			} else {
				foundation := state.CreateFoundations(specs)
				assert.Equal(t, testcase.Number, len(foundation),
					"foundation has incorrect number of elements")
			}
		})
	}
}

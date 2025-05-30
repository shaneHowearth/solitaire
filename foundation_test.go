package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
	"github.com/stretchr/testify/assert"
)

func Test_CreateFoundationss(t *testing.T) {
	BaseCard := solitaire.Ace

	testcases := map[string]struct {
		Number       int
		WillPanic    bool
		PanicMessage string
	}{
		"GIVEN that the number of foundations requested is Zero " +
			"WHEN Create Foundations is called " +
			"THEN the code will panic.": {
			Number: 0, WillPanic: true,
			PanicMessage: "Cannot have zero foundations"},
		"GIVEN that the number of foundation requested is Two " +
			"WHEN Create foundation is called " +
			"THEN the code will panic.": {
			Number: 2, WillPanic: true,
			PanicMessage: "Number of foundations must be a multiple of the number of suits in a deck"},
		"GIVEN that the number of foundation requested is Four " +
			"WHEN Create foundation is called " +
			"THEN Four foundation will be created.": {Number: 4},
		"GIVEN that the number of foundation requested is Four " +
			"WHEN Create foundation is called " +
			"THEN Four foundations will be created.": {Number: 8},
	}
	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			if testCase.WillPanic {
				assert.PanicsWithValue(t, testCase.PanicMessage,
					func() {
						solitaire.CreateFoundations(
							testCase.Number,
							BaseCard,
							func(
								solitaire.Foundation,
								solitaire.SuitedCard,
							) bool {
								return false
							},
						)
					},
				)
			}

			if !testCase.WillPanic {
				foundation := solitaire.CreateFoundations(
					testCase.Number,
					BaseCard,
					func(
						solitaire.Foundation,
						solitaire.SuitedCard,
					) bool {
						return false
					},
				)

				// foundation has the correct number of elements.
				assert.Equalf(t, testCase.Number, len(foundation),
					"foundation has incorrect number of elements, want: %d, got: %d", testCase.Number, len(foundation))
			}
		})
	}
}

func Test_Full(t *testing.T) {
	testcases := map[string]struct {
		Count  int
		IsFull bool
	}{

		"One":  {Count: 1},
		"Full": {Count: solitaire.CardCount, IsFull: true},
	}
	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			standardDeck := solitaire.CreateDecks(1)
			card := standardDeck.Deal()

			foundation := solitaire.CreateFoundations(
				solitaire.SuitCount,
				solitaire.Ace,
				func(
					solitaire.Foundation,
					solitaire.SuitedCard,
				) bool {
					return true
				},
			)

			for x := 0; x < testCase.Count; x++ {
				foundation[0].Add(card, true)
			}

			assert.Equalf(
				t, testCase.IsFull,
				foundation[0].Full(),
				"Foundation full error got %t want %t",
				foundation[0].Full(),
				testCase.IsFull,
			)
		})
	}
}

package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
	"github.com/stretchr/testify/assert"
)

func Test_CreateDecks(t *testing.T) {
	standardDeck := solitaire.CreateDecks(1)

	testcases := map[string]struct {
		Number       int
		WillPanic    bool
		PanicMessage string
	}{
		"GIVEN that the number of decks requested is Zero " +
			"WHEN Create Decks is called " +
			"THEN the code will panic.": {Number: 0, WillPanic: true, PanicMessage: "cannot create zero decks"},
		"GIVEN that the number of decks requested is One " +
			"WHEN Create Decks is called " +
			"THEN One deck will be created.": {Number: 1},
		"GIVEN that the number of decks requested is Two " +
			"WHEN Create Decks is called " +
			"THEN Two decks will be created.": {Number: 2},
		"GIVEN that the number of decks requested is Ten " +
			"WHEN Create Decks is called " +
			"THEN Ten decks will be created.": {Number: 10},
	}
	for name, testCase := range testcases {
		t.Run(name, func(t *testing.T) {
			deck := solitaire.Deck{}

			if testCase.WillPanic {
				assert.PanicsWithValue(t, testCase.PanicMessage, func() { solitaire.CreateDecks(testCase.Number) })
			}

			if !testCase.WillPanic {
				deck = solitaire.CreateDecks(testCase.Number)
			}

			expected := testCase.Number * 52

			// Deck has the correct number of elements.
			assert.Equalf(t, expected, len(deck),
				"Deck has incorrect number of elements, want: %d, got: %d", expected, len(deck))

			// Deck has the right cards in each subdeck.
			for i := 1; i < testCase.Number; i++ {
				start := (i - 1) * 52
				stop := i * 52
				assert.ElementsMatchf(t, standardDeck, deck[start:stop],
					"Deck [%d:%d] has incorrect elements %#v", deck[start:stop])
			}
		})
	}
}

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
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			var deck *solitaire.Deck

			if testcase.WillPanic {
				assert.PanicsWithValue(t, testcase.PanicMessage, func() { solitaire.CreateDecks(testcase.Number) })
			}

			if !testcase.WillPanic {
				deck = solitaire.CreateDecks(testcase.Number)
			}

			expected := testcase.Number * 52

			// Deck has the correct number of elements.
			assert.Equalf(t, expected, deck.Len(),
				"Deck has incorrect number of elements, want: %d, got: %d", expected, deck.Len())

			// Deck has the right cards in each subdeck.
			for i := 1; i < testcase.Number; i++ {
				start := (i - 1) * 52
				stop := i * 52
				assert.ElementsMatchf(t, *standardDeck, (*deck)[start:stop],
					"Deck [%d:%d] has incorrect elements %#v", (*deck)[start:stop])
			}
		})
	}
}

func Test_ShuffleDecks(t *testing.T) {
	testcases := map[string]struct {
		Number int
	}{
		"one": {Number: 1},
		"two": {Number: 2},
		"ten": {Number: 10},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			standardDeck := solitaire.CreateDecks(testcase.Number)

			shuffledDeck := solitaire.CreateDecks(testcase.Number)
			expected := testcase.Number * 52

			shuffledDeck.Shuffle()

			// Deck has the correct number of elements.
			assert.Equalf(t, expected, shuffledDeck.Len(),
				"Deck has incorrect number of elements after shuffle, want: %d, got: %d", expected, shuffledDeck.Len())

			// Look for a value out of position.
			// Only one value needs to be out of position for the deck to be
			// considered shuffled for the purposes of the test.
			// TODO: Is there a stronger test available?
			shuffled := false

			for i := 0; i < testcase.Number*52; i++ {
				if (*shuffledDeck)[i] != (*standardDeck)[i] {
					shuffled = true
					break
				}
			}

			assert.True(t, shuffled, "Shuffled deck wasn't shuffled")
		})
	}
}

func Test_DeckDeal(t *testing.T) {
	testcases := map[string]struct {
		Number int
	}{
		"one": {Number: 1},
		"two": {Number: 2},
		"ten": {Number: 10},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			shuffledDeck := solitaire.CreateDecks(testcase.Number)

			shuffledDeck.Shuffle()

			for idx := 1; idx <= testcase.Number*solitaire.RankCount*solitaire.SuitCount; idx++ {
				expected := shuffledDeck.Top()

				assert.Equal(t, solitaire.RankCount*solitaire.SuitCount*testcase.Number-(idx-1), shuffledDeck.Len())

				actual := shuffledDeck.Deal()
				assert.EqualExportedValues(t, expected, actual, "Wrong card returned got %v want %v", actual, expected)

				// Check that the length of the deck reduces after each Deal.
				assert.Equal(t, solitaire.RankCount*solitaire.SuitCount*testcase.Number-(idx), shuffledDeck.Len())
			}

			// Deck should be empty.
			expected := solitaire.SuitedCard{}
			actual := shuffledDeck.Deal()
			assert.EqualExportedValues(t, expected, actual, "Wrong card returned got %v want %v", actual, expected)

			assert.Equal(t, 0, shuffledDeck.Len())
		})
	}
}

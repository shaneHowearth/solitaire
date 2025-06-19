package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
)

func Test_CardToString(t *testing.T) {
	t.Parallel()
	// Ace.
	expectedAceString := "Ace"
	if state.Ace.String() != expectedAceString {
		t.Errorf(
			"unexpected result Ace.String got %q, want %q",
			state.Ace.String(),
			expectedAceString,
		)
	}
	// Two.
	expectedTwoString := "2"
	if state.Two.String() != expectedTwoString {
		t.Errorf(
			"unexpected result Two.String got %q, want %q",
			state.Two.String(),
			expectedTwoString,
		)
	}
	// Three.
	expectedThreeString := "3"
	if state.Three.String() != expectedThreeString {
		t.Errorf(
			"unexpected result Three.String got %q, want %q",
			state.Three.String(),
			expectedThreeString,
		)
	}
	// Four.
	expectedFourString := "4"
	if state.Four.String() != expectedFourString {
		t.Errorf(
			"unexpected result Four.String got %q, want %q",
			state.Four.String(),
			expectedFourString,
		)
	}
	// Five.
	expectedFiveString := "5"
	if state.Five.String() != expectedFiveString {
		t.Errorf(
			"unexpected result Five.String got %q, want %q",
			state.Five.String(),
			expectedFiveString,
		)
	}
	// Six.
	expectedSixString := "6"
	if state.Six.String() != expectedSixString {
		t.Errorf(
			"unexpected result Six.String got %q, want %q",
			state.Six.String(),
			expectedSixString,
		)
	}
	// Seven.
	expectedSevenString := "7"
	if state.Seven.String() != expectedSevenString {
		t.Errorf(
			"unexpected result Seven.String got %q, want %q",
			state.Seven.String(),
			expectedSevenString,
		)
	}
	// Eight.
	expectedEightString := "8"
	if state.Eight.String() != expectedEightString {
		t.Errorf(
			"unexpected result Eight.String got %q, want %q",
			state.Eight.String(),
			expectedEightString,
		)
	}
	// Nine.
	expectedNineString := "9"
	if state.Nine.String() != expectedNineString {
		t.Errorf(
			"unexpected result Nine.String got %q, want %q",
			state.Nine.String(),
			expectedNineString,
		)
	}
	// Ten.
	expectedTenString := "10"
	if state.Ten.String() != expectedTenString {
		t.Errorf(
			"unexpected result Ten.String got %q, want %q",
			state.Ten.String(),
			expectedTenString,
		)
	}
	// Jack.
	expectedJackString := "Jack"
	if state.Jack.String() != expectedJackString {
		t.Errorf(
			"unexpected result Jack.String got %q, want %q",
			state.Jack.String(),
			expectedJackString,
		)
	}
	// Queen.
	expectedQueenString := "Queen"
	if state.Queen.String() != expectedQueenString {
		t.Errorf(
			"unexpected result Queen.String got %q, want %q",
			state.Queen.String(),
			expectedQueenString,
		)
	}
	// King.
	expectedKingString := "King"
	if state.King.String() != expectedKingString {
		t.Errorf(
			"unexpected result King.String got %q, want %q",
			state.King.String(),
			expectedKingString,
		)
	}
}

func Test_SuitToString(t *testing.T) {
	t.Parallel()
	// Spades.
	expectedSpadeString := "♠"
	if state.Spades.String() != expectedSpadeString {
		t.Errorf(
			"unexpected result Spades.String got %q, want %q",
			state.Spades.String(),
			expectedSpadeString,
		)
	}
	// Hearts.
	expectedHeartString := "♥"
	if state.Hearts.String() != expectedHeartString {
		t.Errorf(
			"unexpected result Hearts.String got %q, want %q",
			state.Hearts.String(),
			expectedHeartString,
		)
	}
	// Clubs.
	expectedClubString := "♣"
	if state.Clubs.String() != expectedClubString {
		t.Errorf(
			"unexpected result Clubs.String got %q, want %q",
			state.Clubs.String(),
			expectedClubString,
		)
	}
	// Diamonds.
	expectedDiamondString := "♦"
	if state.Diamonds.String() != expectedDiamondString {
		t.Errorf(
			"unexpected result Diamonds.String got %q, want %q",
			state.Diamonds.String(),
			expectedDiamondString,
		)
	}
}

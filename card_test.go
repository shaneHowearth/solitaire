package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
)

func Test_CardToString(t *testing.T) {
	t.Parallel()
	// Ace.
	expectedAceString := "A"
	if solitaire.Ace.String() != expectedAceString {
		t.Errorf(
			"unexpected result Ace.String got %q, want %q",
			solitaire.Ace.String(),
			expectedAceString,
		)
	}
	// Two.
	expectedTwoString := "2"
	if solitaire.Two.String() != expectedTwoString {
		t.Errorf(
			"unexpected result Two.String got %q, want %q",
			solitaire.Two.String(),
			expectedTwoString,
		)
	}
	// Three.
	expectedThreeString := "3"
	if solitaire.Three.String() != expectedThreeString {
		t.Errorf(
			"unexpected result Three.String got %q, want %q",
			solitaire.Three.String(),
			expectedThreeString,
		)
	}
	// Four.
	expectedFourString := "4"
	if solitaire.Four.String() != expectedFourString {
		t.Errorf(
			"unexpected result Four.String got %q, want %q",
			solitaire.Four.String(),
			expectedFourString,
		)
	}
	// Five.
	expectedFiveString := "5"
	if solitaire.Five.String() != expectedFiveString {
		t.Errorf(
			"unexpected result Five.String got %q, want %q",
			solitaire.Five.String(),
			expectedFiveString,
		)
	}
	// Six.
	expectedSixString := "6"
	if solitaire.Six.String() != expectedSixString {
		t.Errorf(
			"unexpected result Six.String got %q, want %q",
			solitaire.Six.String(),
			expectedSixString,
		)
	}
	// Seven.
	expectedSevenString := "7"
	if solitaire.Seven.String() != expectedSevenString {
		t.Errorf(
			"unexpected result Seven.String got %q, want %q",
			solitaire.Seven.String(),
			expectedSevenString,
		)
	}
	// Eight.
	expectedEightString := "8"
	if solitaire.Eight.String() != expectedEightString {
		t.Errorf(
			"unexpected result Eight.String got %q, want %q",
			solitaire.Eight.String(),
			expectedEightString,
		)
	}
	// Nine.
	expectedNineString := "9"
	if solitaire.Nine.String() != expectedNineString {
		t.Errorf(
			"unexpected result Nine.String got %q, want %q",
			solitaire.Nine.String(),
			expectedNineString,
		)
	}
	// Ten.
	expectedTenString := "10"
	if solitaire.Ten.String() != expectedTenString {
		t.Errorf(
			"unexpected result Ten.String got %q, want %q",
			solitaire.Ten.String(),
			expectedTenString,
		)
	}
	// Jack.
	expectedJackString := "J"
	if solitaire.Jack.String() != expectedJackString {
		t.Errorf(
			"unexpected result Jack.String got %q, want %q",
			solitaire.Jack.String(),
			expectedJackString,
		)
	}
	// Queen.
	expectedQueenString := "Q"
	if solitaire.Queen.String() != expectedQueenString {
		t.Errorf(
			"unexpected result Queen.String got %q, want %q",
			solitaire.Queen.String(),
			expectedQueenString,
		)
	}
	// King.
	expectedKingString := "K"
	if solitaire.King.String() != expectedKingString {
		t.Errorf(
			"unexpected result King.String got %q, want %q",
			solitaire.King.String(),
			expectedKingString,
		)
	}
}

func Test_SuitToString(t *testing.T) {
	t.Parallel()
	// Spade.
	expectedSpadeString := "♠"
	if solitaire.Spade.String() != expectedSpadeString {
		t.Errorf(
			"unexpected result Spade.String got %q, want %q",
			solitaire.Spade.String(),
			expectedSpadeString,
		)
	}
	// Heart.
	expectedHeartString := "♥"
	if solitaire.Heart.String() != expectedHeartString {
		t.Errorf(
			"unexpected result Heart.String got %q, want %q",
			solitaire.Heart.String(),
			expectedHeartString,
		)
	}
	// Club.
	expectedClubString := "♣"
	if solitaire.Club.String() != expectedClubString {
		t.Errorf(
			"unexpected result Club.String got %q, want %q",
			solitaire.Club.String(),
			expectedClubString,
		)
	}
	// Diamond.
	expectedDiamondString := "♦"
	if solitaire.Diamond.String() != expectedDiamondString {
		t.Errorf(
			"unexpected result Diamond.String got %q, want %q",
			solitaire.Diamond.String(),
			expectedDiamondString,
		)
	}
}

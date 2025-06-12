package solitaire

// SuitCount - the number of suits.
const SuitCount = 4

// RankCount - the number of unsuited cards.
const RankCount = 13

// Rank - Enumeration of cards.
type Rank int

//nolint:revive // ignore need to comment this block of exported consts.
const (
	Ace Rank = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

func (rank Rank) String() string {
	return [...]string{
		"A",
		"2",
		"3",
		"4",
		"5",
		"6",
		"7",
		"8",
		"9",
		"10",
		"J",
		"Q",
		"K",
	}[rank]
}

// Suit - Enumeration of suits.
type Suit int

//nolint:revive // ignore need to comment this block of exported consts.
const (
	Hearts Suit = iota
	Clubs
	Diamonds
	Spades
)

func (suit Suit) String() string {
	return [...]string{
		"♥",
		"♣",
		"♦",
		"♠",
	}[suit]
}

// SuitedCard - A type that encapsulates that a hard has a rank/name and a
// suit.
//
//nolint:unused // The fields are used elsewhere.
type SuitedCard struct {
	Rank    Rank
	Suit    Suit
	Visible bool
}

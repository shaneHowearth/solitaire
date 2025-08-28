package game

import "github.com/shanehowearth/solitaire/state"

// Acme - https://en.wikipedia.org/wiki/Acme_(card_game)
type Acme struct{}

// Ensure that Acme implements game.Variant.
var _ Variant = (*Acme)(nil)

// Name - name of the variant.
func (*Acme) Name() string {
	return "Acme"
}

// TableauGridSize - The size of the grid required by acme.
func (*Acme) TableauGridSize() (int, int) {
	const height = 1
	const numAcmeTableau = 4

	return height, numAcmeTableau
}

// Decks - How many decks of cards are required to play acme.
func (*Acme) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*Acme) Reserves() []state.StackSpec {
	return []state.StackSpec{
		{
			AddRule: func(*state.Stack, state.SuitedCard) bool {
				// Nothing can be added to a reserve.
				return false
			},
			CardCount: [2]int{13, 1},
		},
	}
}

// Tableau - how the tableau are defined.
func (acme *Acme) Tableau() []state.StackSpec {

	return []state.StackSpec{
		{
			AddRule:   acme.tableauRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   acme.tableauRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   acme.tableauRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
		{
			AddRule:   acme.tableauRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{Rank: state.King},
		},
	}
}

func (*Acme) tableauRule(tableau *state.Stack, card state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// Get the card currently at the top of the tableau.
	topCard, err := tableau.Top()
	if err != nil {
		return false
	}

	// If the card is the same suit, and is one down in rank then it can go onto
	// the tableau.
	if (card.Suit == topCard.Suit) && (topCard.Rank-card.Rank) == 1 {
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// Foundations - how the foundations are defined.
func (*Acme) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

// HowToPlay - Tell the player how to play the game.
func (*Acme) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The four foundations (rectangles in the upper right of the board) are built up by suit from Ace (low in this game) to King, and the tableau piles can be built down by suit.
`,
		`The aim of the game is to build up four stacks of cards starting with Ace and ending with King all of the same suit, on one of the four foundations, at which time the player would have won.
`,
		`There is only one (1) redeal available in this game.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*Acme) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, foundation := range foundations {
		if foundation.Len() != state.RankCount {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*Acme) MaxRedeals() int {
	// Only one redeal is allowed.
	return 1
}

// Move -
func (*Acme) Move(source, destination *state.Stack, _ []*state.Tableau) bool {
	return Move(source, destination)
}

// Compact
func (*Acme) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Acme) Talon() bool {
	return true
}

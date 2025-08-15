package game

import "github.com/shanehowearth/solitaire/state"

// Acme - https://en.wikipedia.org/wiki/Acme_(card_game)
type Acme struct{}

var _ Variant = (*Acme)(nil)

// Name - name of the variant.
func (*Acme) Name() string {
	return "Acme"
}

// numAcmeTableau - Strictly speaking there are 4 tableau and one depot for
// Acme.
const numAcmeTableau = 5

// TableauGridSize - The size of the grid required by acme.
func (*Acme) TableauGridSize() (int, int) {
	const height = 1

	return height, numAcmeTableau
}

// Decks - How many decks of cards are required to play acme.
func (*Acme) Decks() int {
	return 1
}

// Tableau - how the tableau are defined.
func (acme *Acme) Tableau() (
	number int,
	basecard state.Rank,
	addRule func(*state.Tableau, state.SuitedCard) bool,
) {
	return numAcmeTableau, state.King, acme.tableauRule
}

func (*Acme) tableauRule(tableau *state.Tableau, card state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if (*tableau).Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// Get the card currently at the top of the tableau.
	topCard, err := (*tableau).Top()
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

// TableauPosition - Not currently used.
func (*Acme) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0

	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
func (*Acme) Foundations() (
	number int,
	basecard state.Rank,
	addRule func(state.Foundation, state.SuitedCard) bool,
) {
	const foundationCount = 4
	return foundationCount, state.Ace, PlusOneRule
}

// SetupDealCardCounts - Should return a list of ints, the first int will be the
// number of cards going into the first tableau, the second will be how many
// cards are visible in that tableau. The third and fourth ints will apply
// to the second tableau, etc.
func (*Acme) SetupDealCardCounts() []int {
	//nolint:revive // Ignore the constant complaint.
	return []int{13, 1, 1, 1, 1, 1, 1, 1, 1, 1}
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
	// Allow an unlimited number of redeals.
	return 1
}

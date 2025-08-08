package game

import "github.com/shanehowearth/solitaire/state"

// Variant - The variant of solitaire being defined.
type Variant interface {
	// Name of the Variant.
	Name() string

	// TableauGridSize - How big is the grid that the tableau needs.
	TableauGridSize() (height, width int)

	// Decks - How many decks of cards are required to play the variant.
	Decks() int

	// Tableau.
	Tableau() (
		number int,
		basecard state.Rank,
		addRule func(
			*state.Tableau,
			state.SuitedCard,
		) bool,
	)

	// TableauPosition - the position of the the tableau.
	// The number of each tableau is passed to a function that returns the
	// x, y position on the grid, and the orientation of the pile (in degrees).
	// tableaus are to be 1 indexed.
	TableauPosition(tableauNumber int) (
		x,
		y,
		orientation int,
	)

	// Foundations - how many, what is the first card to go on one if the
	// tableau is empty, and what rule is to be applied when deciding if a new
	// card can be added to the tableau.
	Foundations() (
		number int,
		basecard state.Rank,
		addRule func(
			state.Foundation,
			state.SuitedCard,
		) bool,
	)

	// SetupDealCardCounts - Should return a list of ints, the first int will be the
	// number of cards going into the first tableau, the second will be how many
	// cards are visible in that tableau. The third and fourth ints will apply
	// to the second tableau, etc.
	SetupDealCardCounts() []int
}

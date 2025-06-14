package game

import "github.com/shanehowearth/solitaire"

// Variant - The variant of solitaire being defined.
type Variant interface {
	// Name of the Variant.
	Name() string

	// GridSize - How big is the grid that the tableau needs.
	GridSize() (height, width int)

	// Decks - How many decks of cards are required to play the variant.
	Decks() int

	// Tableau.
	Tableau() (
		number int,
		basecard solitaire.Rank,
		addRule func(
			solitaire.Tableau,
			solitaire.SuitedCard,
		) bool,
	)

	// Layout of tableau.
	// The number of each tableau is passed to a function that returns the
	// x, y position on the grid, and the orientation of the pile (in degrees).
	// tableaus are to be 1 indexed.
	Layout(tableauNumber int) (
		x,
		y,
		orientation int,
	)

	// Foundations - how many, what is the first card to go on one if the
	// tableau is empty, and what rule is to be applied when deciding if a new
	// card can be added to the tableau.
	Foundations() (
		number int,
		basecard solitaire.Rank,
		addRule func(
			solitaire.Foundation,
			solitaire.SuitedCard,
		) bool,
	)
}

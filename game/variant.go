package game

import "github.com/shanehowearth/solitaire/state"

// Variant - The variant of solitaire being defined.
type Variant interface {
	// Name of the Variant.
	// This name is what is displayed to the user.
	Name() string

	// Decks - How many decks of cards are required to play the variant.
	Decks() int

	// Tableau - An arrangement of cards on the table, typically comprising
	// several depots i.e. places where columns of overlapping cards may be
	// formed, the packing taking place on the available cards on the columns.
	// It is thus distinct from a layout, reserve, talon or wastepile.[2] The
	// main part of the layout on the table.[9] Sometimes equated, confusingly,
	// to layout. https://en.wikipedia.org/wiki/Glossary_of_patience_terms#tableau
	Tableau() (
		number int,
		base state.Rank,
		addRule func(
			*state.Tableau,
			state.SuitedCard,
		) bool,
	)

	// TableauGridSize - How big is the grid that the tableau needs.
	TableauGridSize() (height, width int)

	// TableauPosition - the position of the the tableau.
	// The number of each tableau is passed to a function that returns the
	// x, y position on the grid, and the orientation of the pile (in degrees).
	// tableaus are to be 1 indexed.
	TableauPosition(tableauNumber int) (
		x,
		y,
		orientation int,
	)

	// Reserves - this gives how many reserves are required, and their
	// configuration.
	// A reserve is cards available for play that are not part of the
	// foundations, talon, tableau or discard piles.
	// https://en.wikipedia.org/wiki/Glossary_of_patience_terms#reservehttps://en.wikipedia.org/wiki/Glossary_of_patience_terms#reserve
	Reserves() (
		number int,
		cardCount [][2]int, // A description of how many cards each reserve receives at the beginning.
		addRule func(
			*state.Reserve,
			state.SuitedCard,
		) bool,
	)

	// Foundations - this gives how many foundations are required, and their
	// configuration.
	// A foundation is a pile of cards, typically squared and face-up, and built
	// on the bottom card which is the foundation card. As the tableau is
	// cleared, cards are moved to the foundations.
	Foundations() (
		number int,
		basecard state.Rank,
		addRule func(
			state.Foundation,
			state.SuitedCard,
		) bool,
	)

	// SetupTableauCardCounts - Should return a list of [2]int, the first int
	// will be the number of cards going into the first tableau, the second will
	// be how many cards are visible in that tableau. The third and fourth ints
	// will apply to the second tableau, etc.
	SetupTableauCardCounts() [][2]int
	// SetupDealCardCounts() []int

	// HowToPlay - Explains to the player how the game is played.
	HowToPlay() []string

	// HasWon - Checks if the Game has been Won.
	HasWon([]*state.Tableau, []*state.Foundation) bool

	// MaxRedeals - Rule for how many times the stock can be dealt back to the
	// talon/stock -1 inidcates unlimited.
	MaxRedeals() int
}

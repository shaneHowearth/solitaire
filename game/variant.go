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
	Tableau() []state.StackSpec

	// TableauGridSize - How big is the grid that the tableau needs.
	TableauGridSize() (height, width int)

	// Reserves - this gives how many reserves are required, and their
	// configuration.
	// A reserve is cards available for play that are not part of the
	// foundations, talon, tableau or discard piles.
	// https://en.wikipedia.org/wiki/Glossary_of_patience_terms#reservehttps://en.wikipedia.org/wiki/Glossary_of_patience_terms#reserve
	Reserves() []state.StackSpec

	// Foundations - this gives how many foundations are required, and their
	// configuration.
	// A foundation is a pile of cards, typically squared and face-up, and built
	// on the bottom card which is the foundation card. As the tableau is
	// cleared, cards are moved to the foundations.
	Foundations() []state.StackSpec

	// HowToPlay - Explains to the player how the game is played.
	HowToPlay() []string

	// HasWon - Checks if the Game has been Won.
	HasWon([]*state.Tableau, []*state.Foundation) bool

	// MaxRedeals - Rule for how many times the stock can be dealt back to the
	// talon/stock -1 inidcates unlimited.
	MaxRedeals() int

	// Move - how cards are moved from one stack to another.
	Move(source, destination *state.Stack, maxRedeals int) bool
}

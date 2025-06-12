package game

import "github.com/shanehowearth/solitaire"

type Variant interface {
	// Name of the Variant.
	Name() string

	// Foundations.
	Foundations(
		number int,
		basecard solitaire.Rank,
		addRule func(solitaire.SuitedCard) bool,
	)

	// Tableau.
	Tableau(
		number int,
		basecard solitaire.Rank,
		addRule func(solitaire.SuitedCard) bool,
	)
}

package state

// State -
type State struct {
	Deck        *Deck
	Tableau     []*Tableau
	Talon       *Talon
	Foundations []*Foundation
	Reserves    []*Reserve
}

// New - Create a new set of stacks.
func New(
	// Number of decks.
	decks int,
	// Foundation Setup.
	foundationsSpec []StackSpec,
	// Tableau Setup.
	tableauSpec []StackSpec,
	// Reserve Setup.
	reserveSpec []StackSpec,
	// Talon Setup.
	dealCount int,
	perDealCount int,
	talonRule func(SuitedCard) bool,
) *State {
	return &State{
		Deck: CreateDecks(decks),
		Foundations: CreateFoundations(
			foundationsSpec,
		),
		Tableau: CreateTableaus(
			tableauSpec,
		),
		Reserves: CreateReserves(
			reserveSpec,
		),
		Talon: NewTalon(
			dealCount,
			perDealCount,
			talonRule,
		),
	}
}

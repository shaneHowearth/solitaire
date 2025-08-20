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
	numFoundations int,
	foundationBase Rank,
	foundationRule func(Foundation, SuitedCard) bool,
	// Tableau Setup.
	numTableau int,
	tableauBase Rank,
	tableauRule func(*Tableau, SuitedCard) bool,
	// Reserve Setup.
	numReserve int,
	reserveRule func(*Reserve, SuitedCard) bool,
	// Talon Setup.
	dealCount int,
	perDealCount int,
	talonRule func(SuitedCard) bool,
) *State {
	return &State{
		Deck: CreateDecks(decks),
		Foundations: CreateFoundations(
			numFoundations,
			foundationBase,
			foundationRule,
		),
		Tableau: CreateTableaus(
			numTableau,
			tableauBase,
			tableauRule,
		),
		Reserves: CreateReserves(
			numReserve,
			reserveRule,
		),
		Talon: NewTalon(
			dealCount,
			perDealCount,
			talonRule,
		),
	}
}

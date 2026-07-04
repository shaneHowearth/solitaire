package game

import (
	"github.com/shanehowearth/solitaire/state"
)

type FreeCell struct {
	reserves []state.Reserve
	tableaus []state.Tableau
}

var _ Variant = (*FreeCell)(nil)

func (*FreeCell) Name() string { return "FreeCell" }

func (*FreeCell) Category() Category { return CatKlondike }

func (*FreeCell) Description() string {
	return "A solitaire game where nearly all deals are solvable. All cards are dealt face-up into eight cascades, and four open cells are available for temporary storage."
}

func (*FreeCell) TableauGridSize() (int, int) { return 1, 8 }

func (*FreeCell) Decks() int { return 1 }

func (*FreeCell) Reserves() []state.StackSpec {
	rule := func(reserve *state.Stack, _ state.SuitedCard) bool {
		return reserve.Len() == 0
	}
	return []state.StackSpec{
		{AddRule: rule, CardCount: [2]int{0, 0}},
		{AddRule: rule, CardCount: [2]int{0, 0}},
		{AddRule: rule, CardCount: [2]int{0, 0}},
		{AddRule: rule, CardCount: [2]int{0, 0}},
	}
}

func (*FreeCell) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{7, 7}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{7, 7}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{7, 7}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{7, 7}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{6, 6}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{6, 6}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{6, 6}},
		{AddRule: MinusOneAlternatingColorRule, CardCount: [2]int{6, 6}},
	}
}

func (*FreeCell) Fanned() bool { return true }

func (*FreeCell) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}, AddRule: PlusOneRule},
		{BaseCard: state.SuitedCard{Rank: state.Ace, Suit: state.Spades}, AddRule: PlusOneRule},
	}
}

func (*FreeCell) HowToPlay() []string {
	return []string{
		`Build tableaus down by alternating colors.`,
		`Foundations build up by suit from Ace to King.`,
		`Any cell card or top card of any cascade may be moved.`,
		`Empty cells and empty cascades can be used as intermediate storage.`,
		`Number of cards movable at once = (empty free cells + 1) × 2 per empty column.`,
	}
}

func (*FreeCell) HasWon(_ []*state.Tableau, foundations []*state.Foundation) bool {
	for _, f := range foundations {
		if f.Len() != state.RankCount {
			return false
		}
	}
	return true
}

func (*FreeCell) MaxRedeals() int { return 0 }

func (f *FreeCell) freeCellCanMove(
	source, destination *state.Stack,
	reserves []*state.Reserve,
	tableaus []*state.Tableau,
) (bool, int) {
	emptyCells := 0
	for _, r := range reserves {
		if r.Stack.Len() == 0 {
			emptyCells++
		}
	}
	emptyTableaus := 0
	for _, t := range tableaus {
		if t.Stack.Len() == 0 && t.Stack != destination {
			emptyTableaus++
		}
	}
	maxMovable := emptyCells + 1
	for i := 0; i < emptyTableaus; i++ {
		maxMovable *= 2
	}
	canMove, numCards := CanMove(source, destination, true)
	if !canMove || numCards > maxMovable {
		return false, 0
	}
	return true, numCards
}

// func (f *FreeCell) freeCellCanMove(source, destination *state.Stack) (bool, int) {
// 	emptyCells := 0
// 	for _, r := range f.reserves {
// 		if r.Stack.Len() == 0 {
// 			emptyCells++
// 		}
// 	}

// 	emptyTableaus := 0
// 	for _, t := range f.tableaus {
// 		if t.Stack.Len() == 0 && t.Stack != destination {
// 			emptyTableaus++
// 		}
// 	}

// 	maxMovable := emptyCells + 1
// 	for i := 0; i < emptyTableaus; i++ {
// 		maxMovable *= 2
// 	}

// 	canMove, numCards := CanMove(source, destination, true)
// 	if !canMove || numCards > maxMovable {
// 		return false, 0
// 	}

// 	return true, numCards
// }

func (f *FreeCell) Move(source, destination *state.Stack, tableaus []*state.Tableau, reserves []*state.Reserve) bool {
	if destination.Type == state.StackReserve {
		if destination.Len() == 0 {
			card, err := source.Top()
			if err != nil {
				return false
			}
			_, _ = source.Deal()
			destination.Add(card, true)
			return true
		}
		return false
	}

	canMove, numCards := f.freeCellCanMove(source, destination, reserves, tableaus)
	if !canMove || numCards == 0 {
		return false
	}

	// Physical transfer - bypass engine Move to enforce numCards limit
	temp := make([]state.SuitedCard, 0, numCards)
	for i := 0; i < numCards; i++ {
		card, err := source.Top()
		if err != nil {
			break
		}
		temp = append(temp, card)
		_, _ = source.Deal()
	}
	// Reverse onto destination
	for i := len(temp) - 1; i >= 0; i-- {
		destination.Add(temp[i], true)
	}
	// Reveal new top of source
	if source.Len() > 0 {
		top, err := source.Top()
		if err == nil && !top.Visible {
			_, _ = source.Deal()
			source.Add(top, true)
		}
	}
	return true
}

func (f *FreeCell) Compact(_ *state.Stack, _ *state.Stack, _ []*state.Tableau) {}

func (*FreeCell) Talon() bool { return false }

func (*FreeCell) Redeal(_ *state.Talon, _ []*state.Tableau) {}

func (*FreeCell) FoundationBase() bool { return false }

func (f *FreeCell) AvailableMoves(
	tableaus []*state.Tableau,
	foundations []*state.Foundation,
	_ []*state.Talon,
	reserves []*state.Reserve,
) []state.Move {
	var moves []state.Move

	// Tableau to Foundation
	for i := range tableaus {
		for j := range foundations {
			if m := checkMove(tableaus[i].Stack, foundations[j].Stack, false, true); m.NumberMoving > 0 {
				moves = append(moves, m)
			}
		}
	}

	// Tableau to Tableau
	for i := range tableaus {
		for j := range tableaus {
			if i == j {
				continue
			}
			canMove, numCards := f.freeCellCanMove(tableaus[i].Stack, tableaus[j].Stack, reserves, tableaus)
			if canMove && numCards > 0 {
				m := checkMove(tableaus[i].Stack, tableaus[j].Stack, true, true)
				if m.NumberMoving > 0 && m.NumberMoving <= numCards {
					moves = append(moves, m)
				}
			}
		}
	}

	// Reserve to Foundation
	for i := range reserves {
		top, err := reserves[i].Stack.Top()
		if err != nil {
			continue
		}
		for j := range foundations {
			if foundations[j].Stack.Rule(top) {
				moves = append(moves, state.Move{
					Source:        *reserves[i].Stack,
					Destination:   *foundations[j].Stack,
					NumberMoving:  1,
					SourceCardTop: top,
				})
			}
		}
	}

	// Reserve to Tableau
	for i := range reserves {
		top, err := reserves[i].Stack.Top()
		if err != nil {
			continue
		}
		for j := range tableaus {
			if tableaus[j].Stack.Rule(top) {
				moves = append(moves, state.Move{
					Source:        *reserves[i].Stack,
					Destination:   *tableaus[j].Stack,
					NumberMoving:  1,
					SourceCardTop: top,
				})
			}
		}
	}

	// Tableau to Reserve (single card only)
	for i := range tableaus {
		top, err := tableaus[i].Stack.Top()
		if err != nil {
			continue
		}
		for j := range reserves {
			if reserves[j].Stack.Len() == 0 {
				moves = append(moves, state.Move{
					Source:        *tableaus[i].Stack,
					Destination:   *reserves[j].Stack,
					NumberMoving:  1,
					SourceCardTop: top,
				})
				break // one empty cell is enough to show the move exists
			}
		}
	}

	return moves
}

var MinusOneAlternatingColorRule = func(tableau *state.Stack, card state.SuitedCard) bool {
	if (*tableau).Len() == 0 {
		return true
	}
	topCard, err := (*tableau).Top()
	if err != nil {
		return false
	}
	return (card.Suit+topCard.Suit)%2 == 1 && topCard.Rank-card.Rank == 1
}

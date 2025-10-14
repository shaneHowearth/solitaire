package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// AcesSquare - https://en.wikipedia.org/wiki/Monte_Carlo_(card_game)
type AcesSquare struct{}

// Ensure that AcesSquare implements game.Variant.
var _ Variant = (*AcesSquare)(nil)

// Name - name of the variant.
func (*AcesSquare) Name() string {
	return "Aces Square"
}

const numAcesSquareRows = 5
const numAcesSquareCols = 5

// TableauGridSize - The size of the grid required by Aces Square.
func (*AcesSquare) TableauGridSize() (int, int) {
	return numAcesSquareRows, numAcesSquareCols
}

// Decks - How many decks of cards are required to play Aces Square.
func (*AcesSquare) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
// Note that there are no reserves required in a game of Aces Square.
func (*AcesSquare) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (*AcesSquare) Tableau() []state.StackSpec {

	return []state.StackSpec{
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
		{
			AddRule: func(stack *state.Stack, card state.SuitedCard) bool {
				// Cards can only be added if the tableau is empty.
				return stack.Len() < 1
			},
			CardCount: [2]int{1, 1},
		},
	}
}

// Foundations - how the foundations are defined.
func (*AcesSquare) Foundations() []state.StackSpec {
	return []state.StackSpec{}
}

// HowToPlay - Tell the player how to play the game.
func (*AcesSquare) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The game is set up by laying out 25 cards so that they form a 5x5 grid. The rest of the pack is set aside as the stock.
`,
		`Cards that make up a pair (such as two Kings or two Sixes) are removed when they are immediately next to each other horizontally, vertically, or diagonally. Once the pair has been removed, the cards are consolidated, i.e. moving cards to the left as if towards the upper left corner to fill any gaps left behind by the discarded pair. New cards are then laid out from the stock to form a fresh layout of 25 cards.
`,
		`This process is repeated, and continues until it is no longer possible to remove pairs (e.g. in the finishing stages of the game one might be stuck with "4-6-4-6."). The game is won if all cards are successfully removed.
`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
func (*AcesSquare) HasWon(tableaus []*state.Tableau, _ []*state.Foundation) bool {
	// The game is won when all of the tableaus are mpty.
	for _, tableau := range tableaus {
		if tableau.Len() != 0 {
			return false
		}
	}

	return true
}

// MaxRedeals - how many redeals are allowed.
func (*AcesSquare) MaxRedeals() int {
	// There are no redeals.
	return 0
}

// Move -
func (*AcesSquare) Move(first, second *state.Stack, _ []*state.Tableau) bool {
	if second.Type == state.StackWaste && first.Type != state.StackTalon {
		return false
	}

	// Are the first and second "touching"
	// first.
	firstCard, err := first.Top()
	if err != nil {
		return false
	}
	secondCard, err := second.Top()
	if err != nil {
		return false
	}

	// Both card need to have the same rank.
	if firstCard.Rank != secondCard.Rank {
		return false
	}

	// First and Second have to be 'next' to each other.
	// They can be on top of one another, side by side, or on an angle.

	// Convert 1D indices to 2D coordinates
	row1, col1 := first.TableauPosition/numAcesSquareCols, first.TableauPosition%numAcesSquareCols
	row2, col2 := second.TableauPosition/numAcesSquareCols, second.TableauPosition%numAcesSquareCols

	rowDiff := row1 - row2
	if rowDiff < 0 {
		rowDiff = -rowDiff
	}

	colDiff := col1 - col2
	if colDiff < 0 {
		colDiff = -colDiff
	}

	if rowDiff <= 1 && colDiff <= 1 {
		_, _ = first.Deal()
		_, _ = second.Deal()
	}

	return true
}

func (AcesSquare) Compact(stock, waste *state.Stack, tableaus []*state.Tableau) {
	for readIdx := 0; readIdx < numAcesSquareCols*numAcesSquareRows; readIdx++ {
		if tableaus[readIdx].Len() == 0 {
			sourceIdx := -1
			for j := readIdx + 1; j < numAcesSquareCols*numAcesSquareRows; j++ {
				if tableaus[j].Len() > 0 {
					sourceIdx = j
					break
				}
			}

			// If no non-empty tableau found after this position, we're done
			if sourceIdx == -1 {
				break
			}

			// Shift everything from sourceIdx down to readIdx
			for j := sourceIdx; j > readIdx; j-- {
				if tableaus[j].Len() > 0 {
					card, err := tableaus[j].Top()
					if err == nil {
						tableaus[j-1].Stack.Add(card, true)
						_, _ = tableaus[j].Stack.Deal()
					}
				}
			}
			readIdx--
		}
	}

	// Take the waste card and put it onto the second to last grid cell
	wasteCard, err := waste.Deal()
	if err != nil {
		return
	}
	tableaus[numAcesSquareRows*numAcesSquareCols-2].Stack.Add(wasteCard, true)
	// Take the top stock card and put it onto the last grid cell
	stockCard, err := stock.Deal()
	if err != nil {
		return
	}
	tableaus[numAcesSquareRows*numAcesSquareCols-1].Stack.Add(stockCard, true)

	// Take the next stock card and put it onto the waste.
	stockCard, err = stock.Deal()
	if err != nil {
		return
	}
	waste.Add(stockCard, true)
}

// Talon
func (*AcesSquare) Talon() bool {
	return true
}

// Redeal
func (*AcesSquare) Redeal(_ *state.Talon, _ []*state.Tableau) {}

// FoundationBase
func (*AcesSquare) FoundationBase() bool {
	return false
}

// AvailableMoves - return a list of the available moves.
func (*AcesSquare) AvailableMoves([]state.Tableau, []state.Foundation, []state.Talon, []state.Reserve) []string {
	// TODO
	return []string{}
}

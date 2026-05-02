package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Addiction - https://en.wikipedia.org/wiki/Addiction.
type Addiction struct{}

// Ensure that Addiction implements game.Variant.
var _ Variant = (*Addiction)(nil)

// Name - name of the variant.
func (*Addiction) Name() string {
	return "Addiction"
}

func (*Addiction) Category() Category {
	return CatSpecialty
}

func (*Addiction) Description() string {
	return "A high-speed rearrangement game. Shift cards into empty gaps to create sequential rows by suit before you run out of reshuffles."
}

const addictionColumns = 13
const addictionRows = 4

// TableauGridSize - The size of the grid required by acme.
func (*Addiction) TableauGridSize() (int, int) {
	return addictionRows, addictionColumns
}

// Decks - How many decks of cards are required to play acme.
func (*Addiction) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*Addiction) Reserves() []state.StackSpec {
	// There are no reserves.
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (addiction *Addiction) Tableau() []state.StackSpec {
	// There are 52 tableau, each has one card that is open.
	tableau := make([]state.StackSpec, 0, state.RankCount*state.SuitCount)

	for i := 0; i < state.RankCount*state.SuitCount; i++ {
		stateStackSpec := state.StackSpec{
			AddRule:   addiction.tableauRule,
			CardCount: [2]int{1, 1},
			SkipCards: map[state.SuitedCard]struct{}{
				{Rank: state.Ace, Suit: state.Hearts}:   {},
				{Rank: state.Ace, Suit: state.Diamonds}: {},
				{Rank: state.Ace, Suit: state.Clubs}:    {},
				{Rank: state.Ace, Suit: state.Spades}:   {},
			},
		}
		tableau = append(tableau, stateStackSpec)
	}

	return tableau
}

func (*Addiction) Fanned() bool { return false }

func (*Addiction) tableauRule(tableau *state.Stack, _ state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// Foundations - how the foundations are defined.
func (*Addiction) Foundations() []state.StackSpec {
	// There are no foundations.
	return []state.StackSpec{}
}

// HowToPlay - Tell the player how to play the game.
func (*Addiction) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The cards are dealt into four rows of thirteen. The aces are removed and discarded from further play. The addiction that they leave behind are filled by cards that are the same suit and a rank higher than the card on the left of the gap. (For example, 4♣ can be placed beside 3♣.) However, any gap at the right of a King is considered dead and no card can fill it.
`,

		`Any gap on the left hand side of the row should be placed by a deuce and the row should be built up by suit beside the deuce (i. e. 2-3-4-5, etc.). It is the discretion of the player on which suit would occupy which row.
`,
		`When there are no more possible moves, the round is over. The cards that are not in order are gathered, making sure to leave any suit sequence (e.g., ♥2-3-4-5) behind. In some variations, the cards are then shuffled;[1][2] in others, they are not. The cards are then redealt, making sure there is a gap in each row at the immediate right of each suit sequence or at the extreme left of the row if no suit sequence is formed in that row. There are only three rounds in the usual rules. But some variants allow four, and some have no limit.
`,
		`The game is won when all 48 cards are arranged in numerical order and in suits, with the addiction of each row beside the Kings at the extreme right hand of the row.
`,
		`There are three redeals available in this game.`,
	}

	return lines
}

// MaxRedeals - how many redeals are allowed.
func (*Addiction) MaxRedeals() int {
	// Only one redeal is allowed.
	return 3
}

// Move -.
func (addiction *Addiction) Move(source, destination *state.Stack, tableau []*state.Tableau) bool {
	// Nothing to move.
	if source.Len() == 0 {
		return false
	}

	// Waste can only accept cards from the talon.
	if destination.Type == state.StackWaste && source.Type != state.StackTalon {
		return false
	}

	// Only one card per tableau.
	if destination.Len() > 0 {
		return false
	}

	sourceTop, err := source.Top()
	if err != nil {
		return false
	}

	if addiction.checkMove(destination, sourceTop, tableau) {
		// Move the card.
		destination.Add(sourceTop, true)

		_, _ = source.Deal()

		return true
	}

	return false
}

func (*Addiction) checkMove(
	destination *state.Stack,
	sourceTop state.SuitedCard,
	tableau []*state.Tableau,
) bool {
	if destination.TableauPosition%addictionColumns != 0 {
		neighbourTop, err := tableau[destination.TableauPosition-1].Top()
		if err != nil {
			return false
		}

		// Nothing can be placed to the right of a King.
		if neighbourTop.Rank == state.King {
			return false
		}

		// Can only place a card of the same suit as the the card to the left of.
		// the gap.
		if sourceTop.Suit != neighbourTop.Suit {
			return false
		}

		// Card being moved must have a rank one higher than the neighbour.
		if sourceTop.Rank-neighbourTop.Rank != 1 {
			return false
		}

		return true
	}

	// Destination is at the front of the row, which can only take a two (any.
	// two).
	if sourceTop.Rank == state.Two {
		return true
	}

	return false
}

// Compact - .
func (*Addiction) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon - .
func (*Addiction) Talon() bool {
	return false
}

// Redeal - .
func (*Addiction) Redeal(talon *state.Talon, tableau []*state.Tableau) {
	if talon.Stock.CanReceiveMore() {
		GapsRedeal(tableau, addictionRows, addictionColumns)
	}
}

// HasWon - How to tell if the game has been won.
func (*Addiction) HasWon(tableau []*state.Tableau, _ []*state.Foundation) bool {
	for row := 0; row < 4; row++ {
		rowStart := row * addictionColumns
		// Get the suit from the first non-empty card in the row.
		// (since the last position should be empty).
		rowSuit := state.Undefined

		// Find the suit for this row from any non-empty card.
		for col := 0; col < addictionColumns-1; col++ { // Skip last position (should be empty)
			stackIndex := rowStart + col

			card, err := tableau[stackIndex].Top()
			if err != nil {
				continue
			}

			rowSuit = card.Suit

			break
		}

		// If all cards in row are empty (should never happen) fail.
		if rowSuit == state.Undefined {
			return false
		}

		// Check each position in the row.
		for col := 0; col < addictionColumns; col++ {
			stackIndex := rowStart + col

			card, err := tableau[stackIndex].Top()
			if err != nil {
				// An empty stack should only be at the 13th position.
				if (col+1)%addictionColumns != 0 {
					return false
				}

				continue
			}

			expectedRank := state.Rank((col + 1) % addictionColumns)

			// Check rank matches.
			if card.Rank != expectedRank {
				return false
			}

			// Check suit matches (except for empty cards).
			if card.Suit != rowSuit {
				return false
			}
		}
	}

	return true
}

// FoundationBase.
func (*Addiction) FoundationBase() bool {
	return false
}

// AvailableMoves - return a list of the available moves.
func (addiction *Addiction) AvailableMoves(
	tableau []state.Tableau,
	_ []state.Foundation,
	_ []state.Talon,
	_ []state.Reserve,
) []state.Move {
	moves := []state.Move{}

	// create a map of all tableau, with a key of their current value.
	tableauCards := map[state.SuitedCard]state.Tableau{}
	empties := make([]int, 0, 4)

	for idx := range tableau {
		card, err := tableau[idx].Top()
		if err != nil {
			empties = append(empties, idx)
			continue
		}

		tableauCards[card] = tableau[idx]
	}

	for emptyIdx := range empties {
		card := state.SuitedCard{}

		var err error

		if emptyIdx > 0 {
			card, err = tableau[empties[emptyIdx]-1].Top()
			if err != nil {
				// second empty?
				continue
			}

			if card.Rank == state.King {
				// nothing to add.
				continue
			}
		}

		// If the tableau is in the far left position for that row, then any of the.
		// 2s not in the 0th position of any row can be moved.
		// collect the moves.
		// TODO fix the pointer to tableau problem.
		hackedTableau := make([]*state.Tableau, 0, len(tableau))
		for tableauIdx := range tableau {
			hackedTableau = append(hackedTableau, &tableau[tableauIdx])
		}

		for tableauIdx := range tableau {
			sourceTop, err := tableau[tableauIdx].Top()
			if err != nil {
				continue
			}

			if !addiction.checkMove(tableau[empties[emptyIdx]].Stack, sourceTop, hackedTableau) {
				continue
			}

			move := state.Move{
				Source:        *tableau[tableauIdx].Stack,
				Destination:   *tableau[empties[emptyIdx]].Stack,
				NumberMoving:  tableau[tableauIdx].Len(),
				SourceCardTop: sourceTop,
				// DestinationCardBottom will always be empty, the game won't
				// allow a move to any thing other than an empty Destination.
			}

			moves = append(moves, move)
		}
	}

	return moves
}

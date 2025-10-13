package game

import (
	"fmt"
	"log"

	"github.com/shanehowearth/solitaire/state"
)

// Gaps - https://en.wikipedia.org/wiki/Gaps
type Gaps struct{}

// Ensure that Gaps implements game.Variant.
var _ Variant = (*Gaps)(nil)

// Name - name of the variant.
func (*Gaps) Name() string {
	return "Gaps"
}

const gapsColumns = 13
const gapsRows = 4

// TableauGridSize - The size of the grid required by acme.
func (*Gaps) TableauGridSize() (int, int) {
	return gapsRows, gapsColumns
}

// Decks - How many decks of cards are required to play acme.
func (*Gaps) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*Gaps) Reserves() []state.StackSpec {
	// There are no reserves.
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (gaps *Gaps) Tableau() []state.StackSpec {
	// There are 52 tableau, each has one card that is open.
	tableau := make([]state.StackSpec, 0, 52)
	for i := 0; i < 52; i++ {
		t := state.StackSpec{
			AddRule:   gaps.tableauRule,
			CardCount: [2]int{1, 1},
			SkipCards: map[state.SuitedCard]struct{}{
				state.SuitedCard{Rank: state.Ace, Suit: state.Hearts}:   struct{}{},
				state.SuitedCard{Rank: state.Ace, Suit: state.Diamonds}: struct{}{},
				state.SuitedCard{Rank: state.Ace, Suit: state.Clubs}:    struct{}{},
				state.SuitedCard{Rank: state.Ace, Suit: state.Spades}:   struct{}{},
			},
		}
		tableau = append(tableau, t)
	}
	return tableau
}

func (*Gaps) tableauRule(tableau *state.Stack, card state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// Foundations - how the foundations are defined.
func (*Gaps) Foundations() []state.StackSpec {
	// There are no foundations.
	return []state.StackSpec{}
}

// HowToPlay - Tell the player how to play the game.
func (*Gaps) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The cards are dealt into four rows of thirteen. The aces are removed and discarded from further play. The gaps that they leave behind are filled by cards that are the same suit and a rank higher than the card on the left of the gap. (For example, 4♣ can be placed beside 3♣.) However, any gap at the right of a King is considered dead and no card can fill it.
`,

		`Any gap on the left hand side of the row should be placed by a deuce and the row should be built up by suit beside the deuce (i. e. 2-3-4-5, etc.). It is the discretion of the player on which suit would occupy which row.
`,
		`When there are no more possible moves, the round is over. The cards that are not in order are gathered, making sure to leave any suit sequence (e.g., ♥2-3-4-5) behind. In some variations, the cards are then shuffled;[1][2] in others, they are not. The cards are then redealt, making sure there is a gap in each row at the immediate right of each suit sequence or at the extreme left of the row if no suit sequence is formed in that row. There are only three rounds in the usual rules. But some variants allow four, and some have no limit.
`,
		`The game is won when all 48 cards are arranged in numerical order and in suits, with the gaps of each row beside the Kings at the extreme right hand of the row.
`,
		`There is only one (1) redeal available in this game.`,
	}

	return lines
}

// MaxRedeals - how many redeals are allowed.
func (*Gaps) MaxRedeals() int {
	// Only one redeal is allowed.
	return 1
}

// Move -
func (gaps *Gaps) Move(source, destination *state.Stack, tableau []*state.Tableau) bool {
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

	if gaps.checkMove(destination, sourceTop, tableau) {
		// Move the card.
		destination.Add(sourceTop, true)
		_, _ = source.Deal()

		return true
	}

	return false
}

func (*Gaps) checkMove(
	destination *state.Stack,
	sourceTop state.SuitedCard,
	tableau []*state.Tableau,
) bool {

	if destination.TableauPosition%gapsColumns != 0 {
		neighbourTop, err := tableau[destination.TableauPosition-1].Top()
		if err != nil {
			return false
		}

		// Nothing can be placed to the right of a King.
		if neighbourTop.Rank == state.King {
			return false
		}

		// Can only place a card of the same suit as the the card to the left of
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

	// destination is at the front of the row, which can only take a two (any
	// two).
	if sourceTop.Rank == state.Two {
		return true
	}

	return false
}

// Compact
func (*Gaps) Compact(_, _ *state.Stack, _ []*state.Tableau) {}

// Talon
func (*Gaps) Talon() bool {
	return false
}

// Redeal
func (gaps *Gaps) Redeal(talon *state.Talon, tableau []*state.Tableau) {
	if talon.Stock.CanReceiveMore() {
		GapsRedeal(tableau, gapsRows, gapsColumns)
	}
}

// HasWon - How to tell if the game has been won.
func (*Gaps) HasWon(tableau []*state.Tableau, _ []*state.Foundation) bool {

	for row := 0; row < 4; row++ {
		rowStart := row * gapsColumns
		// Get the suit from the first non-empty card in the row
		// (since the last position should be empty)
		rowSuit := state.Undefined

		// Find the suit for this row from any non-empty card
		for col := 0; col < gapsColumns-1; col++ { // Skip last position (should be empty)
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

		// Check each position in the row
		for col := 0; col < gapsColumns; col++ {
			stackIndex := rowStart + col
			card, err := tableau[stackIndex].Top()
			if err != nil {
				// An empty stack should only be at the 13th position.
				if (col+1)%13 != 0 {
					return false
				}

				continue
			}

			expectedRank := state.Rank((col + 1) % 13)

			// Check rank matches
			if card.Rank != expectedRank {
				return false
			}

			// Check suit matches (except for empty cards)
			if card.Suit != rowSuit {
				return false
			}
		}
	}

	return true
}

// FoundationBase
func (*Gaps) FoundationBase() bool {
	return false
}

// AvailableMoves - return a list of the available moves.
func (gaps *Gaps) AvailableMoves(
	tableau []state.Tableau,
	foundations []state.Foundation,
	talon []state.Talon,
	reserves []state.Reserve,
) []string {
	moves := []string{}
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

	log.Printf("Empties %#v", empties)
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

		// If the tableau is in the far left position for that row, then any of the
		// 2s not in the 0th position of any row can be moved.
		// collect the moves
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
			if !gaps.checkMove(tableau[empties[emptyIdx]].Stack, sourceTop, hackedTableau) {
				continue
			}

			cardStr := fmt.Sprintf("%s %s can be moved from row %d column %d to row %d column %d",
				sourceTop.Rank.String(), sourceTop.Suit.String(),
				tableauIdx/gapsColumns+1, tableauIdx%gapsColumns+1,
				empties[emptyIdx]/gapsColumns+1, empties[emptyIdx]%gapsColumns+1,
			)
			moves = append(moves, cardStr)
		}
	}

	return moves
}

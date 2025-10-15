package game

import (
	"fmt"

	"github.com/shanehowearth/solitaire/state"
)

// Accordian - https://en.wikipedia.org/wiki/Accordion_(card_game)
type Accordian struct{}

// Ensure that Accordian implements game.Variant.
var _ Variant = (*Accordian)(nil)

// Name - name of the variant.
func (*Accordian) Name() string {
	return "Accordian"
}

// TableauGridSize - The size of the grid required by acme.
func (*Accordian) TableauGridSize() (int, int) {
	const accordianColumns = 13
	const accordianRows = 4

	return accordianRows, accordianColumns
}

// Decks - How many decks of cards are required to play acme.
func (*Accordian) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*Accordian) Reserves() []state.StackSpec {
	// There are no reserves.
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (accordian *Accordian) Tableau() []state.StackSpec {
	// There are 52 tableau, each has one card that is open.
	tableau := make([]state.StackSpec, 0, 52)
	for i := 0; i < 52; i++ {
		t := state.StackSpec{
			AddRule:   accordian.tableauRule,
			CardCount: [2]int{1, 1},
			ShowCount: 1,
		}
		tableau = append(tableau, t)
	}
	return tableau
}

func (*Accordian) tableauRule(tableau *state.Stack, card state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// Foundations - how the foundations are defined.
func (*Accordian) Foundations() []state.StackSpec {
	// There are no foundations.
	return []state.StackSpec{}
}

// HowToPlay - Tell the player how to play the game.
func (*Accordian) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`The cards from the entire deck are spread out in a single line (displayed as 4 lines of 13).
`,

		`A pile can be moved on top of another pile immediately to its left or moved three piles to its left if the top cards of each pile have the same suit or rank.
`,
		`Gaps left behind are filled by moving piles to the left. The player is not required to make a particular move if they prefer not to
`,
		`There is no redeal available in this game.`,
	}

	return lines
}

// MaxRedeals - how many redeals are allowed.
func (*Accordian) MaxRedeals() int {
	// Only one redeal is allowed.
	return 1
}

// Move -
func (accordian *Accordian) Move(source, destination *state.Stack, tableau []*state.Tableau) bool {

	if !accordian.checkMove(source, destination) {
		return false
	}
	// Move the cards.
	// Need to move the whole stack.

	// Step 1. Reverse the source cards (need the top of this pile to be the top
	// of the new pile).
	source.Reverse()

	// Step 2. Move to the destination.
	for {
		sourceTop, err := source.Top()
		if err != nil {
			break
		}

		destination.Add(sourceTop, true)
		_, _ = source.Deal()
	}

	return true
}

// Compact
func (*Accordian) Compact(_, _ *state.Stack, tableau []*state.Tableau) {
	// Move the card(s) to the left of an empty stack into the empty stack.
	for readIdx := range tableau {
		if tableau[readIdx].Len() == 0 {
			sourceIdx := -1
			for j := readIdx + 1; j < len(tableau); j++ {
				if tableau[j].Len() > 0 {
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
				if tableau[j].Len() > 0 {
					tableau[j].Stack.Reverse()

					// // Step 2. Move to the destination.
					for {
						sourceTop, err := tableau[j].Top()
						if err != nil {
							break
						}

						tableau[j-1].Stack.Add(sourceTop, true)
						_, _ = tableau[j].Stack.Deal()
					}
				}
			}
			readIdx--
		}
	}
}

// Talon
func (*Accordian) Talon() bool {
	return false
}

// Redeal
func (accordian *Accordian) Redeal(_ *state.Talon, _ []*state.Tableau) {}

// HasWon - How to tell if the game has been won.
func (*Accordian) HasWon(tableau []*state.Tableau, _ []*state.Foundation) bool {
	// All the cards should have accumulated into the first tableau.
	return tableau[0].Len() == 52
}

// FoundationBase
func (*Accordian) FoundationBase() bool {
	return false
}

// AvailableMoves - return a list of the available moves.
func (*Accordian) AvailableMoves([]state.Tableau, []state.Foundation, []state.Talon, []state.Reserve) []string {
	// TODO
	return []string{}
}

func (*Accordian) checkMove(source, destination *state.Stack) bool {
	// Nothing to move.
	if source.Len() == 0 || destination.Len() == 0 {
		return false
	}

	// destination location must be 1 or 3 places to the left, no more, no less.
	diffPosition := source.TableauPosition - destination.TableauPosition
	if diffPosition != 1 && diffPosition != 3 {
		return false
	}

	topDestination, err := destination.Top()
	if err != nil {
		return false
	}

	topSource, err := source.Top()
	if err != nil {
		return false
	}

	if topDestination.Rank != topSource.Rank && topDestination.Suit != topSource.Suit {
		return false
	}

	return true
}

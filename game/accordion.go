package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// Accordion - https://en.wikipedia.org/wiki/Accordion_(card_game).
type Accordion struct{}

// Ensure that Accordion implements game.Variant.
var _ Variant = (*Accordion)(nil)

// Name - name of the variant.
func (*Accordion) Name() string {
	return "Accordion"
}

func (*Accordion) Category() Category {
	return CatSpecialty
}

func (*Accordion) Description() string {
	return "A unique one-row game where you stack cards or piles on top of neighbors if they match suit or rank."
}

// TableauGridSize - The size of the grid required by acme.
func (*Accordion) TableauGridSize() (int, int) {
	const (
		accordionColumns = 13
		accordionRows    = 4
	)

	return accordionRows, accordionColumns
}

func (*Accordion) Fanned() bool { return false }

// Decks - How many decks of cards are required to play acme.
func (*Accordion) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*Accordion) Reserves() []state.StackSpec {
	// There are no reserves.
	return []state.StackSpec{}
}

// Tableau - how the tableau are defined.
func (accordion *Accordion) Tableau() []state.StackSpec {
	// There are 52 tableau, each has one card that is open.
	tableau := make([]state.StackSpec, 0, state.RankCount*state.SuitCount)

	for i := 0; i < state.RankCount*state.SuitCount; i++ {
		t := state.StackSpec{
			AddRule:   accordion.tableauRule,
			CardCount: [2]int{1, 1},
			ShowCount: 1,
		}
		tableau = append(tableau, t)
	}

	return tableau
}

func (*Accordion) tableauRule(tableau *state.Stack, _ state.SuitedCard) bool {
	// Handle when the tableau is empty.
	if tableau.Len() == 0 {
		// Anything can be put onto an empty tableau.
		return true
	}

	// All other cases the card should not be added to the tableau.
	return false
}

// Foundations - how the foundations are defined.
func (*Accordion) Foundations() []state.StackSpec {
	// There are no foundations.
	return []state.StackSpec{}
}

// HowToPlay - Tell the player how to play the game.
func (*Accordion) HowToPlay() []string {
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
func (*Accordion) MaxRedeals() int {
	// Only one redeal is allowed.
	return 1
}

// Move -.
func (accordion *Accordion) Move(source, destination *state.Stack, _ []*state.Tableau, _ []*state.Reserve) bool {
	if !accordion.checkMove(source, destination) {
		return false
	}
	// Move the cards.
	// Need to move the whole stack.

	// Step 1. Reverse the source cards (need the top of this pile to be the top.
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

// Compact.
func (*Accordion) Compact(_, _ *state.Stack, tableau []*state.Tableau) {
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

			// If no non-empty tableau found after this position, we're done.
			if sourceIdx == -1 {
				break
			}

			// Shift everything from sourceIdx down to readIdx.
			for shiftIdx := sourceIdx; shiftIdx > readIdx; shiftIdx-- {
				if tableau[shiftIdx].Len() > 0 {
					tableau[shiftIdx].Stack.Reverse()

					// Move to the destination.
					for {
						sourceTop, err := tableau[shiftIdx].Top()
						if err != nil {
							break
						}

						tableau[shiftIdx-1].Stack.Add(sourceTop, true)
						_, _ = tableau[shiftIdx].Stack.Deal()
					}
				}
			}
		}
	}
}

// Talon.
func (*Accordion) Talon() bool {
	return false
}

// Redeal.
func (*Accordion) Redeal(_ *state.Talon, _ []*state.Tableau) {}

// HasWon - How to tell if the game has been won.
func (*Accordion) HasWon(tableau []*state.Tableau, _ []*state.Foundation) bool {
	// All the cards should have accumulated into the first tableau.
	return tableau[0].Len() == state.RankCount*state.SuitCount
}

// FoundationBase.
func (*Accordion) FoundationBase() bool {
	return false
}

func (*Accordion) checkMove(source, destination *state.Stack) bool {
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

// AvailableMoves - return a list of the available moves.
func (accordion *Accordion) AvailableMoves(
	tableau []*state.Tableau,
	_ []*state.Foundation,
	_ []*state.Talon,
	_ []*state.Reserve,
) []state.Move {
	moves := []state.Move{}
	// create a map of all tableau, with a key of their current value.
	for idx := range tableau {
		// There are two possible stacks that need to be checked, +1 (the stack.
		// next to this one) and +3, the stack 3 across.
		// Only checking those gives us a big o of O(2n).
		checkOne := idx + 1
		checkThree := idx + 3

		if checkOne >= len(tableau) {
			break
		}

		destinationTop, err := tableau[idx].Top()
		if err != nil {
			continue
		}

		if accordion.checkMove(tableau[checkOne].Stack, tableau[idx].Stack) {
			sourceTop, _ := tableau[checkOne].Top()
			moves = append(moves,
				state.Move{
					Source:                *tableau[checkOne].Stack,
					Destination:           *tableau[idx].Stack,
					NumberMoving:          tableau[checkOne].Len(),
					SourceCardTop:         sourceTop,
					DestinationCardBottom: destinationTop,
				})
		}

		if checkThree >= len(tableau) {
			continue
		}

		if accordion.checkMove(tableau[checkThree].Stack, tableau[idx].Stack) {
			sourceTop, _ := tableau[checkThree].Top()
			moves = append(moves,
				state.Move{
					Source:                *tableau[checkThree].Stack,
					Destination:           *tableau[idx].Stack,
					NumberMoving:          tableau[checkThree].Len(), // The whole source is always moving.
					SourceCardTop:         sourceTop,
					DestinationCardBottom: destinationTop,
				})
		}
	}

	return moves
}

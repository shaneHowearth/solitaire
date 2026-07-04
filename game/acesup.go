package game

import (
	"github.com/shanehowearth/solitaire/state"
)

// AcesUp - https://en.wikipedia.org/wiki/Aces_Up.
type AcesUp struct{}

// Ensure that AcesUp implements game.Variant.
var _ Variant = (*AcesUp)(nil)

// Name - name of the variant.
func (*AcesUp) Name() string {
	return "Aces Up"
}

func (*AcesUp) Category() Category {
	return CatPairing // Or CatSpecialty, but since it involves rank comparison, Pairing fits well.
}

func (*AcesUp) Description() string {
	return "A fast-paced elimination game. If two top cards have the same suit, discard the one with the lower rank. Try to leave only the four Aces."
}

// TableauGridSize - The size of the grid required by Aces Up.
func (*AcesUp) TableauGridSize() (int, int) {
	const (
		height         = 1
		numAcesUpPiles = 4
	)

	return height, numAcesUpPiles
}

// Decks - How many decks of cards are required to play Aces Up.
func (*AcesUp) Decks() int {
	return 1
}

// Reserves - how the reserves are defined.
func (*AcesUp) Reserves() []state.StackSpec {
	return []state.StackSpec{}
}

// AcesUpRule - A rule for Aces Up tableau piles where only single cards
// can be placed on empty piles.
var AcesUpRule = func(tableau *state.Stack, card state.SuitedCard) bool {
	// Can only place a card on an empty tableau pile
	if tableau.Len() == 0 {
		return true
	}

	// Cannot stack cards on top of each other in Aces Up
	return false
}

// DiscardRule checks if a card can be moved to the discard pile.
var DiscardRule = func(tableauPiles []*state.Stack, card state.SuitedCard) bool {
	// Treat Ace as 14 for comparison
	cardRank := card.Rank
	if cardRank == 1 {
		cardRank = 14
	}

	for _, stack := range tableauPiles {
		if stack.Len() == 0 {
			continue
		}

		top, err := stack.Top()
		if err != nil {
			continue
		}

		topRank := top.Rank
		if topRank == 1 {
			topRank = 14
		}

		// Rule: Same suit, and there is a higher card elsewhere
		if top.Suit == card.Suit && topRank > cardRank {
			return true
		}
	}
	return false
}

// Tableau - how the tableau are defined.
func (*AcesUp) Tableau() []state.StackSpec {
	return []state.StackSpec{
		{
			AddRule:   AcesUpRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{},
		},
		{
			AddRule:   AcesUpRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{},
		},
		{
			AddRule:   AcesUpRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{},
		},
		{
			AddRule:   AcesUpRule,
			CardCount: [2]int{1, 1},
			BaseCard:  state.SuitedCard{},
		},
	}
}

func (*AcesUp) Fanned() bool { return false }

// TableauPosition - Where does each tableau go in the grid, and what angle (relative to
// straight up and down) should the tableau be twisted.
func (*AcesUp) TableauPosition(tableauNumber int) (int, int, int) {
	const x = 0
	const angle = 0

	return x, tableauNumber, angle
}

// Foundations - how the foundations are defined.
// Aces Up has one foundation that acts as a discard pile.
// Foundations - Aces Up has one discard pile, but engine requires
// multiples of suits (4).
func (*AcesUp) Foundations() []state.StackSpec {
	return []state.StackSpec{
		{
			BaseCard: state.SuitedCard{},
			AddRule:  func(s *state.Stack, c state.SuitedCard) bool { return true },
		},
		{
			BaseCard: state.SuitedCard{},
			AddRule:  func(s *state.Stack, c state.SuitedCard) bool { return true },
		},
		{
			BaseCard: state.SuitedCard{},
			AddRule:  func(s *state.Stack, c state.SuitedCard) bool { return true },
		},
		{
			BaseCard: state.SuitedCard{},
			AddRule:  func(s *state.Stack, c state.SuitedCard) bool { return true },
		},
	}
}

// HowToPlay - Tell the player how to play the game.
func (*AcesUp) HowToPlay() []string {
	//nolint:lll // Long line needed to preserve readability of game rules text.
	lines := []string{
		`Deal four cards in a row face up. If there are two or more cards of the same suit, discard all but the highest-ranked card of that suit to the foundation. Aces rank high.`,
		`Whenever there are empty spaces, you may move the top card of another pile into the empty space.`,
		`When there are no more cards to move or remove, deal out the next four cards from the deck face-up onto each pile.`,
		`The goal is to discard all cards except the four aces. The game ends when all cards have been dealt and no more moves are possible.`,
		`Your score is the number of discarded cards. The maximum score (and winning score) is 48.`,
	}

	return lines
}

// HasWon - How to tell if the game has been won.
// Win condition: only 4 cards remain in tableau (the four aces).
func (*AcesUp) HasWon(tableau []*state.Tableau, foundations []*state.Foundation) bool {
	const winningCards = 4
	const totalCards = 52

	totalTableau := 0
	for _, p := range tableau {
		totalTableau += p.Len()
	}

	totalFoundation := 0
	for _, f := range foundations {
		totalFoundation += f.Len()
	}

	return totalTableau == winningCards && totalFoundation == (totalCards-winningCards)
}

// MaxRedeals - how many redeals are allowed.
func (*AcesUp) MaxRedeals() int {
	return 0
}

// Move - handles moving cards between piles.
func (a *AcesUp) Move(source, destination *state.Stack, allTableaus []*state.Tableau, _ []*state.Reserve) bool {
	card, _ := source.Top()

	// 1. DISCARD LOGIC
	if destination.Type == state.StackFoundation {
		// Find if ANY other tableau has a higher card of the same suit
		for _, t := range allTableaus {
			if t.Stack.Len() == 0 {
				continue
			}
			top, _ := t.Stack.Top()

			// Check: Same suit AND other card is stronger
			if top.Suit == card.Suit && rankStrength(top.Rank) > rankStrength(card.Rank) {
				// EXECUTE MOVE
				c, _ := source.Deal()
				destination.Add(c, true)
				return true
			}
		}
		return false // No higher card found, discard blocked
	}

	// 2. SPACE FILLING LOGIC
	if destination.Type == state.StackTableau && destination.Len() == 0 {
		c, _ := source.Deal()
		destination.Add(c, true)
		return true
	}

	return false
}

// Compact.
func (*AcesUp) Compact(_, _ *state.Stack, _ []*state.Tableau) {
}

// Talon.
func (*AcesUp) Talon() bool {
	return true
}

// Redeal - deals 4 cards to the 4 piles.
func (*AcesUp) Redeal(talon *state.Talon, tableau []*state.Tableau) {
	for idx := range tableau {
		if talon.Stock.Len() > 0 {
			card, _ := talon.Stock.Deal()
			tableau[idx].Stack.Add(card, true)
		}
	}
}

// FoundationBase.
func (*AcesUp) FoundationBase() bool {
	return false
}
func rankStrength(r state.Rank) int {
	// If it's an Ace (0), give it the highest strength (13)
	if r == state.Ace {
		return 13
	}
	// Otherwise, use its iota value (1-12)
	return int(r)
}

// AvailableMoves - return a list of the available moves.
func (*AcesUp) AvailableMoves(
	tableau []*state.Tableau,
	foundations []*state.Foundation,
	_ []*state.Talon,
	_ []*state.Reserve,
) []state.Move {
	moves := []state.Move{}

	// Check for cards that can be discarded to foundation.
	// Find cards of same suit where a higher card exists.
	for sourceIdx := range tableau {
		sourceTop, err := tableau[sourceIdx].Stack.Top()
		if err != nil {
			continue
		}

		canDiscard := false
		for checkIdx := range tableau {
			if sourceIdx == checkIdx {
				continue
			}

			checkTop, err := tableau[checkIdx].Stack.Top()
			if err != nil {
				continue
			}

			// Compare using virtual strength
			if sourceTop.Suit == checkTop.Suit && rankStrength(checkTop.Rank) > rankStrength(sourceTop.Rank) {
				canDiscard = true
				break
			}
		}

		if canDiscard {
			moves = append(moves, state.Move{
				Source:        *tableau[sourceIdx].Stack, // The lower card is the source
				Destination:   *foundations[0].Stack,
				NumberMoving:  1,
				SourceCardTop: sourceTop,
			})
		}
	}

	// Check for moves to empty tableau positions.
	emptyPiles := []int{}
	for idx := range tableau {
		if tableau[idx].Len() == 0 {
			emptyPiles = append(emptyPiles, idx)
		}
	}

	for _, emptyIdx := range emptyPiles {
		for sourceIdx := range tableau {
			if sourceIdx == emptyIdx || tableau[sourceIdx].Len() == 0 {
				continue
			}

			sourceTop, _ := tableau[sourceIdx].Stack.Top()
			move := state.Move{
				Source:        *tableau[sourceIdx].Stack,
				Destination:   *tableau[emptyIdx].Stack,
				NumberMoving:  1,
				SourceCardTop: sourceTop,
			}
			moves = append(moves, move)
		}
	}

	return moves
}

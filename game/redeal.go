package game

import "github.com/shanehowearth/solitaire/state"

func GapsRedeal(tableau []*state.Tableau, rowCount, colCount int) {

	// Step 1: Find where each row's unsolved section starts
	unsolvedStartPositions := make([]int, rowCount)
	cardsToShuffle := state.NewStack(
		rowCount,
		state.SuitedCard{},
		func(state.SuitedCard) bool { return true },
		state.StackUndefined,
	)

	for row := 0; row < rowCount; row++ {
		rowStart := row * colCount
		unsolvedStartPositions[row] = findUnsolvedStart(tableau, rowCount, colCount, row)

		// Collect cards from unsolved positions for shuffling
		for col := unsolvedStartPositions[row]; col < colCount; col++ {
			stackIndex := rowStart + col
			card, err := tableau[stackIndex].Top()
			if err == nil { // If there's a card here
				cardsToShuffle.Add(card, true)
				// Remove the card from the tableau
				_, _ = tableau[stackIndex].Stack.Deal()
			}
		}
	}

	// Step 2: Shuffle the collected cards
	cardsToShuffle.Shuffle()

	// Step 3: Redistribute shuffled cards back to unsolved positions
	// but leave the gaps (empty positions) where they should be
	for row := 0; row < rowCount; row++ {
		rowStart := row * colCount

		unsolvedStart := unsolvedStartPositions[row]
		// The gap should be at position unsolvedStart (right after solved cards)
		gapPosition := unsolvedStart

		for col := unsolvedStart; col < colCount; col++ {
			stackIndex := rowStart + col

			// Skip the gap position (leave it empty)
			if col == gapPosition {
				continue
			}

			if cardsToShuffle.Len() > 0 {
				top, err := cardsToShuffle.Top()
				if err == nil {
					tableau[stackIndex].Stack.Add(top, true)
					cardsToShuffle.Deal()
				}
			}
		}
	}
}

// findUnsolvedStart finds where the unsolved section starts in a row
func findUnsolvedStart(tableau []*state.Tableau, rowCount, colCount, row int) int {
	rowStart := row * colCount

	// Get the expected suit for this row from the first correctly placed card
	expectedSuit := getRowSuit(tableau, row, colCount)
	if expectedSuit == state.Undefined {
		return 0 // Entire row needs to be reshuffled
	}

	expectedRanks := []state.Rank{
		state.Two, state.Three, state.Four, state.Five, state.Six,
		state.Seven, state.Eight, state.Nine, state.Ten, state.Jack,
		state.Queen, state.King,
	}
	// Find first incorrect position
	for col := 0; col < colCount; col++ {
		stackIndex := rowStart + col
		card, err := tableau[stackIndex].Top()

		if col == colCount-1 {
			// Last position should be empty
			if err == nil {
				return col // Found a card where there should be empty space
			}
			continue // This position is correctly empty
		}

		// Check if card is in correct position
		expectedRank := expectedRanks[col]
		if err != nil || card.Rank != expectedRank || card.Suit != expectedSuit {
			return col // This is where unsolved section starts
		}
	}

	return colCount // Entire row is solved
}

// getRowSuit determines the suit for a row based on correctly placed cards
func getRowSuit(tableau []*state.Tableau, row, colCount int) state.Suit {
	rowStart := row * colCount
	expectedRanks := []state.Rank{
		state.Two, state.Three, state.Four, state.Five, state.Six,
		state.Seven, state.Eight, state.Nine, state.Ten, state.Jack,
		state.Queen, state.King,
	}

	// Look for the first correctly placed card to determine suit
	for col := 0; col < colCount-1; col++ { // Skip last position (should be empty)
		stackIndex := rowStart + col
		card, err := tableau[stackIndex].Top()
		if err == nil && card.Rank == expectedRanks[col] {
			return card.Suit
		}
	}

	return state.Undefined
}

package state

import (
	"errors"
	"fmt"
	"testing"
)

// --- MOCK/HELPER DEFINITIONS FOR TESTING ---
// Since the original code did not define Stack, SuitedCard, or Len(),
// these mocks are required to make the tests runnable.

// var ErrEmpty = errors.New("Empty") // Mock the error from the original package

// newTestStack creates a new Stack instance with a deep copy of the initial cards.
func newTestStack(initialCards []SuitedCard) *Stack {
	cardsCopy := make([]SuitedCard, len(initialCards))
	copy(cardsCopy, initialCards)
	return &Stack{cards: &cardsCopy}
}

func TestStack_Add(t *testing.T) {
	initialCards := []SuitedCard{{Rank: Ace, Suit: Hearts}}
	stack := newTestStack(initialCards)
	initialLen := stack.Len()

	newCard := SuitedCard{Rank: Ten, Suit: Spades}

	// Test Add with visible = true
	stack.Add(newCard, true)
	if stack.Len() != initialLen+1 {
		t.Errorf("Add failed: expected length %d, got %d", initialLen+1, stack.Len())
	}

	topCard, _ := stack.Top()
	if topCard.Rank != Ten || topCard.Visible != true {
		t.Errorf("Add failed: expected card {10, Spades, true}, got %+v", topCard)
	}

	// Test Add with visible = false
	newCard2 := SuitedCard{Rank: Five, Suit: Hearts}
	stack.Add(newCard2, false)
	topCard2, _ := stack.Top()
	if topCard2.Visible != false {
		t.Errorf("Add failed: expected Visible=false, got Visible=%t", topCard2.Visible)
	}
}

func TestStack_Deal(t *testing.T) {
	cards := []SuitedCard{
		{Rank: Ace, Suit: Hearts},
		{Rank: Two, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
	}
	stack := newTestStack(cards)
	initialLen := stack.Len() // 3

	// 1. Test successful deal
	card, err := stack.Deal()
	if err != nil {
		t.Fatalf("Deal failed unexpectedly: %v", err)
	}
	if card.Rank != Three || card.Suit != Clubs {
		t.Errorf("Deal returned wrong card: expected {3, C}, got %+v", card)
	}
	if stack.Len() != initialLen-1 {
		t.Errorf("Deal failed: expected length %d, got %d", initialLen-1, stack.Len())
	}

	// 2. Deal remaining cards
	_, _ = stack.Deal() // Rank 2
	_, _ = stack.Deal() // Rank 1

	// 3. Test deal on empty stack (error case)
	_, err = stack.Deal()
	if !errors.Is(err, ErrEmpty) {
		t.Errorf("Deal on empty stack: expected ErrEmpty, got %v", err)
	}
}

func TestStack_Reverse(t *testing.T) {
	cards := []SuitedCard{
		{Rank: Ace, Suit: Hearts},
		{Rank: Two, Suit: Diamonds},
		{Rank: Three, Suit: Clubs},
	}
	stack := newTestStack(cards)

	stack.Reverse()

	// Check the new order
	reverseTopZero, _ := stack.Deal()
	reverseTopOne, _ := stack.Deal()
	reverseTopTwo, _ := stack.Deal()
	if reverseTopZero.Rank != Ace || reverseTopOne.Rank != Two || reverseTopTwo.Rank != Three {
		t.Errorf("Reverse failed. Expected order [2, 1, 0], got [%d, %d, %d]",
			reverseTopZero.Rank, reverseTopOne.Rank, reverseTopTwo.Rank)
	}

	// Test case with a single card
	singleCardStack := newTestStack(
		[]SuitedCard{
			{Rank: Ten, Suit: Hearts},
		})
	singleCardStack.Reverse()
	singleTop, _ := singleCardStack.Top()
	if singleCardStack.Len() != 1 || singleTop.Rank != Ten {
		t.Error("Reverse failed on single card stack")
	}

	// Test case with empty stack (should not panic)
	emptyStack := newTestStack([]SuitedCard{})
	emptyStack.Reverse()
	if emptyStack.Len() != 0 {
		t.Error("Reverse failed on empty stack")
	}
}

func TestStack_Clone(t *testing.T) {
	originalCards := []SuitedCard{
		{Rank: Ten, Visible: true},
		{Rank: Eight, Visible: false},
	}
	originalStack := newTestStack(originalCards)

	cloneStack := originalStack.Clone()

	// 1. Check lengths and initial values are the same
	if originalStack.Len() != cloneStack.Len() {
		t.Fatalf("Clone failed: Lengths do not match. Original: %d, Clone: %d", originalStack.Len(), cloneStack.Len())
	}
	originalTop, _ := originalStack.Top()
	clonedTop, _ := cloneStack.Top()
	if originalTop != clonedTop {
		t.Fatal("Clone failed: Card values do not match initially.")
	}

	// 2. Crucial test: Check for deep copy (modifying clone should not affect original)

	// Modify the clone
	cloneStack.Reverse()

	// Check that the original remains unchanged
	originalTop, _ = originalStack.Top()
	if originalTop.Rank != Eight {
		t.Errorf("Clone failed: Deep copy was not performed. Modifying clone modified original, got %s.", originalTop.Rank.String())
	}
	if originalStack.Len() != 2 {
		t.Errorf("Clone failed: Adding to clone incorrectly changed original length. Expected 2, got %d", originalStack.Len())
	}

	// Check that the clone was indeed modified
	clonedTop, _ = cloneStack.Top()
	if cloneStack.Len() != 2 {
		t.Errorf("Clone modification failed: expected clone length 2, got %d", cloneStack.Len())
	}
	if clonedTop.Rank != Ten {
		t.Errorf("Clone modification failed: Rank not updated in clone, got %s.", clonedTop.Rank.String())
	}

	// Ensure an Add to the clone does not affect the original.
	cloneStack.Add(SuitedCard{Rank: Jack, Suit: Spades}, true)
	// Check that the original remains unchanged
	originalTop, _ = originalStack.Top()
	if originalTop.Rank != Eight {
		t.Errorf("Clone failed: Deep copy was not performed. Modifying clone modified original, got %s.", originalTop.Rank.String())
	}
	if originalStack.Len() != 2 {
		t.Errorf("Clone failed: Adding to clone incorrectly changed original length. Expected 2, got %d", originalStack.Len())
	}

	// Check that the clone was indeed modified
	clonedTop, _ = cloneStack.Top()
	if cloneStack.Len() != 3 {
		t.Errorf("Clone modification failed: expected clone length 3, got %d", cloneStack.Len())
	}
	if clonedTop.Rank != Jack {
		t.Errorf("Clone modification failed: Rank not updated in clone, got %s.", clonedTop.Rank.String())
	}

	// 3. Pointer isolation check
	if fmt.Sprintf("%p", originalStack.cards) == fmt.Sprintf("%p", cloneStack.cards) {
		t.Error("Clone failed: Both stacks share the same memory address.")
	}
}

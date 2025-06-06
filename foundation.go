package solitaire

// Foundation - The final place for cards. Cards are built up in piles from the
// (nominated) base card.
// There are at least 4 Foundations per game, sometimes more (there will always
// be a multiple of 4 foundations though).
type Foundation struct {
	Stack *Stack
	Base  SuitedCard
	Rule  func(SuitedCard) bool
}

// CreateFoundations - Create the foundations that will host the cards.
func CreateFoundations(number int, base Card, rule func(Foundation, SuitedCard) bool) []Foundation {
	if number < 1 {
		panic("Cannot have zero foundations")
	}

	if number%SuitCount != 0 {
		panic("Number of foundations must be a multiple of the number of suits in a deck")
	}

	if rule == nil {
		panic("Cannot create foundations without a rule.")
	}

	foundations := make([]Foundation, 0, SuitCount*number)

	for i := 0; i < number; i++ {
		foundation := Foundation{
			Base: SuitedCard{Card: base, Suit: Suit(i)},
		}

		stack := NewStack(CardCount)
		stack.Rule = func(card SuitedCard) bool {
			return rule(foundation, card)
		}

		foundation.Rule = func(card SuitedCard) bool {
			return rule(foundation, card)
		}
		foundation.Stack = stack

		foundations = append(foundations, foundation)
	}

	return foundations
}

// Full - the foundation is full.
func (foundation Foundation) Full() bool {
	return foundation.Len() == CardCount
}

// Add - Add a card to the foundation.
func (foundation Foundation) Add(card SuitedCard, visible bool) {
	if foundation.Stack.Rule(card) {
		foundation.Stack.Add(card, visible)
	}
}

// Len - the length of the stack inside the foundation.
func (foundation Foundation) Len() int {
	return foundation.Stack.Len()
}

// Top - the top most card on the stack inside the foundation.
func (foundation Foundation) Top() (SuitedCard, error) {
	return foundation.Stack.Top()
}

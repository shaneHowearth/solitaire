package state

type StackSpec struct {
	BaseCard  SuitedCard
	AddRule   func(*Stack, SuitedCard) bool
	CardCount [2]int // [2]int{number in stack total, number where visible is set to true}
	SkipCards map[SuitedCard]struct{}
	ShowCount int // Number of cards in a stack that should be shown.
}

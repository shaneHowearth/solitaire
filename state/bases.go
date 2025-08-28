package state

type StackSpec struct {
	BaseCard  SuitedCard
	AddRule   func(*Stack, SuitedCard) bool
	CardCount [2]int
	SkipCards map[SuitedCard]struct{}
}

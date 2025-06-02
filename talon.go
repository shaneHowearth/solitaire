package solitaire

import "fmt"

// Talon -
// The remaining stack of cards, typically squared and face-down, that is left
// after the layout has been populated. These cards can be turned over into the
// waste, usually one-by-one, but sometimes in groups of two or three (depending
// on rules), whenever the player wishes. Also stock. Sometimes equated,
// confusingly, to waste pile.
type Talon struct {
	Stock        *Stack
	Waste        *Stack
	DealCount    int // How many times can the stack of cards be dealt.
	PerDealCount int // How many cards are dealt per deal.
}

// Deal - Deal the top card of the Stock to the Waste pile.
func (talon Talon) Deal() {
	// If there are cards left on the Stock pile, shift the top one onto the
	// Waste pile.
	if talon.Stock.Len() > 0 {
		card, err := talon.Stock.Deal()
		if err != nil {
			// Well hot diggity, something's wrong.
			fmt.Println("An error occurred where one shouldn't", err.Error())
		}

		talon.Waste.Add(card, true)

		return
	}

	// If the Count of deals remaining is greater than 1, and there are cards on
	// the Waste pile, put them all into the Stock Pile.
	// Note that if there is only one card on the waste pile that there is no
	// need to re-deal.
	if talon.DealCount > 0 && talon.Waste.Len() > 1 {
		// Move the waste to the stock.
		for {
			wasteCard, err := talon.Waste.Deal()
			if err != nil {
				// Nothing left to do.
				break
			}

			talon.Stock.Add(wasteCard, false)
		}

		talon.DealCount--

		// Deal the first card onto the Waste.
		talon.Deal()
	}
}

// Top - the top card on the waste stack.
func (talon Talon) Top() (SuitedCard, error) {
	return talon.Waste.Top()
}

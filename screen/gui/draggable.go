package gui

import (
	"fyne.io/fyne/v2"
	"github.com/shanehowearth/solitaire/screen"
	"github.com/shanehowearth/solitaire/state"
)

var _ screen.Display = (*DisplayGUI)(nil)

func (d *DisplayGUI) removeCardFromSource(card *CardWidget) {
	// Find and remove the card from its current pile
	// Implementation depends on your data structure
}

func (d *DisplayGUI) getAllPiles() []*PileWidget {
	// Return all piles in your game
	// Implementation depends on your game structure
	return nil
}

// Helper methods for game rules
func (d *DisplayGUI) isSameSuit(card1, card2 string) bool {
	if len(card1) == 0 || len(card2) == 0 {
		return false
	}
	return card1[len(card1)-1] == card2[len(card2)-1]
}

func (d *DisplayGUI) isOppositeColor(card1, card2 string) bool {
	color1 := d.getCardColor(card1)
	color2 := d.getCardColor(card2)
	return color1 != color2
}

func (d *DisplayGUI) getCardColor(card string) string {
	if len(card) == 0 {
		return ""
	}
	// TODO: split this properly
	suit := string(card[len(card)-1])
	if suit == "♥" || suit == "♦" {
		return "red"
	}
	return "black"
}

func (d *DisplayGUI) isNextRank(current, next string) bool {
	// Check if next is the next rank after current (for foundations)
	// Implementation depends on your rank representation
	return false // Placeholder
}

func (d *DisplayGUI) isPreviousRank(current, previous string) bool {
	// Check if previous is the previous rank before current (for tableau)
	// Implementation depends on your rank representation
	return false // Placeholder
}

func (d *DisplayGUI) layoutPile(pile *PileWidget) {
	// Layout cards in the pile with proper positioning
	for i, card := range pile.cards {
		// Stack cards with slight offset for tableau, or exact overlap for other piles
		if pile.stackType == state.StackTableau {
			// Tableau cards should cascade down
			card.Move(fyne.NewPos(
				pile.Position().X,
				pile.Position().Y+float32(i*20), // 20 pixel offset per card
			))
		} else {
			// Other piles stack exactly on top
			card.Move(pile.Position())
		}
	}
	pile.Refresh()
}

package solitaire_test

import (
	"testing"

	"github.com/shanehowearth/solitaire"
	"github.com/stretchr/testify/assert"
)

func Test_Deal(t *testing.T) {
	stockSize := 52
	dealCount := 1
	perDealCount := 3

	testcases := map[string]struct {
		StockCount      int
		FinalStockCount int
		WasteCount      int
		FinalWasteCount int
		Output          bool
		DealCount       int
		Rule            func(solitaire.SuitedCard) bool
	}{
		"Stack has multiple cards, Waste is empty": {
			StockCount:      5,
			FinalStockCount: 4,
			FinalWasteCount: 1,
			DealCount:       1,
			Output:          true,
			Rule:            func(solitaire.SuitedCard) bool { return true },
		},
		"Stack and waste both have multiple cards": {
			StockCount:      5,
			FinalStockCount: 4,
			WasteCount:      5,
			FinalWasteCount: 6,
			DealCount:       1,
			Output:          true,
			Rule:            func(solitaire.SuitedCard) bool { return true },
		},
		"Stock is empty and Waste has multiple cards": {
			StockCount:      0,
			FinalStockCount: 4,
			WasteCount:      5,
			FinalWasteCount: 1,
			DealCount:       1,
			Output:          true,
			Rule:            func(solitaire.SuitedCard) bool { return true },
		},
		"Stock is empty and Waste only has one card": {
			StockCount:      0,
			FinalStockCount: 0,
			WasteCount:      1,
			FinalWasteCount: 1,
			DealCount:       1,
			Output:          false,
			Rule:            func(solitaire.SuitedCard) bool { return true },
		},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			standardDeck := solitaire.CreateDecks(1)
			talon := solitaire.NewTalon(stockSize, dealCount, perDealCount, testcase.Rule)

			// Add cards to the talon stock pile.
			for sc := 0; sc < testcase.StockCount; sc++ {
				card := standardDeck.Deal()

				talon.Stock.Add(card, false)
			}

			// Add cards to the talon waste pile.
			for sc := 0; sc < testcase.WasteCount; sc++ {
				card := standardDeck.Deal()

				talon.Waste.Add(card, false)
			}

			for dc := 1; dc < testcase.DealCount; dc++ {
				// No-op, purely to use up dealcount
				_ = talon.Deal()
			}

			output := talon.Deal()

			assert.Equalf(t, testcase.Output, output,
				"Unexpected output got %t, want %t", output, testcase.Output)

			assert.Equalf(t, testcase.FinalStockCount, talon.Stock.Len(),
				"Unexpected final stock count got %d, want %d",
				talon.Stock.Len(), testcase.FinalStockCount,
			)

			assert.Equalf(t, testcase.FinalWasteCount, talon.Waste.Len(),
				"Unexpected final waste count got %d, want %d",
				talon.Waste.Len(), testcase.FinalWasteCount,
			)
		})
	}
}

package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func Test_Deal(t *testing.T) {
	// A standard rule factory for the waste pile
	ruleFactory := func(s *state.Stack) func(state.SuitedCard) bool {
		return func(state.SuitedCard) bool { return true }
	}

	testcases := map[string]struct {
		StockCount      int
		WasteCount      int
		InitialDeals    int
		ExpectedOutput  bool
		FinalStockCount int
		FinalWasteCount int
		FinalDeals      int
	}{
		"Stock has cards, moves one to Waste": {
			StockCount:      5,
			WasteCount:      0,
			InitialDeals:    1,
			ExpectedOutput:  true,
			FinalStockCount: 4,
			FinalWasteCount: 1,
			FinalDeals:      1, // No recycle yet
		},
		"Stock empty, Waste recycles to Stock": {
			StockCount:      0,
			WasteCount:      5,
			InitialDeals:    1,
			ExpectedOutput:  true,
			FinalStockCount: 4, // 5 cards recycled, then 1 dealt immediately
			FinalWasteCount: 1,
			FinalDeals:      0, // DealCount decremented
		},
		"Everything empty returns false": {
			StockCount:      0,
			WasteCount:      0,
			InitialDeals:    1,
			ExpectedOutput:  false,
			FinalStockCount: 0,
			FinalWasteCount: 0,
			FinalDeals:      1,
		},
	}

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			talon := state.NewTalon(tc.InitialDeals, 1, ruleFactory)
			deck := state.CreateDecks(1)

			// Setup state
			for i := 0; i < tc.StockCount; i++ {
				talon.Stock.Add(deck.Deal(), false)
			}
			for i := 0; i < tc.WasteCount; i++ {
				talon.Waste.Add(deck.Deal(), false)
			}

			result := talon.Deal()

			assert.Equal(t, tc.ExpectedOutput, result)
			assert.Equal(t, tc.FinalStockCount, talon.Stock.Len(), "Stock length mismatch")
			assert.Equal(t, tc.FinalWasteCount, talon.Waste.Len(), "Waste length mismatch")
			assert.Equal(t, tc.FinalDeals, talon.DealCount, "Remaining deals mismatch")
		})
	}
}

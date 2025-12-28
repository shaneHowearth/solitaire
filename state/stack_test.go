package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func Test_Add(t *testing.T) {
	shuffledDeck := state.CreateDecks(1)
	shuffledDeck.Shuffle()

	testcases := map[string]struct {
		Number  int
		Visible bool
	}{
		"Add one":         {Number: 1},
		"Add ten":         {Number: 10},
		"Add ten visible": {Number: 10, Visible: true},
	}
	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			stack := state.NewStack(
				testcase.Number,
				state.SuitedCard{},
				func(*state.Stack) func(state.SuitedCard) bool {
					return func(state.SuitedCard) bool {
						return true
					}
				},
				state.StackUndefined,
			)

			for x := 1; x <= testcase.Number; x++ {
				card := shuffledDeck.Deal()
				card.Visible = testcase.Visible
				stack.Add(card, testcase.Visible)
				assert.Equalf(t, x, stack.Len(), "Stack length error got %d want %d", stack.Len(), x)
				topCard, err := stack.Top()
				assert.Nilf(t, err, "Getting the top card errored %v", err)

				assert.EqualExportedValuesf(t,
					topCard, card,
					"Got different cards got %#v want %#v",
					topCard, card)
			}
		})
	}
}

// func Test_Move(t *testing.T) {.
// 	// Standard rule factory: always allows adding.
// 	allowRule := func(s *state.Stack) func(state.SuitedCard) bool {.
// 		return func(state.SuitedCard) bool { return true }.
// 	}.

// 	// Deny rule factory: never allows adding.
// 	denyRule := func(s *state.Stack) func(state.SuitedCard) bool {.
// 		return func(state.SuitedCard) bool { return false }.
// 	}.

// 	testcases := map[string]struct {.
// 		Number         int.
// 		SourceCount    int.
// 		BuildDest      func() *state.Stack.
// 		ExpectedOutput bool.
// 		ExpectedSource int.
// 		ExpectedDest   int.
// 	}{.
// 		"Should not be able to move a card to the same stack that it came from": {.
// 			Number:         1,.
// 			SourceCount:    5,.
// 			BuildDest:      nil, // Logic: if nil, use source as dest.
// 			ExpectedOutput: false,.
// 			ExpectedSource: 5,.
// 		},.
// 		"Move to an empty stack where the rule allows the move": {.
// 			Number:      1,.
// 			SourceCount: 5,.
// 			BuildDest: func() *state.Stack {.
// 				return state.NewStack(0, state.SuitedCard{}, allowRule, state.StackTableau).
// 			},.
// 			ExpectedOutput: true,.
// 			ExpectedSource: 4,.
// 			ExpectedDest:   1,.
// 		},.
// 		"Move to an empty stack where the rule denies the move": {.
// 			Number:      1,.
// 			SourceCount: 5,.
// 			BuildDest: func() *state.Stack {.
// 				return state.NewStack(0, state.SuitedCard{}, denyRule, state.StackTableau).
// 			},.
// 			ExpectedOutput: false,.
// 			ExpectedSource: 5,.
// 			ExpectedDest:   0,.
// 		},.
// 		"Move multiple cards (e.g. Klondike build)": {.
// 			Number:      3,.
// 			SourceCount: 5,.
// 			BuildDest: func() *state.Stack {.
// 				return state.NewStack(0, state.SuitedCard{}, allowRule, state.StackTableau).
// 			},.
// 			ExpectedOutput: true,.
// 			ExpectedSource: 2,.
// 			ExpectedDest:   3,.
// 		},.
// 	}.

// 	for name, tc := range testcases {.
// 		t.Run(name, func(t *testing.T) {.
// 			deck := state.CreateDecks(1).
// 			source := state.NewStack(0, state.SuitedCard{}, allowRule, state.StackTableau).

// 			for i := 0; i < tc.SourceCount; i++ {.
// 				source.Add(deck.Deal(), true).
// 			}.

// 			var dest *state.Stack.
// 			if tc.BuildDest == nil {.
// 				dest = source.
// 			} else {.
// 				dest = tc.BuildDest().
// 			}.

// 			output := source.Move(dest, tc.Number).

// 			assert.Equal(t, tc.ExpectedOutput, output, "Movement success/failure mismatch").
// 			assert.Equal(t, tc.ExpectedSource, source.Len(), "Source stack count mismatch").
// 			if dest != source && dest != nil {.
// 				assert.Equal(t, tc.ExpectedDest, dest.Len(), "Destination stack count mismatch").
// 			}.
// 		}).
// 	}.
// }.

// func Test_Top(t *testing.T) {.
// 	testcases := map[string]struct {.
// 		stackLen int.
// 		card     state.SuitedCard.
// 		err      error.
// 	}{.
// 		"empty stack":          {stackLen: 0, err: state.ErrEmpty},.
// 		"more than 2 in stack": {stackLen: 2},.
// 		"1 in stack":           {stackLen: 1},.
// 	}.
// 	for name, testcase := range testcases {.
// 		t.Run(name, func(t *testing.T) {.
// 			deck := state.CreateDecks(1).
// 		}).
// 	}.
// }.

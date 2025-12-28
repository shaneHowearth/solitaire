package state_test

import (
	"testing"

	"github.com/shanehowearth/solitaire/state"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTableaus(t *testing.T) {
	// A dummy rule to use for valid test cases
	standardRule := func(t *state.Stack, c state.SuitedCard) bool {
		return true
	}

	testcases := map[string]struct {
		WillPanic    bool
		PanicMessage string
		// BuildSpecs is a helper to generate the input for the specific case
		BuildSpecs func() []state.StackSpec
		Expected   int
	}{
		"Zero tableaus": {
			WillPanic:    true,
			PanicMessage: "Cannot have zero tableaus",
			BuildSpecs: func() []state.StackSpec {
				return []state.StackSpec{}
			},
		},
		"No rule": {
			WillPanic:    true,
			PanicMessage: "Cannot create tableau 0 without a rule.",
			BuildSpecs: func() []state.StackSpec {
				return []state.StackSpec{
					{AddRule: nil}, // One spec, but missing the rule
				}
			},
		},
		"Seven tableaus (klondike)": {
			Expected: 7,
			BuildSpecs: func() []state.StackSpec {
				specs := make([]state.StackSpec, 7)
				for i := 0; i < 7; i++ {
					specs[i] = state.StackSpec{AddRule: standardRule}
				}
				return specs
			},
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			specs := testcase.BuildSpecs()

			if testcase.WillPanic {
				assert.PanicsWithValue(t, testcase.PanicMessage,
					func() {
						state.CreateTableaus(specs)
					},
				)
			} else {
				tableaus := state.CreateTableaus(specs)

				// Check length
				assert.Equal(t, testcase.Expected, len(tableaus),
					"tableau has incorrect number of elements")

				// Optional: Check that the rule was assigned correctly to the underlying stack
				if len(tableaus) > 0 {
					assert.NotNil(t, tableaus[0].Stack.Rule, "Stack rule should be initialized")
				}
			}
		})
	}
}

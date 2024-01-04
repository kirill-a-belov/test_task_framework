package tracer

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStart(t *testing.T) {
	testCaseList := []struct {
		name string
		args func() (input context.Context, wantResult context.Context)
	}{
		{
			name: "Success",
			args: func() (input context.Context, wantResult context.Context) {
				input = context.Background()
				wantResult = input

				return
			},
		},
	}
	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			input, wantResult := tc.args()

			resultCTX, resultSpan := Start(input, "example")
			require.NotPanics(t, func() {
				resultSpan.End()
			})
			assert.Equal(t, wantResult, resultCTX)

		})
	}
}

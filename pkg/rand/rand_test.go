package rand

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRand(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      func() (int, int64)
		wantError bool
	}{
		{
			name: "Success",
			args: func() (int, int64) {
				return 100, 1024
			},
			wantError: false,
		},
		{
			name: "Invalid len",
			args: func() (int, int64) {
				return -100, 1024
			},
			wantError: true,
		},
		{
			name: "Invalid max",
			args: func() (int, int64) {
				return 100, -1024
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			l, m := tc.args()

			result, err := Rand(l, m)

			if tc.wantError {
				assert.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, l, len(result))
		})
	}
}

package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum(t *testing.T) {
	testCaseList := []struct {
		name       string
		args       func() interface{}
		wantResult interface{}
	}{
		{
			name: "Signed",
			args: func() interface{} {
				return Sum([]int{-1, 2, 3}...)
			},
			wantResult: 4,
		},
		{
			name: "Unsigned",
			args: func() interface{} {
				return Sum([]uint{1, 2, 3}...)
			},
			wantResult: uint(6),
		},
		{
			name: "Float",
			args: func() interface{} {
				return Sum([]float64{1.1, 2.2, 3.3}...)
			},
			wantResult: 6.6,
		},
		{
			name: "Complex",
			args: func() interface{} {
				return Sum([]complex128{1 + 1i, 2 + 2i, 3 + 3i}...)
			},
			wantResult: 6 + 6i,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.args()

			assert.Equal(t, tc.wantResult, result)
		})
	}
}

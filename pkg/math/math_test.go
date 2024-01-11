package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/kirill-a-belov/test_task_framework/pkg/rand"
)

func TestSum(t *testing.T) {
	t.Run("Basic cases", func(t *testing.T) {
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
	})

	t.Run("Special cases", func(t *testing.T) {
		t.Run("Args mutation", func(t *testing.T) {
			reference := []int{1, 2, 3}
			input := make([]int, len(reference))
			_ = copy(input, reference)

			result := Sum(input...)
			require.Equal(t, 6, result)
			require.NotEqual(t, reference, input)
			assert.Equal(t, result, input[0])
		})

		t.Run("Zero alloc", func(t *testing.T) {
			res := testing.Benchmark(BenchmarkSum)
			assert.Equal(t, int64(0), res.AllocsPerOp())
		})
	})
}

func BenchmarkSum(b *testing.B) {
	testSet, err := rand.Rand(1024, 1024)
	require.NoError(b, err)

	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = Sum(testSet...)
	}
}

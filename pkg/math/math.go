package math

import "golang.org/x/exp/constraints"

type number interface {
	constraints.Signed | constraints.Unsigned | constraints.Float | constraints.Complex
}

// Sum zero alloc NOTE: args will be mutated
func Sum[N number](args ...N) N {
	return sum(args, 1)
}

func sum[N number](args []N, pos int) N {
	if pos > len(args)-1 {
		return args[0]
	}

	args[0] += args[pos]
	pos++

	return sum(args, pos)
}

package math

import "golang.org/x/exp/constraints"

type number interface {
	constraints.Signed | constraints.Unsigned | constraints.Float | constraints.Complex
}

func Sum[N number](args ...N) N {
	var result N

	for _, item := range args {
		result += item
	}

	return result
}

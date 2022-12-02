package math

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Integer | constraints.Float
}

func Max[N Number](a, b N) N {
	if a > b {
		return a
	}

	return b
}

func Min[N Number](a, b N) N {
	if a < b {
		return a
	}

	return b
}

func Clamp[N Number](v, lo, hi N) N {
	return Min(Max(v, lo), hi)
}

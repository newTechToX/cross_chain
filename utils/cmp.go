package utils

import "golang.org/x/exp/constraints"

func Max[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	max := s[0]
	for _, e := range s {
		if e > max {
			max = e
		}
	}
	return max
}

func Min[T constraints.Ordered](s ...T) T {
	if len(s) == 0 {
		var zero T
		return zero
	}
	min := s[0]
	for _, e := range s {
		if e < min {
			min = e
		}
	}
	return min
}

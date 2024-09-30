// Package ordered contains generic utility functions over ordered types, i.e. types
// that support the operators <, >, == etc.
package ordered

import "golang.org/x/exp/constraints"

// Max returns the largest of one or more ordered values. Note the values can be strings as well
// as numeric types.
func Max[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		panic(ErrEmptySlice)
	}
	max := xs[0]
	for _, n := range xs[1:] {
		if n > max {
			max = n
		}
	}
	return max
}

// Max returns the smallest of one or more ordered values. Note the values can be strings as well
// as numeric types.
func Min[T constraints.Ordered](xs ...T) T {
	if len(xs) == 0 {
		panic(ErrEmptySlice)
	}
	min := xs[0]
	for _, n := range xs[1:] {
		if n < min {
			min = n
		}
	}
	return min
}

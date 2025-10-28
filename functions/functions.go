package functions

import (
	"golang.org/x/exp/constraints"
)

type Numeric interface {
	constraints.Integer | constraints.Float | constraints.Complex
}

// Identity function - returns the value passed.
func Id[T any](v T) T {
	return v
}

// Sum sums the two values passed and returns the result. The values must be of
// the same type and must support the + operator. Such types include all numeric
// types, including complex numbers and string.
func Sum[T Numeric | ~string](a, b T) T {
	return a + b
}

// Product multiplies the two values passed and returns the result. The values must
// be of the same type and must support the * operator. These are the numeric types,
// including complex numbers.
func Product[T Numeric](a, b T) T {
	return a * b
}

type Enum[T any] struct {
	Index int
	Value T
}

// Returns a pointer to a variable whose value is initialized to v.
func Ref[T any](v T) *T {
	return &v
}

// Ternary logic function. If `cond` is true, `v` is returned.
// Otherwise `alt` is returned.
func IfElse[T any](cond bool, v T, alt T) T {
	if cond {
		return v
	} else {
		return alt
	}
}

// Ternary logic function. If `cond` is true, the value obtained by
// evaluating `f` is returned. Otherwise `alt` is evaluated and returned.
func IfElseF[T any](cond bool, f func() T, alt func() T) T {
	if cond {
		return f()
	} else {
		return alt()
	}
}

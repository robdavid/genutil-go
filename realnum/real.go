// Package realnum contains generic utility functions for real numbers, where real numbers
// are any non-complex primitive numeric types spanning signed and unsigned integers and
// floats of all sizes.
package realnum

import (
	"errors"
	"unsafe"

	"golang.org/x/exp/constraints"
)

var ErrUnknownType = errors.New("uknown type")

// Scalar numeric type constraint.
type Real interface {
	constraints.Float | constraints.Integer
}

// Max returns the largest of two ordered values. Note the values can be strings as well
// as numeric types.
func Max[T constraints.Ordered](x, y T) T {
	if y > x {
		return y
	}
	return x
}

// Min returns the smallest of two ordered values. Note the values can be strings as well
// as numeric types.
func Min[T constraints.Ordered](x, y T) T {
	if y < x {
		return y
	}
	return x
}

// Abs returns the absolute value of a non-complex numeric type
func Abs[T Real](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

// IsInteger returns true for instances that are signed or unsigned integers
func IsInteger[T Real](n T) bool {
	var one T = 1
	var two T = 2
	return one/two == 0
}

// Sub returns the difference between two real types, whilst casting
// to a new type.
func Sub[R Real, T Real](x, y T) R {
	return R(x) - R(y)
}

// Sub returns the sum of two real types, whilst casting
// to a (larger) type.
func Add[R Real, T Real](x, y T) R {
	return R(x) + R(y)
}

// Precision returns the number of bits of precision. For
// integers, this is simply the bit size of the integer
// (including the sign bit if present). For floating
func Precision[T Real](v T) int {
	bits := unsafe.Sizeof(v)
	if IsInteger(v) {
		return int(bits)
	} else {
		switch bits {
		case 32:
			return 25 // float32, including implicit leading bit and sign bit.
		case 64:
			return 54 // float64, including implicit leading bit sign bit.
		default:
			panic(ErrUnknownType)
		}
	}
}

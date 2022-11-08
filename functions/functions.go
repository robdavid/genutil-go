package functions

// Identity function - returns the value passed
func Id[T any](v T) T {
	return v
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

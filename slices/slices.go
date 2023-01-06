package slices

// Concatenates a list of list of items into a list of items
func Concat[T any](ss ...[]T) (result []T) {
	cap := Fold(0, ss, func(a int, s []T) int { return a + len(s) })
	result = make([]T, 0, cap)
	for i := range ss {
		result = append(result, ss[i]...)
	}
	return
}

// Returns true if predicate returns true for all elements in
// slice.
func All[T any](slice []T, predicate func(v T) bool) bool {
	for i := range slice {
		if !predicate(slice[i]) {
			return false
		}
	}
	return true
}

// Returns true if predicate returns true for at least one element in
// slice. This is a variation on All in which the predicate function
// takes a pointer to the element.
func AllRef[T any](slice []T, predicate func(v *T) bool) bool {
	for i := range slice {
		if !predicate(&slice[i]) {
			return false
		}
	}
	return true
}

// Returns true if predicate returns true for at least one element in
// slice.
func Any[T any](slice []T, predicate func(v T) bool) bool {
	for i := range slice {
		if predicate(slice[i]) {
			return true
		}
	}
	return false
}

// Returns true if predicate returns true for at least one element in
// slice. This is a variation on Any in which the predicate function
// takes a pointer to the element.
func AnyRef[T any](slice []T, predicate func(v *T) bool) bool {
	for i := range slice {
		if predicate(&slice[i]) {
			return true
		}
	}
	return false
}

// Returns true if slice contains value
func Contains[T comparable](slice []T, value T) bool {
	return Find(slice, value) != -1
}

// Returns the smallest index in slice for which the element equals value, or -1
// none do.
func Find[T comparable](slice []T, value T) int {
	return FindFrom(0, slice, value)
}

// Returns the smallest index in slice, greater than or equal to start,
// for which the element equals value, or -1 if none do.
func FindFrom[T comparable](start int, slice []T, value T) int {
	for i := start; i < len(slice); i++ {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Returns the smallest index in slice for which the element satisfies the predicate,
// or -1 if none do.
func FindUsing[T any](slice []T, predicate func(T) bool) int {
	return FindFromUsing(0, slice, predicate)
}

// Returns the smallest index in slice for which the element satisfies the predicate,
// or -1 if none do.
// This is a variation on FindUsing where the element is passed to the predicate
// by reference.
func FindUsingRef[T any](slice []T, predicate func(*T) bool) int {
	return FindFromUsingRef(0, slice, predicate)
}

// Returns the first index in slice greater than or equal to start,
// for which the element satisfies predicate, or -1 if none do.
func FindFromUsing[T any](start int, slice []T, predicate func(T) bool) int {
	for i := start; i < len(slice); i++ {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the first index in slice greater than or equal to start,
// for which the element satisfies predicate, or -1 if none do.
// This is a variation on FindFromUsing in which each element is passed
// to the predicate by reference.
func FindFromUsingRef[T any](start int, slice []T, predicate func(*T) bool) int {
	for i := start; i < len(slice); i++ {
		if predicate(&slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the largest index, less than or equal to start
// for which the element in slice equals value, or -1 if none do.
func RFindFrom[T comparable](start int, slice []T, value T) int {
	for i := start; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return RFindFromUsingRef(start, slice, func(v *T) bool { return *v == value })
}

// Returns the largest index of slice for which the element equals value,
// or -1 if none do.
func RFind[T comparable](slice []T, value T) int {
	return RFindFrom(len(slice)-1, slice, value)
}

// Returns the largest index in slice, less than or equal to start,
// where the element satisfies predicate, or -1 if none do.
func RFindFromUsing[T any](start int, slice []T, predicate func(T) bool) int {
	for i := start; i >= 0; i-- {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the largest index in slice, less than or equal to start,
// for which the element satisfies predicate, or -1 if none do.
// This is a variation of RFindFromUsing where each element value is passed to the
// predicate by reference.
func RFindFromUsingRef[T any](start int, slice []T, predicate func(*T) bool) int {
	for i := start; i >= 0; i-- {
		if predicate(&slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the largest index in slice for which the element satisfies predicate,
// or -1 if none do.
func RFindUsing[T any](slice []T, predicate func(T) bool) int {
	return RFindFromUsing(len(slice)-1, slice, predicate)
}

// Returns the largest index in slice for which the element satisfies predicate,
// or -1 if none do.
// This is a variation of RFindUsing where each element value is passed to the
// predicate by reference.
func RFindUsingRef[T any](slice []T, predicate func(*T) bool) int {
	return RFindFromUsingRef(len(slice)-1, slice, predicate)
}

// Generates sliceOut from sliceIn, by applying function f to each element of
// sliceIn.
func Map[T any, U any](sliceIn []T, f func(T) U) (sliceOut []U) {
	sliceOut = make([]U, len(sliceIn))
	for i := range sliceIn {
		sliceOut[i] = f(sliceIn[i])
	}
	return
}

// Generates sliceOut from sliceIn, by applying function f to
// the address of each element of sliceIn.
func MapRef[T any, U any](sliceIn []T, f func(*T) U) (sliceOut []U) {
	sliceOut = make([]U, len(sliceIn))
	for i := range sliceIn {
		sliceOut[i] = f(&sliceIn[i])
	}
	return
}

// Applies a function f to an accumulator, with initial value
// a, and a slice element, returning a new accumulator, for each element
// in the slice s. The final accumulator value is returned.
func Fold[A any, T any](a A, s []T, f func(A, T) A) A {
	for i := range s {
		a = f(a, s[i])
	}
	return a
}

// Applies a function f to a reference to an accumulator, with initial value a,
// and a reference to slice element, mutating the accumulator, for every element in the
// slice s. The final value of the accumulator is returned.
func FoldRef[A any, T any](a A, s []T, f func(*A, *T)) A {
	result := a
	for i := range s {
		f(&result, &s[i])
	}
	return result
}

package slices

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

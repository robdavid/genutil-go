package slice

// Returns true if `predicate` returns true for all elements in
// `slice`. This function short circuits and does not run in constant
// time
func All[T any](slice []T, predicate func(v T) bool) bool {
	for i := range slice {
		if !predicate(slice[i]) {
			return false
		}
	}
	return true
}

// Returns true if `predicate` returns true for at least one element in
// `slice`. This function short circuits and does not run in constant
// time. This is a variation on `All` in which the predicate function
// takes a pointer to the element.
func AllRef[T any](slice []T, predicate func(v *T) bool) bool {
	for i := range slice {
		if !predicate(&slice[i]) {
			return false
		}
	}
	return true
}

// Returns true if `predicate` returns true for at least one element in
// `slice`. This function short circuits and does not run in constant
// time
func Any[T any](slice []T, predicate func(v T) bool) bool {
	for i := range slice {
		if predicate(slice[i]) {
			return true
		}
	}
	return false
}

// Returns true if `predicate` returns true for at least one element in
// `slice`. This function short circuits and does not run in constant
// time. This is a variation on `Any` in which the predicate function
// takes a pointer to the element.
func AnyRef[T any](slice []T, predicate func(v *T) bool) bool {
	for i := range slice {
		if predicate(&slice[i]) {
			return true
		}
	}
	return false
}

// Returns true if `slice` contains `value`
func Contains[T comparable](slice []T, value T) bool {
	return Find(slice, value) != -1
}

// Returns the first index from the left in `slice` for which the element equals `value`
func Find[T comparable](slice []T, value T) int {
	return Select(slice, func(v T) bool { return v == value })
}

// Returns the first index from the left in `slice`, greater than or equal to `start`,
// for which the element equals `value`
func FindFrom[T comparable](start int, slice []T, value T) int {
	return SelectFrom(start, slice, func(v T) bool { return v == value })
}

func Select[T any](slice []T, predicate func(T) bool) int {
	return SelectFrom(0, slice, predicate)
}

// Returns the first index from the left, greater than or equal to  `start` in `slice`,
// for which applying the predicate to the value returns true.
func SelectFrom[T any](start int, slice []T, predicate func(T) bool) int {
	for i := start; i < len(slice); i++ {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the first index from the right, less than or equal to  `start` in `slice`,
// for which applying the predicate to the value returns true.
func RSelectFrom[T any](start int, slice []T, predicate func(T) bool) int {
	for i := start; i >= 0; i-- {
		if predicate(slice[i]) {
			return i
		}
	}
	return -1
}

// Returns the first index from the right in `slice` which contains `value`
func RFind[T comparable](slice []T, value T) int {
	return RSelectFrom(len(slice)-1, slice, func(v T) bool { return v == value })
}

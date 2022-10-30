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

// Returns the first index from the left in `slice` which contains `value`
func Find[T comparable](slice []T, value T) int {
	for i := range slice {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Returns the first index from the right in `slice` which contains `value`
func RFind[T comparable](slice []T, value T) int {
	for i := len(slice) - 1; i >= 0; i-- {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

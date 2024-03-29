package slices

import "sort"

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

// Returns true if the two slices provided are the same length
// and contain the same elements in the same order. The nil slice
// and the empty slice are regarded as equivalent and therefore
// equal.
func Equal[T comparable](left []T, right []T) bool {
	var i int
	if len(left) != len(right) {
		return false
	}
	for i = 0; i < len(left); i++ {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}

// Types that have a well defined ordering, comparable with
// `<` and `>` operators.
type OrderComparable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64
}

// Compare two slices of `OrderComparable` elements, most significant item first.
// It will return a value less than 0 if the left slice is smaller than the right
// slice, a value greater than 0 if the left slice is greater than the right slice,
// and 0 if they are equal. The nil slice and the empty slice are regarded as equivalent,
// and therefore equal.
// Examples:
//
//	slices.Compare([]int{1,2},[]int{1,3}) < 0 // true
//	slices.Compare([]int{1,3},[]int{1,3}) == 0 // true
//	slices.Compare([]int{1,3,4},[]int{1,3}) > 0 // true
func Compare[T OrderComparable](left []T, right []T) int {
	var i int
	lenR := len(right)
	lenL := len(left)
	for i = 0; i < lenL && i < lenR; i++ {
		if left[i] < right[i] {
			return -1
		} else if left[i] > right[i] {
			return 1
		}
	}
	if lenL < lenR {
		return -1
	} else if lenL > lenR {
		return 1
	}
	return 0
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

// Applies a mapping function to each element of a slice in place. The mapping function is
// from type T to type T.
func MapI[T any](slice []T, f func(T) T) {
	for i := range slice {
		slice[i] = f(slice[i])
	}
}

// Applies a mapping function to each element of a slice in place. The mapping function is
// from type T to type T, and each element is passed by reference.
func MapRefI[T any](slice []T, f func(*T) T) {
	for i := range slice {
		slice[i] = f(&slice[i])
	}
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

// Accept only the elements of s that satisfy the predicate function f,
// returning a new slice containing those elements.
func Filter[T any](s []T, f func(T) bool) (result []T) {
	result = make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return
}

// Accept only the elements of s that satisfy the predicate function f,
// returning a new slice containing those elements. Similar to Filter()
// except that the elements are passed to the predicate function by
// reference.
func FilterRef[T any](s []T, f func(*T) bool) (result []T) {
	result = make([]T, 0, len(s))
	for i := range s {
		if f(&s[i]) {
			result = append(result, s[i])
		}
	}
	return
}

// A filter function that alters the provided slice in place.
// The slice referenced by s is altered so that it contains only
// elements that satisfy the predicate function f.
//
// eg.
//
//	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	FilterI(&slice, func(i int) bool { return i%2 == 0 })
//	fmt.Printf("%v",slice) // [2 4 6 8]
func FilterI[T any](s *[]T, f func(T) bool) {
	j := 0
	for i := range *s {
		if f((*s)[i]) {
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

// A filter function that alters the provided slice in place.
// The slice referenced by s is altered so that it contains only
// elements that satisfy the predicate function f. Each element is
// passed to f by reference.
//
// eg.
//
//	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
//	FilterRefI(&slice, func(i *int) bool { return (*i)%2 == 0 })
//	fmt.Printf("%v",slice) // [2 4 6 8]
func FilterRefI[T any](s *[]T, f func(*T) bool) {
	j := 0
	for i := range *s {
		if f(&(*s)[i]) {
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

// A type constraint for types that can be compared
// via the < operator
type Sortable interface {
	int | int8 | int16 | int32 | int64 |
		uint | uint8 | uint16 | uint32 | uint64 |
		float32 | float64 |
		string
}

// A wrapper type around a slice that satisfies the
// sort.Interface interface. The element type of the
// slice must satisfy Sortable, meaning that the elements
// must be comparable by the < operator
type SortableSlice[T Sortable] []T

func (ss SortableSlice[T]) Len() int {
	return len(ss)
}

func (ss SortableSlice[T]) Less(i, j int) bool {
	return ss[i] < ss[j]
}

func (ss SortableSlice[T]) Swap(i, j int) {
	tmp := ss[i]
	ss[i] = ss[j]
	ss[j] = tmp
}

// Sorts slice in place
func Sort[T Sortable](slice []T) {
	sort.Sort(SortableSlice[T](slice))
}

// Creates a copy of slice, sorted. The
// slice parameter remains unchanged.
func Sorted[T Sortable](slice []T) []T {
	sorted := make([]T, len(slice))
	copy(sorted, slice)
	sort.Sort(SortableSlice[T](sorted))
	return sorted
}

// A implementation of the `sort.Interface` interface
// based on a slice and a predicate over two elements
// that should return true if the first element should
// be ordered before the second
type SortableByFunction[T any] struct {
	Slice     []T
	Predicate func(T, T) bool
}

func (sbf SortableByFunction[T]) Len() int {
	return len(sbf.Slice)
}

func (sbf SortableByFunction[T]) Swap(i, j int) {
	tmp := sbf.Slice[i]
	sbf.Slice[i] = sbf.Slice[j]
	sbf.Slice[j] = tmp
}

func (sbf SortableByFunction[T]) Less(i, j int) bool {
	return sbf.Predicate(sbf.Slice[i], sbf.Slice[j])
}

// Sorts a slice by using a function to determine ordering. The
// less parameter is a function of two elements of type T that should
// return true if the first is less then (should be ordered before) the second.
// The slice is sorted in place.
func SortUsing[T any](slice []T, less func(T, T) bool) {
	sort.Sort(SortableByFunction[T]{slice, less})
}

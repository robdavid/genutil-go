package slices

import (
	"cmp"
	"errors"
	"fmt"
	"runtime"
	stdslices "slices"
	"sort"
	"sync"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/internal/rangehelper"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/opt"
	"github.com/robdavid/genutil-go/ordered"
)

var ErrInvalidRange = errors.New("invalid range")
var ErrInvalidNumCPU = errors.New("invalid number of CPUs")

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
// The result may have additional unused capacity.
// The result preserves the nilness of s.
// This function is pulled through from the standard library
func Clone[S ~[]E, E any](s S) S {
	return stdslices.Clone(s)
}

// BinarySearch searches for target in a sorted slice and returns the earliest
// position where target is found, or the position where target would appear
// in the sort order; it also returns a bool saying whether the target is
// really found in the slice. The slice must be sorted in increasing order.
// This function is pulled through from the standard library.
func BinarySearch[S ~[]E, E cmp.Ordered](x S, target E) (int, bool) {
	return stdslices.BinarySearch(x, target)
}

// BinarySearchUsing works like [BinarySearch], but uses a custom comparison
// function. The slice must be sorted in increasing order, where "increasing"
// is defined by cmp. cmp should return 0 if the slice element matches
// the target, a negative number if the slice element precedes the target,
// or a positive number if the slice element follows the target.
// cmp must implement the same ordering as the slice, such that if
// cmp(a, t) < 0 and cmp(b, t) >= 0, then a must precede b in the slice.
// This function is pulled through from the standard library and renamed for
// consistency.
func BinarySearchUsing[S ~[]E, E, T any](x S, target T, cmp func(E, T) int) (int, bool) {
	return stdslices.BinarySearchFunc(x, target, cmp)
}

// Concat concatenates a slice of slices of items into a new slice of items
// This function is pulled through from the standard library.
func Concat[S ~[]E, E any](slices ...S) S {
	return stdslices.Concat(slices...)
}

// Affix concatenates a slice with a number of elements, creating a new slice.
// This is similar to [append] but differs in that it always allocates a new
// slice for the result.
func Affix[S ~[]E, E any](slice S, items ...E) S {
	return stdslices.Concat(slice, items)
}

// New creates a new slice with the given elements
func New[T any](items ...T) []T {
	return items
}

// New creates a new slice with the given elements, wrapping
// the slice in explicit type S.
func NewAs[S ~[]T, T any](items ...T) S {
	return S(items)
}

func reverse[S ~[]T, T any](from, to S) {
	len := len(from)
	for i := 0; i < len/2; i++ {
		to[i], to[len-1-i] = from[len-1-i], from[i]
	}
	if len&1 != 0 && &from[0] != &to[0] {
		mid := len / 2
		to[mid] = from[mid]
	}
}

/*
ReverseI reverses the elements of s in place.

	s := []int{1,2,3,4,5}
	slices.Reverse(s)
	s == []int{5,4,3,2,1} // true
*/
func ReverseI[S ~[]T, T any](s S) {
	reverse(s, s)
}

/*
Reverse creates a new slice containing all the elements
of f in reverse order. If f is nil, nil will be returned.

	slices.Reverse([]int{1,2,3,4,5}) // []int{5,4,3,2,1}
*/
func Reverse[S ~[]T, T any](f S) (t S) {
	if f != nil {
		t = make(S, len(f))
		reverse(f, t)
	}
	return
}

func parChunks[T any](slice []T, minPar int, maxCpu int) (slices [][]T) {
	if l := len(slice); l <= minPar {
		return [][]T{slice}
	} else {
		var chunkSize int
		if maxCpu < 1 {
			panic(fmt.Errorf("%w: %d", ErrInvalidNumCPU, maxCpu))
		}
		idealCpu := l / minPar
		if idealCpu <= maxCpu {
			chunkSize = l / idealCpu
		} else {
			chunkSize = l / maxCpu
		}
		if l%chunkSize != 0 {
			chunkSize++
		}
		slices = make([][]T, 0, l/chunkSize)
		for pos := 0; pos < l; {
			npos := pos + chunkSize
			chunk := functions.IfElseF(npos < l,
				func() []T { return slice[pos:npos] },
				func() []T { return slice[pos:] })
			slices = append(slices, chunk)
			pos = npos
		}
		return
	}
}

func sliceFill[S ~[]T, T ordered.Real](start, aStep T, desc bool, slice S) {
	v := start
	if desc {
		for i := range slice {
			slice[i] = v
			v -= aStep
		}
	} else {
		for i := range slice {
			slice[i] = v
			v += aStep
		}
	}
}

func parSliceFill[T ordered.Real](start, aStep T, desc bool, chunks [][]T) {
	cstart := start
	var wg sync.WaitGroup
	for i := range chunks {
		wg.Add(1)
		n := i
		cs := cstart
		go func() {
			defer wg.Done()
			sliceFill(cs, aStep, desc, chunks[n])
		}()
		chunkStep := aStep * T(len(chunks[i]))
		if desc {
			cstart -= chunkStep
		} else {
			cstart += chunkStep
		}
	}
	wg.Wait()
}

func rangeBy[T ordered.Real, S ordered.Real](start, end T, step S, inclusive bool, parThreshold int, maxCpu int) (result []T) {
	if start == end {
		if inclusive {
			return []T{start}
		} else {
			return []T{}
		}
	}
	if T(step) == 0 {
		panic(fmt.Errorf("%w: step is zero", ErrInvalidRange))
	}
	if (step > 0 && end < start) || (step < 0 && end > start) {
		panic(fmt.Errorf("%w: negative step or inverse range (but not both)", ErrInvalidRange))
	}
	size, aStep := rangehelper.RangeSize(start, end, step, inclusive)
	result = make([]T, size)
	if parThreshold < 1 {
		sliceFill(start, aStep, step < 0, result)
	} else {
		chunks := parChunks(result, parThreshold, maxCpu)
		if len(chunks) == 1 {
			sliceFill(start, aStep, step < 0, result)
		} else {
			parSliceFill(start, aStep, step < 0, chunks)
		}
	}
	return
}

// RangeBy generates a slice consisting of a sequence of real numbers. The
// range starts at start and extends up to, but does not include, end.
// The difference between each number will be determined by step.
//
// e.g.
//
//	slices.RangeBy(0, 5, 1) // []int{0, 1, 2, 3, 4}
//
// If the range is to be in descending order, the step should be negative
// and start should be larger than end.
//
//	slices.RangeBy[uint](6, 0, -2) // []uint{6, 4, 2}
//
// If start is larger than end whilst step is positive, or if end is larger
// than start whilst step is negative, the function panics. If start is equal to
// end then an empty slice is returned.
func RangeBy[T ordered.Real, S ordered.Real](start, end T, step S) (result []T) {
	return rangeBy[T, S](start, end, step, false, 0, 0)
}

// IncRangeBy generates a slice consisting of a sequence of real numbers. The
// range starts at start and extends up to, and includes, end.
// The difference between each number will be determined by step.
//
// e.g.
//
//	slices.IncRangeBy(0, 5, 1) // []int{0, 1, 2, 3, 4, 5}
//
// If the range is to be in descending order, the step should be negative
// and start should be larger than end.
//
//	slices.IncRangeBy[uint](6, 0, -2) // []uint{6, 4, 2, 0}
//
// If start is larger than end whilst step is positive, or if end is larger
// than start whilst step is negative, the function panics. If start is equal to
// end then a slice containing that value as its only element is returned.
func IncRangeBy[T ordered.Real, S ordered.Real](start, end T, step S) (result []T) {
	return rangeBy[T, S](start, end, step, true, 0, 0)
}

// Range generates a slice consisting of a sequence of real numbers. The
// sequence begins at start and extends up to, but does not include, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 if end is less than start. If start and end are equal, an empty slice
// is returned.
//
// e.g.
//
//	slices.Range(0, 5)         // []int{0, 1, 2, 3, 4}
//	slices.RangeBy[uint](5, 0) // []uint{5, 4, 3, 2, 1}
func Range[T ordered.Real](start, end T) []T {
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), false, 0, 0)
}

// IncRange generates a slice consisting of a sequence of real numbers. The
// sequence begins at start and extends up to, and includes, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 if end is less than start. If start and end are equal, a slice
// containing that value as its only element is returned.
//
// e.g.
//
//	slices.IncRange(0, 5)       // []int{0, 1, 2, 3, 4, 5}
//	slices.IncRange[uint](5, 0) // []uint{5, 4, 3, 2, 1, 0}
func IncRange[T ordered.Real](start, end T) []T {
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), true, 0, 0)
}

type parOptions struct {
	threshold int
	maxCpu    int
}

type ParOption func(*parOptions)

func combineParOptions(opts []ParOption) parOptions {
	result := parOptions{100000, runtime.NumCPU()}
	for _, opt := range opts {
		opt(&result)
	}
	return result
}

// ParThreshold is an option that defines the minimum size of a slice beyond which multiple goroutine
// threads will be used to populate a slice in a Par*Range* function. Defaults to 100000.
func ParThreshold(threshold int) ParOption { return func(o *parOptions) { o.threshold = threshold } }

// ParMaxCpu is the maximum number of goroutine threads to use to populate a slice range in parallel
// in a Par*Range* function. Defaults to `runtime.MaxCPU()`.
func ParMaxCpu(maxCpu int) ParOption { return func(o *parOptions) { o.maxCpu = maxCpu } }

// ParRange generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The sequence begins at start and extends up to, but does not include, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 if end is less than start. If start and end are equal, an empty slice
// is returned.
//
// e.g.
//
//	slices.ParRange(0, 400000)    // []int{0, 1, 2, 3, 4, ..., 399999}
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter is used to control this threshold and the maximum number of goroutines
// used. The number of goroutines ought not exceed the number of logical CPU cores available.
//
// e.g.
//
//	slices.ParRange(0, 400000, ParThreshold(100000), ParMaxCpu(4))
func ParRange[T ordered.Real](start, end T, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), false, opts.threshold, opts.maxCpu)
}

// ParIncRange generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The sequence begins at start and extends up to, and includes, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 if end less than start. If start and end are equal, a slice
// containing that value as its only element is returned.
//
// e.g.
//
//	slices.ParIncRange(0, 400000)    // []int{0, 1, 2, 3, 4, ..., 400000}
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter is used to control this threshold and the maximum number of goroutines
// used. The number of goroutines ought not exceed the number of logical CPU cores available.
//
// e.g.
//
//	slices.ParIncRange(0, 400000, ParThreshold(100000), ParMaxCpu(4))
func ParIncRange[T ordered.Real](start, end T, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), true, opts.threshold, opts.maxCpu)
}

// ParRangeBy generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The sequence begins at start and extends up to, but does not include, end.
// The difference between each number will be determined by step.
//
// e.g.
//
//	slices.ParRangeBy(0, 400000, 1) // []int{0, 1, 2, 3, ..., 399998, 399999}
//
// If the range is to be in descending order, the step should be negative
// and start should be larger than end.
//
//	slices.ParRangeBy[uint](400000, 0, -2) // []uint{400000, 399998, 399996, ..., 4, 2}
//
// If start is larger than end whilst step is positive, or if end is larger
// than start whilst step is negative, the function panics. If start and end are equal,
// an empty slice is returned.
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter can be used to control this threshold and the maximum number of goroutines
// used. The number of goroutines ought not exceed the number of logical CPU cores available.
//
// e.g.
//
//	slices.ParRangeBy(0, 400000, 1, ParThreshold(100000), ParMaxCpu(4))
func ParRangeBy[T ordered.Real, S ordered.Real](start, end T, step S, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, step, false, opts.threshold, opts.maxCpu)
}

// ParIncRangeBy generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The range begins at start and extends up to, and including, end.
// The difference between each number will be determined by step.
//
// e.g.
//
//	slices.ParIncRangeBy(0, 400000, 1) // []int{0, 1, 2, 3, ..., 399999, 400000}
//
// If the range is to be in descending order, the step should be negative
// and start should be larger than end.
//
//	slices.ParIncRangeBy[uint](400000, 0, -2) // []uint{400000, 399998, 399996, ..., 2, 0}
//
// If start is larger than end whilst step is positive, or if end is larger
// than start whilst step is negative, the function panics. If start and end are equal,
// a slice containing that value as its only element is returned.
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter can be used to control this threshold and the maximum number of goroutines
// used. The number of goroutines ought not exceed the number of logical CPU cores available.
//
// e.g.
//
//	slices.ParIncRangeBy(0, 400000, 1, ParThreshold(100000), ParMaxCpu(4))
func ParIncRangeBy[T ordered.Real, S ordered.Real](start, end T, step S, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, step, true, opts.threshold, opts.maxCpu)
}

// Returns true if predicate returns true for all elements in
// slice.
func All[S ~[]T, T any](slice S, predicate func(v T) bool) bool {
	for i := range slice {
		if !predicate(slice[i]) {
			return false
		}
	}
	return true
}

// Returns true if predicate returns true for all the elements in
// slice. This is a variation on All in which the predicate function
// takes a pointer to the element.
func AllRef[S ~[]T, T any](slice S, predicate func(v *T) bool) bool {
	for i := range slice {
		if !predicate(&slice[i]) {
			return false
		}
	}
	return true
}

// Generate a function equivalent to indexing slice, mapping indexes to values.
func AsFunc[S ~[]T, T any](slice S) func(int) T {
	return func(i int) (v T) {
		return slice[i]
	}
}

// Returns true if predicate returns true for at least one element in
// slice.
func Any[S ~[]T, T any](slice S, predicate func(v T) bool) bool {
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
func AnyRef[S ~[]T, T any](slice S, predicate func(v *T) bool) bool {
	for i := range slice {
		if predicate(&slice[i]) {
			return true
		}
	}
	return false
}

// Returns true if slice contains value
func Contains[S ~[]T, T comparable](slice S, value T) bool {
	return Find(slice, value) != -1
}

// Returns true if the two slices provided are the same length
// and contain the same elements in the same order. The nil slice
// and the empty slice are regarded as equivalent and therefore
// equal.
func Equal[S ~[]T, T comparable](left S, right S) bool {
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

// Compare two slices of `OrderComparable` elements, most significant item first.
// It will return a value less than 0 if the left slice is smaller than the right
// slice, a value greater than 0 if the left slice is greater than the right slice,
// and 0 if they are equal. The nil slice and the empty slice are regarded as equivalent,
// and therefore equal.
func Compare[S ~[]T, T cmp.Ordered](left S, right S) int {
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
func Find[S ~[]T, T comparable](slice S, value T) int {
	return FindFrom(0, slice, value)
}

// Returns the smallest index in slice, greater than or equal to start,
// for which the element equals value, or -1 if none do.
func FindFrom[S ~[]T, T comparable](start int, slice S, value T) int {
	for i := start; i < len(slice); i++ {
		if slice[i] == value {
			return i
		}
	}
	return -1
}

// Returns the smallest index in slice for which the element satisfies the predicate,
// or -1 if none do.
func FindUsing[S ~[]T, T any](slice S, predicate func(T) bool) int {
	return FindFromUsing(0, slice, predicate)
}

// Returns the smallest index in slice for which the element satisfies the predicate,
// or -1 if none do.
// This is a variation on FindUsing where the element is passed to the predicate
// by reference.
func FindUsingRef[S ~[]T, T any](slice S, predicate func(*T) bool) int {
	return FindFromUsingRef(0, slice, predicate)
}

// Returns the first index in slice greater than or equal to start,
// for which the element satisfies predicate, or -1 if none do.
func FindFromUsing[S ~[]T, T any](start int, slice S, predicate func(T) bool) int {
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
func Map[S ~[]T, T any, U any](sliceIn S, f func(T) U) (sliceOut []U) {
	sliceOut = make([]U, len(sliceIn))
	for i := range sliceIn {
		sliceOut[i] = f(sliceIn[i])
	}
	return
}

func MapAs[SO ~[]U, S ~[]T, T, U any](sliceIn S, f func(T) U) SO {
	return SO(Map(sliceIn, f))
}

// Generates sliceOut from sliceIn, by applying function f to
// the address of each element of sliceIn.
func MapRef[S ~[]T, T any, U any](sliceIn S, f func(*T) U) (sliceOut []U) {
	sliceOut = make([]U, len(sliceIn))
	for i := range sliceIn {
		sliceOut[i] = f(&sliceIn[i])
	}
	return
}

// Applies a mapping function to each element of a slice in place. The mapping function is
// from type T to type T.
func MapI[S ~[]T, T any](slice S, f func(T) T) {
	for i := range slice {
		slice[i] = f(slice[i])
	}
}

// Applies a mapping function to each element of a slice in place. The mapping function is
// from type T to type T, and each element is passed by reference.
func MapRefI[S ~[]T, T any](slice S, f func(*T) T) {
	for i := range slice {
		slice[i] = f(&slice[i])
	}
}

// Applies a function f to an accumulator, with initial value
// a, and a slice element, returning a new accumulator, for each element
// in the slice s. The final accumulator value is returned.
func Fold[S ~[]T, A any, T any](s S, a A, f func(A, T) A) A {
	for i := range s {
		a = f(a, s[i])
	}
	return a
}

// Applies a function f to a reference to an accumulator, with initial value a,
// and a reference to slice element, mutating the accumulator, for every element in the
// slice s. The final value of the accumulator is returned.
func FoldRef[S ~[]T, A any, T any](s S, a A, f func(*A, *T)) A {
	result := a
	for i := range s {
		f(&result, &s[i])
	}
	return result
}

// Accept only the elements of s that satisfy the predicate function f,
// returning a new slice containing those elements.
func Filter[S ~[]T, T any](s S, f func(T) bool) (result S) {
	result = make(S, 0, len(s))
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
func FilterRef[S ~[]T, T any](s S, f func(*T) bool) (result S) {
	result = make(S, 0, len(s))
	for i := range s {
		if f(&s[i]) {
			result = append(result, s[i])
		}
	}
	return
}

// FilterMap Combines Filter and Map functionality. The mapping function f is applied
// to each element in the slice of type T, and returns a value of type U and a boolean. If the
// boolean is true, the value of type U is appended to the result slice. Otherwise,
// the element is skipped.
func FilterMap[S ~[]T, T any, U any](s S, f func(T) (U, bool)) (result []U) {
	result = make([]U, 0, len(s))
	for _, v := range s {
		if o, ok := f(v); ok {
			result = append(result, o)
		}
	}
	return
}

// FilterMapOpt Combines Filter and Map functionality. The mapping function f is applied
// to each element in the slice of type T, and returns an Option value of type U. If the
// option is non-empty, the value of type U is appended to the result slice. Otherwise,
// the element is skipped.
func FilterMapOpt[S ~[]T, T any, U any](s S, f func(T) opt.Opt[U]) (result []U) {
	result = make([]U, 0, len(s))
	for _, v := range s {
		if o := f(v); o.HasValue() {
			result = append(result, o.Get())
		}
	}
	return
}

// FilterMapRef combines Filter and Map functionality. The mapping function f is applied to a reference
// to each element in the slice of type T, and returns an Option value of type U. If the option
// is non-empty, the value of type U is appended to the result slice. Otherwise,
// the element is skipped.
func FilterMapRef[S ~[]T, T any, U any](s S, f func(*T) (U, bool)) (result []U) {
	result = make([]U, 0, len(s))
	for i := range s {
		if o, ok := f(&s[i]); ok {
			result = append(result, o)
		}
	}
	return
}

// FilterMapRefOpt combines Filter and Map functionality. The mapping function f is applied to a reference
// to each element in the slice of type T, and returns an Option value of type U. If the option
// is non-empty, the value of type U is appended to the result slice. Otherwise,
// the element is skipped.
func FilterMapRefOpt[S ~[]T, T any, U any](s S, f func(*T) opt.Opt[U]) (result []U) {
	result = make([]U, 0, len(s))
	for i := range s {
		if o := f(&s[i]); o.HasValue() {
			result = append(result, o.Get())
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
func FilterI[S ~[]T, T any](s *S, f func(T) bool) {
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
func FilterRefI[S ~[]T, T any](s *S, f func(*T) bool) {
	j := 0
	for i := range *s {
		if f(&(*s)[i]) {
			(*s)[j] = (*s)[i]
			j++
		}
	}
	*s = (*s)[:j]
}

// A wrapper type around a slice that satisfies the
// sort.Interface interface. The element type of the
// slice must satisfy [cmp.Ordered], meaning that the elements
// must be comparable by the < operator
type SortableSlice[T cmp.Ordered] []T

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
func Sort[T cmp.Ordered](slice []T) {
	sort.Sort(SortableSlice[T](slice))
}

// Creates a copy of slice, sorted. The
// slice parameter remains unchanged.
func Sorted[T cmp.Ordered](slice []T) []T {
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
// return true if the first is less than (should be ordered before) the second.
// The slice is sorted in place.
func SortUsing[T any](slice []T, less func(T, T) bool) {
	sort.Sort(SortableByFunction[T]{slice, less})
}

// Fill fills an existing slice with a specified value
func Fill[S ~[]T, T any](slice S, value T) {
	for i := range slice {
		slice[i] = value
	}
}

// Filled returns a new slice, of the given size, filled with a
// given value.
func Filled[T any](size int, value T) []T {
	slice := make([]T, size)
	Fill(slice, value)
	return slice
}

// Iter makes an [iterator.Iterator][T] from slice []T, containing all the elements
// from the slice in order.
func Iter[T any](slice []T) iterator.Iterator[T] {
	return iterator.NewDefaultIterator(iterator.NewSliceCoreIterator(&slice))
}

// IterRef makes an [iterator.Iterator][*T] from slice []T, containing references to all the elements
// in the slice in order.
func IterRef[T any](slice []T) iterator.Iterator[*T] {
	return iterator.NewDefaultIterator(iterator.NewSliceCoreIteratorRef(&slice))
}

// IterMut makes an [iterator.MutableIterator][T] from slice reference *[]T, containing all the elements
// from the slice in order.
func IterMut[T any](slice *[]T) iterator.MutableIterator[T] {
	return iterator.NewDefaultMutableIterator(iterator.NewSliceCoreIterator(slice))
}

// IterRefMut makes a MutableIterator[*T] from slice reference *[]T, containing references to all the elements
// in the slice in order.
func IterRefMut[T any](slice []T) iterator.MutableIterator[*T] {
	return iterator.NewDefaultMutableIterator(iterator.NewSliceCoreIteratorRef(&slice))
}

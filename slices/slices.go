package slices

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"sync"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/internal/rangehelper"
	"github.com/robdavid/genutil-go/realnum"
)

var ErrInvalidRange = errors.New("invalid range")
var ErrInvalidNumCPU = errors.New("invalid number of CPUs")

// Concatenates a list of list of items into a list of items
func Concat[T any](ss ...[]T) (result []T) {
	cap := Fold(0, ss, func(a int, s []T) int { return a + len(s) })
	result = make([]T, 0, cap)
	for i := range ss {
		result = append(result, ss[i]...)
	}
	return
}

func reverse[T any](from, to []T) {
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
func ReverseI[T any](s []T) {
	reverse(s, s)
}

/*
Reverse creates a new slice containing all the elements
of f in reverse order. If f is nil, nil will be returned.

	slices.Reverse([]int{1,2,3,4,5}) // []int{5,4,3,2,1}
*/
func Reverse[T any](f []T) (t []T) {
	if f != nil {
		t = make([]T, len(f))
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

func sliceFill[T realnum.Real](start, aStep T, desc bool, slice []T) {
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

func parSliceFill[T realnum.Real](start, aStep T, desc bool, chunks [][]T) {
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

func rangeBy[T realnum.Real, S realnum.Real](start, end T, step S, inclusive bool, parThreshold int, maxCpu int) (result []T) {
	if start == end && !inclusive {
		return []T{}
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

// RangeBy generates a slice consisting of a range of a sequence of numbers. The
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
// than start whilst step is negative, the function panics.
func RangeBy[T realnum.Real, S realnum.Real](start, end T, step S) (result []T) {
	return rangeBy[T, S](start, end, step, false, 0, 0)
}

// IncRangeBy generates a slice consisting of a sequence of numbers. The
// range starts at start and extends up to, and including, end.
// The difference between each number will be determined by step.
//
// e.g.
//
//	slices.IncRangeBy(0, 5, 1) // []int{0, 1, 2, 3, 4, 5}
//
// If the range is to be in descending order, the step should be negative
// and start should be larger than end.
//
//	slices.IncRangeBy[uint](4, 0, -2) // []uint{4, 2, 0}
//
// If start is larger than end whilst step is positive, or if end is larger
// than start whilst step is negative, the function panics.
func IncRangeBy[T realnum.Real, S realnum.Real](start, end T, step S) (result []T) {
	return rangeBy[T, S](start, end, step, true, 0, 0)
}

// Range generates a slice consisting of a sequence of real numbers. The
// sequence begins at start and extends up to, but does not include, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 otherwise.
//
// e.g.
//
//	slices.Range(0, 5)         // []int{0, 1, 2, 3, 4}
//	slices.RangeBy[uint](5, 0) // []uint{5, 4, 3, 2, 1}
func Range[T realnum.Real](start, end T) []T {
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), false, 0, 0)
}

// Range generates a slice consisting of a sequence of real numbers. The
// sequence begins at start and extends up to, and including, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 otherwise.
//
// e.g.
//
//	slices.IncRange(0, 5)       // []int{0, 1, 2, 3, 4, 5}
//	slices.IncRange[uint](5, 0) // []uint{5, 4, 3, 2, 1, 0}
func IncRange[T realnum.Real](start, end T) []T {
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
// threads will be used to populate a slice. Defaults to 100000.
func ParThreshold(threshold int) ParOption { return func(o *parOptions) { o.threshold = threshold } }

// ParMaxCpu is the maximum number of goroutine threads to use to populate a slice range in parallel.
// This should not exceed the number processor cores available. Defaults to `runtime.MaxCPU()`.
func ParMaxCpu(maxCpu int) ParOption { return func(o *parOptions) { o.maxCpu = maxCpu } }

// ParRange generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The sequence begins at start and extends up to, but does not include, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 otherwise.
//
// e.g.
//
//	slices.ParRange(0, 400000)    // []int{0, 1, 2, 3, 4, ..., 399999}
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter is used to control this threshold and the maximum number of goroutines
// used, which should not exceed the number of CPU cores available.
//
// e.g.
//
//	slices.ParRange(0, 400000, ParThreshold(100000), ParMaxCpu(4))
func ParRange[T realnum.Real](start, end T, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), false, opts.threshold, opts.maxCpu)
}

// ParIncRange generates a slice consisting of a sequence of real numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.
//
// The sequence begins at start and extends up to, and includes, end.
// Consecutive numbers in the sequence differ by 1 if end is greater than start,
// and by -1 otherwise.
//
// e.g.
//
//	slices.ParIncRange(0, 400000)    // []int{0, 1, 2, 3, 4, ..., 400000}
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter is used to control this threshold and the maximum number of goroutines
// used, which should not exceed the number of CPU cores available.
//
// e.g.
//
//	slices.ParIncRange(0, 400000, ParThreshold(100000), ParMaxCpu(4))
func ParIncRange[T realnum.Real](start, end T, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, functions.IfElse(end < start, -1, 1), true, opts.threshold, opts.maxCpu)
}

// ParRangeBy generates a slice consisting of a range of a sequence of numbers, potentially using
// multiple parallel go routines to accelerate the process on multi core systems.

// The range starts at start and extends up to, but does not include, end.
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
// than start whilst step is negative, the function panics.
//
// If the number of requested elements exceeds a given threshold - by default 100,000 - multiple
// goroutines will be launched in parallel, each tasked with filling a different part of the slice.
// The parOpts variadic parameter is used to control this threshold and the maximum number of goroutines
// used, which should not exceed the number of CPU cores available.
//
// e.g.
//
//	slices.ParRangeBy(0, 400000, 1, ParThreshold(100000), ParMaxCpu(4))

func ParRangeBy[T realnum.Real, S realnum.Real](start, end T, step S, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, step, false, opts.threshold, opts.maxCpu)
}

func ParIncRangeBy[T realnum.Real, S realnum.Real](start, end T, step S, parOpts ...ParOption) []T {
	opts := combineParOptions(parOpts)
	return rangeBy(start, end, step, true, opts.threshold, opts.maxCpu)
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

// Returns true if predicate returns true for all the elements in
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

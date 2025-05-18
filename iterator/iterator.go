// Iterators and generators
package iterator

import (
	"errors"
	"fmt"
	"iter"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/internal/rangehelper"
	"github.com/robdavid/genutil-go/option"
	"github.com/robdavid/genutil-go/ordered"
	"github.com/robdavid/genutil-go/tuple"
)

// The largest slice capacity we are prepared to allocate to collect
// iterators of uncertain size.
const maxUncertainAllocation = 100000

var ErrAllocationSizeInfinite = errors.New("cannot allocate storage for an infinite iterator")
var ErrInvalidIteratorSizeType = errors.New("invalid iterator size type")
var ErrInvalidIteratorRange = errors.New("invalid iterator range")

// SimpleIterator supports a simple sequence of elements
type SimpleIterator[T any] interface {
	// Next sets the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned only when there are no more elements (the current value remains unchanged)
	Next() bool
	// Value gets the current iterator value.
	Value() T
	// Abort stops the iterator; subsequent calls to Next() will return false.
	Abort()
}

// Seq is a generic iter.Seq implementation for any simple iterator
func Seq[T any](i SimpleIterator[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i.Next() {
			if !yield(i.Value()) {
				i.Abort()
				break
			}
		}
	}
}

// Seq2 is a generic iter.Seq2 implementation for any simple iterator returning a 2-tuple
func Seq2[K any, V any](i SimpleIterator[tuple.Tuple2[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i.Next() {
			if !yield(i.Value().First, i.Value().Second) {
				i.Abort()
				break
			}
		}
	}
}

type OptSeq[T any] = option.Option[iter.Seq[T]]

func OptSeqFrom[T any](seq iter.Seq[T]) OptSeq[T] {
	return option.From(seq)
}

// CoreIterator is an extension of SimpleIterator that also holds some additional information
// about what is being iterated over.
type CoreIterator[T any] interface {
	SimpleIterator[T]
	// Seq returns the iterator as a Go iter.Seq function.
	Seq() iter.Seq[T]
	// Size is an estimate, where possible, of the number of elements remaining.
	Size() IteratorSize
	// PreferSeq returns true if the underlying iterator does not have the most efficient implementation
	// of the simple iterator. For example, an iterator based purely on iter.Seq must use iter.Pull
	// to create an implementation of Next(), Value() etc., which carries a performance hit. An iterator
	// of this type would return true for this method.
	PreferSeq() bool
}

type IteratorExtensions[T any] interface {
	// Chan returns the iterator as a channel.
	Chan() <-chan T
	Enumerate() Iterator2[int, T]
}

// Generic iterator
type Iterator[T any] interface {
	CoreIterator[T]
	IteratorExtensions[T]
}

type CoreIterator2[K any, V any] interface {
	CoreIterator[V]
	Seq2() iter.Seq2[K, V]
	Key() K
}
type Iterator2[K any, V any] interface {
	CoreIterator2[K, V]
	IteratorExtensions[V]
}

type SimpleCoreIterator[T any] struct {
	SimpleIterator[T]
	size func() IteratorSize
}

func NewSimpleCoreIterator[T any](itr SimpleIterator[T]) *SimpleCoreIterator[T] {
	return &SimpleCoreIterator[T]{SimpleIterator: itr}
}

func NewSimpleCoreIteratorWithSize[T any](itr SimpleIterator[T], size func() IteratorSize) *SimpleCoreIterator[T] {
	return &SimpleCoreIterator[T]{SimpleIterator: itr, size: size}
}

func (itr *SimpleCoreIterator[T]) Size() IteratorSize {
	if itr.size == nil {
		return NewSizeUnknown()
	} else {
		return itr.size()
	}
}

func (itr *SimpleCoreIterator[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for itr.Next() {
			if !yield(itr.Value()) {
				itr.Abort()
				break
			}
		}
	}
}

func (itr *SimpleCoreIterator[T]) PreferSeq() bool {
	return false
}

type DefaultIterator[T any] struct {
	CoreIterator[T]
}

func NewDefaultIterator[T any](citr CoreIterator[T]) DefaultIterator[T] {
	return DefaultIterator[T]{CoreIterator: citr}
}

func (di DefaultIterator[T]) Chan() <-chan T {
	return Chan(di)
}

func (di DefaultIterator[T]) Enumerate() Iterator2[int, T] {
	return Enumerate(di)
}

type Indexed[T any] = tuple.Tuple2[int, T]

func IndexValue[T any](index int, value T) Indexed[T] {
	return tuple.Of2(index, value)
}

type enumeratedCoreIterator[T any] struct {
	CoreIterator[T]
	key, index int
}

func newEnumeratedCoreIterator[T any](citr CoreIterator[T]) *enumeratedCoreIterator[T] {
	return &enumeratedCoreIterator[T]{CoreIterator: citr, key: 0, index: 0}
}
func (eci *enumeratedCoreIterator[T]) Next() bool {
	if !eci.CoreIterator.Next() {
		return false
	}
	eci.key = eci.index
	eci.index++
	return true
}

func (eci *enumeratedCoreIterator[T]) Key() int {
	return eci.key
}

func (eci *enumeratedCoreIterator[T]) Seq2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for v := range eci.Seq() {
			eci.key = eci.index
			if !yield(eci.key, v) {
				eci.Abort()
				break
			}
			eci.index++
		}
	}
}

type SeqCoreIterator[T any] struct {
	seq   iter.Seq[T]
	size  func() IteratorSize
	stop  func()
	next  func() (T, bool)
	value T
}

func (si *SeqCoreIterator[T]) Seq() iter.Seq[T] {
	return si.seq
}

func (si *SeqCoreIterator[T]) Size() IteratorSize {
	if si.size == nil {
		return NewSizeUnknown()
	} else {
		return si.size()
	}
}

func (si *SeqCoreIterator[T]) PreferSeq() bool {
	return true
}

func (si *SeqCoreIterator[T]) Next() (ok bool) {
	if si.next == nil {
		si.next, si.stop = iter.Pull(si.Seq())
	}
	si.value, ok = si.next()
	return
}

func (si *SeqCoreIterator[T]) Value() T {
	return si.value
}

func (si *SeqCoreIterator[T]) Abort() {
	if si.next == nil {
		si.next, si.stop = iter.Pull(si.Seq())
	}
	si.stop()
}

func NewSeqCoreIterator[T any](seq iter.Seq[T]) *SeqCoreIterator[T] {
	return &SeqCoreIterator[T]{seq: seq}
}

func NewSeqCoreIteratorWithSize[T any](seq iter.Seq[T], size func() IteratorSize) *SeqCoreIterator[T] {
	return &SeqCoreIterator[T]{seq: seq, size: size}
}

func New[T any](seq iter.Seq[T]) Iterator[T] {
	return &DefaultIterator[T]{CoreIterator: NewSeqCoreIterator(seq)}
}

func NewWithSize[T any](seq iter.Seq[T], size func() IteratorSize) Iterator[T] {
	return &DefaultIterator[T]{CoreIterator: NewSeqCoreIteratorWithSize(seq, size)}
}

type SeqCoreIterator2[K any, V any] struct {
	CoreIterator[V]
	seq2 iter.Seq2[K, V]
	key  K
}

func NewSeqCoreIterator2[K any, V any](seq2 iter.Seq2[K, V]) CoreIterator2[K, V] {
	itr2 := SeqCoreIterator2[K, V]{
		seq2: seq2,
	}
	seq := func(yield func(V) bool) {
		for _, v := range itr2.Seq2() {
			if !yield(v) {
				break
			}
		}
	}
	itr2.CoreIterator = NewSeqCoreIterator(seq)
	return &itr2
}

func (si *SeqCoreIterator2[K, V]) Seq2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		var value V
		for si.key, value = range si.seq2 {
			if !yield(si.key, value) {
				break
			}
		}
	}
}

func (si *SeqCoreIterator2[K, V]) Key() K {
	return si.key
}

type DefaultIterator2[K any, V any] struct {
	CoreIterator2[K, V]
	DefaultIterator[V]
}

func NewDefaultIterator2[K any, V any](core CoreIterator2[K, V]) DefaultIterator2[K, V] {
	return DefaultIterator2[K, V]{DefaultIterator: DefaultIterator[V]{CoreIterator: core}, CoreIterator2: core}
}

func New2[K any, V any](seq2 iter.Seq2[K, V]) Iterator2[K, V] {
	return NewDefaultIterator2(NewSeqCoreIterator2(seq2))
}

type IteratorSizeType int

const (
	SizeUnknown IteratorSizeType = iota
	SizeKnown
	SizeAtMost
	SizeInfinite
)

// IteratorSize holds iterator sizing information
type IteratorSize struct {
	Type IteratorSizeType
	Size int
}

func (isz IteratorSize) Allocate() int {
	switch isz.Type {
	case SizeUnknown:
		return 0
	case SizeKnown:
		return isz.Size
	case SizeInfinite:
		panic(ErrAllocationSizeInfinite)
	case SizeAtMost:
		{
			sz := isz.Size / 2
			if sz >= maxUncertainAllocation {
				sz = maxUncertainAllocation
			}
			return sz
		}
	}
	panic(ErrInvalidIteratorSizeType)
}

func (isz IteratorSize) Subset() IteratorSize {
	switch isz.Type {
	case SizeUnknown, SizeInfinite, SizeAtMost:
		return isz
	case SizeKnown:
		return IteratorSize{SizeAtMost, isz.Size}
	}
	panic(ErrInvalidIteratorSizeType)
}

// Iterator sizing information; size is unknown
func NewSizeUnknown() IteratorSize {
	return IteratorSize{Type: SizeUnknown}
}

// IsSizeUnknown returns true if the given IteratorSize instance represents
// an unknown size
func IsSizeUnknown(size IteratorSize) bool {
	return size.Type == SizeUnknown
}

// NewSize creates an `IteratorSize` implementation that has a fixed size of `n`.
func NewSize(n int) IteratorSize { return IteratorSize{SizeKnown, n} }

// IsSizeKnown returns true if the iterator size is one whose actual size is known.
func IsSizeKnown(size IteratorSize) bool {
	return size.Type == SizeKnown
}

// NewSizeAtMost creates an `IteratorSize` implementation that has a size no more than n.
func NewSizeAtMost(n int) IteratorSize {
	return IteratorSize{SizeAtMost, n}
}

// IsSizeAtMost returns true if the iterator size is one whose maximum size is known.
func IsSizeAtMost(size IteratorSize) bool {
	return size.Type == SizeAtMost
}

func NewSizeInfinite() IteratorSize {
	return IteratorSize{SizeInfinite, -1}
}

func IsSizeInfinite(size IteratorSize) bool {
	return size.Type == SizeInfinite
}

// funcNext is a function supporting a transforming operation by consuming
// all or part of an iterator, returning the next value
type funcNext[T any, U any] func(Iterator[T]) (bool, U)

func safeClose[T any](ch chan T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	close(ch)
	return true
}

func safeSend[T any](ch chan<- T, val T) (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	ch <- val
	return true
}

// Generic channel implementation. Produces a channel yielding
// values from the iterator
func Chan[T any](itr CoreIterator[T]) (out chan T) {
	out = make(chan T)
	go func() {
		defer safeClose(out)
		if _, ok := itr.(*SimpleCoreIterator[T]); ok {
			for itr.Next() {
				if !safeSend(out, itr.Value()) {
					break
				}
			}
		} else {
			for v := range itr.Seq() {
				if !safeSend(out, v) {
					itr.Abort()
					break
				}
			}
		}
	}()
	return
}

func Enumerate[T any](itr CoreIterator[T]) Iterator2[int, T] {
	return NewDefaultIterator2(newEnumeratedCoreIterator(itr))
}

// MakeIterator creates a generic iterator from a simple iterator. Provides an implementation
// of additional Iterator methods.
func MakeIterator[T any](base SimpleIterator[T]) Iterator[T] {
	return DefaultIterator[T]{NewSimpleCoreIterator(base)}
}

// mapIter wraps an iterator and adds a mapping function
type mapIter[T, U any] struct {
	base     Iterator[T]
	mapping  funcNext[T, U]
	value    U
	sizeFunc func(IteratorSize) IteratorSize
}

func (i *mapIter[T, U]) Next() bool {
	ok, value := i.mapping(i.base)
	if !ok {
		return false
	} else {
		i.value = value
		return true
	}
}

func (i *mapIter[T, U]) Value() U {
	return i.value
}

func (i *mapIter[T, U]) Abort() {
	i.base.Abort()
}

func (i *mapIter[T, U]) Size() IteratorSize {
	return i.sizeFunc(i.base.Size())
}

func (i *mapIter[T, U]) PreferSeq() bool { return false }

func (i *mapIter[T, U]) Seq() iter.Seq[U] {
	return func(yield func(U) bool) {
		for i.Next() {
			if !yield(i.Value()) {
				break
			}
		}
	}
}

// wrapFunc creates a new iterator from an existing iterator and a function that consumes it, yielding
// one element at a time.
func wrapFunc[T any, U any](iterator Iterator[T], f funcNext[T, U], sizeFunc func(sz IteratorSize) IteratorSize) Iterator[U] {
	return NewDefaultIterator(&mapIter[T, U]{base: iterator, mapping: f, sizeFunc: sizeFunc})
}

// Iterator over a slice
type sliceIter[T any] struct {
	slice []T
	index int
	value T
}

func (si *sliceIter[T]) Next() bool {
	if si.index < len(si.slice) {
		si.value = si.slice[si.index]
		si.index++
		return true
	} else {
		return false
	}
}

func (si *sliceIter[T]) Value() T {
	return si.value
}

func (si *sliceIter[T]) Abort() {
	si.index = len(si.slice)
}

func (si *sliceIter[T]) Size() IteratorSize {
	return NewSize(len(si.slice) - si.index)
}

func (si *sliceIter[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		defer si.Abort()
		for si.index = 0; si.index < len(si.slice); si.index++ {
			si.value = si.slice[si.index]
			if !yield(si.value) {
				break
			}
		}
	}
}

func (si *sliceIter[T]) PreferSeq() bool { return false }

// Slice makes an Iterator[T] from slice []T, containing all the elements
// from the slice in order.
func Slice[T any](slice []T) Iterator[T] {
	iter := &sliceIter[T]{slice: slice, index: 0}
	return NewDefaultIterator(iter)
}

// Of makes an Iterator[T] containing the variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	return Slice(elements)
}

type rangeIter[T ordered.Real, S ordered.Real] struct {
	index, to T
	by        S
	value     T
	inclusive bool
}

func (ri *rangeIter[T, S]) incdec() {
	if ri.by < 0 {
		ri.index -= T(-ri.by) // T might not be signed
	} else {
		ri.index += T(ri.by)
	}
}

func (ri *rangeIter[T, S]) validateRange() {
	if ri.by == 0 && ri.index != ri.to {
		panic(fmt.Errorf("%w: step is zero", ErrInvalidIteratorRange))
	}
	if (ri.by > 0 && ri.to < ri.index) || (ri.by < 0 && ri.to > ri.index) {
		panic(fmt.Errorf("%w: negative step or inverse range (but not both)", ErrInvalidIteratorRange))
	}
}

func (ri *rangeIter[T, S]) Next() bool {
	if ri.index == ri.to {
		// Handles the case where by is zero, which is valid if index is at the end
		if ri.inclusive {
			ri.value = ri.index
			ri.inclusive = false // Causes iterator to terminate next time
			return true
		} else {
			return false
		}
	}
	if (ri.by < 0 && ri.index < ri.to) || (ri.by > 0 && ri.index > ri.to) {
		return false
	}
	ri.value = ri.index
	ri.incdec()
	return true
}

func (ri *rangeIter[T, S]) Value() T {
	return ri.value
}

func (ri *rangeIter[T, S]) Abort() {
	ri.index = ri.to
	ri.inclusive = false
}

func (ri *rangeIter[T, S]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		defer ri.Abort()
		if ri.index == ri.to {
			if ri.inclusive {
				yield(ri.index)
			}
			return
		}
		size, aStep := rangehelper.RangeSize(ri.index, ri.to, ri.by, ri.inclusive)
		if ri.by < 0 {
			for range size {
				index := ri.index
				ri.index -= aStep
				if !yield(index) {
					break
				}
			}
		} else {
			for range size {
				index := ri.index
				ri.index += aStep
				if !yield(index) {
					break
				}
			}
		}

	}
}

func (ri *rangeIter[T, S]) Size() IteratorSize {
	var size int
	if ri.index == ri.to {
		size = functions.IfElse(ri.inclusive, 1, 0)
	} else if (ri.index > ri.to && ri.by > 0) || (ri.index < ri.to && ri.by < 0) {
		size = 0
	} else {
		size, _ = rangehelper.RangeSize(ri.index, ri.to, ri.by, ri.inclusive)
	}
	return NewSize(size)
}

func (ri *rangeIter[T, S]) PreferSeq() bool { return false }

func newRangeIter[T ordered.Real, S ordered.Real](from, upto T, by S, inclusive bool) Iterator[T] {
	itr := rangeIter[T, S]{index: from, to: upto, by: by, inclusive: inclusive}
	itr.validateRange()
	return NewDefaultIterator(&itr)
}

// Range creates an iterator that ranges from `from` to
// `upto` exclusive
func Range[T ordered.Real](from, upto T) Iterator[T] {
	return newRangeIter(from, upto, functions.IfElse(upto < from, -1, 1), false)
}

// Range creates an iterator that ranges from `from` to
// `upto` inclusive
func IncRange[T ordered.Real](from, upto T) Iterator[T] {
	return newRangeIter(from, upto, functions.IfElse(upto < from, -1, 1), true)
}

// RangeBy creates an iterator that ranges from `from` up to
// `upto` exclusive, incrementing by `by` each step.
// This can be negative (and `upto` should be less than `from`),
// but it cannot be zero unless from == upto, in which case
// an empty iterator is returned.
func RangeBy[T ordered.Real, S ordered.Real](from, upto T, by S) Iterator[T] {
	return newRangeIter(from, upto, by, false)
}

// RangeBy creates an iterator that ranges from `from` up to
// `upto` inclusive, incrementing by `by` each step.
// This can be negative (and `upto` should be less than `from`),
// but it cannot be zero unless from == upto, in which case
// an iterator with a single value is returned.
func IncRangeBy[T ordered.Real, S ordered.Real](from, upto T, by S) Iterator[T] {
	return newRangeIter(from, upto, by, true)
}

type emptyIter[T any] struct{}

func (emptyIter[T]) Next() bool         { return false }
func (emptyIter[T]) Value() T           { var zero T; return zero }
func (emptyIter[T]) Size() IteratorSize { return NewSize(0) }
func (emptyIter[T]) Abort()             {}
func (emptyIter[T]) PreferSeq() bool    { return false }
func (emptyIter[T]) Seq() iter.Seq[T]   { return func(yield func(T) bool) {} }

// Empty creates an iterator that returns no items.
func Empty[T any]() Iterator[T] {
	return NewDefaultIterator(emptyIter[T]{})
}

// Map applies function `mapping` of type `func(T) U` to each value, producing
// a new iterator over `U`.
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	mapNext := func(iterator Iterator[T]) (ok bool, value U) {
		if ok = iterator.Next(); ok {
			value = mapping(iterator.Value())
		}
		return
	}
	return wrapFunc(iter, mapNext, functions.Id)
}

// Filter applies a filter function `predicate` of type `func(T) bool`, producing
// a new iterator containing only the elements than satisfy the function.
func Filter[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	filterNext := func(i Iterator[T]) (ok bool, value T) {
		for {
			if ok = i.Next(); ok {
				if !predicate(i.Value()) {
					continue
				} else {
					value = i.Value()
				}
			}
			return
		}
	}
	return wrapFunc(iter, filterNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterMap applies both transformation and filtering logic to an iterator. The function `mapping` is
// applied to each element of type `T`, producing either an option value of type `U` or an empty
// option. The result is an iterator over `U` drawn from only the non-empty options
// returned.
func FilterMap[T any, U any](iter Iterator[T], mapping func(T) option.Option[U]) Iterator[U] {
	filterMapNext := func(i Iterator[T]) (ok bool, value U) {
		for {
			if ok = i.Next(); ok {
				if value, ok = mapping(i.Value()).ToRef().GetOK(); !ok {
					continue
				}
			}
			return
		}
	}
	return wrapFunc(iter, filterMapNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterValues takes an iterator of results and returns an iterator of the underlying
// result value type for only those results that have no error.
func FilterValues[T any](iter Iterator[result.Result[T]]) Iterator[T] {
	return FilterMap(iter, func(res result.Result[T]) option.Option[T] {
		if res.IsError() {
			return option.Empty[T]()
		} else {
			return option.Value(res.Get())
		}
	})
}

// CollectInto collects all elements from an iterator into a pointer to a slice.
// The slice referenced may be reallocated as the append function is used to add
// elements to the slice. The slice may be a nil slice.
func CollectInto[T any](iter CoreIterator[T], slice *[]T) []T {
	if iter.PreferSeq() {
		for v := range iter.Seq() {
			*slice = append(*slice, v)
		}
	} else {
		for iter.Next() {
			*slice = append(*slice, iter.Value())
		}
	}
	return *slice
}

// Collect collects all elements from an iterator into a slice.
func Collect[T any](iter CoreIterator[T]) []T {
	result := make([]T, 0, iter.Size().Allocate())
	return CollectInto(iter, &result)
}

// CollectResults collects all elements from an iterator of results into a result of slice of the iterator's underlying type
// If the iterator returns an error result at any point, this call will terminate and return that error, along with the elements
// collected thus far.
func CollectResults[T any](iter CoreIterator[result.Result[T]]) ([]T, error) {
	collectResult := make([]T, 0, iter.Size().Allocate())
	for res := range iter.Seq() {
		if res.IsError() {
			return collectResult, res.GetErr()
		}
		collectResult = append(collectResult, res.Get())
	}
	return collectResult, nil
}

// PartitionResults collects the elements from an iterator of result types into two slices, one of
// successful (nil error) values, and the other of error values.
func PartitionResults[T any](iter CoreIterator[result.Result[T]]) ([]T, []error) {
	values := make([]T, 0, iter.Size().Allocate())
	var errs []error
	for res := range iter.Seq() {
		if res.IsError() {
			errs = append(errs, res.GetErr())
		} else {
			values = append(values, res.Must())
		}
	}
	return values, errs
}

// All returns true if `predicate` returns true for every value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns false.
func All[T any](iter CoreIterator[T], predicate func(v T) bool) bool {
	for v := range iter.Seq() {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if `predicate` returns true for any value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns true.
func Any[T any](iter CoreIterator[T], predicate func(v T) bool) bool {
	for v := range iter.Seq() {
		if predicate(v) {
			return true
		}
	}
	return false
}

// A core iterator that obtains values (or an error) from
// a channel
type genIter[T any] struct {
	source chan T
	value  T
}

func newGenIter[T any](source chan T) *genIter[T] {
	return &genIter[T]{source: source}
}

func (pi *genIter[T]) Next() bool {
	var ok bool
	pi.value, ok = <-pi.source
	return ok
}

func (pi *genIter[T]) Value() T {
	return pi.value
}

func (pi *genIter[T]) Chan() <-chan T {
	return pi.source
}

func (pi *genIter[T]) Abort() {
	safeClose(pi.source)
}

func (pi *genIter[T]) Size() IteratorSize {
	return NewSizeUnknown()
}

func (pi *genIter[T]) Seq() iter.Seq[T] {
	return Seq(pi)
}

func (pi *genIter[T]) PreferSeq() bool {
	return false
}

// AbortGenerator is a panic type that will be raised if a Generator function is to be
// aborted.
type AbortGenerator struct{}

// Consumer is a type, an instance of which is passed to a Generator generator
// function. Values from the function can be yielded to the generator
// via the Yield method (or an error via the YieldError method).
type Consumer[T any] struct {
	sink chan T
}

// Yield yields the next value to the generator
func (y Consumer[T]) Yield(t T) {
	if !safeSend(y.sink, t) {
		panic(AbortGenerator{})
	}
}

// ResultConsumer is a variation on `Consumer` which is used to yield only result types. It adds
// dedicated methods to yield non-error values and errors.
type ResultConsumer[T any] Consumer[result.Result[T]]

// Yield yields the next result to the result consumer
func (yr *ResultConsumer[T]) Yield(value result.Result[T]) {
	(*Consumer[result.Result[T]])(yr).Yield(value)
}

// YieldValue yields the next successful value to the consumer
func (yr *ResultConsumer[T]) YieldValue(value T) {
	yr.Yield(result.Value(value))
}

// YieldError yields an error to the consumer
func (yr *ResultConsumer[T]) YieldError(err error) {
	yr.Yield(result.Error[T](err))
}

// GeneratorPanic is an error type indicating that a generator iterator function has panicked
type GeneratorPanic struct {
	panic any
}

func (pp GeneratorPanic) Error() string {
	return fmt.Sprintf("panic in generator: %#v", pp.panic)
}

func (pp GeneratorPanic) Unwrap() error {
	if err, ok := pp.panic.(error); ok {
		return err
	} else {
		return nil
	}
}

// Generator is a function taking a Consumer. The function is expected to yield values to the consumer.
type Generator[T any] func(Consumer[T])

func runGenerator[T any](c Consumer[T], activity Generator[T]) {
	defer safeClose(c.sink)
	defer func() {
		if p := recover(); p != nil {
			if _, abort := p.(AbortGenerator); !abort {
				panic(p)
			}
		}
	}()
	activity(c)
}

// ResultGenerator is a function taking a ResultConsumer object to which results may be yielded.
// If a non-nil error is returned, it will be yielded as an error result.
type ResultGenerator[T any] func(ResultConsumer[T]) error

func runResultGenerator[T any](c ResultConsumer[T], activity ResultGenerator[T]) {
	defer safeClose(c.sink)
	defer func() {
		if p := recover(); p != nil {
			if _, abort := p.(AbortGenerator); !abort {
				c.YieldError(GeneratorPanic{p})
			}
		}
	}()
	defer eh.Handle(func(err error) { c.YieldError(err) })
	eh.Check(activity(c))
}

/*
Generate creates an Iterator from a Generator function. A Consumer is created and passed to the function.
The function is run in a separate goroutine, and its yielded values are sent over a channel
to the iterator where they may be consumed in an iterative way by calls to Next() and Value().
Alternatively, the channel itself is available via the Chan() method.
A call to Abort() will cause the channel to close and no further elements will be produced by
Next() or a read of the channel. Any attempt to subsequently yield a value in the generator
will cause it to terminate, via an AbortGenerator panic.
*/
func Generate[T any](generator Generator[T]) Iterator[T] {
	ch := make(chan T)
	yield := Consumer[T]{ch}
	go runGenerator(yield, generator)
	return NewDefaultIterator(newGenIter(ch))
}

// GenerateResults is a variation on Generate that produces an iterator of result types. If the
// generator function panics, an error result of type GeneratorPanic is produced prior to closing
// the consumer channel.
func GenerateResults[T any](generator ResultGenerator[T]) Iterator[result.Result[T]] {
	ch := make(chan result.Result[T])
	yield := ResultConsumer[T](Consumer[result.Result[T]]{ch})
	go runResultGenerator(yield, generator)
	return NewDefaultIterator(newGenIter(ch))
}

// Iterators over iterators
type takeIterator[T any] struct {
	count, max int
	aborted    bool
	iterator   Iterator[T]
	value      T
}

func (ti *takeIterator[T]) Value() T {
	ti.value = ti.iterator.Value()
	return ti.value
}

func (ti *takeIterator[T]) Abort() {
	if !ti.aborted {
		ti.iterator.Abort()
	}
	ti.aborted = true
}

func (ti *takeIterator[T]) Next() bool {
	if !ti.aborted && ti.count < ti.max {
		ti.count++
		return ti.iterator.Next()
	} else {
		return false
	}
}

// Non-seq is more efficient  here
func (ti *takeIterator[T]) PreferSeq() bool { return false }

func (ti *takeIterator[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if ti.count < ti.max && !ti.aborted {
			next, _ := iter.Pull(ti.iterator.Seq())
			for ti.count < ti.max {
				var ok bool
				if ti.value, ok = next(); ok {
					if !yield(ti.value) {
						ti.aborted = true
						break
					}
					ti.count++
				} else {
					break
				}
			}
		}
	}
}

func (ti *takeIterator[T]) Size() IteratorSize {
	itrSize := ti.iterator.Size()
	remain := ti.max - ti.count
	switch itrSize.Type {
	case SizeKnown:
		return NewSize(min(remain, itrSize.Size))
	case SizeUnknown:
		return NewSizeAtMost(remain)
	case SizeAtMost:
		return NewSizeAtMost(min(remain, itrSize.Size))
	case SizeInfinite:
		return NewSize(remain)
	default:
		panic(ErrInvalidIteratorSizeType)
	}
}

// Take transforms an iterator into an iterator the returns the
// first n elements of the original iterator. If there are less
// than n elements available, they are all returned.
func Take[T any](n int, iter Iterator[T]) Iterator[T] {
	return NewDefaultIterator(&takeIterator[T]{iterator: iter, max: n})
}

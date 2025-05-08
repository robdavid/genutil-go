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

// SizedIterator is an extension of SimpleIterator that also holds some sizing information
type SizedIterator[T any] interface {
	SimpleIterator[T]
	// Size is an estimate, where possible, of the number of elements remaining.
	Size() IteratorSize
}

// Generic iterator
type Iterator[T any] interface {
	SizedIterator[T]
	// Chan returns iterator as a channel.
	Chan() <-chan T
	Seq() iter.Seq[T]
}

type DefaultIterator[T any] struct {
	SimpleIterator[T]
}

func (di DefaultIterator[T]) Chan() <-chan T {
	return Chan(di)
}

func (di DefaultIterator[T]) Seq() iter.Seq[T] {
	return Seq(di)
}

func (di DefaultIterator[T]) Size() IteratorSize {
	return NewSizeUnknown()
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

// FuncNext is a function supporting a transforming operation by consuming
// all or part of an iterator, returning the next value
type FuncNext[T any, U any] func(Iterator[T]) (bool, U)

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
func Chan[T any](iter SimpleIterator[T]) (out chan T) {
	out = make(chan T)
	go func() {
		defer safeClose(out)
		for iter.Next() {
			if !safeSend(out, iter.Value()) {
				break
			}
		}
	}()
	return
}

type sizedIterator[T any] struct {
	SimpleIterator[T]
	size IteratorSize
}

func (si sizedIterator[T]) Size() IteratorSize {
	return si.size
}

// MakeIterator creates a generic iterator from a simple iterator. Provides an implementation
// of additional Iterator methods.
func MakeIterator[T any](base SimpleIterator[T]) Iterator[T] {
	return DefaultIterator[T]{base}
}

// mapIter wraps an iterator and adds a mapping function
type mapIter[T, U any] struct {
	base    Iterator[T]
	mapping FuncNext[T, U]
	value   option.Option[U]
	outChan chan U
	size    IteratorSize
}

func (i *mapIter[T, U]) Next() bool {
	ok, value := i.mapping(i.base)
	if !ok {
		return false
	} else {
		i.value.Set(value)
		return true
	}
}

func (i *mapIter[T, U]) Value() U {
	return i.value.GetOrZero()
}

func (i *mapIter[T, U]) Abort() {
	if i.outChan != nil {
		safeClose(i.outChan)
		i.outChan = nil
	}
	i.base.Abort()
}

func (i *mapIter[T, U]) Chan() <-chan U {
	if i.outChan == nil {
		i.outChan = Chan[U](i)
	}
	return i.outChan
}

func (i *mapIter[T, U]) Size() IteratorSize {
	return i.size
}

func (i *mapIter[T, U]) Seq() iter.Seq[U] {
	return Seq(i)
}

// wrapFunc creates a new iterator from an existing operator and a function that consumes it, yielding
// one element at a time.
func wrapFunc[T any, U any](iterator Iterator[T], f FuncNext[T, U], size IteratorSize) Iterator[U] {
	return &mapIter[T, U]{base: iterator, mapping: f, value: option.Empty[U](), size: size}
}

// Iterator over a slice
type sliceIter[T any] struct {
	DefaultIterator[T]
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

// Slice makes an Iterator[T] from slice []T, containing all the elements
// from the slice in order.
func Slice[T any](slice []T) Iterator[T] {
	iter := &sliceIter[T]{slice: slice, index: 0}
	iter.SimpleIterator = iter
	return iter
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

func (ri *rangeIter[T, S]) Chan() <-chan T {
	return Chan(ri)
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
				if !yield(ri.index) {
					break
				}
				ri.index -= aStep
			}
		} else {
			for range size {
				if !yield(ri.index) {
					break
				}
				ri.index += aStep
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

func newRangeIter[T ordered.Real, S ordered.Real](from, upto T, by S, inclusive bool) *rangeIter[T, S] {
	itr := rangeIter[T, S]{index: from, to: upto, by: by, inclusive: inclusive}
	itr.validateRange()
	return &itr
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
func (emptyIter[T]) Chan() <-chan T {
	c := make(chan T)
	close(c)
	return c
}
func (emptyIter[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		return
	}
}

// Empty creates an iterator that returns no items.
func Empty[T any]() Iterator[T] { return emptyIter[T]{} }

// Map applies function `mapping` of type `func(T) U` to each value, producing
// a new iterator over `U`.
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	mapNext := func(iterator Iterator[T]) (ok bool, value U) {
		if ok = iterator.Next(); ok {
			value = mapping(iterator.Value())
		}
		return
	}
	return wrapFunc(iter, mapNext, iter.Size())
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
	return wrapFunc(iter, filterNext, iter.Size().Subset())
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
	return wrapFunc(iter, filterMapNext, iter.Size().Subset())
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
func CollectInto[T any](iter Iterator[T], slice *[]T) []T {
	for iter.Next() {
		*slice = append(*slice, iter.Value())
	}
	return *slice
}

// Collect collects all elements from an iterator into a slice.
func Collect[T any](iter Iterator[T]) []T {
	result := make([]T, 0, iter.Size().Allocate())
	return CollectInto(iter, &result)
}

// CollectResults collects all elements from an iterator of results into a result of slice of the iterator's underlying type
// If the iterator returns an error result at any point, this call will terminate and return that error in the
// result, along with the elements collected thus far.
func CollectResults[T any](iter Iterator[result.Result[T]]) ([]T, error) {
	collectResult := make([]T, 0, iter.Size().Allocate())
	for iter.Next() {
		res := iter.Value()
		if res.IsError() {
			iter.Abort()
			return collectResult, res.GetErr()
		}
		collectResult = append(collectResult, res.Get())
	}
	return collectResult, nil
}

// PartitionResults collects the elements from an iterator of result types into two slices, one of
// successful (nil error) values, and the other of error values.
func PartitionResults[T any](iter Iterator[result.Result[T]]) ([]T, []error) {
	values := make([]T, 0, iter.Size().Allocate())
	var errs []error
	for iter.Next() {
		if res := iter.Value(); res.IsError() {
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
func All[T any](iter Iterator[T], predicate func(v T) bool) bool {
	for iter.Next() {
		if !predicate(iter.Value()) {
			iter.Abort()
			return false
		}
	}
	return true
}

// Any returns true if `predicate` returns true for any value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns true.
func Any[T any](iter Iterator[T], predicate func(v T) bool) bool {
	for iter.Next() {
		if predicate(iter.Value()) {
			iter.Abort()
			return true
		}
	}
	return false
}

// An iterator that obtains values (or an error) from
// a channel
type genIter[T any] struct {
	source chan T
	value  T
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
	return &genIter[T]{source: ch}
}

// GenerateResults is a variation on Generate that produces an iterator of result types. If the
// generator function panics, an error result of type GeneratorPanic is produced prior to closing
// the consumer channel.
func GenerateResults[T any](generator ResultGenerator[T]) Iterator[result.Result[T]] {
	ch := make(chan result.Result[T])
	yield := ResultConsumer[T](Consumer[result.Result[T]]{ch})
	go runResultGenerator(yield, generator)
	return &genIter[result.Result[T]]{source: ch}
}

// Iterators over iterators
type takeIterator[T any] struct {
	DefaultIterator[T]
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
	i := &takeIterator[T]{iterator: iter, max: n}
	i.SimpleIterator = i
	return i
}

// SeqIterator is an iterator based on a Go iter.Seq
type SeqIterator[T any] struct {
	DefaultIterator[T]
	seq     iter.Seq[T]
	seqNext func() (T, bool)
	seqStop func()
	value   T
}

// FromSeq creates an iterator from a Go iter.Seq function
func FromSeq[T any](seq iter.Seq[T]) Iterator[T] {
	var si SeqIterator[T]
	si.SimpleIterator = &si
	si.seq = seq
	si.seqNext, si.seqStop = iter.Pull(seq)
	return &si
}

// FromSeq creates an iterator from a Go iter.Seq2 function
func FromSeq2[K any, V any](seq iter.Seq2[K, V]) Iterator[tuple.Tuple2[K, V]] {
	var si SeqIterator[tuple.Tuple2[K, V]]
	si.seq = func(yield func(tuple.Tuple2[K, V]) bool) {
		for k, v := range seq {
			if !yield(tuple.Pair(k, v)) {
				break
			}
		}
	}
	var seq2Next func() (K, V, bool)
	seq2Next, si.seqStop = iter.Pull2(seq)
	si.seqNext = func() (tuple.Tuple2[K, V], bool) {
		k, v, b := seq2Next()
		return tuple.Pair(k, v), b
	}
	return &si
}

func (si *SeqIterator[T]) Next() (next bool) {
	si.value, next = si.seqNext()
	return
}

func (si *SeqIterator[T]) Value() T {
	return si.value
}

func (si *SeqIterator[T]) Abort() {
	si.seqStop()
}

func (si *SeqIterator[T]) Seq() iter.Seq[T] {
	return si.seq
}

func (si *SeqIterator[T]) Size() IteratorSize {
	return NewSizeUnknown()
}

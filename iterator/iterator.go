// Iterators and generators
package iterator

import (
	"fmt"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/option"
)

// The largest slice capacity we are prepared to allocate to collect
// iterators of uncertain size.
const maxUncertainAllocation = 100000

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
}

// IteratorSize holds iterator sizing information
type IteratorSize interface {
	// Allocate returns a value for the size of slice required to hold all elements returned by the iterator.
	// This may be an exact value, an estimated value or unknown (in which case 0 is returned)
	Allocate() int
	// Subset returns a new iterator size based on some unknown subset of iterator values.
	Subset() IteratorSize
}

// Iterator sizing information; size is unknown
type sizeUnknown struct{}

func (sizeUnknown) Allocate() int        { return 0 }
func (sizeUnknown) Subset() IteratorSize { return SizeUnknown }

// SizeUnknown is a value implementing IteratorSize, representing an iterator of unknown size
var SizeUnknown = sizeUnknown{}

// IsSizeUnknown returns true if the given IteratorSize instance represents
// an unknown size
func IsSizeUnknown(size IteratorSize) bool {
	return size == SizeUnknown
}

// SizeKnown holds Iterator sizing information where the size is known with certainty.
type SizeKnown struct {
	Size int
}

func (sk SizeKnown) Allocate() int        { return sk.Size }
func (sk SizeKnown) Subset() IteratorSize { return SizeAtMost(sk) }

// NewSize creates an `IteratorSize` implementation that has a fixed size of `n`.
func NewSize(n int) IteratorSize { return SizeKnown{n} }

// IsSizeKnown returns true if the iterator size is one whose actual size is known.
func IsSizeKnown(size IteratorSize) bool {
	_, ok := size.(SizeKnown)
	return ok
}

// SizeAtMost holds Iterator sizing information where only an upper bound is known.
type SizeAtMost struct {
	Size int
}

func (sm SizeAtMost) Allocate() int {
	alloc := sm.Size / 2
	if alloc > maxUncertainAllocation {
		alloc = maxUncertainAllocation
	}
	return alloc
}
func (sm SizeAtMost) Subset() IteratorSize { return sm }

// NewSizeAtMost creates an `IteratorSize` implementation that has a size no more than n.
func NewSizeAtMost(n int) IteratorSize {
	return SizeAtMost{n}
}

// IsSizeAtMost returns true if the iterator size is one whose maximum size is known.
func IsSizeAtMost(size IteratorSize) bool {
	_, ok := size.(SizeAtMost)
	return ok
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
func iterChan[T any](iter SimpleIterator[T]) (out chan T) {
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

type unknownSizeIterator[T any] struct {
	SimpleIterator[T]
}

func (usi unknownSizeIterator[T]) Size() IteratorSize {
	return SizeUnknown
}

// MakeSizedIterator creates a sized iterator from a simple iterator; the iterator size will
// be unknown.
func MakeSizedIterator[T any](si SimpleIterator[T]) SizedIterator[T] {
	return unknownSizeIterator[T]{si}
}

type sizedIterator[T any] struct {
	SimpleIterator[T]
	size IteratorSize
}

func (si sizedIterator[T]) Size() IteratorSize {
	return si.size
}

// MakeIteratorOfSize creates a sized iterator from a simple iterator; the iterator size will be the
// size provided.
func MakeIteratorOfSize[T any](si SimpleIterator[T], size IteratorSize) SizedIterator[T] {
	return sizedIterator[T]{si, size}
}

type autoChannelIterator[T any] struct {
	SizedIterator[T]
	OutChan chan T
}

// MakeIterator creates a generic iterator from a simple iterator. Provides an implementation
// of a source of elements over a channel.
func MakeIterator[T any](base SizedIterator[T]) Iterator[T] {
	return &autoChannelIterator[T]{base, nil}
}

// MakeIteratorFromSimple creates a generic iterator from a sized iterator. Provides an implementation
// of a source of elements over a channel.
func MakeIteratorFromSimple[T any](base SimpleIterator[T]) Iterator[T] {
	return MakeIterator(MakeSizedIterator(base))
}

func (ei *autoChannelIterator[T]) Chan() <-chan T {
	if ei.OutChan == nil {
		ei.OutChan = iterChan[T](ei)
	}
	return ei.OutChan
}

func (ei *autoChannelIterator[T]) Abort() {
	ei.SizedIterator.Abort()
	if ei.OutChan != nil {
		safeClose(ei.OutChan)
		ei.OutChan = nil
	}
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
		i.outChan = iterChan[U](i)
	}
	return i.outChan
}

func (i *mapIter[T, U]) Size() IteratorSize {
	return i.size
}

// wrapFunc creates a new iterator from an existing operator and a function that consumes it, yielding
// one element at a time.
func wrapFunc[T any, U any](iterator Iterator[T], f FuncNext[T, U], size IteratorSize) Iterator[U] {
	return &mapIter[T, U]{base: iterator, mapping: f, value: option.Empty[U](), size: size}
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
	return SizeKnown{len(si.slice) - si.index}
}

// Slice makes an Iterator[T] from slice []T, containing all the elements
// from the slice in order.
func Slice[T any](slice []T) Iterator[T] {
	var t T
	return MakeIterator[T](&sliceIter[T]{slice, 0, t})
}

// Of makes an Iterator[T] containing the variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	return Slice(elements)
}

// Numeric types with ordering
type Numeric interface {
	uint | uint8 | uint16 | uint32 | uint64 |
		int | int8 | int16 | int32 | int64 |
		float32 | float64
}

type rangeIter[T Numeric] struct {
	index, to, by, value T
}

func (ri *rangeIter[T]) Next() bool {
	if ri.by < 0 {
		if ri.index <= ri.to {
			return false
		}
	} else if ri.index >= ri.to {
		return false
	}
	ri.value = ri.index
	ri.index += ri.by
	return true
}

func (ri *rangeIter[T]) Value() T {
	return ri.value
}

func (ri *rangeIter[T]) Abort() {
	ri.index = ri.to
}

func (ri *rangeIter[T]) Chan() <-chan T {
	return iterChan[T](ri)
}

func (ri *rangeIter[T]) Size() IteratorSize {
	size := int((ri.to - ri.index) / ri.by)
	if size < 0 {
		size = 0
	}
	return SizeKnown{size}
}

// Range creates an iterator that ranges from `from` to
// `upto` exclusive
func Range[T Numeric](from, upto T) Iterator[T] {
	return &rangeIter[T]{from, upto, 1, 0}
}

// RangeBy creates an iterator that ranges from `from` up to
// `upto` exclusive, incrementing by `by` each step.
// This can be negative (and `upto` should be less than `from`),
// but it cannot be zero.
func RangeBy[T Numeric](from, upto, by T) Iterator[T] {
	if by == 0 {
		panic("Illegal range by zero")
	}
	return &rangeIter[T]{from, upto, by, 0}
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

// FilterResults takes an iterator of results and returns an iterator of the underlying
// result type for only those results that have no error.
func FilterResults[T any](iter Iterator[result.Result[T]]) Iterator[T] {
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
	return CollectInto[T](iter, &result)
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
	return SizeUnknown
}

type AbortGenerator struct{}

// Generator is a type, an instance of which is passed to a GenFunc generator
// function. Values from the function can be yielded to the generator
// via the Yield method (or an error via the YieldError method).
type Generator[T any] struct {
	sink chan T
}

// Yield yields the next value to the generator
func (y Generator[T]) Yield(t T) {
	if !safeSend(y.sink, t) {
		panic(AbortGenerator{})
	}
}

// ResultGenerator is a variation on `Yield` which is used to yield only result types. It adds
// dedicated methods to yield non-error values and errors.
type ResultGenerator[T any] Generator[result.Result[T]]

// Yield yields the next result to the result generator
func (yr *ResultGenerator[T]) Yield(value result.Result[T]) {
	(*Generator[result.Result[T]])(yr).Yield(value)
}

// YieldValue yields the next successful value to the generator
func (yr *ResultGenerator[T]) YieldValue(value T) {
	yr.Yield(result.Value(value))
}

// YieldError yields an error to the generator
func (yr *ResultGenerator[T]) YieldError(err error) {
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

// GenFunc is a function taking a Generator, to which values can be yielded.
type GenFunc[T any] func(Generator[T])

func runGenerator[T any](y Generator[T], activity GenFunc[T]) {
	defer safeClose(y.sink)
	defer func() {
		if p := recover(); p != nil {
			if _, abort := p.(AbortGenerator); !abort {
				panic(p)
			}
		}
	}()
	activity(y)
}

// GenResultFunc is a function taking a ResultGenerator object to which results to be yielded.
// If a non-nil error is returned, it will be yielded as an error result.
type GenResultFunc[T any] func(ResultGenerator[T]) error

func runResultGenerator[T any](y ResultGenerator[T], activity GenResultFunc[T]) {
	defer safeClose(y.sink)
	defer func() {
		if p := recover(); p != nil {
			if _, abort := p.(AbortGenerator); !abort {
				y.YieldError(GeneratorPanic{p})
			}
		}
	}()
	defer eh.Handle(func(err error) { y.YieldError(err) })
	eh.Check(activity(y))
}

// Generate creates an iterator from a function, passed in activity, that yields a sequence of values
// by making calls to Yield(), on a Generator object that will be passed to the function.
// The function is run in a separate goroutine, and its yielded values are sent over a channel
// to the iterator where they can be consumed in an iterative way by calls to Next() and Value().
// The channel itself is available via the Chan() method. A call to Abort() will cause the channel
// to close and no further elements will be produced by Next() or a read of the channel. Any attempt
// to subsequently yield a value in the activity function will cause it to abort (via a panic).
// The activity parameter is of type PipeFunc[T], which is an alias for func(Yield[T])
func Generate[T any](activity GenFunc[T]) Iterator[T] {
	ch := make(chan T)
	yield := Generator[T]{ch}
	go runGenerator(yield, activity)
	return &genIter[T]{source: ch}
}

// GenerateResults is a variation on Generate that produces an iterator of result types.
func GenerateResults[T any](activity GenResultFunc[T]) Iterator[result.Result[T]] {
	ch := make(chan result.Result[T])
	yield := ResultGenerator[T](Generator[result.Result[T]]{ch})
	go runResultGenerator(yield, activity)
	return &genIter[result.Result[T]]{source: ch}
}

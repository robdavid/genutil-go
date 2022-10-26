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

// Generic iterator
type Iterator[T any] interface {
	// Set the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned when there are no more elements (the current value remains unchanged)
	Next() bool
	Value() T           // Get the current iterator value.
	Abort()             // Stop the iterator; subsequent calls to Next() will return false.
	Chan() <-chan T     // Return iterator as a channel.
	Size() IteratorSize // Size estimate, where possible, of the number of elements remaining.
}

// Iterator sizing information
type IteratorSize interface {
	Allocate() int          // The best guess of size to allocate to collect this iterator
	Filtered() IteratorSize // A new size estimate of an iterator filtered from this one
}

// Iterator sizing information; size is unknown
type SizeUnknown struct{}

func (SizeUnknown) Allocate() int          { return 0 }
func (SizeUnknown) Filtered() IteratorSize { return SizeUnknown{} }

// Iterator sizing information; size is known with certainty
type SizeKnown struct {
	Size int
}

func (sk SizeKnown) Allocate() int          { return sk.Size }
func (sk SizeKnown) Filtered() IteratorSize { return SizeAtMost(sk) }

// Iterator sizing information; maximum size is known with certainty
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
func (sm SizeAtMost) Filtered() IteratorSize { return sm }

// A function supporting a transforming operation by consuming
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

// Generic channel implementation
func iterChan[T any](iter Iterator[T]) (out chan T) {
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

// Wraps an iterator and adds a mapping function
type Func[T, U any] struct {
	base    Iterator[T]
	mapping FuncNext[T, U]
	value   option.Option[U]
	outChan chan U
	size    IteratorSize
}

func (i *Func[T, U]) Next() bool {
	ok, value := i.mapping(i.base)
	if !ok {
		return false
	} else {
		i.value.Set(value)
		return true
	}
}

func (i *Func[T, U]) Value() U {
	return i.value.GetOrZero()
}

func (i *Func[T, U]) Abort() {
	if i.outChan != nil {
		safeClose(i.outChan)
		i.outChan = nil
	}
	i.base.Abort()
}

func (i *Func[T, U]) Chan() <-chan U {
	if i.outChan == nil {
		i.outChan = iterChan[U](i)
	}
	return i.outChan
}

func (i *Func[T, U]) Size() IteratorSize {
	return i.size
}

// Create a new iterator from an existing operator and a function that consumes it, yielding
// one element at a time.
func WrapFunc[T any, U any](iterator Iterator[T], f FuncNext[T, U], size IteratorSize) Iterator[U] {
	return &Func[T, U]{base: iterator, mapping: f, value: option.Empty[U](), size: size}
}

// Iterator over a slice
type SliceIter[T any] struct {
	slice   []T
	index   int
	value   T
	outChan chan T
}

func (si *SliceIter[T]) Next() bool {
	if si.index < len(si.slice) {
		si.value = si.slice[si.index]
		si.index++
		return true
	} else {
		return false
	}
}

func (si *SliceIter[T]) Value() T {
	return si.value
}

func (si *SliceIter[T]) Abort() {
	si.index = len(si.slice)
	if si.outChan != nil {
		safeClose(si.outChan)
	}
}

func (si *SliceIter[T]) Chan() <-chan T {
	if si.outChan == nil {
		si.outChan = iterChan[T](si)
	}
	return si.outChan
}

func (si *SliceIter[T]) Size() IteratorSize {
	return SizeKnown{len(si.slice) - si.index}
}

// Makes an Iterator[T] from slice []T, containing all the elements
// from the slice in order.
func Slice[T any](slice []T) Iterator[T] {
	var t T
	return &SliceIter[T]{slice, 0, t, nil}
}

// Makes an Iterator[T] containing the variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	var t T
	return &SliceIter[T]{elements, 0, t, nil}
}

type RangeIter struct {
	index, to, by, value int
}

func (ri *RangeIter) Next() bool {
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

func (ri *RangeIter) Value() int {
	return ri.value
}

func (ri *RangeIter) Abort() {
	ri.index = ri.to
}

func (ri *RangeIter) Chan() <-chan int {
	return iterChan[int](ri)
}

func (ri *RangeIter) Size() IteratorSize {
	size := (ri.to - ri.index) / ri.by
	if size < 0 {
		size = 0
	}
	return SizeKnown{size}
}

// Create an iterator that ranges from `from` up to
// `upto` exclusive
func Range(from, upto int) Iterator[int] {
	return &RangeIter{from, upto, 1, 0}
}

// Create an iterator that ranges from `from` up to
// `upto` exclusive, incrementing by `by` each step.
// This can be negative (and `upto` should be less that `from`),
// but it cannot be zero.
func RangeBy(from, upto, by int) Iterator[int] {
	if by == 0 {
		panic("Illegal range by zero")
	}
	return &RangeIter{from, upto, by, 0}
}

// Applies function `mapping` of type `func(T) U` to each value, producing
// a new iterator over `U`.
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	mapNext := func(iterator Iterator[T]) (ok bool, value U) {
		if ok = iterator.Next(); ok {
			value = mapping(iterator.Value())
		}
		return
	}
	return WrapFunc(iter, mapNext, iter.Size())
}

// Applies a filter function `predicate` of type `func(T) bool`, producing
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
	return WrapFunc(iter, filterNext, iter.Size().Filtered())
}

// Applies both transformation and filtering logic to an iterator. The function `mapping` is
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
	return WrapFunc(iter, filterMapNext, iter.Size().Filtered())
}

// Takes an iterator of results and returns an iterator of the underlying
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

// Collects all elements from an iterator into a slice.
func Collect[T any](iter Iterator[T]) []T {
	result := make([]T, 0, iter.Size().Allocate())
	for iter.Next() {
		result = append(result, iter.Value())
	}
	return result
}

// Collects all elements from an iterator of results into a result of slice of the iterator's underlying type
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

// Collect the elements from an iterator of result types into two slices, one of
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

// An iterator that obtains values (or an error) from
// a channel
type GenIter[T any] struct {
	source chan T
	value  T
}

func (pi *GenIter[T]) Next() bool {
	var ok bool
	pi.value, ok = <-pi.source
	return ok
}

func (pi *GenIter[T]) Value() T {
	return pi.value
}

func (pi *GenIter[T]) Chan() <-chan T {
	return pi.source
}

func (pi *GenIter[T]) Abort() {
	safeClose(pi.source)
}

func (pi *GenIter[T]) Size() IteratorSize {
	return SizeUnknown{}
}

type AbortGenerator struct{}

// An object to which consecutive iterator values are supplied via the Yield
// method (or an error via the Error method), sent via a channel
type Yield[T any] struct {
	sink chan T
}

// Yield the next value to the iterator
func (y Yield[T]) Yield(t T) {
	if !safeSend(y.sink, t) {
		panic(AbortGenerator{})
	}
}

// A variation on `Yield` which is used to yield only result types. It adds
// dedicated methods to yield non-error values and errors.
type YieldResult[T any] Yield[result.Result[T]]

// Yield the next result to the iterator
func (yr *YieldResult[T]) Yield(value result.Result[T]) {
	(*Yield[result.Result[T]])(yr).Yield(value)
}

// Yield the next successful value to the iterator
func (yr *YieldResult[T]) YieldValue(value T) {
	yr.Yield(result.Value(value))
}

// Yield an error to the iterator
func (yr *YieldResult[T]) YieldError(err error) {
	yr.Yield(result.Error[T](err))
}

// An error type indicating that a generator iterator function has panicked
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

// A function taking a Yield object, to which values can be yielded.
type GenFunc[T any] func(Yield[T])

func runGenerator[T any](y Yield[T], activity GenFunc[T]) {
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

// A type alias for a function taking a ResultYield object to which results to be yielded.
// If a non-nil error is returned, it will be yielded as an error result.
type GenResultFunc[T any] func(YieldResult[T]) error

func runResultGenerator[T any](y YieldResult[T], activity GenResultFunc[T]) {
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

// Create an iterator from a function, passed in activity, that yields a sequence of values
// by making calls to Yield(), on a Yield object that will be passed to the function.
// The function is run in a separate goroutine, and its yielded values are sent over a channel
// to the iterator where they can be consumed in the normal way by calls to Next() and Value().
// The channel itself is available via the Chan() method. A call to Abort() will cause the channel
// to close and no further elements will be produced by Next() or a read of the channel.
// The activity parameter is of type PipeFunc[T], which is an alias for func(Yield[T])
func Generate[T any](activity GenFunc[T]) Iterator[T] {
	ch := make(chan T)
	yield := Yield[T]{ch}
	go runGenerator(yield, activity)
	return &GenIter[T]{source: ch}
}

// A variation on Generate that produces an iterator of result types.
func GenerateResults[T any](activity GenResultFunc[T]) Iterator[result.Result[T]] {
	ch := make(chan result.Result[T])
	yield := YieldResult[T](Yield[result.Result[T]]{ch})
	go runResultGenerator(yield, activity)
	return &GenIter[result.Result[T]]{source: ch}
}

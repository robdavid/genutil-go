// Iterators and generators
package iterator

import (
	"fmt"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/option"
)

// Generic iterator
type Iterator[T any] interface {
	// Set the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned when there are no more elements (the current value remains unchanged)
	Next() bool
	Value() (T, error)             // Get the current iterator value, or an error if the last call to Next() has failed
	Must() T                       // Get the current iterator value. Panics if the last call to Next() has failed.
	Try() T                        // Get the current iterator value. Creates a Try() panic if the the value has an error.
	Abort()                        // Stop the iterator; subsequent calls to Next() will return false.
	Chan() <-chan result.Result[T] // Return iterator as a channel
}

// A function supporting a transforming operation by consuming
// all or part of an iterator, returning the next value
type FuncNext[T any, U any] func(Iterator[T]) (bool, U, error)

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
func iterChan[T any](iter Iterator[T]) (out chan result.Result[T]) {
	out = make(chan result.Result[T])
	go func() {
		defer safeClose(out)
		for iter.Next() {
			if !safeSend(out, result.From(iter.Value())) {
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
	value   option.Option[result.Result[U]]
	outChan chan result.Result[U]
}

func (i *Func[T, U]) Next() bool {
	ok, value, err := i.mapping(i.base)
	if !ok {
		i.value.Clear()
		return false
	} else {
		i.value.Set(result.From(value, err))
		return true
	}
}

func (i *Func[T, U]) Value() (u U, err error) {
	return i.value.GetOrZero().ToRef().Return()
}

func (i *Func[T, U]) Must() U {
	return i.value.GetOrZero().ToRef().Must()
}

func (i *Func[T, U]) Try() U {
	return i.value.GetOrZero().ToRef().Try()
}

func (i *Func[T, U]) Abort() {
	if i.outChan != nil {
		safeClose(i.outChan)
		i.outChan = nil
	}
	i.base.Abort()
}

func (i *Func[T, U]) Chan() <-chan result.Result[U] {
	if i.outChan == nil {
		i.outChan = iterChan[U](i)
	}
	return i.outChan
}

func WrapFunc[T any, U any](iterator Iterator[T], f FuncNext[T, U]) Iterator[U] {
	return &Func[T, U]{base: iterator, mapping: f, value: option.Empty[result.Result[U]]()}
}

// Iterator over a slice
type SliceIter[T any] struct {
	slice   []T
	index   int
	value   T
	outChan chan result.Result[T]
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

func (si *SliceIter[T]) Must() T {
	return si.value
}

func (si *SliceIter[T]) Try() T {
	return si.value
}

func (si *SliceIter[T]) Value() (T, error) {
	return si.value, nil
}

func (si *SliceIter[T]) Abort() {
	si.index = len(si.slice)
	if si.outChan != nil {
		safeClose(si.outChan)
	}
}

func (si *SliceIter[T]) Chan() <-chan result.Result[T] {
	if si.outChan == nil {
		si.outChan = iterChan[T](si)
	}
	return si.outChan
}

// Makes an Iterator[T] from slice []T
func Slice[T any](slice []T) Iterator[T] {
	var t T
	return &SliceIter[T]{slice, 0, t, nil}
}

// Makes an Iterator[T] from variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	var t T
	return &SliceIter[T]{elements, 0, t, nil}
}

func wrapMap[T any, U any](mapping func(T) (U, error)) FuncNext[T, U] {
	return func(iterator Iterator[T]) (ok bool, value U, err error) {
		if ok = iterator.Next(); ok {
			var vt T
			if vt, err = iterator.Value(); err == nil {
				value, err = mapping(vt)
			}
		}
		return
	}
}

// Wraps an iterator with a mapping function that may return an error, producing a new iterator
func DoMap[T any, U any](iter Iterator[T], mapping func(T) (U, error)) Iterator[U] {
	return WrapFunc[T, U](iter, wrapMap(mapping))
}

// Wraps an iterator with a mapping function, producing a new iterator
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	action := func(t T) (U, error) {
		return mapping(t), nil
	}
	return WrapFunc[T, U](iter, wrapMap(action))
}

// Collects all elements from an iterator into a slice. If the iterator
// returns an error, this call will panic.
func MustCollect[T any](iter Iterator[T]) []T {
	return handler.Must(Collect(iter))
}

// Collects all elements from an iterator into a slice. If the iterator
// returns an error, this call will create handler.Try type panic.
func TryCollect[T any](iter Iterator[T]) []T {
	return handler.Try(Collect(iter))
}

// Collects all elements from an iterator into a slice. If the iterator
// returns an error, this call will terminate and return that error.
func Collect[T any](iter Iterator[T]) ([]T, error) {
	result := make([]T, 0)
	for iter.Next() {
		if val, err := iter.Value(); err != nil {
			return result, err
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

// An iterator that obtains values (or an error) from
// a channel
type PipeIter[T any] struct {
	source chan result.Result[T]
	value  result.Result[T]
}

func (pi *PipeIter[T]) Next() bool {
	var ok bool
	pi.value, ok = <-pi.source
	return ok
}

func (pi *PipeIter[T]) Value() (T, error) {
	return pi.value.Return()
}

func (pi *PipeIter[T]) Must() T {
	return pi.value.Must()
}

func (pi *PipeIter[T]) Try() T {
	return pi.value.Must()
}

func (pi *PipeIter[T]) Chan() <-chan result.Result[T] {
	return pi.source
}

func (pi *PipeIter[T]) Abort() {
	safeClose(pi.source)
}

type AbortPipe struct{}

// An object to which consecutive iterator values are supplied via the Yield
// method (or an error via the Error method), sent via a channel
type Yield[T any] struct {
	sink chan result.Result[T]
}

func (y Yield[T]) Yield(t T) {
	if !safeSend(y.sink, result.Value(t)) {
		panic(AbortPipe{})
	}
}

func (y Yield[T]) Error(err error) {
	safeSend(y.sink, result.Error[T](err))
}

// An error type indicating that a Pipe iterator function has panicked
type PipePanic struct {
	panic interface{}
}

func (pp PipePanic) Error() string {
	return fmt.Sprintf("panic during pipe iterator: %#v", pp.panic)
}

// A type alias for a function taking a Yield object and returning
// an error.
type PipeFunc[T any] func(Yield[T]) error

func runPipe[T any](y Yield[T], activity PipeFunc[T]) {
	defer safeClose(y.sink)
	defer func() {
		if p := recover(); p != nil {
			if _, abort := p.(AbortPipe); !abort {
				y.Error(PipePanic{p})
			}
		}
	}()
	if err := activity(y); err != nil {
		y.Error(err)
	}
}

// Create an iterator from a function, passed in activity, that yields a sequence of results
// by making calls to Yield(), or Error(), on an object that will be passed to the function.
// The function is run in a separate goroutine, and its yielded results are sent over a channel
// to the iterator where they can be consumed in the normal way by calls to Next() and Value().
// Any error result yielded, or returned by the function, cause that error to be returned by Value()
// in the iterator.
//
// The activity parameter is of type PipeFunc[T], which is an alias for func(Yield[T]) error
func Pipe[T any](activity PipeFunc[T]) Iterator[T] {
	ch := make(chan result.Result[T])
	yield := Yield[T]{ch}
	go runPipe(yield, activity)
	return &PipeIter[T]{source: ch}
}

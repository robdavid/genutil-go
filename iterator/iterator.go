package iterator

import (
	"fmt"

	"github.com/robdavid/genutil-go/option"
	"github.com/robdavid/genutil-go/result"
)

// Generic iterator
type Iterator[T any] interface {
	// Set the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned when there are no more elements (the current value remains unchanged)
	Next() bool
	Value() (T, error) // Get the current iterator value, or an error if the last call to Next() has failed
	MustValue() T      // Get the current iterator value. Panics if the last call to Next() has failed.
}

// Wraps an iterator and adds a mapping function
type Func[T, U any] struct {
	base    Iterator[T]
	mapping func(T) (U, error)
	value   option.Option[result.Result[U]]
}

func (i *Func[T, U]) Next() bool {
	i.value.Clear()
	return i.base.Next()
}

func (i *Func[T, U]) Value() (u U, err error) {
	if i.value.IsEmpty() {
		var t T
		if t, err = i.base.Value(); err != nil {
			return
		}
		u, err = i.mapping(t)
		i.value.Set(result.Pair(u, err))
	} else {
		u, err = i.value.Ref().Pair()
	}
	return
}

func (i *Func[T, U]) MustValue() U {
	if u, err := i.Value(); err != nil {
		panic(err)
	} else {
		return u
	}
}

type SliceIter[T any] struct {
	slice []T
	index int
	value T
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

func (si *SliceIter[T]) MustValue() T {
	return si.value
}

func (si *SliceIter[T]) Value() (T, error) {
	return si.value, nil
}

// Makes an Iterator[T] from slice []T
func Slice[T any](slice []T) Iterator[T] {
	var t T
	return &SliceIter[T]{slice, 0, t}
}

// Makes an Iterator[T] from variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	var t T
	return &SliceIter[T]{elements, 0, t}
}

// Wraps an iterator with a mapping function that may return an error, producing a new iterator
func DoMap[T any, U any](iter Iterator[T], mapping func(T) (U, error)) Iterator[U] {
	return &Func[T, U]{iter, mapping, option.Empty[result.Result[U]]()}
}

// Wraps an iterator with a mapping function, producing a new iterator
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	action := func(t T) (U, error) {
		return mapping(t), nil
	}
	return &Func[T, U]{iter, action, option.Empty[result.Result[U]]()}
}

// Collects all elements from an iterator into a slice. If the iterator
// returns an error, this call will panic.
func MustCollect[T any](iter Iterator[T]) []T {
	if result, err := Collect(iter); err != nil {
		panic(err)
	} else {
		return result
	}
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
	source <-chan result.Result[T]
	value  result.Result[T]
}

func (pi *PipeIter[T]) Next() bool {
	var ok bool
	pi.value, ok = <-pi.source
	return ok
}

func (pi *PipeIter[T]) Value() (T, error) {
	return pi.value.Pair()
}

func (pi *PipeIter[T]) MustValue() T {
	if v, err := pi.value.Pair(); err != nil {
		panic(err)
	} else {
		return v
	}
}

// An object to which consecutive iterator values are supplied via the Yield
// method (or an error via the Error method), sent via a channel
type Yield[T any] struct {
	sink chan<- result.Result[T]
}

func (y Yield[T]) Yield(t T) {
	y.sink <- result.Value(t)
}

func (y Yield[T]) Error(err error) {
	y.sink <- result.Error[T](err)
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
	defer close(y.sink)
	defer func() {
		if p := recover(); p != nil {
			y.Error(PipePanic{p})
		}
	}()
	if err := activity(y); err != nil {
		y.Error(err)
	}
}

// Create an iterator from a function that yields a sequence of results
// by making calls to Yield(), or Error(), on an object that will be passed to the function.
// The function is run in a separate goroutine, and its yielded results are sent over a channel
// to the iterator where they can be consumed in the normal way by calls to Next() and Value().
// Any error result yielded, or returned by the function, cause that error to be returned by Value()
// in the iterator.
//
// PipeFunc[T] is an alias for func(Yield[T]) error
func Pipe[T any](activity PipeFunc[T]) Iterator[T] {
	ch := make(chan result.Result[T], 1)
	yield := Yield[T]{ch}
	go runPipe(yield, activity)
	return &PipeIter[T]{source: ch}
}

// Iterators and generators
package iterator

import (
	"fmt"

	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/option"
)

// Generic iterator
type Iterator[T any] interface {
	// Set the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned when there are no more elements (the current value remains unchanged)
	Next() bool
	Value() T       // Get the current iterator value, or an error if the last call to Next() has failed
	Abort()         // Stop the iterator; subsequent calls to Next() will return false.
	Chan() <-chan T // Return iterator as a channel
}

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
}

func (i *Func[T, U]) Next() bool {
	ok, value := i.mapping(i.base)
	if !ok {
		i.value.Clear()
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

func WrapFunc[T any, U any](iterator Iterator[T], f FuncNext[T, U]) Iterator[U] {
	return &Func[T, U]{base: iterator, mapping: f, value: option.Empty[U]()}
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

func wrapMap[T any, U any](mapping func(T) U) FuncNext[T, U] {
	return func(iterator Iterator[T]) (ok bool, value U) {
		if ok = iterator.Next(); ok {
			value = mapping(iterator.Value())
		}
		return
	}
}

// Wraps an iterator with a mapping function, producing a new iterator
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	action := wrapMap(func(t T) U { return mapping(t) })
	return WrapFunc(iter, action)
}

// Collects all elements from an iterator into a slice.
func Collect[T any](iter Iterator[T]) []T {
	result := make([]T, 0)
	for iter.Next() {
		result = append(result, iter.Value())
	}
	return result
}

// Collects all elements from an iterator of results into a result of slice of the iterator's underlying type
// If the iterator returns an error result at any point, this call will terminate and return that error in the
// result, along with the elements collected thus far.
func CollectResults[T any](iter Iterator[result.Result[T]]) result.Result[[]T] {
	collectResult := make([]T, 0)
	for iter.Next() {
		res := iter.Value()
		if res.IsError() {
			iter.Abort()
			return result.From(collectResult, res.GetErr())
		}
		collectResult = append(collectResult, res.Get())
	}
	return result.Value(collectResult)
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

type AbortGenerator struct{}

// An object to which consecutive iterator values are supplied via the Yield
// method (or an error via the Error method), sent via a channel
type Yield[T any] struct {
	sink chan T
}

func (y Yield[T]) Yield(t T) {
	if !safeSend(y.sink, t) {
		panic(AbortGenerator{})
	}
}

// An error type indicating that a Pipe iterator function has panicked
type GeneratorPanic struct {
	panic interface{}
}

func (pp GeneratorPanic) Error() string {
	return fmt.Sprintf("panic during pipe iterator: %#v", pp.panic)
}

// A type alias for a function taking a Yield object and returning
// an error.
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

// Create an iterator from a function, passed in activity, that yields a sequence of results
// by making calls to Yield(), on an object that will be passed to the function.
// The function is run in a separate goroutine, and its yielded results are sent over a channel
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

package iterator

import (
	"fmt"
	"iter"

	. "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
)

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

// Reset is the same as abort for this iterator
func (pi *genIter[T]) Reset() {
	safeClose(pi.source)
}

func (pi *genIter[T]) Size() IteratorSize {
	return NewSizeUnknown()
}

func (pi *genIter[T]) Seq() iter.Seq[T] {
	return Seq(pi)
}

func (pi *genIter[T]) SeqOK() bool {
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
	defer Handle(func(err error) { c.YieldError(err) })
	Check(activity(c))
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

// Chan takes a CoreIterator and produces a channel yielding
// values from the iterator.
func Chan[T any](itr CoreIterator[T]) (out chan T) {
	out = make(chan T)
	go func() {
		defer safeClose(out)
		if !itr.SeqOK() {
			for itr.Next() {
				if !safeSend(out, itr.Value()) {
					itr.Abort()
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

// Chan2 takes a CoreIterator2 and produces a channel yielding
// key and value pairs from the iterator.
func Chan2[K any, V any](itr CoreIterator2[K, V]) (out chan KeyValue[K, V]) {
	out = make(chan KeyValue[K, V])
	go func() {
		defer safeClose(out)
		if !itr.SeqOK() {
			for itr.Next() {
				if !safeSend(out, KVOf(itr.Key(), itr.Value())) {
					itr.Abort()
					break
				}
			}
		} else {
			for k, v := range itr.Seq2() {
				if !safeSend(out, KVOf(k, v)) {
					itr.Abort()
					break
				}
			}
		}
	}()
	return
}

package result

import (
	"fmt"

	"github.com/robdavid/genutil-go/tuple"
)

//import "github.com"

// Contains a value of type T, or an error
type Result[T any] struct {
	value T
	err   error
}

type Result2[A any, B any] struct{ Result[tuple.Tuple2[A, B]] }

// Creates a new non-error result
func Value[T any](t T) Result[T] {
	return Result[T]{t, nil}
}

// Creates an error result
func Error[T any](err error) Result[T] {
	var t T
	return Result[T]{t, err}
}

// Creates a result from both error and non-error
// values. Typical use case is to wrap the return values
// from a function that returns err,
// e.g.
//
//	res := result.From(os.Open("myfile"))
func From[T any](t T, err error) Result[T] {
	return Result[T]{t, err}
}

func From2[A any, B any](a A, b B, err error) Result2[A, B] {
	return Result2[A, B]{Result[tuple.Tuple2[A, B]]{tuple.Pair(a, b), err}}
}

// Returns a result as a pair of values, suitable to be
// used as function return values.
func (r *Result[T]) Return() (T, error) {
	return r.value, r.err
}

func (r *Result2[A, B]) Return() (A, B, error) {
	return r.value.First, r.value.Second, r.err
}

func (r *Result[T]) GetErr() error {
	return r.err
}

func (r *Result[T]) Get() T {
	return r.value
}

func (r *Result[T]) Must() T {
	if r.err != nil {
		panic(r.err)
	}
	return r.value
}

func (r *Result[T]) IsError() bool {
	return r.err != nil
}

func (r *Result[T]) IsValue() bool {
	return r.err == nil
}

func (r *Result[T]) String() string {
	if r.IsError() {
		return r.err.Error()
	} else if str, ok := any(r.value).(fmt.Stringer); ok {
		return str.String()
	} else if str, ok := any(&r.value).(fmt.Stringer); ok {
		return str.String()
	} else {
		return fmt.Sprintf("%v", r.value)
	}
}

func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		return Value(f(r.Get()))
	}
}

func MapErr[T, U any](r Result[T], f func(T) (U, error)) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		if u, err := f(r.Get()); err != nil {
			return Error[U](err)
		} else {
			return Value(u)
		}
	}
}

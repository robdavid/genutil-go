package result

import (
	"fmt"

	"github.com/robdavid/genutil-go/errors/handler"
)

//import "github.com"

// Contains a value of type T, or an error
type Result[T any] struct {
	value T
	err   error
}

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

func New[T any](t T, err error) *Result[T] {
	return &Result[T]{t, err}
}

func (r Result[T]) ToRef() *Result[T] {
	return &r
}

// Returns a result as a pair of values, suitable to be
// used as function return values.
func (r *Result[T]) Return() (T, error) {
	return r.value, r.err
}

func (r *Result[T]) GetErr() error {
	return r.err
}

func (r *Result[T]) Get() T {
	return r.value
}

func (r *Result[T]) Must() T {
	return handler.Must(r.value, r.err)
}

func (r *Result[T]) Try() T {
	return handler.Try(r.value, r.err)
}

// Perform a transformation to the error part of an error result.
// If the result is not an error, return the original result.
func (r *Result[T]) MapErr(f func(error) error) *Result[T] {
	if r.err != nil {
		res := From(r.value, f(r.err))
		return &res
	} else {
		return r
	}
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

//go:generate code-template --set max_params=9 result.tmpl

// Apply a map function `f` to the value part of a result `r`.
// If `r` is an error return an error result, with a zero value.
// Otherwise return a new successful result with the value transformed by `f`
func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		return Value(f(r.Get()))
	}
}

// Apply a map function `f`, which may return an error, to the value part of a result `r`.
// If `r` is an error, or `f` returns an error, return an error result, with a zero value
// Otherwise return a new successful result with the value transformed by `f`
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

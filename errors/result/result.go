// A value type that contains a value plus an error, typically used to represent the
// return value of a function, including its error component. It has convenience methods
// for constructing an instance from a function return, e.g.
//
//	r := result.From(os.Open(file))
package result

import (
	"fmt"
	"reflect"

	"github.com/robdavid/genutil-go/errors/handler"
)

var packageName string

// PackageName returns the pull path of this package
func PackageName() string {
	if packageName == "" {
		packageName = reflect.TypeOf(Result[bool]{}).PkgPath()
	}
	return packageName
}

// Contains a value of type T, and/or an error
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

// Like From, except returns a reference to the resulting
// Result value.
func New[T any](t T, err error) *Result[T] {
	return &Result[T]{t, err}
}

// Returns a reference to the given result object
func (r Result[T]) ToRef() *Result[T] {
	return &r
}

// Returns a result as a pair of values, suitable to be
// used as function return values.
func (r *Result[T]) Return() (T, error) {
	return r.value, r.err
}

// Returns the associated error, or nil if there is no error
func (r *Result[T]) GetErr() error {
	return r.err
}

// Returns the underlying value. Does not panic.
func (r *Result[T]) Get() T {
	return r.value
}

// Returns the underlying value, or panics if the error is
// non-nil
func (r *Result[T]) Must() T {
	return handler.Must(r.value, r.err)
}

// Returns the underlying value, or panics with a try value
// if the error is. This can be handled by the Catch or Handle
// functions in the handler package
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

// Returns true if the Result holds an error
func (r *Result[T]) IsError() bool {
	return r.err != nil
}

// Returns true if the Result contains no error
func (r *Result[T]) IsValue() bool {
	return r.err == nil
}

// Renders either the value or the error in the case that the
// Result contains an error.
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

func Then[T, U any](r Result[T], f func(T) U) Result[U] {
	return Map(r, f)
}

func AndThen[T, U any](r Result[T], f func(T) Result[U]) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		return f(r.Get())
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

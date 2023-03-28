// Some error handling functions and types to more ergonomically assist
// with writing tests against functions that may return errors
package test

import (
	"fmt"
	"strings"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/tuple"
)

// An interface implemented by multiple types in the "testing" package
type TestReporting interface {
	Error(args ...any)
	FailNow()
	Helper()
}

// A wrapper around result.Result that supports test assertions.
type TestableResult[T any] struct {
	result.Result[T]
}

func resultFrom[T any](value T, err error) TestableResult[T] {
	return TestableResult[T]{result.From(value, err)}
}

// Creates a TestableResult from a an error only,
// e.g.
//
//	r := test.Result0(os.Rename(oldfile,newfile))
func Result0(err error) *TestableResult[tuple.Tuple0] {
	return Result(tuple.Of0(), err)
}

// Creates a TestableResult from a an error only. Alias for Result0.
// e.g.
//
//	r := test.Status(os.Rename(oldfile,newfile))
func Status(err error) *TestableResult[tuple.Tuple0] {
	return Result0(err)
}

// Creates a TestableResult from a return value and an error
// e.g.
//
//	r := test.Result(os.Open(myfile))
func Result[T any](value T, err error) *TestableResult[T] {
	return &TestableResult[T]{result.From(value, err)}
}

// Checks if a given error value is nil. If not, report the
// error and fail the test.
func Check(t TestReporting, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// Either return the value of a TestableResult, or if
// an error has occurred, report it and fail the test.
func (r *TestableResult[T]) Must(t TestReporting) T {
	t.Helper()
	if r.IsError() {
		Check(t, r.GetErr())
	}
	return r.Get()
}

// Expects an error; marks the test as failed if the result is not an error
func (r *TestableResult[T]) Fails(t TestReporting) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	}
	return r
}

// Expects a specific error; marks the test as failed if the result is not the error
// provided in expected.
func (r *TestableResult[T]) FailsWith(t TestReporting, expected error) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	} else if expected != r.GetErr() {
		t.Error(fmt.Errorf("expected error '%s', but got '%s'", expected, r.GetErr()))
	}
	return r
}

// Expects a specific error; marks the test as failed if the result is not an
// error whose Error() return contains the string provided in expected.
func (r *TestableResult[T]) FailsContaining(t TestReporting, expected string) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	} else if !strings.Contains(r.GetErr().Error(), expected) {
		t.Error(fmt.Errorf("expected error to contain '%s', but was '%s'", expected, r.GetErr().Error()))
	}
	return r
}

// Reports any error encountered in a call to Try().
// Should be used as part of a defer call.
// e.g.
//
//	defer ReportErr(t)
//	f := Try(os.Open(myfile))
func ReportErr(t TestReporting) {
	t.Helper()
	if err := recover(); err != nil {
		if tryErr, ok := err.(handler.TryError); ok {
			t.Error(tryErr.Error)
		} else {
			panic(err)
		}
	}
}

//go:generate code-template test.tmpl

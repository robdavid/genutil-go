// Some error handling functions and types to more ergonomically assist
// with writing tests against functions that may return errors
package test

import (
	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
)

// An interface implemented by multiple types in the "testing" package
type TestReporting interface {
	Error(args ...any)
	FailNow()
}

// A wrapper around result.Result that supports test assertions.
type TestableResult[T any] struct {
	result.Result[T]
}

func resultFrom[T any](value T, err error) TestableResult[T] {
	return TestableResult[T]{result.From(value, err)}
}

// Creates a TestableResult from a return value and an error
// e.g.
//
//	r := Result(os.Open(myfile))
func Result[T any](value T, err error) *TestableResult[T] {
	return &TestableResult[T]{result.From(value, err)}
}

// Checks if a given error value is nil. If not, report the
// error and fail the test.
func Check(t TestReporting, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// Either return the value of a TestableResult, or if
// an error has occurred, report it and fail the test.
func (r *TestableResult[T]) Must(t TestReporting) T {
	if r.IsError() {
		Check(t, r.GetErr())
	}
	return r.Get()
}

// Reports any error encountered in a call to Try().
// Should be used as part of a defer call.
// e.g.
//
//	defer ReportErr(t)
//	f := Try(os.Open(myfile))
func ReportErr(t TestReporting) {
	if err := recover(); err != nil {
		if tryErr, ok := err.(handler.TryError); ok {
			t.Error(tryErr.Error)
		} else {
			panic(err)
		}
	}
}

//go:generate code-template test.tmpl

// Some error handling functions and type to more ergonomically assist
// with writing tests against functions that return errors
package test

import (
	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
)

type TestReporting interface {
	Error(args ...any)
	FailNow()
}

type TestableResult[T any] struct {
	result.Result[T]
}

func resultFrom[T any](value T, err error) TestableResult[T] {
	return TestableResult[T]{result.From(value, err)}
}

func Result[T any](value T, err error) *TestableResult[T] {
	return &TestableResult[T]{result.From(value, err)}
}

func Check(t TestReporting, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func (r *TestableResult[T]) Must(t TestReporting) T {
	if r.IsError() {
		Check(t, r.GetErr())
	}
	return r.Get()
}

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

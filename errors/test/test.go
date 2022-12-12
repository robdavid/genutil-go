// Some error handling functions and type to more ergonomically assist
// with writing tests against functions that return errors
package test

import (
	"testing"

	"github.com/robdavid/genutil-go/errors/result"
)

type TestableResult[T any] struct {
	result.Result[T]
}

func resultFrom[T any](value T, err error) TestableResult[T] {
	return TestableResult[T]{result.From(value, err)}
}

func Result[T any](value T, err error) *TestableResult[T] {
	return &TestableResult[T]{result.From(value, err)}
}

func Check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func (r *TestableResult[T]) Must(t *testing.T) T {
	if r.IsError() {
		Check(t, r.GetErr())
	}
	return r.Get()
}

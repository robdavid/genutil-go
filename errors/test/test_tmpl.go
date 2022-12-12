package test

import (
	"testing"

	"github.com/robdavid/genutil-go/tuple"
)

type TestableResult2[T1 any, T2 any] struct {
	TestableResult[tuple.Tuple2[T1, T2]]
}

func Result2[T1 any, T2 any](v1 T1, v2 T2, err error) *TestableResult2[T1, T2] {
	return &TestableResult2[T1, T2]{resultFrom(tuple.Of2(v1, v2), err)}
}

func (r *TestableResult2[T1, T2]) Must(t *testing.T) (T1, T2) {
	v := r.TestableResult.Must(t)
	return v.First, v.Second
}

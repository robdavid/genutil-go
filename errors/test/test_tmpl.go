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
	return v.Return()
}
type TestableResult3[T1 any, T2 any, T3 any] struct {
	TestableResult[tuple.Tuple3[T1, T2, T3]]
}

func Result3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error) *TestableResult3[T1, T2, T3] {
	return &TestableResult3[T1, T2, T3]{resultFrom(tuple.Of3(v1, v2, v3), err)}
}

func (r *TestableResult3[T1, T2, T3]) Must(t *testing.T) (T1, T2, T3) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult4[T1 any, T2 any, T3 any, T4 any] struct {
	TestableResult[tuple.Tuple4[T1, T2, T3, T4]]
}

func Result4[T1 any, T2 any, T3 any, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) *TestableResult4[T1, T2, T3, T4] {
	return &TestableResult4[T1, T2, T3, T4]{resultFrom(tuple.Of4(v1, v2, v3, v4), err)}
}

func (r *TestableResult4[T1, T2, T3, T4]) Must(t *testing.T) (T1, T2, T3, T4) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	TestableResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
}

func Result5[T1 any, T2 any, T3 any, T4 any, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, err error) *TestableResult5[T1, T2, T3, T4, T5] {
	return &TestableResult5[T1, T2, T3, T4, T5]{resultFrom(tuple.Of5(v1, v2, v3, v4, v5), err)}
}

func (r *TestableResult5[T1, T2, T3, T4, T5]) Must(t *testing.T) (T1, T2, T3, T4, T5) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct {
	TestableResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
}

func Result6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, err error) *TestableResult6[T1, T2, T3, T4, T5, T6] {
	return &TestableResult6[T1, T2, T3, T4, T5, T6]{resultFrom(tuple.Of6(v1, v2, v3, v4, v5, v6), err)}
}

func (r *TestableResult6[T1, T2, T3, T4, T5, T6]) Must(t *testing.T) (T1, T2, T3, T4, T5, T6) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct {
	TestableResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
}

func Result7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, err error) *TestableResult7[T1, T2, T3, T4, T5, T6, T7] {
	return &TestableResult7[T1, T2, T3, T4, T5, T6, T7]{resultFrom(tuple.Of7(v1, v2, v3, v4, v5, v6, v7), err)}
}

func (r *TestableResult7[T1, T2, T3, T4, T5, T6, T7]) Must(t *testing.T) (T1, T2, T3, T4, T5, T6, T7) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct {
	TestableResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
}

func Result8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, err error) *TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8]{resultFrom(tuple.Of8(v1, v2, v3, v4, v5, v6, v7, v8), err)}
}

func (r *TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8]) Must(t *testing.T) (T1, T2, T3, T4, T5, T6, T7, T8) {
	v := r.TestableResult.Must(t)
	return v.Return()
}
type TestableResult9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct {
	TestableResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
}

func Result9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9, err error) *TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return &TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{resultFrom(tuple.Of9(v1, v2, v3, v4, v5, v6, v7, v8, v9), err)}
}

func (r *TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Must(t *testing.T) (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	v := r.TestableResult.Must(t)
	return v.Return()
}

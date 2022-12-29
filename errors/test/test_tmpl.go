package test

import (
	"github.com/robdavid/genutil-go/tuple"
)

// A wrapper type for a TestableResult that contains a value of type type.Tuple2
type TestableResult2[T1 any, T2 any] struct {
	TestableResult[tuple.Tuple2[T1, T2]]
}

// A constructor for TestableResult2 from the values returned by a function that returns 2 values
// plus an error.
func Result2[T1 any, T2 any](v1 T1, v2 T2, err error) *TestableResult2[T1, T2] {
	return &TestableResult2[T1, T2]{resultFrom(tuple.Of2(v1, v2), err)}
}

// A variation of Must that returns 2 non-error values
func (r *TestableResult2[T1, T2]) Must2(t TestReporting) (T1, T2) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 2 non-error values
func (r *TestableResult2[T1, T2]) Try2() (T1, T2) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple3
type TestableResult3[T1 any, T2 any, T3 any] struct {
	TestableResult[tuple.Tuple3[T1, T2, T3]]
}

// A constructor for TestableResult3 from the values returned by a function that returns 3 values
// plus an error.
func Result3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, err error) *TestableResult3[T1, T2, T3] {
	return &TestableResult3[T1, T2, T3]{resultFrom(tuple.Of3(v1, v2, v3), err)}
}

// A variation of Must that returns 3 non-error values
func (r *TestableResult3[T1, T2, T3]) Must3(t TestReporting) (T1, T2, T3) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 3 non-error values
func (r *TestableResult3[T1, T2, T3]) Try3() (T1, T2, T3) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple4
type TestableResult4[T1 any, T2 any, T3 any, T4 any] struct {
	TestableResult[tuple.Tuple4[T1, T2, T3, T4]]
}

// A constructor for TestableResult4 from the values returned by a function that returns 4 values
// plus an error.
func Result4[T1 any, T2 any, T3 any, T4 any](v1 T1, v2 T2, v3 T3, v4 T4, err error) *TestableResult4[T1, T2, T3, T4] {
	return &TestableResult4[T1, T2, T3, T4]{resultFrom(tuple.Of4(v1, v2, v3, v4), err)}
}

// A variation of Must that returns 4 non-error values
func (r *TestableResult4[T1, T2, T3, T4]) Must4(t TestReporting) (T1, T2, T3, T4) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 4 non-error values
func (r *TestableResult4[T1, T2, T3, T4]) Try4() (T1, T2, T3, T4) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple5
type TestableResult5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	TestableResult[tuple.Tuple5[T1, T2, T3, T4, T5]]
}

// A constructor for TestableResult5 from the values returned by a function that returns 5 values
// plus an error.
func Result5[T1 any, T2 any, T3 any, T4 any, T5 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, err error) *TestableResult5[T1, T2, T3, T4, T5] {
	return &TestableResult5[T1, T2, T3, T4, T5]{resultFrom(tuple.Of5(v1, v2, v3, v4, v5), err)}
}

// A variation of Must that returns 5 non-error values
func (r *TestableResult5[T1, T2, T3, T4, T5]) Must5(t TestReporting) (T1, T2, T3, T4, T5) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 5 non-error values
func (r *TestableResult5[T1, T2, T3, T4, T5]) Try5() (T1, T2, T3, T4, T5) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple6
type TestableResult6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct {
	TestableResult[tuple.Tuple6[T1, T2, T3, T4, T5, T6]]
}

// A constructor for TestableResult6 from the values returned by a function that returns 6 values
// plus an error.
func Result6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, err error) *TestableResult6[T1, T2, T3, T4, T5, T6] {
	return &TestableResult6[T1, T2, T3, T4, T5, T6]{resultFrom(tuple.Of6(v1, v2, v3, v4, v5, v6), err)}
}

// A variation of Must that returns 6 non-error values
func (r *TestableResult6[T1, T2, T3, T4, T5, T6]) Must6(t TestReporting) (T1, T2, T3, T4, T5, T6) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 6 non-error values
func (r *TestableResult6[T1, T2, T3, T4, T5, T6]) Try6() (T1, T2, T3, T4, T5, T6) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple7
type TestableResult7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct {
	TestableResult[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]]
}

// A constructor for TestableResult7 from the values returned by a function that returns 7 values
// plus an error.
func Result7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, err error) *TestableResult7[T1, T2, T3, T4, T5, T6, T7] {
	return &TestableResult7[T1, T2, T3, T4, T5, T6, T7]{resultFrom(tuple.Of7(v1, v2, v3, v4, v5, v6, v7), err)}
}

// A variation of Must that returns 7 non-error values
func (r *TestableResult7[T1, T2, T3, T4, T5, T6, T7]) Must7(t TestReporting) (T1, T2, T3, T4, T5, T6, T7) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 7 non-error values
func (r *TestableResult7[T1, T2, T3, T4, T5, T6, T7]) Try7() (T1, T2, T3, T4, T5, T6, T7) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple8
type TestableResult8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct {
	TestableResult[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]]
}

// A constructor for TestableResult8 from the values returned by a function that returns 8 values
// plus an error.
func Result8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, err error) *TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8]{resultFrom(tuple.Of8(v1, v2, v3, v4, v5, v6, v7, v8), err)}
}

// A variation of Must that returns 8 non-error values
func (r *TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8]) Must8(t TestReporting) (T1, T2, T3, T4, T5, T6, T7, T8) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 8 non-error values
func (r *TestableResult8[T1, T2, T3, T4, T5, T6, T7, T8]) Try8() (T1, T2, T3, T4, T5, T6, T7, T8) {
	v := r.Try()
	return v.Return()
}

// A wrapper type for a TestableResult that contains a value of type type.Tuple9
type TestableResult9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct {
	TestableResult[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]]
}

// A constructor for TestableResult9 from the values returned by a function that returns 9 values
// plus an error.
func Result9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](v1 T1, v2 T2, v3 T3, v4 T4, v5 T5, v6 T6, v7 T7, v8 T8, v9 T9, err error) *TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return &TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{resultFrom(tuple.Of9(v1, v2, v3, v4, v5, v6, v7, v8, v9), err)}
}

// A variation of Must that returns 9 non-error values
func (r *TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Must9(t TestReporting) (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	v := r.Must(t)
	return v.Return()
}

// A variation of Try that returns 9 non-error values
func (r *TestableResult9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Try9() (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	v := r.Try()
	return v.Return()
}


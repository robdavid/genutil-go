package result
import "github.com/robdavid/genutil-go/tuple"


// Arity 2 result values
type Result2[T1 any, T2 any] struct{ Result[tuple.Tuple2[T1, T2]] }

func Value2[T1 any, T2 any](v tuple.Tuple2[T1, T2]) Result2[T1, T2] {
	return Result2[T1, T2]{Value(v)}
}

func Error2[T1 any, T2 any](err error) Result2[T1, T2] {
	var zero tuple.Tuple2[T1, T2]
	return Result2[T1, T2]{From(zero, err)}
}

func From2[T1 any, T2 any](t1 T1, t2 T2, err error) Result2[T1, T2] {
	return Result2[T1, T2]{From(tuple.Of2(t1, t2), err)}
}

func New2[T1 any, T2 any](t1 T1, t2 T2, err error) *Result2[T1, T2] {
	r := From2(t1, t2, err)
	return &r
}

func (r *Result2[T1, T2]) Return() (T1, T2, error) {
	return r.value.First, r.value.Second, r.err
}

func (r Result2[T1, T2]) ToRef() *Result2[T1, T2] {
	return &r
}

// Arity 3 result values
type Result3[T1 any, T2 any, T3 any] struct{ Result[tuple.Tuple3[T1, T2, T3]] }

func Value3[T1 any, T2 any, T3 any](v tuple.Tuple3[T1, T2, T3]) Result3[T1, T2, T3] {
	return Result3[T1, T2, T3]{Value(v)}
}

func Error3[T1 any, T2 any, T3 any](err error) Result3[T1, T2, T3] {
	var zero tuple.Tuple3[T1, T2, T3]
	return Result3[T1, T2, T3]{From(zero, err)}
}

func From3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, err error) Result3[T1, T2, T3] {
	return Result3[T1, T2, T3]{From(tuple.Of3(t1, t2, t3), err)}
}

func New3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, err error) *Result3[T1, T2, T3] {
	r := From3(t1, t2, t3, err)
	return &r
}

func (r *Result3[T1, T2, T3]) Return() (T1, T2, T3, error) {
	return r.value.First, r.value.Second, r.value.Third, r.err
}

func (r Result3[T1, T2, T3]) ToRef() *Result3[T1, T2, T3] {
	return &r
}

// Arity 4 result values
type Result4[T1 any, T2 any, T3 any, T4 any] struct{ Result[tuple.Tuple4[T1, T2, T3, T4]] }

func Value4[T1 any, T2 any, T3 any, T4 any](v tuple.Tuple4[T1, T2, T3, T4]) Result4[T1, T2, T3, T4] {
	return Result4[T1, T2, T3, T4]{Value(v)}
}

func Error4[T1 any, T2 any, T3 any, T4 any](err error) Result4[T1, T2, T3, T4] {
	var zero tuple.Tuple4[T1, T2, T3, T4]
	return Result4[T1, T2, T3, T4]{From(zero, err)}
}

func From4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) Result4[T1, T2, T3, T4] {
	return Result4[T1, T2, T3, T4]{From(tuple.Of4(t1, t2, t3, t4), err)}
}

func New4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) *Result4[T1, T2, T3, T4] {
	r := From4(t1, t2, t3, t4, err)
	return &r
}

func (r *Result4[T1, T2, T3, T4]) Return() (T1, T2, T3, T4, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.err
}

func (r Result4[T1, T2, T3, T4]) ToRef() *Result4[T1, T2, T3, T4] {
	return &r
}

// Arity 5 result values
type Result5[T1 any, T2 any, T3 any, T4 any, T5 any] struct{ Result[tuple.Tuple5[T1, T2, T3, T4, T5]] }

func Value5[T1 any, T2 any, T3 any, T4 any, T5 any](v tuple.Tuple5[T1, T2, T3, T4, T5]) Result5[T1, T2, T3, T4, T5] {
	return Result5[T1, T2, T3, T4, T5]{Value(v)}
}

func Error5[T1 any, T2 any, T3 any, T4 any, T5 any](err error) Result5[T1, T2, T3, T4, T5] {
	var zero tuple.Tuple5[T1, T2, T3, T4, T5]
	return Result5[T1, T2, T3, T4, T5]{From(zero, err)}
}

func From5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) Result5[T1, T2, T3, T4, T5] {
	return Result5[T1, T2, T3, T4, T5]{From(tuple.Of5(t1, t2, t3, t4, t5), err)}
}

func New5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) *Result5[T1, T2, T3, T4, T5] {
	r := From5(t1, t2, t3, t4, t5, err)
	return &r
}

func (r *Result5[T1, T2, T3, T4, T5]) Return() (T1, T2, T3, T4, T5, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.err
}

func (r Result5[T1, T2, T3, T4, T5]) ToRef() *Result5[T1, T2, T3, T4, T5] {
	return &r
}

// Arity 6 result values
type Result6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct{ Result[tuple.Tuple6[T1, T2, T3, T4, T5, T6]] }

func Value6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](v tuple.Tuple6[T1, T2, T3, T4, T5, T6]) Result6[T1, T2, T3, T4, T5, T6] {
	return Result6[T1, T2, T3, T4, T5, T6]{Value(v)}
}

func Error6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](err error) Result6[T1, T2, T3, T4, T5, T6] {
	var zero tuple.Tuple6[T1, T2, T3, T4, T5, T6]
	return Result6[T1, T2, T3, T4, T5, T6]{From(zero, err)}
}

func From6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, err error) Result6[T1, T2, T3, T4, T5, T6] {
	return Result6[T1, T2, T3, T4, T5, T6]{From(tuple.Of6(t1, t2, t3, t4, t5, t6), err)}
}

func New6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, err error) *Result6[T1, T2, T3, T4, T5, T6] {
	r := From6(t1, t2, t3, t4, t5, t6, err)
	return &r
}

func (r *Result6[T1, T2, T3, T4, T5, T6]) Return() (T1, T2, T3, T4, T5, T6, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.err
}

func (r Result6[T1, T2, T3, T4, T5, T6]) ToRef() *Result6[T1, T2, T3, T4, T5, T6] {
	return &r
}

// Arity 7 result values
type Result7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct{ Result[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]] }

func Value7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](v tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]) Result7[T1, T2, T3, T4, T5, T6, T7] {
	return Result7[T1, T2, T3, T4, T5, T6, T7]{Value(v)}
}

func Error7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](err error) Result7[T1, T2, T3, T4, T5, T6, T7] {
	var zero tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]
	return Result7[T1, T2, T3, T4, T5, T6, T7]{From(zero, err)}
}

func From7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, err error) Result7[T1, T2, T3, T4, T5, T6, T7] {
	return Result7[T1, T2, T3, T4, T5, T6, T7]{From(tuple.Of7(t1, t2, t3, t4, t5, t6, t7), err)}
}

func New7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, err error) *Result7[T1, T2, T3, T4, T5, T6, T7] {
	r := From7(t1, t2, t3, t4, t5, t6, t7, err)
	return &r
}

func (r *Result7[T1, T2, T3, T4, T5, T6, T7]) Return() (T1, T2, T3, T4, T5, T6, T7, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.err
}

func (r Result7[T1, T2, T3, T4, T5, T6, T7]) ToRef() *Result7[T1, T2, T3, T4, T5, T6, T7] {
	return &r
}

// Arity 8 result values
type Result8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct{ Result[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]] }

func Value8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](v tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{Value(v)}
}

func Error8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](err error) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	var zero tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{From(zero, err)}
}

func From8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, err error) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{From(tuple.Of8(t1, t2, t3, t4, t5, t6, t7, t8), err)}
}

func New8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, err error) *Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	r := From8(t1, t2, t3, t4, t5, t6, t7, t8, err)
	return &r
}

func (r *Result8[T1, T2, T3, T4, T5, T6, T7, T8]) Return() (T1, T2, T3, T4, T5, T6, T7, T8, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.value.Eighth, r.err
}

func (r Result8[T1, T2, T3, T4, T5, T6, T7, T8]) ToRef() *Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &r
}

// Arity 9 result values
type Result9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct{ Result[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]] }

func Value9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](v tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{Value(v)}
}

func Error9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](err error) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	var zero tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{From(zero, err)}
}

func From9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, err error) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{From(tuple.Of9(t1, t2, t3, t4, t5, t6, t7, t8, t9), err)}
}

func New9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, err error) *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	r := From9(t1, t2, t3, t4, t5, t6, t7, t8, t9, err)
	return &r
}

func (r *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Return() (T1, T2, T3, T4, T5, T6, T7, T8, T9, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.value.Eighth, r.value.Ninth, r.err
}

func (r Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) ToRef() *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return &r
}


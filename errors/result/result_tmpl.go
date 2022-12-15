package result
import "github.com/robdavid/genutil-go/tuple"


// A wrapper around Result that contains a Tuple2 value.
type Result2[T1 any, T2 any] struct{ 
	Result[tuple.Tuple2[T1, T2]] 
}

// A non-error constructor that builds a Result2 value from 
// 2 parameters.
func Value2[T1 any, T2 any](v tuple.Tuple2[T1, T2]) Result2[T1, T2] {
	return Result2[T1, T2]{Value(v)}
}

// A non-error constructor that builds a Result2 value from 
// an error parameter
func Error2[T1 any, T2 any](err error) Result2[T1, T2] {
	var zero tuple.Tuple2[T1, T2]
	return Result2[T1, T2]{From(zero, err)}
}

// A constructor that builds a Result2 from 2 parameters and an error value.
// Can be used to create a Result2 from a function that returns 2 values
// and an error, as in:
//
//   result.From2(functionReturning2ParamsAndError())
func From2[T1 any, T2 any](t1 T1, t2 T2, err error) Result2[T1, T2] {
	return Result2[T1, T2]{From(tuple.Of2(t1, t2), err)}
}

// Like From2 except that a reference to the constructed Result2 is returned
func New2[T1 any, T2 any](t1 T1, t2 T2, err error) *Result2[T1, T2] {
	r := From2(t1, t2, err)
	return &r
}

// Returns the 2 values and error held in the result.
func (r *Result2[T1, T2]) Return() (T1, T2, error) {
	return r.value.First, r.value.Second, r.err
}

// Returns the underlying tuple value as a sequence of 2 elements.
// Panics with a try error if the result has an error
func (r *Result2[T1, T2]) Try2() (T1, T2) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 2 elements.
// Panics if the result has an error
func (r *Result2[T1, T2]) Must2() (T1, T2) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result2 
func (r Result2[T1, T2]) ToRef() *Result2[T1, T2] {
	return &r
}

// A wrapper around Result that contains a Tuple3 value.
type Result3[T1 any, T2 any, T3 any] struct{ 
	Result[tuple.Tuple3[T1, T2, T3]] 
}

// A non-error constructor that builds a Result3 value from 
// 3 parameters.
func Value3[T1 any, T2 any, T3 any](v tuple.Tuple3[T1, T2, T3]) Result3[T1, T2, T3] {
	return Result3[T1, T2, T3]{Value(v)}
}

// A non-error constructor that builds a Result3 value from 
// an error parameter
func Error3[T1 any, T2 any, T3 any](err error) Result3[T1, T2, T3] {
	var zero tuple.Tuple3[T1, T2, T3]
	return Result3[T1, T2, T3]{From(zero, err)}
}

// A constructor that builds a Result3 from 3 parameters and an error value.
// Can be used to create a Result3 from a function that returns 3 values
// and an error, as in:
//
//   result.From3(functionReturning3ParamsAndError())
func From3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, err error) Result3[T1, T2, T3] {
	return Result3[T1, T2, T3]{From(tuple.Of3(t1, t2, t3), err)}
}

// Like From3 except that a reference to the constructed Result3 is returned
func New3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3, err error) *Result3[T1, T2, T3] {
	r := From3(t1, t2, t3, err)
	return &r
}

// Returns the 3 values and error held in the result.
func (r *Result3[T1, T2, T3]) Return() (T1, T2, T3, error) {
	return r.value.First, r.value.Second, r.value.Third, r.err
}

// Returns the underlying tuple value as a sequence of 3 elements.
// Panics with a try error if the result has an error
func (r *Result3[T1, T2, T3]) Try3() (T1, T2, T3) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 3 elements.
// Panics if the result has an error
func (r *Result3[T1, T2, T3]) Must3() (T1, T2, T3) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result3 
func (r Result3[T1, T2, T3]) ToRef() *Result3[T1, T2, T3] {
	return &r
}

// A wrapper around Result that contains a Tuple4 value.
type Result4[T1 any, T2 any, T3 any, T4 any] struct{ 
	Result[tuple.Tuple4[T1, T2, T3, T4]] 
}

// A non-error constructor that builds a Result4 value from 
// 4 parameters.
func Value4[T1 any, T2 any, T3 any, T4 any](v tuple.Tuple4[T1, T2, T3, T4]) Result4[T1, T2, T3, T4] {
	return Result4[T1, T2, T3, T4]{Value(v)}
}

// A non-error constructor that builds a Result4 value from 
// an error parameter
func Error4[T1 any, T2 any, T3 any, T4 any](err error) Result4[T1, T2, T3, T4] {
	var zero tuple.Tuple4[T1, T2, T3, T4]
	return Result4[T1, T2, T3, T4]{From(zero, err)}
}

// A constructor that builds a Result4 from 4 parameters and an error value.
// Can be used to create a Result4 from a function that returns 4 values
// and an error, as in:
//
//   result.From4(functionReturning4ParamsAndError())
func From4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) Result4[T1, T2, T3, T4] {
	return Result4[T1, T2, T3, T4]{From(tuple.Of4(t1, t2, t3, t4), err)}
}

// Like From4 except that a reference to the constructed Result4 is returned
func New4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4, err error) *Result4[T1, T2, T3, T4] {
	r := From4(t1, t2, t3, t4, err)
	return &r
}

// Returns the 4 values and error held in the result.
func (r *Result4[T1, T2, T3, T4]) Return() (T1, T2, T3, T4, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.err
}

// Returns the underlying tuple value as a sequence of 4 elements.
// Panics with a try error if the result has an error
func (r *Result4[T1, T2, T3, T4]) Try4() (T1, T2, T3, T4) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 4 elements.
// Panics if the result has an error
func (r *Result4[T1, T2, T3, T4]) Must4() (T1, T2, T3, T4) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result4 
func (r Result4[T1, T2, T3, T4]) ToRef() *Result4[T1, T2, T3, T4] {
	return &r
}

// A wrapper around Result that contains a Tuple5 value.
type Result5[T1 any, T2 any, T3 any, T4 any, T5 any] struct{ 
	Result[tuple.Tuple5[T1, T2, T3, T4, T5]] 
}

// A non-error constructor that builds a Result5 value from 
// 5 parameters.
func Value5[T1 any, T2 any, T3 any, T4 any, T5 any](v tuple.Tuple5[T1, T2, T3, T4, T5]) Result5[T1, T2, T3, T4, T5] {
	return Result5[T1, T2, T3, T4, T5]{Value(v)}
}

// A non-error constructor that builds a Result5 value from 
// an error parameter
func Error5[T1 any, T2 any, T3 any, T4 any, T5 any](err error) Result5[T1, T2, T3, T4, T5] {
	var zero tuple.Tuple5[T1, T2, T3, T4, T5]
	return Result5[T1, T2, T3, T4, T5]{From(zero, err)}
}

// A constructor that builds a Result5 from 5 parameters and an error value.
// Can be used to create a Result5 from a function that returns 5 values
// and an error, as in:
//
//   result.From5(functionReturning5ParamsAndError())
func From5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) Result5[T1, T2, T3, T4, T5] {
	return Result5[T1, T2, T3, T4, T5]{From(tuple.Of5(t1, t2, t3, t4, t5), err)}
}

// Like From5 except that a reference to the constructed Result5 is returned
func New5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, err error) *Result5[T1, T2, T3, T4, T5] {
	r := From5(t1, t2, t3, t4, t5, err)
	return &r
}

// Returns the 5 values and error held in the result.
func (r *Result5[T1, T2, T3, T4, T5]) Return() (T1, T2, T3, T4, T5, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.err
}

// Returns the underlying tuple value as a sequence of 5 elements.
// Panics with a try error if the result has an error
func (r *Result5[T1, T2, T3, T4, T5]) Try5() (T1, T2, T3, T4, T5) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 5 elements.
// Panics if the result has an error
func (r *Result5[T1, T2, T3, T4, T5]) Must5() (T1, T2, T3, T4, T5) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result5 
func (r Result5[T1, T2, T3, T4, T5]) ToRef() *Result5[T1, T2, T3, T4, T5] {
	return &r
}

// A wrapper around Result that contains a Tuple6 value.
type Result6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct{ 
	Result[tuple.Tuple6[T1, T2, T3, T4, T5, T6]] 
}

// A non-error constructor that builds a Result6 value from 
// 6 parameters.
func Value6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](v tuple.Tuple6[T1, T2, T3, T4, T5, T6]) Result6[T1, T2, T3, T4, T5, T6] {
	return Result6[T1, T2, T3, T4, T5, T6]{Value(v)}
}

// A non-error constructor that builds a Result6 value from 
// an error parameter
func Error6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](err error) Result6[T1, T2, T3, T4, T5, T6] {
	var zero tuple.Tuple6[T1, T2, T3, T4, T5, T6]
	return Result6[T1, T2, T3, T4, T5, T6]{From(zero, err)}
}

// A constructor that builds a Result6 from 6 parameters and an error value.
// Can be used to create a Result6 from a function that returns 6 values
// and an error, as in:
//
//   result.From6(functionReturning6ParamsAndError())
func From6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, err error) Result6[T1, T2, T3, T4, T5, T6] {
	return Result6[T1, T2, T3, T4, T5, T6]{From(tuple.Of6(t1, t2, t3, t4, t5, t6), err)}
}

// Like From6 except that a reference to the constructed Result6 is returned
func New6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, err error) *Result6[T1, T2, T3, T4, T5, T6] {
	r := From6(t1, t2, t3, t4, t5, t6, err)
	return &r
}

// Returns the 6 values and error held in the result.
func (r *Result6[T1, T2, T3, T4, T5, T6]) Return() (T1, T2, T3, T4, T5, T6, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.err
}

// Returns the underlying tuple value as a sequence of 6 elements.
// Panics with a try error if the result has an error
func (r *Result6[T1, T2, T3, T4, T5, T6]) Try6() (T1, T2, T3, T4, T5, T6) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 6 elements.
// Panics if the result has an error
func (r *Result6[T1, T2, T3, T4, T5, T6]) Must6() (T1, T2, T3, T4, T5, T6) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result6 
func (r Result6[T1, T2, T3, T4, T5, T6]) ToRef() *Result6[T1, T2, T3, T4, T5, T6] {
	return &r
}

// A wrapper around Result that contains a Tuple7 value.
type Result7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct{ 
	Result[tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]] 
}

// A non-error constructor that builds a Result7 value from 
// 7 parameters.
func Value7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](v tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]) Result7[T1, T2, T3, T4, T5, T6, T7] {
	return Result7[T1, T2, T3, T4, T5, T6, T7]{Value(v)}
}

// A non-error constructor that builds a Result7 value from 
// an error parameter
func Error7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](err error) Result7[T1, T2, T3, T4, T5, T6, T7] {
	var zero tuple.Tuple7[T1, T2, T3, T4, T5, T6, T7]
	return Result7[T1, T2, T3, T4, T5, T6, T7]{From(zero, err)}
}

// A constructor that builds a Result7 from 7 parameters and an error value.
// Can be used to create a Result7 from a function that returns 7 values
// and an error, as in:
//
//   result.From7(functionReturning7ParamsAndError())
func From7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, err error) Result7[T1, T2, T3, T4, T5, T6, T7] {
	return Result7[T1, T2, T3, T4, T5, T6, T7]{From(tuple.Of7(t1, t2, t3, t4, t5, t6, t7), err)}
}

// Like From7 except that a reference to the constructed Result7 is returned
func New7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, err error) *Result7[T1, T2, T3, T4, T5, T6, T7] {
	r := From7(t1, t2, t3, t4, t5, t6, t7, err)
	return &r
}

// Returns the 7 values and error held in the result.
func (r *Result7[T1, T2, T3, T4, T5, T6, T7]) Return() (T1, T2, T3, T4, T5, T6, T7, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.err
}

// Returns the underlying tuple value as a sequence of 7 elements.
// Panics with a try error if the result has an error
func (r *Result7[T1, T2, T3, T4, T5, T6, T7]) Try7() (T1, T2, T3, T4, T5, T6, T7) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 7 elements.
// Panics if the result has an error
func (r *Result7[T1, T2, T3, T4, T5, T6, T7]) Must7() (T1, T2, T3, T4, T5, T6, T7) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result7 
func (r Result7[T1, T2, T3, T4, T5, T6, T7]) ToRef() *Result7[T1, T2, T3, T4, T5, T6, T7] {
	return &r
}

// A wrapper around Result that contains a Tuple8 value.
type Result8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct{ 
	Result[tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]] 
}

// A non-error constructor that builds a Result8 value from 
// 8 parameters.
func Value8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](v tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{Value(v)}
}

// A non-error constructor that builds a Result8 value from 
// an error parameter
func Error8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](err error) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	var zero tuple.Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{From(zero, err)}
}

// A constructor that builds a Result8 from 8 parameters and an error value.
// Can be used to create a Result8 from a function that returns 8 values
// and an error, as in:
//
//   result.From8(functionReturning8ParamsAndError())
func From8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, err error) Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Result8[T1, T2, T3, T4, T5, T6, T7, T8]{From(tuple.Of8(t1, t2, t3, t4, t5, t6, t7, t8), err)}
}

// Like From8 except that a reference to the constructed Result8 is returned
func New8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, err error) *Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	r := From8(t1, t2, t3, t4, t5, t6, t7, t8, err)
	return &r
}

// Returns the 8 values and error held in the result.
func (r *Result8[T1, T2, T3, T4, T5, T6, T7, T8]) Return() (T1, T2, T3, T4, T5, T6, T7, T8, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.value.Eighth, r.err
}

// Returns the underlying tuple value as a sequence of 8 elements.
// Panics with a try error if the result has an error
func (r *Result8[T1, T2, T3, T4, T5, T6, T7, T8]) Try8() (T1, T2, T3, T4, T5, T6, T7, T8) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 8 elements.
// Panics if the result has an error
func (r *Result8[T1, T2, T3, T4, T5, T6, T7, T8]) Must8() (T1, T2, T3, T4, T5, T6, T7, T8) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result8 
func (r Result8[T1, T2, T3, T4, T5, T6, T7, T8]) ToRef() *Result8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return &r
}

// A wrapper around Result that contains a Tuple9 value.
type Result9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct{ 
	Result[tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]] 
}

// A non-error constructor that builds a Result9 value from 
// 9 parameters.
func Value9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](v tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{Value(v)}
}

// A non-error constructor that builds a Result9 value from 
// an error parameter
func Error9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](err error) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	var zero tuple.Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{From(zero, err)}
}

// A constructor that builds a Result9 from 9 parameters and an error value.
// Can be used to create a Result9 from a function that returns 9 values
// and an error, as in:
//
//   result.From9(functionReturning9ParamsAndError())
func From9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, err error) Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{From(tuple.Of9(t1, t2, t3, t4, t5, t6, t7, t8, t9), err)}
}

// Like From9 except that a reference to the constructed Result9 is returned
func New9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9, err error) *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	r := From9(t1, t2, t3, t4, t5, t6, t7, t8, t9, err)
	return &r
}

// Returns the 9 values and error held in the result.
func (r *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Return() (T1, T2, T3, T4, T5, T6, T7, T8, T9, error) {
	return r.value.First, r.value.Second, r.value.Third, r.value.Forth, r.value.Fifth, r.value.Sixth, r.value.Seventh, r.value.Eighth, r.value.Ninth, r.err
}

// Returns the underlying tuple value as a sequence of 9 elements.
// Panics with a try error if the result has an error
func (r *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Try9() (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	v := r.Try()
	return v.Return()
}

// Returns the underlying tuple value as a sequence of 9 elements.
// Panics if the result has an error
func (r *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Must9() (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	v := r.Must()
	return v.Return()
}

// Returns a reference to the Result9 
func (r Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) ToRef() *Result9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return &r
}


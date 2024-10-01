package option

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	eh "github.com/robdavid/genutil-go/errors/handler"
)

var ErrOptionIsEmpty = errors.New("option is empty")

// A container for a value that might be empty
type Option[T any] struct {
	value    T
	nonEmpty bool
}

// Creates a non-empty option from a value
func Value[T any](v T) Option[T] {
	return Option[T]{v, true}
}

// From creates an option from a nilable value,
// which will be empty if the value is nil.
func From[T any](v T) Option[T] {
	return Option[T]{v, !isNil(v)}
}

// Safe is similar to From except a nil slice
// is considered equivalent to an empty slice.
func Safe[T any](v T) Option[T] {
	return Option[T]{v, !isUnsafe(v)}
}

// Creates an option from a pointer to a value; a nil
// pointer results in an empty option
func Ref[T any](v *T) *Option[T] {
	if v == nil {
		var zero T
		return &Option[T]{zero, false}
	} else {
		return &Option[T]{*v, true}
	}
}

// Create an empty option
func Empty[T any]() Option[T] {
	return Option[T]{}
}

// EmptyRef creates a reference to an empty option.
func EmptyRef[T any]() *Option[T] {
	return &Option[T]{}
}

// New creates a reference to a new zero value Option
func New[T any]() *Option[T] {
	var zero T
	return &Option[T]{zero, true}
}

// Returns true if the option is empty and has no value
func (o *Option[T]) IsEmpty() bool {
	return !o.nonEmpty
}

// Returns true if the option is non-empty and has a value
func (o *Option[T]) HasValue() bool {
	return o.nonEmpty
}

// Returns true if the passed value is nil
func isNil[T any](value T) bool {
	typ := reflect.TypeOf(value)
	if typ == nil { // Returns nil if nil interface
		return true
	}
	switch typ.Kind() {
	case reflect.Pointer, reflect.Chan, reflect.Slice, reflect.Interface, reflect.Map, reflect.Func:
		return reflect.ValueOf(value).IsNil()
	default:
		return false
	}
}

// Returns true if the passed value is an "unsafe" nil; e.g. a nil slice is
// considered safe as it is equivalent to an empty slice in terms of safe operations.
func isUnsafe[T any](value T) bool {
	typ := reflect.TypeOf(value)
	if typ == nil { // Returns nil if nil interface
		return true
	}
	switch typ.Kind() {
	case reflect.Pointer, reflect.Chan, reflect.Interface, reflect.Map, reflect.Func:
		return reflect.ValueOf(value).IsNil()
	default:
		return false
	}
}

// Get the options' value. If the option is empty, this call will panic
func (o Option[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return o.value
	}
}

// Get a pointer to the options' value. If the option is empty, this call will panic
func (o *Option[T]) GetRef() *T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return &o.value
	}
}

// Returns a pointer to the value in the option. If the value is empty,
// nil will be returned.
func (o *Option[T]) Ref() *T {
	if o.IsEmpty() {
		return nil
	} else {
		return &o.value
	}
}

// Either return a value if non-empty, or return the value's
// zero value. Note this call does not discriminate between
// an empty option and an option that contains the zero value.
func (o Option[T]) GetOrZero() T {
	return o.value // o.value is set to the zero value if o.nonEmpty is false
}

// Either return a pointer to the value if non-empty, or a pointer to a
// zero value. Note this call does not discriminate between
// an empty option and an option that contains the zero value.
func (o *Option[T]) GetOrZeroRef() *T {
	return &o.value // o.value is set to the zero value if o.nonEmpty is false
}

// Either return a value if non-empty, or return the default
// provided in the parameter.
func (o Option[T]) GetOr(def T) T {
	if o.IsEmpty() {
		return def
	} else {
		return o.value
	}
}

// Either return a pointer to the value if non-empty, or return the
// default pointer provided in the parameter.
func (o *Option[T]) GetOrRef(def *T) *T {
	if o.IsEmpty() {
		return def
	} else {
		return &o.value
	}
}

// Return the current value and a boolean flag which is true if the option is non-empty.
// If the option is empty. the value returned will be the zero value.
func (o Option[T]) GetOK() (T, bool) {
	return o.value, o.HasValue()
}

// Return a pointer to the current value and a boolean flag which is true if
// the option is non-empty. If the option is empty the pointer returned will
// be nil.
func (o *Option[T]) RefOK() (*T, bool) {
	if o.IsEmpty() {
		return nil, false
	} else {
		return &o.value, true
	}
}

// Try returns the option's non-empty value. If the option is empty,
// this call will panic with a try value that can be caught with Catch or Handle in
// errors/handler package.
// e.g.
//
//	 func tryOption() (err error) {
//		  defer handler.Catch(&err) // err set to ErrOptionIsEmpty
//		  ov := Empty[int]()
//		  v := ov.Try()
//	 }
func (o Option[T]) Try() T {
	return o.TryErr(ErrOptionIsEmpty)
}

// TryErr returns the option's non-empty value. If the option is empty,
// this call will panic with a try value, wapping the error supplied in err. This panic
// can be caught with Catch or Handle. If err is nil, there will be no panic and a
// zero value will be returned.
// e.g.
//
//	func tryOption() (err error) {
//	  myerr := errors.New("test error")
//	  defer handler.Catch(&err) // err set to myerr
//	  ov := Empty[int]()
//	  v := ov.TryErr(myerr)
//	}
func (o Option[T]) TryErr(err error) T {
	if o.IsEmpty() {
		eh.Check(err)
	}
	return o.value
}

// TryErrF returns the option's non-empty value. If the option is empty,
// the user supplied error function will be invoked and TryErrF will panic with
// a try value wrapping this error. This panic can be caught with Catch or Handle in
// errors/handler package. If the user supplied error function returns a nil,
// there will be no panic and TryErrF will return a zero value.
// e.g.
//
//	func tryOption() (err error) {
//	 myerr := errors.New("test error")
//	 fnErr := func() error { return myerr }
//	 defer handler.Catch(&err) // err set to myerr
//	 ov := Empty[int]()
//	 v := ov.TryErrF(fnErr)
//	}
func (o Option[T]) TryErrF(err func() error) T {
	if o.IsEmpty() {
		eh.Check(err())
	}
	return o.value
}

// Returns a pointer to the value in the option. If the value is empty,
// the method will panic with a try value that can be caught with
// handler.Catch() or handler.Handle().
func (o *Option[T]) TryRef() *T {
	if o.IsEmpty() {
		eh.Check(ErrOptionIsEmpty)
	}
	return &o.value
}

// Convert an option to a pointer to an option. Sometimes useful for fluent
// method chaining. E.g
//
//	var slice []int = []int{6, 7}
//	append42 := func(s *[]int) { *s = append(*s, 42) }
//	option.Value(slice).ToRef().Mutate(append42).Get() // []int{6, 7, 42}
func (o Option[T]) ToRef() *Option[T] {
	return &o
}

// Render an option to a string. An empty option
// results in an empty string.
func (o Option[T]) String() string {
	if o.IsEmpty() {
		return ""
	} else {
		return fmt.Sprintf("%v", o.value)
	}
}

// RefOrNil returns a pointer to the value in the option. If the value is empty,
// nil will be returned.
func (o *Option[T]) RefOrNil() *T {
	return o.RefOr(nil)
}

// RefOr returns a pointer to the value in the option. If the value is empty,
// the default pointer will be returned. The primary use case is to allow
// mutation of the value held in the option.
func (o *Option[T]) RefOr(def *T) *T {
	if o.IsEmpty() {
		return def
	} else {
		return &o.value
	}
}

// Set sets the value in the option to a new value. The option will then be
// non-empty.
func (o *Option[T]) Set(v T) {
	o.value = v
	o.nonEmpty = true
}

// SafeSet sets the value in the option to a new value. If the value of v is nil,
// the option will be empty. Otherwise it will be non-empty.
func (o *Option[T]) SafeSet(v T) {
	o.value = v
	o.nonEmpty = !isNil(o.value)
}

// SetRef sets the value in the option to the value pointed to by the parameter.
// If this is nil, the option will be set empty. Otherwise it will be
// non-empty and contain the referenced value.
func (o *Option[T]) SetRef(v *T) {
	if v == nil {
		o.nonEmpty = false
		var zero T
		o.value = zero
	} else {
		o.value = *v
		o.nonEmpty = true
	}
}

// Clear sets the option empty
func (o *Option[T]) Clear() {
	var zero T
	o.value = zero
	o.nonEmpty = false
}

// Map applies a function to the non-empty value of an Option.
// If the option is non-empty, the function is applied
// to it's value, and the result wrapped in an Option
// and returned. Otherwise, an empty option is returned.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if val, ok := o.GetOK(); !ok {
		return Empty[U]()
	} else {
		return Value(f(val))
	}
}

// FlatMap applies a function returning a new option to the
// non-empty value of an Option.
// If the option is non-empty, the function is applied
// to it's value, and the result is returned.
// Otherwise, an empty option is returned.
func FlatMap[T, U any](o Option[T], f func(T) Option[U]) Option[U] {
	if val, ok := o.GetOK(); !ok {
		return Empty[U]()
	} else {
		return f(val)
	}
}

// MapRef is a variation of Map() in which the mapping function takes and
// returns pointers to values. A pointer to the resultant
// option type is returned.
func MapRef[T, U any](o *Option[T], f func(*T) *U) *Option[U] {
	var result *Option[U]
	if r := o.RefOrNil(); r == nil {
		result = Ref[U](nil)
	} else {
		result = Ref(f(r))
	}
	return result
}

// FlatMapRef is a variation of FlatMap() in which the mapping function takes and
// returns pointers to values. A pointer to the resultant
// option type is returned.
func FlatMapRef[T, U any](o *Option[T], f func(*T) *Option[U]) *Option[U] {
	var result *Option[U]
	if r := o.RefOrNil(); r == nil {
		result = Ref[U](nil)
	} else {
		result = f(r)
	}
	return result
}

// Then invokes the supplied function with the Option's value
// if the Option is non-empty. Otherwise, this is a no-op. It
//
//	always returns the option instance it wal called with.
func (o Option[T]) Then(f func(T)) Option[T] {
	if o.nonEmpty {
		f(o.value)
	}
	return o
}

// Then invokes the supplied function with a pointer to the Option's value
// if the Option is non-empty. Otherwise, this is a no-op. It always returns
// the pointer to the Option it is called with
func (o *Option[T]) ThenRef(f func(*T)) *Option[T] {
	if o.nonEmpty {
		f(&o.value)
	}
	return o
}

// Else invokes the provided function if the Option passed is empty.
func (o Option[T]) Else(f func()) {
	if !o.nonEmpty {
		f()
	}
}

// ElseRef invokes the provided function if the Option whose address is passed
// is empty.
func (o *Option[T]) ElseRef(f func()) {
	if !o.nonEmpty {
		f()
	}
}

// Morph, inspired by the concept of [Endomorphism]: https://en.wikipedia.org/wiki/Endomorphism,
// maps an Option value, if non-empty, to another value of the same type, wrapped in an Option.
// Mapping to values of different types via methods is not possible due to limitations in Go
// generics. For this use the option.Map function.
func (o Option[T]) Morph(f func(T) T) Option[T] {
	if o.nonEmpty {
		return Value(f(o.value))
	} else {
		return o
	}
}

// MorphRef, inspired by the concept of [Endomorphism]: https://en.wikipedia.org/wiki/Endomorphism,
// maps a pointer to an Option value, if non-empty, to another value of the same type, wrapped in
// in Option, returned as a pointer to the Option. The mapping function takes and returns pointers
// to the underlying value, where present.
// Mapping to values of different types via methods is not possible due to limitations in Go
// generics. For this use the option.Map function.
func (o *Option[T]) MorphRef(f func(*T) *T) *Option[T] {
	if o.nonEmpty {
		return Ref(f(&o.value))
	} else {
		return o
	}
}

// Mutate applies an in place mutation function to an
// option's value. It is a no-op if the option is empty.
// A pointer to the original option is returned.
func (o *Option[T]) Mutate(f func(*T)) *Option[T] {
	if o.nonEmpty {
		f(&o.value)
	}
	return o
}

// Marshalling / unmarshalling support //

// JSON marshalling of an option. Empty options are
// marshalled as "null". Non-empty options are marshalled
// as their underlying value.
func (o *Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return []byte("null"), nil
	} else {
		return json.Marshal(o.value)
	}
}

// JSON un-marshalling of an option.
func (o *Option[T]) UnmarshalJSON(j []byte) error {
	if len(j) == 0 || string(j) == "null" {
		o.Clear()
	} else {
		if err := json.Unmarshal(j, &o.value); err != nil {
			return err
		}
		o.nonEmpty = true
	}
	return nil
}

// Returns true if the option is empty. Used by the YAML
// marshalling/un-marshalling interface.
func (o Option[T]) IsZero() bool {
	return o.IsEmpty()
}

func (o Option[T]) MarshalYAML() (interface{}, error) {
	return o.value, nil
}

func (o *Option[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&o.value); err != nil {
		return err
	}
	o.nonEmpty = true
	return nil
}

package option

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	eh "github.com/robdavid/genutil-go/errors/handler"
)

var ErrOptionIsEmpty = errors.New("option is empty")
var ErrOptionIsNil = errors.New("option is nil")

// A container for a value that might be empty
type Option[T any] struct {
	value    T
	nonEmpty bool
}

// Creates a non-empty option from a value
func Value[T any](v T) Option[T] {
	return Option[T]{v, true}
}

// Creates an option from a pointer to a value; a nil
// pointer results in an empty option
func Ref[T any](v *T) Option[T] {
	if v == nil {
		var zero T
		return Option[T]{zero, false}
	} else {
		return Option[T]{*v, true}
	}
}

// Create an empty option
func Empty[T any]() Option[T] {
	return Option[T]{}
}

// Returns true if the option is empty and has no value
func (o *Option[T]) IsEmpty() bool {
	return !o.nonEmpty
}

// Returns true if the option is non-empty and has a value
func (o *Option[T]) HasValue() bool {
	return o.nonEmpty
}

func (o *Option[T]) NilErr() error {
	if o.IsEmpty() {
		return ErrOptionIsEmpty
	} else if o.IsNil() {
		return ErrOptionIsNil
	} else {
		return nil
	}

}

// Returns true if the option is empty or the value
// held is nil
func (o *Option[T]) IsNil() bool {
	if o.nonEmpty {
		val := reflect.ValueOf(o.value)
		switch val.Kind() {
		case reflect.Pointer, reflect.Chan, reflect.Slice, reflect.Map, reflect.Interface, reflect.Func:
			return val.IsNil()
		default:
			return false
		}
	}
	return true
}

// Returns true iff the option is not empty and the contained
// value is not nil.
func (o *Option[T]) NonNil() bool {
	return !o.IsNil()
}

// Either return a value if non-empty, or return the value's
// zero value. Note this call does not discriminate between
// an empty option and an option that contains the zero value.
func (o *Option[T]) GetOrZero() T {
	return o.value // o.value is set to the zero value if o.nonEmpty is false
}

// Either return a value if non-empty, or return the default
// provided in the parameter.
func (o *Option[T]) GetOr(def T) T {
	if o.IsEmpty() {
		return def
	} else {
		return o.value
	}
}

// Return the current value and a boolean flag which is true if the option is non-empty.
// If the option is empty. the value returned will be the zero value.
func (o *Option[T]) GetOK() (T, bool) {
	return o.value, o.HasValue()
}

// Get the options' value. If the option is empty, this call will panic
func (o *Option[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return o.value
	}
}

// Get the options' value. If the option is empty, this call will panic
// with a try object that can be caught with Catch or Handle in
// errors/handler package.
func (o *Option[T]) Try() T {
	if o.IsEmpty() {
		eh.Check(ErrOptionIsEmpty)
	}
	return o.value
}

// Get the value of an option. Panics if the option is empty or nil.
// Only useful if the type the option value is nillable, such as a pointer
// interface, or slice. Otherwise, just use Get.
func (o *Option[T]) GetNonNil() T {
	if err := o.NilErr(); err != nil {
		panic(err)
	}
	return o.value
}

func (o *Option[T]) TryNonNil() T {
	eh.Check(o.NilErr())
	return o.value
}

func (o *Option[T]) GetNonNilOK() (T, bool) {
	return o.value, !o.IsNil()
}

// Convert an option to a pointer to an option. Sometimes useful for fluent
// method chaining.
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

// Returns a pointer to the value in the option. If the value is empty,
// the method will panic. The primary use case is to allow mutation of
// the value held in the option.
func (o *Option[T]) Ref() *T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return &o.value
	}
}

// Returns a pointer to the value in the option. If the value is empty,
// the method will panic with a try value that can be caught with
// handler.Catch() or handler.Handle(). The primary use case is to allow
// mutation of the value held in the option.
func (o *Option[T]) TryRef() *T {
	if o.IsEmpty() {
		eh.Check(ErrOptionIsEmpty)
	}
	return &o.value
}

// Returns a pointer to the value in the option. If the value is empty,
// nil will be returned.
func (o *Option[T]) RefOrNil() *T {
	return o.RefOr(nil)
}

// Returns a pointer to the value in the option. If the value is empty,
// a default pointer will be returned.The primary use case is to allow
// mutation of the value held in the option.
func (o *Option[T]) RefOr(def *T) *T {
	if o.IsEmpty() {
		return def
	} else {
		return &o.value
	}
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

func (o *Option[T]) Set(v T) {
	o.value = v
	o.nonEmpty = true
}

func (o *Option[T]) SetRef(v *T) {
	if v == nil {
		o.nonEmpty = false
		var value T
		o.value = value
	} else {
		o.value = *v
		o.nonEmpty = true
	}
}

func (o *Option[T]) Clear() {
	var v T
	o.value = v
	o.nonEmpty = false
}

func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if val, ok := o.GetOK(); !ok {
		return Empty[U]()
	} else {
		return Value(f(val))
	}
}

func MapRef[T, U any](o *Option[T], f func(*T) *U) *Option[U] {
	var result Option[U]
	if r := o.RefOrNil(); r == nil {
		result = Ref[U](nil)
	} else {
		result = Ref(f(r))
	}
	return &result
}

// Marshalling / unmarshalling support //

func (o Option[T]) IsZero() bool {
	return o.IsEmpty()
}

func (o *Option[T]) MarshalJSON() ([]byte, error) {
	if o.IsEmpty() {
		return []byte("null"), nil
	} else {
		return json.Marshal(o.value)
	}
}

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

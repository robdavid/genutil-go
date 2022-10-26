package option

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

var ErrOptionIsEmpty = errors.New("option is empty")

type Option[T any] struct {
	value    T
	nonEmpty bool
}

func Value[T any](v T) Option[T] {
	return Option[T]{v, true}
}

func Ref[T any](v *T) Option[T] {
	if v == nil {
		var zero T
		return Option[T]{zero, false}
	} else {
		return Option[T]{*v, true}
	}
}

func Empty[T any]() Option[T] {
	return Option[T]{}
}

func (o *Option[T]) IsEmpty() bool {
	return !o.nonEmpty
}

func (o *Option[T]) HasValue() bool {
	return o.nonEmpty
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

func (o *Option[T]) GetOrZero() T {
	return o.value
}

func (o *Option[T]) GetOr(def T) T {
	if o.IsEmpty() {
		return def
	} else {
		return o.value
	}
}

func (o *Option[T]) GetOK() (T, bool) {
	return o.value, o.HasValue()
}

func (o *Option[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return o.value
	}
}

func (o Option[T]) ToRef() *Option[T] {
	return &o
}

func (o Option[T]) String() string {
	if o.IsEmpty() {
		return ""
	} else {
		return fmt.Sprintf("%v", o.value)
	}
}

func (o *Option[T]) Ref() *T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return &o.value
	}
}

func (o *Option[T]) RefOrNil() *T {
	return o.RefOr(nil)
}

func (o *Option[T]) RefOr(def *T) *T {
	if o.IsEmpty() {
		return def
	} else {
		return &o.value
	}
}

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

package option

import (
	"encoding/json"
	"errors"
	"reflect"
)

var ErrOptionIsEmpty = errors.New("option is empty")

type Option[T any] struct {
	value    T
	nonEmpty bool
}

func isNil(pv any) bool {
	v := reflect.Indirect(reflect.ValueOf(pv))
	switch v.Kind() {
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Pointer, reflect.Func, reflect.Interface:
		return v.IsNil()
	default:
		return false
	}
}

func Value[T any](v T) Option[T] {
	return Option[T]{v, !isNil(&v)}
}

func Ref[T any](v *T) Option[T] {
	return Option[T]{*v, !isNil(&v)}
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

func (o Option[T]) GetOrZero() T {
	return o.value
}

func (o Option[T]) GetOr(def T) T {
	if o.IsEmpty() {
		return def
	} else {
		return o.value
	}
}

func (o Option[T]) GetOK() (T, bool) {
	return o.value, o.HasValue()
}

func (o Option[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return o.value
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
	o.nonEmpty = !isNil(&v)
}

func (o *Option[T]) SetRef(v *T) {
	o.value = *v
	o.nonEmpty = !isNil(v)
}

func (o *Option[T]) Clear() {
	var v T
	o.value = v
	o.nonEmpty = false
}

func Equal[T comparable](o, p *Option[T]) bool {
	if o.value == p.value {
		return true
	} else if o.HasValue() && p.HasValue() && o.value == p.value {
		return true
	}
	return false
}

func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if o.IsEmpty() {
		return Empty[U]()
	} else {
		u := f(o.value)
		return Value(u)
	}
}

func MapRef[T, U any](o *Option[T], f func(*T) U) Option[U] {
	if o.IsEmpty() {
		return Empty[U]()
	} else {
		u := f(&o.value)
		return Value(u)
	}
}

// Marshalling / unmarshalling support //

func (o *Option[T]) IsZero() bool {
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
		o.nonEmpty = !isNil(&o.value)
	}
	return nil
}

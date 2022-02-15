package option

import (
	"errors"
	"reflect"
)

var ErrOptionIsEmpty = errors.New("option is empty")

type Option[T any] struct {
	value     T
	nonEmpty  bool
	isNilable bool
}

func isNilable(v any) bool {
	switch reflect.ValueOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
		return true
	default:
		return false
	}
}

func Value[T any](v T) Option[T] {
	return Option[T]{v, true, isNilable(v)}
}

func Ref[T any](v *T) Option[T] {
	return Option[T]{*v, true, isNilable(*v)}
}

func Empty[T any]() Option[T] {
	return Option[T]{}
}

func (o *Option[T]) IsEmpty() bool {
	return !o.nonEmpty || (o.isNilable && reflect.ValueOf(o.value).IsNil())
}

func (o *Option[T]) HasValue() bool {
	return !o.IsEmpty()
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
	o.nonEmpty = true
	o.isNilable = isNilable(v)
}

func (o *Option[T]) SetRef(v *T) {
	o.value = *v
	o.nonEmpty = true
	o.isNilable = isNilable(*v)
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

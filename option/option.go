package option

import (
	"encoding/json"
	"errors"
)

var ErrOptionIsEmpty = errors.New("option is empty")

type IOption[T any] interface {
	IsEmpty() bool
	HasValue() bool
	GetOrZero() T
	GetOr(T) T
	GetOK() (T, bool)
	Get() T
}

type IOptionRef[T any] interface {
	IOption[T]
	Ref() *T
	RefOrNil() *T
	RefOr(*T) *T
	RefOK() (*T, bool)
	Set(T)
	SetRef(*T)
	Clear()
}

type Option[T any] struct {
	value    T
	nonEmpty bool
}

type OptionRef[T any] struct {
	ref *T
}

func Value[T any](v T) Option[T] {
	return Option[T]{v, true}
}

func Ref[T any](v *T) OptionRef[T] {
	return OptionRef[T]{v}
}

func Empty[T any]() Option[T] {
	return Option[T]{}
}

func (o Option[T]) IsEmpty() bool {
	return !o.nonEmpty
}

func (o OptionRef[T]) IsEmpty() bool {
	return o.ref == nil
}

func (o Option[T]) HasValue() bool {
	return o.nonEmpty
}

func (o OptionRef[T]) HasValue() bool {
	return o.ref != nil
}

func (o Option[T]) GetOrZero() T {
	return o.value
}

func (o OptionRef[T]) GetOrZero() T {
	if o.ref != nil {
		return *o.ref
	} else {
		var zero T
		return zero
	}
}

func (o Option[T]) GetOr(def T) T {
	if o.IsEmpty() {
		return def
	} else {
		return o.value
	}
}

func (o OptionRef[T]) GetOr(def T) T {
	if o.ref != nil {
		return *o.ref
	} else {
		return def
	}
}

func (o Option[T]) GetOK() (T, bool) {
	return o.value, o.HasValue()
}

func (o OptionRef[T]) GetOK() (value T, ok bool) {
	if o.ref != nil {
		ok = true
		value = *o.ref
	}
	return
}

func (o Option[T]) Get() T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return o.value
	}
}

func (o OptionRef[T]) Get() T {
	if o.ref == nil {
		panic(ErrOptionIsEmpty)
	} else {
		return *o.ref
	}
}

func (o *Option[T]) Ref() *T {
	if o.IsEmpty() {
		panic(ErrOptionIsEmpty)
	} else {
		return &o.value
	}
}

func (o *OptionRef[T]) Ref() *T {
	if o.ref == nil {
		panic(ErrOptionIsEmpty)
	} else {
		return o.ref
	}
}

func (o *Option[T]) RefOrNil() *T {
	return o.RefOr(nil)
}

func (o *OptionRef[T]) RefOrNil() *T {
	return o.ref
}

func (o *Option[T]) RefOr(def *T) *T {
	if o.IsEmpty() {
		return def
	} else {
		return &o.value
	}
}

func (o *OptionRef[T]) RefOr(def *T) *T {
	if o.ref == nil {
		return def
	} else {
		return o.ref
	}
}

func (o *Option[T]) RefOK() (*T, bool) {
	if o.IsEmpty() {
		return nil, false
	} else {
		return &o.value, true
	}
}

func (o *OptionRef[T]) RefOK() (ref *T, ok bool) {
	if o != nil {
		ref = o.ref
		ok = true
	}
	return
}

func (o *Option[T]) Set(v T) {
	o.value = v
	o.nonEmpty = true
}

func (o *OptionRef[T]) Set(v T) {
	o.ref = &v
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

func (o *OptionRef[T]) SetRef(v *T) {
	o.ref = v
}

func (o *Option[T]) Clear() {
	var v T
	o.value = v
	o.nonEmpty = false
}

func (o *OptionRef[T]) Clear() {
	o.ref = nil
}

func Equal[T comparable](o, p IOption[T]) bool {
	if o.IsEmpty() && p.IsEmpty() {
		return true
	} else if o.IsEmpty() != p.IsEmpty() {
		return false
	} else {
		return o.Get() == p.Get()
	}
}

func Map[T, U any](o IOption[T], f func(T) U) Option[U] {
	if o.IsEmpty() {
		return Empty[U]()
	} else {
		u := f(o.Get())
		return Value(u)
	}
}

func MapRef[T, U any](o IOptionRef[T], f func(*T) *U) OptionRef[U] {
	if r := o.RefOrNil(); r == nil {
		return Ref[U](nil)
	} else {
		return Ref(f(r))
	}
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

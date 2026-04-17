package opt

import "errors"

var ErrOptionIsEmpty = errors.New("option is empty")

type Val[T any] struct {
	value    T
	nonEmpty bool
}

type Ref[T any] struct {
	reference *T
}

type Option[T any] interface {
	IsEmpty() bool
	HasValue() bool
	Get() T
	Ref() *T
	// GetOr returns the option value if present, or otherwise returns the fallback value.
	GetOr(fallback T) T
	// RefOr returns a reference to the option value if present, or otherwise returns the fallback
	// reference
	RefOr(*T) *T
}

func (v Val[T]) Get() T {
	if !v.nonEmpty {
		panic(ErrOptionIsEmpty)
	} else {
		return v.value
	}
}

func (r Ref[T]) Get() T {
	if r.reference == nil {
		panic(ErrOptionIsEmpty)
	} else {
		return *r.reference
	}
}

func (v Val[T]) GetOr(fallback T) T {
	if !v.nonEmpty {
		return fallback
	} else {
		return v.value
	}
}

func (r Ref[T]) GetOr(fallback T) T {
	if r.reference == nil {
		return fallback
	} else {
		return *r.reference
	}
}

func (v *Val[T]) Ref() *T {
	if !v.nonEmpty {
		panic(ErrOptionIsEmpty)
	} else {
		return &v.value
	}
}

func (r Ref[T]) Ref() *T {
	if r.reference == nil {
		panic(ErrOptionIsEmpty)
	} else {
		return r.reference
	}
}

func (v *Val[T]) RefOr(fallback *T) *T {
	if !v.nonEmpty {
		return fallback
	} else {
		return &v.value
	}
}

func (r Ref[T]) RefOr(fallback *T) *T {
	if r.reference == nil {
		return fallback
	} else {
		return r.reference
	}
}

// IsEmpty returns true if there is no value
func (v Val[T]) IsEmpty() bool {
	return !v.nonEmpty
}

// Has value returns true if there is a value, and it is
// safe to access the value via functions such as [Val.Get].
func (v Val[T]) HasValue() bool {
	return v.nonEmpty
}

func (r Ref[T]) IsEmpty() bool {
	return r.reference == nil
}

func (r Ref[T]) HasValue() bool {
	return r.reference != nil
}

// FromVal creates an Option[T] from a value of type T.
// It returns a Val[T] instance.
func FromVal[T any](obj T) Val[T] {
	return Val[T]{value: obj, nonEmpty: true}
}

// FromRef creates an Option[T] from a pointer of type *T.
// It returns a Ref[T] instance.
func FromRef[T any](obj *T) Ref[T] {
	return Ref[T]{reference: obj}
}

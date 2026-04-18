package opt

import (
	"errors"
	"fmt"

	"github.com/robdavid/genutil-go/errors/handler"
)

var ErrOptionIsEmpty = errors.New("option is empty")

type Val[T any] struct {
	value    T
	nonEmpty bool
}

type Ref[T any] struct {
	reference *T
}

// Value creates an Option[T] from a value of type T.
// It returns a Val[T] instance.
func Value[T any](obj T) Val[T] {
	return Val[T]{value: obj, nonEmpty: true}
}

// Reference creates an Option[T] from a pointer of type *T.
// It returns a Ref[T] instance.
func Reference[T any](obj *T) Ref[T] {
	return Ref[T]{reference: obj}
}

func Empty[T any]() Val[T] {
	return Val[T]{}
}

func EmptyRef[T any]() Ref[T] {
	return Ref[T]{}
}

// Option represents an object that may or may not contain an underlying value of type T.
// This may be held by value (implemented by *[Val][T]) or by reference (implemented by [Ref][T])
// The various method allows access to the underlying by value or reference, or detect if the
// underlying value exists at all.
type Option[T any] interface {

	// IsEmpty returns true if there is no underlying value. Attempts to access the value will fail,
	// typically leading to panic conditions.
	IsEmpty() bool

	// HasValue returns true if the underlying value is present. When true, it is safe to access the
	// underlying value directly.
	HasValue() bool

	// Get returns the option value if present. Otherwise the method will panic.
	Get() T

	// GetOK either returns the underlying option value and true if the value is present, or
	// the zero value for T and false if not.
	GetOK() (T, bool)

	// Try returns the option value if present, or else will panic, similar to the [Option.Get] method.
	// However, the panic raised is one that can be recovered via [handler.Catch] or [handler.Handle]
	// functions.
	Try() T

	// Ref returns a reference to the underlying option value if there is one. If not, the function
	// will panic.
	Ref() *T

	// RefOK either returns a reference to the underlying value and true if
	// the value is present, or a nil pointer and false if not.
	RefOK() (*T, bool)

	// TryRef returns a reference to the underlying option value if there is one. If not,
	// it will panic, similar to [Option.Ref]. However, the panic raised is one that can be recovered via
	// [handler.Catch] or [handler.Handle] functions.
	TryRef() *T

	// GetOr returns the option value if present, or otherwise returns the fallback value.
	GetOr(fallback T) T

	// RefOr returns a reference to the option value if present, or otherwise returns the fallback
	// reference
	RefOr(*T) *T

	// GetOrF returns the option value if present, or otherwise invokes the provided function
	// and returns the value obtained
	GetOrF(fallbackFn func() T) T
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

// ToRef makes a copy of the value option and returns a reference to it. This is useful
// for fluent method chaining.
func (v Val[T]) ToRef() *Val[T] {
	return &v
}

// AsRef converts the [Val][T] instance to a [Ref][T] that references the value held if present.
// Otherwise it returns an empty [Ref][T].
func (v *Val[T]) AsRef() Ref[T] {
	if v.nonEmpty {
		return Reference(&v.value)
	} else {
		return EmptyRef[T]()
	}
}

// GetOK returns the underlying value if present, along with a boolean flag set true
// indicating that the value is present. Otherwise, it returns the zero value
// for T and false.
func (v Val[T]) GetOK() (val T, ok bool) {
	if !v.nonEmpty {
		var zero T
		return zero, false
	}
	return v.value, true
}

// GetOK returns the referenced value if present, along with a boolean value of true
// indicating that the value is present. Otherwise, it returns the zero value
// for T and false.
func (r Ref[T]) GetOK() (val T, ok bool) {
	if r.reference == nil {
		var zero T
		return zero, false
	}
	return *r.reference, true
}

// GetOr returns the value of the option if it exists, otherwise returns the
func (v Val[T]) GetOr(fallback T) T {
	if !v.nonEmpty {
		return fallback
	} else {
		return v.value
	}
}

// GetOr returns the value of the option if it exists, otherwise returns the
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

// RefOr returns a pointer to the value stored in the Val option if it exists,
func (v *Val[T]) RefOr(fallback *T) *T {
	if !v.nonEmpty {
		return fallback
	} else {
		return &v.value
	}
}

// RefOr returns a pointer to the value referenced by the Ref option if it
func (r Ref[T]) RefOr(fallback *T) *T {
	if r.reference == nil {
		return fallback
	} else {
		return r.reference
	}
}

// GetOrF returns the value of the option if it exists, otherwise invokes the
func (v Val[T]) GetOrF(fallbackFn func() T) T {
	if !v.nonEmpty {
		return fallbackFn()
	} else {
		return v.value
	}
}

// GetOrF returns the value of the option if it exists, otherwise invokes the
func (r Ref[T]) GetOrF(fallbackFn func() T) T {
	if r.reference == nil {
		return fallbackFn()
	} else {
		return *r.reference
	}
}

// Try returns the value of the option if it exists, or otherwise raises an
// error that can be recovered via [handler.Catch] or [handler.Handle]
// functions.
func (v Val[T]) Try() T {
	if !v.nonEmpty {
		handler.Raise(ErrOptionIsEmpty)
	}
	return v.value
}

// Try returns the value of the option if it exists, or otherwise raises an
// error that can be recovered via [handler.Catch] or [handler.Handle]
// functions.
func (r Ref[T]) Try() T {
	if r.reference == nil {
		handler.Raise(ErrOptionIsEmpty)
	}
	return *r.reference
}

// TryRef returns a reference to the underlying option value if there is one. If not,
// it will panic, similar to [Option.Ref]. However, the panic raised is one that can be recovered via
// [handler.Catch] or [handler.Handle] functions.
func (v Val[T]) TryRef() *T {
	if !v.nonEmpty {
		handler.Raise(ErrOptionIsEmpty)
	}
	return &v.value
}

// TryRef returns a reference to the underlying option value if there is one. If not,
// it will panic, similar to [Option.Ref]. However, the panic raised is one that can be recovered via
// [handler.Catch] or [handler.Handle] functions.
func (r Ref[T]) TryRef() *T {
	if r.reference == nil {
		handler.Raise(ErrOptionIsEmpty)
	}
	return r.reference
}

// RefOK returns a reference to the underlying value and true if the value is present,
// or a nil pointer and false if not.
func (v Val[T]) RefOK() (*T, bool) {
	if !v.nonEmpty {
		return nil, false
	}
	return &v.value, true
}

// RefOK returns a reference to the underlying value and true if the value is present,
// or a nil pointer and false if not.
func (r Ref[T]) RefOK() (*T, bool) {
	if r.reference == nil {
		return nil, false
	}
	return r.reference, true
}

func (v Val[T]) String() string {
	if v.nonEmpty {
		return fmt.Sprint(v.value)
	} else {
		return ""
	}
}

func (r Ref[T]) String() string {
	if r.reference != nil {
		return fmt.Sprint(*r.reference)
	} else {
		return ""
	}
}

// IsEmpty returns true if there is no value
func (v Val[T]) IsEmpty() bool {
	return !v.nonEmpty
}

func (r Ref[T]) IsEmpty() bool {
	return r.reference == nil
}

// Has value returns true if there is a value, and it is
// safe to access the value via functions such as [Val.Get].
func (v Val[T]) HasValue() bool {
	return v.nonEmpty
}

func (r Ref[T]) HasValue() bool {
	return r.reference != nil
}

// Map applies a function to the non-empty value of an Option.
// If the option is non-empty, the function is applied
// to it's value, and the result wrapped in an Option
// and returned. Otherwise, an empty option is returned.
func Map[T, U any](o Option[T], f func(T) U) Option[U] {
	if val, ok := o.GetOK(); !ok {
		return Empty[U]().ToRef()
	} else {
		return Value(f(val)).ToRef()
	}
}

// MapRef is a variation of Map() in which the mapping function takes and
// returns pointers to values. A pointer to the resultant
// option type is returned.
func MapRef[T, U any](o Option[T], f func(*T) *U) Option[U] {
	if r := o.RefOr(nil); r == nil {
		return EmptyRef[U]()
	} else {
		return Reference(f(r))
	}
}

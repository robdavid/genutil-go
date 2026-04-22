package opt

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/robdavid/genutil-go/errors/handler"
)

// ErrOptionIsEmpty is an error raised (via a panic) when an option is empty.
var ErrOptionIsEmpty = errors.New("option is empty")

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

	// Mutate applies an in place mutation function to an
	// option's value. It is a no-op if the option is empty.
	// The mutated option is returned.
	Mutate(f func(*T)) Option[T]

	// Ensure makes sure the option is non-empty. If it is already
	// non-empty, it is a no-op. Otherwise it is mutated to be populated
	// with the zero value. The mutated or original option is returned.
	Ensure() Option[T]

	// AsRef converts the underlying option implementation to a [Ref][T] and
	// returns it. It's a no-op returning the receiver if the implementation is
	// already a [Ref][T].
	AsRef() Ref[T]

	// Morph, inspired by the concept of [Endomorphism]:
	// https://en.wikipedia.org/wiki/Endomorphism, maps an [Option] value, if
	// non-empty, to another value of the same type, wrapped in an [Val]. If the
	// passed [Option] is empty, an empty [Val][T] is returned. Mapping to
	// values of different types via methods is not possible due to limitations
	// in Go generics. For this use the [Map] function.
	Morph(func(T) T) Val[T]

	// MorphRef, inspired by the concept of [Endomorphism]:
	// https://en.wikipedia.org/wiki/Endomorphism, maps a Option value, if
	// non-empty, to another value of the same type, wrapped in in a [Ref][T].
	// The mapping function takes and returns pointers to the underlying values.
	// If the passed Option is empty, an empty [Ref][T] is returned.
	// Mapping to values of different types via methods is not possible due to
	// limitations in Go generics. For this use the option.Map function.
	MorphRef(func(*T) *T) Ref[T]

	// Then invokes the supplied function with the Option's value
	// if the Option is non-empty. Otherwise, this is a no-op. It
	// always returns the option instance it was called with.
	Then(func(T)) Option[T]

	// Else invokes the provided function if the Option is empty.
	Else(func()) Option[T]
}

// Val is an [Option] implementation which consists of a member of type T, and a
// boolean flag which is true if the option holds a value. It is suitable for
// primitive values, such as int or string, or small structures for which
// copying represents a small overhead.
type Val[T any] struct {
	value    T
	nonEmpty bool
}

// Ref is an [Option] implementation which is simply a pointer to the underlying
// value, which is nil if there is no value present. It is suitable for larger
// structures the copying of which would represent a significant overhead, or
// for situations where a reference is required (e.g. to support mutability).
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

// Empty creates a [Val][T] with no value.
func Empty[T any]() Val[T] {
	return Val[T]{}
}

// EmptyRef creates a [Ref][T] with no value.
func EmptyRef[T any]() Ref[T] {
	return Ref[T]{}
}

// Get returns the underlying value if there is one, or else the function panics with the
// value of [ErrOptionIsEmpty].
func (v Val[T]) Get() T {
	if !v.nonEmpty {
		panic(ErrOptionIsEmpty)
	} else {
		return v.value
	}
}

// Ref returns the underlying referenced value, if there is one , or else the function panics
// with the value of [ErrOptionIsEmpty].
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

// AsRef converts the [Ref][T] instance to a [Ref][T] simply by returning the receiver.
func (r Ref[T]) AsRef() Ref[T] {
	return r
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

// IsEmpty returns true if there is no value
func (r Ref[T]) IsEmpty() bool {
	return r.reference == nil
}

// Has value returns true if there is a value, and it is
// safe to access the value via functions such as [Val.Get].
func (v Val[T]) HasValue() bool {
	return v.nonEmpty
}

// Has value returns true if there is a value, and it is
// safe to access the value via functions such as [Ref.Get].
func (r Ref[T]) HasValue() bool {
	return r.reference != nil
}

func (v *Val[T]) Mutate(f func(*T)) Option[T] {
	if v.nonEmpty {
		f(&v.value)
	}
	return v
}

func (r Ref[T]) Mutate(f func(*T)) Option[T] {
	if r.reference != nil {
		f(r.reference)
	}
	return r
}

// Ensure ensures that the option is non-empty. If it is already non-empty, it is a no-op.
// Otherwise, it is mutated to be populated with the zero value. The mutated or original option is returned.
//
// For [Val][T], this method uses a pointer receiver to allow for fluent method chaining.
func (v *Val[T]) Ensure() Option[T] {
	if !v.nonEmpty {
		var zero T
		v.value = zero
		v.nonEmpty = true
	}
	return v
}

// Ensure ensures that the option is non-empty. If it is already non-empty, it is a no-op.
// Otherwise, it is mutated to be populated with the zero value. The mutated or original option is returned.
func (r Ref[T]) Ensure() Option[T] {
	if r.reference == nil {
		var zero T
		r.reference = &zero
	}
	return r
}

// Morph transforms the underlying value, if present, by means of the supplied
// function f. If the underlying value is present, it is passed to the function
// and the resulting value is wrapped in a [Val][T]. If there is no underlying
// value, an empty Val[T] is returned.
func (v Val[T]) Morph(f func(T) T) Val[T] {
	if v.nonEmpty {
		return Value(f(v.value))
	} else {
		return Empty[T]()
	}
}

// Morph transforms the underlying value, if present, by means of the supplied
// function f. If the underlying value is present, it is passed to the function
// and the resulting value is wrapped in a [Val][T]. If there is no underlying
// value, an empty Val[T] is returned.
func (r Ref[T]) Morph(f func(T) T) Val[T] {
	if r.reference != nil {
		return Value(f(*r.reference))
	} else {
		return Empty[T]()
	}
}

// MorphRef transforms the underlying value, if present, by means of the supplied
// function f. If the underlying value is present, a reference to it is passed to
// the function, and the resulting value reference is wrapped in a [Ref][T]. If
// there is no underlying value, and empty [Ref][T] is returned.
func (v *Val[T]) MorphRef(f func(*T) *T) Ref[T] {
	if v.nonEmpty {
		return Reference(f(&v.value))
	} else {
		return EmptyRef[T]()
	}
}

// MorphRef transforms the underlying value, if present, by means of the supplied
// function f. If the underlying value is present, a reference to it is passed to
// the function, and the resulting value reference is wrapped in a [Ref][T]. If
// there is no underlying value, and empty [Ref][T] is returned.
func (r Ref[T]) MorphRef(f func(*T) *T) Ref[T] {
	if r.reference != nil {
		return Reference(f(r.reference))
	} else {
		return EmptyRef[T]()
	}
}

// Then invokes the supplied function with the Option's value if the Option is
// non-empty. Otherwise, this is a no-op. It always returns the option instance
// it was called with.
func (v Val[T]) Then(f func(T)) Option[T] {
	if v.nonEmpty {
		f(v.value)
	}
	return &v
}

// Then invokes the supplied function with the Option's value if the Option is
// non-empty. Otherwise, this is a no-op. It always returns the option instance
// it was called with.
func (r Ref[T]) Then(f func(T)) Option[T] {
	if r.reference != nil {
		f(*r.reference)
	}
	return r
}

// Else invokes the provided function if the Option passed is empty.
func (v Val[T]) Else(f func()) Option[T] {
	if !v.nonEmpty {
		f()
	}
	return &v
}

// Else invokes the provided function if the Option passed is empty.
func (r Ref[T]) Else(f func()) Option[T] {
	if r.reference == nil {
		f()
	}
	return r
}

// ThenRef invokes the supplied function with a reference to the Option's value
// if the Option is non-empty. Otherwise, this is a no-op. It always returns the
// option instance it was called with.
func (v Val[T]) ThenRef(f func(*T)) Option[T] {
	if v.nonEmpty {
		f(&v.value)
	}
	return &v
}

// ThenRef invokes the supplied function with a reference to the Option's value
// if the Option is non-empty. Otherwise, this is a no-op. It always returns the
// option instance it was called with.
func (r Ref[T]) ThenRef(f func(*T)) Option[T] {
	if r.reference != nil {
		f(r.reference)
	}
	return r
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

// Marshalling / unmarshaling support //

// MarshalJSON implements JSON marshaling of a [Val][T] object. Empty options
// are marshaled as "null". Non-empty options are marshaled as their
// underlying value.
func (v *Val[T]) MarshalJSON() ([]byte, error) {
	if !v.nonEmpty {
		return []byte("null"), nil
	} else {
		return json.Marshal(v.value)
	}
}

// UnmarshalJSON implements JSON umarshaling into a [Val][T] object. An
// input of null or zero length unmarshals as an empty value. Otherwise,
// the input in unmarshaled as the underlying type.
func (v *Val[T]) UnmarshalJSON(j []byte) error {
	if len(j) == 0 || string(j) == "null" {
		*v = Empty[T]()
	} else {
		if err := json.Unmarshal(j, &v.value); err != nil {
			return err
		}
		v.nonEmpty = true
	}
	return nil
}

// Returns true if the option is empty. Used by the YAML
// marshaling/un-marshaling interface, and by the standard
// library JSON v2 marshaling if using "omitzero".
func (v Val[T]) IsZero() bool {
	return !v.nonEmpty
}

// MarshalYAML implements YAML marshaling of a [Val][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. An empty value
// is marshaled as it's zero value. Otherwise it is simply marshaled
// and the underlying value. Note that if "omitempty" is used, this
// function won't be called for empty values, as it should be guarded
// due the the [Val.IsZero] method.
func (v Val[T]) MarshalYAML() (any, error) {
	if !v.nonEmpty {
		var zero T
		return zero, nil
	}
	return v.value, nil
}

// UnmarshalYAML implements YAML unmarshaling into a [Val][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. Input is unmarshaled
// into the underlying value, and the [Val] will always be non-empty, unless
// an error is
func (v *Val[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&v.value); err != nil {
		return err
	}
	v.nonEmpty = true
	return nil
}

// MarshalJSON implements JSON marshaling of a [Ref][T] object. Empty options
// are marshaled as "null". Non-empty options are marshaled as their
// underlying value.
func (r Ref[T]) MarshalJSON() ([]byte, error) {
	if r.reference == nil {
		return []byte("null"), nil
	} else {
		return json.Marshal(r.reference)
	}
}

// UnmarshalJSON implements JSON umarshaling into a [Ref][T] object. An
// input of null or zero length unmarshals as an empty value. Otherwise,
// the input in unmarshaled as the underlying type.
func (r *Ref[T]) UnmarshalJSON(j []byte) error {
	if len(j) == 0 || string(j) == "null" {
		r.reference = nil
	} else {
		r.reference = new(T)
		if err := json.Unmarshal(j, r.reference); err != nil {
			return err
		}
	}
	return nil
}

// Returns true if the option is empty. Used by the YAML
// marshaling/un-marshaling interface, and by the standard
// library JSON v2 marshaling if using "omitzero".
func (r Ref[T]) IsZero() bool {
	return r.reference == nil
}

// MarshalYAML implements YAML marshaling of a [Ref][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. An empty value
// is marshaled as it's zero value. Otherwise it is simply marshaled
// and the underlying value. Note that if "omitempty" is used, this
// function won't be called for empty values, as it should be guarded
// due the the [Ref.IsZero] method.
func (r Ref[T]) MarshalYAML() (any, error) {
	if r.reference == nil {
		var zero T
		return zero, nil
	}
	return r.reference, nil
}

// UnmarshalYAML implements YAML unmarshaling into a [Ref][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. Input is unmarshaled
// into the underlying value, and the [Ref] will always be non-empty, unless
// an error is
func (r *Ref[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	r.reference = new(T)
	if err := unmarshal(r.reference); err != nil {
		return err
	}
	return nil
}

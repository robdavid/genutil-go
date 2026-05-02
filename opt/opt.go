package opt

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/robdavid/genutil-go/errors/handler"
)

// ErrOptionIsEmpty is an error raised (via a panic) when an option is empty.
var ErrOptionIsEmpty = errors.New("optional value not present")

// Opt represents a wrapper around an optional value of type T.
// It can hold a value by value ([Val][T]) or by reference ([Ref][T]). The various
// methods provide safe and direct ways to access the underlying value or detect its presence.
type Opt[T any] interface {

	// IsEmpty returns true if no value is present. Accessing the underlying
	// value when empty will typically result in a panic.
	IsEmpty() bool

	// HasValue returns true if the underlying value is present. When true, it
	// is safe to access the underlying value directly, e.g. via [Get].
	HasValue() bool

	// IsRef returns true if the [Option] implementation is reference based.
	IsRef() bool

	// Get returns the option value if present. If empty, this function panics
	// with an error containing [ErrOptionIsEmpty].
	Get() T

	// GetOK either returns the underlying option value and true if the value is present, or
	// the zero value for T and false if not. This approach avoids panics.
	GetOK() (T, bool)

	// Try returns the option value if present, or else will panic, similar to
	// [Option.Get] method. However, the panic raised is one that can be
	// recovered via [handler.Catch] or [handler.Handle] functions.
	Try() T

	// Ref returns a reference to the option value if there is one.
	// If not, the function will panic with an error containing
	// [ErrOptionIsEmpty].
	Ref() *T

	// RefOK either returns a reference to the option value and true if
	// the value is present, or a nil pointer and false if not.
	RefOK() (*T, bool)

	// TryRef returns a reference to the option value if there is
	// one. If not, it will panic, similar to [Option.Ref]. However, the panic
	// raised is one that can be recovered via [handler.Catch] or
	// [handler.Handle] functions.
	TryRef() *T

	// GetOr returns the value if present; otherwise, it returns the provided
	// fallback value.
	GetOr(fallback T) T

	// RefOr returns a reference to the option value if present; otherwise, it
	// returns the fallback reference.
	RefOr(*T) *T

	// GetOrF returns the value if present, otherwise it executes and returns
	// the result of the provided function.
	GetOrF(fallbackFn func() T) T

	// Morph, inspired by the concept of [Endomorphism]:
	// https://en.wikipedia.org/wiki/Endomorphism, maps an [Option] value. If
	// non-empty, it applies f(T) and wraps the result in a [Val][T]. If empty,
	// an empty [Val][T] is returned. Mapping to any type other than [Val][T]
	// requires the use of the [Map]() function.
	Morph(func(T) T) Opt[T]

	// MorphRef, inspired by the concept of [Endomorphism]:
	// https://en.wikipedia.org/wiki/Endomorphism, maps a Option value. If
	// non-empty, it applies f(*T) and wraps the resulting pointer in Ref[T]. If
	// empty, an empty Ref[T] is returned. Mapping to any type other than [Ref][T]
	// requires the use of the [Map]() function.
	MorphRef(func(*T) *T) Opt[T]

	// Then executes the supplied function if the Option is non-empty. It always
	// returns the option instance it was called with.
	Then(func(T)) Opt[T]

	// Else executes the provided function if the Option is empty. It always
	// returns the option instance it was called with.
	Else(func()) Opt[T]
}

// MutOpt is an extension of [Opt] which provides methods for mutation of option
// values. [MutOpt][T] is implemented by *[Val][T] and *[Ref][T]. Methods in
// this interface (that is, excluding those in [Opt]) are able to mutate the
// underlying value held by the option wrapper ([Val] or [Ref]). Methods in
// [Opt], by contrast, may or may not be able to obtain a reference to the
// underlying value depending on which option type is being used. In the case of
// [Val] instances, only a reference to a copy of the data can be obtained.
type MutOpt[T any] interface {
	Opt[T]

	// Mutate applies an in place mutation function to an option's value. It is
	// a no-op if the option is empty. The mutated option is returned.
	Mutate(f func(*T)) MutOpt[T]

	// Ensure ensures that the option is non-empty. If it is already non-empty,
	// it is a no-op. Otherwise, it is populated with the zero value. The
	// mutated or original option is returned.
	Ensure() MutOpt[T]

	// Set sets the underlying value, mutating the option. It will cause the
	// option to become non-empty, if it isn't already. The mutated option is
	// returned.
	Set(value T) MutOpt[T]

	// SetRef sets the underlying value by reference, mutating the option. For a
	// [Val], the underlying value is set by copying the value pointed to by
	// reference, unless reference is nil. For a [Ref] the underlying reference
	// becomes the one provided. If reference is nil, the option becomes empty,
	// otherwise it becomes non-empty. The mutated option is returned.
	SetRef(reference *T) MutOpt[T]

	// Causes the option to become empty, if it isn't already, mutating it. The
	// mutated option is returned.
	Unset() MutOpt[T]

	// SetFrom mutates the option to copy the value or reference from an [Opt][T] in opt. If opt is
	// empty, the resulting option will be empty.
	SetFrom(opt Opt[T]) MutOpt[T]
}

// Val is an [Opt] implementation which consists of a member of type T, and a
// boolean flag indicating presence. It is suitable for primitive values (int, string)
// or small structures where copying overhead is negligible.
type Val[T any] struct {
	value    T
	nonEmpty bool
}

// Ref is an [Opt] implementation which is simply contains a pointer to the
// underlying value, which is nil if there is no value present. It is suitable
// for larger structures or situations requiring mutability.
type Ref[T any] struct {
	reference *T
}

// Value creates a [Val][T] instance from a value of type T.
func Value[T any](obj T) Val[T] {
	return Val[T]{value: obj, nonEmpty: true}
}

// Reference creates a [Ref][T] instance from a pointer of type *T.
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

// ValFrom creates a [Val][T] from the value obtained from an [Opt][T]. If the
// [Opt] is empty, the result will be empty.
func ValFrom[T any](opt Opt[T]) Val[T] {
	if v, ok := opt.GetOK(); ok {
		return Value(v)
	} else {
		return Empty[T]()
	}
}

// RefFrom creates a [Ref][T] from the reference obtained from an Opt[T]. If the
// [Opt] is empty, the result will be empty.
func RefFrom[T any](opt Opt[T]) Ref[T] {
	return Reference(opt.RefOr(nil))
}

func (v *Val[T]) error() error {
	typ := reflect.TypeFor[T]()
	return fmt.Errorf("error in opt.Val[%s]: %w", typ.Name(), ErrOptionIsEmpty)
}

func (r Ref[T]) error() error {
	typ := reflect.TypeFor[T]()
	return fmt.Errorf("error in opt.Ref[%s]: %w", typ.Name(), ErrOptionIsEmpty)
}

// IsEmpty returns true if there is no value
func (v Val[T]) IsEmpty() bool {
	return !v.nonEmpty
}

// IsEmpty returns true if there is no value
func (r Ref[T]) IsEmpty() bool {
	return r.reference == nil
}

// HasValue returns true if there is a value, and it is safe to access the value
// via functions such as [Val.Get].
func (v Val[T]) HasValue() bool {
	return v.nonEmpty
}

// HasValue returns true if there is a value, and it is safe to access the value
// via functions such as [Ref.Get].
func (r Ref[T]) HasValue() bool {
	return r.reference != nil
}

// IsRef is always false for [Val] types
func (Val[T]) IsRef() bool {
	return false
}

// IsRef is always true for [Ref] types
func (Ref[T]) IsRef() bool {
	return true
}

// Get returns the underlying value if there is one, or else this function
// panics with an error containing [ErrOptionIsEmpty].
func (v Val[T]) Get() T {
	if !v.nonEmpty {
		panic(v.error())
	} else {
		return v.value
	}
}

// Get returns the underlying referenced value if there is one, or else this
// function panics with an error containing [ErrOptionIsEmpty].
func (r Ref[T]) Get() T {
	if r.reference == nil {
		panic(r.error())
	} else {
		return *r.reference
	}
}

// AsRef converts the [Val][T] instance to a [Ref][T] that references the value
// held if present. Otherwise it returns an empty [Ref][T].
func (v *Val[T]) AsRef() Ref[T] {
	if v.nonEmpty {
		return Reference(&v.value)
	} else {
		return EmptyRef[T]()
	}
}

func (r Ref[T]) AsVal() Val[T] {
	if r.reference != nil {
		return Value(*r.reference)
	} else {
		return Empty[T]()
	}
}

// GetOK returns the underlying value and a true boolean if present. It returns
// the zero value for T and false if not present.
func (v Val[T]) GetOK() (val T, ok bool) {
	if !v.nonEmpty {
		var zero T
		return zero, false
	}
	return v.value, true
}

// GetOK returns the underlying value and a true boolean if present presence. It returns
// the zero value for T and false if not present.
func (r Ref[T]) GetOK() (val T, ok bool) {
	if r.reference == nil {
		var zero T
		return zero, false
	}
	return *r.reference, true
}

// GetOr returns the value of the option if it exists; otherwise, it returns the
// provided fallback value.
func (v Val[T]) GetOr(fallback T) T {
	if !v.nonEmpty {
		return fallback
	} else {
		return v.value
	}
}

// GetOr returns the value of the option if it exists; otherwise, it returns the
// provided fallback value.
func (r Ref[T]) GetOr(fallback T) T {
	if r.reference == nil {
		return fallback
	} else {
		return *r.reference
	}
}

// Ref returns a pointer to a copy of the underlying value if present. If the
// option is empty, it panics with an error containing [ErrOptionIsEmpty].
func (v Val[T]) Ref() *T {
	if !v.nonEmpty {
		panic(v.error())
	} else {
		return &v.value
	}
}

// Ref returns a pointer to the underlying value if present. If empty, it panics
// with an error containing [ErrOptionIsEmpty].
func (r Ref[T]) Ref() *T {
	if r.reference == nil {
		panic(r.error())
	} else {
		return r.reference
	}
}

// RefOr returns a pointer to a copy of the underlying value if present. If empty it
// returns the fallback pointer.
func (v Val[T]) RefOr(fallback *T) *T {
	if !v.nonEmpty {
		return fallback
	} else {
		return &v.value
	}
}

// RefOr returns a pointer to the value referenced by the Ref option if it
// exists, otherwise returns the fallback pointer.
func (r Ref[T]) RefOr(fallback *T) *T {
	if r.reference == nil {
		return fallback
	} else {
		return r.reference
	}
}

// GetOrF returns the value of the option if it exists, otherwise it executes and
// returns the result of the provided function.
func (v Val[T]) GetOrF(fallbackFn func() T) T {
	if !v.nonEmpty {
		return fallbackFn()
	} else {
		return v.value
	}
}

// GetOrF returns the value of the option if it exists, otherwise executes and
// returns the result of the provided function.
func (r Ref[T]) GetOrF(fallbackFn func() T) T {
	if r.reference == nil {
		return fallbackFn()
	} else {
		return *r.reference
	}
}

// Try returns the value of the option if it exists, or otherwise raises an
// error that can be recovered via [handler.Catch] or [handler.Handle].
func (v Val[T]) Try() T {
	if !v.nonEmpty {
		handler.Raise(v.error())
	}
	return v.value
}

// Try returns the value of the option if it exists, or otherwise raises an
// error that can be recovered via [handler.Catch] or [handler.Handle].
func (r Ref[T]) Try() T {
	if r.reference == nil {
		handler.Raise(r.error())
	}
	return *r.reference
}

// Try returns the value of the option if it exists, or otherwise raises the
// error value provided in the err parameter which can be recovered via
// [handler.Catch] or [handler.Handle].
func (v Val[T]) TryOr(err error) T {
	if !v.nonEmpty {
		handler.Raise(err)
	}
	return v.value
}

// Try returns the value of the option if it exists, or otherwise raises the
// error value provided in the err parameter which can be recovered via
// [handler.Catch] or [handler.Handle].
func (r Ref[T]) TryOr(err error) T {
	if r.reference == nil {
		handler.Raise(err)
	}
	return *r.reference
}

// TryRef returns a reference to the underlying option value if there is one. If
// not, it will panic, similar to [Opt.Ref]. However, the panic raised is one
// that can be recovered via [handler.Catch] or [handler.Handle] functions.
func (v Val[T]) TryRef() *T {
	if !v.nonEmpty {
		handler.Raise(v.error())
	}
	return &v.value
}

// TryRef returns a reference to the underlying option value if there is one. If
// not, it will panic, similar to [Opt.Ref]. However, the panic raised is one
// that can be recovered via [handler.Catch] or [handler.Handle] functions.
func (r Ref[T]) TryRef() *T {
	if r.reference == nil {
		handler.Raise(r.error())
	}
	return r.reference
}

// TryRef returns a reference to the underlying option value if there is one. If
// not, it will panic, similar to [Opt.Ref]. However, the panic raised is one
// that wraps the error provided in err and can be recovered via [handler.Catch]
// or [handler.Handle] functions.
func (v Val[T]) TryRefOr(err error) *T {
	if !v.nonEmpty {
		handler.Raise(err)
	}
	return &v.value
}

// TryRef returns a reference to the underlying option value if there is one. If
// not, it will panic, similar to [Opt.Ref]. However, the panic raised is one
// that wraps the error provided in err and can be recovered via [handler.Catch]
// or [handler.Handle] functions.
func (r Ref[T]) TryRefOr(err error) *T {
	if r.reference == nil {
		handler.Raise(err)
	}
	return r.reference
}

// RefOK returns a reference to the underlying value and true if the value is
// present, or a nil pointer and false if not.
func (v Val[T]) RefOK() (*T, bool) {
	if !v.nonEmpty {
		return nil, false
	}
	return &v.value, true
}

// RefOK returns a reference to the underlying value and true if the value is
// present, or a nil pointer and false if not.
func (r Ref[T]) RefOK() (*T, bool) {
	if r.reference == nil {
		return nil, false
	}
	return r.reference, true
}

// String returns a string representation of the underlying value if present,
// or an empty string if the option is empty.
func (v Val[T]) String() string {
	if v.nonEmpty {
		return fmt.Sprint(v.value)
	} else {
		return ""
	}
}

// String returns a string representation of the underlying value if present,
// or an empty string if the option is empty.
func (r Ref[T]) String() string {
	if r.reference != nil {
		return fmt.Sprint(*r.reference)
	} else {
		return ""
	}
}

// Mutate applies function f to a reference to a copy of the underlying value if
// present, returning a modified [Val][T] object. If there is no value present,
// the method is a no-op and the receiver is returned.
func (v *Val[T]) Mutate(f func(*T)) MutOpt[T] {
	if v.nonEmpty {
		f(&v.value)
	}
	return v
}

// Mutate applies function f to a reference to the underlying value if present.
// The function may alter the value via this pointer. If there is no value
// present, the method is a no-op. The receiver is always returned.
func (r *Ref[T]) Mutate(f func(*T)) MutOpt[T] {
	if r.reference != nil {
		f(r.reference)
	}
	return r
}

// Ensure ensures that the option is non-empty. If it is already non-empty, it
// returns the receiver. Otherwise, a new empty Val[T] is returned.
func (v *Val[T]) Ensure() MutOpt[T] {
	if !v.nonEmpty {
		var zero T
		*v = Value(zero)
	}
	return v
}

// Ensure ensures that the option is non-empty. If it is already non-empty, it
// is a no-op. Otherwise, it is mutated to be populated with the zero value. The
// mutated or original option is returned.
func (r *Ref[T]) Ensure() MutOpt[T] {
	if r.reference == nil {
		var zero T
		r.reference = &zero
	}
	return r
}

// Set sets the underlying value, mutating the option. It will cause the
// option to become non-empty, if it isn't already. The mutated option is
// returned.
func (v *Val[T]) Set(value T) MutOpt[T] {
	*v = Value(value)
	return v
}

// Set sets the underlying value, mutating the option. It will cause the
// option to become non-empty, if it isn't already. The mutated option is
// returned.
func (r *Ref[T]) Set(value T) MutOpt[T] {
	if r.reference == nil {
		r.reference = &value
	} else {
		*r.reference = value
	}
	return r
}

// Unset causes the option to become empty, if it isn't already, mutating it. The
// mutated option is returned.
func (v *Val[T]) Unset() MutOpt[T] {
	*v = Empty[T]()
	return v
}

// Unset causes the option to become empty, if it isn't already, mutating it. The
// mutated option is returned.
func (r *Ref[T]) Unset() MutOpt[T] {
	r.reference = nil
	return r
}

// SetRef sets the underlying value by reference, mutating the option. If
// reference is nil, the option becomes empty. Otherwise the underlying value is
// set by copying the value pointed to by reference, and the value becomes
// non-empty. The mutated option is returned.
func (v *Val[T]) SetRef(reference *T) MutOpt[T] {
	if reference == nil {
		*v = Empty[T]()
	} else {
		*v = Value(*reference)
	}
	return v
}

// SetRef sets the underlying value by reference, mutating the option. The
// underlying reference becomes the one provided. If the option was empty, it
// becomes non empty, unless reference is nil. The mutated option is returned.
func (r *Ref[T]) SetRef(reference *T) MutOpt[T] {
	r.reference = reference
	return r
}

// SetFrom sets the underlying value from the value obtained from opt, mutating
// the option. If opt is empty, the option will be empty.
func (v *Val[T]) SetFrom(opt Opt[T]) MutOpt[T] {
	if optv, ok := opt.GetOK(); ok {
		*v = Value(optv)
	} else {
		*v = Empty[T]()
	}
	return v
}

// SetFrom sets the underlying reference from the reference obtained from opt,
// mutating the option. If opt is empty, the option will be empty.
func (r *Ref[T]) SetFrom(opt Opt[T]) MutOpt[T] {
	r.reference = opt.RefOr(nil)
	return r
}

// Morph transforms the underlying value, if present, by means of the supplied
// function f. If non-empty, the function is applied to the value, and the
// result is wrapped in a [Val][T]. If empty, an empty Val[T] is returned.
func (v Val[T]) Morph(f func(T) T) Opt[T] {
	if v.nonEmpty {
		return Value(f(v.value))
	} else {
		return Empty[T]()
	}
}

// Morph transforms the underlying value, if present, by means of the supplied
// function f. If non-empty, it applies f to the pointer and wraps the result in
// a [Ref][T]. If empty, an empty Ref[T] is returned.
func (r Ref[T]) Morph(f func(T) T) Opt[T] {
	if r.reference != nil {
		value := f(*r.reference)
		return Reference(&value)
	} else {
		return EmptyRef[T]()
	}
}

// MorphRef transforms the underlying value, if present, by means of the
// supplied function f. If non-empty, a reference to a copy of it is passed to
// the function, and the resulting pointer is wrapped in a [Val][T]. If empty,
// an empty [Val][T] is returned. The receiver is not modified.
func (v Val[T]) MorphRef(f func(*T) *T) Opt[T] {
	if v.nonEmpty {
		return Value(*f(&v.value))
	} else {
		return Empty[T]()
	}
}

// MorphRef transforms the underlying value, if present, by means of the
// supplied function f. If non-empty, a reference to it is passed to the
// function, and the resulting pointer is wrapped in a [Ref][T]. If empty, an
// empty [Ref][T] is returned.
func (r Ref[T]) MorphRef(f func(*T) *T) Opt[T] {
	if r.reference != nil {
		return Reference(f(r.reference))
	} else {
		return EmptyRef[T]()
	}
}

// Then executes the supplied function with the value held by v, if v is non-empty. Otherwise, this is a
// no-op. It always returns a pointer to v.
func (v Val[T]) Then(f func(T)) Opt[T] {
	if v.nonEmpty {
		f(v.value)
	}
	return &v
}

// Then executes the supplied function with the value referenced by r if r is non-empty. Otherwise, this is a
// no-op. It always returns r.
func (r Ref[T]) Then(f func(T)) Opt[T] {
	if r.reference != nil {
		f(*r.reference)
	}
	return r
}

// Else executes the provided function if v is empty. It always
// returns a pointer to v.
func (v Val[T]) Else(f func()) Opt[T] {
	if !v.nonEmpty {
		f()
	}
	return &v
}

// Else executes the provided function if r is empty. It always
// returns r.
func (r Ref[T]) Else(f func()) Opt[T] {
	if r.reference == nil {
		f()
	}
	return r
}

// ThenRef invokes the supplied function with a reference to v's value if v is
// non-empty. Otherwise, this is a no-op. It always returns a pointer to v.
func (v Val[T]) ThenRef(f func(*T)) Opt[T] {
	if v.nonEmpty {
		f(&v.value)
	}
	return &v
}

// ThenRef invokes the supplied function with r's pointer to its value if r is non-empty.
// Otherwise, this is a no-op. It always returns r
func (r Ref[T]) ThenRef(f func(*T)) Opt[T] {
	if r.reference != nil {
		f(r.reference)
	}
	return r
}

// Map applies a function to the non-empty value of an [Opt]. If the option
// is non-empty, the function is applied to its value, and the result is wrapped
// as an [Opt][U] and returned. Otherwise, an empty option is returned.
func Map[T, U any](o Opt[T], f func(T) U) Opt[U] {
	if val, ok := o.GetOK(); !ok {
		return Empty[U]()
	} else {
		return Value(f(val))
	}
}

// MapRef is a variation of [Map]() in which the mapping function takes and
// returns pointers to values. The referenced computed value is returned as an
// [Opt][U].
func MapRef[T, U any](o Opt[T], f func(*T) *U) Opt[U] {
	if r := o.RefOr(nil); r == nil {
		return EmptyRef[U]()
	} else {
		return Reference(f(r))
	}
}

// Equal compares two Option values for equality. It checks if both options are
// empty or both are non-empty. If one is empty and the other is not, it returns
// false. If both are empty, it returns true. If both are non-empty, it
// dereferences the underlying values and compares them using the == operator.
// The type T must be comparable for this function to work.
func Equal[T comparable](o1 Opt[T], o2 Opt[T]) bool {
	if o1.IsEmpty() != o2.IsEmpty() {
		return false
	} else if o1.IsEmpty() && o2.IsEmpty() {
		return true
	} else {
		return *o1.Ref() == *o2.Ref()
	}
}

// DeepEqual compares two Option values for deep equality. It checks if both
// options are empty or both are non-empty. If one is empty and the other is
// not, it returns false. If both are empty, it returns true. If both are
// non-empty, it uses reflect.DeepEqual to compare the underlying values
// (including nested structures, slices, maps, and pointers).
func DeepEqual[T any](o1 Opt[T], o2 Opt[T]) bool {
	if o1.IsEmpty() != o2.IsEmpty() {
		return false
	} else if o1.IsEmpty() && o2.IsEmpty() {
		return true
	} else {
		return reflect.DeepEqual(o1.Ref(), o2.Ref())
	}
}

// Marshalling / unmarshaling support //

// IsZero returns true if the option is empty. Used by the YAML
// marshaling/un-marshaling interface, and by the standard library JSON v2
// marshaling if using "omitzero".
func (v Val[T]) IsZero() bool {
	return !v.nonEmpty
}

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

// UnmarshalJSON implements JSON unmarshaling into a [Val][T] object. An
// input of null or zero length unmarshals as an empty value. Otherwise,
// the input is unmarshaled into the underlying type.
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

// MarshalYAML implements YAML marshaling of a [Val][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. An empty value
// is marshaled as its zero value. Otherwise it is simply marshaled
// as the underlying value. Note that if "omitempty" is used, this
// function won't be called for empty values, as it should be guarded
// by the [Val.IsZero] method.
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
// an error occurs during parsing.
func (v *Val[T]) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&v.value); err != nil {
		return err
	}
	v.nonEmpty = true
	return nil
}

// IsZero returns true if the option is empty. Used by the YAML
// marshaling/un-marshaling interface, and by the standard
// library JSON v2 marshaling if using "omitzero".
func (r Ref[T]) IsZero() bool {
	return r.reference == nil
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

// UnmarshalJSON implements JSON unmarshaling into a [Ref][T] object. An
// input of null or zero length unmarshals as an empty value. Otherwise,
// the input is unmarshaled into the underlying type.
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

// MarshalYAML implements YAML marshaling of a [Ref][T] for the
// https://pkg.go.dev/gopkg.in/yaml.v2 YAML parser. An empty value
// is marshaled as its zero value. Otherwise it is simply marshaled
// as the underlying value. Note that if "omitempty" is used, this
// function won't be called for empty values, as it should be guarded
// by the [Ref.IsZero] method.
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
// an error occurs during parsing.
func (r *Ref[T]) UnmarshalYAML(unmarshal func(any) error) error {
	r.reference = new(T)
	if err := unmarshal(r.reference); err != nil {
		return err
	}
	return nil
}

package types

import "github.com/robdavid/genutil-go/option"

// Attempts a type assertion of a to type T. If
// successful, an option of value T is returned.
// Otherwise, an empty option is returned.
func As[T any, U any](a U) option.Option[T] {
	v, ok := any(a).(T)
	if ok {
		return option.Value(v)
	} else {
		return option.Empty[T]()
	}
}

// Attempts a type assertion of a to type T. If
// successful, and a is non-nil, an option of value T is returned.
// Otherwise, an empty option is returned.
func AsRef[T any, U any](a U) *option.Option[T] {
	v, _ := any(a).(*T)
	return option.Ref(v)
}

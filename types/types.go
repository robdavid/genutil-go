package types

import (
	"github.com/robdavid/genutil-go/opt"
)

// Attempts a type assertion of a to type T. If
// successful, an option of value T is returned.
// Otherwise, an empty option is returned.
func As[T any, U any](a U) opt.Val[T] {
	v, ok := any(a).(T)
	if ok {
		return opt.Value(v)
	} else {
		return opt.Empty[T]()
	}
}

// Attempts a type assertion of a to type T. If
// successful, and a is non-nil, an option of value T is returned.
// Otherwise, an empty option is returned.
func AsRef[T any, U any](a U) opt.Ref[T] {
	v, _ := any(a).(*T)
	return opt.Reference(v)
}

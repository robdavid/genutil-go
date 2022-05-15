package result

// Contains a value of type T, or an error
type Result[T any] struct {
	value T
	err   error
}

func (r *Result[T]) Pair() (T, error) {
	return r.value, r.err
}

func (r *Result[T]) GetErr() error {
	return r.err
}

func (r *Result[T]) Get() T {
	return r.value
}

func (r *Result[T]) IsError() bool {
	return r.err != nil
}

func (r *Result[T]) IsValue() bool {
	return r.err == nil
}

// Creates a new non-error result
func Value[T any](t T) Result[T] {
	return Result[T]{t, nil}
}

// Creates an error result
func Error[T any](err error) Result[T] {
	var t T
	return Result[T]{t, err}
}

func Pair[T any](t T, err error) Result[T] {
	return Result[T]{t, err}
}

func Map[T, U any](r Result[T], f func(T) U) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		return Value(f(r.Get()))
	}
}

func MapErr[T, U any](r Result[T], f func(T) (U, error)) Result[U] {
	if r.IsError() {
		return Error[U](r.GetErr())
	} else {
		if u, err := f(r.Get()); err != nil {
			return Error[U](err)
		} else {
			return Value(u)
		}
	}
}

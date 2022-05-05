package iterator

import "github.com/robdavid/genutil-go/option"

type Iterator[T any] interface {
	Next() bool
	Value() (T, error)
	MustValue() T
}

type Func[T, U any] struct {
	base    Iterator[T]
	mapping func(T) U
	value   option.OptionRef[U]
}

func (i *Func[T, U]) Next() bool {
	i.value.Clear()
	return i.base.Next()
}

func (i *Func[T, U]) Value() (u U, err error) {
	if i.value.IsEmpty() {
		var t T
		if t, err = i.base.Value(); err != nil {
			return
		}
		u = i.mapping(t)
		i.value.Set(u)
	} else {
		u = i.value.Get()
	}
	return
}

func (i *Func[T, U]) MustValue() U {
	if u, err := i.Value(); err != nil {
		panic(err)
	} else {
		return u
	}
}

type SliceIter[T any] struct {
	slice []T
	index int
	value T
}

func (si *SliceIter[T]) Next() bool {
	if si.index < len(si.slice) {
		si.value = si.slice[si.index]
		si.index++
		return true
	} else {
		return false
	}
}

func (si *SliceIter[T]) MustValue() T {
	return si.value
}

func (si *SliceIter[T]) Value() (T, error) {
	return si.value, nil
}

func Slice[T any](slice []T) Iterator[T] {
	var t T
	return &SliceIter[T]{slice, 0, t}
}

func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	return &Func[T, U]{iter, mapping, option.Ref[U](nil)}
}

func MustCollect[T any](iter Iterator[T]) []T {
	if result, err := Collect(iter); err != nil {
		panic(err)
	} else {
		return result
	}
}

func Collect[T any](iter Iterator[T]) ([]T, error) {
	result := make([]T, 0)
	for iter.Next() {
		if val, err := iter.Value(); err != nil {
			return result, err
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

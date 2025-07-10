package iterator

import (
	"iter"

	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/option"
)

// funcNext is a function supporting a transforming operation by consuming
// all or part of an iterator, returning the next value
type funcNext[T any, U any] func(T) (U, bool)

// mapIter wraps an iterator and adds a mapping function
type mapIter[T, U any] struct {
	base     Iterator[T]
	mapping  funcNext[T, U]
	value    U
	sizeFunc func(IteratorSize) IteratorSize
}

func (i *mapIter[T, U]) Next() bool {
	for {
		if ok := i.base.Next(); !ok {
			return false
		}
		if value, ok := i.mapping(i.base.Value()); ok {
			i.value = value
			return true
		}
	}
}

func (i *mapIter[T, U]) Value() U {
	return i.value
}

func (i *mapIter[T, U]) Abort() {
	i.base.Abort()
}

func (i *mapIter[T, U]) Reset() {
	i.base.Reset()
}

func (i *mapIter[T, U]) Size() IteratorSize {
	return i.sizeFunc(i.base.Size())
}

func (i *mapIter[T, U]) SeqOK() bool { return i.base.SeqOK() }

func (i *mapIter[T, U]) Seq() iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range i.base.Seq() {
			if mv, ok := i.mapping(v); ok {
				if !yield(mv) {
					break
				}
			}
		}
	}
}

// wrapFunc creates a new iterator from an existing iterator and a function that consumes it, yielding
// one element at a time.
func wrapFunc[T any, U any](iterator Iterator[T], f funcNext[T, U], sizeFunc func(sz IteratorSize) IteratorSize) Iterator[U] {
	return NewDefaultIterator(&mapIter[T, U]{base: iterator, mapping: f, sizeFunc: sizeFunc})
}

// Map applies function `mapping` of type `func(T) U` to each value, producing
// a new iterator over `U`.
func Map[T any, U any](iter Iterator[T], mapping func(T) U) Iterator[U] {
	mapNext := func(value T) (U, bool) {
		return mapping(value), true
	}
	return wrapFunc(iter, mapNext, functions.Id)
}

// Filter applies a filter function `predicate` of type `func(T) bool`, producing
// a new iterator containing only the elements than satisfy the function.
func Filter[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	filterNext := func(value T) (T, bool) {
		return value, predicate(value)
	}
	return wrapFunc(iter, filterNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterMap applies both transformation and filtering logic to an iterator. The function `mapping` is
// applied to each element of type `T`, producing either an option value of type `U` or an empty
// option. The result is an iterator over `U` drawn from only the non-empty options
// returned.
func FilterMap[T any, U any](iter Iterator[T], mapping func(T) option.Option[U]) Iterator[U] {
	filterMapNext := func(value T) (U, bool) {
		return mapping(value).GetOK()
	}
	return wrapFunc(iter, filterMapNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterValues takes an iterator of results and returns an iterator of the underlying
// result value type for only those results that have no error.
func FilterValues[T any](iter Iterator[result.Result[T]]) Iterator[T] {
	return FilterMap(iter, func(res result.Result[T]) option.Option[T] {
		if res.IsError() {
			return option.Empty[T]()
		} else {
			return option.Value(res.Get())
		}
	})
}

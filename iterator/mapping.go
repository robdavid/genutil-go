package iterator

import (
	"iter"

	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/option"
)

// funcNext is a function supporting a transforming operation by consuming
// the next value, returning optionally a new value; if the boolean return is false
// no value has been returned and function should be called again with the next value.
type funcNext[T any, U any] func(T) (U, bool)

// funcNext is a function supporting a transforming operation by consuming the next key
// and value, returning optionally a new key and value; if the boolean return is false no
// key or value has been returned and function should be called again with the next kay
// and value from the iterator.
type func2Next[K, V, X, Y any] func(K, V) (X, Y, bool)

type mapIterBase[T, U any] struct {
	base     CoreIterator[T]
	value    U
	sizeFunc func(IteratorSize) IteratorSize
}

func newMapIterBase[T, U any](iterator CoreIterator[T], sizeFunc func(sz IteratorSize) IteratorSize) mapIterBase[T, U] {
	return mapIterBase[T, U]{base: iterator, sizeFunc: sizeFunc}
}

// mapIter wraps an Iterator and adds a mapping function
type mapIter[T, U any] struct {
	mapIterBase[T, U]
	mapping funcNext[T, U]
}

// wrapFunc creates a new Iterator from an existing iterator and a function that consumes it, yielding
// (at most) one element at a time.
func wrapFunc[T any, U any](iterator CoreIterator[T], f funcNext[T, U], sizeFunc func(sz IteratorSize) IteratorSize) Iterator[U] {
	return NewDefaultIterator(&mapIter[T, U]{mapIterBase: newMapIterBase[T, U](iterator, sizeFunc), mapping: f})
}

// mapIter2 wraps an Iterator2 and adds a mapping function of type func2Next.
type mapIter2[K, V, X, Y any] struct {
	mapIterBase[V, Y]
	key      X
	base2    CoreIterator2[K, V]
	mapping2 func2Next[K, V, X, Y]
}

// wrapFunc2 creates a new Iterator2 from an existing Iterator2 and a function that
// consumes it, yielding (at most) one key and value at a time.
func wrapFunc2[K, V, X, Y any](iterator CoreIterator2[K, V], f func2Next[K, V, X, Y], sizeFunc func(sz IteratorSize) IteratorSize) Iterator2[X, Y] {
	return NewDefaultIterator2(&mapIter2[K, V, X, Y]{mapIterBase: newMapIterBase[V, Y](iterator, sizeFunc), mapping2: f, base2: iterator})
}

func (i *mapIterBase[T, U]) Value() U {
	return i.value
}

func (i *mapIterBase[T, U]) Abort() {
	i.base.Abort()
}

func (i *mapIterBase[T, U]) Reset() {
	i.base.Reset()
}

func (i *mapIterBase[T, U]) Size() IteratorSize {
	return i.sizeFunc(i.base.Size())
}

func (i *mapIterBase[T, U]) SeqOK() bool { return i.base.SeqOK() }

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

func (i *mapIter2[K, V, X, Y]) Key() X {
	return i.key
}

func (i *mapIter2[K, V, X, Y]) Next() bool {
	for {
		if ok := i.base2.Next(); !ok {
			return false
		}
		if key, value, ok := i.mapping2(i.base2.Key(), i.base2.Value()); ok {
			i.value = value
			i.key = key
			return true
		}
	}
}

func (i *mapIter2[K, V, X, Y]) Seq() iter.Seq[Y] {
	return func(yield func(Y) bool) {
		for k, v := range i.base2.Seq2() {
			if _, mv, ok := i.mapping2(k, v); ok {
				if !yield(mv) {
					break
				}
			}
		}
	}
}

func (i *mapIter2[K, V, X, Y]) Seq2() iter.Seq2[X, Y] {
	return func(yield func(X, Y) bool) {
		for k, v := range i.base2.Seq2() {
			if mk, mv, ok := i.mapping2(k, v); ok {
				if !yield(mk, mv) {
					break
				}
			}
		}
	}
}

// Map applies function mapping of type func(T) U to each value, producing
// a new iterator over U.
func Map[T any, U any](iter CoreIterator[T], mapping func(T) U) Iterator[U] {
	mapNext := func(value T) (U, bool) {
		return mapping(value), true
	}
	return wrapFunc(iter, mapNext, functions.Id)
}

// Map2 applies function mapping of type func(K, V) (X,Y) to each key and value pair, producing
// a new Iterator2 over X and Y.
func Map2[K, V, X, Y any](iter CoreIterator2[K, V], mapping func(K, V) (X, Y)) Iterator2[X, Y] {
	mapNext2 := func(key K, value V) (X, Y, bool) {
		k, v := mapping(key, value)
		return k, v, true
	}
	return wrapFunc2(iter, mapNext2, functions.Id)
}

// Filter applies a filter function predicate of type func(T) bool, producing
// a new iterator containing only the elements that satisfy the function.
func Filter[T any](iter CoreIterator[T], predicate func(T) bool) Iterator[T] {
	filterNext := func(value T) (T, bool) {
		return value, predicate(value)
	}
	return wrapFunc(iter, filterNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// Filter2 applies a filter function of type func(K, V) bool over each key and value pair,
// producing a new iterator containing only the elements that satisfy the function.
func Filter2[K, V any](iter CoreIterator2[K, V], predicate func(K, V) bool) Iterator2[K, V] {
	filterNext := func(key K, value V) (K, V, bool) {
		return key, value, predicate(key, value)
	}
	return wrapFunc2(iter, filterNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterMap applies both transformation and filtering logic to an Iterator. The function
// mapping is applied to each value of type T, producing either a new value of type U
// and a true boolean, or undefined value and a false boolean. Values are taken from
// the value when the boolean is true to produce a new Iterator; returns when
// the boolean is false are ignored.
func FilterMap[T, U any](iter CoreIterator[T], mapping func(T) (U, bool)) Iterator[U] {
	return wrapFunc(iter, mapping, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterMapOpt applies both transformation and filtering logic to an iterator. The
// function mapping is applied to each element of type T, producing either an option
// value of type U or an empty option. The result is an iterator over U drawn from
// only the non-empty options returned.
func FilterMapOpt[T any, U any](iter CoreIterator[T], mapping func(T) option.Option[U]) Iterator[U] {
	filterMapNext := func(value T) (U, bool) {
		return mapping(value).GetOK()
	}
	return wrapFunc(iter, filterMapNext, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterMap2 applies both transformation and filtering logic to an Iterator2. The
// function mapping is applied to each key and value pair, of type K and V respectively,
// producing either a new key and value pairs (of type X and Y) and a true boolean, or
// undefined key and value and a false boolean. Values are taken from the keys and values
// when the boolean is true to produce a new Iterator2.
func FilterMap2[K, V, X, Y any](iter CoreIterator2[K, V], mapping func(K, V) (X, Y, bool)) Iterator2[X, Y] {
	return wrapFunc2(iter, mapping, func(sz IteratorSize) IteratorSize { return sz.Subset() })
}

// FilterValues takes an iterator of results and returns an iterator of the underlying
// result value type for only those results that have no error.
func FilterValues[T any](iter CoreIterator[result.Result[T]]) Iterator[T] {
	return FilterMapOpt(iter, func(res result.Result[T]) option.Option[T] {
		if res.IsError() {
			return option.Empty[T]()
		} else {
			return option.Value(res.Get())
		}
	})
}

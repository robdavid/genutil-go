package iterator

import (
	"iter"

	"github.com/robdavid/genutil-go/tuple"
)

type Indexed[T any] = tuple.Tuple2[int, T]

func IndexValue[T any](index int, value T) Indexed[T] {
	return tuple.Of2(index, value)
}

type enumeratedCoreIterator[T any] struct {
	CoreIterator[T]
	key, index int
}

func (eci *enumeratedCoreIterator[T]) Next() bool {
	if !eci.CoreIterator.Next() {
		return false
	}
	eci.key = eci.index
	eci.index++
	return true
}

func (eci *enumeratedCoreIterator[T]) Key() int {
	return eci.key
}

func (eci *enumeratedCoreIterator[T]) Seq2() iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for v := range eci.Seq() {
			eci.key = eci.index
			if !yield(eci.key, v) {
				eci.Abort()
				break
			}
			eci.index++
		}
	}
}

func newEnumeratedCoreIterator[T any](citr CoreIterator[T]) *enumeratedCoreIterator[T] {
	return &enumeratedCoreIterator[T]{CoreIterator: citr, key: 0, index: 0}
}

// Enumerate takes a CoreIterator and builds an Iterator2 that returns the pair of
// the index of each element (starting at zero) and the original element.
func Enumerate[T any](itr CoreIterator[T]) Iterator2[int, T] {
	return NewDefaultIterator2(newEnumeratedCoreIterator(itr))
}

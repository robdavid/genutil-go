package iterator

import "iter"

// Iterator over a slice
type sliceIter[T any] struct {
	slice *[]T
	index int
	ref   *T
}

func (si *sliceIter[T]) Next() bool {
	if si.index < len(*si.slice) {
		si.ref = &(*si.slice)[si.index]
		si.index++
		return true
	} else {
		return false
	}
}

func (si *sliceIter[T]) Value() T {
	if si.ref == nil {
		var zero T
		return zero
	}
	return *si.ref
}

func (si *sliceIter[T]) Set(e T) {
	if si.ref != nil {
		*si.ref = e
	}
}

func (si *sliceIter[T]) Delete() {
	if si.index > len(*si.slice) || si.index < 1 || si.ref == nil {
		return
	}
	*si.slice = append((*si.slice)[:si.index-1], (*si.slice)[si.index:]...)
	si.index--
	si.ref = nil
}

func (si *sliceIter[T]) Abort() {
	si.index = len(*si.slice)
}

func (si *sliceIter[T]) Reset() {
	si.index = 0
}

func (si *sliceIter[T]) Size() IteratorSize {
	return NewSize(len(*si.slice) - si.index)
}

func (si *sliceIter[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		defer si.Abort()
		for si.index = 0; si.index < len(*si.slice); {
			si.ref = &(*si.slice)[si.index]
			si.index++
			if !yield(*si.ref) {
				break
			}
		}
	}
}

func (si *sliceIter[T]) SeqOK() bool { return si.index == 0 }

func NewSliceCoreIterator[T any](slice *[]T) CoreMutableIterator[T] {
	return &sliceIter[T]{slice: slice, index: 0}
}

type sliceIterRef[T any] struct {
	sliceIter[T]
}

func (sir *sliceIterRef[T]) Value() *T {
	return sir.ref
}

func (sir *sliceIterRef[T]) Set(e *T) {
	*sir.ref = *e
}

func (sir *sliceIterRef[T]) Seq() iter.Seq[*T] {
	return func(yield func(*T) bool) {
		defer sir.Abort()
		for sir.index = 0; sir.index < len(*sir.slice); {
			sir.ref = &(*sir.slice)[sir.index]
			sir.index++
			if !yield(sir.ref) {
				break
			}
		}
	}
}

// Slice makes an Iterator[T] from slice []T, containing all the elements
// from the slice in order.
//
// Deprecated: use slices.Iter()
func Slice[T any](slice []T) Iterator[T] {
	iter := &sliceIter[T]{slice: &slice, index: 0}
	return NewDefaultIterator(iter)
}

// MutSlice makes a MutableIterator[T] from slice []T, containing all the elements
// from the slice in order.
//
// Deprecated: use slices.IterMut()
func MutSlice[T any](slice *[]T) MutableIterator[T] {
	iter := &sliceIter[T]{slice: slice, index: 0}
	return NewDefaultMutableIterator(iter)
}

// Of makes an Iterator[T] containing the variadic arguments of type T
func Of[T any](elements ...T) Iterator[T] {
	return Slice(elements)
}

package iterator

import "iter"

type emptyIter[T any] struct{}

type emptyIter2[K any, V any] struct {
	emptyIter[V]
}

func (emptyIter[T]) Next() bool                { return false }
func (emptyIter[T]) Value() T                  { var zero T; return zero }
func (emptyIter[T]) Size() IteratorSize        { return NewSize(0) }
func (emptyIter[T]) Abort()                    {}
func (emptyIter[T]) Reset()                    {}
func (emptyIter[T]) SeqOK() bool               { return false }
func (emptyIter[T]) Seq() iter.Seq[T]          { return func(yield func(T) bool) {} }
func (emptyIter2[K, V]) Key() K                { var zero K; return zero }
func (emptyIter2[K, V]) Seq2() iter.Seq2[K, V] { return func(yield func(K, V) bool) {} }

// Empty creates an iterator that returns no items.
func Empty[T any]() Iterator[T] {
	return NewDefaultIterator(emptyIter[T]{})
}

// Empty2 creates an Iterator2 that returns no items.
func Empty2[K any, V any]() Iterator2[K, V] {
	return NewDefaultIterator2(emptyIter2[K, V]{})
}

func EmptySeq[T any]() iter.Seq[T] {
	return func(func(T) bool) {}
}

func EmptySeq2[K any, V any]() iter.Seq2[K, V] {
	return func(func(K, V) bool) {}
}

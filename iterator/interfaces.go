package iterator

import "iter"

// SimpleIterator supports a simple sequence of elements
type SimpleIterator[T any] interface {
	// Next sets the iterator's current value to be the first, and subsequent, iterator elements.
	// False is returned only when there are no more elements (the current value remains unchanged)
	Next() bool
	// Value gets the current iterator value.
	Value() T
	// Abort stops the iterator; subsequent calls to Next() will return false.
	Abort()
	// Reset stops the iterator; subsequent calls to Next() will begin the iterator from the start.
	// Note not all iterators are guaranteed to return the same sequence again, for example iterators
	// that perform IO may not read the same data again.
	Reset()
}

type SimpleMutableIterator[T any] interface {
	SimpleIterator[T]
	// Set allows a value to be modified in place
	Set(T)
	// Delete deletes the current value, which must be the last value returned by Next(). This
	// function may not be implemented for all iterator types, in which case it will return an
	// ErrDeleteNotImplemented error.
	Delete() error
}

// CoreIterator is an extension of SimpleIterator that in aggregate provides the minimum set of methods
// that are intrinsic to an iterator implementation.
type CoreIterator[T any] interface {
	SimpleIterator[T]
	// Seq returns the iterator as a Go iter.Seq iterator.
	Seq() iter.Seq[T]
	// Size is an estimate, where possible, of the number of elements remaining.
	Size() IteratorSize
	// SeqOK returns true if the Seq() method should be used to perform iterations.
	// Generally, using Seq() is the preferred method for efficiency reasons. However
	// there are situations where this is not the case and this method will return false.
	// For example, if the underlying iterator is based on a simple iterator, it is
	// slightly more efficient to stick to the simple iterator methods. Also, if simple
	// iterator methods have already been called against a Seq based iterator, calling
	// Seq() will cause inconsistent results, as it will restart the iterator from the
	// beginning, and so in these cases, SeqOK() will return false.
	SeqOK() bool
}

// CoreMutableIterator is an extension of CoreIterator which adds methods to facilitate
// iterator mutation.
type CoreMutableIterator[T any] interface {
	CoreIterator[T]
	// Set modifies the current value in place.
	Set(T)
	// Delete deletes the current value, which must be the last value returned by Next(). This
	// function may not be implemented for all iterator types, in which case it will panic.
	Delete()
}

type CoreIterator2[K any, V any] interface {
	CoreIterator[V]
	Seq2() iter.Seq2[K, V]
	Key() K
}

type CoreMutableIterator2[K any, V any] interface {
	CoreIterator2[K, V]
	Set(V)
	Delete()
}

// IteratorExtensions defines methods available to all iterators beyond the core functionality provided by CoreIterator.
type IteratorExtensions[T any] interface {
	// Chan returns the iterator as a channel.
	Chan() <-chan T
	// Collect collects all elements from the iterator into a slice.
	Collect() []T
	// Enumerate returns an iterator that enumerates the elements of this iterator, returning a tuple of the index and the value.
	Enumerate() Iterator2[int, T]
	Filter(func(T) bool) Iterator[T]
	FilterMorph(func(T) (T, bool)) Iterator[T]
	Morph(func(T) T) Iterator[T]
	Take(int) Iterator[T]
}

// Iterator2Extensions defines additional iterator methods that are specific
// to Iterator2.
type Iterator2Extensions[K any, V any] interface {
	Collect2() []KeyValue[K, V]
	Chan2() <-chan KeyValue[K, V]
	Take2(int) Iterator2[K, V]
	Morph2(func(K, V) (K, V)) Iterator2[K, V]
	FilterMorph2(func(K, V) (K, V, bool)) Iterator2[K, V]
}

// Generic iterator
type Iterator[T any] interface {
	CoreIterator[T]
	IteratorExtensions[T]
}

// Generic mutable iterator
type MutableIterator[T any] interface {
	CoreMutableIterator[T]
	IteratorExtensions[T]
}

type Iterator2[K any, V any] interface {
	CoreIterator2[K, V]
	IteratorExtensions[V]
	Iterator2Extensions[K, V]
}

type MutableIterator2[K any, V any] interface {
	CoreMutableIterator2[K, V]
	IteratorExtensions[V]
	Iterator2Extensions[K, V]
}

package iterator

import "iter"

// SimpleIterator defines a core set of methods for iterating over a collection of elements, of type T.
// More complete Iterator implementations can be built on this core set of methods.
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
	// that perform IO may not read the same data again, or may return no data at all.
	Reset()
}

// SimpleMutableIterator extends SimpleIterator by adding methods to support element mutation.
// More complete MutableIterator implementations can be built on this core set of methods.
type SimpleMutableIterator[T any] interface {
	SimpleIterator[T]
	// Set allows the value at the iterator's current position to be modified in place
	Set(T)
	// Delete deletes the current value, which must be the last value arrived at by Next().
	Delete()
}

// CoreIterator is an extension of SimpleIterator that in aggregate provides the minimum set of methods
// that are intrinsic to an iterator implementation. These methods are those that are concerned with
// handling the underlying data.
type CoreIterator[T any] interface {
	SimpleIterator[T]
	// Seq returns the iterator as a Go `iter.Seq` iterator. The iterator may be backed by
	// an `iter.Seq[T]` object, in which case that iterator object will typically be returned
	// directly. Otherwise, an `iter.Seq[T]` will be synthesised from the underlying iterator, typically
	// a SimpleIterator.
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
	// beginning, and so in these cases, SeqOK() should return false.
	SeqOK() bool
}

// CoreMutableIterator is an extension of CoreIterator which adds methods to facilitate
// iterator mutation.
type CoreMutableIterator[T any] interface {
	CoreIterator[T]
	// Set modifies the current value, the last value arrived at by a call to Next(), in place.
	Set(T)
	// Delete deletes the current value, which must be the last value returned by Next(). This
	// function may not be implemented for all iterator types, in which case it will panic.
	Delete()
}

// CoreIterator2 is an extension of CoreIterator that adds support for a second variable of type
// K (the "key") in addition to the existing value, of type V.
type CoreIterator2[K any, V any] interface {
	CoreIterator[V]
	// Seq returns the iterator as a Go `iter.Seq2` iterator. The iterator may be backed by
	// an `iter.Seq2[T]` object, in which case that iterator object will typically be returned
	// directly. Otherwise, an `iter.Seq2[T]` will be synthesised from the underlying iterator.
	Seq2() iter.Seq2[K, V]
	// Key returns the current iterator key.
	Key() K
}

// CoreMutableIterator2 is an extension of CoreIterator2 that adds support for mutability. The
// iterator value may be changed, and the current item may be deleted. There is no support
// for modifing the key.
type CoreMutableIterator2[K any, V any] interface {
	CoreIterator2[K, V]
	// Set will modify the current iterator value.
	Set(V)
	// Delete will remove the current iterator item. Calliing Next() is still required to advance
	// to the next item.
	Delete()
}

// IteratorExtensions defines methods available to all iterators beyond the core functionality provided by CoreIterator.
type IteratorExtensions[T any] interface {
	// Chan returns the iterator as a channel.
	Chan() <-chan T
	// Collect collects all elements from the iterator into a slice.
	Collect() []T
	// Enumerate returns an iterator that enumerates the elements of this iterator, returning an Iterator2 of the index and the value.
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
	Filter2(func(K, V) bool) Iterator2[K, V]
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

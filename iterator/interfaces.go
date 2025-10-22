package iterator

import "iter"

// SimpleIterator defines a core set of methods for iterating over a collection
// of elements, of type T. More complete iterator implementations can be built
// on this core set of methods.
type SimpleIterator[T any] interface {

	// Next sets the iterator's current value to be the first, and subsequent,
	// iterator elements. False is returned only when there are no more elements
	// (the current value remains unchanged).
	Next() bool

	// Value gets the current iterator value.
	Value() T

	// Abort stops the iterator; subsequent calls to Next() will return false.
	Abort()

	// Reset stops the iterator; subsequent calls to Next() will begin the
	// iterator from the start. Note not all iterators are guaranteed to return
	// the same sequence again, for example iterators that perform IO may not
	// read the same data again, or may return no data at all.
	Reset()
}

// SimpleMutableIterator extends [SimpleIterator] by adding methods to support
// element mutation. More complete [MutableIterator] implementations can be
// built on this core set of methods.
type SimpleMutableIterator[T any] interface {
	SimpleIterator[T]

	// Set modifies a value in place in the underlying collection.
	Set(T)

	// Delete deletes the current value, which must be the last value returned
	// by Next(). This function may not be implemented for all iterator types,
	// in which case it will return an [ErrDeleteNotImplemented] error.
	Delete()
}

// CoreIterator is an extension of [SimpleIterator] that in aggregate provides
// the minimum set of methods that are intrinsic to an iterator implementation,
// i.e. those methods that are concerned with interacting the underlying data.
type CoreIterator[T any] interface {
	SimpleIterator[T]

	// Seq returns the iterator as a Go [iter.Seq] iterator. The iterator may be
	// backed by an iter.Seq[T] object, in which case that iterator object will
	// typically be returned directly. Otherwise, an iter.Seq[T] will be
	// synthesised from the underlying iterator, typically a [SimpleIterator].
	Seq() iter.Seq[T]

	// Size is an estimate, where possible, of the number of elements remaining.
	Size() IteratorSize

	// SeqOK returns true if the Seq() method should be used to perform
	// iterations. Generally, using Seq() is the preferred method for efficiency
	// reasons. However there are situations where this is not the case and in
	// those cases this method will return false. For example, if the underlying
	// iterator is based on a simple iterator, it is slightly more efficient to
	// stick to the simple iterator methods. Also, if simple iterator methods
	// have already been called against a Seq based iterator, calling Seq() will
	// cause inconsistent results, as it will restart the iterator from the
	// beginning, and so in these cases, SeqOK() will return false.
	SeqOK() bool
}

// CoreMutableIterator is an extension of [CoreIterator] which adds methods to
// facilitate iterator mutation.
type CoreMutableIterator[T any] interface {
	CoreIterator[T]
	// Set modifies the current value, the last value arrived at by a call to
	// Next(), in place.
	Set(T)

	// Delete deletes the current value, which must be the last value returned
	// by Next(). This function may not be implemented for all iterator types,
	// in which case it will panic.
	Delete()
}

// CoreIterator2 is an extension of [CoreIterator] that adds support for a
// second variable of type K (the "key") in addition to the existing value, of
// type V.
type CoreIterator2[K any, V any] interface {
	CoreIterator[V]
	// Seq returns the iterator as a Go [iter.Seq2] iterator. The iterator may
	// be backed by an iter.Seq2[T] object, in which case that iterator object
	// will typically be returned directly. Otherwise, an iter.Seq2[T] will be
	// synthesized from the underlying iterator.
	Seq2() iter.Seq2[K, V]
	// Key returns the current iterator key.
	Key() K
}

// CoreMutableIterator2 is an extension of [CoreIterator2] that adds support for
// mutability. The iterator value may be changed, and the current item may be
// deleted. There is no support for modifying the key.
type CoreMutableIterator2[K any, V any] interface {
	CoreIterator2[K, V]
	// Set will modify the current iterator value.
	Set(V)
	// Delete will remove the current iterator item. Calling Next() is still
	// required to advance to the next item.
	Delete()
}

// IteratorExtensions defines methods available to all iterators beyond the core
// functionality provided by CoreIterator.
type IteratorExtensions[T any] interface {

	// Chan returns the iterator as a channel.
	Chan() <-chan T

	// Collect collects all elements from the iterator into a slice.
	Collect() []T

	// CollectInto adds all elements from the iterator into an existing slice.
	CollectInto(*[]T) []T

	// CollectIntoCap add elements from the iterator into an existing slice
	// until either the capacity of the slice is filled, or the iterator is
	// exhausted, which ever is first.
	CollectIntoCap(*[]T) []T

	// Enumerate returns an iterator that enumerates the elements of this
	// iterator, returning an Iterator2 of the index and the value.
	Enumerate() Iterator2[int, T]

	// Filter is a filtering method that creates a new iterator which contains a
	// subset of elements contained in the current one. This function takes a
	// predicate function p and only elements for which this function returns
	// true should be included in the filtered iterator.
	Filter(p func(T) bool) Iterator[T]

	// Morph is a mapping function that creates a new iterator which contains
	// the elements of the current iterator with the supplied function m applied
	// to each one. The type of the return value of m must be the same as that
	// of the elements of the current iterator. This is because of limitations
	// of Go generics. To apply a mapping that changes the type, see the
	// [iterator.Map] function.
	Morph(m func(T) T) Iterator[T]

	// FilterMorph is a filtering and mapping method that creates a new iterator
	// from an existing one by simultaneously transforming and filtering each
	// iterator element. The method takes a mapping function f that transforms
	// and filters each element. It does this by taking an input element value
	// and returning a new element value and a boolean flag. Only elements for
	// which this flag is true are included in the new iterator. E.g.
	//  itr := iterator.Of(0,1,2,3,4)
	//  itrMorph := itr.FilterMorph(func(v int) (int, bool) { return v*2, v%2 == 0})
	//  result := itrMorph.Collect() // []int{0,4,8}
	// Note that this function is not able to map elements to values of a
	// different type due to limitations of Go generics. For a filter mapping
	// function that can change the type, see the [iterator.FilterMap] function.
	FilterMorph(f func(T) (T, bool)) Iterator[T]

	// Take returns a variant of the current iterator that which returns at most
	// n elements. If the current iterator has less than or exactly n elements,
	// the returned iterator is equivalent to the input iterator.
	Take(n int) Iterator[T]

	// Any returns true if p returns true for at least one element in the iterator.
	Any(p func(T) bool) bool

	// All returns true if p returns true for all the elements in the iterator.
	All(p func(T) bool) bool

	// Fold1 combines the elements of the iterator into a single value using an
	// accumulation function. It takes an initial value init and an accumulation
	// function f. The iterator must have at least one element, or this method will
	// panic with [ErrEmptyIterator].  If there is only one element, this element be
	// returned. Otherwise, the first two elements are fed into the accumulation
	// function f. The result of this is combined with the next element to get the
	// next result and so on until the iterator is consumed. The final result is
	// returned. If the iterator is known to be of infinite size, this method will
	// panic with [ErrSizeInfinite].
	Fold1(f func(a, e T) T) T

	// Fold combines the elements of the iterator into a single value using an
	// accumulation function. It takes an initial value init and the accumulation
	// function f. If the iterator is empty, init is returned. Otherwise the initial
	// value is fed into the accumulation function f along with the first element
	// from the iterator. The result is then fed back into f along with the second
	// element, and so on until the iterator is consumed. The final result is the
	// final value returned by the function. If the iterator is known to be of
	// infinite size, this method will panic with [ErrSizeInfinite].
	Fold(init T, f func(a, e T) T) T

	// Intercalate1 is a variation on [Fold1] which combines the elements of the
	// iterator into a single value using an accumulation function and an
	// interspersed value. The iterator must have at least one element, otherwise it
	// will panic with [ErrEmptyIterator]. The accumulated result is initially set
	// to the first element, and is combined with subsequent elements as follows:
	//
	//	acc = f(f(acc,inter),e)
	//
	// where acc is he accumulated value, e is the next element and inter is the
	// inter parameter. The final value of acc will be the value returned. If the
	// iterator is known to be of infinite size, this function will panic with
	// [ErrSizeInfinite].
	Intercalate1(inter T, f func(a, e T) T) T

	// Intercalate is a variation on [Fold] which combines the elements of the
	// iterator into a single value using an accumulation function and an
	// interspersed value. If the iterator has no elements, the empty parameter is
	// returned. Otherwise, the accumulated result is initially set to the first
	// element, and is combined with subsequent elements as follows:
	//
	//	acc = f(f(acc,inter),e)
	//
	// where acc is he accumulated value, e is the next element and inter is the
	// inter parameter, and the final value of acc will be the value returned. If
	// the iterator is known to be of infinite size, this function will panic with
	// [ErrSizeInfinite].
	Intercalate(empty T, inter T, f func(a, e T) T) T
}

// Iterator2Extensions defines additional iterator methods that are specific
// to Iterator2.
type Iterator2Extensions[K any, V any] interface {

	// Collect2 collects all the element pairs from the iterator into a slice of
	// KeyValue objects.
	Collect2() []KeyValue[K, V]

	// Collect2Into collects all the element pairs from the iterator into the
	// slice of KeyValue objects pointed to by s. The final slice is returned.
	Collect2Into(s *[]KeyValue[K, V]) []KeyValue[K, V]

	// Collect2IntoCap collects all the element pairs from the iterator into the
	// slice of KeyValue objects pointed to by s, up to but not exceeding the
	// capacity of *s. The final slice is returned.
	Collect2IntoCap(s *[]KeyValue[K, V]) []KeyValue[K, V]

	// Chan2 returns the iterator as a channel of KeyValue objects. The iterator
	// is consumed in a goroutine which yields results to the channel.
	Chan2() <-chan KeyValue[K, V]

	// Filter2 is a filtering method that creates a new iterator which contains
	// a subset of element pairs contained by the current one. This function
	// takes a predicate function p and only element pairs for which this
	// function returns true should be included in the filtered iterator.
	Filter2(p func(K, V) bool) Iterator2[K, V]

	// Morph2 is a mapping function that creates a new iterator which contains
	// pairs of elements of the current iterator with the supplied function m
	// applied to each key and value. The type of the return value and key of m
	// must be of the same type as the kay and value of the current iterator.
	// This is because of limitations of Go generics. To apply a mapping that
	// changes the type, see the [iterator.Map2] function.
	Morph2(m func(K, V) (K, V)) Iterator2[K, V]

	// FilterMorph2 is a filtering and mapping method that creates a new
	// iterator from an existing one by simultaneously transforming and
	// filtering each iterator element pair. The method takes a mapping function
	// f that transforms and filters each element pair. It does this by taking
	// an input element key and value and returning a new element key and value
	// and a boolean flag. Only elements for which this flag is true are
	// included in the new iterator. E.g.
	//  inputMap := map[int]int{0: 2, 1: 4, 2: 6, 3: 8}
	//  itr := maps.Iter(inputMap).FilterMorph2(func(k, v int) (int, int, bool) {
	//    return k + 1, v * 2, (k+v)%2 == 0
	//  })
	//  output := iterator.CollectMap(itr) // map[int]int{1: 4, 3: 12}
	// Note that this function is not able to map element keys or values to
	// different types due to limitations of Go generics. For a filter mapping
	// function that can map to different types, see the [iterator.FilterMap2]
	// function.
	FilterMorph2(f func(K, V) (K, V, bool)) Iterator2[K, V]

	// Take2 returns a variant of the current iterator that which returns at
	// most n pairs of elements. If the current iterator has less than or
	// exactly n element pairs, the returned iterator is equivalent to the input
	// iterator.
	Take2(int) Iterator2[K, V]
}

// Top level iterator types

// Iterator is a generic iterator type, facilitating iteration over single
// elements of a generic type plus some utility methods. It consists of methods
// from [CoreIterator], plus the ones from [IteratorExtensions].
type Iterator[T any] interface {
	CoreIterator[T]
	IteratorExtensions[T]
}

// MutableIterator is a generic iterator type, facilitating iteration over
// single elements of a generic type that also supports mutation of the
// underlying value. Elements can be modified in place or removed from their
// underlying collection. This type also includes some utility methods. It
// consists of methods from [CoreMutableIterator], plus the ones from
// [IteratorExtensions].
type MutableIterator[T any] interface {
	CoreMutableIterator[T]
	IteratorExtensions[T]
}

// Iterator2 is a generic iterator type, facilitating iteration over pairs of
// elements of different generic types. One of these is the "value", and the
// other is the "key" which may be something like a map key or slice index. It
// also has some utility methods. It consists of methods from [CoreIterator2],
// [IteratorExtensions] and [Iterator2Extensions].
type Iterator2[K any, V any] interface {
	CoreIterator2[K, V]
	IteratorExtensions[V]
	Iterator2Extensions[K, V]
}

// MutableIterator2 is a generic iterator type, facilitating iteration over
// pairs of elements of different generic types. One of these is the "value",
// and the other is the "key" which may be something like a map key or slice
// index. Two mutation operations are supported; the modification of the the
// "value" element, and removal of the element pair from the underlying
// collection. It also has some utility methods. It consists of methods from
// [CoreIterator2], [IteratorExtensions] and [Iterator2Extensions].
type MutableIterator2[K any, V any] interface {
	CoreMutableIterator2[K, V]
	IteratorExtensions[V]
	Iterator2Extensions[K, V]
}

package iterator

import (
	"errors"
	"iter"

	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/option"
)

var ErrSizeInfinite = errors.New("cannot consume an infinite iterator")
var ErrInvalidIteratorSizeType = errors.New("invalid iterator size type")
var ErrInvalidIteratorRange = errors.New("invalid iterator range")
var ErrDeleteNotImplemented = errors.New("delete not implemented")
var ErrEmptyIterator = errors.New("iterator is empty")

// Seq transforms any generic [SimpleIterator] into a native [iter.Seq] iterator.
func Seq[T any](i SimpleIterator[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for i.Next() {
			if !yield(i.Value()) {
				i.Abort()
				break
			}
		}
	}
}

// Seq2 transforms a [SimpleIterator] of [KeyValue] pairs into a native [iter.Seq2] iterator.
func Seq2[K any, V any](i SimpleIterator[KeyValue[K, V]]) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for i.Next() {
			if !yield(i.Value().Key, i.Value().Value) {
				i.Abort()
				break
			}
		}
	}
}

// DefaultIterator wraps a [CoreIterator] and provides default implementations
// of the iterator methods defined in [IteratorExtensions] to provide a complete
// [Iterator] implementation. All the extension methods are written purely in
// terms of the methods of [CoreIterator], any typically delegate to the
// function in the iterator package of the same name.
type DefaultIterator[T any] struct {
	CoreIterator[T]
}

// NewDefaultIterator builds an [Iterator] from a [CoreIterator] by adding the methods defined in
// [IteratorExtensions].
func NewDefaultIterator[T any](citr CoreIterator[T]) DefaultIterator[T] {
	return DefaultIterator[T]{CoreIterator: citr}
}

// Chan returns the iterator as a channel.  The iterator is consumed in a
// goroutine which yields results to the channel.
func (di DefaultIterator[T]) Chan() <-chan T {
	return Chan(di)
}

// Collect collects all elements from the iterator into a slice.
func (di DefaultIterator[T]) Collect() []T {
	return Collect(di)
}

// CollectInto adds all elements from the iterator into an existing slice.
func (di DefaultIterator[T]) CollectInto(slice *[]T) []T {
	return CollectInto(di, slice)
}

// CollectIntoCap add elements from the iterator into an existing slice
// until either the capacity of the slice is filled, or the iterator is
// exhausted, which ever is first.
func (di DefaultIterator[T]) CollectIntoCap(slice *[]T) []T {
	return CollectIntoCap(di, slice)
}

// Enumerate returns an iterator that enumerates the elements of this
// iterator, returning an Iterator2 of the index and the value.
func (di DefaultIterator[T]) Enumerate() Iterator2[int, T] {
	return Enumerate(di)
}

// Filter is a filtering method that creates a new iterator which contains a
// subset of elements contained in the current one. This function takes a
// predicate function p and only elements for which this function returns
// true should be included in the filtered iterator.
func (di DefaultIterator[T]) Filter(predicate func(T) bool) Iterator[T] {
	return Filter(di, predicate)
}

// Morph is a mapping function that creates a new iterator which contains
// the elements of the current iterator with the supplied function m applied
// to each one. The type of the return value of m must be the same as that
// of the elements of the current iterator. This is because of limitations
// of Go generics. To apply a mapping that changes the type, see the
// [Map] function.
func (di DefaultIterator[T]) Morph(mapping func(T) T) Iterator[T] {
	return Map(di, mapping)
}

// FilterMorph is a filtering and mapping method that creates a new iterator
// from an existing one by simultaneously transforming and filtering each
// iterator element. The method takes a mapping function f that transforms and
// filters each element. It does this by taking an input element value and
// returning a new element value and a boolean flag. Only elements for which
// this flag is true are included in the new iterator. E.g.
//
//	itr := iterator.Of(0,1,2,3,4)
//	itrMorph := itr.FilterMorph(func(v int) (int, bool) { return v*2, v%2 == 0})
//	result := itrMorph.Collect() // []int{0,4,8}
//
// Note that this function is not able to map elements to values of a different
// type due to limitations of Go generics. For a filter mapping function that
// can change the type, see the [iterator.FilterMap] function.
func (di DefaultIterator[T]) FilterMorph(mapping func(T) (T, bool)) Iterator[T] {
	return FilterMap(di, mapping)
}

// Take returns a variant of the current iterator that which returns at most
// n elements. If the current iterator has less than or exactly n elements,
// the returned iterator is equivalent to the input iterator.
func (di DefaultIterator[T]) Take(n int) Iterator[T] {
	return Take(n, di)
}

// All returns true if p returns true for all the elements in the iterator.
// This method short circuits and does not execute in constant time; the
// iterator is aborted after the first value for which the predicate returns
// false.
func (di DefaultIterator[T]) All(predicate func(T) bool) bool {
	return All(di, predicate)
}

// Any returns true if p returns true for at least one element in the iterator.
// This method short circuits and does not execute in constant time; the
// iterator is aborted after the first value for which the predicate returns
// true.
func (di DefaultIterator[T]) Any(predicate func(T) bool) bool {
	return Any(di, predicate)
}

// Fold1 combines the elements of the iterator into a single value using an
// accumulation function. It takes an initial value init and an accumulation
// function f. The iterator must have at least one element, or this method will
// panic with [ErrEmptyIterator].  If there is only one element, this element be
// returned. Otherwise, the first two elements are fed into the accumulation
// function f. The result of this is combined with the next element to get the
// next result and so on until the iterator is consumed. The final result is
// returned. If the iterator is known to be of infinite size, this method will
// panic with [ErrSizeInfinite].
func (di DefaultIterator[T]) Fold1(f func(a, e T) T) T {
	return Fold1(di, f)
}

// Fold combines the elements of the iterator into a single value using an
// accumulation function. It takes an initial value init and the accumulation
// function f. If the iterator is empty, init is returned. Otherwise the initial
// value is fed into the accumulation function f along with the first element
// from the iterator. The result is then fed back into f along with the second
// element, and so on until the iterator is consumed. The final result is the
// final value returned by the function. If the iterator is known to be of
// infinite size, this method will panic with [ErrSizeInfinite].
func (di DefaultIterator[T]) Fold(init T, f func(a, e T) T) T {
	return Fold(di, init, f)
}

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
func (di DefaultIterator[T]) Intercalate(empty T, inter T, f func(a, e T) T) T {
	return Intercalate(di, empty, inter, f)
}

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
func (di DefaultIterator[T]) Intercalate1(inter T, f func(a, e T) T) T {
	return Intercalate1(di, inter, f)
}

// DefaultMutableIterator wraps a CoreMutableIterator together with a DefaultIterator to provide
// and implementation of MutableIterator.
type DefaultMutableIterator[T any] struct {
	CoreMutableIterator[T]
	DefaultIterator[T]
}

// NewDefaultMutableIterator builds a MutableIterator from a CoreMutableIterator by adding the methods
// of IteratorExtensions.
func NewDefaultMutableIterator[T any](citr CoreMutableIterator[T]) DefaultMutableIterator[T] {
	return DefaultMutableIterator[T]{CoreMutableIterator: citr, DefaultIterator: DefaultIterator[T]{CoreIterator: citr}}
}

// DefaultMutableIterator2 wraps a CoreMutableIterator together with a DefaultIterator to provide
// and implementation of MutableIterator.
type DefaultMutableIterator2[K any, V any] struct {
	CoreMutableIterator2[K, V]
	DefaultIterator2[K, V]
}

func NewDefaultMutableIterator2[K any, V any](citr CoreMutableIterator2[K, V]) DefaultMutableIterator2[K, V] {
	return DefaultMutableIterator2[K, V]{CoreMutableIterator2: citr, DefaultIterator2: NewDefaultIterator2(citr)}
}

// New builds an Iterator from a standard library [iter.Seq]
func New[T any](seq iter.Seq[T]) Iterator[T] {
	return NewDefaultIterator(NewSeqCoreIterator(seq))
}

// New builds an Iterator from a standard library iter.Seq plus a function that
// returns the number of items left in the iterator.
func NewWithSize[T any](seq iter.Seq[T], size func() IteratorSize) Iterator[T] {
	return NewDefaultIterator(NewSeqCoreIteratorWithSize(seq, size))
}

// NewFromSimple builds an Iterator from a SimpleIterator.
func NewFromSimple[T any](simple SimpleIterator[T]) Iterator[T] {
	return NewDefaultIterator(NewSimpleCoreIterator(simple))
}

// NewFromSimpleWithSize builds an Iterator from a SimpleIterator plus a function that
// returns the number of items left in the iterator.
func NewFromSimpleWithSize[T any](simple SimpleIterator[T], size func() IteratorSize) Iterator[T] {
	return NewDefaultIterator(NewSimpleCoreIteratorWithSize(simple, size))
}

// DefaultIterator2 is an [Iterator2] implementation which embeds [CoreIterator2] and [IteratorExtensions]
// interfaces, and implements the methods of [IteratorExtensions2] in terms of [CoreIterator2]. It can be
// constructed solely from a [CoreIterator2] implementation by the [NewDefaultIterator2] function.
type DefaultIterator2[K any, V any] struct {
	CoreIterator2[K, V]
	IteratorExtensions[V]
}

// Collect2 collects all the element pairs from the iterator into a slice of
// [KeyValue] objects.
func (di2 DefaultIterator2[K, V]) Collect2() []KeyValue[K, V] {
	return Collect2(di2.CoreIterator2)
}

// Collect2Into collects all the element pairs from the iterator into the
// slice of [KeyValue] objects pointed to by s. The final slice is returned.
func (di2 DefaultIterator2[K, V]) Collect2Into(s *[]KeyValue[K, V]) []KeyValue[K, V] {
	return Collect2Into(di2.CoreIterator2, s)
}

// Collect2IntoCap collects all the element pairs from the iterator into the
// slice of KeyValue objects pointed to by s, up to but not exceeding the
// capacity of *s. The final slice is returned.
func (di2 DefaultIterator2[K, V]) Collect2IntoCap(s *[]KeyValue[K, V]) []KeyValue[K, V] {
	return Collect2IntoCap(di2.CoreIterator2, s)
}

// Chan2 returns the iterator as a channel of KeyValue objects. The iterator
// is consumed in a goroutine which yields results to the channel.
func (di2 DefaultIterator2[K, V]) Chan2() <-chan KeyValue[K, V] {
	return Chan2(di2.CoreIterator2)
}

// Take2 returns a variant of the current iterator that which returns at
// most n pairs of elements. If the current iterator has less than or
// exactly n element pairs, the returned iterator is equivalent to the input
// iterator.
func (di2 DefaultIterator2[K, V]) Take2(n int) Iterator2[K, V] {
	return Take2(n, di2.CoreIterator2)
}

// Filter2 is a filtering method that creates a new iterator which contains
// a subset of element pairs contained by the current one. This function
// takes a predicate function p and only element pairs for which this
// function returns true should be included in the filtered iterator.
func (di2 DefaultIterator2[K, V]) Filter2(f func(k K, v V) bool) Iterator2[K, V] {
	return Filter2(di2, f)
}

// Morph2 is a mapping function that creates a new iterator which contains
// pairs of elements of the current iterator with the supplied function m
// applied to each key and value. The type of the return value and key of m
// must be of the same type as the kay and value of the current iterator.
// This is because of limitations of Go generics. To apply a mapping that
// changes the type, see the [iterator.Map2] function.
func (di2 DefaultIterator2[K, V]) Morph2(f func(k K, v V) (K, V)) Iterator2[K, V] {
	return Map2(di2, f)
}

// FilterMorph2 is a filtering and mapping method that creates a new
// iterator from an existing one by simultaneously transforming and
// filtering each iterator element pair. The method takes a mapping function
// f that transforms and filters each element pair. It does this by taking
// an input element key and value and returning a new element key and value
// and a boolean flag. Only elements for which this flag is true are
// included in the new iterator.
//
// Note that this function is not able to map element keys or values to
// different types due to limitations of Go generics. For a filter mapping
// function that can map to different types, see the [iterator.FilterMap2]
// function.
func (di2 DefaultIterator2[K, V]) FilterMorph2(f func(K, V) (K, V, bool)) Iterator2[K, V] {
	return FilterMap2(di2, f)
}

// NewDefaultIterator2 constructs an [Iterator2] by wrapping a [CoreIterator2]
// with a [DefaultIterator2] to add the additional  [IteratorExtensions] and
// [IteratorExtensions2] methods to provide the complete [Iterator2]
// implementation
func NewDefaultIterator2[K any, V any](core CoreIterator2[K, V]) DefaultIterator2[K, V] {
	return DefaultIterator2[K, V]{
		IteratorExtensions: NewDefaultIterator(core),
		CoreIterator2:      core,
	}
}

// New2 creates an Iterator2 from a standard library iter.Seq2. Iterator Size is unknown.
func New2[K any, V any](seq2 iter.Seq2[K, V]) Iterator2[K, V] {
	return NewDefaultIterator2(CoreIterator2[K, V](NewSeqCoreIterator2(seq2)))
}

// New2WithSize creates an Iterator2 from a standard library iter.Seq2 and a size function that
// returns the remaining items in the iterator.
func New2WithSize[K any, V any](seq2 iter.Seq2[K, V], size func() IteratorSize) Iterator2[K, V] {
	return NewDefaultIterator2(NewSeqCoreIterator2WithSize(seq2, size))
}

// MakeIterator creates a generic iterator from a simple iterator. Provides an implementation
// of additional Iterator methods.
//
// Deprecated: use [NewFromSimple]
func MakeIterator[T any](base SimpleIterator[T]) Iterator[T] {
	return DefaultIterator[T]{NewSimpleCoreIterator(base)}
}

func NewSliceCoreIteratorRef[T any](slice *[]T) CoreMutableIterator[*T] {
	return &sliceIterRef[T]{sliceIter: sliceIter[T]{slice: slice, index: 0}}
}

// CollectInto collects all elements from an iterator into a slice a pointer to which is passed.
// Elements are appended to any existing data in the slice. The slice referenced may be reallocated
// as the append function is used to add elements to the slice. The slice may be a nil slice. For
// convenience, the final slice is returned. If the iterator is known to have an infinite size, this
// function will panic.
func CollectInto[T any](iter CoreIterator[T], slice *[]T) []T {
	if iter.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	if iter.SeqOK() {
		for v := range iter.Seq() {
			*slice = append(*slice, v)
		}
	} else {
		for iter.Next() {
			*slice = append(*slice, iter.Value())
		}
	}
	return *slice
}

// CollectIntoCap appends elements from an iterator into a slice a pointer to which is passed, up to
// a maximum of the capacity for the slice. Elements are not added beyond the capacity. Elements are
// appended to any existing data in the slice. The slice may be a nil slice, in which case no
// elements will be added. For convenience, tbe final slice is returned.
func CollectIntoCap[T any](iter CoreIterator[T], slice *[]T) []T {
	max := cap(*slice) - len(*slice)
	if max > 0 {
		if iter.SeqOK() {
			for v := range iter.Seq() {
				if max <= 0 {
					break
				}
				*slice = append(*slice, v)
				max--
			}
		} else {
			for iter.Next() {
				if max <= 0 {
					break
				}
				*slice = append(*slice, iter.Value())
				max--
			}
		}
	}
	return *slice
}

// Collect2Into collects all element pairs from an Iterator2 into a slice of KeyValue pairs, a
// pointer to which is passed. Pairs are appended to any existing data in the slice. The slice
// referenced may be reallocated as the append function is used to add pairs to the slice. The slice
// may be a nil slice. For convenience, the final slice is returned. If the iterator is known to
// have an infinite size, this function will panic.
func Collect2Into[K any, V any](iter CoreIterator2[K, V], slice *[]KeyValue[K, V]) []KeyValue[K, V] {
	if iter.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	if iter.SeqOK() {
		for k, v := range iter.Seq2() {
			*slice = append(*slice, KVOf(k, v))
		}
	} else {
		for iter.Next() {
			*slice = append(*slice, KVOf(iter.Key(), iter.Value()))
		}
	}
	return *slice
}

// Collect2IntoCap collects all element pairs from an Iterator2 into a slice of KeyValue pairs, a
// pointer to which is passed. Pairs are appended to any existing data in the slice. The slice
// referenced may be reallocated as the append function is used to add pairs to the slice. The slice
// may be a nil slice. For convenience, the final slice is returned. If the iterator is known to
// have an infinite size, this function will panic.
func Collect2IntoCap[K any, V any](iter CoreIterator2[K, V], slice *[]KeyValue[K, V]) []KeyValue[K, V] {
	if iter.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	max := cap(*slice) - len(*slice)
	if max > 0 {
		if iter.SeqOK() {
			for k, v := range iter.Seq2() {
				if max <= 0 {
					break
				}
				*slice = append(*slice, KVOf(k, v))
				max--
			}
		} else {
			for iter.Next() {
				if max <= 0 {
					break
				}
				*slice = append(*slice, KVOf(iter.Key(), iter.Value()))
				max--
			}
		}
	}
	return *slice
}

// CollectIntoMap collects all element pairs from an Iterator2 into the map passed, populating it
// with key and value pairs. All keys should be unique; any duplicate keys will result in only the
// latest key/value pair with that key being preserved. The map may be nil, in which case a new map
// will be created. The final map is returned. If the iterator is known to have an infinite size,
// this function will panic.
func CollectIntoMap[K comparable, V any](iter CoreIterator2[K, V], m map[K]V) map[K]V {
	if iter.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	if m == nil {
		m = make(map[K]V)
	}
	if iter.SeqOK() {
		for k, v := range iter.Seq2() {
			m[k] = v
		}
	} else {
		for iter.Next() {
			m[iter.Key()] = iter.Value()
		}
	}
	return m
}

// Collect all elements from an iterator into a slice. If the iterator is known to be
// of infinite size, this function will panic.
func Collect[T any](iter CoreIterator[T]) []T {
	result := make([]T, 0, iter.Size().Allocate())
	return CollectInto(iter, &result)
}

// Collect2 collects all elements from an Iterator2 into a slice of KeyValue pairs. If the iterator
// is known to be of infinite size, this function will panic.
func Collect2[K any, V any](itr CoreIterator2[K, V]) []KeyValue[K, V] {
	result := make([]KeyValue[K, V], 0, itr.Size().Allocate())
	return Collect2Into(itr, &result)
}

// CollectMap collects all elements from an Iterator2 into a map of key and value pairs. All keys
// should be unique; any duplicate keys will result in only the latest key/value pair with that key
// being preserved. If the iterator is known to be of infinite size, this function will panic.
func CollectMap[K comparable, V any](itr CoreIterator2[K, V]) map[K]V {
	return CollectIntoMap(itr, nil)
}

// CollectResults collects all elements from an iterator of results into a result of slice of the iterator's underlying type
// If the iterator returns an error result at any point, this call will terminate and return that error, along with the elements
// collected thus far.
func CollectResults[T any](iter CoreIterator[result.Result[T]]) ([]T, error) {
	collectResult := make([]T, 0, iter.Size().Allocate())
	for res := range iter.Seq() {
		if res.IsError() {
			return collectResult, res.GetErr()
		}
		collectResult = append(collectResult, res.Get())
	}
	return collectResult, nil
}

// PartitionResults collects the elements from an iterator of result types into two slices, one of
// successful (nil error) values, and the other of error values.
func PartitionResults[T any](iter CoreIterator[result.Result[T]]) ([]T, []error) {
	values := make([]T, 0, iter.Size().Allocate())
	var errs []error
	for res := range iter.Seq() {
		if res.IsError() {
			errs = append(errs, res.GetErr())
		} else {
			values = append(values, res.Must())
		}
	}
	return values, errs
}

// All returns true if predicate returns true for every value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns false.
func All[T any](itr CoreIterator[T], predicate func(v T) bool) bool {
	if itr.SeqOK() {
		for v := range itr.Seq() {
			if !predicate(v) {
				return false
			}
		}
	} else {
		for itr.Next() {
			if !predicate(itr.Value()) {
				return false
			}
		}
	}
	return true
}

// Any returns true if predicate returns true for any value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns true.
func Any[T any](itr CoreIterator[T], predicate func(v T) bool) bool {
	if itr.SeqOK() {
		for v := range itr.Seq() {
			if predicate(v) {
				return true
			}
		}
	} else {
		for itr.Next() {
			if predicate(itr.Value()) {
				return true
			}
		}
	}
	return false
}

// Fold combines the elements of an iterator into a single value using an
// accumulation function. It takes an iterator, an initial value init and an
// accumulation function f. If the iterator is empty, init is returned.
// Otherwise the initial value is fed into the accumulation function f along
// with the first element from the iterator. The result is then fed back into f
// along with the second element, and so on. The final result is the final value
// returned by the function. If the iterator is known to be of infinite size,
// this function will panic with [ErrSizeInfinite].
func Fold[U any, T any](itr CoreIterator[T], init U, f func(a U, e T) U) U {
	if itr.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	acc := init
	if itr.SeqOK() {
		for e := range itr.Seq() {
			acc = f(acc, e)
		}
	} else {
		for itr.Next() {
			acc = f(acc, itr.Value())
		}
	}
	return acc
}

// Fold1 combines the elements of an iterator into a single value using an
// accumulation function. It takes an iterator, an initial value init and an
// accumulation function f. The iterator must have at least one element, or this
// function will panic with [ErrEmptyIterator].  If there is only one element,
// this element be returned. Otherwise, the first two elements are fed into the
// accumulation function f. The result of this is combined with the next element
// to get the next result and so on until the iterator is consumed. The final
// result is returned. If the iterator is known to be of infinite size, this
// function will panic with [ErrSizeInfinite].
func Fold1[T any](itr CoreIterator[T], f func(a, e T) T) T {
	if itr.Size().IsInfinite() {
		panic(ErrSizeInfinite)
	}
	if itr.SeqOK() {
		acc := option.Empty[T]()
		for e := range itr.Seq() {
			if acc.HasValue() {
				acc.Set(f(acc.Get(), e))
			} else {
				acc.Set(e)
			}
		}
		if acc.IsEmpty() {
			panic(ErrEmptyIterator)
		}
		return acc.Get()
	} else {
		if !itr.Next() {
			panic(ErrEmptyIterator)
		}
		acc := itr.Value()
		for itr.Next() {
			acc = f(acc, itr.Value())
		}
		return acc
	}
}

// Intercalate1 is a variation on [Fold1] which combines the elements of an
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
func Intercalate1[T any](itr CoreIterator[T], inter T, f func(a, e T) T) T {
	interFold := func(a, e T) T { return f(f(a, inter), e) }
	return Fold1(itr, interFold)
}

// Intercalate is a variation on [Fold] which combines the elements of an
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
func Intercalate[T any](itr CoreIterator[T], empty T, inter T, f func(a, e T) T) T {
	first := true
	interFold := func(a, e T) T {
		if first {
			first = false
			return e
		} else {
			return f(f(a, inter), e)
		}
	}
	return Fold(itr, empty, interFold)
}

func Cycle[T any](itr CoreIterator[T]) Iterator[T] {
	seq := func(yield func(T) bool) {
		for {
			hasItems := false
			for item := range itr.Seq() {
				if !yield(item) {
					break
				}
				hasItems = true
			}
			if !hasItems {
				break
			}
			itr.Reset()
		}
	}
	return New(seq)
}

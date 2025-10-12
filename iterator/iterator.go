// Iterators and generators
package iterator

import (
	"errors"
	"iter"

	"github.com/robdavid/genutil-go/errors/result"
)

// The largest slice capacity we are prepared to allocate to collect
// iterators of uncertain size.
const maxUncertainAllocation = 100000

var ErrAllocationSizeInfinite = errors.New("cannot allocate storage for an infinite iterator")
var ErrInvalidIteratorSizeType = errors.New("invalid iterator size type")
var ErrInvalidIteratorRange = errors.New("invalid iterator range")
var ErrDeleteNotImplemented = errors.New("delete not implemented")

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

// KeyValue holds a key value pair
type KeyValue[K any, V any] struct {
	Key   K
	Value V
}

// KVOf constructs a key value pair
func KVOf[K any, V any](key K, value V) KeyValue[K, V] {
	return KeyValue[K, V]{key, value}
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

// DefaultIterator wraps a CoreIterator and provides the additional functions defined in IteratorExtensions
// to provide a complete Iterator implementation.
type DefaultIterator[T any] struct {
	CoreIterator[T]
}

// NewDefaultIterator builds an [Iterator] from a [CoreIterator] by adding the methods defined in
// [IteratorExtensions].
func NewDefaultIterator[T any](citr CoreIterator[T]) DefaultIterator[T] {
	return DefaultIterator[T]{CoreIterator: citr}
}

func (di DefaultIterator[T]) Chan() <-chan T {
	return Chan(di)
}

func (di DefaultIterator[T]) Collect() []T {
	return Collect(di)
}

func (di DefaultIterator[T]) CollectInto(slice *[]T) []T {
	return CollectInto(di, slice)
}

func (di DefaultIterator[T]) CollectIntoCap(slice *[]T) []T {
	return CollectIntoCap(di, slice)
}

func (di DefaultIterator[T]) Enumerate() Iterator2[int, T] {
	return Enumerate(di)
}

func (di DefaultIterator[T]) Filter(predicate func(T) bool) Iterator[T] {
	return Filter(di, predicate)
}

func (di DefaultIterator[T]) FilterMorph(mapping func(T) (T, bool)) Iterator[T] {
	return FilterMap(di, mapping)
}

func (di DefaultIterator[T]) Morph(mapping func(T) T) Iterator[T] {
	return Map(di, mapping)
}

func (di DefaultIterator[T]) Take(n int) Iterator[T] {
	return Take(n, di)
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

// NewFromSimple builds an Iterator from a SimpleIterator
func NewFromSimple[T any](simple SimpleIterator[T]) Iterator[T] {
	return NewDefaultIterator(NewSimpleCoreIterator(simple))
}

// NewFromSimpleWithSize builds an Iterator from a SimpleIterator plus a function that
// returns the number of items left in the iterator.
func NewFromSimpleWithSize[T any](simple SimpleIterator[T], size func() IteratorSize) Iterator[T] {
	return NewDefaultIterator(NewSimpleCoreIteratorWithSize(simple, size))
}

type DefaultIterator2[K any, V any] struct {
	CoreIterator2[K, V]
	IteratorExtensions[V]
}

func (di2 DefaultIterator2[K, V]) Collect2() []KeyValue[K, V] {
	return Collect2(di2.CoreIterator2)
}

func (di2 DefaultIterator2[K, V]) Collect2Into(s *[]KeyValue[K, V]) []KeyValue[K, V] {
	return Collect2Into(di2.CoreIterator2, s)
}

func (di2 DefaultIterator2[K, V]) Collect2IntoCap(s *[]KeyValue[K, V]) []KeyValue[K, V] {
	return Collect2IntoCap(di2.CoreIterator2, s)
}

func (di2 DefaultIterator2[K, V]) Chan2() <-chan KeyValue[K, V] {
	return Chan2(di2.CoreIterator2)
}

func (di2 DefaultIterator2[K, V]) Take2(n int) Iterator2[K, V] {
	return Take2(n, di2.CoreIterator2)
}

func (di2 DefaultIterator2[K, V]) Filter2(f func(k K, v V) bool) Iterator2[K, V] {
	return Filter2(di2, f)
}

func (di2 DefaultIterator2[K, V]) Morph2(f func(k K, v V) (K, V)) Iterator2[K, V] {
	return Map2(di2, f)
}

func (di2 DefaultIterator2[K, V]) FilterMorph2(f func(K, V) (K, V, bool)) Iterator2[K, V] {
	return FilterMap2(di2, f)
}

func NewDefaultIterator2[K any, V any](core CoreIterator2[K, V]) DefaultIterator2[K, V] {
	return DefaultIterator2[K, V]{IteratorExtensions: NewDefaultIterator(core), CoreIterator2: core}
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
// Deprecated: use NewFromSimple
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
		panic(ErrAllocationSizeInfinite)
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
		panic(ErrAllocationSizeInfinite)
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
		panic(ErrAllocationSizeInfinite)
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
		panic(ErrAllocationSizeInfinite)
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

// All returns true if `predicate` returns true for every value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns false.
func All[T any](iter CoreIterator[T], predicate func(v T) bool) bool {
	for v := range iter.Seq() {
		if !predicate(v) {
			return false
		}
	}
	return true
}

// Any returns true if `predicate` returns true for any value returned
// by the iterator. This function short circuits and does not
// execute in constant time; the iterator is aborted after the
// first value for which the predicate returns true.
func Any[T any](iter CoreIterator[T], predicate func(v T) bool) bool {
	for v := range iter.Seq() {
		if predicate(v) {
			return true
		}
	}
	return false
}

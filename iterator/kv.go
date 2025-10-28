package iterator

import "iter"

// KeyValue holds a key value pair
type KeyValue[K any, V any] struct {
	Key   K
	Value V
}

// KVOf constructs a key value pair
func KVOf[K any, V any](key K, value V) KeyValue[K, V] {
	return KeyValue[K, V]{key, value}
}

type kvIter[K any, V any] struct {
	base CoreIterator2[K, V]
}

func (pi *kvIter[K, V]) Next() bool {
	return pi.base.Next()
}

func (pi *kvIter[K, V]) Value() KeyValue[K, V] {
	return KVOf(pi.base.Key(), pi.base.Value())
}

func (pi *kvIter[K, V]) Size() IteratorSize {
	return pi.base.Size()
}

func (pi *kvIter[K, V]) Abort() {
	pi.base.Abort()
}

func (pi *kvIter[K, V]) Reset() {
	pi.base.Reset()
}

func (pi *kvIter[K, V]) Seq() iter.Seq[KeyValue[K, V]] {
	return func(yield func(KeyValue[K, V]) bool) {
		for k, v := range pi.base.Seq2() {
			if !yield(KVOf(k, v)) {
				break
			}
		}
	}
}

func (pi *kvIter[K, V]) Chan() <-chan KeyValue[K, V] {
	return Chan2(pi.base)
}

func (pi *kvIter[K, V]) SeqOK() bool {
	return pi.base.SeqOK()
}

// AsKV takes a CoreIterator2 and returns an Iterator where each value is
// a KeyValue pair.
func AsKV[K any, V any](iter2 CoreIterator2[K, V]) Iterator[KeyValue[K, V]] {
	return NewDefaultIterator(&kvIter[K, V]{iter2})
}

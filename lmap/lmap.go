package lmap

import (
	"iter"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/list"
)

type value[K comparable, V any] struct {
	node  *list.Node[K]
	value V
}

// makeValue creates a new value struct with the given node and value.
func makeValue[K comparable, V any](k *list.Node[K], v V) value[K, V] {
	return value[K, V]{k, v}
}

// LinkedMap is map combined with linked list of keys which maintains a consistent
// key order. Keys are placed in the order they are first inserted.
type LinkedMap[K comparable, V any] struct {
	kv   map[K]value[K, V]
	keys list.List[K]
}

// Make creates a new empty LinkedMap instance.
func Make[K comparable, V any]() LinkedMap[K, V] {
	return LinkedMap[K, V]{
		kv:   make(map[K]value[K, V]),
		keys: list.Make[K](),
	}
}

// FromSeq2 creates a new LinkedMap instance populated with keys and values taken from
// the provided [iter.Seq2][K,V] iterator.
func FromSeq2[K comparable, V any](itr iter.Seq2[K, V]) LinkedMap[K, V] {
	result := Make[K, V]()
	for k, v := range itr {
		result.Put(k, v)
	}
	return result
}

// From creates a new LinkedMap instance populated with keys and values taken from
// the provided [iterator.Iterator2][K,V] iterator.
func From[K comparable, V any](itr iterator.Iterator2[K, V]) LinkedMap[K, V] {
	return FromSeq2(itr.Seq2())
}

// FromSeqKeys creates a new LinkedMap instance populated with keys from the
// supplied [iter.Seq][K] iterator with each associated value computed via the
// supplied function.
func FromSeqKeys[K comparable, V any](iterKeys iter.Seq[K], valueFn func(K) V) LinkedMap[K, V] {
	result := Make[K, V]()
	for k := range iterKeys {
		result.Put(k, valueFn(k))
	}
	return result
}

// FromSeqKeys creates a new LinkedMap instance populated with keys from the
// supplied [iterator.Iterator][K] iterator with each associated value computed
// via the supplied function.
func FromIterKeys[K comparable, V any](iterKeys iterator.Iterator[K], valueFn func(K) V) LinkedMap[K, V] {
	return FromSeqKeys(iterKeys.Seq(), valueFn)
}

// FromKeys creates a new LinkedMap instance populated with keys from the
// supplied slice of keys with each associated value computed via the supplied
// function.
func FromKeys[K comparable, V any](keys []K, valueFn func(K) V) LinkedMap[K, V] {
	result := Make[K, V]()
	for _, k := range keys {
		result.Put(k, valueFn(k))
	}
	return result
}

// Len returns the number of elements in the map.
func (lm LinkedMap[K, V]) Len() int {
	return len(lm.kv)
}

// IsEmpty returns true if there are no elements in the map.
func (lm LinkedMap[K, V]) IsEmpty() bool {
	return len(lm.kv) == 0
}

// Put places a key and value pair into the map, either adding it as a new
// entry if the key is not already in the map, or replacing an existing one.
func (lm *LinkedMap[K, V]) Put(k K, v V) {
	if current, ok := lm.kv[k]; ok {
		current.value = v
		lm.kv[k] = current
	} else {
		if lm.kv == nil {
			lm.kv = make(map[K]value[K, V])
		}
		lm.keys.Append(k)
		lm.kv[k] = makeValue(lm.keys.Last(), v)
	}
}

// Get returns the value in the map stored for key k. If key k is not present,
// the zero value of type V is returned.
func (lm LinkedMap[K, V]) Get(k K) V {
	return lm.kv[k].value
}

// Get returns the value in the map stored for key k along with an indicator
// flag. If k is present in the map the associated value is returned along with
// a true flag value. Otherwise if key k is not present, the zero value of type
// V is returned along with a false flag value.
func (lm LinkedMap[K, V]) GetOk(k K) (V, bool) {
	val, ok := lm.kv[k]
	return val.value, ok
}

// SeqKeys returns an [iter.Seq][K] iterator over the keys in the map.
func (lm LinkedMap[K, V]) SeqKeys() iter.Seq[K] {
	return lm.keys.Seq()
}

// Seq2 returns an [iter.Seq2][K,V] iterator over the key-value pairs in the map.
//
// Deprecated: use [Seq]
func (lm LinkedMap[K, V]) Seq2() iter.Seq2[K, V] {
	return lm.Seq()
}

// Seq returns an [iter.Seq2][K,V] iterator over the key-value pairs in the map.
func (lm LinkedMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key := range lm.SeqKeys() {
			if !yield(key, lm.Get(key)) {
				break
			}
		}
	}
}

// Seq returns an [iter.Seq][K,V] iterator over the values in the map.
func (lm LinkedMap[K, V]) SeqValues() iter.Seq[V] {
	return func(yield func(V) bool) {
		for key := range lm.SeqKeys() {
			if !yield(lm.kv[key].value) {
				break
			}
		}
	}
}

// IterKeys returns an [iterator.Iterator][K] over the keys in the map.
func (lm LinkedMap[K, V]) IterKeys() iterator.Iterator[K] {
	return lm.keys.Iter()
}

// Iter returns an [iterator.Iterator2][K,V] over the key-value pairs in the map.
func (lm LinkedMap[K, V]) Iter() iterator.Iterator2[K, V] {
	return iterator.New2(lm.Seq())
}

// IterKeys returns an [iterator.Iterator][V] over the values in the map.
func (lm LinkedMap[K, V]) IterValues() iterator.Iterator[V] {
	return iterator.New(lm.SeqValues())
}

// Delete removes the key from the map and returns the associated value and whether the key was present.
func (lm *LinkedMap[K, V]) Delete(k K) (V, bool) {
	val, ok := lm.kv[k]
	if ok {
		lm.keys.Delete(val.node)
	}
	delete(lm.kv, k)
	return val.value, ok
}

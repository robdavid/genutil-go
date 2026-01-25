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

func makeValue[K comparable, V any](k *list.Node[K], v V) value[K, V] {
	return value[K, V]{k, v}
}

// LinkedMap is map combined with linked list of keys which maintains a consistent
// key order. Keys are placed in the order they are first inserted.
type LinkedMap[K comparable, V any] struct {
	kv   map[K]value[K, V]
	keys list.List[K]
}

// Make creates a new LinkedMap instance.
func Make[K comparable, V any]() LinkedMap[K, V] {
	return LinkedMap[K, V]{
		kv:   make(map[K]value[K, V]),
		keys: list.Make[K](),
	}
}

// Size returns the number of elements in the map.
func (lm LinkedMap[K, V]) Size() int {
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

func (lm LinkedMap[K, V]) GetOk(k K) (V, bool) {
	val, ok := lm.kv[k]
	return val.value, ok
}

func (lm LinkedMap[K, V]) SeqKeys() iter.Seq[K] {
	return lm.keys.Seq()
}

func (lm LinkedMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key := range lm.SeqKeys() {
			yield(key, lm.Get(key))
		}
	}
}

func (lm LinkedMap[K, V]) IterKeys() iterator.Iterator[K] {
	return lm.keys.Iter()
}

func (lm LinkedMap[K, V]) Iter() iterator.Iterator2[K, V] {
	return iterator.New2(lm.Seq())
}

func (lm *LinkedMap[K, V]) Delete(k K) (V, bool) {
	val, ok := lm.kv[k]
	if ok {
		lm.keys.Delete(val.node)
	}
	delete(lm.kv, k)
	return val.value, ok
}

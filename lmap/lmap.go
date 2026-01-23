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

type linkedMap[K comparable, V any] struct {
	kv   map[K]value[K, V]
	keys list.List[K]
}

type LinkedMap[K comparable, V any] = *linkedMap[K, V]

func New[K comparable, V any]() LinkedMap[K, V] {
	lm := &linkedMap[K, V]{
		kv:   make(map[K]value[K, V]),
		keys: list.Make[K](),
	}
	return lm
}

func (lm *linkedMap[K, V]) Size() int {
	return len(lm.kv)
}

func (lm *linkedMap[K, V]) IsEmpty() bool {
	return len(lm.kv) == 0
}

func (lm *linkedMap[K, V]) Put(k K, v V) {
	if current, ok := lm.kv[k]; ok {
		current.value = v
		lm.kv[k] = current
	} else {
		lm.keys.Append(k)
		lm.kv[k] = makeValue(lm.keys.Last(), v)
	}
}

func (lm *linkedMap[K, V]) Get(k K) V {
	return lm.kv[k].value
}

func (lm *linkedMap[K, V]) GetOk(k K) (V, bool) {
	val, ok := lm.kv[k]
	return val.value, ok
}

func (lm *linkedMap[K, V]) SeqKeys() iter.Seq[K] {
	return lm.keys.Seq()
}

func (lm *linkedMap[K, V]) Seq() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key := range lm.SeqKeys() {
			yield(key, lm.Get(key))
		}
	}
}

func (lm *linkedMap[K, V]) IterKeys() iterator.Iterator[K] {
	return lm.keys.Iter()
}

func (lm *linkedMap[K, V]) Iter() iterator.Iterator2[K, V] {
	return iterator.New2(lm.Seq())
}

func (lm *linkedMap[K, V]) Delete(k K) (V, bool) {
	val, ok := lm.kv[k]
	if ok {
		lm.keys.Delete(val.node)
	}
	delete(lm.kv, k)
	return val.value, ok
}

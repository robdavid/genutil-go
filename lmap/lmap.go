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

type LinkedMap[K comparable, V any] struct {
	kv   map[K]value[K, V]
	keys list.List[K]
}

func Make[K comparable, V any]() LinkedMap[K, V] {
	return LinkedMap[K, V]{
		kv:   make(map[K]value[K, V]),
		keys: list.Make[K](),
	}
}

func (lm LinkedMap[K, V]) Size() int {
	return len(lm.kv)
}

func (lm LinkedMap[K, V]) IsEmpty() bool {
	return len(lm.kv) == 0
}

func (lm *LinkedMap[K, V]) Put(k K, v V) {
	if lm.kv == nil {
		lm.kv = make(map[K]value[K, V])
		//lm.keys = list.New[K]()
		lm.keys.Append(k)
		lm.kv[k] = makeValue(lm.keys.Last(), v)
	} else if current, ok := lm.kv[k]; ok {
		current.value = v
		lm.kv[k] = current
	} else {
		lm.keys.Append(k)
		lm.kv[k] = makeValue(lm.keys.Last(), v)
	}
}

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

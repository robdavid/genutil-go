package lmap_test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/lmap"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
)

var stringKeys = []string{"zero", "one", "two", "three", "four", "five"}

func TestEmpty(t *testing.T) {
	lm := lmap.Make[string, int]()
	assert.True(t, lm.IsEmpty())
	assert.Equal(t, 0, lm.Len())
}

func TestZero(t *testing.T) {
	var lm lmap.LinkedMap[string, int]
	assert.True(t, lm.IsEmpty())
	assert.Equal(t, 0, lm.Len())
	z, ok := lm.Delete("one")
	assert.False(t, ok)
	assert.Equal(t, 0, z)
	assert.Empty(t, lm.IterKeys().Collect())
	assert.Empty(t, lm.Iter().Collect2())
	lm.Put("one", 1)
	assert.False(t, lm.IsEmpty())
	assert.Equal(t, 1, lm.Len())
}

func TestPutGet(t *testing.T) {
	lm := lmap.Make[string, int]()
	lm.Put("one", 1)
	val := lm.Get("one")
	assert.Equal(t, 1, val)
	val = lm.Get("two")
	assert.Equal(t, 0, val)
	lm.Put("two", 2)
	val, ok := lm.GetOk("two")
	assert.True(t, ok)
	assert.Equal(t, 2, val)
}

func TestInsertOrder(t *testing.T) {
	lm := lmap.Make[string, int]()
	for i, key := range stringKeys {
		lm.Put(key, i)
	}
	assert.Equal(t, stringKeys, lm.IterKeys().Collect())
	for key, value := range lm.Seq2() {
		assert.Equal(t, stringKeys[value], key)
	}
}

func TestValuesOrder(t *testing.T) {
	lm := lmap.Make[string, int]()
	for i, key := range stringKeys {
		lm.Put(key, i)
	}
	assert.Equal(t, slices.Range(0, len(stringKeys)), lm.IterValues().Collect())
}

func TestInsertDelete(t *testing.T) {
	lm := lmap.Make[string, int]()
	for i, key := range stringKeys {
		lm.Put(key, i)
	}
	lm.Delete("two")
	for key, value := range lm.Seq2() {
		assert.NotEqual(t, 2, value)
		assert.Equal(t, stringKeys[value], key)
	}
	assert.Equal(t, append(stringKeys[:2], stringKeys[3:]...), lm.IterKeys().Collect())
}

func TestFromIterator(t *testing.T) {
	const size = 10
	lm := lmap.From(iterator.Range(size, size*2).Enumerate())
	assert.Equal(t, size, lm.Len())
	i := 0
	for k, v := range lm.Seq2() {
		assert.Equal(t, i, k)
		assert.Equal(t, i+size, v)
		i++
	}
}

func TestFromKeysFunc(t *testing.T) {
	const size = 10
	lm := lmap.FromKeys(slices.Range(0, size), func(x int) string { return fmt.Sprintf("Value %d", x) })
	assert.Equal(t, size, lm.Len())
	for i := range size {
		assert.Equal(t, lm.Get(i), fmt.Sprintf("Value %d", i))
	}
}

func TestReplaceValues(t *testing.T) {
	const size = 10
	lm := lmap.FromIterKeys(iterator.Range(0, size), functions.Id)
	assert.Equal(t, size, lm.Len())
	for i := range size {
		assert.Equal(t, lm.Get(i), i)
	}
	for i := range size {
		lm.Put(i, lm.Get(i)+100)
	}
	for k, v := range lm.Seq2() {
		assert.Equal(t, k+100, v)
	}
}

func Benchmark(b *testing.B) {
	lm := lmap.Make[int, int]()
	for i := range b.N {
		lm.Put(i, i*2)
	}
	b.ResetTimer()
	for i := range b.N {
		if lm.Get(i) != i*2 {
			b.Fatalf("Expected %d, got %d", i*2, lm.Get(i))
		}
	}
}

package lmap_test

import (
	"testing"

	"github.com/robdavid/genutil-go/lmap"
	"github.com/stretchr/testify/assert"
)

var keys = []string{"zero", "one", "two", "three", "four", "five"}

func TestEmpty(t *testing.T) {
	lm := lmap.Make[string, int]()
	assert.True(t, lm.IsEmpty())
	assert.Equal(t, 0, lm.Size())
}

func TestZero(t *testing.T) {
	var lm lmap.LinkedMap[string, int]
	assert.True(t, lm.IsEmpty())
	assert.Equal(t, 0, lm.Size())
	z, ok := lm.Delete("one")
	assert.False(t, ok)
	assert.Equal(t, 0, z)
	assert.Empty(t, lm.IterKeys().Collect())
	assert.Empty(t, lm.Iter().Collect2())
	lm.Put("one", 1)
	assert.False(t, lm.IsEmpty())
	assert.Equal(t, 1, lm.Size())
}

func TestInsertOrder(t *testing.T) {
	lm := lmap.Make[string, int]()
	for i, key := range keys {
		lm.Put(key, i)
	}
	assert.Equal(t, keys, lm.IterKeys().Collect())
	for key, value := range lm.Seq() {
		assert.Equal(t, keys[value], key)
	}
}

func TestInsertDelete(t *testing.T) {
	lm := lmap.Make[string, int]()
	for i, key := range keys {
		lm.Put(key, i)
	}
	lm.Delete("two")
	for key, value := range lm.Seq() {
		assert.NotEqual(t, 2, value)
		assert.Equal(t, keys[value], key)
	}
	assert.Equal(t, append(keys[:2], keys[3:]...), lm.IterKeys().Collect())
}

func TestReplaceValues(t *testing.T) {
	//lm := lmap.LinkedMap[string,int])
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

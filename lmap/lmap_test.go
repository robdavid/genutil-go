package lmap_test

import (
	"testing"

	"github.com/robdavid/genutil-go/lmap"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	lm := lmap.New[string, int]()
	assert.True(t, lm.IsEmpty())
	assert.Equal(t, 0, lm.Size())
}

func TestInsertOrder(t *testing.T) {
	lm := lmap.New[string, int]()
	keys := []string{"zero", "one", "two", "three", "four", "five"}
	for i, key := range keys {
		lm.Put(key, i)
	}
	assert.Equal(t, keys, lm.IterKeys().Collect())
}

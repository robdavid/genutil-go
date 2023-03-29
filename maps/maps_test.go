package maps

import (
	"errors"
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/errors/test"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/slices"
	"github.com/robdavid/genutil-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestInsertPathOne(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b"}, 123, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath([]string{"a", "b"}, m)).Must(t))
}

func TestInsertPathTwo(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b"}, 123, m))
	test.Check(t, PutPath([]string{"a", "c"}, 456, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123, "c": 456},
	}, m)
}

func TestInsertPathFourDeep(t *testing.T) {
	m := make(map[string]any)
	path := []string{"a", "b", "c", "d"}
	test.Check(t, PutPath(path, 123, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": 123,
				},
			},
		},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath(path, m)).Must(t))
	res := test.Result(GetPath(slices.Concat(path, []string{"e"}), m))
	assert.True(t, res.IsError())
	assert.True(t, errors.Is(res.GetErr(), ErrKeyError))
}

func TestInsertPathConflictLeaf(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b", "c", "d"}, 123, m))
	err := PutPath([]string{"a", "b", "c"}, 456, m)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

func TestInsertPathConflictInterior(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b", "c"}, 123, m))
	err := PutPath([]string{"a", "b", "c", "d"}, 456, m)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

func TestKeys(t *testing.T) {
	mymap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	keys := Keys(mymap)
	assert.ElementsMatch(t, []string{"one", "two", "three"}, keys)
}

func TestSortedKeys(t *testing.T) {
	mymap := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}
	keys := SortedKeys(mymap)
	assert.Equal(t, []string{"one", "three", "two"}, keys)
}

func TestSortedValuesByKey(t *testing.T) {
	mymap := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}
	values := SortedValuesByKey(mymap)
	assert.Equal(t, []string{"one", "two", "three"}, values)
}

func TestEmptyKeys(t *testing.T) {
	mymap := map[string]int{}
	keys := Keys(mymap)
	assert.Empty(t, keys)
}

func TestValues(t *testing.T) {
	mymap := map[string]int{
		"one":   1,
		"three": 3,
		"two":   2,
	}
	values := Values(mymap)
	assert.ElementsMatch(t, []int{1, 2, 3}, values)
}

func TestIterKeys(t *testing.T) {
	mymap := map[string]int{
		"one":   1,
		"three": 3,
		"two":   2,
	}
	keys := iterator.Collect(IterKeys(mymap))
	assert.ElementsMatch(t, []string{"one", "two", "three"}, keys)
}

func BenchmarkKeyIterator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mymap := make(map[string]int)
		for j := 0; j < 100; j++ {
			mymap[fmt.Sprintf("key-%d", j)] = j
		}
		keys := iterator.Collect(IterKeys(mymap))
		assert.Equal(b, 100, len(keys))
	}
}

func BenchmarkKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mymap := make(map[string]int)
		for j := 0; j < 100; j++ {
			mymap[fmt.Sprintf("key-%d", j)] = j
		}
		keys := Keys(mymap)
		assert.Equal(b, 100, len(keys))
	}
}

func TestSortedItems(t *testing.T) {
	mymap := make(map[string]int)
	for j := 0; j < 5; j++ {
		mymap[fmt.Sprintf("key-%d", j)] = j
	}
	items := SortedItems(mymap)
	expected := []tuple.Tuple2[string, int]{
		tuple.Of2("key-0", 0),
		tuple.Of2("key-1", 1),
		tuple.Of2("key-2", 2),
		tuple.Of2("key-3", 3),
		tuple.Of2("key-4", 4),
	}
	assert.Equal(t, expected, items)
}

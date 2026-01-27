package maps

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"testing"

	"github.com/robdavid/genutil-go/errors/test"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/slices"
	"github.com/robdavid/genutil-go/tuple"
	"github.com/stretchr/testify/assert"
)

func TestGetAs(t *testing.T) {
	m := map[string]any{
		"one":   1,
		"two":   2,
		"hello": "world",
	}
	assert := assert.New(t)
	assert.Equal(1, GetAs[int](m, "one").GetOr(-1))
	assert.Equal(2, GetAs[int](m, "two").GetOr(-1))
	assert.Equal(-1, GetAs[int](m, "three").GetOr(-1))
	assert.Equal("nope", GetAs[string](m, "two").GetOr("nope"))
	assert.Equal("world", GetAs[string](m, "hello").GetOr("nope"))
}

func TestFetchPathNil(t *testing.T) {
	var m map[string]any = nil
	_, e := GetPath(m, []string{"a", "b", "c", "e"})
	assert.ErrorIs(t, e, ErrKeyError)
}

func TestFetchPath(t *testing.T) {
	m := map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": 123,
				},
			},
		},
	}
	v, e := GetPath(m, []string{"a", "b", "c", "d"})
	assert.NoError(t, e)
	assert.Equal(t, 123, v)
	_, e = GetPath(m, []string{"a", "b", "c", "e"})
	assert.ErrorIs(t, e, ErrKeyError)
	_, e = GetPath(m, []string{"a", "b", "c", "d", "e"})
	assert.ErrorIs(t, e, ErrKeyError)
}

func TestInsertPathOne(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b"}, 123))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath(m, []string{"a", "b"})).Must(t))
}

func TestInsertPathTwo(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b"}, 123))
	test.Check(t, PutPath(m, []string{"a", "c"}, 456))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123, "c": 456},
	}, m)
}

func TestInsertPathFourDeep(t *testing.T) {
	m := make(map[string]any)
	path := []string{"a", "b", "c", "d"}
	test.Check(t, PutPath(m, path, 123))
	assert.Equal(t, map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": 123,
				},
			},
		},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath(m, path)).Must(t))
	res := test.Result(GetPath(m, slices.Concat(path, []string{"e"})))
	assert.True(t, res.IsError())
	assert.True(t, errors.Is(res.GetErr(), ErrKeyError))
}

func TestInsertPathConflictLeaf(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "c", "d"}, 123))
	err := PutPath(m, []string{"a", "b", "c"}, 456)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

func TestInsertPathConflictInterior(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "c"}, 123))
	err := PutPath(m, []string{"a", "b", "c", "d"}, 456)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

func TestDeleteTop(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a"}, 123))
	value, ok := test.Result2(DeletePath(m, []string{"a"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, 123, value)
	assert.True(t, len(m) == 0)
}

func TestDeleteDepth2(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b"}, 123))
	value, ok := test.Result2(DeletePath(m, []string{"a", "b"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, 123, value)
	assert.True(t, len(m) == 0)
}

func TestDeleteDepth3(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "c"}, 123))
	value, ok := test.Result2(DeletePath(m, []string{"a", "b", "c"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, 123, value)
	assert.True(t, len(m) == 0)
}

func TestDeleteEmptySubtree(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "x"}, 123))
	test.Check(t, PutPath(m, []string{"y"}, 456))
	value, ok := test.Result2(DeletePath(m, []string{"a", "b", "x"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, 123, value)
	expected := map[string]any{"y": 456}
	assert.Equal(t, expected, m)
}

func TestDeleteLeaf(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "x"}, 123))
	test.Check(t, PutPath(m, []string{"a", "b", "y"}, 456))
	value, ok := test.Result2(DeletePath(m, []string{"a", "b", "x"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, 123, value)
	expected := map[string]any{"a": map[string]any{"b": map[string]any{"y": 456}}}
	assert.Equal(t, expected, m)
}

func TestDeleteSubtree(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath(m, []string{"a", "b", "x"}, 123))
	test.Check(t, PutPath(m, []string{"a", "b", "y"}, 456))
	test.Check(t, PutPath(m, []string{"a", "z"}, 789))
	value, ok := test.Result2(DeletePath(m, []string{"a", "b"})).Must2(t)
	assert.True(t, ok)
	assert.Equal(t, map[string]any{"x": 123, "y": 456}, value)
	expected := map[string]any{"a": map[string]any{"z": 789}}
	assert.Equal(t, expected, m)
}

func TestDeleteExample(t *testing.T) {
	m := map[string]any{
		"one": 1,
		"two": map[string]any{
			"three": 23,
		},
	}
	DeletePath(m, []string{"two", "three"})
	assert.Equal(t, map[string]any{"one": 1}, m)
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

func TestIter(t *testing.T) {
	mymap := map[string]int{
		"one":   1,
		"three": 3,
		"two":   2,
	}
	items := Iter(mymap).Collect2()
	assert.ElementsMatch(t, []iterator.KeyValue[string, int]{iterator.KVOf("one", 1), iterator.KVOf("two", 2), iterator.KVOf("three", 3)}, items)
}

func TestIterSize(t *testing.T) {
	mymap := make(map[int]string)
	const mapsize = 20
	for n := range mapsize {
		mymap[n] = strconv.Itoa(n)
	}
	itr := Iter(mymap)
	count := mapsize
	for range itr.Seq2() {
		count--
		assert.Equal(t, count, itr.Size().Size)
	}
	assert.Zero(t, itr.Size().Size)
}

func TestIterSimple(t *testing.T) {
	mymap := make(map[int]int)
	seen := make(map[int]bool)
	const mapsize = 20
	for n := range mapsize {
		mymap[n] = n
	}
	itr := Iter(mymap)
	count := mapsize
	for itr.Next() {
		count--
		assert.Equal(t, count, itr.Size().Size)
		seen[itr.Key()] = true
		if count <= 5 {
			break
		}
	}
	assert.Equal(t, 5, itr.Size().Size)
	remain := itr.Collect()
	assert.Equal(t, 5, len(remain))
	for _, r := range remain {
		assert.False(t, seen[r])
	}
}

func generateMap(size int) map[string]int {
	mymap := make(map[string]int)
	for j := 0; j < size; j++ {
		mymap[fmt.Sprintf("key-%d", j)] = j
	}
	return mymap
}

type testMapValue struct {
	intValue     int
	floatValue   float64
	strValue     string
	complexValue complex128
}

func generateLargeMap(size int) map[int]testMapValue {
	mymap := make(map[int]testMapValue)
	for j := 0; j < size; j++ {
		f := float64(j)
		mymap[j] = testMapValue{j, f, strconv.Itoa(j), complex(f, f*2)}
	}
	return mymap
}

func TestFindUsing(t *testing.T) {
	mymap := generateMap(10)
	match := FindUsing(mymap, func(k string, v int) bool { return v == 5 })
	assert.Equal(t, match.Get(), tuple.Of2("key-5", 5))
}

func TestFindUsingRef(t *testing.T) {
	mymap := generateMap(10)
	match := FindUsingRef(mymap, func(k *string, v *int) bool { return *v == 5 })
	assert.Equal(t, *match.Get().First, "key-5")
	assert.Equal(t, *match.Get().Second, 5)
}

func BenchmarkKeyIterator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mymap := generateMap(100)
		keys := iterator.Collect(IterKeys(mymap))
		assert.Equal(b, 100, len(keys))
	}
}

func BenchmarkKeys(b *testing.B) {
	for i := 0; i < b.N; i++ {
		mymap := generateMap(100)
		keys := Keys(mymap)
		assert.Equal(b, 100, len(keys))
	}
}

func BenchmarkFindUsing(b *testing.B) {
	mymap := generateLargeMap(100)
	for i := 0; i < b.N; i++ {
		match := FindUsing(mymap, func(k int, v testMapValue) bool { return v.intValue == 50 })
		assert.Equal(b, match.Get().First, 50)
	}
}

func BenchmarkFindUsingRef(b *testing.B) {
	mymap := generateLargeMap(100)
	for i := 0; i < b.N; i++ {
		match := FindUsingRef(mymap, func(k *int, v *testMapValue) bool { return v.intValue == 50 })
		assert.Equal(b, *match.Get().First, 50)
	}
}

func TestSortedItems(t *testing.T) {
	mymap := generateMap(5)
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

func TestIteratorMutations(t *testing.T) {
	m := make(map[int]int)
	for i := range 10 {
		m[i] = i * 2
	}
	iter := IterMut(m)
	assert.Equal(t, len(m), iter.Size().Size)
	for k, v := range iter.Seq2() {
		if k%3 == 0 {
			iter.Set(v / 2 * 3)
		} else if k%2 == 0 {
			iter.Delete()
		}
	}
	keys := make([]int, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	actual := make([]int, 0, len(keys))
	for _, k := range keys {
		actual = append(actual, m[k])
	}
	expected := []int{0, 2, 9, 10, 18, 14, 27}
	assert.Equal(t, expected, actual)
}

func TestIterMutNextCollect(t *testing.T) {
	m := make(map[int]int)
	for i := range 10 {
		m[i] = i
	}
	itr := IterMut(m)
	var collected []int
	count := 0
	assert.True(t, itr.SeqOK())
	for itr.Next() {
		assert.False(t, itr.SeqOK())
		count++
		if count == 5 {
			collected = itr.Collect()
		}
	}
	assert.Equal(t, 5, len(collected))
}

func TestClone(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	c := Clone(m)
	assert.Equal(t, m, c)
	c["a"] = 9
	assert.Equal(t, 1, m["a"])
}

func TestSubAndSubI(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	s := map[string]int{"b": 2}
	SubI(m, s)
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, m)

	m2 := map[string]int{"a": 1, "b": 2, "c": 3}
	r := Sub(m2, s)
	assert.Equal(t, map[string]int{"a": 1, "c": 3}, r)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m2)
}

func TestSubI_LargerS(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	s := map[string]int{"b": 9, "c": 3, "d": 4}
	SubI(m, s)
	assert.Equal(t, map[string]int{"a": 1}, m)
}

func TestSub_LargerS(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	s := map[string]int{"b": 9, "c": 3, "d": 4}
	r := Sub(m, s)
	assert.Equal(t, map[string]int{"a": 1}, r)
	// original map should be unchanged
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, m)
}

func TestAddAndAddI(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2}
	a := map[string]int{"b": 9, "c": 3}
	AddI(m, a)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, m)

	m2 := map[string]int{"a": 1, "b": 2}
	r := Add(m2, a)
	assert.Equal(t, map[string]int{"a": 1, "b": 2, "c": 3}, r)
	assert.Equal(t, map[string]int{"a": 1, "b": 2}, m2)
}

func TestAsFunc(t *testing.T) {
	m := map[int]string{1: "one", 2: "two"}
	f := AsFunc(m)
	assert.Equal(t, "one", f(1))
	assert.Equal(t, "", f(3))
}

package list_test

import (
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/list"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	var empty list.List[int]
	assert.Equal(t, 0, empty.Size())
	assert.True(t, empty.IsEmpty())
	assert.Empty(t, empty.Iter().Collect())
	assert.Empty(t, empty.RevIter().Collect())
	assert.Empty(t, empty.IterNode().Collect())
	assert.Empty(t, empty.RevIterNode().Collect())
	assert.True(t, empty.GetFirst().ToRef().IsEmpty())
	assert.True(t, empty.GetLast().ToRef().IsEmpty())
}

func TestOf(t *testing.T) {
	var lst list.List[int]
	lst = list.Of(1, 2, 3, 4)
	expected := []int{1, 2, 3, 4}
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestAppend(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	for n := range iterator.Range(0, size).Seq() {
		list.Append(n)
	}
	assert.Equal(t, size, list.Size())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestAppendSlice(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	slice := slices.Range(0, size)
	list.Append(slice...)
	assert.Equal(t, size, list.Size())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestPrepend(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	for n := range iterator.IncRange(size-1, 0).Seq() {
		list.Prepend(n)
	}
	assert.Equal(t, size, list.Size())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestPrependSlice(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	slice := slices.Range(0, size)
	list.Prepend(slice...)
	assert.Equal(t, size, list.Size())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsert(t *testing.T) {
	list := list.Of(1, 2, 4, 5)
	list.Insert(list.At(2), 3)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertStart(t *testing.T) {
	list := list.Of(1, 2, 3, 4)
	list.Insert(list.At(0), 0)
	expected := []int{0, 1, 2, 3, 4}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertEnd(t *testing.T) {
	list := list.Of(1, 2, 3, 5)
	list.Insert(list.At(-1), 4)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertAfter(t *testing.T) {
	list := list.Of(1, 2, 3, 5)
	list.InsertAfter(list.At(2), 4)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertNothingAfter(t *testing.T) {
	list := list.Of(1, 2, 3, 4, 5)
	list.InsertAfter(list.At(2))
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertSliceAfter(t *testing.T) {
	list := list.Of(1, 5)
	list.InsertAfter(list.First(), 2, 3, 4)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertAfterStart(t *testing.T) {
	list := list.Of(1, 3, 4, 5)
	list.InsertAfter(list.First(), 2)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestInsertAfterEnd(t *testing.T) {
	list := list.Of(1, 2, 3, 4)
	list.InsertAfter(list.Last(), 5)
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestFrom(t *testing.T) {
	const size = 10
	seq := iterator.Range(0, size)
	lst := list.From(seq)
	expected := slices.Range(0, size)
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestFromSeq(t *testing.T) {
	const size = 10
	seq := iterator.Range(0, size).Seq()
	lst := list.FromSeq(seq)
	expected := slices.Range(0, size)
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestSizeAndCount(t *testing.T) {
	const insSize = 10
	const iterations = 20
	const deletions = 50
	var lst list.List[int]
	insSlice := slices.Range(0, insSize)
	lst.Append(insSlice...)
	for range iterations {
		lst.Insert(lst.At(lst.Size()/2), insSlice...)
		lst.InsertAfter(lst.At(lst.Size()/2), insSlice...)
	}
	for range deletions {
		lst.Delete(lst.At(lst.Size() / 2))
	}
	expectedSize := insSize*(iterations*2+1) - deletions
	assert.Equal(t, expectedSize, lst.Size())
	assert.Equal(t, expectedSize, lst.First().Count())
	assert.Equal(t, expectedSize, lst.Last().RevCount())
}

func TestFirstLast(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	assert.Equal(t, 0, lst.GetFirst().Get())
	assert.Equal(t, size-1, lst.GetLast().Get())
}

func TestGet(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	for i := range size {
		assert.Equal(t, i, lst.Get(i))
	}
	for i := 1; i <= size; i++ {
		assert.Equal(t, size-i, lst.Get(-i))
	}
}

func TestIter(t *testing.T) {
	const size = 5
	lst := list.From(iterator.Range(0, size))
	iter := lst.Iter()

	// Test size
	assert.True(t, iter.Size().IsKnown())
	assert.Equal(t, size, iter.Size().Size)

	// Test manual iteration
	collected := iter.Collect()
	assert.Equal(t, slices.Range(0, size), collected)
	assert.Equal(t, size, cap(collected))

	// Test Collect after manual iteration (should be empty since iterator is consumed)
	collected = iter.Collect()
	assert.Empty(t, collected)

	iter.Reset()

	// Test fresh iterator Collect
	collected = iter.Collect()
	assert.Equal(t, slices.Range(0, size), collected)
}

func TestIterNode(t *testing.T) {
	const size = 5
	lst := list.From(iterator.Range(0, size))
	iter := lst.IterNode()

	// Test size
	assert.True(t, iter.Size().IsKnown())
	assert.Equal(t, size, iter.Size().Size)

	// Test value collection
	collected := iterator.Map(iter, list.NodeToValue).Collect()
	assert.Equal(t, slices.Range(0, size), collected)
	assert.Equal(t, size, cap(collected))

	// Test Collect after manual iteration (should be empty since iterator is consumed)
	collectedNodes := iter.Collect()
	assert.Empty(t, collectedNodes)

	iter.Reset()

	// Test fresh iterator Collect
	collected = iterator.Map(iter, list.NodeToValue).Collect()
	assert.Equal(t, slices.Range(0, size), collected)
}

func TestListSeq(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	c := make([]int, 0, size)
	for v := range lst.Seq() {
		c = append(c, v)
	}
	assert.Equal(t, slices.Range(0, size), c)
}

func TestListRevSeq(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	c := make([]int, 0, size)
	for v := range lst.RevSeq() {
		c = append(c, v)
	}
	assert.Equal(t, slices.Reverse(slices.Range(0, size)), c)
}

func TestNodeFirstLast(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	mid := lst.At(size / 2)
	firstNode := mid.First()
	lastNode := mid.Last()
	assert.Equal(t, 0, firstNode.Get())
	assert.Equal(t, size-1, lastNode.Get())
}

func TestIterRef(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	for ln := range lst.SeqNode() {
		*ln.Ref() += 5
	}
	c := lst.Iter().Collect()
	assert.Equal(t, slices.Range(5, size+5), c)
}

func TestRevIterRef(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	for ln := range lst.RevSeqNode() {
		*ln.Ref() += 5
	}
	c := lst.Iter().Collect()
	assert.Equal(t, slices.Range(5, size+5), c)
}

func TestIterSet(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	for ln := range lst.SeqNode() {
		ln.Set(ln.Get() + 5)
	}
	c := lst.Iter().Collect()
	assert.Equal(t, slices.Range(5, size+5), c)
}

func TestDelete(t *testing.T) {
	const size = 10
	input := slices.Range(0, size)
	lst := list.Of(input...)
	for node := range lst.SeqNode() {
		if node.Get()%2 != 0 {
			lst.Delete(node)
		}
	}
	expected := slices.RangeBy(0, size, 2)
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestDeleteSingleton(t *testing.T) {
	lst := list.Of(1)
	lst.Delete(lst.First())
	assert.True(t, lst.IsEmpty())
	assert.Equal(t, 0, lst.Size())
}

func TestDeleteFirst(t *testing.T) {
	const size = 10
	lst := list.Of(slices.Range(0, size)...)
	lst.Delete(lst.First())
	expected := slices.Range(1, size)
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestDeleteLast(t *testing.T) {
	const size = 10
	lst := list.Of(slices.Range(0, size)...)
	lst.Delete(lst.Last())
	expected := slices.Range(0, size-1)
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func BenchmarkIter(b *testing.B) {
	var lst list.List[int]
	for i := range b.N {
		lst.Append(i)
	}
	b.ResetTimer()
	n := 0
	for i := range lst.Iter().Seq() {
		n += i
	}
	if n != (b.N*(b.N-1))/2 {
		b.Fail()
	}
}

func BenchmarkIterCollect(b *testing.B) {
	var lst list.List[int]
	for i := range b.N {
		lst.Append(i)
	}
	b.ResetTimer()
	c := lst.Iter().Collect()
	if len(c) != b.N || c[0] != 0 || c[b.N-1] != b.N-1 {
		b.Fail()
	}
}

func BenchmarkSeq(b *testing.B) {
	var lst list.List[int]
	for i := range b.N {
		lst.Append(i)
	}
	b.ResetTimer()
	n := 0
	for i := range lst.Seq() {
		n += i
	}
	if n != (b.N*(b.N-1))/2 {
		b.Fail()
	}
}

func BenchmarkSeqCollect(b *testing.B) {
	var lst list.List[int]
	for i := range b.N {
		lst.Append(i)
	}
	b.ResetTimer()
	var c []int
	for v := range lst.Seq() {
		c = append(c, v)
	}
	if len(c) != b.N || c[0] != 0 || c[b.N-1] != b.N-1 {
		b.Fail()
	}
}

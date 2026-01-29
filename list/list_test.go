package list_test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/list"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
)

func TestNil(t *testing.T) {
	var zero list.List[int]
	assert.True(t, zero.IsNil())
	assert.True(t, zero.IsEmpty())
	assert.Equal(t, 0, zero.Len())
	assert.Empty(t, zero.Iter().Collect())
	assert.Empty(t, zero.RevIter().Collect())
	assert.Empty(t, zero.IterNode().Collect())
	assert.Empty(t, zero.RevIterNode().Collect())
	assert.True(t, zero.GetFirst().ToRef().IsEmpty())
	assert.True(t, zero.GetLast().ToRef().IsEmpty())
}

func TestOf(t *testing.T) {
	var lst list.List[int]
	lst = list.Of(1, 2, 3, 4)
	expected := []int{1, 2, 3, 4}
	assert.Equal(t, expected, lst.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), lst.RevIter().Collect())
}

func TestNew(t *testing.T) {
	lst := list.New[int]()
	assert.NotNil(t, lst)
	assert.True(t, lst.IsEmpty())
	assert.Equal(t, 0, lst.Len())
}

func TestClear(t *testing.T) {
	var lst list.List[int]
	lst = list.Of(1, 2, 3, 4)
	lst.Clear()
	assert.True(t, lst.IsEmpty())
	assert.Equal(t, 0, lst.Len())
}

func TestString(t *testing.T) {
	lst := list.Of(1, 2, 3, 4)
	assert.Equal(t, fmt.Sprintf("%v", lst.Iter().Collect()), fmt.Sprintf("%v", lst))
}

func TestAppend(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	for n := range iterator.Range(0, size).Seq() {
		list.Append(n)
	}
	assert.Equal(t, size, list.Len())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestAppendNil(t *testing.T) {
	var lst list.List[int]
	assert.PanicsWithError(t, list.ErrNilList.Error(), func() {
		lst.Append(1)
	})
	lst = lst.Make()
	lst.Append(1)
	assert.Equal(t, []int{1}, lst.Iter().Collect())
}

func TestAppendSlice(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	slice := slices.Range(0, size)
	list.Append(slice...)
	assert.Equal(t, size, list.Len())
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
	assert.Equal(t, size, list.Len())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestPrependNil(t *testing.T) {
	assert.PanicsWithError(t, list.ErrNilList.Error(), func() {
		var lst list.List[int]
		lst.Prepend(1)
	})
}

func TestPrependSlice(t *testing.T) {
	const size = 10
	list := list.Make[int]()
	slice := slices.Range(0, size)
	list.Prepend(slice...)
	assert.Equal(t, size, list.Len())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestCopyClone(t *testing.T) {
	l1 := list.New[string]()
	l2 := l1
	l2.Append("first")
	assert.Equal(t, l2.Get(0), "first")
	assert.Equal(t, l1.Get(0), "first")
	l3 := l1.Clone()
	l1.Append("last")
	l3.Append("second")
	assert.Equal(t, l3.Get(1), "second")
	assert.Equal(t, l1.Get(1), "last")
}

func TestInsertNothing(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	lst.Insert(lst.First())
	assert.Equal(t, slices.Range(0, size), lst.Iter().Collect())
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

func TestLenAndCount(t *testing.T) {
	const insSize = 10
	const iterations = 20
	const deletions = 50
	lst := list.New[int]()
	insSlice := slices.Range(0, insSize)
	lst.Append(insSlice...)
	for range iterations {
		lst.Insert(lst.At(lst.Len()/2), insSlice...)
		lst.InsertAfter(lst.At(lst.Len()/2), insSlice...)
	}
	for range deletions {
		lst.Delete(lst.At(lst.Len() / 2))
	}
	expectedLen := insSize*(iterations*2+1) - deletions
	assert.Equal(t, expectedLen, lst.Len())
	assert.Equal(t, expectedLen, lst.First().Count())
	assert.Equal(t, expectedLen, lst.Last().RevCount())
}

func TestFirstLast(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	assert.Equal(t, 0, lst.GetFirst().Get())
	assert.Equal(t, size-1, lst.GetLast().Get())
}

func TestGetSetAt(t *testing.T) {
	const size = 10
	lst := list.New[int]()
	for range size {
		lst.Append(0)
	}
	for i := range size {
		lst.Set(i, i+10)
	}
	for i := range size {
		assert.Equal(t, i+10, lst.Get(i))
	}
}

func TestListAtPanic(t *testing.T) {
	lst := list.Of(1, 2, 3)
	assert.PanicsWithValue(t, list.ErrIndexError, func() { lst.At(10) })
	assert.PanicsWithValue(t, list.ErrIndexError, func() { lst.At(-10) })
}

func TestInsertPanic(t *testing.T) {
	lst := list.Of[int]()
	assert.PanicsWithValue(t, list.ErrNilNode, func() { lst.InsertAfter(lst.First(), 0) })
	assert.PanicsWithValue(t, list.ErrNilNode, func() { lst.Insert(lst.Last(), 0) })
}

func TestWalkNodes(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	forward := make([]int, 0, size)
	for n := lst.First(); n != nil; n = n.Next() {
		forward = append(forward, n.Get())
	}
	assert.Equal(t, slices.Range(0, size), forward)
	backward := make([]int, 0, size)
	var mid *list.Node[int]
	count := 0
	for n := lst.Last(); n != nil; n = n.Prev() {
		backward = append(backward, n.Get())
		if count == size/2 {
			mid = n
		}
		count++
	}
	assert.Equal(t, slices.Reverse(slices.Range(0, size)), backward)
	assert.Equal(t, lst.First().Get(), mid.First().Get())
	assert.Equal(t, lst.Last().Get(), mid.Last().Get())
}

func TestNodeIter(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	assert.Equal(t, slices.Range(0, size), lst.First().Iter().Collect())
	assert.Equal(t, slices.IncRange(size-1, 0), lst.Last().RevIter().Collect())
	assert.Equal(t, slices.Range(0, size), iterator.Map(lst.First().IterNode(), list.NodeToValue).Collect())
	assert.Equal(t, slices.IncRange(size-1, 0), iterator.Map(lst.Last().RevIterNode(), list.NodeToValue).Collect())

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

func TestNodeAtPanic(t *testing.T) {
	lst := list.Of(1, 2, 3)
	node := lst.First()
	assert.PanicsWithValue(t, list.ErrIndexError, func() { node.At(10) })
	assert.PanicsWithValue(t, list.ErrIndexError, func() { node.At(-10) })
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

func TestNodeRefNil(t *testing.T) {
	var node *list.Node[int]
	assert.Nil(t, node.Ref())
	assert.False(t, node.Set(5))
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

func TestDeleteNil(t *testing.T) {
	lst := list.Of(1, 2, 3)
	original := lst.Iter().Collect()
	lst.Delete(nil)
	assert.Equal(t, original, lst.Iter().Collect())
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

func BenchmarkAppend(b *testing.B) {
	lst := list.Make[int]()
	for i := range b.N {
		lst.Append(i)
	}
	if lst.Len() != b.N {
		b.Errorf("Expecting length %d, got %d", b.N, lst.Len())
	}
}

func BenchmarkIter(b *testing.B) {
	lst := list.Make[int]()
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
	lst := list.New[int]()
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
	lst := list.New[int]()
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
	lst := list.Make[int]()
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

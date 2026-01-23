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
	list := list.New[int]()
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
	list := list.New[int]()
	slice := slices.Range(0, size)
	list.Append(slice...)
	assert.Equal(t, size, list.Size())
	expected := slices.Range(0, size)
	assert.Equal(t, expected, list.Iter().Collect())
	assert.Equal(t, slices.Reverse(expected), list.RevIter().Collect())
}

func TestPrepend(t *testing.T) {
	const size = 10
	list := list.New[int]()
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
	list := list.New[int]()
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

func TestSize(t *testing.T) {
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

/*

func TestInsertSlice(t *testing.T) {
	const size = 10
	list := list.New[int]()
	slice := slices.Range(0, size)
	list.Insert(slice...)
	assert.Equal(t, size, list.Size())
	l := list
	for i := range size {
		n, ok := l.Get().GetOK()
		require.True(t, ok)
		assert.Equal(t, i, n)
		l = l.Next()
		require.Equal(t, i == size-1, l.IsEmpty())
	}
	collected := list.Iter().Collect()
	assert.Equal(t, slice, collected)
}

func TestLastAndFirst(t *testing.T) {
	const size = 10
	list := list.Of(slices.Range(1, size+1)...)
	last := list.Last()
	assert.Equal(t, option.Value(size), last.Get())
	assert.Equal(t, size, last.RevSize())
	assert.Equal(t, option.Value(1), last.First().Get())
}

func TestAt(t *testing.T) {
	const size = 10
	const mid = size / 2
	defer test.ReportErr(t)
	l := list.Of(slices.Range(0, size)...)
	m := l.At(mid)
	assert.Equal(t, 0, l.At(0).Get().Try())
	assert.Equal(t, mid, m.Get().Try())
	assert.Equal(t, 0, m.At(-mid).Get().Try())
	assert.Equal(t, mid*2-1, m.At(mid-1).Get().Try())
}

func TestGetAt(t *testing.T) {
	const size = 10
	const mid = size / 2
	defer test.ReportErr(t)
	l := list.Of(slices.Range(0, size)...)
	m := l.At(mid)
	assert.Equal(t, 0, l.GetAt(0).Try())
	assert.Equal(t, mid, l.GetAt(mid).Try())
	assert.Equal(t, mid, m.GetAt(0).Try())
	assert.Equal(t, 0, m.GetAt(-mid).Try())
	assert.Equal(t, mid*2-1, m.GetAt(mid-1).Try())
}

func TestIter(t *testing.T) {
	const size = 10
	list := list.New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	collected := list.Iter().Collect()
	assert.Equal(t, slice, collected)
	assert.Equal(t, size, cap(collected))
}

func TestIterRef(t *testing.T) {
	const size = 10
	defer test.ReportErr(t)
	lst := list.From(iterator.Range(0, size))
	for lp := range lst.IterRef().Seq() {
		*lp += 5
	}
	c := lst.Iter().Collect()
	assert.Equal(t, slices.Range(5, size+5), c)
}

func TestRevIter(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	last := lst.Last()
	reved := last.RevIter().Collect()
	assert.Equal(t, slices.IncRange(size-1, 0), reved)
}

func TestRevIterRef(t *testing.T) {
	const size = 10
	lst := list.From(iterator.Range(0, size))
	last := lst.Last()
	for lp := range last.RevIterRef().Seq() {
		*lp += 5
	}
	reved := last.RevIter().Collect()
	assert.Equal(t, slices.IncRange(size+4, 5), reved)
}

func TestRef(t *testing.T) {
	const size = 10
	list := list.From(iterator.Range(0, size))
	for itr := range list.SeqList() {
		*itr.Ref()++
	}
	c := list.Iter().Collect()
	assert.Equal(t, slices.Range(1, size+1), c)
}

func TestIterList(t *testing.T) {
	defer test.ReportErr(t)
	const size = 10
	list := list.New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	for i, l := range list.IterList().Enumerate().Seq2() {
		if i > 0 && i < size-1 {
			v := l.Get().Try()
			p := l.Prev().Get().Try()
			n := l.Next().Get().Try()
			assert.Equal(t, i, v)
			assert.Equal(t, i-1, p)
			assert.Equal(t, i+1, n)
		}
	}
	collected := list.Iter().Collect()
	assert.Equal(t, slice, collected)
	assert.Equal(t, size, cap(collected))
}

func TestRevIterList(t *testing.T) {
	defer test.ReportErr(t)
	const size = 10
	list := list.New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	for i, l := range list.Last().RevIterList().Enumerate().Seq2() {
		if i > 0 && i < size-1 {
			v := l.Get().Try()
			p := l.Prev().Get().Try()
			n := l.Next().Get().Try()
			assert.Equal(t, size-i-1, v)
			assert.Equal(t, size-i-2, p)
			assert.Equal(t, size-i, n)
		}
	}
	revSlice := slices.IncRange(size-1, 0)
	collected := list.Last().RevIter().Collect()
	assert.Equal(t, revSlice, collected)
	assert.Equal(t, size, cap(collected))
}

func TestListSet(t *testing.T) {
	defer test.ReportErr(t)
	const size = 10
	list := list.New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	for l := range list.SeqList() {
		l.Set(l.Get().Try() * 2)
	}
	collected := list.Iter().Collect()
	for i := range size {
		assert.Equal(t, slice[i]*2, collected[i])
	}
	assert.Equal(t, size, len(collected))
	assert.Equal(t, size, cap(collected))
}

func TestFrom(t *testing.T) {
	const size = 10
	input := slices.Range(0, size)
	lst := list.Of(input...)
	assert.Equal(t, input, lst.Iter().Collect())
}

*/

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

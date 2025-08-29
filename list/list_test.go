package list

import (
	"testing"

	"github.com/robdavid/genutil-go/errors/test"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/option"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	empty := New[int]()
	assert.Equal(t, 0, empty.Size())
	assert.True(t, empty.Get().ToRef().IsEmpty())
}

func TestAppend(t *testing.T) {
	const size = 10
	list := New[int]()
	for n := range iterator.Range(0, 10).Seq() {
		list.Append(n)
	}
	assert.Equal(t, size, list.Size())
	l := list
	for i := range size {
		var ok bool
		assert.Equal(t, option.Value(i), l.Get())
		l, ok = l.Next().GetOK()
		require.Equal(t, i != size-1, ok)
		if ok {
			require.NotNil(t, l.head)
		}
	}
}

func TestInsert(t *testing.T) {
	const size = 10
	list := New[int]()
	for n := range iterator.Range(0, 10).Seq() {
		list.Insert(n)
	}
	assert.Equal(t, size, list.Size())
	l := list
	for i := range size {
		n, ok := l.Get().GetOK()
		require.True(t, ok)
		assert.Equal(t, size-1-i, n)
		l, ok = l.Next().GetOK()
		require.Equal(t, i != size-1, ok)
		if ok {
			require.NotNil(t, l.head)
		}
	}
}

func TestIterator(t *testing.T) {
	const size = 10
	list := New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	collected := list.Iter().Collect()
	assert.Equal(t, slice, collected)
	assert.Equal(t, size, cap(collected))
}

func TestListIterator(t *testing.T) {
	defer test.ReportErr(t)
	const size = 10
	list := New[int]()
	slice := slices.Range(0, size)
	for _, n := range slice {
		list.Append(n)
	}
	for i, l := range list.IterList().Enumerate().Seq2() {
		if i > 0 && i < size-1 {
			v := l.Get().Try()
			p := l.Prev().Try().Get().Try()
			n := l.Next().Try().Get().Try()
			assert.Equal(t, i, v)
			assert.Equal(t, i-1, p)
			assert.Equal(t, i+1, n)
		}
	}
	collected := list.Iter().Collect()
	assert.Equal(t, slice, collected)
	assert.Equal(t, size, cap(collected))
}

func TestListSet(t *testing.T) {
	defer test.ReportErr(t)
	const size = 10
	list := New[int]()
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
	lst := From(input...)
	assert.Equal(t, input, lst.Iter().Collect())
}

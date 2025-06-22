package iterator_test

import (
	"fmt"
	"iter"
	"sort"
	"strings"
	"testing"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapMutations(t *testing.T) {
	m := make(map[int]int)
	for i := range 10 {
		m[i] = i * 2
	}
	for k := range m {
		if k%3 == 0 {
			m[k] = k * 3
		} else if k%2 == 0 {
			delete(m, k)
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
func TestCollectInto(t *testing.T) {
	iter := iterator.Range(0, 5)
	var output []int = nil
	iterator.CollectInto(iter, &output)
	iter2 := iterator.Range(10, 15)
	iterator.CollectInto(iter2, &output)
	expected := []int{0, 1, 2, 3, 4, 10, 11, 12, 13, 14}
	assert.Equal(t, expected, output)
}

func TestInclusiveCollectInto(t *testing.T) {
	iter := iterator.IncRange(0, 5)
	var output []int = nil
	iterator.CollectInto(iter, &output)
	iter2 := iterator.IncRange(10, 15)
	iterator.CollectInto(iter2, &output)
	expected := []int{0, 1, 2, 3, 4, 5, 10, 11, 12, 13, 14, 15}
	assert.Equal(t, expected, output)
}

func TestFloatingRange(t *testing.T) {
	iter := iterator.Range(0.0, 5.0)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 5, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{0.0, 1.0, 2.0, 3.0, 4.0}
	assert.Equal(t, expected, output)
}

func TestReverseFloatingRange(t *testing.T) {
	iter := iterator.Range(5.0, 0.0)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 5, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{5.0, 4.0, 3.0, 2.0, 1.0}
	assert.Equal(t, expected, output)
}

func TestFloatingRangeBy(t *testing.T) {
	iter := iterator.RangeBy(0.0, 5.0, 0.5)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 10, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{0.0, 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5}
	assert.Equal(t, expected, output)
	assert.Equal(t, 10, cap(output))
}

func TestInclusiveFloatingRangeBy(t *testing.T) {
	iter := iterator.IncRangeBy(0.0, 5.0, 0.5)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 11, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{0.0, 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5, 5.0}
	assert.Equal(t, expected, output)
}

func TestFloatingRangeDesc(t *testing.T) {
	iter := iterator.RangeBy(5.0, 0.0, -0.5)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 10, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{5.0, 4.5, 4.0, 3.5, 3.0, 2.5, 2.0, 1.5, 1.0, 0.5}
	assert.Equal(t, expected, output)
}

func TestInclusiveFloatingRangeDesc(t *testing.T) {
	iter := iterator.IncRangeBy(5.0, 0.0, -0.5)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 11, iter.Size().Allocate())
	output := iterator.Collect(iter)
	expected := []float64{5.0, 4.5, 4.0, 3.5, 3.0, 2.5, 2.0, 1.5, 1.0, 0.5, 0.0}
	assert.Equal(t, expected, output)
}

func TestInvalidRange(t *testing.T) {
	defer func() {
		r := recover()
		require.NotNil(t, r)
		assert.ErrorIs(t, r.(error), iterator.ErrInvalidIteratorRange)
	}()
	iterator.RangeBy(5.0, 0.0, 0.5)
}

func TestZeroRange(t *testing.T) {
	iter := iterator.RangeBy(1, 1, 0)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 0, iter.Size().Allocate())
	output := iterator.Collect(iter)
	assert.Empty(t, output)
}

func TestInclusiveZeroRange(t *testing.T) {
	iter := iterator.IncRangeBy(1, 1, 0)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 1, iter.Size().Allocate())
	output := iterator.Collect(iter)
	assert.Equal(t, []int{1}, output)
}

func TestSliceIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	output := iterator.Collect(iter)
	assert.Equal(t, input, output)
}

func TestTake(t *testing.T) {
	input := slices.Range(0, 10)
	iter := iterator.Take(4, iterator.Slice(input))
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 4, iter.Size().Allocate())
	output := iterator.Collect(iter)
	assert.Equal(t, slices.Range(0, 4), output)
}

func TestTakeNext(t *testing.T) {
	assert := assert.New(t)
	input := slices.Range(0, 10)
	sliceIter := iterator.Slice(input)
	iter := iterator.Take(4, sliceIter)
	assert.True(iter.Next())
	assert.True(iterator.IsSizeKnown(iter.Size()))
	assert.Equal(3, iter.Size().Allocate())
	output := iterator.Collect(iter)
	assert.Equal(slices.Range(1, 4), output)
	assert.True(iterator.IsSizeKnown(sliceIter.Size()))
	assert.Equal(6, sliceIter.Size().Allocate())
	remain := iterator.Collect(sliceIter)
	assert.Equal(slices.Range(4, 10), remain)
	assert.Equal(6, cap(remain))
}

func TestTakeMore(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Take(10, iterator.Slice(input))
	assert.True(t, iterator.IsSizeKnown(iter.Size()))
	assert.Equal(t, 4, iter.Size().Allocate())
	output := iterator.Collect(iter)
	assert.Equal(t, input, output)
}

func TestSliceIterString(t *testing.T) {
	input := []string{"one", "two", "three", "four"}
	iter := iterator.Slice(input)
	output := iterator.Collect(iter)
	assert.Equal(t, input, output)
}

func TestSliceMutIterRef(t *testing.T) {
	input := []string{"one", "two", "three", "four"}
	iter := iterator.MutSlice(&input)
	for iter.Next() {
		iter.Set(strings.ToUpper(iter.Value()))
	}
	expected := []string{"ONE", "TWO", "THREE", "FOUR"}
	assert.Equal(t, expected, input)
}

func TestSliceMutIterDelete(t *testing.T) {
	input := []string{"one", "two", "three", "four"}
	iter := iterator.MutSlice(&input)
	for iter.Next() {
		if iter.Value() == "two" {
			iter.Delete()
		}
	}
	expected := []string{"one", "three", "four"}
	assert.Equal(t, expected, input)
}

func TestSliceIterChan(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	i := 1
	for v := range iter.Chan() {
		assert.Equal(t, v, i)
		i++
	}
}

func TestSliceIterSeq(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	i := 1
	for v := range iter.Seq() {
		assert.Equal(t, v, i)
		i++
	}
}

func TestSeqIterMix(t *testing.T) {
	assert := assert.New(t)
	itr := iterator.Range(0, 10)
	var first []int
	for v := range iterator.Take(5, itr).Seq() {
		first = append(first, v)
	}
	assert.Equal(slices.Range(0, 5), first)
	second := iterator.Collect(itr)
	assert.Equal(slices.Range(5, 10), second)
	assert.Equal(5, cap(second))
}

func TestRange(t *testing.T) {
	r := iterator.Range(0, 10)
	seq := iterator.Collect(r)
	assert.Equal(t, 10, len(seq))
	for i, v := range seq {
		assert.Equal(t, i, v)
	}
}

func TestInclusiveRange(t *testing.T) {
	r := iterator.IncRange(0, 10)
	seq := iterator.Collect(r)
	assert.Equal(t, 11, len(seq))
	for i, v := range seq {
		assert.Equal(t, i, v)
	}
}

func TestRangeChan(t *testing.T) {
	r := iterator.Range(0, 10)
	i := 0
	for v := range r.Chan() {
		assert.Equal(t, i, v)
		i += 1
	}
	assert.Equal(t, 10, i)
}

func TestRangeSeq(t *testing.T) {
	r := iterator.Range(0, 10)
	i := 0
	for v := range r.Seq() {
		assert.Equal(t, i, v)
		i += 1
	}
	assert.Equal(t, 10, i)
}

func TestReverseRangeSeq(t *testing.T) {
	r := iterator.Range(10, 0)
	i := 10
	for v := range r.Seq() {
		assert.Equal(t, i, v)
		i -= 1
	}
	assert.Equal(t, 0, i)
}

func TestInclusiveRangeChan(t *testing.T) {
	r := iterator.IncRange(0, 10)
	i := 0
	for v := range r.Chan() {
		assert.Equal(t, i, v)
		i += 1
	}
	assert.Equal(t, 11, i)
}

func TestInclusiveRangeSeq(t *testing.T) {
	r := iterator.IncRange(0, 10)
	i := 0
	for v := range r.Seq() {
		assert.Equal(t, i, v)
		i += 1
		assert.Equal(t, 11-i, r.Size().Size)
	}
	assert.Equal(t, 11, i)
}

func TestRangeFor(t *testing.T) {
	r := iterator.Range(0, 10)
	i := 0
	for r.Next() {
		assert.Equal(t, i, r.Value())
		i += 1
		assert.Equal(t, 10-i, r.Size().Size)
	}
}

func TestEmptyRange(t *testing.T) {
	r := iterator.Range(10, 10)
	seq := iterator.Collect(r)
	assert.Empty(t, seq)
}

func TestEmptySeq(t *testing.T) {
	e := iterator.Empty[int]()
	for range e.Seq() {
		assert.Fail(t, "empty iterator should produce no values")
	}
	slice := iterator.Collect(e)
	assert.Empty(t, slice)
}

func TestNegativeRange(t *testing.T) {
	r := iterator.RangeBy(9, -1, -1)
	assert.True(t, iterator.IsSizeKnown(r.Size()))
	assert.Equal(t, 10, r.Size().Allocate())
	seq := iterator.Collect(r)
	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	assert.Equal(t, expected, seq)
}

func TestEmptyNegativeRange(t *testing.T) {
	r := iterator.RangeBy(10, 10, -1)
	seq := iterator.Collect(r)
	assert.Empty(t, seq)
}

func TestRangeBy(t *testing.T) {
	r := iterator.RangeBy(0, 9, 3)
	seq := iterator.Collect(r)
	assert.Equal(t, 3, len(seq))
	for i, v := range seq {
		assert.Equal(t, i*3, v)
	}
}

func TestInclusiveRangeBy(t *testing.T) {
	r := iterator.IncRangeBy(0, 9, 3)
	seq := iterator.Collect(r)
	assert.Equal(t, 4, len(seq))
	for i, v := range seq {
		assert.Equal(t, i*3, v)
	}
}

func TestEnumeratedRange(t *testing.T) {
	e := iterator.Range(10, 20).Enumerate()
	i := 0
	for e.Next() {
		assert.Equal(t, i, e.Key())
		assert.Equal(t, i+10, e.Value())
		i++
	}
}

func TestEnumeratedRangeSeq(t *testing.T) {
	e := iterator.Range(10, 20).Enumerate()
	i := 0
	for n, v := range e.Seq2() {
		assert.Equal(t, i, n)
		assert.Equal(t, i+10, v)
		assert.Equal(t, i, e.Key())
		i++
	}
}

func TestSliceIterChanAbort(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	i := 1
	for v := range iter.Chan() {
		assert.Equal(t, v, i)
		i++
		if i == 3 {
			iter.Abort()
		}
	}
	assert.Equal(t, i, 3)
}

func TestSliceIterSeqAbort(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	i := 1
	for v := range iter.Seq() {
		assert.Equal(t, v, i)
		i++
		if i == 3 {
			iter.Abort()
		}
	}
	assert.Equal(t, i, 3)
}

func TestSliceIterSeqBreak(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := iterator.Slice(input)
	i := 1
	for v := range iter.Seq() {
		assert.Equal(t, v, i)
		i++
		if i == 3 {
			break
		}
	}
	assert.Equal(t, i, 3)
	assert.False(t, iter.Next())
}

func TestMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := iterator.Collect(iterator.Map(iterator.Slice(input), func(n int) int { return n * 2 }))
	assert.Equal(t, expected, actual)
}

func TestMapIterChan(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	mi := iterator.Map(iterator.Slice(input), func(n int) int { return n * 2 })
	actual := make([]int, 0, mi.Size().Allocate())
	for v := range mi.Chan() {
		actual = append(actual, v)
	}
	assert.Equal(t, expected, actual)
}

func TestMapIterSeq(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	mi := iterator.Map(iterator.Slice(input), func(n int) int { return n * 2 })
	actual := make([]int, 0, mi.Size().Allocate())
	size := len(input)
	assert.Equal(t, size, mi.Size().Size)
	for v := range mi.Seq() {
		actual = append(actual, v)
		size--
		assert.Equal(t, size, mi.Size().Size)
	}
	assert.Equal(t, expected, actual)
}

func TestMapIterChanAbort(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4}
	mi := iterator.Map(iterator.Slice(input), func(n int) int { return n * 2 })
	actual := make([]int, 0, mi.Size().Allocate())
	i := 0
	for v := range mi.Chan() {
		actual = append(actual, v)
		i++
		if i == 2 {
			mi.Abort()
		}
	}
	assert.Equal(t, expected, actual)
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4}
	iter := iterator.Filter(iterator.Slice(input), func(n int) bool { return n&1 == 0 })
	assert.True(t, iterator.IsSizeAtMost(iter.Size()))
	actual := iterator.Collect(iter)
	assert.Equal(t, expected, actual)
}

func TestDoMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6}
	expectedErr := "Value 4 too large"
	actual, err := iterator.CollectResults(iterator.Map(iterator.Slice(input), func(n int) result.Result[int] {
		if n < 4 {
			return result.Value(n * 2)
		} else {
			return result.Error[int](fmt.Errorf("Value %d too large", n))
		}
	}))
	assert.Equal(t, expected, actual)
	assert.EqualError(t, err, expectedErr)
}

func TestGenerator(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := iterator.Collect(gen)
	expected := iterator.Collect(iterator.Range(0, 10))
	assert.Equal(t, expected, actual)
}

func TestFromSeq(t *testing.T) {
	seq := func(yield func(int) bool) {
		for i := range 5 {
			if !yield(i) {
				break
			}
		}
	}
	itr := iterator.New(seq)
	slice := iterator.Collect(itr)
	assert.Equal(t, slices.Range(0, 5), slice)
}

func TestFromSeqToChan(t *testing.T) {
	seq := func(yield func(int) bool) {
		for i := range 5 {
			if !yield(i) {
				break
			}
		}
	}
	itr := iterator.New(seq)
	slice := make([]int, 0, 5)
	for i := range itr.Chan() {
		slice = append(slice, i)
	}
	assert.Equal(t, slices.Range(0, 5), slice)
}

func fib() iterator.Iterator[int] {
	return iterator.Generate(func(c iterator.Consumer[int]) {
		tail := [2]int{0, 1}
		c.Yield(tail[0])
		c.Yield(tail[1])
		for {
			next := tail[0] + tail[1]
			c.Yield(next)
			tail[0] = tail[1]
			tail[1] = next
		}
	})
}

func fibPureSeq(yield func(int) bool) {
	tail := [2]int{0, 1}
	if !(yield(tail[0]) && yield(tail[1])) {
		return
	}
	for {
		next := tail[0] + tail[1]
		if !yield(next) {
			return
		}
		tail[0] = tail[1]
		tail[1] = next
	}
}

func fibSeq() iterator.Iterator[int] {
	return iterator.New(fibPureSeq)
}

func TestGenerateFib(t *testing.T) {
	result := iterator.Collect(iterator.Take(10, fib()))
	var expected = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	assert.Equal(t, expected, result)
}

func TestGenerateFibChan(t *testing.T) {
	var result []int
	var expected = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	for e := range iterator.Take(10, fib()).Chan() {
		result = append(result, e)
	}
	assert.Equal(t, expected, result)
}

func BenchmarkGenerateSimpleFib(b *testing.B) {
	iter := newFib()
	var sum uint64 = 0
	var count int
	for range b.N {
		iter.Next()
		sum += uint64(iter.Value())
		count++
	}
	assert.Equal(b, b.N, count)
}

func BenchmarkGenerateFib(b *testing.B) {
	iter := fib()
	var sum uint64 = 0
	var count int
	for range b.N {
		iter.Next()
		sum += uint64(iter.Value())
		count++
	}
	assert.Equal(b, b.N, count)
}

func BenchmarkGenerateTakeFib(b *testing.B) {
	iter := iterator.Take(b.N, fib())
	var sum uint64 = 0
	var count int
	for iter.Next() {
		sum += uint64(iter.Value())
		count++
	}
	assert.Equal(b, b.N, count)
}

func BenchmarkGenerateTakeFibSeq(b *testing.B) {
	iter := iterator.Take(b.N, fibSeq())
	var sum uint64 = 0
	for v := range iter.Seq() {
		sum += uint64(v)
	}
}

func BenchmarkGenerateMorphFibSeq(b *testing.B) {
	iter := fibSeq().Morph(func(n int) int { return n * 2 })
	var sum uint64 = 0
	i := 0
	for v := range iter.Seq() {
		if i >= b.N {
			break
		}
		sum += uint64(v)
		i++
	}
}

func BenchmarkGenerateFibSeq(b *testing.B) {
	iter := fibSeq()
	var sum uint64 = 0
	i := 0
	for v := range iter.Seq() {
		if i >= b.N {
			break
		}
		sum += uint64(v)
		i++
	}
}

func BenchmarkGenerateFibSeqPure(b *testing.B) {
	iter := fibPureSeq
	var sum uint64 = 0
	i := 0
	for v := range iter {
		if i >= b.N {
			break
		}
		sum += uint64(v)
		i++
	}
}

func BenchmarkGenerateFibSeqPull(b *testing.B) {
	itr := fibPureSeq
	next, stop := iter.Pull(itr)
	defer stop()
	var sum uint64 = 0
	var count int
	for range b.N {
		if v, ok := next(); !ok {
			break
		} else {
			sum += uint64(v)
			count++
		}
	}
	assert.Equal(b, b.N, count)
}

func BenchmarkGenerateFib2(b *testing.B) {
	iter := fib()
	defer iter.Abort()
	var sum uint64 = 0
	for i := 0; i < b.N && iter.Next(); i++ {
		sum += uint64(iter.Value())
	}
}

func BenchmarkGenerateFibChan(b *testing.B) {
	var sum uint64 = 0
	for v := range iterator.Take(b.N, fib()).Chan() {
		sum += uint64(v)
	}
}

func BenchmarkGenerateFibChan2(b *testing.B) {
	var sum uint64 = 0
	i := 0
	iter := fib()
	for v := range iter.Chan() {
		sum += uint64(v)
		i++
		if i > b.N {
			iter.Abort()
		}
	}
}

type SimpleFib [2]int

func NewSimpleFib() iterator.SimpleIterator[int] {
	return &SimpleFib{0, 1}
}

func (sf *SimpleFib) Value() int {
	return sf[0]
}

func (sf *SimpleFib) Next() bool {
	if sf[1] == 0 {
		return false
	} else {
		sf[0], sf[1] = sf[1], sf[0]+sf[1]
		// *sf = SimpleFib{sf[1], sf[0] + sf[1]}
		return true
	}
}

func (sf *SimpleFib) Abort() {
	sf[1] = 0
}

func newFib() iterator.Iterator[int] {
	return iterator.MakeIterator[int](NewSimpleFib())
}

func TestSimpleFib(t *testing.T) {
	fib := newFib()
	seq := iterator.Collect(iterator.Take(10, fib))
	expected := []int{1, 1, 2, 3, 5, 8, 13, 21, 34, 55}
	assert.Equal(t, expected, seq)
}

func repeatSeq[T any](r int, v T) func(func(T) bool) {
	return func(yield func(T) bool) {
		for range r {
			if !yield(v) {
				break
			}
		}
	}
}

var repeatSeqInt func(int, int) func(func(int) bool) = repeatSeq[int]

func repeatSeqIter[T any](r int, v T) iterator.Iterator[T] {
	return iterator.New(repeatSeq(r, v))
}

type repeatSimpleIter[T any] struct {
	iterator.DefaultIterator[T]
	index, repetitions int
	value              T
}

func (rsi *repeatSimpleIter[T]) Next() bool {
	if rsi.index < rsi.repetitions {
		rsi.index++
		return true
	} else {
		return false
	}
}

func (rsi *repeatSimpleIter[T]) Value() T {
	return rsi.value
}

func (rsi *repeatSimpleIter[T]) Abort() {
	rsi.index = rsi.repetitions
}

func repeatIter[T any](r int, v T) iterator.Iterator[T] {
	rsi := &repeatSimpleIter[T]{repetitions: r, value: v}
	rsi.CoreIterator = rsi
	return rsi
}

func BenchmarkBaseSeq(b *testing.B) {
	var sum uint64
	const value = 3
	seqIter := repeatSeq(b.N, value)
	for v := range seqIter {
		sum += uint64(v)
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func BenchmarkBaseSeqNonopt(b *testing.B) {
	var sum uint64
	const value = 3
	seqIter := repeatSeqInt(b.N, value)
	for v := range seqIter {
		sum += uint64(v)
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func BenchmarkBaseSeqIter(b *testing.B) {
	var sum uint64
	const value = 3
	seqIter := repeatSeqIter(b.N, value)
	for v := range seqIter.Seq() {
		sum += uint64(v)
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func BenchmarkBaseSimple(b *testing.B) {
	var sum uint64
	const value = 3
	itr := repeatIter(b.N, value)
	for itr.Next() {
		sum += uint64(itr.Value())
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func BenchmarkSeqSimple(b *testing.B) {
	var sum uint64
	const value = 3
	itr := repeatIter(b.N, value)
	for v := range itr.Seq() {
		sum += uint64(v)
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func BenchmarkBSimpleSeq(b *testing.B) {
	var sum uint64
	const value = 3
	seqIter := repeatSeqIter(b.N, value)
	for seqIter.Next() {
		sum += uint64(seqIter.Value())
	}
	assert.Equal(b, value*uint64(b.N), sum)
}

func rangeSum(from, to int) uint64 {
	return (uint64(to-from) * (uint64(from) + uint64(to-1))) / 2
}

func BenchmarkSeqRange(b *testing.B) {
	var sum uint64
	for v := range iterator.Range(0, b.N).Seq() {
		sum += uint64(v)
	}
	assert.Equal(b, rangeSum(0, b.N), sum)
}

func BenchmarkSimpleRange(b *testing.B) {
	var sum uint64
	for itr := iterator.Range(0, b.N); itr.Next(); {
		sum += uint64(itr.Value())
	}
	assert.Equal(b, rangeSum(0, b.N), sum)
}

func BenchmarkSeqFromSimpleRange(b *testing.B) {
	var sum uint64
	for v := range iterator.SimpleToSeq(iterator.Range(0, b.N)) {
		sum += uint64(v)
	}
	assert.Equal(b, rangeSum(0, b.N), sum)
}

func TestGeneratorChan(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[int]) {
		for i := range 10 {
			c.Yield(i)
		}
	})
	actual := make([]int, 10)
	expected := iterator.Collect(iterator.Range(0, 10))
	p := 0
	for i := range gen.Chan() {
		actual[p] = i
		p++
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorChanAbort(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[int]) {
		for i := range 10 {
			c.Yield(i)
		}
	})
	actual := make([]int, 5)
	expected := iterator.Collect(iterator.Range(0, 5))
	p := 0
	for i := range gen.Chan() {
		actual[p] = i
		p++
		if p >= len(actual) {
			gen.Abort()
		}
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorIterAbort(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[int]) {
		for i := range 10 {
			c.Yield(i)
		}
	})
	actual := make([]int, 5)
	expected := iterator.Collect(iterator.Range(0, 5))
	p := 0
	for gen.Next() {
		actual[p] = gen.Value()
		p++
		if p >= len(actual) {
			gen.Abort()
		}
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorMap(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[int]) {
		for i := range 10 {
			c.Yield(i)
		}
	})
	actual := iterator.Collect(iterator.Map(gen, func(x int) int { return x * 3 }))
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i * 3
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorError(t *testing.T) {
	gen := iterator.Generate(func(c iterator.Consumer[result.Result[int]]) {
		for i := range 10 {
			c.Yield(result.Value(i))
		}
		c.Yield(result.Error[int](fmt.Errorf("iterator failed")))
	})
	actual, err := iterator.CollectResults(gen)
	expected := iterator.Collect(iterator.Range(0, 10))
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
	assert.EqualError(t, err, "iterator failed")
}

// A result generator will yield an error if the generator
// function returns an error.
func TestGeneratorResultError(t *testing.T) {
	gen := iterator.GenerateResults(func(c iterator.ResultConsumer[int]) error {
		for i := range 10 {
			c.YieldValue(i)
		}
		return fmt.Errorf("iterator failed")
	})
	actual, err := iterator.CollectResults(gen)
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
	assert.EqualError(t, err, "iterator failed")
}

// It's possible to use Try without a handler in a result iterator;
// an error result will be automatically yielded and the iterator
// stopped.
func TestGeneratorResultTry(t *testing.T) {
	validate := func(x int) (int, error) {
		if x == 7 {
			return 0, fmt.Errorf("I don't like %d", x)
		} else {
			return x, nil
		}
	}
	gen := iterator.GenerateResults(func(c iterator.ResultConsumer[int]) error {
		for i := range 10 {
			c.YieldValue(eh.Try(validate(i)))
		}
		return nil
	})
	actual, err := iterator.CollectResults(gen)
	expected := iterator.Collect(iterator.Range(0, 7))
	assert.Equal(t, expected, actual)
	assert.EqualError(t, err, "I don't like 7")
}

// FilterResults will filter out error results
// and just return good values
func TestFilterSuccess(t *testing.T) {
	validate := func(x int) (int, error) {
		if x == 7 {
			return 0, fmt.Errorf("I don't like %d", x)
		} else {
			return x, nil
		}
	}
	gen := iterator.GenerateResults(func(c iterator.ResultConsumer[int]) error {
		for i := 0; i < 10; i++ {
			c.Yield(result.From(validate(i)))
		}
		return nil
	})
	actual := iterator.Collect(iterator.FilterValues(gen))
	expected := []int{0, 1, 2, 3, 4, 5, 6, 8, 9}
	assert.Equal(t, expected, actual)
}

// PartitionResults will collect non-error and error results in two separate
// slices.
func TestPartitionResults(t *testing.T) {
	validate := func(x int) (int, error) {
		if x == 3 || x == 7 {
			return 0, fmt.Errorf("I don't like %d", x)
		} else {
			return x, nil
		}
	}
	gen := iterator.Map(iterator.Range(0, 10), func(i int) result.Result[int] { return result.From(validate(i)) })
	actual, errs := iterator.PartitionResults(gen)
	expected := []int{0, 1, 2, 4, 5, 6, 8, 9}
	expectedErrs := []error{fmt.Errorf("I don't like 3"), fmt.Errorf("I don't like 7")}
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErrs, errs)
}

func TestAll(t *testing.T) {
	trueInput := []rune("---------")
	assert.True(t, iterator.All(iterator.Slice(trueInput), func(r rune) bool {
		return r == '-'
	}))
	falseInput := []rune("-----!----")
	assert.False(t, iterator.All(iterator.Slice(falseInput), func(r rune) bool {
		return r == '-'
	}))
}

func TestAny(t *testing.T) {
	trueInput := []rune("-----!----")
	assert.True(t, iterator.Any(iterator.Slice(trueInput), func(r rune) bool {
		return r == '!'
	}))
	falseInput := []rune("---------")
	assert.False(t, iterator.Any(iterator.Slice(falseInput), func(r rune) bool {
		return r == '!'
	}))
}

func testSeqIter2(from, to int, mapping func(int) int) iterator.Iterator2[int, int] {
	size := to - from
	core := iterator.NewSeqCoreIterator2WithSize(
		func(yield func(int, int) bool) {
			size = to - from
			for n := from; n < to; n++ {
				size--
				if !yield(n, mapping(n)) {
					break
				}
			}
		},
		func() iterator.IteratorSize {
			return iterator.NewSize(size)
		},
	)
	return iterator.NewDefaultIterator2(core)
}

func TestIterator2Seq2(t *testing.T) {
	const size = 10
	itr := testSeqIter2(0, size, func(n int) int { return n + size })
	count := 0
	for k, v := range itr.Seq2() {
		assert.Equal(t, count, k)
		assert.Equal(t, count+size, v)
		count++
		assert.Equal(t, size-count, itr.Size().Size)
	}
}

func TestIterator2Seq(t *testing.T) {
	const size = 10
	itr := testSeqIter2(0, size, func(n int) int { return n + size })
	count := 0
	for v := range itr.Seq() {
		assert.Equal(t, count+size, v)
		count++
		assert.Equal(t, size-count, itr.Size().Size)
	}
}

func TestIterator2Pull(t *testing.T) {
	const size = 10
	itr := testSeqIter2(0, size, func(n int) int { return n + size })
	count := 0
	for itr.Next() {
		assert.Equal(t, count, itr.Key())
		assert.Equal(t, count+size, itr.Value())
		count++
		assert.Equal(t, size-count, itr.Size().Size)
	}
}

func TestIterator2PullAndCollect2(t *testing.T) {
	const size = 10
	itr := testSeqIter2(0, size, func(n int) int { return n + size })
	count := 0
	for count < size/2 && itr.Next() {
		fmt.Println(count)
		count++
	}
	coll := itr.Collect2()
	assert.Equal(t, size/2, len(coll))
	for i, pair := range coll {
		assert.Equal(t, i+size/2, pair.Key)
		assert.Equal(t, i+size*3/2, pair.Value)
	}
}

func TestIterator2PullAndCollect(t *testing.T) {
	const size = 10
	itr := testSeqIter2(0, size, func(n int) int { return n + size })
	count := 0
	for count < size/2 && itr.Next() {
		fmt.Println(count)
		count++
	}
	coll := itr.Collect()
	assert.Equal(t, size/2, len(coll))
	for i, value := range coll {
		assert.Equal(t, i+size*3/2, value)
	}
}

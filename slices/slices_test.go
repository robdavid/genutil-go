package slices

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	trueInput := []rune("---------")
	assert.True(t, All(trueInput, func(r rune) bool {
		return r == '-'
	}))
	falseInput := []rune("-----!----")
	assert.False(t, All(falseInput, func(r rune) bool {
		return r == '-'
	}))
}

func TestAllRef(t *testing.T) {
	trueInput := []rune("---------")
	assert.True(t, AllRef(trueInput, func(r *rune) bool {
		return *r == '-'
	}))
	falseInput := []rune("-----!----")
	assert.False(t, AllRef(falseInput, func(r *rune) bool {
		return *r == '-'
	}))
}

func TestAny(t *testing.T) {
	trueInput := []rune("-----!----")
	assert.True(t, Any(trueInput, func(r rune) bool {
		return r == '!'
	}))
	falseInput := []rune("----------")
	assert.False(t, Any(falseInput, func(r rune) bool {
		return r == '!'
	}))
}

func TestAnyRef(t *testing.T) {
	trueInput := []rune("-----!----")
	assert.True(t, AnyRef(trueInput, func(r *rune) bool {
		return *r == '!'
	}))
	falseInput := []rune("----------")
	assert.False(t, AnyRef(falseInput, func(r *rune) bool {
		return *r == '!'
	}))
}

func TestFind(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 5, Find(input, '!'))
	assert.True(t, Contains(input, '!'))
	inputNF := []rune("----------")
	assert.Equal(t, -1, Find(inputNF, '!'))
	assert.False(t, Contains(inputNF, '!'))
}

func TestFindFrom(t *testing.T) {
	input := []rune("!----!---!-")
	assert.Equal(t, 0, FindFrom(0, input, '!'))
	assert.Equal(t, 5, FindFrom(1, input, '!'))
	assert.Equal(t, 9, FindFrom(9, input, '!'))
	assert.Equal(t, -1, FindFrom(10, input, '!'))
}

func TestFindUsing(t *testing.T) {
	input := []int{1, 3, 5, 8, 9}
	assert.Equal(t, 3, FindUsing(input, func(x int) bool { return x%2 == 0 }))
	assert.Equal(t, -1, FindUsing(input, func(x int) bool { return x%7 == 0 }))
}

func TestFindUsingRef(t *testing.T) {
	input := []int{1, 3, 5, 8, 9}
	assert.Equal(t, 3, FindUsingRef(input, func(x *int) bool { return (*x)%2 == 0 }))
	assert.Equal(t, -1, FindUsingRef(input, func(x *int) bool { return (*x)%7 == 0 }))
}

func TestRFind(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFind(input, '!'))
	inputNF := []rune("----------")
	assert.Equal(t, -1, RFind(inputNF, '!'))
}

func TestRFindUsing(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFindUsing(input, func(r rune) bool { return r != '-' }))
	assert.Equal(t, -1, RFindUsing(input, func(r rune) bool { return r == '*' }))
}

func TestRFindUsingRef(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFindUsingRef(input, func(r *rune) bool { return *r != '-' }))
	assert.Equal(t, -1, RFindUsingRef(input, func(r *rune) bool { return *r == '*' }))
}

func TestMap(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := Map(sliceIn, func(x int) int { return x * 2 })
	assert.Equal(t, expected, actual)
}

func TestMapDifferentTypes(t *testing.T) {
	sliceIn := []string{"apple", "banana", "cherry", "strawberry"}
	expected := []int{5, 6, 6, 10}
	actual := Map(sliceIn, func(s string) int { return len(s) })
	assert.Equal(t, expected, actual)
}

func TestMapRef(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := MapRef(sliceIn, func(x *int) int { return *x * 2 })
	assert.Equal(t, expected, actual)
}

func TestMapI(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	MapI(slice, func(x int) int { return x * 2 })
	assert.Equal(t, expected, slice)
}

func TestMapRefI(t *testing.T) {
	slice := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	MapRefI(slice, func(x *int) int { return *x * 2 })
	assert.Equal(t, expected, slice)
}

func TestFold(t *testing.T) {
	sliceIn := make([]int, 10)
	for i := range sliceIn {
		sliceIn[i] = i + 1
	}
	total := Fold(0, sliceIn, func(a int, t int) int { return a + t })
	assert.Equal(t, 55, total)
}

func TestRef(t *testing.T) {
	sliceIn := make([]int, 10)
	for i := range sliceIn {
		sliceIn[i] = i + 1
	}
	total := FoldRef(0, sliceIn, func(a *int, t *int) { *a += *t })
	assert.Equal(t, 55, total)
}

func TestConcat(t *testing.T) {
	slicesIn := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	sliceOut := Concat(slicesIn...)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, sliceOut)
}

func TestReverse(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4, 5, 6}
	assert.Equal(t, []int{6, 5, 4, 3, 2, 1}, Reverse(sliceIn))
	sliceIn = append(sliceIn, 7)
	assert.Equal(t, []int{7, 6, 5, 4, 3, 2, 1}, Reverse(sliceIn))
}

func TestReverseINil(t *testing.T) {
	var sliceIn []int = nil
	ReverseI(sliceIn)
	assert.Nil(t, sliceIn)
}

func TestReverseI(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4, 5, 6}
	ReverseI(sliceIn)
	assert.Equal(t, []int{6, 5, 4, 3, 2, 1}, sliceIn)
	ReverseI(sliceIn)
	sliceIn = append(sliceIn, 7)
	ReverseI(sliceIn)
	assert.Equal(t, []int{7, 6, 5, 4, 3, 2, 1}, sliceIn)
}

func TestFilterRef(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sliceOut := FilterRef(sliceIn, func(i *int) bool { return (*i)&1 == 0 })
	assert.Equal(t, []int{2, 4, 6, 8}, sliceOut)
}
func TestFilter(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	sliceOut := Filter(sliceIn, func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6, 8}, sliceOut)
}

func TestFilterRefI(t *testing.T) {
	sliceI := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	FilterRefI(&sliceI, func(i *int) bool { return (*i)&1 == 0 })
	assert.Equal(t, []int{2, 4, 6, 8}, sliceI)
}

func TestFilterI(t *testing.T) {
	sliceI := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	FilterI(&sliceI, func(i int) bool { return i%2 == 0 })
	assert.Equal(t, []int{2, 4, 6, 8}, sliceI)
}

func TestSortInt(t *testing.T) {
	sortableSlice := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	Sort(sortableSlice)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, sortableSlice)
}

func TestSortByte(t *testing.T) {
	sortableSlice := []byte{9, 8, 7, 6, 5, 4, 3, 2, 1}
	Sort(sortableSlice)
	assert.Equal(t, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, sortableSlice)
}

func TestSortString(t *testing.T) {
	sortableSlice := []string{"dates", "banana", "apple", "coconut"}
	Sort(sortableSlice)
	assert.Equal(t, []string{"apple", "banana", "coconut", "dates"}, sortableSlice)
}

func TestSortUsing(t *testing.T) {
	dates := []time.Time{
		time.Date(1988, 9, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1985, 9, 1, 0, 0, 0, 0, time.UTC),
		time.Date(1966, time.July, 30, 0, 0, 0, 0, time.UTC),
		time.Date(1974, time.August, 9, 0, 0, 0, 0, time.UTC),
	}
	SortUsing(dates, func(d1, d2 time.Time) bool { return d1.Before(d2) })
	fmt.Printf("%v\n", dates)
	for i := range dates {
		if i > 0 {
			assert.True(t, dates[i-1].Before(dates[i]))
		}
	}
}

var sorted []int

func BenchmarkSortUsing(b *testing.B) {
	var s []int
	for j := 0; j < b.N; j++ {
		items := make([]int, 100)
		for i := range items {
			items[i] = len(items) - i
		}
		SortUsing(items, func(i, j int) bool { return i < j })
		s = items
	}
	sorted = s
}

func BenchmarkSort(b *testing.B) {
	var s []int
	for j := 0; j < b.N; j++ {
		items := make([]int, 100)
		for i := range items {
			items[i] = len(items) - i
		}
		Sort(items)
		s = items
	}
	sorted = s
}

func TestEqual(t *testing.T) {
	var l []int = nil
	var r []int = nil
	assert.True(t, Equal(l, r))
	r = []int{}
	assert.True(t, Equal(l, r))
	r = append(r, 1)
	l = append(l, 1)
	assert.True(t, Equal(l, r))
	r = append(r, 2)
	l = append(l, 2)
	assert.True(t, Equal(l, r))
	r = append(r, 3)
	l = append(l, 4)
	assert.False(t, Equal(l, r))
	assert.False(t, Equal([]int{1, 2, 3}, []int{1, 2}))
}

func TestCompare(t *testing.T) {
	assert.Equal(t, 0, Compare([]int{}, nil))
	assert.Equal(t, 0, Compare([]int{}, []int{}))
	assert.Equal(t, -1, Compare([]int{1, 2}, []int{1, 3}))
	assert.Equal(t, 0, Compare([]int{1, 2}, []int{1, 2}))
	assert.Equal(t, 1, Compare([]int{1, 2, 4}, []int{1, 2, 3}))
}

func TestEmptyRange(t *testing.T) {
	for i := -9; i < 10; i++ {
		assert.Equal(t, []int{}, Range(i, i))
		assert.Equal(t, []int{}, RangeBy(i, i, -1))
		assert.Equal(t, []float64{}, Range(float64(i), float64(i)))
		assert.Equal(t, []int{}, ParRange(i, i))
		assert.Equal(t, []int{}, ParRangeBy(i, i, 0))
		assert.Equal(t, []int{}, ParRange(i, i, ParThreshold(0)))
	}
}

func TestSingletonIncRange(t *testing.T) {
	for i := -9; i < 10; i++ {
		assert.Equal(t, []int{i}, IncRange(i, i))
		assert.Equal(t, []int{i}, IncRangeBy(i, i, -1))
		assert.Equal(t, []float64{float64(i)}, IncRange(float64(i), float64(i)))
		assert.Equal(t, []int{i}, ParIncRange(i, i))
		assert.Equal(t, []int{i}, ParIncRangeBy(i, i, 0))
		assert.Equal(t, []int{i}, ParIncRange(i, i, ParThreshold(0)))
	}
}

func TestSingletonRange(t *testing.T) {
	assert.Equal(t, []int{0}, Range(0, 1))
	assert.Equal(t, []float64{0.0}, Range(0.0, 0.9))
}

func TestSimpleRange(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Range(0, 5))
	assert.Equal(t, []int{-2, -1}, Range(-2, 0))
	assert.Equal(t, []int{-2, -1, 0, 1}, Range(-2, 2))
	assert.Equal(t, []float64{0.0, 1.0, 2.0, 3.0, 4.0}, Range(0.0, 5.0))
	assert.Equal(t, []float64{-2.0, -1.5, -1.0, -0.5}, RangeBy(-2.0, 0.0, 0.5))
	assert.Equal(t, []float64{-1.0, -0.5, 0.0, 0.5}, RangeBy(-1.0, 1.0, 0.5))
}

func TestSimpleInclusiveRange(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2, 3, 4, 5}, IncRange(0, 5))
	assert.Equal(t, []int{-2, -1, 0}, IncRange(-2, 0))
	assert.Equal(t, []int{-2, -1, 0, 1, 2}, IncRange(-2, 2))
	assert.Equal(t, []float64{0.0, 1.0, 2.0, 3.0, 4.0, 5.0}, IncRange(0.0, 5.0))
	assert.Equal(t, []float64{-2.0, -1.5, -1.0, -0.5, 0.0}, IncRangeBy(-2.0, 0.0, 0.5))
	assert.Equal(t, []float64{-1.0, -0.5, 0.0, 0.5, 1.0}, IncRangeBy(-1.0, 1.0, 0.5))
}

func rangeBench(parallel bool, scale int, minPar int, numCpu int) {
	var size = minPar * numCpu
	slice := make([]int, size)
	if parallel {
		for iter := 0; iter < scale; iter++ {
			chunks := parChunks(slice, minPar, numCpu)
			parSliceFill(0, 1, false, chunks)
		}
	} else {
		for iter := 0; iter < scale; iter++ {
			sliceFill(0, 1, false, slice)
		}
	}
}

func BenchmarkParChunk(b *testing.B) {
	const parMin = 100000
	const numCpu = 4
	rangeBench(true, b.N, parMin, numCpu)
}

func BenchmarkOneChunk(b *testing.B) {
	const parMin = 100000
	const numCpu = 4
	rangeBench(false, b.N, parMin, numCpu)
}

func BenchmarkParSmallChunk(b *testing.B) {
	const parMin = 10000
	const numCpu = 4
	rangeBench(true, b.N, parMin, numCpu)
}

func BenchmarkOneSmallChunk(b *testing.B) {
	const parMin = 10000
	const numCpu = 4
	rangeBench(false, b.N, parMin, numCpu)
}

func TestLargeFloatRange(t *testing.T) {
	r := RangeBy(0.0, 1000000.0, 0.25)
	assert.Equal(t, 4000000, len(r))
	v := 0.0
	for _, e := range r {
		assert.Equal(t, v, e)
		v += 0.25
	}
}

func TestLargeIntRange(t *testing.T) {
	r := Range(0, 4000000)
	assert.Equal(t, 4000000, len(r))
	v := 0
	for _, e := range r {
		assert.Equal(t, v, e)
		v++
	}
}

func TestLargeInclusiveFloatRange(t *testing.T) {
	r := IncRangeBy(0.0, 1000000.0, 0.25)
	assert.Equal(t, 4000001, len(r))
	v := 0.0
	for _, e := range r {
		assert.Equal(t, v, e)
		v += 0.25
	}
	assert.Equal(t, iterator.Collect(iterator.IncRangeBy(0.0, 1000000, 0.25)), r)
}

func TestLargeInclusiveIntRange(t *testing.T) {
	r := IncRange(0, 4000000)
	assert.Equal(t, 4000001, len(r))
	v := 0
	for _, e := range r {
		assert.Equal(t, v, e)
		v++
	}
}

func TestInclusiveFullRange(t *testing.T) {
	full := IncRange(0, ^byte(0))
	assert.Equal(t, 256, len(full))
	for i := range full {
		assert.Equal(t, i, int(full[i]))
	}
}

func TestInclusiveSignedFullRange(t *testing.T) {
	full := IncRange(int8(-128), int8(127))
	assert.Equal(t, 256, len(full))
	for i := range full {
		assert.Equal(t, i-128, int(full[i]))
	}
}

func TestSillyStep(t *testing.T) {
	assert.PanicsWithError(t, "invalid range: step is zero", func() {
		RangeBy(0, 2, 0.5)
	})
}

func TestTruncatedStep(t *testing.T) {
	assert.Equal(t, []int{0, 1, 2}, IncRangeBy(0, 2, 1.5))
}

func TestSillyNegativeStep(t *testing.T) {
	assert.PanicsWithError(t, "invalid range: step is zero", func() {
		RangeBy(0, 2, -0.5)
	})
}

func TestIntRangeConsistency(t *testing.T) {
	for size := 0; size < 10000; size += 10 {
		for by := 1; by < 10; by++ {
			irange := iterator.Collect(iterator.RangeBy(0, size, by))
			srange := RangeBy(0, size, by)
			prange := ParRangeBy(0, size, by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.RangeBy(size, 0, -by))
			srange = RangeBy(size, 0, -by)
			prange = ParRangeBy(size, 0, -by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.IncRangeBy(0, size, by))
			srange = IncRangeBy(0, size, by)
			prange = ParIncRangeBy(0, size, by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.IncRangeBy(size, 0, -by))
			srange = IncRangeBy(size, 0, -by)
			prange = ParIncRangeBy(size, 0, -by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
		}
	}
}

func TestFloatRangeConsistency(t *testing.T) {
	for size := 0.0; size < 1000.0; size += 1.0 {
		for by := 1; by < 10; by++ {
			irange := iterator.Collect(iterator.RangeBy(0, size, by))
			srange := RangeBy(0, size, by)
			prange := ParRangeBy(0, size, by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.RangeBy(size, 0, -by))
			srange = RangeBy(size, 0, -by)
			prange = ParRangeBy(size, 0, -by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.IncRangeBy(0, size, by))
			srange = IncRangeBy(0, size, by)
			prange = ParIncRangeBy(0, size, by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
			irange = iterator.Collect(iterator.IncRangeBy(size, 0, -by))
			srange = IncRangeBy(size, 0, -by)
			prange = ParIncRangeBy(size, 0, -by, ParThreshold(1000), ParMaxCpu(4))
			assert.Equal(t, irange, srange)
			assert.Equal(t, irange, prange)
		}
	}
}

func TestParChunks4Core(t *testing.T) {
	inslice := make([]int, 100)
	chunks := parChunks(inslice, 10, 4)
	assert.Equal(t, 4, len(chunks))
	for _, chunk := range chunks {
		assert.Equal(t, 25, len(chunk))
	}
}

func TestParChunks16Core(t *testing.T) {
	inslice := make([]int, 100)
	chunks := parChunks(inslice, 10, 16)
	assert.Equal(t, 10, len(chunks))
	for _, chunk := range chunks {
		assert.Equal(t, 10, len(chunk))
	}
}

func TestUnequalParChunks(t *testing.T) {
	inslice := make([]int, 100)
	chunks := parChunks(inslice, 33, 16)
	assert.Equal(t, 3, len(chunks))
	assert.Equal(t, 34, len(chunks[0]))
	assert.Equal(t, 34, len(chunks[1]))
	assert.Equal(t, 32, len(chunks[2]))
}

func TestMultipleParChunksSize(t *testing.T) {
	inslice := make([]int, 100)
	for i := 1; i <= len(inslice)/2; i++ {
		chunks := parChunks(inslice, i, len(inslice))
		for _, chunk := range chunks {
			ideal := len(inslice) / len(chunks)
			diff := len(chunk) - ideal
			assert.LessOrEqual(t, diff, ideal/2)
		}
		sum := Fold(0, chunks, func(n int, c []int) int { return n + len(c) })
		assert.Equal(t, len(inslice), sum)
	}
}

func TestParRangeExample(t *testing.T) {
	actual := ParRangeBy[uint](400000, 0, -2)
	assert.Equal(t, 200000, len(actual))
	var last uint
	for i, v := range actual {
		assert.Equal(t, 400000-i*2, int(v))
		last = v
	}
	assert.Equal(t, uint(2), last)
}

func TestParIncRangeExample(t *testing.T) {
	actual := ParIncRangeBy[uint](400000, 0, -2)
	assert.Equal(t, 200001, len(actual))
	var last uint
	for i, v := range actual {
		assert.Equal(t, 400000-i*2, int(v))
		last = v
	}
	assert.Equal(t, uint(0), last)
}

func TestInvalidRange(t *testing.T) {
	assert.PanicsWithError(t, "invalid range: negative step or inverse range (but not both)",
		func() { RangeBy(0, 5, -1) },
	)
	assert.PanicsWithError(t, "invalid range: step is zero",
		func() { RangeBy(0.0, 0.5, 0.0) },
	)
}

func TestReverseRange(t *testing.T) {
	assert.Equal(t, []int{5, 4, 3, 2, 1}, Range(5, 0))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, Range(0, 5))
	assert.Equal(t, []int{5, 3, 1}, RangeBy(5, 0, -2))
	assert.Equal(t, []int{5, 3, 1}, IncRangeBy(5, 0, -2))
	assert.Equal(t, []int{4, 2, 0}, IncRangeBy(4, 0, -2))
	assert.Equal(t, IncRange(5, 0), Reverse(IncRange(0, 5)))
	assert.Equal(t, []uint{5, 4, 3, 2, 1}, Range[uint](5, 0))
	assert.Equal(t, []uint{5, 4, 3, 2, 1, 0}, IncRange[uint](5, 0))
	assert.Equal(t, []float64{5.0, 4.0, 3.0, 2.0, 1.0}, Range(5.0, 0.0))
	assert.Equal(t, []float64{132, 99, 66, 33}, Map(RangeBy(1.32, 0.0, -0.33), func(x float64) float64 { return math.Round(x * 100) }))
	assert.Equal(t, []float64{0.0, 0.33, 0.66, 0.99}, RangeBy(0.0, 1.32, 0.33))
}

func TestNonIntegerRange(t *testing.T) {
	assert.Equal(t, []float64{0.0, 0.5, 1.0, 1.5, 2.0}, RangeBy(0.0, 2.5, 0.5))
}

func TestNonIntegerReverseRange(t *testing.T) {
	assert.Equal(t, []float64{0.0, 0.5, 1.0, 1.5, 2.0}, RangeBy(0.0, 2.5, 0.5))
	reversed := Reverse(RangeBy(0.0, 2.5, 0.5))
	assert.Equal(t, reversed, RangeBy(2.0, -0.5, -0.5))
	//assert.Equal(t, reversed, RangeBy(2.5, 0.0, 0.5))
}

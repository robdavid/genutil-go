package iterator

import (
	"fmt"
	"testing"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/stretchr/testify/assert"
)

func TestCollectInto(t *testing.T) {
	iter := Range(0, 5)
	var output []int = nil
	CollectInto(iter, &output)
	iter2 := Range(10, 15)
	CollectInto(iter2, &output)
	expected := []int{0, 1, 2, 3, 4, 10, 11, 12, 13, 14}
	assert.Equal(t, expected, output)
}

func TestFloatingRange(t *testing.T) {
	iter := RangeBy(0.0, 5.0, 0.5)
	assert.True(t, IsSizeKnown(iter.Size()))
	assert.Equal(t, 10, iter.Size().(SizeKnown).Size)
	output := Collect(iter)
	expected := []float64{0.0, 0.5, 1.0, 1.5, 2.0, 2.5, 3.0, 3.5, 4.0, 4.5}
	assert.Equal(t, expected, output)
}

func TestSliceIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
	assert.True(t, IsSizeKnown(iter.Size()))
	output := Collect(iter)
	assert.Equal(t, input, output)
}

func TestTakeMore(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Take(10, Slice(input))
	assert.True(t, IsSizeKnown(iter.Size()))
	output := Collect(iter)
	assert.Equal(t, input, output)
}

func TestSliceIterString(t *testing.T) {
	input := []string{"one", "two", "three", "four"}
	iter := Slice(input)
	output := Collect(iter)
	assert.Equal(t, input, output)
}

func TestSliceIterChan(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
	i := 1
	for v := range iter.Chan() {
		assert.Equal(t, v, i)
		i++
	}
}

func TestRange(t *testing.T) {
	r := Range(0, 10)
	seq := Collect(r)
	for i, v := range seq {
		assert.Equal(t, i, v)
	}
}

func TestRangeChan(t *testing.T) {
	r := Range(0, 10)
	i := 0
	for v := range r.Chan() {
		assert.Equal(t, i, v)
		i += 1
	}
}

func TestRangeFor(t *testing.T) {
	r := Range(0, 10)
	i := 0
	for r.Next() {
		assert.Equal(t, i, r.Value())
		i += 1
	}
}

func TestEmptyRange(t *testing.T) {
	r := Range(10, 9)
	seq := Collect(r)
	assert.Empty(t, seq)
}

func TestNegativeRange(t *testing.T) {
	r := RangeBy(9, -1, -1)
	assert.True(t, IsSizeKnown(r.Size()))
	assert.Equal(t, 10, r.Size().(SizeKnown).Size)
	seq := Collect(r)
	expected := []int{9, 8, 7, 6, 5, 4, 3, 2, 1, 0}
	assert.Equal(t, expected, seq)
}

func TestEmptyNegativeRange(t *testing.T) {
	r := RangeBy(0, 10, -1)
	seq := Collect(r)
	assert.Empty(t, seq)
}

func TestRangeBy(t *testing.T) {
	r := RangeBy(0, 9, 3)
	seq := Collect(r)
	for i, v := range seq {
		assert.Equal(t, i*3, v)
	}
}

func TestSliceIterChanAbort(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
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

func TestMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := Collect(Map(Slice(input), func(n int) int { return n * 2 }))
	assert.Equal(t, expected, actual)
}

func TestFilter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4}
	iter := Filter(Slice(input), func(n int) bool { return n&1 == 0 })
	assert.True(t, IsSizeAtMost(iter.Size()))
	actual := Collect(iter)
	assert.Equal(t, expected, actual)
}

func TestDoMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6}
	expectedErr := "Value 4 too large"
	actual, err := CollectResults(Map(Slice(input), func(n int) result.Result[int] {
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
	gen := Generate(func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := Collect(gen)
	expected := Collect(Range(0, 10))
	assert.Equal(t, expected, actual)
}

func fib() Iterator[int] {
	return Generate(func(c Consumer[int]) {
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

func TestGenerateFib(t *testing.T) {
	result := Collect(Take(10, fib()))
	var expected = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	assert.Equal(t, expected, result)
}

func TestGenerateFibChan(t *testing.T) {
	var result []int
	var expected = []int{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	for e := range Take(10, fib()).Chan() {
		result = append(result, e)
	}
	assert.Equal(t, expected, result)
}

func BenchmarkGenerateFib(b *testing.B) {
	iter := Take(b.N, fib())
	var sum uint64 = 0
	for iter.Next() {
		sum += uint64(iter.Value())
	}
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
	for v := range Take(b.N, fib()).Chan() {
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

func TestGeneratorChan(t *testing.T) {
	gen := Generate(func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := make([]int, 10)
	expected := Collect(Range(0, 10))
	p := 0
	for i := range gen.Chan() {
		actual[p] = i
		p++
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorChanAbort(t *testing.T) {
	gen := Generate(func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := make([]int, 5)
	expected := Collect(Range(0, 5))
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
	gen := Generate(func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := make([]int, 5)
	expected := Collect(Range(0, 5))
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
	gen := Generate(func(c Consumer[int]) {
		for i := 0; i < 10; i++ {
			c.Yield(i)
		}
	})
	actual := Collect(Map(gen, func(x int) int { return x * 3 }))
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i * 3
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorError(t *testing.T) {
	gen := Generate(func(c Consumer[result.Result[int]]) {
		for i := 0; i < 10; i++ {
			c.Yield(result.Value(i))
		}
		c.Yield(result.Error[int](fmt.Errorf("iterator failed")))
	})
	actual, err := CollectResults(gen)
	expected := Collect(Range(0, 10))
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
	assert.EqualError(t, err, "iterator failed")
}

// A result generator will yield an error if the generator
// function returns an error.
func TestGeneratorResultError(t *testing.T) {
	gen := GenerateResults(func(c ResultConsumer[int]) error {
		for i := 0; i < 10; i++ {
			c.YieldValue(i)
		}
		return fmt.Errorf("iterator failed")
	})
	actual, err := CollectResults(gen)
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
	gen := GenerateResults(func(c ResultConsumer[int]) error {
		for i := 0; i < 10; i++ {
			c.YieldValue(eh.Try(validate(i)))
		}
		return nil
	})
	actual, err := CollectResults(gen)
	expected := Collect(Range(0, 7))
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
	gen := GenerateResults(func(c ResultConsumer[int]) error {
		for i := 0; i < 10; i++ {
			c.Yield(result.From(validate(i)))
		}
		return nil
	})
	actual := Collect(FilterValues(gen))
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
	gen := Map(Range(0, 10), func(i int) result.Result[int] { return result.From(validate(i)) })
	actual, errs := PartitionResults(gen)
	expected := []int{0, 1, 2, 4, 5, 6, 8, 9}
	expectedErrs := []error{fmt.Errorf("I don't like 3"), fmt.Errorf("I don't like 7")}
	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedErrs, errs)
}

func TestAll(t *testing.T) {
	trueInput := []rune("---------")
	assert.True(t, All(Slice(trueInput), func(r rune) bool {
		return r == '-'
	}))
	falseInput := []rune("-----!----")
	assert.False(t, All(Slice(falseInput), func(r rune) bool {
		return r == '-'
	}))
}

func TestAny(t *testing.T) {
	trueInput := []rune("-----!----")
	assert.True(t, Any(Slice(trueInput), func(r rune) bool {
		return r == '!'
	}))
	falseInput := []rune("---------")
	assert.False(t, Any(Slice(falseInput), func(r rune) bool {
		return r == '!'
	}))
}

package iterator

import (
	"fmt"
	"testing"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/stretchr/testify/assert"
)

func TestSliceIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
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

func TestEmptyRange(t *testing.T) {
	r := Range(10, 9)
	seq := Collect(r)
	assert.Empty(t, seq)
}

func TestNegativeRange(t *testing.T) {
	r := RangeBy(9, -1, -1)
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
	actual := Collect(Filter(Slice(input), func(n int) bool { return n&1 == 0 }))
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
	gen := Generate(func(y Yield[int]) {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
	})
	actual := Collect(gen)
	expected := Collect(Range(0, 10))
	assert.Equal(t, expected, actual)
}

func TestGeneratorChan(t *testing.T) {
	gen := Generate(func(y Yield[int]) {
		for i := 0; i < 10; i++ {
			y.Yield(i)
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
	gen := Generate(func(y Yield[int]) {
		for i := 0; i < 10; i++ {
			y.Yield(i)
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
	gen := Generate(func(y Yield[int]) {
		for i := 0; i < 10; i++ {
			y.Yield(i)
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
	gen := Generate(func(y Yield[int]) {
		for i := 0; i < 10; i++ {
			y.Yield(i)
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
	gen := Generate(func(y Yield[result.Result[int]]) {
		for i := 0; i < 10; i++ {
			y.Yield(result.Value(i))
		}
		y.Yield(result.Error[int](fmt.Errorf("iterator failed")))
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
	gen := GenerateResults(func(y YieldResult[int]) error {
		for i := 0; i < 10; i++ {
			y.YieldValue(i)
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
	gen := GenerateResults(func(y YieldResult[int]) error {
		for i := 0; i < 10; i++ {
			y.YieldValue(eh.Try(validate(i)))
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
	gen := GenerateResults(func(y YieldResult[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(result.From(validate(i)))
		}
		return nil
	})
	actual := Collect(FilterResults(gen))
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

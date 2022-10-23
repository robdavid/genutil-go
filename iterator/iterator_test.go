package iterator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSliceIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
	output := MustCollect(iter)
	assert.Equal(t, input, output)
}

func TestSliceIterChan(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
	i := 1
	for v := range iter.Chan() {
		assert.Equal(t, v.Get(), i)
		i++
	}
}

func TestSliceIterChanAbort(t *testing.T) {
	input := []int{1, 2, 3, 4}
	iter := Slice(input)
	i := 1
	for v := range iter.Chan() {
		assert.Equal(t, v.Get(), i)
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
	actual := MustCollect(Map(Slice(input), func(n int) int { return n * 2 }))
	assert.Equal(t, expected, actual)
}

func TestDoMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6}
	actual, err := Collect(DoMap(Slice(input), func(n int) (int, error) {
		if n < 4 {
			return n * 2, nil
		} else {
			return 0, fmt.Errorf("Value %d too large", n)
		}
	}))
	assert.Equal(t, expected, actual)
	if assert.Error(t, err) {
		assert.Equal(t, "Value 4 too large", err.Error())
	}
}

func TestGeneratorIter(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		return nil
	})
	actual := MustCollect(pipe)
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorIterChan(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		return nil
	})
	actual := make([]int, 10)
	expected := make([]int, 10)
	p := 0
	for i := range pipe.Chan() {
		actual[p] = i.Get()
		p++
	}
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorIterChanAbort(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		return nil
	})
	actual := make([]int, 5)
	expected := make([]int, 5)
	p := 0
	for i := range pipe.Chan() {
		assert.False(t, i.IsError())
		actual[p] = i.Get()
		p++
		if p >= len(actual) {
			pipe.Abort()
		}
	}
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorIterMap(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		return nil
	})
	actual := MustCollect(Map(pipe, func(x int) int { return x * 3 }))
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i * 3
	}
	assert.Equal(t, expected, actual)
}

func TestGeneratorIterError(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		return fmt.Errorf("iterator failed")
	})
	actual, err := Collect(pipe)
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
	if assert.Error(t, err) {
		assert.Equal(t, "iterator failed", err.Error())
	}
}

func TestGeneratorIterPanic(t *testing.T) {
	pipe := Generate(func(y Yield[int]) error {
		for i := 0; i < 10; i++ {
			y.Yield(i)
		}
		panic("iterator failed")
	})
	actual, err := Collect(pipe)
	expected := make([]int, 10)
	for i := range expected {
		expected[i] = i
	}
	assert.Equal(t, expected, actual)
	if assert.Error(t, err) {
		assert.Equal(t, GeneratorPanic{"iterator failed"}, err)
	}
}

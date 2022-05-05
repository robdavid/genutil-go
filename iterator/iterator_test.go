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

func TestMapIter(t *testing.T) {
	input := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := MustCollect(Map(Slice(input), func(n int) int { return n * 2 }))
	assert.Equal(t, expected, actual)
}

func TestPipeIter(t *testing.T) {
	pipe := Pipe(func(y Yield[int]) error {
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

func TestPipeIterMap(t *testing.T) {
	pipe := Pipe(func(y Yield[int]) error {
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

func TestPipeIterError(t *testing.T) {
	pipe := Pipe(func(y Yield[int]) error {
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

func TestPipeIterPanic(t *testing.T) {
	pipe := Pipe(func(y Yield[int]) error {
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
		assert.Equal(t, PipePanic{"iterator failed"}, err)
	}
}

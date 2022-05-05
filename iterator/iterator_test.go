package iterator

import (
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

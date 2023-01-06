package slices

import (
	"testing"

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

func TestRFind(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFind(input, '!'))
	inputNF := []rune("----------")
	assert.Equal(t, -1, RFind(inputNF, '!'))
}

func TestRFindUsing(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFindUsing(input, func(r rune) bool { return r != '-' }))
}

func TestRFindUsingRef(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFindUsingRef(input, func(r *rune) bool { return *r != '-' }))
}

func TestMap(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := Map(sliceIn, func(x int) int { return x * 2 })
	assert.Equal(t, expected, actual)
}

func TestMapRef(t *testing.T) {
	sliceIn := []int{1, 2, 3, 4}
	expected := []int{2, 4, 6, 8}
	actual := MapRef(sliceIn, func(x *int) int { return *x * 2 })
	assert.Equal(t, expected, actual)
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

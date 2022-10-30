package slice

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

func TestRFind(t *testing.T) {
	input := []rune("-----!---!-")
	assert.Equal(t, 9, RFind(input, '!'))
	inputNF := []rune("----------")
	assert.Equal(t, -1, RFind(inputNF, '!'))
}

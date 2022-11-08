package functions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIfElse(t *testing.T) {
	const yes = "yes"
	const no = "no"
	assert.Equal(t, IfElse(true, yes, no), "yes")
	assert.Equal(t, IfElse(false, yes, no), "no")
}

func TestIfElseF(t *testing.T) {
	const yes = "yes"
	const no = "no"
	assert.Equal(t, IfElseF(true, func() string { return yes }, func() string { return no }), "yes")
	assert.Equal(t, IfElseF(false, func() string { return yes }, func() string { return no }), "no")
}

func TestId(t *testing.T) {
	assert.Equal(t, 6, Id(6))
	assert.Equal(t, 12.3, Id(12.3))
	assert.Equal(t, "hello", Id("hello"))
}

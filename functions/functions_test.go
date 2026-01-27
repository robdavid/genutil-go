package functions_test

import (
	"testing"

	"github.com/robdavid/genutil-go/functions"
	"github.com/stretchr/testify/assert"
)

func TestIfElse(t *testing.T) {
	const yes = "yes"
	const no = "no"
	assert.Equal(t, functions.IfElse(true, yes, no), "yes")
	assert.Equal(t, functions.IfElse(false, yes, no), "no")
}

func TestIfElseF(t *testing.T) {
	const yes = "yes"
	const no = "no"
	assert.Equal(t, functions.IfElseF(true, func() string { return yes }, func() string { return no }), "yes")
	assert.Equal(t, functions.IfElseF(false, func() string { return yes }, func() string { return no }), "no")
}

func TestId(t *testing.T) {
	assert.Equal(t, 6, functions.Id(6))
	assert.Equal(t, 12.3, functions.Id(12.3))
	assert.Equal(t, "hello", functions.Id("hello"))
}

func TestRef(t *testing.T) {
	assert.Equal(t, 6, *functions.Ref(6))
	assert.Equal(t, 12.3, *functions.Ref(12.3))
	assert.Equal(t, "hello", *functions.Ref("hello"))
}

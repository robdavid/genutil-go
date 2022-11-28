package tuple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTuple0(t *testing.T) {
	t0 := Of0()
	assert.Equal(t, t0.Size(), 0)
	assert.Equal(t, t0.String(), "()")
	assert.Panics(t, func() { t0.Get(0) })
}

func TestTuple1(t *testing.T) {
	ta := Of1(10)
	assert.Equal(t, ta.Size(), 1)
	assert.Equal(t, ta.String(), "(10)")
	assert.Equal(t, ta.Get(0), 10)
	assert.Panics(t, func() { ta.Get(1) })
	assert.Equal(t, ta.First, 10)
	tb := Of1(10)
	assert.Equal(t, ta, tb)
	tb = Of1(11)
	assert.NotEqual(t, ta, tb)
}

func TestTuple2(t *testing.T) {
	t2 := Of2(10, 11)
	assert.Equal(t, t2.Size(), 2)
	assert.Equal(t, t2.String(), "(10,11)")
	assert.Equal(t, t2.Get(0), 10)
	assert.Equal(t, t2.Get(1), 11)
	assert.Panics(t, func() { t2.Get(2) })
	assert.Equal(t, t2.First, 10)
	assert.Equal(t, t2.Second, 11)
	assert.True(t, t2 == Of2(10, 11))
	assert.True(t, t2 != Of2(11, 12))
	assert.NotEqual(t, t2, Of1(11))
}

func TestTuple3(t *testing.T) {
	t3 := Of3(10, 11, 12)
	assert.Equal(t, t3.Size(), 3)
	assert.Equal(t, t3.String(), "(10,11,12)")
	assert.Equal(t, t3.Get(0), 10)
	assert.Equal(t, t3.Get(1), 11)
	assert.Equal(t, t3.Get(2), 12)
	assert.Panics(t, func() { t3.Get(3) })
	assert.Equal(t, t3.First, 10)
	assert.Equal(t, t3.Second, 11)
	assert.Equal(t, t3.Third, 12)
	assert.True(t, t3 == Of3(10, 11, 12))
	assert.False(t, t3 == Of3(11, 12, 13))
	assert.NotEqual(t, t3, Of2(10, 11))
	assert.Equal(t, Slice(&t3), []any{10, 11, 12})
}

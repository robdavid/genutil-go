package ordered

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbs(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(123, Abs(-123))
	assert.Equal(1.23, Abs(-1.23))
}

func TestIsInteger(t *testing.T) {
	assert := assert.New(t)
	assert.True(IsInteger(1))
	assert.True(IsInteger(int(0)))
	assert.True(IsInteger(byte(0)))
	assert.True(IsInteger(rune(0)))
	assert.False(IsInteger(1.2))
}

func TestPrecision(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(32, Precision(int32(0)))
	assert.Equal(8, Precision(byte(0)))
	assert.Equal(54, Precision(1.2))
	assert.Equal(25, Precision(float32(1.2)))
}

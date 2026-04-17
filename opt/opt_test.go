package opt_test

import (
	"testing"

	"github.com/robdavid/genutil-go/opt"
	"github.com/stretchr/testify/assert"
)

func TestOptZero(t *testing.T) {
	var intval opt.Val[int]
	var intref opt.Ref[int]
	assert.True(t, intval.IsEmpty())
	assert.True(t, intref.IsEmpty())
}

func TestOption(t *testing.T) {
	assert := assert.New(t)
	var x int = 123
	intval := opt.FromVal(x)
	intref := opt.FromRef(&x)
	assert.Equal(123, intval.Get())
	assert.Equal(123, intref.Get())
	optval := opt.Option[int](&intval)
	optref := opt.Option[int](intref)
	assert.Equal(123, optval.Get())
	assert.Equal(123, optref.Get())
	*optref.Ref() = 456
	assert.Equal(456, x)
	assert.Equal(123, intval.Get())
	*optval.Ref() = 456
	assert.Equal(456, intval.Get())
}

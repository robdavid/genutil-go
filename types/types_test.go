package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAs(t *testing.T) {
	var v  = 123
	opt := As[int](v).Get()
	assert.Equal(t, 123, opt)
	assert.Equal(t, []any(nil), As[[]any](v).GetOr(nil))
}

func TestAsRef(t *testing.T) {
	var n = 123
	var v any = &n
	var i *int = nil
	opt := AsRef[int](&n).Get()
	assert.Equal(t, 123, opt)
	assert.True(t, AsRef[int](i).IsEmpty())
	i = AsRef[int](v).GetRef()
	assert.Equal(t, 123, AsRef[int](i).Get())
}

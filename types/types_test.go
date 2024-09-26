package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAs(t *testing.T) {
	var n = 123
	var v any = &n
	opt := As[int](n).Get()
	assert.Equal(t, 123, opt)
	assert.Equal(t, []any(nil), As[[]any](n).GetOr(nil))
	opt2 := As[*int](v).Get()
	assert.Equal(t, 123, *opt2)
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

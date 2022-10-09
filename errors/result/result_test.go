package result

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/stretchr/testify/assert"
)

func arity2Err(err error) (string, int, error) {
	return "testValue", 123, err
}

func returnResult2[A any, B any](res2 *Result2[A, B]) (A, B, error) {
	return res2.Return()
}

func arity1Err(err error) (string, error) {
	return "testValue", err
}

func returnResult[T any](res *Result[T]) (T, error) {
	return res.Return()
}

func TestGoodResult(t *testing.T) {
	res := From(arity1Err(nil))
	assert.Equal(t, "testValue", res.String())
	// These asserts with == also designed to check compile time types
	assert.True(t, res.Get() == "testValue", "got value %s", res.Get())
	assert.NoError(t, res.GetErr())
	v, e := returnResult(&res)
	assert.True(t, v == "testValue", "got value %s", v)
	assert.True(t, e == nil, "got error %s", e)
}

func TestErrorResult(t *testing.T) {
	res := From(arity1Err(fmt.Errorf("This raises an error")))
	assert.Equal(t, "This raises an error", res.String())
	// These asserts with == also designed to check compile time types
	assert.True(t, res.Get() == "testValue", "got value %s", res.Get())
	assert.EqualError(t, res.GetErr(), "This raises an error")
	assert.Panics(t, func() {
		assert.True(t, res.Must() == "testValue")
	})
	v, e := returnResult(&res)
	assert.True(t, v == "testValue", "got value %s", v)
	assert.EqualError(t, e, "This raises an error")
}

func TestErrorMapTry(t *testing.T) {
	defer handler.Handle(func(err error) {
		assert.EqualError(t, err, "outer error: inner error")
	})
	value := New(arity1Err(fmt.Errorf("inner error"))).
		MapErr(func(err error) error { return fmt.Errorf("outer error: %w", err) }).Try()
	assert.Fail(t, "error not thrown; got value %v", value)
}

func TestErrorMapTrySuccess(t *testing.T) {
	defer handler.Handle(func(err error) {
		assert.NoError(t, err)
	})
	value := New(arity1Err(nil)).
		MapErr(func(err error) error { return fmt.Errorf("outer error: %w", err) }).Try()
	assert.Equal(t, "testValue", value)
}

func TestGoodResult2(t *testing.T) {
	res := From2(arity2Err(nil))
	assert.Equal(t, "(testValue,123)", res.String())
	// These asserts with == also designed to check compile time types
	assert.True(t, res.Get().First == "testValue", "got value %s", res.Get().First)
	assert.True(t, res.Get().Second == 123, "got value %d", res.Get().Second)
	assert.True(t, res.Must().First == "testValue", "got value %s", res.Get().First)
	assert.True(t, res.Must().Second == 123, "got value %d", res.Get().Second)
	assert.NoError(t, res.GetErr())
	a, b, e := returnResult2(&res)
	assert.True(t, a == "testValue", "got value %s", a)
	assert.True(t, b == 123, "got value %d", b)
	assert.True(t, e == nil, "got error %s", e)
}

func TestErrorResult2(t *testing.T) {
	res := From2(arity2Err(fmt.Errorf("This raises an error")))
	assert.Equal(t, "This raises an error", res.String())
	assert.EqualError(t, res.GetErr(), "This raises an error")
	// These asserts with == also designed to check compile time types
	assert.Panics(t, func() {
		assert.True(t, res.Must().First == "testValue")
		assert.True(t, res.Must().Second == 123)
	})
	a, b, e := returnResult2(&res)
	assert.True(t, a == "testValue", "got value %s", a)
	assert.True(t, b == 123, "got value %d", b)
	assert.EqualError(t, e, "This raises an error")
}

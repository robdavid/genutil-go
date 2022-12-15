package test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/tuple"
	"github.com/stretchr/testify/assert"
)

func dummySuccessFunction() (string, error) {
	return "success", nil
}

func dummySuccessFunction2() (string, int, error) {
	return "success", 123, nil
}

func dummyFailureFunction2() (string, int, error) {
	return "success", 123, fmt.Errorf("dummyFailureFunction2 has failed")
}

func dummySuccessFunction5() (string, int, float64, bool, rune, error) {
	return "success", 123, 4.56, true, 'x', nil
}

func TestSuccess(t *testing.T) {
	result := Result(dummySuccessFunction()).Must(t)
	assert.Equal(t, "success", result)
}

func BenchmarkSuccessful(b *testing.B) {
	result := Result(dummySuccessFunction()).Must(b)
	assert.Equal(b, "success", result)
}

func BenchmarkFailure2(b *testing.B) {
	var s string
	var i int
	var err error
	defer func() {
		assert.ErrorContains(b, err, "dummyFailureFunction2 has failed")
	}()
	defer handler.Handle(func(e error) {
		assert.Equal(b, "", s)
		assert.Equal(b, 0, i)
		err = e
	})
	s, i = Result2(dummyFailureFunction2()).Try2()
}

func TestSuccess5(t *testing.T) {
	s, i, f, b, r := Result5(dummySuccessFunction5()).Must5(t)
	assert.Equal(t, "success", s)
	assert.Equal(t, 123, i)
	assert.Equal(t, 4.56, f)
	assert.Equal(t, true, b)
	assert.Equal(t, 'x', r)
}

func TestATestableResult2IsAResult(t *testing.T) {
	result := Result2(dummySuccessFunction2())
	assert.False(t, result.IsError())
	assert.Equal(t, tuple.Of2("success", 123), result.Get())
}

func TestATestableErrorResult2IsAResult(t *testing.T) {
	var err error
	defer func() {
		assert.ErrorContains(t, err, "dummyFailureFunction2 has failed")
	}()
	defer handler.Catch(&err)
	res := Result2(dummyFailureFunction2())
	r := res.Result
	assert.True(t, r.IsError())
	assert.Equal(t, tuple.Of2("success", 123), r.Get())
	assert.True(t, res.IsError())
	assert.Equal(t, tuple.Of2("success", 123), res.Get())
	res.Try()
}

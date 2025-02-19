package test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/tuple"
	"github.com/stretchr/testify/assert"
)

var errTestErr = errors.New("this is my test error")
var errTestBottom = errors.New("hit the bottom")

func dummySuccessFunction() (string, error) {
	return "success", nil
}

func dummySuccessNoResultFunction() error {
	return nil
}

// func dummySuccessFunction2() (string, int, error) {
// 	return "success", 123, nil
// }

func dummyFailureFunction2() (string, int, error) {
	return "success", 123, fmt.Errorf("dummyFailureFunction2 has failed")
}

// func dummySuccessFunction5() (string, int, float64, bool, rune, error) {
// 	return "success", 123, 4.56, true, 'x', nil
// }

func dummyErrFunction() error {
	return errTestErr
}

type mockTesting struct {
	lastError []any
	errOutput []string
	fatal     bool
}

func (mt *mockTesting) Error(args ...any) {
	mt.lastError = args
}

func (mt *mockTesting) Errorf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	fmt.Println(message)
	mt.errOutput = append(mt.errOutput, message)
}

func (mt *mockTesting) FailNow() {
	mt.fatal = true
}

func (mt *mockTesting) Helper() {
}

func (mt *mockTesting) message() string {
	return fmt.Sprint(mt.lastError...)
}

//go:generate code-template test_test.tmpl -o tmpl_test.go

func TestSuccess(t *testing.T) {
	result := Result(dummySuccessFunction()).Must(t)
	assert.Equal(t, "success", result)
}

func TestSuccessNoResult(t *testing.T) {
	Status(dummySuccessNoResultFunction()).Must(t)
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

func TestChainedFailureAssertion(t *testing.T) {
	Status(dummyErrFunction()).Fails(t).FailsWith(t, errTestErr).FailsContaining(t, "test error")
}

func TestUnsatisfiedFailureAssertion(t *testing.T) {
	var mt, mt2, mt3 mockTesting
	var noterr error = errors.New("other error")
	Status(dummySuccessNoResultFunction()).Fails(&mt)
	assert.Equal(t, "an error was expected, but did not occur", mt.message())
	Status(dummyErrFunction()).FailsContaining(&mt2, "something else")
	assert.Equal(t, "expected error to contain 'something else', but was 'this is my test error'", mt2.message())
	Status(dummyErrFunction()).FailsWith(&mt3, noterr)
	assert.Equal(t, "expected error 'other error', but got 'this is my test error'", mt3.message())
}

func TestErrorReporting(t *testing.T) {
	var mt mockTesting
	defer func() {
		assert.Equal(t, "this is my test error", mt.message())
	}()
	defer ReportErr(&mt)
	Status(dummyErrFunction()).Try()
}

var dummyErrorRaiseLine int
var dummyErrorRaiseFile string

func dummyErrorRecurse(level int, catch bool) (err error) {
	if catch {
		defer handler.Catch(&err)
	}
	if level > 0 {
		dummyErrorRecurse(level-1, false)
	} else {

		_, dummyErrorRaiseFile, dummyErrorRaiseLine, _ = runtime.Caller(0)
		dummyErrorRaiseLine += 2
		handler.Check(errTestBottom)
	}
	return
}

func TestErrorTrace(t *testing.T) {
	var mt mockTesting
	var line int
	var file string
	defer func() {
		assert.Equal(t, "hit the bottom", mt.message())
		errLines := strings.Split(mt.errOutput[0], "\n")
		assert.Equal(t, "stack trace", errLines[0])
		parts := strings.Split(errLines[1], " ")
		assert.Equal(t, fmt.Sprintf("%s:%d", file, line+1), parts[0])
	}()
	defer ReportErr(&mt)
	_, file, line, _ = runtime.Caller(0)
	Status(dummyErrorRecurse(0, true)).Try()
}

func TestHandlerTryTrace(t *testing.T) {
	var mt mockTesting
	const depth = 10
	defer func() {
		assert.Equal(t, "hit the bottom", mt.message())
		errLines := strings.Split(mt.errOutput[0], "\n")
		assert.Greater(t, len(errLines), depth)
		assert.Equal(t, "stack trace", errLines[0])
		parts := strings.Split(errLines[1], " ")
		assert.Equal(t, fmt.Sprintf("%s:%d", dummyErrorRaiseFile, dummyErrorRaiseLine), parts[0])
	}()
	defer ReportErr(&mt)
	handler.Check(dummyErrorRecurse(depth, false))
}

func TestNoError(t *testing.T) {
	defer ReportErr(t)
}

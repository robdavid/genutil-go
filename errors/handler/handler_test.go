package handler

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readFileTest(fname string) (content []byte, err error) {
	defer Catch(&err)        // If a try fails, the wrapped error raised is set in the return value here
	f := Try(os.Open(fname)) // Try removes err part of return, and panics with special wrapper if err != nil
	defer f.Close()
	content = Try(io.ReadAll(f))
	return
}

func readFileWrapErr(fname string) (content []byte, err error) {
	defer Handle(func(e error) {
		err = fmt.Errorf("Error reading %s: %w", fname, e)
	})
	content = Try(os.ReadFile(fname))
	return
}

func testPanic() (err error) {
	defer Catch(&err)
	var empty []string
	empty[0] = "hello"
	return
}

func synthetic0(err error) error { return err }

func try0Failure() (err error) {
	defer Catch(&err)
	Try0(synthetic0(fmt.Errorf("simple error")))
	return
}

func try9Success() (p1, p2, p3, p4, p5, p6, p7, p8, p9 int, err error) {
	return 1, 2, 3, 4, 5, 6, 7, 8, 9, nil
}

func TestError(t *testing.T) {
	_, err := readFileTest("nosuchfile")
	switch runtime.GOOS {
	case "linux", "unix":
		assert.ErrorContains(t, err, "open nosuchfile: no such file or directory")
	case "windows":
		assert.ErrorContains(t, err, "open nosuchfile: The system cannot find the file specified.")
	default:
		assert.ErrorContains(t, err, "open nosuchfile:")
	}
}

func TestWrappedError(t *testing.T) {
	_, err := readFileWrapErr("nosuchfile")
	switch runtime.GOOS {
	case "linux", "unix":
		assert.ErrorContains(t, err, "Error reading nosuchfile: open nosuchfile: no such file or directory")
	case "windows":
		assert.ErrorContains(t, err, "Error reading nosuchfile: open nosuchfile: The system cannot find the file specified.")
	default:
		assert.ErrorContains(t, err, "Error reading nosuchfile: open nosuchfile: ")
	}
}

func TestSuccess(t *testing.T) {
	content, err := readFileTest("handler_test.go")
	assert.NoError(t, err)
	assert.NotEmpty(t, content)
}

func TestTry0Error(t *testing.T) {
	err := try0Failure()
	assert.EqualError(t, err, "simple error")
}

func TestTry9Error(t *testing.T) {
	p1, p2, p3, p4, p5, p6, p7, p8, p9 := Try9(try9Success())
	assert.Equal(t, 1, p1)
	assert.Equal(t, 2, p2)
	assert.Equal(t, 3, p3)
	assert.Equal(t, 4, p4)
	assert.Equal(t, 5, p5)
	assert.Equal(t, 6, p6)
	assert.Equal(t, 7, p7)
	assert.Equal(t, 8, p8)
	assert.Equal(t, 9, p9)
}

// Check that a panic raised in a function that defers to Catch
// will just panic as expected, with a somewhat intelligible call
// trace
func TestActualPanic(t *testing.T) {
	defer func() {
		if pnk := recover(); pnk == nil {
			assert.Fail(t, "function did not panic")
		} else {
			pnkText := fmt.Sprintf("%v", pnk)
			assert.Contains(t, pnkText, "index out of range")
			var callers []uintptr = make([]uintptr, 30)
			ncallers := runtime.Callers(0, callers)
			callers = callers[0:ncallers]
			frames := runtime.CallersFrames(callers)
			foundPanicSite := false
			for frame, more := frames.Next(); more; frame, more = frames.Next() {
				if frame.Function == "github.com/robdavid/genutil-go/errors/handler.testPanic" {
					foundPanicSite = true
					break
				}
			}
			assert.True(t, foundPanicSite, "couldn't not find original panic site in stack trace")
		}
	}()
	testPanic()
}

func recurseErrorCatch(depth int) (err error) {
	defer Catch(&err)
	if depth > 0 {
		Check(recurseErrorCatch(depth - 1))
	} else {
		err = fmt.Errorf("Hit bottom")
	}
	return
}

func recurseNoErrorCatch(depth int) (total int, err error) {
	defer Catch(&err)
	if depth > 0 {
		total = 1 + Try(recurseNoErrorCatch(depth-1))
	}
	return
}

func decrement(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("Will not decrement below zero")
	} else {
		return n - 1, nil
	}
}

func loopNoErrorCatch(iterations int) (count int, err error) {
	defer Catch(&err)
	for iterations > 0 {
		iterations = Try(decrement(iterations))
		count++
	}
	return
}

func loopNoErrorNoCatch(iterations int) (count int, err error) {
	for iterations > 0 {
		iterations = Try(decrement(iterations))
		count++
	}
	return
}

func loopNoError(iterations int) (count int, err error) {
	for iterations > 0 {
		if iterations, err = decrement(iterations); err != nil {
			return
		} else {
			count++
		}
	}
	return
}

func recurseErrorCatchOnce(depth int, maxDepth int) (err error) {
	if depth == 0 {
		defer Catch(&err)
	}
	if depth < maxDepth {
		recurseErrorCatchOnce(depth+1, maxDepth)
	} else {
		Raise(fmt.Errorf("Hit bottom"))
	}
	return
}

func recurseErrorReturn(depth int) (err error) {
	if depth > 0 {
		return recurseErrorReturn(depth - 1)
	} else {
		return fmt.Errorf("Hit bottom")
	}
}

// BenchmarkRewindTime times how long it takes to unwind a 1000
// deep stack of recursive functions with each frame calling
// the next via Check() and having a handler like:
//
//	defer Catch(&err)
//
// when the deepest call raises an error.
func BenchmarkRewindTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := recurseErrorCatch(1000)
		assert.Error(b, err, "Hit bottom")
	}
}

// BenchmarkNoErrorReturnTime times how long it takes to return from
// a 1000 deep stack of recursive functions with each frame calling
// the next via Try() and having a handler like:
//
//	defer Catch(&err)
//
// when the deepest call returns a value without an error.
func BenchmarkNoErrorReturnTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		total, err := recurseNoErrorCatch(1000)
		assert.Equal(b, 1000, total)
		assert.NoError(b, err)
	}
}

// BenchmarkCatchOnce times how long it takes to unwind a 1000
// deep stack of recursive functions with only the first frame having
// a handler like:
//
//	defer Catch(&err)
//
// when the deepest call raises an error.
func BenchmarkCatchOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := recurseErrorCatchOnce(0, 1000)
		assert.EqualError(b, err, "Hit bottom")
	}
}

// BenchmarkReturnTime times how long it takes to unwind a 1000
// deep stack of recursive functions which employs traditional
// error handling, when the deepest call returns an error.
func BenchmarkReturnTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		err := recurseErrorReturn(1000)
		assert.EqualError(b, err, "Hit bottom")
	}
}

// BenchmarkNoErrorCatch times how long it takes to process
// a loop of 1000 successful Try() calls in a function with
// a deferred error handler.
func BenchmarkNoErrorCatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		count, err := loopNoErrorCatch(1000)
		require.Equal(b, count, 1000)
		require.Nil(b, err)
	}
}

// BenchmarkNoErrorNoCatch times how long it takes to process
// a loop of 1000 successful Try() calls in a function with
// no error handler.
func BenchmarkNoErrorNoCatch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		count, err := loopNoErrorNoCatch(1000)
		require.Equal(b, count, 1000)
		require.Nil(b, err)
	}
}

// BenchmarkNoError times how long it takes to process
// a loop of 1000 successful function calls using
// traditional error testing only.
func BenchmarkNoError(b *testing.B) {
	for i := 0; i < b.N; i++ {
		count, err := loopNoError(1000)
		require.Equal(b, count, 1000)
		require.Nil(b, err)
	}
}

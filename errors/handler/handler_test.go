package handler

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
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

func recurseErrorCatchOnce(depth int, maxDepth int) (err error) {
	if depth == 0 {
		defer Catch(&err)
	}
	if depth < maxDepth {
		err = recurseErrorCatchOnce(depth+1, maxDepth)
	} else {
		err = fmt.Errorf("Hit bottom")
		Check(err)
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

func BenchmarkRewindTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recurseErrorCatch(1000)
		//assert.EqualError(b, err, "Hit bottom")
	}
}

func BenchmarkCatchOnce(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recurseErrorCatchOnce(0, 1000)
		//assert.EqualError(b, err, "Hit bottom")
	}
}

func BenchmarkReturnTime(b *testing.B) {
	for i := 0; i < b.N; i++ {
		recurseErrorReturn(1000)
		//assert.EqualError(b, err, "Hit bottom")
	}
}

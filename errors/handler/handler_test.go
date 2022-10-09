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

func TestError(t *testing.T) {
	_, err := readFileTest("nosuchfile")
	assert.Error(t, err)
}

func TestWrappedError(t *testing.T) {
	_, err := readFileWrapErr("nosuchfile")
	assert.ErrorContains(t, err, "Error reading nosuchfile: open nosuchfile: ")
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

// Some error handling functions and types to more ergonomically assist
// with writing tests against functions that may return errors
package test

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/errors/result"
	"github.com/robdavid/genutil-go/tuple"
)

// An interface implemented by multiple types in the "testing" package
type TestReporting interface {
	Error(args ...any)
	Errorf(format string, args ...any)
	FailNow()
	Helper()
}

// A wrapper around result.Result that supports test assertions.
type TestableResult[T any] struct {
	result.Result[T]
}

func resultFrom[T any](value T, err error) TestableResult[T] {
	return TestableResult[T]{result.From(value, err)}
}

// Creates a TestableResult from a an error only,
// e.g.
//
//	r := test.Result0(os.Rename(oldfile,newfile))
func Result0(err error) *TestableResult[tuple.Tuple0] {
	return Result(tuple.Of0(), err)
}

// Creates a TestableResult from a an error only. Alias for Result0.
// e.g.
//
//	r := test.Status(os.Rename(oldfile,newfile))
func Status(err error) *TestableResult[tuple.Tuple0] {
	return Result0(err)
}

// Creates a TestableResult from a return value and an error
// e.g.
//
//	r := test.Result(os.Open(myfile))
func Result[T any](value T, err error) *TestableResult[T] {
	return &TestableResult[T]{result.From(value, err)}
}

// Check checks if a given error value is nil. If not, report the
// error and fail the test.
func Check(t TestReporting, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}

// Must either returns the value of a TestableResult, or if
// an error has occurred, report it and fail the test.
func (r *TestableResult[T]) Must(t TestReporting) T {
	t.Helper()
	if r.IsError() {
		Check(t, r.GetErr())
	}
	return r.Get()
}

// Fails expects an error; it marks the test as failed if the result is not an error
func (r *TestableResult[T]) Fails(t TestReporting) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	}
	return r
}

// FailsWith expects a specific error; it marks the test as failed if the result is
// not, or does not wrap, the expected error.
func (r *TestableResult[T]) FailsWith(t TestReporting, expected error) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	} else if !errors.Is(r.GetErr(), expected) {
		t.Error(fmt.Errorf("expected error '%s', but got '%s'", expected, r.GetErr()))
	}
	return r
}

// Expects a specific error; marks the test as failed if the result is not an
// error whose Error() return contains the string provided in expected.
func (r *TestableResult[T]) FailsContaining(t TestReporting, expected string) *TestableResult[T] {
	t.Helper()
	if !r.IsError() {
		t.Error(fmt.Errorf("an error was expected, but did not occur"))
	} else if !strings.Contains(r.GetErr().Error(), expected) {
		t.Error(fmt.Errorf("expected error to contain '%s', but was '%s'", expected, r.GetErr().Error()))
	}
	return r
}

type stackFrame struct {
	file     string
	line     int
	function string
}

func (frame *stackFrame) String() string {
	return fmt.Sprintf("%s:%d %s", frame.file, frame.line, frame.function)
}

func (frame *stackFrame) packageName() string {
	packageLast := strings.LastIndex(frame.function, "/")
	if packageLast < 0 {
		return frame.function
	}
	functionStart := strings.Index(frame.function[packageLast:], ".")
	if functionStart < 0 {
		return frame.function
	}
	return frame.function[:packageLast+functionStart]
}

func callStack(size int) []stackFrame {
	callers := make([]uintptr, size)
	n := runtime.Callers(1, callers)
	frames := runtime.CallersFrames(callers[:n])
	stackFrames := make([]stackFrame, 0, n)
	for frame, _ := frames.Next(); frame.PC != 0; frame, _ = frames.Next() {
		stackFrames = append(stackFrames, stackFrame{
			file:     frame.File,
			line:     frame.Line,
			function: frame.Function,
		})
	}
	return stackFrames
}

func trimCallStack(frames []stackFrame) []stackFrame {
	const (
		lookForTry = iota
		lookForEndOfHandling
	)
	lookFor := lookForTry
	tryFuntionPath := handler.PackageName() + ".Try"
	for i, frame := range frames {
		switch lookFor {
		case lookForTry:
			// Look for Try frame
			if strings.HasPrefix(frame.function, tryFuntionPath) {
				lookFor = lookForEndOfHandling
			}
		case lookForEndOfHandling:
			// Look for end of calls in handler or result package
			pkg := frame.packageName()
			if pkg != handler.PackageName() &&
				pkg != result.PackageName() {
				return frames[i:]
			}
		}
	}
	return frames
}

const maxStackDepth = 255

func logStack(t TestReporting) {
	frames := callStack(maxStackDepth)
	frames = trimCallStack(frames)
	var trace strings.Builder
	for _, frame := range frames {
		fmt.Fprintln(&trace, frame.String())
	}
	t.Errorf("stack trace\n%s\n", trace.String())
}

// Reports any error encountered in a call to Try().
// Should be used as part of a defer call.
// e.g.
//
//	defer ReportErr(t)
//	f := Try(os.Open(myfile))
func ReportErr(t TestReporting) {
	t.Helper()
	logStack(t)
	if err := recover(); err != nil {
		if tryErr, ok := err.(handler.TryError); ok {
			t.Error(tryErr.Error)
		} else {
			panic(err)
		}
	}
}

//go:generate code-template test.tmpl

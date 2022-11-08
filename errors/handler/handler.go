// A Try/Handle/Catch error handling mechanism implemented
// using generics, somewhat similar to previous proposals to
// extend the language with error handling such as
// https://github.com/golang/go/issues/32437
// This package is intended to be imported unqualified, e.g.
//
//	import . "github.com/robdavid/genutil-go/errors/handler"
package handler

// Removes the error value from a return value.
// Panics if err is non-nil, otherwise returns the value
// e.g.
//
//	f := Must(os.Open("myfile"))
func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

type tryError struct {
	err error
}

// Removes the error component of a function's return value. If there is
// no error, the non-error value is returned. Otherwise
// the function panics with a value wrapping the error.
// The error can be recovered and returned by your function in conjunction
// with defer and the Catch or Handle functions.
// e.g.
//
//	f := Try(os.Open("myfile"))
func Try[T any](t T, err error) T {
	if err != nil {
		panic(tryError{err})
	}
	return t
}

// An alias for Try
func Try1[T any](t T, err error) T { return Try(t, err) }

// Zero argument variant of Try
func Try0(err error) {
	if err != nil {
		panic(tryError{err})
	}
}

// An alias for Try0
func Check(err error) { Try0(err) }

// Two non-error argument variant of Try()
func Try2[A any, B any](a A, b B, err error) (A, B) {
	if err != nil {
		panic(tryError{err})
	}
	return a, b
}

// Three non-error argument variant of Try()
func Try3[A any, B any, C any](a A, b B, c C, err error) (A, B, C) {
	if err != nil {
		panic(tryError{err})
	}
	return a, b, c
}

// A function that will recover a panic created by Try. This
// should be called in a defer statement prior to calls to Try.
// The parameter should be a pointer to the calling function's
// error return value which will be set to the error intercepted by
// Try.
//
// e.g.
//
//	  func readFileTest(fname string) (content []byte, err error) {
//		   defer Catch(&err)
//		   f := Try(os.Open(fname))
//		   defer f.Close()
//		   content = Try(io.ReadAll(f))
//		   return
//	  }
func Catch(retErr *error) {
	if err := recover(); err != nil {
		if tryErr, ok := err.(tryError); ok {
			*retErr = tryErr.err
		} else {
			panic(err)
		}
	}
}

// A function that will recover a panic created by Try. This
// should be called in a defer statement prior to calls to Try.
// The handler parameter is an error handling function that can
// be used to place error handling in one place in your function,
// such as wrapping the error in another error type.
// e.g.
//
//	  func readFileWrapErr(fname string) (content []byte, err error) {
//		   defer Handle(func(e error) {
//			 err = fmt.Errorf("Error reading %s: %w", fname, e)
//		   })
//		   content = Try(os.ReadFile(fname))
//		   return
//	  }
func Handle(handler func(err error)) {
	if err := recover(); err != nil {
		if tryErr, ok := err.(tryError); ok {
			handler(tryErr.err)
		} else {
			panic(err)
		}
	}
}

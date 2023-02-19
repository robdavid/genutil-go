# genutil-go - A generics utility library for Go
---
# DOCUMENTATION WORK IN PROGRES

---

A library of utility functions made possible by Go generics, providing features missing from the standard libraries. This library is still in its early stages, and breaking changes are still possible. Additional functionality is likely to be added.

The library falls into a number of categories, subdivided into separate packages.
- [Errors](#errors)
  - [Handler](#handler)
    - [Example](#example)


## Errors
The `errors` package has a number of subpackages related to error handling.

### Handler
The `errors.handler` package provides a way to handle errors in Go more ergonomically, at the expense of less efficient runtime handling of error cases. It works by turning errors into `panic`s which can later be recovered easily. This is less efficient that the traditional Go form of manual checking and returning of error values. If the error condition is likely to be anything other than an infrequent occurance, the traditional method is more appropriate. That said, it is likely no less efficient than languages that provide exception handling.

The advantage is that it can condense repetitive multi line error checking boilerplate to a single call to a `Try` function.

#### Example

Consider the following function to read the contents of a file.
```go
func readFileContent(fname string) (content []byte, err error){
  var f *os.File
  if f,err  = os.Open(fname); err != nil {
    return
  }
  defer f.Close()
  if content, err = io.ReadAll(f); err != nil {
    return
  }
  return
}
```
This is using typical Go error handling, testing each error value explicitly. The following version uses the `Try` and `Catch` error handling functions.
```go
import . "github.com/robdavid/genutil-go/errors/handler"

func readFileContent(fname string) (content []byte, err error) {
  defer Catch(&err)
  f := Try(os.Open(fname))
  defer f.Close()
  content = Try(io.ReadAll(f))
  return
}
```

Here we use `Try` to wrap function calls that return a value plus an error. `Try` removes the error part and returns just the single value. However, if the error is non-nil it will panic with a `TryError`. The `Catch` deferred function is used to recover from any panic created by `Try` and in this example will populate the `err` variable with the error that triggered the panic. This becomes our error return, as it is a named return value.

If you want to augment the error, or perform other processing on the error, the `Handle` function can be used.
```go
import . "github.com/robdavid/genutil-go/errors/handler"

func readFileContent(fname string) (content []byte, err error) {
  defer Handle(func(e error) {
    err = fmt.Errorf("%w: whilst opening %s", e, fname)
  })
  f := Try(os.Open(fname))
  defer f.Close()
  content = Try(io.ReadAll(f))
  return
}
```
See the [API Documentation](https://pkg.go.dev/github.com/robdavid/genutil-go/errors/handler) for more details.




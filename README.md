# genutil-go - A generics utility library for Go

---

# DOCUMENTATION WORK IN PROGRESS

---

A library of utility functions made possible by Go generics, providing features missing from the standard libraries. This library is still in its early stages, and breaking changes are still possible. Additional functionality is likely to be added.

See the [API Documentation](https://pkg.go.dev/github.com/robdavid/genutil-go/errors/handler) for more details.

The library falls into a number of categories, subdivided into separate packages.

- [Tuple](#tuple)
- [Errors](#errors)
  - [Handler](#handler)
    - [Example](#example)
  - [Result](#result)

## Tuple

Tuple generic types exist for tuples of size 0 to 9.

Tuples can be created by one of the `Of` constructor methods. Eg. `Of2` constructs a tuple of 2 elements.  

A tuple of size 0 has no type parameters, has only one value, and is also known as a unit.

```go
  t0 := tuple.Of0()
  u := tuple.Unit()
  fmt.Println(t0 == u) // true
```

For tuple sizes greater than zero, the generic type of the elements are inferred by the constructor. Each element in the tuple can be references by members named `First`, `Second`, `Third` etc.

```go
  t1 := tuple.Of1(3.141)
  fmt.Printf("%f %T\n", t1.First, t1) // 3.141000 tuple.Tuple1[float64]
```

A tuple of size 2 is also known as a Pair.

```go
  t2 := tuple.Of2(1, "one")
  p := tuple.Pair(1, "one")
  fmt.Println(t2 == p) // true
```

All tuple references implement a general `Tuple` interface.

```go
type Tuple interface {
 // Get the nth element of the tuple
 Get(int) any
 // Return the number of elements in the tuple
 Size() int
 // Return the tuple as a string, formatted (e1,e2,...)
 String() string
 // Tuple of first size-1 elements
 Pre() Tuple
 // Return the last element in the tuple
 Last() any
}
```

For example

```go
  t3 := tuple.Of3(1, "two", 3.1)
  t2 := tuple.Of2(1, "two")

  fmt.Printf("%d, %#v, %#v\n", t3.Size(), t3.Get(1), t3.Last()) // 3, "two", 3.1

  pre := t3.Pre().(*tuple.Tuple2[int, string])

  fmt.Println(t2 == *pre) // true
  fmt.Println(pre.String()) // (1,two)
  fmt.Println(&t3) // (1,"two",3.1)
  fmt.Println(t3) // {1 two 3.1}
```

## Errors

The `errors` package has a number of subpackages related to error handling.

### Handler

The `errors.handler` package provides a way to handle errors in Go more ergonomically, at the expense of less efficient runtime handling of error cases. It works by turning errors into `panic`s which can later be recovered easily. This is less efficient that the traditional Go form of manual checking and returning of error values. If the error condition is likely to be anything other than an infrequent occurrence, the traditional method is more appropriate. That said, it is likely no less efficient than languages that have native try/catch exception handling.

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

Here the `Try` function is used to wrap the function calls that return a value plus an error. `Try` removes the error part and returns just the single value. However, if the error is non-nil it will panic with a `TryError` value, wrapping the error. The `Catch` deferred function will recover from this type of panic and in this example will populate the `err` variable with the wrapped error which will be returned to the caller.

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

### Result

The `errors.result` package defines a `result.Result` type that contains a value plus an error, typically used to represent the return value of a function, including its error component. It has convenience methods for constructing an instance from a function return, e.g.

```go
import "github.com/robdavid/genutil-go/errors/result"
r := result.From(os.Open(file))
```

It has `Get` and `GetError` methods to get the value part and error part of the result, either or both of which may be present.

```go
if (r.GetErr() != nil) {
  return nil, r.GetErr()
}
return io.ReadAll(r.Get())

```

The `Result` type also supports a `Try` method similar to the `Try` method in error handler package. This method transforms the result instance to the underlying value only, if the error is nil. Otherwise, if the error is non-nil, the function creates a panic that can be handled using the error handling methods of that package.

```go
import (
  . "github.com/robdavid/genutil-go/errors/handler"
  "github.com/robdavid/genutil-go/errors/result"
  "fmt"
  "os"
)

func openFile(file string) result.Result[*os.File] {
 return result.From(os.Open(file))
}

func printFile(file string) (err error) {
 defer Catch(&err)
 f := openFile(file)                          // Returns result.Result[*os.File]
 fmt.Printf("%s\n", Try(io.ReadAll(f.Try()))) // Call Try on result f 
 return nil
}
```

Results that contain more than one value are covered by the variants of `result.Result`; `result.Result2`, `result.Result3` etc. Each of these hold a `tuple.Tuple` value of the appropriate size.

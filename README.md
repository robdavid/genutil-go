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
  - [Test](#test)
- [Iterator](#iterator)
  - [Constructing iterators](#constructing-iterators)
    - [Slices](#slices)
    - [Ranges](#ranges)
  - [Simple Iterator](#simple-iterator)

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

The `errors.handler` package provides a way to handle errors in Go more ergonomically, at the potential expense of less efficient runtime handling when error cases do occur. It is thus most suitable for use cases where errors are expected to occur infrequently.

It works by removing the error component from a function call's return values, converting it to a `panic` if it is non-nill. This panic can later be recovered easily. Repetitive multi line error checking boilerplate can be condensed to a single call to a `Try` function.

#### Example

Consider the following function to read the contents of a file. This is using typical Go error handling patterns, with explicit testing for non-nil error values.

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

The following version instead uses the `Try` and `Catch` error handling functions.

```go
import . "github.com/robdavid/genutil-go/errors/handler"

func readFileContent(fname string) (content []byte, err error) {
  defer Catch(&err) // Any panic raised by Try is recovered here
  f := Try(os.Open(fname)) // Panics if the error is non-nil
  defer f.Close()
  content = Try(io.ReadAll(f))
  return
}
```

Here the `Try` generic function is used to strip the error part from the io function returns, leaving just a simple value. However, if the error is non-nil it will panic with a `TryError` value, wrapping the error. The `Catch` deferred function will recover from this type of panic and in this example will populate the `err` return value with the original error, thus causing it to be returned to the caller of our function.

If you want to augment the error, or perform other processing on the error, the `Handle` deferred function can be used instead of `Catch`.

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

The `Result` type also supports a `Try` method similar to the `Try` method in error handler package. This method transforms the result instance to the underlying value only, if the error is nil. Otherwise, if the error is non-nil, the function creates a panic that can be handled using the error handling package's error handling methods, such as `Catch` or `Handle` .

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

Results that contain more than one value are covered by the variants of `result.Result`; `result.Result2`, `result.Result3` etc. Each of these hold a `tuple.Tuple2` or `tuple.Tuple3` etc. value respectively. There is also a `result.Result0` type for results that consist of an error only.

### Test

The `test` package contains some error reporting methods to help with unit tests that need to assert that an error should or should not occur. It builds on top of `result.Result` to create a `test.TestableResult` type that can assert against and report errors in a test.

```go
import (
  "github.com/robdavid/genutil-go/errors/test"
  "testing"
  "os"
)

func TestOpen(t *testing.T) {
  f := test.Result(os.Open("myfile")).Must(t)
  // Test assertions
}

```

The above builds a `test.TestResult` value from the return value of the call to `os.Open`. It then calls a method `Must` that asserts the result must have a nil error. If it is non-nil, the error is reported to the test framework, and the test is terminated.

Various other methods and types exist to handle return values with errors only or multiple non-error values, such as `test.Result0` and `test.Result2`.

## Iterator

An `Iterator` is a generic type equivalent to the following definition

```go
type Iterator[T any] interface {
  // Set the iterator's current value to be the first, and subsequent, iterator elements.
  // False is returned when there are no more elements (the current value remains unchanged)
  Next() bool
  // Get the current iterator value.
  Value() T
  // Stop the iterator; subsequent calls to Next() will return false.
  Abort()
  // Size estimate, where possible, of the number of elements remaining.
  Size() IteratorSize
  // Return iterator as a channel.
  Chan() <-chan T
}
```

Iterators can be consumed in a `for` loop in two ways. The first is to use `Next()` and `Value()`.

```go
  var iter iterator.Iterator[int] // Iterator of integers
  // instantiate iterator
  for iter.Next() {
    fmt.Sprintf("%d",iter.Value())
  }
```

The other is to range over the channel that the iterator provides. Each element in the iterator is sent over the channel in sequence, and closed when the iterator has no more elements.

```go
  var iter iterator.Iterator[int] // Iterator of integers
  // instantiate iterator
  for v := range iter.Chan() {
    fmt.Sprintf("%d",v)
  }
```

The `Abort` method can be used to stop the iterator; once called the `Next` method will return `false` and the channel (if used) will be closed. Eg.

```go
  for v := range iter.Chan() {
    fmt.Sprintf("%d",v)
    if v == 0 {
      iter.Abort() // loop will end after current iteration
    } 
  }
```

### Constructing iterators

Aside from just implementing the `Iterator` interface, there are a number of ways available for constructing iterators.

#### Slices

An iterator over a slice of values is easily created with the `Slice` function.

```go
input := []int{1, 2, 3, 4}
iter := iterator.Slice(input)
for iter.Next() {
  fmt.Sprintf("%d ",iter.Value()) // 1 2 3 4
}
```

An iterator can also be collected into a slice with `Collect`

```go
input := []int{1, 2, 3, 4}
iter := iterator.Slice(input)
output := iterator.Collect(iter) // output is equal to input
```

#### Ranges

An iterator over a range of scalar numeric values can be built using the `Range` function.

```go
iter := iterator.Range(1,5)
slice := iterator.Collect(iter) // []int{1,2,3,4}
```

Ranges can be built over any scalar numeric type, including float

```go
iter := iterator.Range(0.0, 5.0)
slice := iterator.Collect(iter) // []float64{0.0, 1.0, 2.0, 3.0, 4.0}
```

The `RangeBy` method creates a range with a specific increment.

```go
iter := iterator.RangeBy(0.0, 2.0, 0.5)
slice := iterator.Collect(iter) // []float64{0.0, 0.5, 1.0, 1.5, 2.0 }
```

The increment may be negative, in which case the `from` value must be less 
than the `upto` value.

```go
iter := iterator.RangeBy(5.0, 0.0, -0.5)
slice := iterator.Collect(iter) // []float64{5.0, 4.5, 4.0, 3.5, 3.0, 2.5, 2.0, 1.5, 1.0, 0.5}
```

### Simple Iterator

An iterator can be built via a simplified SimpleIterator interface.

```go
type SimpleIterator[T any] interface {
  // Next sets the iterator's current value to be the first, and subsequent, iterator elements.
  // False is returned only when there are no more elements (the current value remains unchanged)
  Next() bool
  // Value gets the current iterator value.
  Value() T
  // Abort stops the iterator; subsequent calls to Next() will return false.
  Abort()
}
```

Instances implementing this interface can be transformed to a full `Iterator[T]` via one of the utility methods.

To make an `Iterator[T]` of indeterminate size, use

```go
func MakeIteratorFromSimple[T any](base SimpleIterator[T]) Iterator[T]
```

or to make one with a given size use

```go
func MakeIteratorOfSizeFromSimple[T any](base SimpleIterator[T], size IteratorSize) Iterator[T]
```

The following example illustrates how an Iterator over a slice can be created by implementing only
the SimpleIterator interface. 

```go

// iterSlice is a SimpleIterator over a slice; it implements Next, Value and Abort methods
type iterSlice[T any] struct {
  slice []T
  index int
}

// Advance to first/next element
func (is *iterSlice[T]) Next() bool {
  is.index++
  return is.index < len(is.slice)
}

// Return current element
func (is *iterSlice[T]) Value() T {
  return is.slice[is.index]
}

// Move index after last element; ensures next Next() call returns false
func (is *iterSlice[T]) Abort() {
  is.index = len(is.slice)
}

// newIterSlice creates an Iterator over a slice
func newIterSlice[T any](slice []T) Iterator[T] {
  // First, create the SimpleIterator over the slice.
  // Index starts at -1 because Next() is called for the first element.
  simpleIter := &iterSlice[T]{slice, -1}
  // Then create an Iterator from the SimpleIterator, with known size (the slice's length)
  return iterator.MakeIteratorOfSizeFromSimple[T](simpleIter, iterator.NewSize(len(slice)))
}
```


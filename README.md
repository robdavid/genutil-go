# genutil-go - A generics utility library for Go

A library of utility functions made possible by Go generics, providing features previously missing from the standard libraries. This library is still in its early stages, and breaking changes are still possible. Additional functionality is likely to be added.

See the [API Documentation](https://pkg.go.dev/github.com/robdavid/genutil-go) for more details.

The library falls into a number of categories, subdivided into separate packages.

- [Tuple](#tuple)
- [Errors](#errors)
  - [Handler](#handler)
    - [Example](#example)
  - [Result](#result)
  - [Test](#test)
- [Iterator](#iterator)
  - [Constructing iterators](#constructing-iterators)
    - [Iterators over slices](#iterators-over-slices)
    - [Ranges](#ranges)
  - [Simple Iterator](#simple-iterator)
- [Maps](#maps)
  - [Keys, Values and Items](#keys-values-and-items)
    - [Iterators](#iterators)
    - [Sorted slices](#sorted-slices)
  - [Nested maps](#nested-maps)
    - [Inserting values](#inserting-values)
    - [Fetching values](#fetching-values)
    - [Deleting values](#deleting-values)
- [Slices](#slices)
  - [Functional primitives](#functional-primitives)
    - [Predicate functions](#predicate-functions)
    - [Transformations](#transformations)
  - [Range functions](#range-functions)
- [Option type](#option-type)
  - [Usage](#usage)
  - [Zero value](#zero-value)
  - [Comparisons](#comparisons)
  - [Mutations](#mutations)
  - [Marshalling and Unmarshalling](#marshalling-and-unmarshalling)
    - [Unmarshalling](#unmarshalling)

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

A number of generic iterator types are provided with some useful abilities such as filtering, mapping, enumeration, sizing information, and mutation of underlying values. They can be consumed via `for` loops, or via collected into slices or maps, converted to and from Go native `iter.Seq` and `iter.Seq2` types, or their elements sent from a goroutine over a channel.

The principle abstract iterator types are:

| Interface Type | Description |
|------|-------------|
| `iterator.Iterator[T]` | An immutable iterator over elements of type T |
| `iterator.MutableIterator[T]` | A mutable iterator over elements of type T that supports methods for modifying or deleting the current value. This type supports methods to modify the element in place, or to remove it from an underlying collection. |
| `iterator.Iterator2[K,V]` | An immutable iterator over pairs of element index or key of type K, and element value of type V.|
| `iterator.ImmutableIterator2[K,V]` | A mutable iterator over pairs of element index or key of type K, and element value of type V. This type supports methods to modify the value of an element in place, or to remove it (and the key where applicable) from an underlying collection. Note there is no method to modify the key value. |

The following sections give you an overview of the iterator types and their usage. For full documentation see the API reference.

### Creation

Iterators can be created in various ways.

#### From values

Most simply, an iterator can be created from an explicit list of values:

```go
intIter := iterator.Of(1,1,2,3,5,8)          // iterator.Iterator[int]
strIter := iterator.Of("red","green","blue") // iterator.Iterator[string]
```

#### From numeric value ranges

Iterators can be created as a range of numeric values:

```go
intIter := iterator.Range(0,5)      // iterator.Iterator[int]
                                    // 0,1,2,3,4
fltIter := iterator.Range(0.0, 5.0) // iterator.Iterator[float64]
                                    // 0.0, 1.0, 2.0, 3.0, 4.0
```

#### From slices

An iterator can be created over a slice. Such an iterator carries sizing information:

```go
myslice := []int { 1, 2, 3, 4 }
intIter := slices.Iter(myslice) // iterator.Iterator[int]
intIter.Size().IsKnownToBe(4)   // true
```

In addition, a mutable iterator can be created, which allows modification of the underlying slice (see further below).

```go
myslice := []int { 1, 2, 3, 4 }
intIter := slices.IterMut(myslice) // iterator.MutableIterator[int]
```

#### From native Go iterators

An iterator can be created from a native Go iterator. This underlying iterator can be accessed directly.

```go
	// fib returns a native Go iterator (fibonacci sequence).
	fib := func(yield func(int) bool) {
		tail := [2]int{0, 1}
		for {
			if !yield(tail[1]) {
				return
			}
			tail[0], tail[1] = tail[1], tail[0]+tail[1]
		}
	}

	fibItr := iterator.New(fib) // iterator.Iterator[int]
	fibSeq := fibItr.Seq()      // iter.Seq[int]
```

#### From maps

An iterator can be constructed over both keys and values of a map:

```go
m := map[int]string{ 1: "one", 2: "two", 3: "three" }
itr := maps.Iter(m) // Iterator2[int,string]
```

It's also possible to create a mutable iterator, that supports modification of the underlying map (see further below)

```go
m := map[int]string{ 1: "one", 2: "two", 3: "three" }
itr := maps.IterMut(m) // MutableIterator2[int,string]
```

### Consumption

Iterators can be consumed in a number of ways. 

#### For loops

The most straight forward way to consume an iterator is using a `for` loop. The recommended way is by converting to a native Go iterator with the `Seq` method.

```go
for n := range iterator.Range(0,10).Seq() {
  fmt.Printf("%d\n",n) // 0..1..2..
}
```

It's also possible to use `Next()` and `Value()` methods provided by iterators in a `for` loop as follows:

```go
for itr := iterator.Range(0, 10); itr.Next(); {
  fmt.Fprintf(&buffer, "%d\n", itr.Value()) // 0..1..2.. .. ..9
}
```

Generally speaking, the `Seq()` method is preferred since using `Next()` against an iterator that is backed by an `iter.Seq` native iterator incurs a performance penalty (due to use of `iter.Pull`). 

#### Collection

Iterators have a `Collect()` method that allows elements to be collected in to a slice.

```go
c := iterator.Range(0,5).Collect() // []int{0, 1, 2 ,3, 4}
```

Iterators of element pairs can be collected into a map, provided the key value is comparable.

```go
m := iterator.CollectMap(iterator.Of("zero", "one","two","three").Enumerate()) // map[int]string{0, "zero", 1: "one", 2: "two", 3: "three"}
```

The `Enumerate()` method turns an `iterator.Iterator[T]` into an `iterator.Iterator2[int,T]` by adding a counter key starting at zero. The `iterator.CollectMap` function is a function rather than a method because the comparable constraint needs to be enforced, which cannot be done in a method.

### Mutability

Some iterators support the mutation of the underlying collection from which their elements are drawn. Out of the box, an `iterator.MutableIterator` can be constructed over slices, and an `iterator.MutableIterator2` can be constructed over maps. Both iterators have a `Set(v T)` method which provides for the mutation of the current element (not the key), and a `Delete()` method which removes the current element (or element pair) from the collection.

#### Mutability over slices

In order to support mutability over slices, especially the removal of elements, the iterator needs to operate on a pointer to a slice; the removal of an element may lead to reallocation of the slice to a new location. The following example builds a slice of ints from 0...9 (inclusive), and iterates mutably over it, deleting elements that are odd, whilst dividing even numbers by 2.

```go
s := slices.Range(0,10)
itr := slices.IterMut(&s)
for n := range itr.Seq() {
  if n%2 == 1 {
    itr.Delete()
  } else {
    itr.Set(n/2)
  }
}
fmt.Println(s) // [0 1 2 3 4]
```

#### Mutability over maps



```go
// Make a map
m := make(map[int]int)
for i := range 10 {
  m[i] = i + 10
}

// Iterate
itr := maps.IterMut(m)
for k, v := range itr.Seq2() {
  if k%2 == 1 {
    itr.Delete()
  } else {
    itr.Set(v / 2)
  }
}
fmt.Println(m) // map[0:5 2:6 4:7 6:8 8:9]
```


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

#### Iterators over slices

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

The `RangeBy` method creates a range with a given increment.

```go
iter := iterator.RangeBy(0.0, 2.0, 0.5)
slice := iterator.Collect(iter) // []float64{0.0, 0.5, 1.0, 1.5, 2.0 }
```

The increment may be negative, in which case the `from` value must be less than the `upto` value.

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
func NewFromSimple[T any](base SimpleIterator[T]) Iterator[T]
```

or to make one with a given size use

```go
func NewFromSimpleWithSize[T any](base SimpleIterator[T], size func() IteratorSize) Iterator[T]
```

The following example illustrates how an Iterator over a slice can be created by implementing only
the SimpleIterator interface.

```go

// iterSlice is a SimpleIterator over a slice; it implements Next, Value and Abort methods
type iterSlice[T any] struct {
  slice []T
  index int
  value T
}

// Advance to first/next element
func (is *iterSlice[T]) Next() bool {
  if is.index < len(is.slice) {
    is.value = is.slice[is.index]
    is.index++
    return true
  } else {
    return false
  }
}

// Return current element
func (is *iterSlice[T]) Value() T {
  return is.slice.value
}

// Move index after last element; ensures next Next() call returns false
func (is *iterSlice[T]) Abort() {
  is.index = len(is.slice)
}

// Reset index to the beginning
func (is *iterSlice[T]) Reset() {
  is.index = 0
}


// newIterSlice creates an Iterator over a slice
func newIterSlice[T any](slice []T) Iterator[T] {
  // First, create the SimpleIterator over the slice.
  simpleIter := &iterSlice[T]{slice, 0}
  // Then create an Iterator from the SimpleIterator, with known size (the slice's length)
  return iterator.NewFromSimpleWithSize[T](
    simpleIter, 
    func() iterator.IteratorSize { 
      return iterator.NewSize(len(simpleIter.slice)-simpleIter.index) 
    },
  )
}
```

## Maps

The `maps` package contains a number of utility functions that work over maps, including getting a slice of the Keys or Values of a map.

### Keys, Values and Items

The `Keys` function can be used to collect the keys of a map into a slice, e.g:

```go
m := map[string]int{"one": 1, "two": 2}
k := maps.Keys(m) // []string{"one","two"}
```

Similarly, the `Values` function will collect the values:

```go
v := maps.Values(m) // []int{1,2}
```

If you need both keys and values, the `Items` function will return a slice of `tuple.Tuple2` values with each tuple holding a key/value pair, e.g:

```go
i := maps.Items(m) // []tuple.Tuple2[string,int] { {"one",1}, {"two",2} }
```

Note that in all three cases, the ordering of the slice returned is undefined.

#### Iterators

For each of the these three functions, there exists three variants,  `IterKeys`, `IterValues` and `IterItems`, which return iterators rather than slices. The ordering for these iterators is also undefined.

#### Sorted slices

The slice returning functions also have ordered alternatives, `SortedKeys`, `SortedValuesByKey` and `SortedItems`, which return keys, values and items sorted in key order.

### Nested maps

A group of functions are available for managing nested maps, that is maps whose values may also be maps, and which have the signature `map[K comparable]any`. A common concrete example is `map[string]any`, which is useful for un-marshaling arbitrary YAML or JSON documents.

All the functions take a map with the generic signature above, and a list of elements of type K which represent a path into the map. For example a list consisting of `[]string{"a","b","c"}` describes the value found by first looking up "a" in a map with `string` keys, expecting to find another map value of the same type, then looking up "b" in that map, again expecting a map result, and then finally looking up "c" in that final map.

#### Inserting values

The `PutPath` function will insert or mutate a key in the map. Any missing intermediate levels of map will be created as necessary, except the top level; the map provided cannot be nil. For example:

```go
m := make(map[string]any)
maps.PutPath(m, []string{"a", "b"}, 123)
// m is map[string]any{"a": map[string]any{"b": 123}}
```

Once an item has been established as either a map or non-map value, it cannot be replaced by a value of the opposite kind, for example:

```go
err := maps.PutPath(m, []string{"a"}, 456)
errors.Is(err,maps.ErrPathConflict) // true
```

#### Fetching values

The `GetPath` function will fetch a value at a location in the nested map, defined by a slice of keys. It returns the value found and an error.

```go
m := map[string]any {"a": map[string]any { "b": 123 }}
v, _ := maps.GetPath(m, []string{"a","b"}) // v == 123
```

If the specified path does not exist, then a `maps.ErrKeyError` error will be returned.

```go
_, err := maps.GetPath(m, []string{"a","c"})     // errors.Is(err,maps.ErrKeyError)
_, err := maps.GetPath(m, []string{"a","b","c"}) // errors.Is(err,maps.ErrKeyError)
```

#### Deleting values

The `DeletePath` function will delete an item from a nested map, located by a path consisting of a slice of keys. It can delete a leaf value or an interior map, thereby removing a subtree. If a map becomes empty as a result of deleting a key from it, it itself is deleted from the parent map. This process recurses towards the root of the tree as many times as necessary.

```go
m := map[string]any{
  "one": 1,
  "two": map[string]any{
    "three": 23,
  },
}
maps.DeletePath(m,[]string{"two","three"}) // m == map[string]any{"one": 1 }
```

## Slices


A variety of functions that work over slices are included in the `slices` package. Some examples are covered here. See the [documentation](https://pkg.go.dev/github.com/robdavid/genutil-go/slices) for full details.

### Functional primitives

Some "functional style" operations on slices are available.

#### Predicate functions

Elements in a slice can be tested with predicate functions. The `All` and `Any` functions test the elements in a slice with a given predicate function and determine if all the elements or at least one of them are true under the predicate respectively.

```go
input1 := []rune("---------")
All(input1, func(r rune) bool { return r == '-'}) // true
Any(input1, func(r rune) bool { return r == '!'}) // false

input2 := []rune("-----!----") 
All(input2, func(r rune) bool { return r == '-'}) // false
Any(input2, func(r rune) bool { return r == '!'}) // true

```

#### Transformations

The functional primitives of `Map`, `Filter` and `Fold` are available.

The `Map` function creates a new slice by transforming all the elements of an existing slice by applying a function to each element.

```go
input := []int{1, 2, 3, 4}
actual := slices.Map(input, func(x int) int { return x * 2 }) // []int{2, 4, 6, 8}
```

The `Filter` function creates a new slice by selecting element to retain from an existing slice based on a predicate function.

```go
input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
slices.Filter(input, func(i int) bool { return i%2 == 0 }) // []int{2, 4, 6, 8}
```

The `Fold` function reduces all the elements of a slice down to a single value, using a function to combine elements.

### Range functions

A number of functions are available for generating a slice consisting of a sequence of numbers of various types, including floats. For example, the following call generates a slice consisting of the numbers from 0 to 4:

```go
slices.Range(0,5) // []int{0, 1, 2, 3 ,4}
```

This is an exclusive range which goes up to, but does not include the second parameter value. To generate an inclusive range, the `IncRange` function can be used, e.g.:

```go
slices.IncRange(0, 5) // []int{0, 1, 2, 3 ,4, 5}
```

The difference between each number is 1, unless the second parameter value is less than the first, in which case it is -1.

```go
slices.IncRange(5, 0) // []int{5, 4, 3, 2, 1, 0}
```

Floating point values can also be used in ranges:

```go
slices.IncRange(0.0, 5.0) // []float64{0.0, 1.0, 2.0, 3.0 ,4.0, 5.0}
```

If a non-unity difference between each slice element is required, this can be achieved with `RangeBy` or `IncRangeBy` functions.

```go
slices.RangeBy(0.0, 2.0, 0.5) // []float64{0.0, 0.5, 1.0, 1.5}
```

If the range is descending, a negative step is required, otherwise the function will panic:

```go
slices.RangeBy(2.0, 0.0, -0.5) // []float64{2.0, 1.5, 1.0, 0.5}

```

For very large ranges, if needed, functions are available for generating different parts of the range across multiple processor cores in parallel.  The `ParRange` function works like range, except it will try to accelerate it's execution for large ranges, across multiple cores.

```go
slices.ParRange(0, 400000) // []int{0, 1, 2, ..., 399999}
```

The function takes some optional parameters that control how the activities are parallelised. 

```go
slices.ParRange(0, 400000, ParThreshold(100000), ParMaxCpu(4))
```

The `ParThreshold` function controls the threshold beyond which the population of the slice is broken up in to parallel chunks; a range size below this value will be handled as a single chunk. The default value is 100000. The `ParMaxCpu` function controls the maximum number of parallel chunks. By default this is the number of CPU cores has; a number larger than this will typically result in lower performance.

As well as `ParRange` there are parallel range functions for each of the non-parallel ones, i.e. the following functions exist:

* `ParRange`
* `ParIncRange`
* `ParRangeBy`
* `ParIncRangeBy`

## Option type

Option types are used to hold "nullable" values whilst providing a greater degree of null safety than simple pointers. An option either holds a value (referred to as non-empty) or holds nothing (an empty option). Option types avoid potential additional heap allocation by avoiding the use of a pointer; they are implemented as an underlying value plus a boolean flag. They are particularly useful for representing nullable basic types.

### Usage

A simple option value can be created with the `Value` function.

```go
optInt := option.Value(10) // optInt is a non-empty option.Option[int] type
```

This creates a "non-empty" option holding an `int` value. Options implement the `Stringer` interface, so they can be printed directly.

```go
message := fmt.Sprintf("Hello %s", option.Value("world")) // message == "Hello world"
```

However, the underlying value cannot be accessed directly but only via access methods. This is to encourage the programmer to give adequate attention to the empty case. The `Get` method will return the value of a non-empty option but it will panic if the option is empty. The `IsEmpty` method can be used to detect the empty case. Consider the following function adding a number to an option.Option[int].

```go
func optAdd(o option.Option[int], n int) option.Option[int] {
  if o.IsEmpty() {
    return option.Empty[int]()
  } else {
    return option.Value(o.Get() + n)
  }
}
```

The function is using `IsEmpty` to validate the option is non-empty before attempting to access the value with `Get`. It returns an option as a result as it may return no value if there is no value in the option supplied. If a value is supplied, an option is returned with result of the addition. This kind of pattern is actually quite common and the option library has utility functions to perform this kind of option chaining. The effect of the above function can be achieved with the `Morph` method that applies a function to a non empty option.

```go
func optAdd(o option.Option[int], n int) option.Option[int] {
  return o.Morph(func (x int) int { return x + n })
}
```

An alternative approach, if the option is expected to be non-empty and is an error otherwise, is to access the value with a `Try` method.

```go
func optAdd(o option.Option[int], n int) (result int, err error) {
  defer Catch(&err)
  result = o.Try() + n
  return
}
```

Here the `Catch` function in the error handling package is used to ensure the function just returns an error value in response to an unexpectedly empty option rather than causing panics. The error returned will be `option.ErrOptionIsEmpty`. This kind of approach can be especially useful in functions that process a number of options which should not be empty. It is not recommended for options for which empty values are a non-exceptional condition due to the extra overhead of handling the error processing path.

### Zero value

The zero value of an option is empty. For example.

```go
var zero Option[int]
fmt.Println(zero.IsEmpty()) // true
```

### Comparisons

Two options of the same type can be compared successfully with `==` provided the underlying types can be likewise compared. An empty option will always compare as not equal to a non-empty one.

### Mutations

If you are wrapping a `struct` inside an option, there are methods that allow mutation in place, avoiding any need for wholesale copying of the struct data. The `Ensure` method will set a given option to non-empty, if it isn't already, and `Mutate` allows updates to be performed against a non-empty option value (against an empty option, it is a no-op). These together allows an empty option containing no data to be mutated to contain any desired values.

```go
type nv struct{ name, value string }
opt := Empty[nv]()
opt.Ensure().Mutate(func(n *nv) {
  n.name = "name"
  n.value = "value"
})
```

### Marshalling and Unmarshalling

Option types support marshalling and unmarshalling via the `encoding/json` or `gopkg.in/yaml.v2` packages. Note that `yaml.v3` is not yet supported. A non-empty option is marshalled as simply the value it contains in both JSON and YAML, e.g.

```go
type testOptMarshall struct {
  Name  Option[string] `json:"name,omitempty" yaml:"name,omitempty"`
  Value Option[int]    `json:"value,omitempty" yaml:"value,omitempty"`
}
```

```go
testData := testOptMarshall{
  Name:  Value("a name"),
  Value: Value(123),
}
y := Try(json.Marshal(&testData))
text := string(y) // "{\"name\":\"a name\",\"value\":123}"
```

Empty options are rendered as "null":

```go
testData := testOptMarshall{
  Name:  Value("a name"),
  Value: Empty[int](),
}
y := Try(json.Marshal(&testData))
text := string(y) // "{\"name\":\"a name\",\"value\":null}"
```

However, when rendering YAML, and the `omitempty` annotation is present, any empty values will be omitted:

```go
testData := testOptMarshall{
  Name:  Value("a name"),
  Value: Empty[int](),
}
y := Try(json.Marshal(&testData))
text := string(y) // "name: a name\n"
```

This is unfortunately not true of JSON marshalling, due to limitations of the `json.Marshalling` interface.
#### Unmarshalling

When unmarshalling to option values, omitted keys and `null` values in JSON or YAML are unmarshalled as an empty option.

Eg unmarshalling YAML:

```go
const input := "name: a name\n"
var unmarshalledData testOptMarshall
Try0(yaml.Unmarshall([]byte(input), &unmarshalledData))
unmarshalled.Name.HasValue()  // true
unmarshalled.Value.HasValue() // false
```

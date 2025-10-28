/*
The iterator package provides implementation for the creation, consumption and
transformation of various types kinds of generic iterators.

# Features

Iterators in this package have the following features:

  - Support for iteration over single elements or key/value pairs.

  - Constructable from [iter.Seq] and [iter.Seq2] objects.

  - Constructable over slices, maps and from literal elements.

  - Support for mutation of elements and underlying collections.

  - Convertible to [iter.Seq] or [iter.Seq2] objects

  - Collectable into slices or maps.

  - Transformable via various mapping, filtering and reducing methods.

# Types of iterators

There are four main iterator types, encapsulated in four interfaces.

  - [iterator.Iterator] - A standard generic iterator that yields single values.
  - [iterator.Iterator2] - An iterator that yields pairs of values, typically a
    value plus a key.
  - [iterator.MutableIterator] - An extension of [iterator.Iterator] that allows
    modification and removal of values from some underlying collection (such as
    slice).
  - [iterator.MutableIterator2] - An extension of [iterator.Iterator2] that
    allows modification and removal of key value pairs from some underlying
    collection (such as map). Note that only the value can be modified, not the
    key.

The above top level interfaces are composite definitions such that instances of any
of them will also implement other key interfaces.

Any [iterator.Iterator][T] also implements:
 - [iterator.CoreIterator][T]
 - [iterator.IteratorExtensions][T]
 - [iterator.SimpleIterator][T]

Any [iterator.Iterator2][K,V] also implements:
 - [iterator.Iterator][V]
 - [iterator.CoreIterator2][K,V]
 - [iterator.Iterator2Extensions][K,V]

Any [iterator.MutableIterator][T] also implements:
 - [iterator.Iterator][T]
 - [iterator.CoreMutableIterator][T]

Any [iterator.MutableIterator2][K,V] also implements:
 - [iterator.MutableIterator][V]
 - [iterator.Iterator2][K,V]

# Construction

Out of the box methods exist to produce iterators:

  - Over slices via [github.com/robdavid/genutil-go/slices.Iter],
    [github.com/robdavid/genutil-go/slices.IterRef] or
    [github.com/robdavid/genutil-go/slices.IterMut] methods.
  - Over maps via [github.com/robdavid/genutil-go/maps.Iter],
    [github.com/robdavid/genutil-go/maps.IterMut].
  - Over number ranges via [iterator.Range], [iterator.IncRange],
    [iterator.RangyBy] or [iterator.IncRangeBy] functions.
  - Over an explicit list of elements via the [iterator.Of] function.

# User iterators

New iterators can be created a number of ways.

  - An [iter.Seq] function can be transformed to an [iterator.Iterator] via the
    [iterator.New] or [iterator.NewWithSize] functions. This is the simplest and
    most recommended way to build an iterator.
  - An [iter.Seq2] function can be transformed to an [iterator.Iterator2] via
    the [iterator.New2] or [iterator.New2WithSize] functions. This is the
    simplest and most recommended way to build an iterator of key/value pairs.
  - A user may create an implementation of [iterator.SimpleIterator] and convert
    it to an [iterator.Iterator] with [iterator.NewFromSimple] or
    [iterator.NewFromSimpleWithSize].
  - A user may create an implementation of [iterator.SimpleMutableIterator] and
    convert it to an [iterator.MutableIterator] with
    [iterator.NewFromSimpleMutable] or [iterator.NewFromSimpleMutableWithSize].
  - A user may create an implementation of [iterator.CoreIterator] and convert
    it to an [iterator.Iterator] with [iterator.NewDefaultIterator].
  - A user may create an implementation of [iterator.CoreMutableIterator] and
    convert it to an [iterator.MutableIterator] with
    [iterator.NewDefaultMutableIterator].
  - A user may create an implementation of [iterator.CoreIterator2] and convert
    it to an [iterator.Iterator2] with [iterator.NewDefaultIterator2].
  - A user may create an implementation of [iterator.CoreMutableIterator2] and
    convert it to an [iterator.MutableIterator2] with
    [iterator.NewDefaultMutableIterator2].

# Consumption

The most straightforward and recommended way to consume elements from an
iterator is to use the [iterator.CoreIterator.Seq] or
[iterator.CoreIterator2.Seq2] methods to convert the iterator to an [iter.Seq]
or [iter.Seq2] object, and use a for loop, eg:
*/
/*
  for n := range iterator.Range(0, 5).Seq() {
    fmt.Printf("%d\n", n)
  }
  // Output:
  // 0
  // 1
  // 2
  // 3
  // 4
*/
/*
This is usually the most efficient approach, especially if the iterator is build
on top of an [iter.Seq] or [iter.Seq2] already. However, you can also loop over
the iterator using the [iterator.SimpleIterator.Next] and
[iterator.SimpleIterator.Value] methods, eg:
*/
/*
	for itr := iterator.Range(0, 5); itr.Next(); {
		fmt.Printf("%d\n", itr.Value())
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
*/
/*
It is also possible to collect all the elements in an iterator into a slice or a
map.

An [iterator.CoreIterator] may be collected into a slice using one of the
following methods:
 - [iterator.DefaultIterator.Collect]
 - [iterator.DefaultIterator.CollectInto]
 - [iterator.DefaultIterator.CollectIntoCap]

An [iterator.CoreIterator2] may be collected into a slice of [iterator.KeyValue]
elements via one of the following methods:
 - [iterator.DefaultIterator.Collect2]
 - [iterator.DefaultIterator.Collect2Into]
 - [iterator.DefaultIterator.Collect2IntoCap]

An [iterator.CoreIterator2] may be collected into a map via one of the following
functions:
 - [iterator.CollectMap]
 - [iterator.CollectIntoMap]

# Transformations

Iterator types support a number of transformation methods and functions covering
such things as filtering, mapping and folding over elements. These are defined
in interface definitions [iterator.IteratorExtensions] and
[iterator.Iterator2Extensions]. Default implementations are provided by
[iterator.DefaultIterator] and [iterator.DefaultIterator2]. These are listed
below.

  - [iterator.DefaultIterator.All]
  - [iterator.DefaultIterator.Any]
  - [iterator.DefaultIterator.Chan]
  - [iterator.DefaultIterator2.Chan2]
  - [iterator.DefaultIterator.Collect]
  - [iterator.DefaultIterator2.Collect2]
  - [iterator.DefaultIterator2.Collect2Into]
  - [iterator.DefaultIterator2.Collect2IntoCap]
  - [iterator.DefaultIterator.CollectInto]
  - [iterator.DefaultIterator.CollectIntoCap]
  - [iterator.DefaultIterator.Enumerate]
  - [iterator.DefaultIterator.Filter]
  - [iterator.DefaultIterator2.Filter2]
  - [iterator.DefaultIterator.FilterMorph]
  - [iterator.DefaultIterator2.FilterMorph2]
  - [iterator.DefaultIterator.Fold]
  - [iterator.DefaultIterator.Fold1]
  - [iterator.DefaultIterator.Intercalate]
  - [iterator.DefaultIterator.Intercalate1]
  - [iterator.DefaultIterator.Morph]
  - [iterator.DefaultIterator2.Morph2]
  - [iterator.DefaultIterator.Take]
  - [iterator.DefaultIterator2.Take2]

There are also corresponding functions in the iterator package, which include
some additional or modified functions that can't be expressed as methods due to
limitations of generics. All methods defined in [iterator.DefaultIterator] or
[iterator.DefaultIterator2] delegate to one of these functions.

  - [iterator.All]
  - [iterator.Any]
  - [iterator.Chan]
  - [iterator.Chan2]
  - [iterator.Take]
  - [iterator.Collect]
  - [iterator.Collect2]
  - [iterator.Collect2Into]
  - [iterator.Collect2IntoCap]
  - [iterator.CollectInto]
  - [iterator.CollectIntoCap]
  - [iterator.CollectIntoMap]
  - [iterator.CollectMap]
  - [iterator.Enumerate]
  - [iterator.Filter]
  - [iterator.Filter2]
  - [iterator.FilterMap]
  - [iterator.FilterMap2]
  - [iterator.Fold]
  - [iterator.Fold1]
  - [iterator.Intercalate]
  - [iterator.Intercalate1]
  - [iterator.Map]
  - [iterator.Map2]
  - [iterator.Seq]
  - [iterator.Seq2]
  - [iterator.Take]
  - [iterator.Take2]

**/
package iterator

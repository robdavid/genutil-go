/*
Package iterator provides implementation for the creation, consumption and
transformation of various types kinds of generic iterators.

# Features

Iterators in this package have the following features:

  - Constructable from [iter.Seq] and [iter.Seq2] objects.
  - Constructable over slices, maps and from explicit elements.
  - Convertible to [iter.Seq] or [iter.Seq2] objects.
  - Collectable into slices or maps.
  - Support for mutability.
  - Transformable via various mapping, filtering and reducing methods.
*/
package iterator

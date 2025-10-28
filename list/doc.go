/*
The list package provides a double linked list. Compared with some other
list or vector implementations where the list behaves as if an array, indexed at
zero, this has some different characteristics.

  - Any [List] object represents a list. The list may be empty.
  - A [List] object does not necessarily refer to the start of a list; it may
    be pointing at any point within a list, or even the end.
  - It it possible to iterate both backwards and forwards over the list. See [List.Next], [List.Prev] or
    [List.Seq] and [List.RevSeq] for example.
  - Items can be numerically indexed, both by positive or negative indexes. Access to items at index n
    or -n takes place in O(n) time.
  - As well as methods to access elements, many access methods return new lists as a result. See, for example,
    the [List.At] method.
*/
package list

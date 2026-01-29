/*
The list package provides a double linked list. The list consists of a series of
linked nodes, which may be accessed, traversed, retained etc. independently of
the list in which they are a member. The caller may obtain a pointer to a node
from an arbitrary position in the list and this node will remain valid unless
deleted from the list. Traversal operations can be performed with a node pointer
alone; mutations require the containing list object as well.

Values in the list may be accessed by index or iteration; note that index access
is an O(n) operation requiring traversal of the list (n being the index value).

The list is "pointer-like" in that a value copy of a list actually refers to the
original list, with mutations being reflected in both copies. The zero value is
a nil list that contains no elements that can be read like an empty list but
cannot be mutated. This is similar to the native map object. The rationale for
this is that a "value-like" copy can result in an inconsistent state once
mutations are performed, and a zero value can only carry a nil pointer. Non-nil
lists are created with the Make function (or other constructors).
*/
package list

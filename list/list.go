package list

import (
	"errors"
	"iter"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/option"
)

var ErrIndexError = errors.New("index out of bounds")
var ErrNilNode = errors.New("node is nil")

// Node is an element within a doubly linked list.
type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

// Count returns the number of elements found by iterating forwards from the node.
func (node *Node[T]) Count() int {
	total := 0
	for n := node; n != nil; n = n.next {
		total++
	}
	return total
}

// RevCount returns the number of elements found by iterating backwards from the node.
func (node *Node[T]) RevCount() int {
	total := 0
	for n := node; n != nil; n = n.prev {
		total++
	}
	return total
}

// Seq returns a native iter.Seq iterator over element values moving forwards in the list.
func (node *Node[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := node; n != nil; n = n.next {
			if !yield(n.value) {
				break
			}
		}
	}
}

// SeqNode returns a native iter.Seq iterator over element nodes moving forwards in the list.
func (node *Node[T]) SeqNode() iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		for n := node; n != nil; n = n.next {
			if !yield(n) {
				break
			}
		}
	}
}

// RevSeq returns a native iter.Seq iterator over element values moving backwards in the list.
func (node *Node[T]) RevSeq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := node; n != nil; n = n.prev {
			if !yield(n.value) {
				break
			}
		}
	}
}

// SeqNode returns a native iter.Seq iterator over element nodes moving backwards in the list.
func (node *Node[T]) RevSeqNode() iter.Seq[*Node[T]] {
	return func(yield func(*Node[T]) bool) {
		for n := node; n != nil; n = n.prev {
			if !yield(n) {
				break
			}
		}
	}
}

// Iter returns an iterator over element values moving forwards in the list.
func (node *Node[T]) Iter() iterator.Iterator[T] {
	return iterator.New(node.Seq())
}

// IterNode returns an iterator over element nodes moving forwards in the list.
func (node *Node[T]) IterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(node.SeqNode())
}

// RevIter returns an iterator over element values moving backwards in the list.
func (node *Node[T]) RevIter() iterator.Iterator[T] {
	return iterator.New(node.RevSeq())
}

// RevIterNode returns an iterator over element nodes moving backwards in the list.
func (node *Node[T]) RevIterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(node.RevSeqNode())
}

// NodeToValue is a helper function that extracts the value from a node.
func NodeToValue[T any](node *Node[T]) T { return node.value }

// List is a doubly linked list. More precisely it may reference any given element within such a list, which
// may be at the beginning, the end, or any point in between. The list can be traversed in either
// direction. The list may be empty, in which case it references nothing.

type List[T any] struct {
	first, last *Node[T]
	size        int
}

// Make creates a new empty list
func Make[T any]() List[T] {
	return List[T]{nil, nil, 0}
}

// New returns a pointer to a new empty list
func New[T any]() *List[T] {
	return &List[T]{nil, nil, 0}
}

// Of creates a new list whose elements are taken from the variadic args of the function
func Of[T any](slice ...T) List[T] {
	lst := Make[T]()
	lst.Append(slice...)
	//lst.first, lst.last = Chain(slice...)
	return lst
}

// Len returns the number of elements in the list
func (lst List[T]) Len() int {
	return lst.size
}

// Clear removes all elements from the list
func (lst *List[T]) Clear() {
	lst.first = nil
	lst.last = nil
	lst.size = 0
}

// Append adds the provided elements to the end of the list
func (lst *List[T]) Append(slice ...T) {
	start, end := Chain(slice...)
	if lst.first == nil {
		lst.first = start
		lst.last = end
	} else {
		lst.last.next = start
		start.prev = lst.last
		lst.last = end
	}
	lst.size += len(slice)
}

// Prepend adds the provided elements to the start of the list
func (lst *List[T]) Prepend(slice ...T) {
	start, end := Chain(slice...)
	if lst.first == nil {
		lst.first = start
		lst.last = end
	} else {
		lst.first.prev = end
		end.next = lst.first
		lst.first = start
	}
	lst.size += len(slice)
}

// Seq returns a native iter.Seq iterator over element values moving forwards in the list.
func (lst List[T]) Seq() iter.Seq[T] {
	return lst.first.Seq()
}

// RevSeq returns a native iter.Seq iterator over element values moving backwards in the list.
func (lst List[T]) RevSeq() iter.Seq[T] {
	return lst.last.RevSeq()
}

// SeqNode returns a native iter.Seq iterator over element nodes moving forwards in the list.
func (lst List[T]) SeqNode() iter.Seq[*Node[T]] {
	return lst.first.SeqNode()
}

// RevSeqNode returns a native iter.Seq iterator over element nodes moving backwards in the list.
func (lst List[T]) RevSeqNode() iter.Seq[*Node[T]] {
	return lst.last.RevSeqNode()
}

// Iter returns an iterator over element values moving forwards in the list.
func (lst List[T]) Iter() iterator.Iterator[T] {
	remain := lst.size
	return iterator.NewWithSize(
		func(yield func(T) bool) {
			for n := lst.first; n != nil; n = n.next {
				if !yield(n.value) {
					break
				}
				remain--
			}
		},
		func() iterator.IteratorSize { return iterator.NewSize(remain) },
	)
}

// IterNode returns an iterator over element nodes moving forwards in the list.
func (lst List[T]) IterNode() iterator.Iterator[*Node[T]] {
	remain := lst.size
	return iterator.NewWithSize(
		func(yield func(*Node[T]) bool) {
			for n := lst.first; n != nil; n = n.next {
				if !yield(n) {
					break
				}
				remain--
			}
		},
		func() iterator.IteratorSize { return iterator.NewSize(remain) },
	)
}

// RevIter returns an iterator over element values moving backwards through the list
func (lst List[T]) RevIter() iterator.Iterator[T] {
	return iterator.Map(lst.RevIterNode(), NodeToValue)
}

// RevIter returns an iterator over element nodes moving backwards through the list
func (lst List[T]) RevIterNode() iterator.Iterator[*Node[T]] {
	remain := lst.size
	return iterator.NewWithSize(
		func(yield func(*Node[T]) bool) {
			for n := lst.last; n != nil; n = n.prev {
				if !yield(n) {
					break
				}
				remain--
			}
		},
		func() iterator.IteratorSize { return iterator.NewSize(remain) },
	)
}

// FromSeq creates a new list whose elements are taken from the supplied iter.Seq iterator.
func FromSeq[T any](itr iter.Seq[T]) List[T] {
	lst := Make[T]()
	for e := range itr {
		lst.Append(e)
	}
	return lst
}

// From creates a new list whose elements are taken from the supplied iterator.
func From[T any](itr iterator.Iterator[T]) List[T] {
	if itr.SeqOK() {
		return FromSeq(itr.Seq())
	} else {
		lst := Make[T]()
		for itr.Next() {
			lst.Append(itr.Value())
		}
		return lst
	}
}

func As[U ~struct{ List[T] }, T any](l List[T]) U {
	return U{l}
}

// IsEmpty returns true if the list is empty
func (lst List[T]) IsEmpty() bool {
	return lst.first == nil
}

// First returns the first node in the list that this node is a member of, found
// by iterating backwards. If the node is nil, nil is returned.
func (node *Node[T]) First() *Node[T] {
	var first *Node[T]
	for n := node; n != nil; n = n.prev {
		first = n
	}
	return first
}

// Last returns the last node in the list that this node is a member of, found
// by iterating forwards. If node is nil, nil is returned.
func (node *Node[T]) Last() *Node[T] {
	var last *Node[T]
	for n := node; n != nil; n = n.next {
		last = n
	}
	return last
}

// Next returns the node following the current node. If there are no nodes following nil is returned.
func (node *Node[T]) Next() *Node[T] {
	return node.next
}

// Prev returns the node prior to the current node. If there is no element node nil is returned.
func (node *Node[T]) Prev() *Node[T] {
	return node.prev
}

// Get returns the value at the current node
func (node *Node[T]) Get() T {
	return node.value
}

// At returns the node n elements away from the current node; n may be negative as
// well as positive. If this moves beyond the limits of the list, it will panic with
// [ErrIndexError].
func (node *Node[T]) At(n int) *Node[T] {
	for n != 0 && node != nil {
		if n < 0 {
			node = node.prev
			n++
		} else {
			node = node.next
			n--
		}
	}
	if node == nil {
		panic(ErrIndexError)
	}
	return node
}

// GetAt returns the element n elements away from the start or the end of the
// list; n may be negative as well as positive, or zero (the first element). If
// negative, it counts as an offset from the end of the list plus one.
func (lst List[T]) Get(n int) T {
	return lst.At(n).value
}

// At returns the node n elements away from the start or the end of the list; n
// may be negative as well as positive, or zero (the first element). If
// negative, it counts as an offset from the end of the list plus one.
func (lst List[T]) At(n int) *Node[T] {
	if n < 0 {
		return lst.last.At(n + 1)
	} else {
		return lst.first.At(n)
	}
}

// Set sets the value at the element n elements away from the start or the end
// of the list; n may be negative as well as positive, or zero (the first
// element). If negative, it counts as an offset from the end of the list plus
// one.
func (lst *List[T]) Set(n int, value T) {
	lst.At(n).Set(value)
}

// GetFirst returns the element at the head of the given list. If the list is
// empty, an empty option is returned; otherwise it will hold the element's
// value.
func (lst List[T]) GetFirst() option.Option[T] {
	if lst.first == nil {
		return option.Empty[T]()
	} else {
		return option.Value(lst.first.value)
	}
}

// GetLast returns the element at the end of the given list. If the list is
// empty, an empty option is returned; otherwise it will hold the element's
// value.
func (lst List[T]) GetLast() option.Option[T] {
	if lst.last == nil {
		return option.Empty[T]()
	} else {
		return option.Value(lst.last.value)
	}
}

// First returns the first node of the list.
func (lst List[T]) First() *Node[T] {
	return lst.first
}

// Last returns the last node of the list.
func (lst List[T]) Last() *Node[T] {
	return lst.last
}

// Insert inserts new elements before the given node, moving the original element to following
// position.
func (lst *List[T]) Insert(node *Node[T], values ...T) {
	start, end := Chain(values...)
	if start == nil || end == nil {
		return
	}
	if node == nil {
		panic(ErrNilNode)
	}
	if node.prev == nil {
		if lst.first == node {
			// Make this conditional in case node has been deleted
			lst.first = start
		}
	} else {
		node.prev.next = start
		start.prev = node.prev
	}
	end.next = node
	node.prev = end
	lst.size += len(values)
}

// Insert inserts new elements after the given node, moving the previously following elements to after
// the inserted ones.
func (lst *List[T]) InsertAfter(node *Node[T], values ...T) {
	start, end := Chain(values...)
	if start == nil || end == nil {
		return
	}
	if node == nil {
		panic(ErrNilNode)
	}
	if node.next == nil {
		if lst.last == node {
			// Make this conditional in case node has been deleted
			lst.last = end
		}
	} else {
		node.next.prev = end
		end.next = node.next
	}
	start.prev = node
	node.next = start
	lst.size += len(values)
}

// Ref returns a reference to the element at the current node. It returns nil if
// the node pointer is nil
func (node *Node[T]) Ref() *T {
	if node == nil {
		return nil
	} else {
		return &node.value
	}
}

// Set sets the value at the current node. If the node pointer is nil, no action is taken
// and false is returned. Otherwise the element value is set and true is returned.
func (node *Node[T]) Set(value T) bool {
	if node == nil {
		return false
	} else {
		node.value = value
		return true
	}
}

// Delete removes a node from the list. If there is only one element in the list, the
// result will be the empty list.
func (lst *List[T]) Delete(node *Node[T]) {
	if node == nil {
		return
	}
	if node.prev == nil {
		if lst.first == node {
			lst.first = node.next
		}
	} else {
		node.prev.next = node.next
	}
	if node.next == nil {
		if lst.last == node {
			lst.last = node.prev
		}
	} else {
		node.next.prev = node.prev
	}
	lst.size--
}

// Chain links the provided into a series of nodes, returning the first and last nodes.
func Chain[T any](values ...T) (first *Node[T], last *Node[T]) {
	if len(values) == 0 {
		return nil, nil
	}
	var prev *Node[T]
	for _, value := range values {
		last = &Node[T]{value: value, prev: prev}
		if first == nil {
			first = last
		} else {
			prev.next = last
		}
		prev = last
	}
	return
}

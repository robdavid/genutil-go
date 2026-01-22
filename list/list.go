package list

import (
	"iter"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/option"
)

type Node[T any] struct {
	value T
	prev  *Node[T]
	next  *Node[T]
}

// List is a doubly linked list. More precisely it may reference any given element within such a list, which
// may be at the beginning, the end, or any point in between. The list can be traversed in either
// direction. The list may be empty, in which case it references nothing.

type list[T any] struct {
	first, last *Node[T]
}

type List[T any] = *list[T]

// type List[T any] struct {
// 	first **Node[T]
// 	last  **Node[T]
// }

// New creates a new empty list
func New[T any]() List[T] {
	return &list[T]{nil, nil}
}

// Of creates a new list whose elements are taken from the variadic args of the function
func Of[T any](slice ...T) List[T] {
	lst := New[T]()
	lst.Append(slice...)
	//lst.first, lst.last = Chain(slice...)
	return lst
}

// Size returns the number of elements in the list by iterating through the list
func (l *list[T]) Size() int {
	total := 0
	for n := l.first; n != nil; n = n.next {
		total++
	}
	return total
}

func (lst *list[T]) Append(slice ...T) {
	start, end := Chain(slice...)
	if lst.first == nil {
		lst.first = start
		lst.last = end
	} else {
		lst.last.next = start
		start.prev = lst.last
		lst.last = end
	}
}

func (lst *list[T]) Prepend(slice ...T) {
	start, end := Chain(slice...)
	if lst.first == nil {
		lst.first = start
		lst.last = end
	} else {
		lst.first.prev = end
		end.next = lst.first
		lst.first = start
	}
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

func (lst *list[T]) Seq() iter.Seq[T] {
	return lst.first.Seq()
}

func (lst *list[T]) RevSeq() iter.Seq[T] {
	return lst.last.RevSeq()
}

func (lst *list[T]) SeqNode() iter.Seq[*Node[T]] {
	return lst.first.SeqNode()
}

func (lst *list[T]) RevSeqNode() iter.Seq[*Node[T]] {
	return lst.last.RevSeqNode()
}

func (node *Node[T]) Iter() iterator.Iterator[T] {
	return iterator.New(node.Seq())
}

func (node *Node[T]) IterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(node.SeqNode())
}

func (node *Node[T]) RevIter() iterator.Iterator[T] {
	return iterator.New(node.RevSeq())
}

func (node *Node[T]) RevIterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(node.RevSeqNode())
}

func (lst *list[T]) Iter() iterator.Iterator[T] {
	return iterator.New(lst.Seq())
}

func (lst *list[T]) IterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(lst.SeqNode())
}

func (lst *list[T]) RevIter() iterator.Iterator[T] {
	return iterator.New(lst.RevSeq())
}

func (lst *list[T]) RevIterNode() iterator.Iterator[*Node[T]] {
	return iterator.New(lst.RevSeqNode())
}

/*
// From creates a new list whose elements are taken from the supplied iterator.
func From[T any](itr iterator.Iterator[T]) List[T] {
	return FromSeq(itr.Seq())
}

// FromSeq creates a new list whose elements are taken from the supplied iter.Seq iterator.
func FromSeq[T any](itr iter.Seq[T]) List[T] {
	last := New[T]()
	var first List[T]
	for e := range itr {
		last.Append(e)
		if n := last.Next(); n.head != nil {
			last = n
		} else {
			first = last
		}
	}
	return first
}

func As[U ~struct{ List[T] }, T any](l List[T]) U {
	return U{l}
}
*/

// IsEmpty returns true if the list is empty
func (lst *list[T]) IsEmpty() bool {
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

// GetFirst returns the element at the head of the given list. If the list is empty, an empty
// option is returned; otherwise it will hold the element's value.
func (lst *list[T]) GetFirst() option.Option[T] {
	if lst.first == nil {
		return option.Empty[T]()
	} else {
		return option.Value(lst.first.value)
	}
}

func (lst *list[T]) GetLast() option.Option[T] {
	if lst.first == nil {
		return option.Empty[T]()
	} else {
		return option.Value(lst.first.value)
	}
}

func (lst *list[T]) First() *Node[T] {
	return lst.first
}

func (lst *list[T]) Last() *Node[T] {
	return lst.last
}

// At returns the node n elements away from the current node; n may be negative as
// well as positive. If this moves beyond the limits of the list, nil is returned
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
	return node
}

// GetAt returns the element n elements away from the start or the end of the list; n may be negative as well as
// positive, or zero (the first element). If negative, it counts as an offset from the end of the list plus one.
func (lst *list[T]) Get(n int) option.Option[T] {
	node := lst.At(n)
	if node == nil {
		return option.Empty[T]()
	} else {
		return option.Value(node.value)
	}
}

func (lst *list[T]) At(n int) *Node[T] {
	if n < 0 {
		return lst.last.At(n + 1)
	} else {
		return lst.first.At(n)
	}
}

// Insert inserts new elements before the given node, moving the original element to following
// position.
func (lst *list[T]) Insert(node *Node[T], values ...T) {
	start, end := Chain(values...)
	if start == nil || end == nil {
		return
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
}

// Insert inserts new elements after the given node, moving the previously following elements to after
// the inserted ones.
func (lst *list[T]) InsertAfter(node *Node[T], values ...T) {
	start, end := Chain(values...)
	if start == nil || end == nil {
		return
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
func (lst *list[T]) Delete(node *Node[T]) {
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

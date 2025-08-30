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

// List is a doubly linked list. More precisely it references an item within such a list, which
// may be at the beginning, the end, or any point in between. The list can be traversed in either
// direction. The list may be empty, in which case it references nothing.
type List[T any] struct {
	head *Node[T]
}

// New creates a new empty list
func New[T any]() List[T] {
	return List[T]{nil}
}

// Of creates a new list whose elements are taken from the variadic args of the function
func Of[T any](slice ...T) List[T] {
	lst := New[T]()
	l := len(slice)
	for i := l - 1; i >= 0; i-- {
		lst.Insert(slice[i])
	}
	return lst
}

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

// IsEmpty returns true if the list is empty
func (l List[T]) IsEmpty() bool {
	return l.head == nil
}

// Get returns the element at the head of the given list. If the list is empty, an empty
// option is returned; otherwise it will hold the element's value.
func (l List[T]) Get() option.Option[T] {
	if l.head == nil {
		return option.Empty[T]()
	} else {
		return option.Value(l.head.value)
	}
}

// Ref returns a reference to the element at the head of the given list. It returns nil if
// the list is empty.
func (l List[T]) Ref() *T {
	if l.head == nil {
		return nil
	} else {
		return &l.head.value
	}
}

// Set sets the value at the head of the given list. If the list is empty, no action is taken
// and false is returned. Otherwise the element value is set and true is returned.
func (l List[T]) Set(value T) bool {
	if l.head == nil {
		return false
	} else {
		l.head.value = value
		return true
	}
}

// Next returns the list following the head element. If there are not elements following, or if
// the list is empty, an empty list is returned.
func (l List[T]) Next() List[T] {
	head := l.head
	if head != nil {
		head = head.next
	}
	return List[T]{head}
}

// Prev returns the list prior to the head element. If there is no element prior, or if the list is
// empty, an empty list is returned.
func (l List[T]) Prev() List[T] {
	head := l.head
	if head != nil {
		head = head.prev
	}
	return List[T]{head}
}

// First returns the list starting at the earliest element in the list, found by iterating backwards.
// If the list is empty, an empty list is returned.
func (l List[T]) First() List[T] {
	var first *Node[T]
	for n := l.head; n != nil; n = n.prev {
		first = n
	}
	return List[T]{first}
}

// Last returns the list starting at the final element in the list, found by iterating backwards.
// If the list is empty, an empty list is returned.
func (l List[T]) Last() List[T] {
	var last *Node[T]
	for n := l.head; n != nil; n = n.next {
		last = n
	}
	return List[T]{last}
}

// Size returns the number of elements in the list counting from the current element and moving
// forward.
func (l List[T]) Size() int {
	total := 0
	for n := l.head; n != nil; n = n.next {
		total++
	}
	return total
}

// RevSize returns the number of elements in the list, counting from the current element and
// moving backward.
func (l List[T]) RevSize() int {
	total := 0
	for n := l.head; n != nil; n = n.prev {
		total++
	}
	return total
}

// Insert inserts a new element in the current position, moving the original element to following
// position.
func (l *List[T]) Insert(value T) {
	if l.head == nil {
		l.head = &Node[T]{value: value}
	} else {
		node := &Node[T]{value: value, next: l.head, prev: l.head.prev}
		l.head.prev = node
		l.head = node
	}
}

// Append adds an element to the end of the list
func (l *List[T]) Append(value T) {
	var node *Node[T]
	for node = l.head; node != nil && node.next != nil; node = node.next {
	}
	newNode := Node[T]{value: value, prev: node}
	if node == nil {
		l.head = &newNode
	} else {
		node.next = &newNode
	}
}

// SeqList returns a native iter.Seq iterator over items moving forward in the list, returning
// a new List starting at each element found.
func (l List[T]) SeqList() iter.Seq[List[T]] {
	return func(yield func(List[T]) bool) {
		for n := l.head; n != nil; n = n.next {
			if !yield(List[T]{n}) {
				break
			}
		}
	}
}

// RevSeqList returns a native iter.Seq iterator over elements moving backwards in the list, returning
// a new List starting at each element found.
func (l List[T]) RevSeqList() iter.Seq[List[T]] {
	return func(yield func(List[T]) bool) {
		for n := l.head; n != nil; n = n.prev {
			if !yield(List[T]{n}) {
				break
			}
		}
	}
}

// Seq returns a native iter.Seq iterator over element values moving forwards in the list.
func (l List[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := l.head; n != nil; n = n.next {
			if !yield(n.value) {
				break
			}
		}
	}
}

// RevSeq returns a native iter.Seq iterator over element values moving backwards in the list.
func (l List[T]) RevSeq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := l.head; n != nil; n = n.prev {
			if !yield(n.value) {
				break
			}
		}
	}
}

// IterList returns a iterator.Iterator over items moving forward in the list, returning
// a new List starting at each element found.
func (l List[T]) IterList() iterator.Iterator[List[T]] {
	return iterator.NewWithSize(l.SeqList(),
		func() iterator.IteratorSize { return iterator.NewSize(l.Size()) })
}

// RevIterList returns a iterator.Iterator over items moving forward in the list, returning
// a new List starting at each element found.
func (l List[T]) RevIterList() iterator.Iterator[List[T]] {
	return iterator.NewWithSize(l.RevSeqList(),
		func() iterator.IteratorSize { return iterator.NewSize(l.RevSize()) })
}

// Iter returns an iterator.Iterator over element values, moving forwards in the list.
func (l List[T]) Iter() iterator.Iterator[T] {
	return iterator.NewWithSize(l.Seq(),
		func() iterator.IteratorSize { return iterator.NewSize(l.Size()) })
}

// Iter returns an iterator.Iterator over element values, moving backwards in the list.
func (l List[T]) RevIter() iterator.Iterator[T] {
	return iterator.NewWithSize(l.RevSeq(),
		func() iterator.IteratorSize { return iterator.NewSize(l.RevSize()) })
}

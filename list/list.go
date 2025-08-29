package list

import (
	"iter"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/option"
)

type Node[T any] struct {
	Value T
	Prev  *Node[T]
	Next  *Node[T]
}

type List[T any] struct {
	head *Node[T]
}

func New[T any]() List[T] {
	return List[T]{nil}
}

func From[T any](slice ...T) List[T] {
	lst := New[T]()
	l := len(slice)
	for i := l - 1; i >= 0; i-- {
		lst.Insert(slice[i])
	}
	return lst
}

func (l List[T]) Get() option.Option[T] {
	if l.head == nil {
		return option.Empty[T]()
	} else {
		return option.Value(l.head.Value)
	}
}

func (l List[T]) Ref() *T {
	if l.head == nil {
		return nil
	} else {
		return &l.head.Value
	}
}

func (l List[T]) Set(value T) bool {
	if l.head == nil {
		return false
	} else {
		l.head.Value = value
		return true
	}
}

func (l List[T]) Next() option.Option[List[T]] {
	if l.head == nil || l.head.Next == nil {
		return option.Empty[List[T]]()
	} else {
		return option.Value(List[T]{l.head.Next})
	}
}

func (l List[T]) Prev() option.Option[List[T]] {
	if l.head == nil || l.head.Prev == nil {
		return option.Empty[List[T]]()
	} else {
		return option.Value(List[T]{l.head.Prev})
	}
}

func (l List[T]) Size() int {
	total := 0
	for n := l.head; n != nil; n = n.Next {
		total++
	}
	return total
}

func (l *List[T]) Insert(value T) {
	if l.head == nil {
		l.head = &Node[T]{Value: value}
	} else {
		l.head = &Node[T]{Value: value, Next: l.head, Prev: l.head.Prev}
	}
}

func (l *List[T]) Append(value T) {
	var node *Node[T]
	for node = l.head; node != nil && node.Next != nil; node = node.Next {
	}
	newNode := Node[T]{Value: value, Prev: node}
	if node == nil {
		l.head = &newNode
	} else {
		node.Next = &newNode
	}
}

func (l List[T]) SeqList() iter.Seq[List[T]] {
	return func(yield func(List[T]) bool) {
		for n := l.head; n != nil; n = n.Next {
			if !yield(List[T]{n}) {
				break
			}
		}
	}
}

func (l List[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		for n := l.head; n != nil; n = n.Next {
			if !yield(n.Value) {
				break
			}
		}
	}
}

func (l List[T]) IterList() iterator.Iterator[List[T]] {
	return iterator.NewWithSize(l.SeqList(),
		func() iterator.IteratorSize { return iterator.NewSize(l.Size()) })
}

func (l List[T]) Iter() iterator.Iterator[T] {
	return iterator.NewWithSize(l.Seq(),
		func() iterator.IteratorSize { return iterator.NewSize(l.Size()) })
}

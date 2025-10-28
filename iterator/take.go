package iterator

import "iter"

type takeIterator[T any] struct {
	count, max int
	aborted    bool
	iterator   CoreIterator[T]
}

func (ti *takeIterator[T]) Value() T {
	return ti.iterator.Value()
}

func (ti *takeIterator[T]) Abort() {
	if !ti.aborted {
		ti.iterator.Abort()
	}
	ti.aborted = true
}

func (ti *takeIterator[T]) Reset() {
	ti.count = 0
	ti.iterator.Reset()
}

func (ti *takeIterator[T]) Next() bool {
	if !ti.aborted && ti.count < ti.max {
		ti.count++
		return ti.iterator.Next()
	} else {
		return false
	}
}

// Non-seq is more efficient  here
func (ti *takeIterator[T]) SeqOK() bool { return false }

func (ti *takeIterator[T]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		if ti.count < ti.max && !ti.aborted {
			next, _ := iter.Pull(ti.iterator.Seq())
			for ti.count < ti.max && !ti.aborted {
				if value, ok := next(); ok {
					if !yield(value) {
						ti.aborted = true
						break
					}
					ti.count++
				} else {
					break
				}
			}
		}
	}
}

func (ti *takeIterator[T]) Size() IteratorSize {
	itrSize := ti.iterator.Size()
	remain := ti.max - ti.count
	switch itrSize.Type {
	case SizeKnown:
		return NewSize(min(remain, itrSize.Size))
	case SizeUnknown:
		return NewSizeMax(remain)
	case SizeAtMost:
		return NewSizeMax(min(remain, itrSize.Size))
	case SizeInfinite:
		return NewSize(remain)
	default:
		panic(ErrInvalidIteratorSizeType)
	}
}

type takeIterator2[K any, V any] struct {
	takeIterator[V]
	iterator2 CoreIterator2[K, V]
}

func (ti2 *takeIterator2[K, V]) Key() K {
	return ti2.iterator2.Key()
}

func (ti *takeIterator2[K, V]) Seq2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		if ti.count < ti.max && !ti.aborted {
			next, _ := iter.Pull2(ti.iterator2.Seq2())
			for ti.count < ti.max && !ti.aborted {
				if key, value, ok := next(); ok {
					if !yield(key, value) {
						ti.aborted = true
						break
					}
					ti.count++
				} else {
					break
				}
			}
		}
	}
}

// Take transforms a [CoreIterator] into an [Iterator] that returns at most n
// elements of the original iterator. Note that any [Iterator] is also a
// [CoreIterator]. If there are less than or exactly n elements available, the
// returned iterator is equivalent to the input iterator.
func Take[T any](n int, iter CoreIterator[T]) Iterator[T] {
	return NewDefaultIterator(&takeIterator[T]{iterator: iter, max: n})
}

// Take2 transforms a [CoreIterator2] into an [Iterator2] that returns at most n
// pairs of elements. Note that any [Iterator2] is also a [CoreIterator2]. If
// there are less than or exactly n elements available, the returned iterator is
// equivalent to the input iterator.
func Take2[K any, V any](n int, iter CoreIterator2[K, V]) Iterator2[K, V] {
	return NewDefaultIterator2(&takeIterator2[K, V]{takeIterator: takeIterator[V]{iterator: iter, max: n}, iterator2: iter})
}

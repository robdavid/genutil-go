package iterator

import (
	"fmt"
	"iter"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/internal/rangehelper"
	"github.com/robdavid/genutil-go/ordered"
)

type rangeIterInitial[T ordered.Real] struct {
	index     T
	inclusive bool
}
type rangeIter[T ordered.Real, S ordered.Real] struct {
	index, to T
	by        S
	value     T
	inclusive bool
	initial   rangeIterInitial[T]
}

func (ri *rangeIter[T, S]) incdec() {
	if ri.by < 0 {
		ri.index -= T(-ri.by) // T might not be signed
	} else {
		ri.index += T(ri.by)
	}
}

func (ri *rangeIter[T, S]) validateRange() {
	if ri.by == 0 && ri.index != ri.to {
		panic(fmt.Errorf("%w: step is zero", ErrInvalidIteratorRange))
	}
	if (ri.by > 0 && ri.to < ri.index) || (ri.by < 0 && ri.to > ri.index) {
		panic(fmt.Errorf("%w: negative step or inverse range (but not both)", ErrInvalidIteratorRange))
	}
}

func (ri *rangeIter[T, S]) Next() bool {
	if ri.index == ri.to {
		// Handles the case where by is zero, which is valid if index is at the end
		if ri.inclusive {
			ri.value = ri.index
			ri.inclusive = false // Causes iterator to terminate next time
			return true
		} else {
			return false
		}
	}
	if (ri.by < 0 && ri.index < ri.to) || (ri.by > 0 && ri.index > ri.to) {
		return false
	}
	ri.value = ri.index
	ri.incdec()
	return true
}

func (ri *rangeIter[T, S]) Value() T {
	return ri.value
}

func (ri *rangeIter[T, S]) Abort() {
	ri.index = ri.to
	ri.inclusive = false
}

func (ri *rangeIter[T, S]) Reset() {
	ri.index = ri.initial.index
	ri.inclusive = ri.initial.inclusive
}

func (ri *rangeIter[T, S]) Seq() iter.Seq[T] {
	return func(yield func(T) bool) {
		defer ri.Abort()
		if ri.index == ri.to {
			if ri.inclusive {
				yield(ri.index)
			}
			return
		}
		size, aStep := rangehelper.RangeSize(ri.index, ri.to, ri.by, ri.inclusive)
		if ri.by < 0 {
			for range size {
				index := ri.index
				ri.index -= aStep
				if !yield(index) {
					break
				}
			}
		} else {
			for range size {
				index := ri.index
				ri.index += aStep
				if !yield(index) {
					break
				}
			}
		}

	}
}

func (ri *rangeIter[T, S]) Size() IteratorSize {
	var size int
	if ri.index == ri.to {
		size = functions.IfElse(ri.inclusive, 1, 0)
	} else if (ri.index > ri.to && ri.by > 0) || (ri.index < ri.to && ri.by < 0) {
		size = 0
	} else {
		size, _ = rangehelper.RangeSize(ri.index, ri.to, ri.by, ri.inclusive)
	}
	return NewSize(size)
}

func (ri *rangeIter[T, S]) SeqOK() bool { return false }

func newRangeIter[T ordered.Real, S ordered.Real](from, upto T, by S, inclusive bool) Iterator[T] {
	itr := rangeIter[T, S]{index: from, to: upto, by: by, inclusive: inclusive,
		initial: rangeIterInitial[T]{index: from, inclusive: inclusive}}
	itr.validateRange()
	return NewDefaultIterator(&itr)
}

// Range creates an iterator that produces a sequence of numeric values that
// range between from and upto exclusive. This sequence may be descending if
// upto is less than from. If from == upto, an empty iterator is returned.
func Range[T ordered.Real](from, upto T) Iterator[T] {
	return newRangeIter(from, upto, functions.IfElse(upto < from, -1, 1), false)
}

// Range creates an iterator that produces a sequence of numeric values that
// range between from and upto inclusive. This sequence may be descending if
// upto is less than from. If from == upto, an iterator yielding a single value
// is returned.
func IncRange[T ordered.Real](from, upto T) Iterator[T] {
	return newRangeIter(from, upto, functions.IfElse(upto < from, -1, 1), true)
}

// RangeBy creates an iterator that produces a sequence of numeric values that
// range between from and upto exclusive, with a difference between each value
// of by. The value of by can be negative, in which case upto should be less
// than from, but it cannot be zero unless from == upto, in which case an empty
// iterator is returned.
func RangeBy[T ordered.Real, S ordered.Real](from, upto T, by S) Iterator[T] {
	return newRangeIter(from, upto, by, false)
}

// RangeBy creates an iterator that produces a sequence of numeric values that
// range between from and upto inclusive, with a difference between each value
// of by. The value of by can be negative, in which case upto should be less
// than from, but it cannot be zero unless from == upto, in which case an
// iterator yielding a single value is returned.
func IncRangeBy[T ordered.Real, S ordered.Real](from, upto T, by S) Iterator[T] {
	return newRangeIter(from, upto, by, true)
}

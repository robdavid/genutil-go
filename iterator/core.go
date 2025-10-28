package iterator

import "iter"

// SimpleCoreIterator wraps a [SimpleIterator] and provides methods to implement [CoreIterator]
type SimpleCoreIterator[T any] struct {
	SimpleIterator[T]
	size func() IteratorSize
}

// NewSimpleCoreIterator builds a [CoreIterator] from a [SimpleIterator]
func NewSimpleCoreIterator[T any](itr SimpleIterator[T]) *SimpleCoreIterator[T] {
	return &SimpleCoreIterator[T]{SimpleIterator: itr}
}

// NewSimpleCoreIteratorWithSize builds a [CoreIterator] from a [SimpleIterator] plus a function
// that returns the remaining number of items in the iterator.
func NewSimpleCoreIteratorWithSize[T any](itr SimpleIterator[T], size func() IteratorSize) *SimpleCoreIterator[T] {
	return &SimpleCoreIterator[T]{SimpleIterator: itr, size: size}
}

// SimpleCoreMutableIterator wraps a [SimpleMutableIterator] and provides methods to implement [CoreMutableIterator]
type SimpleCoreMutableIterator[T any] struct {
	SimpleMutableIterator[T]
	size func() IteratorSize
}

// NewSimpleCoreMutableIterator builds a [CoreMutableIterator] from a [SimpleMutableIterator].
func NewSimpleCoreMutableIterator[T any](itr SimpleMutableIterator[T]) *SimpleCoreMutableIterator[T] {
	return &SimpleCoreMutableIterator[T]{SimpleMutableIterator: itr}
}

// NewSimpleCoreMutableIterator builds a [CoreMutableIterator] from a [SimpleMutableIterator] plus a function
// that returns the number of items remaining in the iterator.
func NewSimpleCoreMutableIteratorWithSize[T any](itr SimpleMutableIterator[T], size func() IteratorSize) *SimpleCoreMutableIterator[T] {
	return &SimpleCoreMutableIterator[T]{SimpleMutableIterator: itr, size: size}
}

// NewFromSimpleMutable builds a [MutableIterator] from a [SimpleMutableIterator].
func NewFromSimpleMutable[T any](itr SimpleMutableIterator[T]) MutableIterator[T] {
	return NewDefaultMutableIterator(NewSimpleCoreMutableIterator(itr))
}

// NewFromSimpleMutableWithSize builds a [MutableIterator] from a [SimpleMutableIterator] and a size
// function that returns the number of items remaining items in the iterator.
func NewFromSimpleMutableWithSize[T any](itr SimpleMutableIterator[T], size func() IteratorSize) MutableIterator[T] {
	return NewDefaultMutableIterator(NewSimpleCoreMutableIteratorWithSize(itr, size))
}

// SeqCoreIterator wraps an iter.Seq iterator, plus an optional size function, providing methods
// to implement CoreIterator.
type SeqCoreIterator[T any] struct {
	seq       iter.Seq[T]
	size      func() IteratorSize
	stop      func()
	next      func() (T, bool)
	value     T
	seqCalled bool
}

// NewSeqCoreIterator builds a [CoreIterator] from a standard library [iter.Seq]
func NewSeqCoreIterator[T any](seq iter.Seq[T]) *SeqCoreIterator[T] {
	return &SeqCoreIterator[T]{seq: seq}
}

// NewSeqCoreIterator builds a [CoreIterator] from a standard library [iter.Seq] and a function that
// returns the number of items in the iterator.
func NewSeqCoreIteratorWithSize[T any](seq iter.Seq[T], size func() IteratorSize) *SeqCoreIterator[T] {
	return &SeqCoreIterator[T]{seq: seq, size: size}
}

// SeqCoreIterator2 wraps an [iter.Seq2], providing methods to implement [CoreIterator2]
type SeqCoreIterator2[K any, V any] struct {
	*SeqCoreIterator[V]
	seq2 iter.Seq2[K, V]
	key  K
}

// NewSeqCoreIterator2 builds a [CoreIterator2] implementation from an [iter.Seq2] iterator
func NewSeqCoreIterator2[K any, V any](seq2 iter.Seq2[K, V]) *SeqCoreIterator2[K, V] {
	return NewSeqCoreIterator2WithSize(seq2, nil)
}

// NewSeqCoreIterator2WithSize builds a [CoreIterator2] implementation from an [iter.Seq2] and a function that returns the
// size of the remaining items in the iterator.
func NewSeqCoreIterator2WithSize[K any, V any](seq2 iter.Seq2[K, V], size func() IteratorSize) *SeqCoreIterator2[K, V] {
	itr2 := SeqCoreIterator2[K, V]{
		seq2: seq2,
	}
	seq := func(yield func(V) bool) {
		var v V
		for itr2.key, v = range itr2.seq2 {
			if !yield(v) {
				break
			}
		}
	}
	if size == nil {
		itr2.SeqCoreIterator = NewSeqCoreIterator(seq)
	} else {
		itr2.SeqCoreIterator = NewSeqCoreIteratorWithSize(seq, size)
	}
	return &itr2
}

type coreIteratorMutations[T any] struct {
	delete func()
	set    func(T)
}

type SeqCoreMutableIterator2[K any, V any] struct {
	SeqCoreIterator2[K, V]
	coreIteratorMutations[V]
}

func NewSeqCoreMutableIterator2[K any, V any](seq2 iter.Seq2[K, V], delete func(), set func(V)) *SeqCoreMutableIterator2[K, V] {
	return &SeqCoreMutableIterator2[K, V]{SeqCoreIterator2: *NewSeqCoreIterator2(seq2),
		coreIteratorMutations: coreIteratorMutations[V]{delete: delete, set: set}}
}

func NewSeqCoreMutableIterator2WithSize[K any, V any](seq2 iter.Seq2[K, V], delete func(), set func(V), size func() IteratorSize) *SeqCoreMutableIterator2[K, V] {
	return &SeqCoreMutableIterator2[K, V]{SeqCoreIterator2: *NewSeqCoreIterator2WithSize(seq2, size),
		coreIteratorMutations: coreIteratorMutations[V]{delete: delete, set: set}}
}

func (itr *SimpleCoreIterator[T]) Size() IteratorSize {
	if itr.size == nil {
		return NewSizeUnknown()
	} else {
		return itr.size()
	}
}

func (itr *SimpleCoreIterator[T]) Seq() iter.Seq[T] {
	return Seq(itr.SimpleIterator)
}

func (itr *SimpleCoreIterator[T]) SeqOK() bool {
	return false
}

func (itr *SimpleCoreMutableIterator[T]) Size() IteratorSize {
	if itr.size == nil {
		return NewSizeUnknown()
	} else {
		return itr.size()
	}
}

func (si *SimpleCoreMutableIterator[T]) Seq() iter.Seq[T] {
	return Seq(si)
}

func (itr *SimpleCoreMutableIterator[T]) SeqOK() bool {
	return false
}

func (si *SeqCoreIterator[T]) Seq() iter.Seq[T] {
	if si.next != nil {
		return si.pullSeq
	}
	return si.oneSeq
}

func (si *SeqCoreIterator[T]) oneSeq(yield func(T) bool) {
	if !si.seqCalled {
		si.seqCalled = true
		si.seq(yield)
	}
}

func (si *SeqCoreIterator[T]) pullSeq(yield func(T) bool) {
	for {
		v, ok := si.next()
		if !ok {
			break
		}
		if !yield(v) {
			break
		}
	}
}

func (si *SeqCoreIterator[T]) Size() IteratorSize {
	if si.size == nil {
		return NewSizeUnknown()
	} else {
		return si.size()
	}
}

func (si *SeqCoreIterator[T]) SeqOK() bool {
	return si.next == nil
}

func (si *SeqCoreIterator[T]) Next() (ok bool) {
	if si.seqCalled {
		panic("cannot call Next() on iterator after calling Seq() or Seq2()")
	}
	if si.next == nil {
		si.next, si.stop = iter.Pull(si.seq)
	}
	si.value, ok = si.next()
	return
}

func (si *SeqCoreIterator[T]) Value() T {
	return si.value
}

func (si *SeqCoreIterator[T]) Abort() {
	if si.stop == nil {
		si.next, si.stop = iter.Pull(si.seq)
	}
	si.stop()
}

// Reset restarts the iterator
func (si *SeqCoreIterator[T]) Reset() {
	si.Abort()
	si.stop = nil
	si.next = nil
	si.seqCalled = false
}

func (si *SeqCoreIterator2[K, V]) Key() K {
	return si.key
}

func (si *SeqCoreIterator2[K, V]) Seq2() iter.Seq2[K, V] {
	if si.next != nil {
		panic("cannot call Seq2() on iterator after calling Next()")
	}
	si.seqCalled = true
	return si.seq2
}

func (smi *SeqCoreMutableIterator2[K, V]) Delete() {
	smi.delete()
}

func (smi *SeqCoreMutableIterator2[K, V]) Set(v V) {
	smi.set(v)
}

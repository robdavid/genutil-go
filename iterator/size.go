package iterator

// The largest slice capacity we are prepared to allocate to collect
// iterators of uncertain size.
const maxUncertainAllocation = 100000

type IteratorSizeType int

const (
	SizeUnknown  IteratorSizeType = iota // SizeUnknown represents an unknown size.
	SizeKnown                            // SizeKnown represents a completely known size.
	SizeAtMost                           // SizeAtMost represents a number of elements whose upper limit is known.
	SizeInfinite                         // SizeInfinite represents the knowledge that an iterator will not end.
)

// IteratorSize holds iterator sizing information
type IteratorSize struct {
	Type IteratorSizeType
	Size int
}

var (
	SIZE_UNKNOWN  IteratorSize = IteratorSize{Type: SizeUnknown}
	SIZE_INFINITE              = IteratorSize{Type: SizeInfinite, Size: -1}
)

// Allocate returns an estimated allocation size needed to accommodate the remaining elements in the
// iterator. If the iterator size is infinite, the function will panic with
// iterator.ErrAllocationSizeInfinite.
func (isz IteratorSize) Allocate() int {
	switch isz.Type {
	case SizeUnknown:
		return 0
	case SizeKnown:
		return isz.Size
	case SizeInfinite:
		panic(ErrAllocationSizeInfinite)
	case SizeAtMost:
		{
			sz := isz.Size / 2
			if sz >= maxUncertainAllocation {
				sz = maxUncertainAllocation
			}
			return sz
		}
	}
	panic(ErrInvalidIteratorSizeType)
}

// Subset is used to transform an [IteratorSize] to one which is a subset of the
// current one. Given an iterator A whose size is described by the current [IteratorSize],
// this functions returns a size corresponding to the number of elements in iterator B, where
// iterator B contains no more than the number of elements of iterator A.
func (isz IteratorSize) Subset() IteratorSize {
	switch isz.Type {
	case SizeUnknown, SizeInfinite, SizeAtMost:
		return isz
	case SizeKnown:
		return IteratorSize{SizeAtMost, isz.Size}
	}
	panic(ErrInvalidIteratorSizeType)
}

// Iterator sizing information; size is unknown
func NewSizeUnknown() IteratorSize {
	return SIZE_UNKNOWN
}

// IsUnknown returns true if the given IteratorSize instance represents
// an unknown size
func (size IteratorSize) IsUnknown() bool {
	return size.Type == SizeUnknown
}

// NewSize creates an IteratorSize implementation that has a fixed size of n.
func NewSize(n int) IteratorSize { return IteratorSize{SizeKnown, n} }

// IsKnown returns true if the iterator size is one whose actual size is known.
func (size IteratorSize) IsKnown() bool {
	return size.Type == SizeKnown
}

// IsKnownToBe returns true if the iterator size is one whose actual size is known,
// and is equal to the given value.
func (size IteratorSize) IsKnownToBe(n int) bool {
	return size.Type == SizeKnown && size.Size == n
}

// NewSizeMax creates an IteratorSize implementation that has a size no more than n.
func NewSizeMax(n int) IteratorSize {
	return IteratorSize{SizeAtMost, n}
}

// IsMaxKnown returns true if the iterator size is one whose maximum size is known.
func (size IteratorSize) IsMaxKnown() bool {
	return size.Type == SizeAtMost
}

// IsMaxKnownToBe returns true if the iterator size is one whose maximum size is known, and
// is equal to the given value.
func (size IteratorSize) IsMaxKnownToBe(n int) bool {
	return size.Type == SizeAtMost && size.Size == n
}

func NewSizeInfinite() IteratorSize {
	return SIZE_INFINITE
}

func (size IteratorSize) IsInfinite() bool {
	return size.Type == SizeInfinite
}

package iterator

type IteratorSizeType int

const (
	SizeUnknown IteratorSizeType = iota
	SizeKnown
	SizeAtMost
	SizeInfinite
)

// IteratorSize holds iterator sizing information
type IteratorSize struct {
	Type IteratorSizeType
	Size int
}

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
	return IteratorSize{Type: SizeUnknown}
}

// IsSizeUnknown returns true if the given IteratorSize instance represents
// an unknown size
func IsSizeUnknown(size IteratorSize) bool {
	return size.Type == SizeUnknown
}

// NewSize creates an `IteratorSize` implementation that has a fixed size of `n`.
func NewSize(n int) IteratorSize { return IteratorSize{SizeKnown, n} }

// IsSizeKnown returns true if the iterator size is one whose actual size is known.
func IsSizeKnown(size IteratorSize) bool {
	return size.Type == SizeKnown
}

// NewSizeAtMost creates an `IteratorSize` implementation that has a size no more than n.
func NewSizeAtMost(n int) IteratorSize {
	return IteratorSize{SizeAtMost, n}
}

// IsSizeAtMost returns true if the iterator size is one whose maximum size is known.
func IsSizeAtMost(size IteratorSize) bool {
	return size.Type == SizeAtMost
}

func NewSizeInfinite() IteratorSize {
	return IteratorSize{SizeInfinite, -1}
}

func IsSizeInfinite(size IteratorSize) bool {
	return size.Type == SizeInfinite
}

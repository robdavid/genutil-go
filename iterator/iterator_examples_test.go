package iterator_test

import (
	"iter"
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/stretchr/testify/assert"
)

func TestFromSeqExample(t *testing.T) {
	assert := assert.New(t)

	// fib returns a native Go iterator (fibonacci sequence).
	fib := func(yield func(int) bool) {
		tail := [2]int{0, 1}
		for {
			if !yield(tail[1]) {
				return
			}
			tail[0], tail[1] = tail[1], tail[0]+tail[1]
		}
	}

	fibItr := iterator.New(fib) // iterator.Iterator[int]
	fibSeq := fibItr.Seq()      // iter.Seq[int]

	seqCheck := iter.Seq[int](fibSeq) // compile time check
	assert.NotNil(seqCheck)

	i := 0
	expected := []int{1, 1, 2, 3, 5, 8}
	for f := range fibSeq {
		if i >= len(expected) {
			break
		}
		assert.Equal(expected[i], f)
		i++
	}

}

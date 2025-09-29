package iterator_test

import (
	"bytes"
	"fmt"
	"iter"
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/slices"
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

const (
	size          = 10
	expectedPrint = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n"
)

func TestToSeq(t *testing.T) {
	var buffer bytes.Buffer
	for n := range iterator.Range(0, size).Seq() {
		fmt.Fprintf(&buffer, "%d\n", n)
	}
	assert.Equal(t, buffer.String(), expectedPrint)
}

func TestNextValue(t *testing.T) {
	var buffer bytes.Buffer
	for itr := iterator.Range(0, size); itr.Next(); {
		fmt.Fprintf(&buffer, "%d\n", itr.Value())
	}
	assert.Equal(t, buffer.String(), expectedPrint)
}

func TestToChan(t *testing.T) {
	var buffer bytes.Buffer
	for n := range iterator.Range(0, size).Chan() {
		fmt.Fprintf(&buffer, "%d\n", n)
	}
	assert.Equal(t, buffer.String(), expectedPrint)
}

func TestCollectToMap(t *testing.T) {
	m := iterator.CollectMap(iterator.Of("zero", "one", "two", "three").Enumerate()) // map[int]string{0: "zero", 1: "one", 2: "two", 3: "three"}
	assert.Equal(t, map[int]string{0: "zero", 1: "one", 2: "two", 3: "three"}, m)
}

func TestMutableSlice(t *testing.T) {
	s := slices.Range(0, 10)
	itr := slices.IterMut(&s)
	for n := range itr.Seq() {
		if n%2 == 1 {
			itr.Delete()
		} else {
			itr.Set(n / 2)
		}
	}
	fmt.Println(s)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, s)
}

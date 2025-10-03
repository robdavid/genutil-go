package iterator_test

import (
	"bytes"
	"fmt"
	"iter"
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/maps"
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

func TestRangeByExample(t *testing.T) {
	ascending := iterator.IncRangeBy(0, 5, 2) // 0,2,4
	descending := iterator.RangeBy(5, 0, -2)  // 5,3,1
	assert.Equal(t, []int{0, 2, 4}, ascending.Collect())
	assert.Equal(t, []int{5, 3, 1}, descending.Collect())
}

func TestCollectToMap(t *testing.T) {
	m := iterator.CollectMap(iterator.Of("zero", "one", "two", "three").Enumerate()) // map[int]string{0: "zero", 1: "one", 2: "two", 3: "three"}
	assert.Equal(t, map[int]string{0: "zero", 1: "one", 2: "two", 3: "three"}, m)
}

func TestFilterExample(t *testing.T) {
	predicate := func(n int) bool { return n%2 == 0 }
	i := iterator.IncRange(1, 5).Filter(predicate)
	c := i.Collect() // []int{2,4}
	assert.Equal(t, []int{2, 4}, c)
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

func TestMutableMap(t *testing.T) {
	m := make(map[int]int)
	for i := range 10 {
		m[i] = i + 10
	}
	itr := maps.IterMut(m)
	for k, v := range itr.Seq2() {
		if k%2 == 1 {
			itr.Delete()
		} else {
			itr.Set(v / 2)
		}
	}
	fmt.Println(m) // map[0:5 2:6 4:7 6:8 8:9]
	assert.Equal(t, map[int]int{0: 5, 2: 6, 4: 7, 6: 8, 8: 9}, m)
}

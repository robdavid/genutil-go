package iterator_test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/iterator"
	"github.com/stretchr/testify/assert"
)

func ExampleNew() {
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

	i := 0
	const max = 5
	for f := range fibSeq {
		if i > max {
			break
		}
		fmt.Println(f)
		i++
	}
	// Output
	// 1
	// 1
	// 2
	// 3
	// 5
}

const (
	size          = 10
	expectedPrint = "0\n1\n2\n3\n4\n5\n6\n7\n8\n9\n"
)

func ExampleSimpleCoreIterator_Seq() {
	for n := range iterator.Range(0, 5).Seq() {
		fmt.Printf("%d\n", n)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleSeq() {
	for n := range iterator.Seq(iterator.Range(0, 5)) {
		fmt.Printf("%d\n", n)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleSeq2() {
	// Create an iterator of [KeyValue] pairs from enumerated range
	iterKV := iterator.AsKV(iterator.Range(5, 10).Enumerate())
	for k, v := range iterator.Seq2(iterKV) {
		fmt.Printf("%d: %d\n", k, v)
	}
	// Output:
	// 0: 5
	// 1: 6
	// 2: 7
	// 3: 8
	// 4: 9
}

func ExampleSeqCoreIterator_Next() {
	for itr := iterator.Range(0, 5); itr.Next(); {
		fmt.Printf("%d\n", itr.Value())
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleDefaultIterator_Chan() {
	for n := range iterator.Range(0, 5).Chan() {
		fmt.Println(n)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4

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

// counter is a SimpleIterator implementation that produces an
// infinite string of integers, starting from 0.
type counter struct {
	value   int  // value is the current value
	count   int  // count is the next value
	aborted bool // if true, the iterator is aborted
}

// Next sets value to the next count, increments the count, and
// returns true, unless aborted. When aborted, it is a no-op.
func (c *counter) Next() bool {
	if c.aborted {
		return false
	} else {
		c.value = c.count
		c.count++
		return true
	}
}

// Value returns the current value.
func (c *counter) Value() int {
	return c.value
}

// Abort stops the iterator by setting the aborted flag.
func (c *counter) Abort() {
	c.aborted = true
}

// Reset sets the counter back to 0.
func (c *counter) Reset() {
	c.count = 0
}

func ExampleNewFromSimple() {
	i := iterator.NewFromSimple(&counter{})
	c := i.Take(10).Collect()
	fmt.Println(c)
	// Output:
	// [0 1 2 3 4 5 6 7 8 9]
}

func ExampleNewFromSimpleWithSize() {
	i := iterator.NewFromSimpleWithSize(&counter{},
		func() iterator.IteratorSize { return iterator.SIZE_INFINITE })
	func() {
		defer func() { fmt.Println(recover()) }()
		i.Collect() // Attempting to collect the infinite iterator will panic.
	}()
	c := i.Take(10).Collect() // Collecting only the first 10 elements succeeds.
	fmt.Println(c)
	// Output:
	// cannot allocate storage for an infinite iterator
	// [0 1 2 3 4 5 6 7 8 9]
}

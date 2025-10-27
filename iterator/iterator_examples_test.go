package iterator_test

import (
	"fmt"
	"iter"
	"testing"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/robdavid/genutil-go/maps"
	"github.com/robdavid/genutil-go/slices"
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

func ExampleDefaultIterator_Collect() {
	c := iterator.Range(0, 5).Collect()
	fmt.Printf("%#v\n", c)
	// Output: []int{0, 1, 2, 3, 4}
}

func ExampleDefaultIterator_Morph() {
	f := func(n int) int { return n * 2 }
	i := iterator.Range(0, 5).Morph(f)
	c := i.Collect()
	fmt.Printf("%#v\n", c)
	// Output: []int{0, 2, 4, 6, 8}
}

func ExampleDefaultIterator_Filter() {
	predicate := func(n int) bool { return n%2 == 0 }
	i := iterator.IncRange(1, 5).Filter(predicate)
	c := i.Collect()
	fmt.Printf("%#v\n", c)
	// Output: []int{2, 4}
}

func ExampleDefaultIterator_FilterMorph() {
	// Function to filter on even values, doubling each selected value.
	f := func(v int) (int, bool) { return v * 2, v%2 == 0 }
	i := iterator.Of(0, 1, 2, 3, 4).FilterMorph(f)
	c := i.Collect()
	fmt.Printf("%#v\n", c)
	// Output: []int{0, 4, 8}

}

func ExampleDefaultIterator2_FilterMorph2() {
	// Function to filter on even values, doubling each selected value.
	inputMap := map[int]int{0: 2, 1: 4, 2: 6, 3: 8}
	itr := maps.Iter(inputMap).FilterMorph2(func(k, v int) (int, int, bool) {
		return k + 1, v * 2, (k+v)%2 == 0
	})
	c := iterator.CollectMap(itr)
	fmt.Printf("%#v\n", c)
	// Output: map[int]int{1:4, 3:12}
}

func ExampleDefaultIterator_Take() {
	i := iterator.Range(0, 100).Take(5)
	c := i.Collect()
	fmt.Printf("%#v\n", c)
	// Output: []int{0, 1, 2, 3, 4}
}

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

func ExampleCollectMap() {
	i := iterator.Of("zero", "one", "two", "three").Enumerate()
	m := iterator.CollectMap(i)
	fmt.Printf("%#v\n", m)
	// Output: map[int]string{0:"zero", 1:"one", 2:"two", 3:"three"}
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
	// cannot consume an infinite iterator
	// [0 1 2 3 4 5 6 7 8 9]
}

// coreCounter is a CoreIterator implementation that produces an
// infinite string of integers, starting from 0.
type coreCounter struct {
	value   int  // value is the current value
	count   int  // count is the next value
	aborted bool // if true, the iterator is aborted
}

// Next sets value to the next count, increments the count, and
// returns true, unless aborted. When aborted, it is a no-op.
func (c *coreCounter) Next() bool {
	if c.aborted {
		return false
	} else {
		c.value = c.count
		c.count++
		return true
	}
}

// Value returns the current value.
func (c *coreCounter) Value() int {
	return c.value
}

// Abort stops the iterator by setting the aborted flag.
func (c *coreCounter) Abort() {
	c.aborted = true
}

// Reset sets the counter back to 0.
func (c *coreCounter) Reset() {
	c.count = 0
}

// Seq implements the [CoreIterator] method Seq() by delegating to [iterator.Seq].
func (c *coreCounter) Seq() iter.Seq[int] {
	return iterator.Seq(c)
}

// SeqOK returns false since this iterator is not backed by an [iter.Seq], and it's
// slightly more efficient to use Next/Value to consume it.
func (c coreCounter) SeqOK() bool { return false }

// Size returns a value indicating this iterator does not terminate and returns
// an infinite number of items.
func (c coreCounter) Size() iterator.IteratorSize {
	return iterator.SIZE_INFINITE
}

func ExampleNewDefaultIterator() {
	i := iterator.NewDefaultIterator(&coreCounter{})
	func() {
		defer func() { fmt.Println(recover()) }()
		i.Collect() // Attempting to collect the infinite iterator will panic.
	}()
	c := i.Take(10).Collect() // Collecting only the first 10 elements succeeds.
	fmt.Println(c)
	// Output:
	// cannot consume an infinite iterator
	// [0 1 2 3 4 5 6 7 8 9]
}

// prime returns true if n is prime, otherwise false. It tries to find a factor
// by dividing by every number less than itself, and greater than 1.
func prime(n int) bool {
	for f := 2; f < n; f++ {
		if n%f == 0 {
			return false
		}
	}
	return true
}

func ExampleDefaultIterator_Any() {
	fmt.Println(iterator.Range(3, 10).Any(prime))
	fmt.Println(iterator.RangeBy(4, 10, 2).Any(prime))
	// Output:
	// true
	// false
}

func ExampleDefaultIterator_All() {
	fmt.Println(iterator.Range(3, 10).All(prime))
	fmt.Println(iterator.Of(1, 2, 3, 5, 7, 11).All(prime))
	// Output:
	// false
	// true
}

func ExampleFold() {
	add := func(a, b int) int { return a + b }
	s := iterator.Fold(iterator.IncRange(1, 5), 0, add)
	fmt.Println(s)
	// Output: 15
}

func ExampleFold1() {
	mul := func(a, b int) int { return a * b }
	s := iterator.Fold1(iterator.Of(2, 3, 4), mul)
	fmt.Println(s)
	// Output: 24
}

func ExampleIntercalate1() {
	s1 := iterator.Intercalate1(iterator.Of("Hello"), " ", functions.Sum)
	fmt.Println(s1)
	s := iterator.Intercalate1(iterator.Of("Hello", "world"), " ", functions.Sum)
	fmt.Println(s)
	// Output:
	// Hello
	// Hello world
}

func ExampleIntercalate() {
	inputs := []string{"one", "two", "three"}
	for l := range len(inputs) + 1 {
		s := iterator.Intercalate(slices.Iter(inputs[:l]), "", " ", functions.Sum)
		fmt.Printf("%#v\n", s)
	}
	// Output:
	// ""
	// "one"
	// "one two"
	// "one two three"
}

func ExampleDefaultIterator_Intercalate() {
	inputs := []string{"one", "two", "three"}
	for l := range len(inputs) + 1 {
		s := slices.Iter(inputs[:l]).Intercalate("", " ", functions.Sum)
		fmt.Printf("%#v\n", s)
	}
	// Output:
	// ""
	// "one"
	// "one two"
	// "one two three"
}

func ExampleRange() {
	i := iterator.Range(0, 5)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
}

func ExampleRange_descending() {
	i := iterator.Range(5, 0)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 5
	// 4
	// 3
	// 2
	// 1
}

func ExampleRange_empty() {
	c := iterator.Range(0, 0).Collect()
	fmt.Printf("%#v", c)
	// Output: []int{}
}

func ExampleIncRange() {
	i := iterator.IncRange(0, 5)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 0
	// 1
	// 2
	// 3
	// 4
	// 5
}

func ExampleIncRange_single() {
	i := iterator.IncRange(0, 0)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 0
}

func ExampleRangeBy() {
	i := iterator.RangeBy(0, 5, 2)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 0
	// 2
	// 4
}

func ExampleIncRangeBy() {
	i := iterator.IncRangeBy(0, 4, 2)
	for e := range i.Seq() {
		fmt.Println(e)
	}
	// Output:
	// 0
	// 2
	// 4
}

package slices_test

import (
	"fmt"

	"github.com/robdavid/genutil-go/slices"
)

func ExampleIterMut() {
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
	// Output:
	// [0 1 2 3 4]
}

func ExampleCompare() {
	fmt.Println(slices.Compare([]int{1, 2}, []int{1, 3}) < 0)
	fmt.Println(slices.Compare([]int{1, 3}, []int{1, 3}) == 0)
	fmt.Println(slices.Compare([]int{1, 3, 4}, []int{1, 3}) > 0)

	// Output:
	// true
	// true
	// true
}

type ints []int

func ExampleMapAs() {

	sliceIn := []int{1, 2, 3, 4}
	sliceOut := slices.MapAs[ints](sliceIn, func(x int) int { return x * 2 })
	fmt.Printf("%#v", sliceOut)
	// Output:
	// slices_test.ints{2, 4, 6, 8}
}

func ExampleNewAs() {
	myints := slices.NewAs[ints](1, 2, 3, 4)
	fmt.Printf("%#v\n", myints)
	// Output:
	// slices_test.ints{1, 2, 3, 4}
}

func ExampleAffix() {
	myints := slices.NewAs[ints](1, 2, 3, 4)
	extended := slices.Affix(myints, 5, 6, 7)
	fmt.Println(extended)
	// Output:
	// [1 2 3 4 5 6 7]
}

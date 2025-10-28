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

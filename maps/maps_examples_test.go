package maps_test

import (
	"fmt"

	"github.com/robdavid/genutil-go/maps"
)

func ExampleIterMut() {
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
	fmt.Println(m)
	// Output:
	// map[0:5 2:6 4:7 6:8 8:9]
}

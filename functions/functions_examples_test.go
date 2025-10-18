package functions_test

import (
	"fmt"

	"github.com/robdavid/genutil-go/functions"
)

func ExampleSum() {
	fmt.Println(functions.Sum(3, 4))
	fmt.Println(functions.Sum("Hello", " world"))
	// Output:
	// 7
	// Hello world
}

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

func ExampleProduct() {
	p := functions.Product(3+5i, 4+2i)
	fmt.Println(p)
	// Output: (2+26i)
}

func ExampleRef() {
	hp := functions.Ref("hello")
	fmt.Println(*hp)
	// Output: hello
}

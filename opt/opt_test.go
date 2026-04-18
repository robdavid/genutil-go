package opt_test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/opt"
	"github.com/stretchr/testify/assert"
)

func TestOptZero(t *testing.T) {
	var intval opt.Val[int]
	var intref opt.Ref[int]
	assert.True(t, intval.IsEmpty())
	assert.True(t, intref.IsEmpty())
}

func TestOption(t *testing.T) {
	assert := assert.New(t)
	var x int = 123
	intval := opt.Value(x)
	intref := opt.Reference(&x)
	assert.Equal(123, intval.Get())
	assert.Equal(123, intref.Get())
	optval := opt.Option[int](&intval)
	optref := opt.Option[int](intref)
	assert.Equal(123, optval.Get())
	assert.Equal(123, optref.Get())
	*optref.Ref() = 456
	assert.Equal(456, x)
	assert.Equal(123, intval.Get())
	*optval.Ref() = 456
	assert.Equal(456, intval.Get())
}

func TestGetOK(t *testing.T) {
	assert := assert.New(t)

	// Test Val[T] with a value
	var x int = 123
	valWithValue := opt.Value(x)
	value, ok := valWithValue.GetOK()
	assert.Equal(123, value)
	assert.True(ok)

	// Test Val[T] without a value
	var emptyVal opt.Val[int]
	value, ok = emptyVal.GetOK()
	assert.Equal(0, value) // Zero value for int
	assert.False(ok)

	// Test Ref[T] with a reference
	y := 456
	refWithValue := opt.Reference(&y)
	value, ok = refWithValue.GetOK()
	assert.Equal(456, value)
	assert.True(ok)

	// Test Ref[T] without a reference
	var emptyRef opt.Ref[int]
	value, ok = emptyRef.GetOK()
	assert.Equal(0, value) // Zero value for int
	assert.False(ok)
}

func TestRef(t *testing.T) {
	assert := assert.New(t)

	// Test Val[T].Ref() with a value
	valWithValue := opt.Value(123)
	ref := valWithValue.Ref()
	assert.Equal(123, *ref)
	*ref = 789 // Modify the reference
	assert.Equal(789, valWithValue.Get())

	// Test Val[T].Ref() without a value (should panic)
	var emptyVal opt.Val[int]
	assert.Panics(func() {
		_ = emptyVal.Ref()
	})

	// Test Ref[T].Ref() with a reference
	y := 456
	refWithValue := opt.Reference(&y)
	ref = refWithValue.Ref()
	assert.Equal(456, *ref)
	*ref = 999 // Modify the reference
	assert.Equal(999, y)

	// Test Ref[T].Ref() without a reference (should panic)
	var emptyRef opt.Ref[int]
	assert.Panics(func() {
		_ = emptyRef.Ref()
	})
}

func ExampleVal_Try() {
	defer handler.Handle(func(err error) {
		if err == opt.ErrOptionIsEmpty {
			fmt.Println("Access of empty option detected")
		}
	})

	var empty opt.Val[int]
	var present opt.Val[int] = opt.Value(123)

	fmt.Println(present.Try())
	fmt.Println(empty.Try())

	// Output:
	// 123
	// Access of empty option detected
}

func ExampleRef_Try() {
	defer handler.Handle(func(err error) {
		if err == opt.ErrOptionIsEmpty {
			fmt.Println("Access of empty option detected")
		}
	})

	var empty opt.Ref[int]
	var x int = 123
	var present opt.Ref[int] = opt.Reference(&x)

	fmt.Println(present.Try())
	fmt.Println(empty.Try())

	// Output:
	// 123
	// Access of empty option detected
}

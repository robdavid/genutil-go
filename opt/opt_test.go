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
	assert.False(t, intval.HasValue())
	assert.False(t, intref.HasValue())
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

func TestRefOr(t *testing.T) {
	assert := assert.New(t)

	var x int = 123
	valWithValue := opt.Value(x)
	fallbackPtr := &x // Pointer to the same value for testing

	// Test Val[T].RefOr with a fallback pointer
	refFromVal := valWithValue.RefOr(fallbackPtr)
	assert.Equal(&x, refFromVal) // Should return the reference to x

	var emptyVal opt.Val[int]
	nilPtr := (*int)(nil)

	// Test Val[T].RefOr without a value (should return fallback)
	refFromEmptyVal := emptyVal.RefOr(nilPtr)
	assert.Nil(refFromEmptyVal) // Should return nil pointer as fallback

	// Test Ref[T].RefOr with a reference
	y := 456
	refWithValue := opt.Reference(&y)

	// Test Ref[T].RefOr without a reference (should return fallback)
	nilFallback := (*int)(nil)
	refFromEmptyRef := refWithValue.RefOr(nilFallback)
	assert.Equal(&y, refFromEmptyRef) // Should return the reference to y

	var emptyRef opt.Ref[int]
	refFromEmptyRef = emptyRef.RefOr(nilPtr)
	assert.Nil(refFromEmptyRef) // Should return nil pointer as fallback
}

func TestGetOrF(t *testing.T) {
	assert := assert.New(t)

	// Test Val[T].GetOrF with a value present
	var x int = 123
	valWithValue := opt.Value(x)
	result := valWithValue.GetOrF(func() int { return 0 })
	assert.Equal(123, result)

	// Test Val[T].GetOrF without a value (invoke fallback function)
	var emptyVal opt.Val[int]
	result = emptyVal.GetOrF(func() int { return 999 })
	assert.Equal(999, result)

	// Test Ref[T].GetOrF with a reference present
	y := 456
	refWithValue := opt.Reference(&y)
	result = refWithValue.GetOrF(func() int { return 0 })
	assert.Equal(456, result)

	// Test Ref[T].GetOrF without a reference (invoke fallback function)
	var emptyRef opt.Ref[int]
	result = emptyRef.GetOrF(func() int { return 888 })
	assert.Equal(888, result)

	// Test chained GetOrF with side-effect in fallback function
	fallbackCounter := 0
	getValueFunc := func() int {
		fallbackCounter++
		return fallbackCounter * 100
	}
	var emptyVal2 opt.Val[int]
	result = emptyVal2.GetOrF(getValueFunc)
	assert.Equal(100, result)
	assert.Equal(1, fallbackCounter)

	// Verify GetOrF doesn't run fallback when value is present
	fallbackCounter = 0
	valWithValue2 := opt.Value(777)
	result = valWithValue2.GetOrF(getValueFunc)
	assert.Equal(777, result)
	assert.Equal(0, fallbackCounter)
}

func TestString(t *testing.T) {
	val1 := opt.Value(123)
	var x int = 456
	val2 := opt.Reference(&x)
	str := fmt.Sprintf("%s %s", val1, val2)
	assert.Equal(t, "123 456", str)
	str = fmt.Sprintf("%s %s", &val1, &val2)
	assert.Equal(t, "123 456", str)
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

func ExampleVal_AsRef() {
	value := opt.Value(123)
	fmt.Println(value.Get())

	valueRef := value.AsRef()
	*valueRef.Ref() = 456
	fmt.Println(value.Get())

	// Output:
	// 123
	// 456

}

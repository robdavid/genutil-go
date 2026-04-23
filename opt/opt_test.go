//go:build !goexperiment.jsonv2

package opt_test

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	"github.com/robdavid/genutil-go/opt"
	"github.com/stretchr/testify/assert"
)

func checkTry(f func()) (err error) {
	defer handler.Catch(&err)
	f()
	return nil
}

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

func TestGet(t *testing.T) {
	assert := assert.New(t)

	var v opt.Val[int]
	assert.PanicsWithValue(opt.ErrOptionIsEmpty, func() { v.Get() })
	v = opt.Value(123)
	assert.Equal(123, v.Get())

	var r opt.Ref[int]
	assert.PanicsWithValue(opt.ErrOptionIsEmpty, func() { r.Get() })
	var x int = 456
	r = opt.Reference(&x)
	assert.Equal(456, r.Get())
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

func TestRefOK(t *testing.T) {
	assert := assert.New(t)

	// Test Val[T] with a value
	var x int = 123
	valWithValue := opt.Value(x)
	value, ok := valWithValue.RefOK()
	assert.Equal(123, *value)
	assert.True(ok)

	// Test Val[T] without a value
	var emptyVal opt.Val[int]
	value, ok = emptyVal.RefOK()
	assert.Nil(value) // Zero value for int
	assert.False(ok)

	// Test Ref[T] with a reference
	y := 456
	refWithValue := opt.Reference(&y)
	value, ok = refWithValue.RefOK()
	assert.Equal(&y, value)
	assert.True(ok)

	// Test Ref[T] without a reference
	var emptyRef opt.Ref[int]
	value, ok = emptyRef.RefOK()
	assert.Nil(value)
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

func TestGetOr(t *testing.T) {
	assert := assert.New(t)

	empty := opt.Empty[int]()
	v := empty.GetOr(123)
	assert.Equal(123, v)

	value := opt.Value(456)
	v = value.GetOr(123)
	assert.Equal(456, v)

	emptyRef := opt.EmptyRef[int]()
	v = emptyRef.GetOr(123)
	assert.Equal(123, v)

	var x int = 456
	ref := opt.Reference(&x)
	v = ref.GetOr(123)
	assert.Equal(456, v)

}

func TestRefOr(t *testing.T) {
	assert := assert.New(t)

	var x int = 123
	valWithValue := opt.Value(456)
	assert.Equal(456, *valWithValue.RefOr(&x))

	emptyVal := opt.Empty[int]()
	assert.Equal(&x, emptyVal.RefOr(&x))

	var y int = 456
	refWithValue := opt.Reference(&y)
	assert.Equal(&y, refWithValue.RefOr(nil))

	emptyRef := opt.EmptyRef[int]()
	assert.Equal(&y, emptyRef.RefOr(&y)) // Should return nil pointer as fallback
}

func TestGetOrF(t *testing.T) {
	assert := assert.New(t)

	valWithValue := opt.Value(123)
	result := valWithValue.GetOrF(func() int { return 0 })
	assert.Equal(123, result)

	emptyVal := opt.Empty[int]()
	result = emptyVal.GetOrF(func() int { return 999 })
	assert.Equal(999, result)

	y := 456
	refWithValue := opt.Reference(&y)
	result = refWithValue.GetOrF(func() int { return 0 })
	assert.Equal(456, result)

	emptyRef := opt.EmptyRef[int]()
	result = emptyRef.GetOrF(func() int { return 888 })
	assert.Equal(888, result)
}

func TestEqual(t *testing.T) {
	assert := assert.New(t)
	val1 := opt.Value(123)
	val2 := val1
	val3 := opt.Value(456)
	val4 := opt.Empty[int]()
	assert.False(val1.Ref() == val2.Ref())
	assert.True(val1 == val2)
	assert.False(val1 == val3)
	assert.False(val1 == val4)
	ref1 := val1.AsRef()
	ref2 := val2.AsRef()
	assert.False(ref1 == ref2) // Pointers will not be equal
	assert.True(opt.Equal(ref1, ref2))
	assert.True(opt.Equal(&val1, &val2))
	assert.True(opt.Equal(&val1, ref2))
	assert.False(opt.Equal(&val1, &val3))
}

func TestDeepEqual(t *testing.T) {
	assert := assert.New(t)
	type testdata struct {
		name   string
		values []int
	}
	val1 := opt.Value(testdata{
		name:   "val1",
		values: []int{1, 2, 3},
	})
	val2 := opt.Value(testdata{
		name:   "val1",
		values: []int{1, 2, 3},
	})
	val3 := opt.Value(testdata{
		name:   "val1",
		values: []int{4, 5, 6},
	})
	assert.False(val1.Ref() == val2.Ref())
	assert.True(opt.DeepEqual(&val1, &val2))
	assert.True(opt.DeepEqual(&val1, val2.AsRef()))
	assert.False(opt.DeepEqual(&val1, &val3))
}

func TestString(t *testing.T) {
	val1 := opt.Value(123)
	var x int = 456
	val2 := opt.Reference(&x)
	str := fmt.Sprintf("%s %s", val1, val2)
	assert.Equal(t, "123 456", str)
	str = fmt.Sprintf("%s %s", &val1, &val2)
	assert.Equal(t, "123 456", str)
	empty := opt.Empty[int]()
	emptyRef := opt.Empty[int]()
	emptyStr := fmt.Sprintf("%s-%s", empty, emptyRef)
	assert.Equal(t, "-", emptyStr)
	emptyStr = fmt.Sprintf("%s-%s", &empty, &emptyRef)
	assert.Equal(t, "-", emptyStr)
}

func TestTryRef(t *testing.T) {
	assert := assert.New(t)

	emptyVar := opt.Empty[int]()
	err := checkTry(func() { emptyVar.TryRef() })
	assert.ErrorIs(err, opt.ErrOptionIsEmpty)

	nonemptyVar := opt.Value(123)
	assert.Equal(123, *nonemptyVar.TryRef())

	emptyRef := opt.EmptyRef[int]()
	err = checkTry(func() { emptyRef.TryRef() })
	assert.ErrorIs(err, opt.ErrOptionIsEmpty)

	var x int = 456
	nonemptyRef := opt.Reference(&x)
	assert.Equal(&x, nonemptyRef.TryRef())
}

func TestMorph(t *testing.T) {
	assert := assert.New(t)
	v := opt.Value("hello")
	uv := v.Morph(strings.ToUpper)
	assert.Equal("HELLO", uv.Get())

	v = opt.Empty[string]()
	uv = v.Morph(strings.ToUpper)
	assert.True(uv.IsEmpty())

	var hello string = "hello"
	upperRef := func(s *string) *string {
		upper := strings.ToUpper(*s)
		return &upper
	}
	r := opt.Reference(&hello)
	ur := r.MorphRef(upperRef)
	assert.Equal("HELLO", ur.Get())

	r = opt.Reference(&hello)
	assert.Equal("hello", r.Get())
	uv = r.Morph(strings.ToUpper)
	assert.Equal("HELLO", uv.Get())

	r = opt.EmptyRef[string]()
	ur = r.MorphRef(upperRef)
	assert.True(ur.IsEmpty())
	uv = r.Morph(strings.ToUpper)
	assert.True(ur.IsEmpty())

	v = opt.Value("hello")
	ur = v.MorphRef(upperRef)
	assert.Equal("HELLO", ur.Get())

	v = opt.Empty[string]()
	ur = v.MorphRef(upperRef)
	assert.True(ur.IsEmpty())
}

func TestThenElse(t *testing.T) {
	assert := assert.New(t)

	var thenValue int
	var elseTaken bool = false

	v := opt.Value(123)
	v.Then(func(x int) { thenValue = x }).Else(func() { elseTaken = true })
	assert.Equal(123, thenValue)
	assert.False(elseTaken)

	v = opt.Empty[int]()
	v.Then(func(x int) { thenValue = x }).Else(func() { elseTaken = true })
	assert.Equal(123, thenValue)
	assert.True(elseTaken)

	var x int = 456
	elseTaken = false
	r := opt.Reference(&x)
	r.Then(func(x int) { thenValue = x }).Else(func() { elseTaken = true })
	assert.Equal(456, thenValue)
	assert.False(elseTaken)

	r = opt.EmptyRef[int]()
	r.Then(func(x int) { thenValue = x }).Else(func() { elseTaken = true })
	assert.Equal(456, thenValue)
	assert.True(elseTaken)

}

func TestThenRefElse(t *testing.T) {
	assert := assert.New(t)

	var thenValue int
	var elseTaken bool = false

	v := opt.Value(123)
	v.ThenRef(func(x *int) { thenValue = *x }).Else(func() { elseTaken = true })
	assert.Equal(123, thenValue)
	assert.False(elseTaken)

	v = opt.Empty[int]()
	v.ThenRef(func(x *int) { thenValue = *x }).Else(func() { elseTaken = true })
	assert.Equal(123, thenValue)
	assert.True(elseTaken)

	var x int = 456
	elseTaken = false
	r := opt.Reference(&x)
	r.ThenRef(func(x *int) { thenValue = *x }).Else(func() { elseTaken = true })
	assert.Equal(456, thenValue)
	assert.False(elseTaken)

	r = opt.EmptyRef[int]()
	r.ThenRef(func(x *int) { thenValue = *x }).Else(func() { elseTaken = true })
	assert.Equal(456, thenValue)
	assert.True(elseTaken)

}

func TestMap(t *testing.T) {
	assert := assert.New(t)

	v := opt.Value(123)
	ev := opt.Empty[int]()
	assert.Equal("123", opt.Map(&v, strconv.Itoa).Get())
	assert.True(opt.Map(&ev, strconv.Itoa).IsEmpty())

	var x int = 123
	r := opt.Reference(&x)
	er := opt.EmptyRef[int]()
	assert.Equal("123", opt.Map(r, strconv.Itoa).Get())
	assert.True(opt.Map(er, strconv.Itoa).IsEmpty())
}

func TestMapRef(t *testing.T) {
	assert := assert.New(t)

	itoaRef := func(x *int) *string {
		str := strconv.Itoa(*x)
		return &str
	}
	v := opt.Value(123)
	ev := opt.Empty[int]()
	assert.Equal("123", opt.MapRef(&v, itoaRef).Get())
	assert.True(opt.MapRef(&ev, itoaRef).IsEmpty())

	var x int = 123
	r := opt.Reference(&x)
	er := opt.EmptyRef[int]()
	assert.Equal("123", opt.MapRef(r, itoaRef).Get())
	assert.True(opt.MapRef(er, itoaRef).IsEmpty())
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

func ExampleVal_Mutate() {
	type mystruct struct {
		name  string
		value int
	}
	v := opt.Empty[mystruct]()
	v2 := v.Ensure().Mutate(func(m *mystruct) {
		m.name = "two"
		m.value = 2
	})
	fmt.Printf("%s: %d\n", v2.Ref().name, v2.Ref().value)

	// Output:
	// two: 2
}

func ExampleRef_Mutate() {
	type mystruct struct {
		name  string
		value int
	}
	v := opt.EmptyRef[mystruct]()
	v2 := v.Ensure().Mutate(func(m *mystruct) {
		m.name = "two"
		m.value = 2
	})
	fmt.Printf("%s: %d\n", v2.Ref().name, v2.Ref().value)

	// Output:
	// two: 2
}

func ExampleVal_Ensure() {
	type mystruct struct {
		name  string
		value int
	}
	v := opt.Empty[mystruct]()
	v2 := v.Ensure().Mutate(func(m *mystruct) {
		m.name = "two"
		m.value = 2
	})
	fmt.Printf("%s: %d\n", v2.Ref().name, v2.Ref().value)

	// Output:
	// two: 2
}

func ExampleRef_Ensure() {
	type mystruct struct {
		name  string
		value int
	}
	v := opt.EmptyRef[mystruct]()
	v2 := v.Ensure().Mutate(func(m *mystruct) {
		m.name = "two"
		m.value = 2
	})
	fmt.Printf("%s: %d\n", v2.Ref().name, v2.Ref().value)

	// Output:
	// two: 2
}

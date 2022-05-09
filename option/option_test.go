package option

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestValue(t *testing.T) {
	var v Option[int] = Value(5)
	assert.False(t, v.IsEmpty())
	assert.True(t, v.HasValue())
	assert.Equal(t, 5, v.Get())
	assert.Equal(t, 7, Map[int](v, func(x int) int { return x + 2 }).Get())
	assert.Equal(t, 8, MapRef[int](&v, func(x *int) *int { var y = *x + 3; return &y }).Get())
	v.Set(6)
	assert.Equal(t, 6, v.Get())
	assert.Equal(t, 6, *v.Ref())
	v.Clear()
	assert.True(t, v.IsEmpty())
}

func TestRef(t *testing.T) {
	five := 5
	var v OptionRef[int] = Ref(&five)
	assert.False(t, v.IsEmpty())
	assert.True(t, v.HasValue())
	assert.Equal(t, 5, v.Get())
	assert.Equal(t, &five, v.Ref())
	assert.Equal(t, 7, MapRef[int](&v, func(x *int) *int { var y = *x + 2; return &y }).Get())
	assert.Equal(t, 8, Map[int](v, func(x int) int { return x + 3 }).Get())
	v.Set(6)
	assert.Equal(t, 6, v.Get())
	assert.Equal(t, 6, *v.Ref())
	v.Clear()
	assert.True(t, v.IsEmpty())
}

func TestEquality(t *testing.T) {
	six := 6
	val := Value(six)
	ref := Ref(&six)
	assert.True(t, Equal[int](val, ref))
	val.Set(7)
	assert.False(t, Equal[int](val, ref))
	val.Set(6)
	assert.True(t, Equal[int](val, ref))
	val.Clear()
	assert.False(t, Equal[int](val, ref))
}

func TestRefEquality(t *testing.T) {
	six := 6
	val := Value(six)
	ref := Ref(&six)
	assert.True(t, EqualRef[int](&val, &ref))
	val.Set(7)
	assert.False(t, EqualRef[int](&val, &ref))
	val.Set(6)
	assert.True(t, EqualRef[int](&val, &ref))
	val.Clear()
	assert.False(t, EqualRef[int](&val, &ref))
}
func TestEmpty(t *testing.T) {
	v := Empty[int]()
	assert.True(t, v.IsEmpty())
	assert.False(t, v.HasValue())
	assert.Equal(t, 0, v.GetOrZero())
	vm := Map[int](v, func(x int) int { return x + 2 })
	assert.True(t, vm.IsEmpty())
}

type Wrapper[T any] struct {
	Wrapped IOptionRef[T]
}

func TestInterfaceAssignment(t *testing.T) {
	var io IOption[int]
	var ior IOptionRef[int]
	var wrapper Wrapper[int]
	o := Value(123)
	io = o
	ior = &o
	wrapper = Wrapper[int]{&o}
	assert.Equal(t, io.Get(), 123)
	assert.Equal(t, ior.Get(), 123)
	assert.Equal(t, wrapper.Wrapped.Get(), 123)
}

type TestS1 struct {
	name  string
	value int
}

// It should be possile to copy an option
// without being exposed to hidden references
func TestSafeCopy(t *testing.T) {
	t1 := TestS1{"one", 1}
	v1 := Value(t1)
	v2 := v1
	v2.Ref().name = "two"
	v2.Ref().value = 2
	assert.Equal(t, TestS1{"two", 2}, v2.Get())
	assert.Equal(t, TestS1{"one", 1}, v1.Get())
}

// OptionRefs refer to the original item when copied
func TestRefCopy(t *testing.T) {
	t1 := TestS1{"one", 1}
	v1 := Ref(&t1)
	assert.Equal(t, TestS1{"one", 1}, v1.Get())
	v2 := v1
	v2.Ref().name = "two"
	v2.Ref().value = 2
	assert.Equal(t, TestS1{"two", 2}, v2.Get())
	assert.Equal(t, TestS1{"two", 2}, v1.Get())
}

func TestOptionPtr(t *testing.T) {
	var opt OptionRef[int]
	assert.True(t, opt.IsEmpty())
	opt.SetRef(nil)
	assert.True(t, opt.IsEmpty())
	r, ok := opt.GetOK()
	assert.False(t, ok)
	assert.Zero(t, r)
	v := 123
	opt.SetRef(&v)
	assert.False(t, opt.IsEmpty())
	r, ok = opt.GetOK()
	assert.True(t, ok)
	assert.Equal(t, r, 123)
}

func TestOptionInterface(t *testing.T) {
	var opti IOptionRef[int]
	ref := Ref[int](nil)
	opti = &ref
	assert.True(t, opti.IsEmpty())
	r, ok := opti.GetOK()
	assert.False(t, ok)
	assert.Zero(t, r)
	v := 123
	opt := Value(v)
	opti = &opt
	assert.False(t, opti.IsEmpty())
	r, ok = opti.GetOK()
	assert.True(t, ok)
	assert.Equal(t, r, 123)
}

func TestOptionList(t *testing.T) {
	opt := Value[[]int](nil)
	//assert.True(t, opt.IsEmpty())
	opt.Set(append(opt.GetOrZero(), 1))
	assert.Equal(t, []int{1}, opt.Get())
}

type testMarshall struct {
	Name  string          `json:"name" yaml:"name"`
	Value int             `json:"value" yaml:"value"`
	Opt   Option[testOpt] `json:"opt" yaml:"opt,omitempty"`
}

type testOpt struct {
	Metadata string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ItemList []string `json:"itemList,omitempty" yaml:"itemList,omitempty"`
}

func TestPrintFormatting(t *testing.T) {
	actual := fmt.Sprintf("String is %s and num is %s", Value("5 by 5"), Value(25))
	assert.Equal(t, "String is 5 by 5 and num is 25", actual)
	actual = fmt.Sprintf("String is %s and num is %s", Empty[string](), Empty[int]())
	assert.Equal(t, "String is  and num is ", actual)
}

func TestJSONMarshallOmitOption(t *testing.T) {
	testData := testMarshall{
		Name:  "test1",
		Value: 1,
		Opt:   Empty[testOpt](),
	}
	y, err := json.Marshal(&testData)
	if assert.NoError(t, err) {
		text := string(y)
		assert.Contains(t, text, "\"opt\":null")
		var testData2 testMarshall
		assert.NoError(t, json.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

func TestYAMLMarshallOption(t *testing.T) {
	testData := testMarshall{
		Name:  "test1",
		Value: 1,
		Opt:   Value(testOpt{Metadata: "md", ItemList: []string{"item1"}}),
	}
	y, err := yaml.Marshal(&testData)
	if assert.NoError(t, err) {
		text := string(y)
		assert.Contains(t, text, "opt:")
		assert.Contains(t, text, "metadata:")
		assert.Contains(t, text, "itemList:")
		var testData2 testMarshall
		assert.NoError(t, yaml.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

func TestYAMLMarshallOmitOption(t *testing.T) {
	testData := testMarshall{
		Name:  "test1",
		Value: 1,
		Opt:   Empty[testOpt](),
	}
	y, err := yaml.Marshal(&testData)
	if assert.NoError(t, err) {
		text := string(y)
		assert.NotContains(t, text, "opt:")
		assert.NotContains(t, text, "metadata:")
		assert.NotContains(t, text, "itemList:")
		var testData2 testMarshall
		assert.NoError(t, yaml.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

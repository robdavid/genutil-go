package option

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/errors/handler"
	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestValue(t *testing.T) {
	var v Option[int] = Value(5)
	assert.False(t, v.IsEmpty())
	assert.True(t, v.HasValue())
	assert.Equal(t, 5, v.Get())
	assert.Equal(t, Value(7), Map(v, func(x int) int { return x + 2 }))
	assert.Equal(t, 8, MapRef(&v, func(x *int) *int { var y = *x + 3; return &y }).Get())
	v.Set(6)
	assert.Equal(t, 6, v.Get())
	assert.Equal(t, 6, *v.Ref())
	v.Clear()
	assert.True(t, v.IsEmpty())
}

func TestRef(t *testing.T) {
	five := 5
	var v *Option[int] = Ref(&five)
	assert.False(t, v.IsEmpty())
	assert.True(t, v.HasValue())
	assert.Equal(t, 5, v.Get())
	assert.Equal(t, &five, v.Ref())
	assert.Equal(t, 7, MapRef(v, func(x *int) *int { var y = *x + 2; return &y }).Get())
	assert.Equal(t, Value(8), Map(*v, func(x int) int { return x + 3 }))
	v.Set(6)
	assert.Equal(t, 6, v.Get())
	assert.Equal(t, 6, *v.Ref())
	v.Clear()
	assert.True(t, v.IsEmpty())
}

func TestToRefExample(t *testing.T) {
	var slice []int = []int{6, 7}
	append42 := func(s *[]int) { *s = append(*s, 42) }
	opt := Value(slice).ToRef().Mutate(append42).Get() // []int{6, 7, 42}
	assert.Equal(t, []int{6, 7, 42}, opt)
}

func TestEquality(t *testing.T) {
	six := 6
	val := Value(six)
	ref := Ref(&six)
	assert.True(t, val == *ref)
	val.Set(7)
	assert.False(t, val == *ref)
	val.Set(6)
	assert.True(t, val == *ref)
	val.Clear()
	assert.False(t, val == *ref)
}

func TestEmpty(t *testing.T) {
	v := Empty[int]()
	assert.True(t, v.IsEmpty())
	assert.False(t, v.HasValue())
	assert.Equal(t, 0, v.GetOrZero())
	vm := Map(v, func(x int) int { return x + 2 })
	assert.True(t, vm.IsEmpty())
}

func TestZero(t *testing.T) {
	var zero Option[int]
	assert.True(t, zero.IsEmpty())
}

func TestNewStruct(t *testing.T) {
	opt := New[struct {
		num  int
		text string
	}]()
	opt.Ref().num = 123
	opt.Ref().text = "one two three"
	assert.Equal(t, 123, opt.Get().num)
	assert.Equal(t, "one two three", opt.Get().text)
}

func TestSafe(t *testing.T) {
	var myInt int = 32
	v := Safe[*int](nil)
	assert.True(t, v.IsEmpty())
	v = Safe(&myInt)
	assert.False(t, v.IsEmpty())
	vList := From[[]int](nil)
	assert.True(t, vList.IsEmpty())
	vList = Safe([]int{1, 2, 3})
	assert.False(t, vList.IsEmpty())
	var interf any = nil
	v3 := Safe(interf)
	assert.True(t, v3.IsEmpty())
	interf = vList
	v3.SafeSet(interf)
	assert.False(t, v3.IsEmpty())
	var stringInt map[string]int
	v4 := Safe(stringInt)
	assert.True(t, v4.IsEmpty())
	interf = stringInt
	v3 = Safe(interf)
	assert.True(t, v3.IsEmpty())
	stringInt = make(map[string]int)
	v4 = Safe(stringInt)
	assert.False(t, v4.IsEmpty())
	var double func(int) int
	v5 := Safe(double)
	assert.True(t, v5.IsEmpty())
	double = func(x int) int { return x * 2 }
	v5.SafeSet(double)
	assert.False(t, v5.IsEmpty())
	var ch chan int
	v6 := Safe(ch)
	assert.True(t, v6.IsEmpty())
	ch = make(chan int)
	v6.SafeSet(ch)
	assert.False(t, v6.IsEmpty())
	var chin chan<- int
	v7 := Safe(chin)
	assert.True(t, v7.IsEmpty())
	chin = ch
	v7.SafeSet(chin)
	assert.False(t, v7.IsEmpty())
}

func BenchmarkSafeInit(b *testing.B) {
	var ptr []int
	for i := 0; i < b.N; i++ {
		v := Safe(ptr)
		assert.False(b, v.IsEmpty())
	}
}

func BenchmarkFromInit(b *testing.B) {
	var ptr *int
	for i := 0; i < b.N; i++ {
		v := From(ptr)
		assert.True(b, v.IsEmpty())
	}
}

func BenchmarkValueInit(b *testing.B) {
	var ptr *int
	for i := 0; i < b.N; i++ {
		v := Value(ptr)
		assert.False(b, v.IsEmpty())
	}
}

type TestS1 struct {
	name  string
	value int
}

// It should be possible to copy an option
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

func TestOptionPtr(t *testing.T) {
	var opt Option[int]
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

func TestOptionList(t *testing.T) {
	opt := Value[[]int](nil)
	//assert.True(t, opt.IsEmpty())
	opt.Set(append(opt.GetOrZero(), 1))
	assert.Equal(t, []int{1}, opt.Get())
}

func tryOption(o Option[int], errE error, errF func() error) (result int, err error) {
	defer handler.Catch(&err) // err set to ErrOptionIsEmpty
	if errF != nil {
		result = o.TryErrF(errF)
	} else if errE != nil {
		result = o.TryErr(errE)
	} else {
		result = o.Try()
	}
	return
}

func TestOptionTry(t *testing.T) {
	assert := assert.New(t)
	var actual int
	var err error
	actual, err = tryOption(Value(123), nil, nil)
	assert.Equal(123, actual)
	assert.NoError(err)
	_, err = tryOption(Empty[int](), nil, nil)
	assert.ErrorIs(err, ErrOptionIsEmpty)
}

func TestOptionTryErr(t *testing.T) {
	assert := assert.New(t)
	var actual int
	var err error
	errTest := errors.New("test error")
	actual, err = tryOption(Value(123), errTest, nil)
	assert.Equal(123, actual)
	assert.NoError(err)
	_, err = tryOption(Empty[int](), errTest, nil)
	assert.ErrorIs(err, errTest)
}

func TestOptionTryErrF(t *testing.T) {
	assert := assert.New(t)
	var actual int
	var err error
	errTest := errors.New("test error")
	invoked := false
	fnErr := func() error { invoked = true; return errTest }
	actual, err = tryOption(Value(123), nil, fnErr)
	assert.Equal(123, actual)
	assert.NoError(err)
	assert.False(invoked)
	_, err = tryOption(Empty[int](), nil, fnErr)
	assert.ErrorIs(err, errTest)
	assert.True(invoked)
}

func TestOptionPointerTry(t *testing.T) {
	var err error
	defer func() {
		assert.ErrorIs(t, err, ErrOptionIsEmpty)
	}()
	defer eh.Catch(&err)
	eo := Safe[*int](nil)
	assert.Equal(t, 0, eo.Try())
}

func TestMapAndMorph(t *testing.T) {
	assert := assert.New(t)
	opt := Value(123)
	nopt := Empty[int]()
	fdouble := func(n int) float64 { return 2.0 * float64(n) }
	idouble := func(n int) int { return 2 * n }
	assert.Equal(float64(246), Map(opt, fdouble).Get())
	assert.Equal(Value(246), opt.Morph(idouble))
	assert.Equal(float64(0), Map(nopt, fdouble).GetOrZero())
	assert.Equal(0, nopt.Morph(idouble).GetOr(0))
}

func TestMutate(t *testing.T) {
	assert := assert.New(t)
	opt := Value(123)
	nopt := Empty[int]()
	o2 := opt.Mutate(func(n *int) { *n *= 2 })
	assert.Equal(246, opt.Get())
	assert.Equal(&opt, o2)
	nopt.Mutate(func(n *int) { *n *= 2 })
	assert.Equal(0, nopt.GetOrZero())
}

func TestCompare(t *testing.T) {
	assert := assert.New(t)
	type nv struct{ name, value string }
	var opt1, opt2 Option[nv]
	assert.True(opt1 == opt2)
	opt1.Ensure()
	assert.False(opt1 == opt2)
	opt2.Ensure()
	assert.True(opt1 == opt2)
	opt1.Set(nv{"name", "value"})
	assert.False(opt1 == opt2)
	opt2.Set(nv{"name", "value"})
	assert.True(opt1 == opt2)
}

func TestMutateFromEmpty(t *testing.T) {
	assert := assert.New(t)
	type nv struct{ name, value string }
	opt := Empty[nv]()
	opt.Ensure().Mutate(func(n *nv) {
		n.name = "name"
		n.value = "value"
	})
	assert.Equal(Value(nv{"name", "value"}), opt)
}

func TestThenElse(t *testing.T) {
	const (
		none int = iota
		less
		more
	)
	condition := func(opt Option[int]) (result int) {
		opt.Then(func(n int) {
			if n < 100 {
				result = less
			} else {
				result = more
			}
		}).Else(func() { result = none })
		return
	}
	assert := assert.New(t)
	assert.Equal(none, condition(Empty[int]()))
	assert.Equal(less, condition(Value(50)))
	assert.Equal(more, condition(Value(150)))
}

type testMarshal struct {
	Name  string          `json:"name" yaml:"name"`
	Value int             `json:"value" yaml:"value"`
	Opt   Option[testOpt] `json:"opt" yaml:"opt,omitempty"`
}

type testOpt struct {
	Metadata string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ItemList []string `json:"itemList,omitempty" yaml:"itemList,omitempty"`
}

type testOptMarshal struct {
	Name  Option[string] `json:"name,omitempty" yaml:"name,omitempty"`
	Value Option[int]    `json:"value,omitempty" yaml:"value,omitempty"`
}

type testOptNoOmitMarshal struct {
	Name  Option[string] `json:"name" yaml:"name"`
	Value Option[int]    `json:"value" yaml:"value"`
}

func TestPrintFormatting(t *testing.T) {
	actual := fmt.Sprintf("String is %s and num is %s", Value("5 by 5"), Value(25))
	assert.Equal(t, "String is 5 by 5 and num is 25", actual)
	actual = fmt.Sprintf("String is %s and num is %s", Empty[string](), Empty[int]())
	assert.Equal(t, "String is  and num is ", actual)
}

func TestPrintExample(t *testing.T) {
	actual := fmt.Sprintf("Hello %s", Value("world"))
	assert.Equal(t, "Hello world", actual)
}

func TestJSONMarshalOmitOption(t *testing.T) {
	testData := testMarshal{
		Name:  "test1",
		Value: 1,
		Opt:   Empty[testOpt](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(t, err)
	text := string(y)
	assert.Contains(t, text, "\"opt\":null")
	var testData2 testMarshal
	assert.NoError(t, json.Unmarshal(y, &testData2))
	assert.Equal(t, testData, testData2)
}

func TestJSONMarshalOption(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testMarshal{
		Name:  "test1",
		Value: 1,
		Opt:   Value(testOpt{"Hello", nil}),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"opt\":{\"metadata\":\"Hello\"}")
	var testDataMap map[string]any
	require.NoError(json.Unmarshal(y, &testDataMap))
	expected := map[string]any{
		"name":  "test1",
		"value": float64(1),
		"opt": map[string]any{
			"metadata": "Hello",
		},
	}
	assert.Equal(expected, testDataMap)
}

func TestJSONMarshalEmptyOption(t *testing.T) {
	require := require.New(t)
	testData := testMarshal{
		Name:  "test1",
		Value: 1,
		Opt:   Empty[testOpt](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	var testDataMap map[string]any
	require.NoError(json.Unmarshal(y, &testDataMap))
	expected := map[string]any{
		"name":  "test1",
		"value": float64(1),
		"opt":   nil,
	}
	assert.Equal(t, expected, testDataMap)
}

func TestJSONMarshalOptionSimple(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshal{
		Name:  Value("a name"),
		Value: Empty[int](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	var testDataMap map[string]any
	require.NoError(json.Unmarshal(y, &testDataMap))
	expected := map[string]any{
		"name":  "a name",
		"value": nil,
	}
	assert.Equal(expected, testDataMap)
	var testData2 testOptMarshal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONMarshalOptionSimpleNoEmpty(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshal{
		Name:  Value("a name"),
		Value: Value(123),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":123")
	var testData2 testOptMarshal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONUnMarshalOmitOption(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshal
	var testData = `{ "name": "a name" }`
	err := json.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

func TestJSONMarshalNoOmitOption(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptNoOmitMarshal{
		Name:  Value("a name"),
		Value: Empty[int](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":null")
	var testData2 testOptNoOmitMarshal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLMarshalOption(t *testing.T) {
	testData := testMarshal{
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
		var testData2 testMarshal
		assert.NoError(t, yaml.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

func TestYAMLMarshalOmitOption(t *testing.T) {
	testData := testMarshal{
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
		var testData2 testMarshal
		assert.NoError(t, yaml.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

func TestYAMLMarshalOmitSimpleOption(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshal{
		Name:  Value("a name"),
		Value: Empty[int](),
	}
	y, err := yaml.Marshal(&testData)

	require.NoError(err)
	text := string(y)
	assert.Contains(text, "name:")
	assert.NotContains(text, "value:")
	var testData2 testOptMarshal
	require.NoError(yaml.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLUnMarshalOmitOption(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshal
	var testData = `name: "a name"`
	err := yaml.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

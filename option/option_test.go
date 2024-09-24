package option

import (
	"encoding/json"
	"fmt"
	"testing"

	eh "github.com/robdavid/genutil-go/errors/handler"
	"github.com/stretchr/testify/assert"
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

func TestAs(t *testing.T) {
	var v any = 123
	opt := As[int](v).Get()
	assert.Equal(t, 123, opt)
	assert.Equal(t, []any(nil), As[[]any](v).GetOr(nil))
}

func TestAsRef(t *testing.T) {
	var n = 123
	var v any = &n
	var i *int = nil
	opt := AsRef[int](&n).Get()
	assert.Equal(t, 123, opt)
	assert.True(t, AsRef[int](i).IsEmpty())
	i = AsRef[int](v).GetRef()
	assert.Equal(t, 123, AsRef[int](i).Get())
}

func TestToRefExample(t *testing.T) {
	var slice any = []int{1, 2, 3}
	opt := As[[]int](slice).ToRef().GetRef()
	assert.Equal(t, []int{1, 2, 3}, *opt)
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

func TestSafe(t *testing.T) {
	var myInt int = 32
	v := Safe[*int](nil)
	assert.True(t, v.IsEmpty())
	v = Safe(&myInt)
	assert.False(t, v.IsEmpty())
	vList := Safe[[]int](nil)
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

func TestOptionTry(t *testing.T) {
	var err error
	defer func() {
		assert.ErrorIs(t, err, ErrOptionIsEmpty)
	}()
	defer eh.Catch(&err)
	eo := Empty[int]()
	assert.Equal(t, 0, eo.Try())
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

type testMarshall struct {
	Name  string          `json:"name" yaml:"name"`
	Value int             `json:"value" yaml:"value"`
	Opt   Option[testOpt] `json:"opt" yaml:"opt,omitempty"`
}

type testOpt struct {
	Metadata string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ItemList []string `json:"itemList,omitempty" yaml:"itemList,omitempty"`
}

type testOptMarshall struct {
	Name  Option[string] `json:"name" yaml:"name"`
	Value Option[int]    `json:"value" yaml:"value"`
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

func TestJSONMarshallOption(t *testing.T) {
	testData := testMarshall{
		Name:  "test1",
		Value: 1,
		Opt:   Value(testOpt{"Hello", nil}),
	}
	y, err := json.Marshal(&testData)
	if assert.NoError(t, err) {
		text := string(y)
		assert.Contains(t, text, "\"opt\":{\"metadata\":\"Hello\"}")
		var testData2 testMarshall
		assert.NoError(t, json.Unmarshal(y, &testData2))
		assert.Equal(t, testData, testData2)
	}
}

func TestJSONUnMarshallOmitOption(t *testing.T) {
	var unmarshalled testOptMarshall
	var testData = `{ "name": "a name" }`
	err := json.Unmarshal([]byte(testData), &unmarshalled)
	if assert.NoError(t, err) {
		assert.True(t, unmarshalled.Name.HasValue())
		assert.False(t, unmarshalled.Value.HasValue())
		assert.Equal(t, "a name", unmarshalled.Name.Get())
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

func TestYAMLUnMarshallOmitOption(t *testing.T) {
	var unmarshalled testOptMarshall
	var testData = `name: "a name"`
	err := yaml.Unmarshal([]byte(testData), &unmarshalled)
	if assert.NoError(t, err) {
		assert.True(t, unmarshalled.Name.HasValue())
		assert.False(t, unmarshalled.Value.HasValue())
		assert.Equal(t, "a name", unmarshalled.Name.Get())
	}
}

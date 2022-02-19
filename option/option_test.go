package option

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValue(t *testing.T) {
	v := Value(5)
	assert.False(t, v.IsEmpty())
	assert.True(t, v.HasValue())
	assert.Equal(t, 5, v.Get())
	assert.Equal(t, 7, Map(v, func(x int) int { return x + 2 }).Get())
	v.Set(6)
	assert.Equal(t, 6, v.Get())
	assert.Equal(t, 6, *v.Ref())
	v.Clear()
	assert.True(t, v.IsEmpty())
}

func TestEmpty(t *testing.T) {
	v := Empty[int]()
	assert.True(t, v.IsEmpty())
	assert.False(t, v.HasValue())
	assert.Equal(t, 0, v.GetOrZero())
	vm := Map(v, func(x int) int { return x + 2 })
	assert.True(t, vm.IsEmpty())
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
	assert.Equal(t, v2.Get(), TestS1{"two", 2})
	assert.Equal(t, v1.Get(), TestS1{"one", 1})
}

func TestOptionPtr(t *testing.T) {
	var opt Option[*int]
	assert.True(t, opt.IsEmpty())
	opt.Set(nil)
	assert.True(t, opt.IsEmpty())
	r, ok := opt.GetOK()
	assert.False(t, ok)
	assert.Nil(t, r)
	v := 123
	opt.Set(&v)
	assert.False(t, opt.IsEmpty())
	r, ok = opt.GetOK()
	assert.True(t, ok)
	assert.Equal(t, *r, 123)
}

func TestOptionList(t *testing.T) {
	opt := Value[[]int](nil)
	assert.True(t, opt.IsEmpty())
	opt.Set(append(opt.GetOrZero(), 1))
	assert.Equal(t, []int{1}, opt.Get())
}

type testMarshall struct {
	Name  string          `json:"name"`
	Value int             `json:"value"`
	Opt   Option[testOpt] `json:"opt,omitempty"`
}

type testOpt struct {
	Metadata string   //`json:"metadata,omitempty"`
	ItemList []string //`json:"itemList,omitempty"`
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

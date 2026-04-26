package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/robdavid/genutil-go/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type testMarshalVal struct {
	Name  string              `json:"name" yaml:"name"`
	Value int                 `json:"value" yaml:"value"`
	Opt   opt.Val[testOptVal] `json:"opt" yaml:"opt"`
}

type testOptVal struct {
	Metadata string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ItemList []string `json:"itemList,omitempty" yaml:"itemList,omitempty"`
}

type testOptNoOmitMarshalVal struct {
	Name  opt.Val[string] `json:"name" yaml:"name"`
	Value opt.Val[int]    `json:"value" yaml:"value"`
}

type testOptMarshalVal struct {
	Name  opt.Val[string] `json:"name,omitzero" yaml:"name,omitempty"`
	Value opt.Val[int]    `json:"value,omitzero" yaml:"value,omitempty"`
}

func TestJSONMarshalOmitVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testMarshalVal{
		Name:  "test1",
		Value: 1,
		Opt:   opt.Empty[testOptVal](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"opt\":null")
	var testData2 testMarshalVal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONMarshalVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testMarshalVal{
		Name:  "test1",
		Value: 1,
		Opt:   opt.Value(testOptVal{"Hello", nil}),
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

func TestJSONMarshalEmptyVal(t *testing.T) {
	require := require.New(t)
	testData := testMarshalVal{
		Name:  "test1",
		Value: 1,
		Opt:   opt.Empty[testOptVal](),
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

func TestJSONMarshalValSimple(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalVal{
		Name:  opt.Value("a name"),
		Value: opt.Empty[int](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	var testDataMap map[string]any
	require.NoError(json.Unmarshal(y, &testDataMap))
	expected := map[string]any{
		"name": "a name",
	}
	assert.Equal(expected, testDataMap)
	var testData2 testOptMarshalVal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONMarshalPresentZeroVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalVal{
		Value: opt.Value(0),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	var testDataMap map[string]any
	require.NoError(json.Unmarshal(y, &testDataMap))
	expected := map[string]any{
		"value": float64(0),
	}
	assert.Equal(expected, testDataMap)
	var testData2 testOptMarshalVal
	err = json.Unmarshal(y, &testData2)
	require.NoError(err)
	assert.True(testData2.Name.IsEmpty())
	require.False(testData2.Value.IsEmpty())
	assert.Equal(0, testData2.Value.Get())
}

func TestJSONMarshalValSimpleNoEmpty(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalVal{
		Name:  opt.Value("a name"),
		Value: opt.Value(123),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":123")
	var testData2 testOptMarshalVal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONUnMarshalOmitVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshalVal
	var testData = `{ "name": "a name" }`
	err := json.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

func TestJSONMarshalNoOmitVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptNoOmitMarshalVal{
		Name:  opt.Value("a name"),
		Value: opt.Empty[int](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":null")
	var testData2 testOptNoOmitMarshalVal
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLMarshalOmitSimpleVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalVal{
		Name:  opt.Value("a name"),
		Value: opt.Empty[int](),
	}
	y, err := yaml.Marshal(&testData)

	require.NoError(err)
	text := string(y)
	assert.Contains(text, "name:")
	assert.NotContains(text, "value:")
	var testData2 testOptMarshalVal
	require.NoError(yaml.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLUnMarshalOmitVal(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshalVal
	var testData = `name: "a name"`
	err := yaml.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

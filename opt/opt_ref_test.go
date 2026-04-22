package opt_test

import (
	"encoding/json"
	"testing"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/opt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

type testMarshalRef struct {
	Name  string              `json:"name" yaml:"name"`
	Value int                 `json:"value" yaml:"value"`
	Opt   opt.Ref[testOptRef] `json:"opt" yaml:"opt,omitzero"`
}

type testOptRef struct {
	Metadata string   `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	ItemList []string `json:"itemList,omitempty" yaml:"itemList,omitempty"`
}

type testOptNoOmitMarshalRef struct {
	Name  opt.Ref[string] `json:"name" yaml:"name"`
	Value opt.Ref[int]    `json:"value" yaml:"value"`
}

type testOptMarshalRef struct {
	Name  opt.Ref[string] `json:"name,omitempty" yaml:"name,omitempty"`
	Value opt.Ref[int]    `json:"value,omitempty" yaml:"value,omitempty"`
}

func TestJSONMarshalOmitRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testMarshalRef{
		Name:  "test1",
		Value: 1,
		Opt:   opt.EmptyRef[testOptRef](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"opt\":null")
	var testData2 testMarshalRef
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONMarshalRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testMarshalRef{
		Name:  "test1",
		Value: 1,
		Opt:   opt.Reference(&testOptRef{"Hello", nil}),
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

func TestJSONMarshalEmptyRef(t *testing.T) {
	require := require.New(t)
	testData := testMarshalRef{
		Name:  "test1",
		Value: 1,
		Opt:   opt.EmptyRef[testOptRef](),
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

func TestJSONMarshalRefSimple(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalRef{
		Name:  opt.Reference(functions.Ref("a name")),
		Value: opt.EmptyRef[int](),
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
	var testData2 testOptMarshalRef
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONMarshalRefSimpleNoEmpty(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalRef{
		Name:  opt.Reference(functions.Ref("a name")),
		Value: opt.Reference(functions.Ref(123)),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":123")
	var testData2 testOptMarshalRef
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestJSONUnMarshalOmitRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshalRef
	var testData = `{ "name": "a name" }`
	err := json.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

func TestJSONMarshalNoOmitRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptNoOmitMarshalRef{
		Name:  opt.Reference(functions.Ref("a name")),
		Value: opt.EmptyRef[int](),
	}
	y, err := json.Marshal(&testData)
	require.NoError(err)
	text := string(y)
	assert.Contains(text, "\"name\":\"a name\"")
	assert.Contains(text, "\"value\":null")
	var testData2 testOptNoOmitMarshalRef
	assert.NoError(json.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLMarshalOmitSimpleRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	testData := testOptMarshalRef{
		Name:  opt.Reference(functions.Ref("a name")),
		Value: opt.EmptyRef[int](),
	}
	y, err := yaml.Marshal(&testData)

	require.NoError(err)
	text := string(y)
	assert.Contains(text, "name:")
	assert.NotContains(text, "value:")
	var testData2 testOptMarshalRef
	require.NoError(yaml.Unmarshal(y, &testData2))
	assert.Equal(testData, testData2)
}

func TestYAMLUnMarshalOmitRef(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)
	var unmarshalled testOptMarshalRef
	var testData = `name: "a name"`
	err := yaml.Unmarshal([]byte(testData), &unmarshalled)
	require.NoError(err)
	assert.True(unmarshalled.Name.HasValue())
	assert.False(unmarshalled.Value.HasValue())
	assert.Equal("a name", unmarshalled.Name.Get())
}

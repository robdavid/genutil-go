package maps

import (
	"errors"
	"testing"

	"github.com/robdavid/genutil-go/errors/test"
	"github.com/robdavid/genutil-go/slices"
	"github.com/stretchr/testify/assert"
)

func TestInsertPathOne(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b"}, 123, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath([]string{"a", "b"}, m)).Must(t))
}

func TestInsertPathTwo(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b"}, 123, m))
	test.Check(t, PutPath([]string{"a", "c"}, 456, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{"b": 123, "c": 456},
	}, m)
}

func TestInsertPathFourDeep(t *testing.T) {
	m := make(map[string]any)
	path := []string{"a", "b", "c", "d"}
	test.Check(t, PutPath(path, 123, m))
	assert.Equal(t, map[string]any{
		"a": map[string]any{
			"b": map[string]any{
				"c": map[string]any{
					"d": 123,
				},
			},
		},
	}, m)
	assert.Equal(t, 123, test.Result(GetPath(path, m)).Must(t))
	res := test.Result(GetPath(slices.Concat(path, []string{"e"}), m))
	assert.True(t, res.IsError())
	assert.True(t, errors.Is(res.GetErr(), ErrKeyError))
}

func TestInsertPathConflictLeaf(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b", "c", "d"}, 123, m))
	err := PutPath([]string{"a", "b", "c"}, 456, m)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

func TestInsertPathConflictInterior(t *testing.T) {
	m := make(map[string]any)
	test.Check(t, PutPath([]string{"a", "b", "c"}, 123, m))
	err := PutPath([]string{"a", "b", "c", "d"}, 456, m)
	assert.EqualError(t, err, "conflict between object and non-object types at key path [a b c]")
}

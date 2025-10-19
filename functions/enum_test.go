package functions_test

import (
	"testing"

	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/iterator"
	"github.com/stretchr/testify/assert"
)

func TestEnumFold(t *testing.T) {
	i := iterator.Of("one", "two", "three")
	interf := func(a, e string, i int) string {
		if i == 0 {
			return a + e
		} else {
			return a + " " + e
		}
	}
	s := functions.EnumFold[iterator.CoreIterator[string]](iterator.Fold[functions.Enum[string], string], i, "", interf)
	assert.Equal(t, "one two three", s)
}

func TestToEnumFold(t *testing.T) {
	i := iterator.Of("one", "two", "three")
	interf := func(a, e string, i int) string {
		if i == 0 {
			return a + e
		} else {
			return a + " " + e
		}
	}
	efold := functions.ToEnumFold(iterator.Fold[functions.Enum[string], string])
	assert.Equal(t, efold(i, "", interf), "one two three")
}

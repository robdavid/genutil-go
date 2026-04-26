package yamlv3_test

import (
	"testing"

	opt "github.com/robdavid/genutil-go/opt/yamlv3"
	"github.com/stretchr/testify/assert"
)

func TestOptionCompat(t *testing.T) {
	assert := assert.New(t)
	var option opt.Option[int]
	val := opt.Value(123)
	option = &val
	ref := option.AsRef()
	assert.True(opt.Equal(&val, ref))
}

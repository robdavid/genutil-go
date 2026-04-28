package yamlv3_test

import (
	"testing"

	"github.com/robdavid/genutil-go/opt"
	yopt "github.com/robdavid/genutil-go/opt/yamlv3"
	"github.com/stretchr/testify/assert"
)

func TestOptionCompat(t *testing.T) {
	assert := assert.New(t)
	var option yopt.Opt[int]
	val := yopt.Value(123)
	ref := yopt.Reference(val.Ref())
	option = ref
	assert.True(opt.Equal(val, option))
	assert.True(opt.Equal(val.OptVal, ref.OptRef))
}

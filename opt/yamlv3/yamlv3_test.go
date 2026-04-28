package yamlv3_test

import (
	"fmt"
	"testing"

	"github.com/robdavid/genutil-go/opt"
	"github.com/robdavid/genutil-go/opt/yamlv3"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestOptionCompat(t *testing.T) {
	assert := assert.New(t)
	var option yamlv3.Opt[int]
	val := yamlv3.Value(123)
	ref := yamlv3.Reference(val.Ref())
	option = ref
	assert.True(opt.Equal(val, option))
	assert.True(opt.Equal(val.OptVal, ref.OptRef))
}

func Example() {
	type MyStruct struct {
		Name    yamlv3.Val[string] `yaml:"name,omitempty"`
		Value   yamlv3.Val[int]    `yaml:"value,omitempty"`
		Version yamlv3.Val[int]    `yaml:"version,omitempty"`
	}
	data := MyStruct{
		Name:    yamlv3.Value("myname"),
		Value:   yamlv3.Value(11),
		Version: yamlv3.Empty[int](),
	}
	text, err := yaml.Marshal(&data)
	if err == nil {
		fmt.Print(string(text))
	}

	data.Version = yamlv3.Value(3)
	text, err = yaml.Marshal(&data)
	if err == nil {
		fmt.Println("---")
		fmt.Print(string(text))
	}

	// Output:
	// name: myname
	// value: 11
	// ---
	// name: myname
	// value: 11
	// version: 3
}

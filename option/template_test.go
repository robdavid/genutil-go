package option

import (
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

type TestTemplateValues struct {
	Str Option[string]
	Num Option[int]
}

func TestTemplateValue(t *testing.T) {
	tmpl, err := template.New("test").Parse("String is {{.Str}} and num is {{.Num}}")
	if !assert.NoError(t, err) {
		return
	}
	var result strings.Builder
	values := TestTemplateValues{}
	if err = tmpl.Execute(&result, &values); !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "String is  and num is ", result.String())
	values.Str = Value("not empty")
	values.Num = Value(25)
	result.Reset()
	if err = tmpl.Execute(&result, &values); !assert.NoError(t, err) {
		return
	}
	assert.Equal(t, "String is not empty and num is 25", result.String())
}

func TestTemplateConditional(t *testing.T) {
	tmpl, err := template.New("test").Funcs(TmplFunctions).Parse("{{if (not (isZero .Str))}}String is {{.Str}}{{else}}Num is {{.Num}}{{end}}")
	if !assert.NoError(t, err) {
		return
	}
	for _, strval := range []Option[string]{Empty[string](), Value(""), Value("here")} {
		var result strings.Builder
		values := TestTemplateValues{Num: Value(123), Str: strval}
		if err = tmpl.Execute(&result, &values); !assert.NoError(t, err) {
			return
		}
		if strval.GetOrZero() == "" {
			assert.Equal(t, "Num is 123", result.String(), "for values %#v", values)
		} else {
			assert.Equal(t, "String is here", result.String())
		}
	}
}

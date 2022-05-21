package option

import (
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

type TestTemplateValues struct {
	Str Option[string]
	Num Option[int]
}

func testTemplate(t *testing.T, tmplText string) (tm *template.Template) {
	var err error
	tm, err = template.New("test").Funcs(TmplFunctions).Parse(tmplText)
	if !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to parse template")
	}
	return
}

func runTemplate(t *testing.T, tm *template.Template, values any) string {
	var result strings.Builder
	if err := tm.Execute(&result, &values); !assert.NoError(t, err) {
		assert.FailNow(t, "Failed to execute template")
	}
	return result.String()
}

func TestTemplateValue(t *testing.T) {
	tmpl := testTemplate(t, "String is {{.Str}} and num is {{.Num}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, values)
	assert.Equal(t, "String is  and num is ", result)
	values.Str = Value("not empty")
	values.Num = Value(25)
	result = runTemplate(t, tmpl, values)
	assert.Equal(t, "String is not empty and num is 25", result)
}

func TestTemplateValueAlternate(t *testing.T) {
	tmpl := testTemplate(t, "Num is {{optOr .Num 50}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, values)
	assert.Equal(t, "Num is 50", result)
	values.Num = Value(25)
	result = runTemplate(t, tmpl, values)
	assert.Equal(t, "Num is 25", result)
}

func TestTemplateValueZero(t *testing.T) {
	tmpl := testTemplate(t, "Num is {{optOrZero .Num}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, values)
	assert.Equal(t, "Num is 0", result)
	values.Num = Value(25)
	result = runTemplate(t, tmpl, values)
	assert.Equal(t, "Num is 25", result)
}

func TestTemplateValueNonOption(t *testing.T) {
	tmpl := testTemplate(t, "Num is {{optOr 123 50}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, values)
	assert.Equal(t, "Num is 123", result)
}

func TestTemplateConditional(t *testing.T) {
	tmpl := testTemplate(t, "{{if (not (isZero .Str))}}String is {{.Str}}{{else}}Num is {{.Num}}{{end}}")
	for _, strval := range []Option[string]{Empty[string](), Value(""), Value("here")} {
		values := TestTemplateValues{Num: Value(123), Str: strval}
		result := runTemplate(t, tmpl, values)
		if strval.GetOrZero() == "" {
			assert.Equal(t, "Num is 123", result, "for values %#v", values)
		} else {
			assert.Equal(t, "String is here", result)
		}
	}
}

func TestTemplateValueFromInt(t *testing.T) {
	tmpl := testTemplate(t, `{{optFmt "Int value is: %02d" (value 4)}}`)
	result := runTemplate(t, tmpl, make(map[string]any))
	assert.Equal(t, "Int value is: 04", result)
}

func TestTemplateValueFromFloat(t *testing.T) {
	tmpl := testTemplate(t, `{{optFmt "Float value is: %.3f" (value 1.23)}}`)
	result := runTemplate(t, tmpl, make(map[string]any))
	assert.Equal(t, "Float value is: 1.230", result)
}

func TestTemplateEmptyFromInt(t *testing.T) {
	tmpl := testTemplate(t, `{{optFmt "Int value is: %02d" (empty 0)}}`)
	result := runTemplate(t, tmpl, make(map[string]any))
	assert.Equal(t, "", result)
}

func TestTemplateIntAsEmpty(t *testing.T) {
	tmpl := testTemplate(t, `{{optFmt "Int value is: %02d" (asEmpty (value 12))}}`)
	result := runTemplate(t, tmpl, make(map[string]any))
	assert.Equal(t, "", result)
}

func TestTemplateConditionalEmpty(t *testing.T) {
	tmpl := testTemplate(t, "{{if hasValue .Str}}String is {{.Str}}{{end}}")
	for _, strval := range []Option[string]{Empty[string](), Value(""), Value("here")} {
		values := TestTemplateValues{Str: strval}
		result := runTemplate(t, tmpl, values)
		if strval.IsEmpty() {
			assert.Equal(t, "", result, "for values %#v", values)
		} else {
			assert.Equal(t, fmt.Sprintf("String is %s", strval.Get()), result)
		}
	}
}

func TestTemplateDefault(t *testing.T) {
	tmpl := testTemplate(t, "Result is {{optDefault 2.14 .Str .Num}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, &values)
	assert.Equal(t, "Result is 2.14", result)
	values.Num = Value(123)
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "Result is 123", result)
	values.Str = Value("my string")
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "Result is my string", result)
	values.Str = Value("")
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "Result is ", result)
}

func TestTemplatePad(t *testing.T) {
	tmpl := testTemplate(t, "Values {{- pad .Str}} {{- pad .Num}}")
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, &values)
	assert.Equal(t, "Values", result)
	values.Str = Value("one")
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "Values one", result)
	values.Num = Value(2)
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "Values one 2", result)
}

func TestTemplatePadstr(t *testing.T) {
	tmpl := testTemplate(t, `{{pad .Str "" "-"}}{{pad .Num "" "-"}}post`)
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, &values)
	assert.Equal(t, "post", result)
	values.Num = Value(123)
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "123-post", result)
	values.Str = Value("pre")
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "pre-123-post", result)
}

func TestTemplateFmt(t *testing.T) {
	tmpl := testTemplate(t, `{{optFmt "%s-" .Str "%02d-" .Num	}}post`)
	values := TestTemplateValues{}
	result := runTemplate(t, tmpl, &values)
	assert.Equal(t, "post", result)
	values.Num = Value(1)
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "01-post", result)
	values.Str = Value("pre")
	result = runTemplate(t, tmpl, &values)
	assert.Equal(t, "pre-01-post", result)
}

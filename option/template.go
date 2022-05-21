package option

import (
	"fmt"
	"reflect"
	"strings"
	"text/template"
)

type objOption interface {
	GetObj() any
	GetObjOK() (any, bool)
	GetObjOr(any) any
	GetObjOrZero() any
	IsObjEmpty() bool
}

func (o Option[_]) GetObj() any {
	return o.Get()
}

func (o OptionRef[_]) GetObj() any {
	return o.Ref()
}

func (o Option[_]) GetObjOK() (any, bool) {
	return o.value, o.nonEmpty
}

func (o OptionRef[_]) GetObjOK() (any, bool) {
	return o.GetOK()
}

func (o Option[T]) GetObjOr(a any) any {
	return o.GetOr(a.(T))
}

func (o OptionRef[T]) GetObjOr(a any) any {
	return o.GetOr(a.(T))
}

func (o Option[_]) GetObjOrZero() any {
	return o.GetOrZero()
}

func (o OptionRef[_]) GetObjOrZero() any {
	return o.GetOrZero()
}

func (o Option[_]) IsObjEmpty() bool {
	return !o.nonEmpty
}

func (o OptionRef[_]) IsObjEmpty() bool {
	return o.ref == nil
}

func tmplOptionReflect(a any) (b any, isopt bool, ok bool) {
	r := reflect.ValueOf(a)
	t := r.Type()
	m, ok := t.MethodByName("GetOK")
	name := t.Name()
	if ok && (strings.HasPrefix(name, "Option[") || strings.HasPrefix(name, "OptionRef[")) {
		valok := m.Func.Call([]reflect.Value{r})
		isopt = true
		ok = valok[1].Bool()
		b = valok[0].Interface()
		return
	} else {
		return
	}
}

func tmplOption(a any) (b any, isopt bool, ok bool) {
	var c objOption
	if c, isopt = a.(objOption); isopt {
		b, ok = c.GetObjOK()
	}
	return
}

// Some useful template functions relating to options

func TmplIsEmpty(a any) bool {
	switch v := a.(type) {
	case objOption:
		return v.IsObjEmpty()
	case []any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	default:
		return false
	}
}

func TmplHasValue(a any) bool {
	return !TmplIsEmpty(a)
}

func TmplIsZero(a any) bool {
	switch v := a.(type) {
	case objOption:
		w, t := v.GetObjOK()
		if t {
			return TmplIsZero(w)
		} else {
			return true
		}
	default:
		ref := reflect.ValueOf(a)
		return ref.IsZero()
	}
}

func TmplGetOrZero(a any) any {
	switch v := a.(type) {
	case objOption:
		return v.GetObjOrZero()
	default:
		return a
	}
}

func TmplGet(a any) any {
	switch v := a.(type) {
	case objOption:
		return v.GetObj()
	default:
		return a
	}
}

func TmplGetOr(a any, b any) any {
	switch v := a.(type) {
	case objOption:
		return v.GetObjOr(b)
	default:
		return a
	}
}

func TmplGetDefault(d any, a ...any) any {
	for i := range a {
		if !TmplIsEmpty(a[i]) {
			return a[i]
		}
	}
	return d
}

func TmplPad(a any, padstr ...string) string {
	text := ""
	pre := " "
	post := ""
	if l := len(padstr); l > 0 {
		pre = padstr[0]
		if l > 1 {
			post = padstr[1]
		}
	}

	switch v := a.(type) {
	case fmt.Stringer:
		text = v.String()
	case string:
		text = v
	default:
		text = fmt.Sprintf("%v", v)
	}

	if text != "" {
		text = pre + text + post
	}
	return text
}

func TmplFmt(pairs ...any) (string, error) {
	var result strings.Builder
	for i := 0; i+1 < len(pairs); i += 2 {
		switch v := pairs[i+1].(type) {
		case objOption:
			if vv, ok := v.GetObjOK(); ok {
				fmt.Fprintf(&result, pairs[i].(string), vv)
			}
		case bool:
			if v {
				result.WriteString(pairs[i].(string))
			}
		default:
			return result.String(), fmt.Errorf("data values should be either option or bool")
		}
	}
	return result.String(), nil
}

func TmplValue[T any](v T) Option[T] {
	return Value(v)
}

func TmplEmpty[T any](v T) Option[T] {
	return Empty[T]()
}

func TmplAsEmpty[T any](o Option[T]) Option[T] {
	if o.IsEmpty() {
		return o
	} else {
		return Empty[T]()
	}
}

var TmplFunctions template.FuncMap = template.FuncMap{
	"isZero":     TmplIsZero,
	"hasValue":   TmplHasValue,
	"isEmpty":    TmplIsEmpty,
	"opt":        TmplGet,
	"optOr":      TmplGetOr,
	"optOrZero":  TmplGetOrZero,
	"optDefault": TmplGetDefault,
	"pad":        TmplPad,
	"optFmt":     TmplFmt,
	"value":      TmplValue[any],
	"empty":      TmplEmpty[any],
	"asEmpty":    TmplAsEmpty[any],
}

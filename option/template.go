package option

import (
	"reflect"
	"strings"
	"text/template"
)

type objOption interface {
	GetObjOK() (any, bool)
}

func (o Option[_]) GetObjOK() (any, bool) {
	return o.value, o.nonEmpty
}

func (o OptionRef[_]) GetObjOK() (any, bool) {
	return o.GetOK()
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
	case IOption[any]:
		_, t := v.GetOK()
		return !t
	case []any:
		return len(v) == 0
	case map[any]any:
		return len(v) == 0
	default:
		return false
	}
}

func TmplIsZero(a any) bool {
	b, isopt, ok := tmplOption(a)
	if isopt {
		if ok {
			return TmplIsZero(b)
		} else {
			return true
		}
	} else {
		ref := reflect.ValueOf(a)
		return ref.IsZero()
	}
	// switch v := a.(type) {
	// case IOption[_]:
	// 	w, t := v.getObjOK()
	// 	if t {
	// 		return TmplIsZero(w)
	// 	} else {
	// 		return true
	// 	}
	// default:
	// 	ref := reflect.ValueOf(a)
	// 	return ref.IsZero()
	// }
}

func TmplGetOrZero(a any) any {
	switch v := a.(type) {
	case IOption[any]:
		return v.GetOrZero()
	default:
		return a
	}
}

func TmplGet(a any) any {
	switch v := a.(type) {
	case IOption[any]:
		return v.Get()
	default:
		return a
	}
}

func TmplGetOr(a any, b any) any {
	switch v := a.(type) {
	case IOption[any]:
		return v.GetOr(b)
	default:
		return a
	}
}

func TmplGetDefault(d any, a any) any {
	return TmplGetOr(a, d)
}

var TmplFunctions template.FuncMap = template.FuncMap{
	"isZero":     TmplIsZero,
	"isEmpty":    TmplIsEmpty,
	"get":        TmplGet,
	"getOr":      TmplGetOr,
	"getOrZero":  TmplGetOrZero,
	"getDefault": TmplGetDefault,
}

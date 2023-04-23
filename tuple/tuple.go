// Genric tuples, of size 0 to 9.
package tuple

import (
	"errors"
	"fmt"
	"strings"
)

var ErrIndex = errors.New("index error")

/*
General tuple interface, implemented by all *TupleN
types
*/
type Tuple interface {
	// Get the nth element of the tuple
	Get(int) any
	// Return the number of elements in the tuple
	Size() int
	// Return the tuple as a string, formatted (e1,e2,...)
	String() string
	// Tuple of first size-1 elements
	Pre() Tuple
	// Return the last element in the tuple
	Last() any
}

func tupleString(tuple Tuple) string {
	var result strings.Builder
	result.WriteString("(")
	for i := 0; i < tuple.Size(); i++ {
		if i > 0 {
			result.WriteString(",")
		}
		fmt.Fprintf(&result, "%#v", tuple.Get(i))
	}
	result.WriteString(")")
	return result.String()
}

func tupleGet(tuple Tuple, n int) any {
	if n == tuple.Size()-1 {
		return tuple.Last()
	} else {
		return tupleGet(tuple.Pre(), n)
	}
}

func Slice(t Tuple) (result []any) {
	result = make([]any, t.Size())
	for i := 0; i < t.Size(); i++ {
		result[i] = t.Get(i)
	}
	return
}

type Tuple0 struct{}

func Of0() Tuple0 {
	return Tuple0{}
}

func Unit() Tuple0 {
	return Tuple0{}
}

func (*Tuple0) Size() int         { return 0 }
func (*Tuple0) Get(int) any       { panic(ErrIndex) }
func (*Tuple0) Last() any         { panic(ErrIndex) }
func (*Tuple0) Pre() Tuple        { panic(ErrIndex) }
func (t0 *Tuple0) String() string { return tupleString(t0) }

func Pair[A any, B any](a A, b B) Tuple2[A, B] {
	return Of2(a, b)
}

//go:generate code-template tuple.tmpl

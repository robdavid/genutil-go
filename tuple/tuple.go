package tuple

import (
	"errors"
	"fmt"
	"strings"
)

var ErrIndex = errors.New("Index error")

type Tuple interface {
	Get(int) any
	Size() int
	String() string
	Pre() Tuple // Tuple of first size-1 elements
	Last() any
}

func tupleString(tuple Tuple) string {
	var result strings.Builder
	result.WriteString("(")
	for i := 0; i < tuple.Size(); i++ {
		if i > 0 {
			result.WriteString(",")
		}
		fmt.Fprintf(&result, "%v", tuple.Get(i))
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

type Tuple1[A any] struct {
	Tuple0
	First A
}

func Of1[A any](a A) Tuple1[A] {
	return Tuple1[A]{Unit(), a}
}

func (*Tuple1[A]) Size() int         { return 1 }
func (t1 *Tuple1[A]) Pre() Tuple     { return &t1.Tuple0 }
func (t1 *Tuple1[A]) Last() any      { return t1.First }
func (t1 *Tuple1[A]) Get(n int) any  { return tupleGet(t1, n) }
func (t1 *Tuple1[A]) String() string { return tupleString(t1) }

type Tuple2[A any, B any] struct {
	Tuple1[A]
	Second B
}

func Of2[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{Of1(a), b}
}

func Pair[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{Of1(a), b}
}

func (*Tuple2[A, B]) Size() int         { return 2 }
func (t2 *Tuple2[A, B]) Pre() Tuple     { return &t2.Tuple1 }
func (t2 *Tuple2[A, B]) Last() any      { return t2.Second }
func (t2 *Tuple2[A, B]) String() string { return tupleString(t2) }
func (t2 *Tuple2[A, B]) Get(n int) any  { return tupleGet(t2, n) }

type Tuple3[A any, B any, C any] struct {
	Tuple2[A, B]
	Third C
}

func Of3[A any, B any, C any](a A, b B, c C) Tuple3[A, B, C] {
	return Tuple3[A, B, C]{Of2(a, b), c}
}

func (*Tuple3[A, B, C]) Size() int         { return 3 }
func (t3 *Tuple3[A, B, C]) Pre() Tuple     { return &t3.Tuple2 }
func (t3 *Tuple3[A, B, C]) Last() any      { return t3.Third }
func (t3 *Tuple3[A, B, C]) String() string { return tupleString(t3) }
func (t3 *Tuple3[A, B, C]) Get(n int) any  { return tupleGet(t3, n) }

type Tuple4[A any, B any, C any, D any] struct {
	Tuple3[A, B, C]
	Forth D
}

func Of4[A any, B any, C any, D any](a A, b B, c C, d D) Tuple4[A, B, C, D] {
	return Tuple4[A, B, C, D]{Of3(a, b, c), d}
}

func (*Tuple4[A, B, C, D]) Size() int         { return 4 }
func (t4 *Tuple4[A, B, C, D]) Pre() Tuple     { return &t4.Tuple3 }
func (t4 *Tuple4[A, B, C, D]) Last() any      { return t4.Forth }
func (t4 *Tuple4[A, B, C, D]) String() string { return tupleString(t4) }
func (t4 *Tuple4[A, B, C, D]) Get(n int) any  { return tupleGet(t4, n) }

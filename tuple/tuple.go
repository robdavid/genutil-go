package tuple

import "fmt"

type Tuple2[A any, B any] struct {
	First  A
	Second B
}

func NewTuple2[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{a, b}
}

func Pair[A any, B any](a A, b B) Tuple2[A, B] {
	return Tuple2[A, B]{a, b}
}

func (t2 *Tuple2[A, B]) String() string {
	return fmt.Sprintf("(%v,%v)", t2.First, t2.Second)
}

func (t2 *Tuple2[A, B]) Slice() []any {
	return []any{t2.First, t2.Second}
}

func (t2 *Tuple2[A, B]) RefSlice() []any {
	return []any{&t2.First, &t2.Second}
}

type Tuple3[A any, B any, C any] struct {
	Tuple2[A, B]
	Third C
}

func (t3 *Tuple3[A, B, C]) String() string {
	return fmt.Sprintf("(%v,%v,%v)", t3.First, t3.Second, t3.Third)
}

func (t3 *Tuple3[A, B, C]) Slice() []any {
	return []any{t3.First, t3.Second, t3.Third}
}

func (t3 *Tuple3[A, B, C]) RefSlice() []any {
	return []any{&t3.First, &t3.Second, &t3.Third}
}

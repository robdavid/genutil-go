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

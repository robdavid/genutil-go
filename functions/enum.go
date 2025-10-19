package functions

func EnumOf[T any](index int, value T) Enum[T] { return Enum[T]{Index: index, Value: value} }

func EnumFoldFunc[T any, U any](f func(U, T, int) U) func(Enum[U], T) Enum[U] {
	return func(acc Enum[U], e T) Enum[U] {
		u := f(acc.Value, e, acc.Index)
		return EnumOf(acc.Index+1, u)
	}
}

func EnumFold[C any, T any, U any, F func(C, Enum[U], func(Enum[U], T) Enum[U]) Enum[U]](fld F, col C, init U, f func(U, T, int) U) U {
	fu := fld(col, EnumOf(0, init), EnumFoldFunc(f))
	return fu.Value
}

func ToEnumFold[C any, T any, U any, F func(C, Enum[U], func(Enum[U], T) Enum[U]) Enum[U]](fld F) func(C, U, func(U, T, int) U) U {
	return func(source C, init U, f func(U, T, int) U) U {
		return fld(source, EnumOf(0, init), EnumFoldFunc(f)).Value
	}
}

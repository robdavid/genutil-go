package handler

// Variant of try with 2 non-error arguments
func Try2[T1 any, T2 any](p1 T1, p2 T2, err error) (T1, T2) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2
}

// Variant of try with 3 non-error arguments
func Try3[T1 any, T2 any, T3 any](p1 T1, p2 T2, p3 T3, err error) (T1, T2, T3) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3
}

// Variant of try with 4 non-error arguments
func Try4[T1 any, T2 any, T3 any, T4 any](p1 T1, p2 T2, p3 T3, p4 T4, err error) (T1, T2, T3, T4) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4
}

// Variant of try with 5 non-error arguments
func Try5[T1 any, T2 any, T3 any, T4 any, T5 any](p1 T1, p2 T2, p3 T3, p4 T4, p5 T5, err error) (T1, T2, T3, T4, T5) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4, p5
}

// Variant of try with 6 non-error arguments
func Try6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](p1 T1, p2 T2, p3 T3, p4 T4, p5 T5, p6 T6, err error) (T1, T2, T3, T4, T5, T6) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4, p5, p6
}

// Variant of try with 7 non-error arguments
func Try7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](p1 T1, p2 T2, p3 T3, p4 T4, p5 T5, p6 T6, p7 T7, err error) (T1, T2, T3, T4, T5, T6, T7) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4, p5, p6, p7
}

// Variant of try with 8 non-error arguments
func Try8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](p1 T1, p2 T2, p3 T3, p4 T4, p5 T5, p6 T6, p7 T7, p8 T8, err error) (T1, T2, T3, T4, T5, T6, T7, T8) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4, p5, p6, p7, p8
}

// Variant of try with 9 non-error arguments
func Try9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](p1 T1, p2 T2, p3 T3, p4 T4, p5 T5, p6 T6, p7 T7, p8 T8, p9 T9, err error) (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	if err != nil {
		panic(TryError{err})
	}
	return p1, p2, p3, p4, p5, p6, p7, p8, p9
}


package tuple

// Tuple of 1 fields
type Tuple1[T1 any] struct {
	Tuple0
	First T1
}

func Of1[T1 any](t1 T1) Tuple1[T1] {
	return Tuple1[T1]{Of0(), t1}
}

func (*Tuple1[T1]) Size() int         { return 1 }
func (t1 *Tuple1[T1]) Pre() Tuple     { return &t1.Tuple0 }
func (t1 *Tuple1[T1]) Last() any      { return t1.First }
func (t1 *Tuple1[T1]) Get(n int) any  { return tupleGet(t1, n) }
func (t1 *Tuple1[T1]) String() string { return tupleString(t1) }

// Tuple of 2 fields
type Tuple2[T1 any, T2 any] struct {
	Tuple1[T1]
	Second T2
}

func Of2[T1 any, T2 any](t1 T1, t2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2]{Of1(t1), t2}
}

func (*Tuple2[T1, T2]) Size() int         { return 2 }
func (t2 *Tuple2[T1, T2]) Pre() Tuple     { return &t2.Tuple1 }
func (t2 *Tuple2[T1, T2]) Last() any      { return t2.Second }
func (t2 *Tuple2[T1, T2]) Get(n int) any  { return tupleGet(t2, n) }
func (t2 *Tuple2[T1, T2]) String() string { return tupleString(t2) }

// Tuple of 3 fields
type Tuple3[T1 any, T2 any, T3 any] struct {
	Tuple2[T1, T2]
	Third T3
}

func Of3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3]{Of2(t1, t2), t3}
}

func (*Tuple3[T1, T2, T3]) Size() int         { return 3 }
func (t3 *Tuple3[T1, T2, T3]) Pre() Tuple     { return &t3.Tuple2 }
func (t3 *Tuple3[T1, T2, T3]) Last() any      { return t3.Third }
func (t3 *Tuple3[T1, T2, T3]) Get(n int) any  { return tupleGet(t3, n) }
func (t3 *Tuple3[T1, T2, T3]) String() string { return tupleString(t3) }

// Tuple of 4 fields
type Tuple4[T1 any, T2 any, T3 any, T4 any] struct {
	Tuple3[T1, T2, T3]
	Forth T4
}

func Of4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4) Tuple4[T1, T2, T3, T4] {
	return Tuple4[T1, T2, T3, T4]{Of3(t1, t2, t3), t4}
}

func (*Tuple4[T1, T2, T3, T4]) Size() int         { return 4 }
func (t4 *Tuple4[T1, T2, T3, T4]) Pre() Tuple     { return &t4.Tuple3 }
func (t4 *Tuple4[T1, T2, T3, T4]) Last() any      { return t4.Forth }
func (t4 *Tuple4[T1, T2, T3, T4]) Get(n int) any  { return tupleGet(t4, n) }
func (t4 *Tuple4[T1, T2, T3, T4]) String() string { return tupleString(t4) }

// Tuple of 5 fields
type Tuple5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	Tuple4[T1, T2, T3, T4]
	Fifth T5
}

func Of5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Tuple5[T1, T2, T3, T4, T5] {
	return Tuple5[T1, T2, T3, T4, T5]{Of4(t1, t2, t3, t4), t5}
}

func (*Tuple5[T1, T2, T3, T4, T5]) Size() int         { return 5 }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Pre() Tuple     { return &t5.Tuple4 }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Last() any      { return t5.Fifth }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Get(n int) any  { return tupleGet(t5, n) }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) String() string { return tupleString(t5) }

// Tuple of 6 fields
type Tuple6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct {
	Tuple5[T1, T2, T3, T4, T5]
	Sixth T6
}

func Of6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) Tuple6[T1, T2, T3, T4, T5, T6] {
	return Tuple6[T1, T2, T3, T4, T5, T6]{Of5(t1, t2, t3, t4, t5), t6}
}

func (*Tuple6[T1, T2, T3, T4, T5, T6]) Size() int         { return 6 }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Pre() Tuple     { return &t6.Tuple5 }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Last() any      { return t6.Sixth }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Get(n int) any  { return tupleGet(t6, n) }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) String() string { return tupleString(t6) }

// Tuple of 7 fields
type Tuple7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct {
	Tuple6[T1, T2, T3, T4, T5, T6]
	Seventh T7
}

func Of7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) Tuple7[T1, T2, T3, T4, T5, T6, T7] {
	return Tuple7[T1, T2, T3, T4, T5, T6, T7]{Of6(t1, t2, t3, t4, t5, t6), t7}
}

func (*Tuple7[T1, T2, T3, T4, T5, T6, T7]) Size() int         { return 7 }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Pre() Tuple     { return &t7.Tuple6 }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Last() any      { return t7.Seventh }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Get(n int) any  { return tupleGet(t7, n) }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) String() string { return tupleString(t7) }

// Tuple of 8 fields
type Tuple8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct {
	Tuple7[T1, T2, T3, T4, T5, T6, T7]
	Eighth T8
}

func Of8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) Tuple8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]{Of7(t1, t2, t3, t4, t5, t6, t7), t8}
}

func (*Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Size() int         { return 8 }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Pre() Tuple     { return &t8.Tuple7 }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Last() any      { return t8.Eighth }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Get(n int) any  { return tupleGet(t8, n) }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) String() string { return tupleString(t8) }

// Tuple of 9 fields
type Tuple9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct {
	Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]
	Ninth T9
}

func Of9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]{Of8(t1, t2, t3, t4, t5, t6, t7, t8), t9}
}

func (*Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Size() int         { return 9 }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Pre() Tuple     { return &t9.Tuple8 }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Last() any      { return t9.Ninth }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Get(n int) any  { return tupleGet(t9, n) }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) String() string { return tupleString(t9) }


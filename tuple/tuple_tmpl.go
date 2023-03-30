package tuple

// Tuple of 1 fields
type Tuple1[T1 any] struct {
	First T1
}

func Of1[T1 any](t1 T1) Tuple1[T1] {
	return Tuple1[T1] { t1 }
}

// Interface implementation
func (t1 *Tuple1[T1]) Tuple0() Tuple0  { return Tuple0{  } }
func (*Tuple1[T1]) Size() int         { return 1 }
func (t1 *Tuple1[T1]) Pre() Tuple     { return &Tuple0{  } }
func (t1 *Tuple1[T1]) Last() any      { return t1.First }
func (t1 *Tuple1[T1]) Get(n int) any  { return tupleGet(t1, n) }
func (t1 *Tuple1[T1]) String() string { return tupleString(t1) }
func (t1 Tuple1[T1]) ToRef() *Tuple1[T1] { return &t1 }

// Returns the values in the tuple as a sequence of 1 values
func (t1 *Tuple1[T1]) Return() (T1) {
	return t1.First
}
 
// Tuple of 2 fields
type Tuple2[T1 any, T2 any] struct {
	First T1
	Second T2
}

func Of2[T1 any, T2 any](t1 T1, t2 T2) Tuple2[T1, T2] {
	return Tuple2[T1, T2] { t1, t2 }
}

// Interface implementation
func (t2 *Tuple2[T1, T2]) Tuple1() Tuple1[T1]  { return Tuple1[T1]{ t2.First } }
func (*Tuple2[T1, T2]) Size() int         { return 2 }
func (t2 *Tuple2[T1, T2]) Pre() Tuple     { return &Tuple1[T1]{ t2.First } }
func (t2 *Tuple2[T1, T2]) Last() any      { return t2.Second }
func (t2 *Tuple2[T1, T2]) Get(n int) any  { return tupleGet(t2, n) }
func (t2 *Tuple2[T1, T2]) String() string { return tupleString(t2) }
func (t2 Tuple2[T1, T2]) ToRef() *Tuple2[T1, T2] { return &t2 }

// Returns the values in the tuple as a sequence of 2 values
func (t2 *Tuple2[T1, T2]) Return() (T1, T2) {
	return t2.First, t2.Second
}
 
// Tuple of 3 fields
type Tuple3[T1 any, T2 any, T3 any] struct {
	First T1
	Second T2
	Third T3
}

func Of3[T1 any, T2 any, T3 any](t1 T1, t2 T2, t3 T3) Tuple3[T1, T2, T3] {
	return Tuple3[T1, T2, T3] { t1, t2, t3 }
}

// Interface implementation
func (t3 *Tuple3[T1, T2, T3]) Tuple2() Tuple2[T1, T2]  { return Tuple2[T1, T2]{ t3.First, t3.Second } }
func (*Tuple3[T1, T2, T3]) Size() int         { return 3 }
func (t3 *Tuple3[T1, T2, T3]) Pre() Tuple     { return &Tuple2[T1, T2]{ t3.First, t3.Second } }
func (t3 *Tuple3[T1, T2, T3]) Last() any      { return t3.Third }
func (t3 *Tuple3[T1, T2, T3]) Get(n int) any  { return tupleGet(t3, n) }
func (t3 *Tuple3[T1, T2, T3]) String() string { return tupleString(t3) }
func (t3 Tuple3[T1, T2, T3]) ToRef() *Tuple3[T1, T2, T3] { return &t3 }

// Returns the values in the tuple as a sequence of 3 values
func (t3 *Tuple3[T1, T2, T3]) Return() (T1, T2, T3) {
	return t3.First, t3.Second, t3.Third
}
 
// Tuple of 4 fields
type Tuple4[T1 any, T2 any, T3 any, T4 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
}

func Of4[T1 any, T2 any, T3 any, T4 any](t1 T1, t2 T2, t3 T3, t4 T4) Tuple4[T1, T2, T3, T4] {
	return Tuple4[T1, T2, T3, T4] { t1, t2, t3, t4 }
}

// Interface implementation
func (t4 *Tuple4[T1, T2, T3, T4]) Tuple3() Tuple3[T1, T2, T3]  { return Tuple3[T1, T2, T3]{ t4.First, t4.Second, t4.Third } }
func (*Tuple4[T1, T2, T3, T4]) Size() int         { return 4 }
func (t4 *Tuple4[T1, T2, T3, T4]) Pre() Tuple     { return &Tuple3[T1, T2, T3]{ t4.First, t4.Second, t4.Third } }
func (t4 *Tuple4[T1, T2, T3, T4]) Last() any      { return t4.Forth }
func (t4 *Tuple4[T1, T2, T3, T4]) Get(n int) any  { return tupleGet(t4, n) }
func (t4 *Tuple4[T1, T2, T3, T4]) String() string { return tupleString(t4) }
func (t4 Tuple4[T1, T2, T3, T4]) ToRef() *Tuple4[T1, T2, T3, T4] { return &t4 }

// Returns the values in the tuple as a sequence of 4 values
func (t4 *Tuple4[T1, T2, T3, T4]) Return() (T1, T2, T3, T4) {
	return t4.First, t4.Second, t4.Third, t4.Forth
}
 
// Tuple of 5 fields
type Tuple5[T1 any, T2 any, T3 any, T4 any, T5 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
	Fifth T5
}

func Of5[T1 any, T2 any, T3 any, T4 any, T5 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5) Tuple5[T1, T2, T3, T4, T5] {
	return Tuple5[T1, T2, T3, T4, T5] { t1, t2, t3, t4, t5 }
}

// Interface implementation
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Tuple4() Tuple4[T1, T2, T3, T4]  { return Tuple4[T1, T2, T3, T4]{ t5.First, t5.Second, t5.Third, t5.Forth } }
func (*Tuple5[T1, T2, T3, T4, T5]) Size() int         { return 5 }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Pre() Tuple     { return &Tuple4[T1, T2, T3, T4]{ t5.First, t5.Second, t5.Third, t5.Forth } }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Last() any      { return t5.Fifth }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Get(n int) any  { return tupleGet(t5, n) }
func (t5 *Tuple5[T1, T2, T3, T4, T5]) String() string { return tupleString(t5) }
func (t5 Tuple5[T1, T2, T3, T4, T5]) ToRef() *Tuple5[T1, T2, T3, T4, T5] { return &t5 }

// Returns the values in the tuple as a sequence of 5 values
func (t5 *Tuple5[T1, T2, T3, T4, T5]) Return() (T1, T2, T3, T4, T5) {
	return t5.First, t5.Second, t5.Third, t5.Forth, t5.Fifth
}
 
// Tuple of 6 fields
type Tuple6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
	Fifth T5
	Sixth T6
}

func Of6[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6) Tuple6[T1, T2, T3, T4, T5, T6] {
	return Tuple6[T1, T2, T3, T4, T5, T6] { t1, t2, t3, t4, t5, t6 }
}

// Interface implementation
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Tuple5() Tuple5[T1, T2, T3, T4, T5]  { return Tuple5[T1, T2, T3, T4, T5]{ t6.First, t6.Second, t6.Third, t6.Forth, t6.Fifth } }
func (*Tuple6[T1, T2, T3, T4, T5, T6]) Size() int         { return 6 }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Pre() Tuple     { return &Tuple5[T1, T2, T3, T4, T5]{ t6.First, t6.Second, t6.Third, t6.Forth, t6.Fifth } }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Last() any      { return t6.Sixth }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Get(n int) any  { return tupleGet(t6, n) }
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) String() string { return tupleString(t6) }
func (t6 Tuple6[T1, T2, T3, T4, T5, T6]) ToRef() *Tuple6[T1, T2, T3, T4, T5, T6] { return &t6 }

// Returns the values in the tuple as a sequence of 6 values
func (t6 *Tuple6[T1, T2, T3, T4, T5, T6]) Return() (T1, T2, T3, T4, T5, T6) {
	return t6.First, t6.Second, t6.Third, t6.Forth, t6.Fifth, t6.Sixth
}
 
// Tuple of 7 fields
type Tuple7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
	Fifth T5
	Sixth T6
	Seventh T7
}

func Of7[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7) Tuple7[T1, T2, T3, T4, T5, T6, T7] {
	return Tuple7[T1, T2, T3, T4, T5, T6, T7] { t1, t2, t3, t4, t5, t6, t7 }
}

// Interface implementation
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Tuple6() Tuple6[T1, T2, T3, T4, T5, T6]  { return Tuple6[T1, T2, T3, T4, T5, T6]{ t7.First, t7.Second, t7.Third, t7.Forth, t7.Fifth, t7.Sixth } }
func (*Tuple7[T1, T2, T3, T4, T5, T6, T7]) Size() int         { return 7 }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Pre() Tuple     { return &Tuple6[T1, T2, T3, T4, T5, T6]{ t7.First, t7.Second, t7.Third, t7.Forth, t7.Fifth, t7.Sixth } }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Last() any      { return t7.Seventh }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Get(n int) any  { return tupleGet(t7, n) }
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) String() string { return tupleString(t7) }
func (t7 Tuple7[T1, T2, T3, T4, T5, T6, T7]) ToRef() *Tuple7[T1, T2, T3, T4, T5, T6, T7] { return &t7 }

// Returns the values in the tuple as a sequence of 7 values
func (t7 *Tuple7[T1, T2, T3, T4, T5, T6, T7]) Return() (T1, T2, T3, T4, T5, T6, T7) {
	return t7.First, t7.Second, t7.Third, t7.Forth, t7.Fifth, t7.Sixth, t7.Seventh
}
 
// Tuple of 8 fields
type Tuple8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
	Fifth T5
	Sixth T6
	Seventh T7
	Eighth T8
}

func Of8[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8) Tuple8[T1, T2, T3, T4, T5, T6, T7, T8] {
	return Tuple8[T1, T2, T3, T4, T5, T6, T7, T8] { t1, t2, t3, t4, t5, t6, t7, t8 }
}

// Interface implementation
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Tuple7() Tuple7[T1, T2, T3, T4, T5, T6, T7]  { return Tuple7[T1, T2, T3, T4, T5, T6, T7]{ t8.First, t8.Second, t8.Third, t8.Forth, t8.Fifth, t8.Sixth, t8.Seventh } }
func (*Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Size() int         { return 8 }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Pre() Tuple     { return &Tuple7[T1, T2, T3, T4, T5, T6, T7]{ t8.First, t8.Second, t8.Third, t8.Forth, t8.Fifth, t8.Sixth, t8.Seventh } }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Last() any      { return t8.Eighth }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Get(n int) any  { return tupleGet(t8, n) }
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) String() string { return tupleString(t8) }
func (t8 Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) ToRef() *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8] { return &t8 }

// Returns the values in the tuple as a sequence of 8 values
func (t8 *Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]) Return() (T1, T2, T3, T4, T5, T6, T7, T8) {
	return t8.First, t8.Second, t8.Third, t8.Forth, t8.Fifth, t8.Sixth, t8.Seventh, t8.Eighth
}
 
// Tuple of 9 fields
type Tuple9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any] struct {
	First T1
	Second T2
	Third T3
	Forth T4
	Fifth T5
	Sixth T6
	Seventh T7
	Eighth T8
	Ninth T9
}

func Of9[T1 any, T2 any, T3 any, T4 any, T5 any, T6 any, T7 any, T8 any, T9 any](t1 T1, t2 T2, t3 T3, t4 T4, t5 T5, t6 T6, t7 T7, t8 T8, t9 T9) Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9] {
	return Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9] { t1, t2, t3, t4, t5, t6, t7, t8, t9 }
}

// Interface implementation
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Tuple8() Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]  { return Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]{ t9.First, t9.Second, t9.Third, t9.Forth, t9.Fifth, t9.Sixth, t9.Seventh, t9.Eighth } }
func (*Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Size() int         { return 9 }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Pre() Tuple     { return &Tuple8[T1, T2, T3, T4, T5, T6, T7, T8]{ t9.First, t9.Second, t9.Third, t9.Forth, t9.Fifth, t9.Sixth, t9.Seventh, t9.Eighth } }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Last() any      { return t9.Ninth }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Get(n int) any  { return tupleGet(t9, n) }
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) String() string { return tupleString(t9) }
func (t9 Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) ToRef() *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9] { return &t9 }

// Returns the values in the tuple as a sequence of 9 values
func (t9 *Tuple9[T1, T2, T3, T4, T5, T6, T7, T8, T9]) Return() (T1, T2, T3, T4, T5, T6, T7, T8, T9) {
	return t9.First, t9.Second, t9.Third, t9.Forth, t9.Fifth, t9.Sixth, t9.Seventh, t9.Eighth, t9.Ninth
}
 

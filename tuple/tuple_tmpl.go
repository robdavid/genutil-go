package tuple

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


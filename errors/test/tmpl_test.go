package test
import (
  "testing"
  "github.com/robdavid/genutil-go/errors/result"
  "github.com/stretchr/testify/assert"
)




func dummySuccessFunction9() (string, int, float64, bool, rune, byte, int32, float32, complex64, error) {
	return "success", 123, 4.56, true, 'x', 'y', 789, 1.23, 1 + 3i, nil
}


func dummySuccessFunction2() (string, int, error) {
	r := result.From3(dummySuccessFunction3())
	return result.Value2(r.Get().Tuple2).ToRef().Return()
}

func TestSuccess2(t *testing.T) {
  var err error
  a1, a2 := Result2(dummySuccessFunction2()).Must2(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
}


func dummySuccessFunction3() (string, int, float64, error) {
	r := result.From4(dummySuccessFunction4())
	return result.Value3(r.Get().Tuple3).ToRef().Return()
}

func TestSuccess3(t *testing.T) {
  var err error
  a1, a2, a3 := Result3(dummySuccessFunction3()).Must3(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
}


func dummySuccessFunction4() (string, int, float64, bool, error) {
	r := result.From5(dummySuccessFunction5())
	return result.Value4(r.Get().Tuple4).ToRef().Return()
}

func TestSuccess4(t *testing.T) {
  var err error
  a1, a2, a3, a4 := Result4(dummySuccessFunction4()).Must4(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
}


func dummySuccessFunction5() (string, int, float64, bool, rune, error) {
	r := result.From6(dummySuccessFunction6())
	return result.Value5(r.Get().Tuple5).ToRef().Return()
}

func TestSuccess5(t *testing.T) {
  var err error
  a1, a2, a3, a4, a5 := Result5(dummySuccessFunction5()).Must5(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
  assert.Equal(t, rune('x'), a5)
}


func dummySuccessFunction6() (string, int, float64, bool, rune, byte, error) {
	r := result.From7(dummySuccessFunction7())
	return result.Value6(r.Get().Tuple6).ToRef().Return()
}

func TestSuccess6(t *testing.T) {
  var err error
  a1, a2, a3, a4, a5, a6 := Result6(dummySuccessFunction6()).Must6(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
  assert.Equal(t, rune('x'), a5)
  assert.Equal(t, byte('y'), a6)
}


func dummySuccessFunction7() (string, int, float64, bool, rune, byte, int32, error) {
	r := result.From8(dummySuccessFunction8())
	return result.Value7(r.Get().Tuple7).ToRef().Return()
}

func TestSuccess7(t *testing.T) {
  var err error
  a1, a2, a3, a4, a5, a6, a7 := Result7(dummySuccessFunction7()).Must7(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
  assert.Equal(t, rune('x'), a5)
  assert.Equal(t, byte('y'), a6)
  assert.Equal(t, int32(789), a7)
}


func dummySuccessFunction8() (string, int, float64, bool, rune, byte, int32, float32, error) {
	r := result.From9(dummySuccessFunction9())
	return result.Value8(r.Get().Tuple8).ToRef().Return()
}

func TestSuccess8(t *testing.T) {
  var err error
  a1, a2, a3, a4, a5, a6, a7, a8 := Result8(dummySuccessFunction8()).Must8(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
  assert.Equal(t, rune('x'), a5)
  assert.Equal(t, byte('y'), a6)
  assert.Equal(t, int32(789), a7)
  assert.Equal(t, float32(1.23), a8)
}

func TestSuccess9(t *testing.T) {
  var err error
  a1, a2, a3, a4, a5, a6, a7, a8, a9 := Result9(dummySuccessFunction9()).Must9(t)
  assert.Nil(t, err)
  assert.Equal(t, string("success"), a1)
  assert.Equal(t, int(123), a2)
  assert.Equal(t, float64(4.56), a3)
  assert.Equal(t, bool(true), a4)
  assert.Equal(t, rune('x'), a5)
  assert.Equal(t, byte('y'), a6)
  assert.Equal(t, int32(789), a7)
  assert.Equal(t, float32(1.23), a8)
  assert.Equal(t, complex64(1 + 3i), a9)
}


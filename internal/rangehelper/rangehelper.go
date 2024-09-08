package rangehelper

import (
	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/realnum"
)

// RangeSize calculates the number of elements expected in a given range. Signs of step and the
// difference between start and end are ignored. Will panic (divide zero) if step is zero.
func RangeSize[T realnum.Real, S realnum.Real](start, end T, step S, inclusive bool) (int, T) {
	if realnum.IsInteger(start) {
		iRange := realnum.Abs(realnum.Sub[int64](end, start))
		iStep := int64(realnum.Abs(step))
		intervals := int(iRange / iStep)
		return intervals + functions.IfElse(iStep*int64(intervals) < iRange || inclusive, 1, 0), T(iStep)
	} else {
		fRange := realnum.Abs(realnum.Sub[float64](end, start))
		fStep := float64(realnum.Abs(step))
		intervals := int(fRange / fStep)
		return intervals + functions.IfElse(fStep*float64(intervals) < fRange || inclusive, 1, 0), T(fStep)
	}
}

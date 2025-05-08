package rangehelper

import (
	"github.com/robdavid/genutil-go/functions"
	"github.com/robdavid/genutil-go/ordered"
)

// RangeSize calculates the number of elements expected in a given range. Signs of step and the
// difference between start and end are ignored. Will panic (divide zero) if step is zero.
// Returns the number of elements and the absolute value of the step, cast to T.
func RangeSize[T ordered.Real, S ordered.Real](start, end T, step S, inclusive bool) (int, T) {
	if ordered.IsInteger(start) {
		iRange := ordered.Abs(ordered.Sub[int64](end, start))
		iStep := int64(ordered.Abs(step))
		intervals := int(iRange / iStep)
		return intervals + functions.IfElse(iStep*int64(intervals) < iRange || inclusive, 1, 0), T(iStep)
	} else {
		fRange := ordered.Abs(ordered.Sub[float64](end, start))
		fStep := float64(ordered.Abs(step))
		intervals := int(fRange / fStep)
		return intervals + functions.IfElse(fStep*float64(intervals) < fRange || inclusive, 1, 0), T(fStep)
	}
}

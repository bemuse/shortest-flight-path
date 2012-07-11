package interpolate

import (
	"math"
)

func Interpolator(inLow, inHigh, outDesired, precision float64, f func(float64) float64) (inDesired float64, ok bool) {
	outLow := f(inLow)
	outHigh := f(inHigh)
	lastDelta := math.Abs(outHigh - outLow)

	for {
		factor := (outDesired - outLow) / (outHigh - outLow)
		inTest := inLow + factor*(inHigh-inLow)
		outTest := f(inTest)
		delta := math.Abs(outTest - outDesired)
		if delta < precision {
			return inTest, true
		} else if delta > lastDelta {
			break
		}
		lastDelta = delta

		if (outDesired-outTest)*(outHigh-outLow) > 0 {
			inLow = inTest
			outLow = outTest
		} else {
			inHigh = inTest
			outHigh = outTest
		}
	}

	return 0, false
}

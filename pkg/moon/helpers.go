package moon

import "math"

func dsin(x float64) float64 {
	return math.Sin(dtr(x))
}

func dcos(x float64) float64 {
	return math.Cos(dtr(x))
}

package math_helpers

import "math"

func Dsin(x float64) float64 {
	return math.Sin(Dtr(x))
}

func Dcos(x float64) float64 {
	return math.Cos(Dtr(x))
}

// DTR - Degrees to radians
func Dtr(d float64) float64 {
	return (d * math.Pi) / 180.0
}

func Constrain(d float64) float64 {
	t := math.Mod(d, 360)
	if t < 0 {
		t += 360
	}
	return t
}

func GetSignPrefix(sign int) string {
	if sign >= 0 {
		return "+"
	}
	return "-"
}

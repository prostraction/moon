package moon

import "math"

func dsin(x float64) float64 {
	return math.Sin(dtr(x))
}

func dcos(x float64) float64 {
	return math.Cos(dtr(x))
}

// DTR - Degrees to radians
func dtr(d float64) float64 {
	return (d * math.Pi) / 180.0
}

func constrain(d float64) float64 {
	t := math.Mod(d, 360)
	if t < 0 {
		t += 360
	}
	return t
}

func getSignPrefix(sign int) string {
	if sign >= 0 {
		return "+"
	}
	return "-"
}

package moon

import "time"

var monthsGo = []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December}

func JulianDate(t time.Time) float64 {
	m := 1
	for i := range monthsGo {
		if t.Month() == monthsGo[i] {
			m = i + 1
		}
	}
	if m < 3 {
		t = t.AddDate(-1, 0, 0)
		m += 12
	}

	A := float64(t.Year()) / 100
	B := A / 4
	C := 2 - A + B
	E := 365.25 * float64(t.Year()+4716)
	F := 30.6001 * (float64(m) + 1)
	return C + float64(float64(t.Day()+t.Hour()/24.+t.Minute()/60)) + E + F - 1524.5
}

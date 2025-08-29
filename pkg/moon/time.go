package moon

import (
	"math"
	"time"
)

func ToJulianDate(t time.Time) float64 {
	m := 1
	for i := range months {
		if t.Month() == months[i] {
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

func FromJulianDate(j float64) time.Time {
	datey, datem, dated := jyear(j)

	//j1 -= 5.0 / 24.0 // 5 timezones west of UTC
	j += (30.0 / (24 * 60 * 60))

	timeh, timem, times := jhms(j)

	t := time.Date(datey, getMonth(datem), dated, timeh, timem, times, 0, time.UTC)
	return t
}

// JYMD - Convert Julian time to year, months, and days
func jyear(td float64) (int, int, int) {
	td += 0.5
	z := math.Floor(td)
	f := td - z

	var a float64
	if z < 2299161.0 {
		a = z
	} else {
		alpha := math.Floor((z - 1867216.25) / 36524.25)
		a = z + 1 + alpha - math.Floor(alpha/4)
	}

	b := a + 1524
	c := math.Floor((b - 122.1) / 365.25)
	d := math.Floor(365.25 * c)
	e := math.Floor((b - d) / 30.6001)

	mm := int(math.Floor(e))
	if mm >= 14 {
		mm -= 13
	} else {
		mm -= 1
	}

	year := int(math.Floor(c))
	if mm > 2 {
		year -= 4716
	} else {
		year -= 4715
	}

	day := int(math.Floor(b - d - math.Floor(30.6001*e) + f))

	return year, mm, day
}

// JHMS - Convert Julian time to hour, minutes, and seconds
func jhms(j float64) (int, int, int) {
	j += 0.5 // Astronomical to civil
	ij := (j - math.Floor(j)) * 86400.0
	hours := math.Floor(ij / 3600)
	minutes := math.Floor((ij / 60))
	seconds := math.Floor(ij)
	return int(hours), int(math.Mod(minutes, 60)), int(math.Mod(seconds, 60))
}

func getMonth(datem int) time.Month {
	datem = min(max(datem-1, 0), 11)
	return months[datem]
}

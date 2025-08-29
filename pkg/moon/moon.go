package moon

import (
	"math"
	"time"
)

func CalcMoonNumber(yearGiven int) int {
	d := 6 // 2000
	mod := 19
	yearCurrent := 2000

	if yearCurrent < yearGiven {
		for yearCurrent+mod < yearGiven {
			yearCurrent += mod
		}
		for yearCurrent < yearGiven {
			yearCurrent++
			d++
			if d > 19 {
				d = 1
			}
		}
	} else {
		for yearCurrent-mod > yearGiven {
			yearCurrent -= mod
		}
		for yearCurrent > yearGiven {
			yearCurrent--
			d--
			if d < 1 {
				d = 19
			}
		}
	}
	return d
}

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

type MoonTableElement struct {
	TNew  time.Time
	TFull time.Time
	T1    float64
	T2    float64
}

/*
JHMS  --  Convert Julian time to hour, minutes, and seconds,

	returned as three separate values.
*/
func jhms(j float64) (int, int, int) {
	j += 0.5 // Astronomical to civil
	ij := (j - math.Floor(j)) * 86400.0
	hours := math.Floor(ij / 3600)
	minutes := math.Floor((ij / 60))
	seconds := math.Floor(ij)
	return int(hours), int(math.Mod(minutes, 60)), int(math.Mod(seconds, 60))
}

func getIlluminatedFractionOfMoon(jd float64) float64 {
	const toRad = math.Pi / 180.0
	T := (jd - 2451545.) / 36525.0

	D := constrain(297.8501921+445267.1114034*T-0.0018819*T*T+1.0/545868.0*T*T*T-1.0/113065000.0*T*T*T*T) * toRad
	M := constrain(357.5291092+35999.0502909*T-0.0001536*T*T+1.0/24490000.0*T*T*T) * toRad
	Mp := constrain(134.9633964+477198.8675055*T+0.0087414*T*T+1.0/69699.0*T*T*T-1.0/14712000.0*T*T*T*T) * toRad

	i := constrain(180.-D*180./math.Pi-6.289*math.Sin(Mp)+2.1*math.Sin(M)-1.274*math.Sin(2.*D-Mp)-0.658*math.Sin(2.*D)-0.214*math.Sin(2.*Mp)-0.11*math.Sin(D)) * toRad

	k := (1. + math.Cos(i)) / 2.
	return k
}

func Gen(year int, month int, day int, hour int, minute int, second int, offset int) ([]*MoonTableElement, time.Duration, float64, string) {
	var moonDays time.Duration
	moonTable := []*MoonTableElement{}
	tGiven := time.Date(year, getMonth(month), day, hour-offset, minute, second, 0, time.UTC)

	var l int
	var k1, mtime float64
	var minx int
	var phaset []float64

	phaset = make([]float64, 0)

	// Tabulate new and full moons surrounding the year
	k1 = math.Floor((float64(year) - 1900) * 12.3685) // - 4
	minx = 0
	isNext := true
	for l = 0; ; l++ {
		mtime = truephase(k1, float64(l&1)*0.5)
		datey, _, _ := jyear(mtime)
		if datey >= year {
			minx++
		}
		phaseSign := mtime
		if (l & 1) == 0 {
			phaseSign = -mtime
		}
		phaset = append(phaset, phaseSign)
		if !isNext {
			break
		}
		if datey > year {
			minx++
			isNext = false
		}
		if (l & 1) != 0 {
			k1 += 1
		}
	}

	var lastnew float64
	for l = 0; l < minx; l++ {
		elem := &MoonTableElement{}

		mp := phaset[l]
		if mp < 0 {
			mp = -mp

			elem.T1 = mp
			elem.T2 = lastnew

			lastnew = mp
		}

		elem.T1 = mp
		elem.T2 = lastnew

		elem.TNew = cuzcoDateTime(lastnew)
		elem.TFull = cuzcoDateTime(mp)

		if elem.T1 != elem.T2 {
			moonTable = append(moonTable, elem)
			if tGiven.After(elem.TNew) && tGiven.Before(elem.TFull) {
				moonDays = tGiven.Sub(elem.TNew)
			}
		}
	}

	jdIllumination := getIlluminatedFractionOfMoon(JulianDate(tGiven))
	zodiacPosition := int((jdIllumination*360)/30) % 12

	return moonTable, moonDays, jdIllumination, getZodiacSign(zodiacPosition)
}

func getZodiacSign(position int) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer",
		"Leo", "Virgo", "Libra", "Scorpio",
		"Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	return signs[position]
}

func truephase(k, phase float64) float64 {
	var t, t2, t3, pt, m, mprime, f float64
	SynMonth := 29.53058868 // Synodic month (mean time from new to next new Moon)

	k += phase           // Add phase to new moon time
	t = k / 1236.85      // Time in Julian centuries from 1900 January 0.5
	t2 = t * t           // Square for frequent use
	t3 = t2 * t          // Cube for frequent use
	pt = 2415020.75933 + // Mean time of phase
		SynMonth*k +
		0.0001178*t2 -
		0.000000155*t3 +
		0.00033*dsin(166.56+132.87*t-0.009173*t2)

	m = 359.2242 + // Sun's mean anomaly
		29.10535608*k -
		0.0000333*t2 -
		0.00000347*t3
	mprime = 306.0253 + // Moon's mean anomaly
		385.81691806*k +
		0.0107306*t2 +
		0.00001236*t3
	f = 21.2964 + // Moon's argument of latitude
		390.67050646*k -
		0.0016528*t2 -
		0.00000239*t3

	if (phase < 0.01) || (math.Abs(phase-0.5) < 0.01) {
		// Corrections for New and Full Moon
		pt += (0.1734-0.000393*t)*dsin(m) +
			0.0021*dsin(2*m) -
			0.4068*dsin(mprime) +
			0.0161*dsin(2*mprime) -
			0.0004*dsin(3*mprime) +
			0.0104*dsin(2*f) -
			0.0051*dsin(m+mprime) -
			0.0074*dsin(m-mprime) +
			0.0004*dsin(2*f+m) -
			0.0004*dsin(2*f-m) -
			0.0006*dsin(2*f+mprime) +
			0.0010*dsin(2*f-mprime) +
			0.0005*dsin(m+2*mprime)
	} else if (math.Abs(phase-0.25) < 0.01) || (math.Abs(phase-0.75) < 0.01) {
		pt += (0.1721-0.0004*t)*dsin(m) +
			0.0021*dsin(2*m) -
			0.6280*dsin(mprime) +
			0.0089*dsin(2*mprime) -
			0.0004*dsin(3*mprime) +
			0.0079*dsin(2*f) -
			0.0119*dsin(m+mprime) -
			0.0047*dsin(m-mprime) +
			0.0003*dsin(2*f+m) -
			0.0004*dsin(2*f-m) -
			0.0006*dsin(2*f+mprime) +
			0.0021*dsin(2*f-mprime) +
			0.0003*dsin(m+2*mprime) +
			0.0004*dsin(m-2*mprime) -
			0.0003*dsin(2*m+mprime)

		if phase < 0.5 {
			// First quarter correction
			pt += 0.0028 - 0.0004*dcos(m) + 0.0003*dcos(mprime)
		} else {
			// Last quarter correction
			pt += -0.0028 + 0.0004*dcos(m) - 0.0003*dcos(mprime)
		}
	}
	return pt
}

func getMonth(datem int) time.Month {
	datem = datem - 1
	if datem < 0 {
		datem = 0
	}
	if datem > 11 {
		datem = 11
	}

	return monthsGo[datem]
}

func cuzcoDateTime(j float64) time.Time {
	datey, datem, dated := jyear(j)
	//t.AddDate(datey, datem, dated)

	j1 := j
	//j1 -= 5.0 / 24.0 // 5 timezones west of UTC
	j1 += (30.0 / (24 * 60 * 60))

	timeh, timem, times := jhms(j1)

	t := time.Date(datey, getMonth(datem), dated, timeh, timem, times, 0, time.UTC)
	return t
}

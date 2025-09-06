package moon

import (
	"math"
	"time"
)

func (c *Cache) CreateMoonTable(timeGiven time.Time) []*MoonTableElement {
	t := time.Date(timeGiven.Year(), 0, 0, 0, 0, 0, 0, timeGiven.Location())
	if c.tables != nil && c.tables[t.String()] != nil {
		return c.tables[t.String()]
	}
	moonTable := []*MoonTableElement{}

	var l int
	var k1, mtime float64
	var minx int
	var phaset []float64

	phaset = make([]float64, 0)

	// Tabulate new and full moons surrounding the year
	k1 = math.Floor((float64(timeGiven.Year()) - 1900) * 12.3685)
	minx = 0
	isNext := true
	for l = 0; ; l++ {
		mtime = truephase(k1, float64(l&1)*0.5)
		datey, _, _ := jyear(mtime)
		if datey >= timeGiven.Year() {
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
		if datey > timeGiven.Year() {
			minx++
			isNext = false
		}
		if (l & 1) != 0 {
			k1 += 1
		}
	}

	var lastnew float64
	for l = range minx {
		elem := &MoonTableElement{}
		mp := phaset[l]
		if mp < 0 {
			mp = -mp

			elem.t1 = mp
			elem.t2 = lastnew

			lastnew = mp
		}

		elem.t1 = mp
		elem.t2 = lastnew

		firstQuarterTime := lastnew
		firstQuarterIllum := 0.
		for firstQuarterIllum < 0.5 {
			firstQuarterTime += 0.001
			firstQuarterIllum = GetCurrentMoonIllumination(FromJulianDate(firstQuarterTime, timeGiven.Location()), timeGiven.Location())
		}
		elem.FirstQuarter = FromJulianDate(firstQuarterTime, timeGiven.Location())

		lastQuarterTime := mp
		lastQuarterIllum := 1.
		for lastQuarterIllum > 0.5000 {
			lastQuarterTime += 0.001
			lastQuarterIllum = GetCurrentMoonIllumination(FromJulianDate(lastQuarterTime, timeGiven.Location()), timeGiven.Location())
		}
		elem.LastQuarter = FromJulianDate(lastQuarterTime, timeGiven.Location())

		elem.NewMoon = FromJulianDate(lastnew, timeGiven.Location())
		elem.FullMoon = FromJulianDate(mp, timeGiven.Location())

		if elem.t1 != elem.t2 {
			moonTable = append(moonTable, elem)
		}

		if elem.LastQuarter.Year() > timeGiven.Year() {
			break
		}
	}
	if c.tables == nil {
		c.tables = make(map[string][]*MoonTableElement)
	}
	if c.tables[t.String()] == nil {
		c.tables[t.String()] = moonTable
	}
	return moonTable
}

func (c *Cache) BeginMoonDayToEarthDay(tGiven time.Time, duration time.Duration) time.Time {
	moonTable := c.CreateMoonTable(tGiven)
	for i := range moonTable {
		elem := moonTable[i]
		if elem.t1 != elem.t2 {
			if tGiven.After(elem.NewMoon) && tGiven.Before(elem.LastQuarter) {
				t := elem.NewMoon
				t = t.Add(duration)
				return t
			}
			if i < len(moonTable) {
				elem2 := moonTable[i+1]
				if tGiven.After(elem.LastQuarter) && tGiven.Before(elem2.NewMoon) {
					t := elem.NewMoon
					t = t.Add(duration)
					return t
				}
			}
		}
	}
	return time.Time{} // fix
}

func GetMoonDays(tGiven time.Time, table []*MoonTableElement) time.Duration {
	var moonDays time.Duration
	for i := range table {
		elem := table[i]

		if elem.t1 != elem.t2 {
			if tGiven.After(elem.NewMoon) /*&& tGiven.Before(elem.TFull)*/ {
				moonDays = tGiven.Sub(elem.NewMoon)
			}
		}
	}
	return moonDays
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

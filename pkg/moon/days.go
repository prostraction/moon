package moon

import (
	"math"
	"time"

	il "moon/pkg/illumination"
	jt "moon/pkg/julian-time"
	phase "moon/pkg/phase"
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
		mtime = phase.Truephase(k1, float64(l&1)*0.5)
		datey, _, _ := jt.Jyear(mtime)
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

		firstQuarterTime := il.BinarySearchIllumination(lastnew, mp, timeGiven.Location(), true)
		elem.FirstQuarter = jt.FromJulianDate(firstQuarterTime, timeGiven.Location())

		lastQuarterTime := il.BinarySearchIllumination(mp, mp+10, timeGiven.Location(), false)
		elem.LastQuarter = jt.FromJulianDate(lastQuarterTime, timeGiven.Location())

		elem.NewMoon = jt.FromJulianDate(lastnew, timeGiven.Location())
		elem.FullMoon = jt.FromJulianDate(mp, timeGiven.Location())

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

func BeginMoonDayToEarthDay(tGiven time.Time, duration time.Duration, moonTable []*MoonTableElement) time.Time {
	if len(moonTable) == 0 {
		return time.Time{}
	}
	for i := range moonTable {
		elem := moonTable[i]
		if elem.t1 != elem.t2 {
			if tGiven.After(elem.NewMoon) && tGiven.Before(elem.LastQuarter) {
				t := elem.NewMoon
				t = t.Add(duration)
				return t
			}
			if i < len(moonTable)-1 {
				elem2 := moonTable[i+1]
				if tGiven.After(elem.LastQuarter) && tGiven.Before(elem2.NewMoon) {
					t := elem.NewMoon
					t = t.Add(duration)
					return t
				}
			}
		}
	}
	return time.Time{}
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

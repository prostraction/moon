package moon

import (
	pos "moon/pkg/position"
	"time"
)

func (c *Cache) CurrentMoonDays(tGiven time.Time, loc *time.Location) (time.Duration, time.Duration, time.Duration) {
	if loc == nil {
		loc = time.UTC
	}

	dayBeginTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayEndTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonTable := c.CreateMoonTable(tGiven)
	beginMoonDays := GetMoonDays(dayBeginTime, moonTable)
	currentMoonDays := GetMoonDays(tGiven, moonTable)
	endMoonDays := GetMoonDays(dayEndTime, moonTable)

	return beginMoonDays, currentMoonDays, endMoonDays
}

func (c *Cache) MoonDetailed(tGiven time.Time, loc *time.Location, lang string, longitude float64, latitude float64) *MoonDaysDetailed {
	if loc == nil {
		loc = time.UTC
	}

	dayYesterday := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()-1, 0, 0, 0, 0, loc)
	dayToday := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayTomorrow := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonDaysDetailed := new(MoonDaysDetailed)
	moonDaysDetailed.Day = make([]MoonDay, 2)
	moonDaysDetailed.Count = 2

	moonRiseYesterday, err1 := pos.GetRisesDay(dayYesterday.Year(), int(dayYesterday.Month()), dayYesterday.Day(), loc, 2, longitude, latitude)
	moonRiseToday, err2 := pos.GetRisesDay(dayToday.Year(), int(dayToday.Month()), dayToday.Day(), loc, 2, longitude, latitude)
	moonRiseTomorrow, err3 := pos.GetRisesDay(dayTomorrow.Year(), int(dayTomorrow.Month()), dayTomorrow.Day(), loc, 2, longitude, latitude)

	if err1 == nil && err2 == nil {
		if moonRiseYesterday.IsMoonRise {
			moonDaysDetailed.Day[0].Begin = new(time.Time)
			*moonDaysDetailed.Day[0].Begin = moonRiseYesterday.Moonrise.TimeISO
			moonDaysDetailed.Day[0].IsBeginExists = true
		}
		if moonRiseToday.IsMoonRise {
			moonDaysDetailed.Day[0].End = new(time.Time)
			*moonDaysDetailed.Day[0].End = moonRiseToday.Moonrise.TimeISO
			moonDaysDetailed.Day[0].IsEndExists = true
		}
	}
	if err2 == nil && err3 == nil {
		if moonRiseToday.IsMoonRise {
			moonDaysDetailed.Day[1].Begin = new(time.Time)
			*moonDaysDetailed.Day[1].Begin = moonRiseToday.Moonrise.TimeISO
			moonDaysDetailed.Day[1].IsBeginExists = true
		}
		if moonRiseTomorrow.IsMoonRise {
			moonDaysDetailed.Day[1].End = new(time.Time)
			*moonDaysDetailed.Day[1].End = moonRiseTomorrow.Moonrise.TimeISO
			moonDaysDetailed.Day[1].IsEndExists = true
		}
	}

	if !(moonDaysDetailed.Day[0].IsBeginExists && moonDaysDetailed.Day[0].IsEndExists) {
		moonDaysDetailed.Count = 1
	}

	return moonDaysDetailed
}

func (c *Cache) GenerateMoonTable(tGiven time.Time) []*MoonTableElement {
	return c.CreateMoonTable(tGiven)
}

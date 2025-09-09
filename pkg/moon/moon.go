package moon

import (
	jt "moon/pkg/julian-time"
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

func (c *Cache) MoonDetailed(tGiven time.Time, loc *time.Location, lang string) *MoonDaysDetailed {
	if loc == nil {
		loc = time.UTC
	}

	dayBeginTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayEndTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonTable := c.CreateMoonTable(tGiven)
	beginMoonDays := GetMoonDays(dayBeginTime, moonTable).Minutes() / jt.Fminute
	endMoonDays := GetMoonDays(dayEndTime, moonTable).Minutes() / jt.Fminute

	moonDaysDetailed := new(MoonDaysDetailed)
	tFirstDayBegin := BeginMoonDayToEarthDay(tGiven, time.Duration(int(time.Minute)+int(time.Hour)*24*int(beginMoonDays)), moonTable)
	tFirstDayEnd := BeginMoonDayToEarthDay(tGiven, time.Duration(int(time.Minute)+int(time.Hour)*24*int(beginMoonDays+1)), moonTable)

	if int(endMoonDays) == 0 {
		endMoonDays += (beginMoonDays + 1)
	}

	tSecondDayBegin := BeginMoonDayToEarthDay(tGiven, time.Duration(int(time.Minute)+int(time.Hour)*24*int(endMoonDays)), moonTable)
	tSecondDayEnd := BeginMoonDayToEarthDay(tGiven, time.Duration(int(time.Minute)+int(time.Hour)*24*int(endMoonDays+1)), moonTable)

	// regular case or 29 day -> 0 day
	if int(endMoonDays)-int(beginMoonDays) <= 1 || (int(beginMoonDays) != 0 && int(endMoonDays) == 0) {
		moonDaysDetailed.Day = make([]MoonDay, 2)
		moonDaysDetailed.Count = 2

		moonDaysDetailed.Day[1].Begin = tSecondDayBegin
		moonDaysDetailed.Day[1].End = tSecondDayEnd
	} else {
		// 3 days
		moonDaysDetailed.Day = make([]MoonDay, 3)
		moonDaysDetailed.Count = 3

		moonDaysDetailed.Day[2].Begin = tSecondDayBegin
		moonDaysDetailed.Day[2].End = tSecondDayEnd

		moonDaysDetailed.Day[1].Begin = tFirstDayEnd
		moonDaysDetailed.Day[1].End = tSecondDayBegin

	}

	moonDaysDetailed.Day[0].Begin = tFirstDayBegin
	moonDaysDetailed.Day[0].End = tFirstDayEnd

	return moonDaysDetailed
}

func (c *Cache) GenerateMoonTable(tGiven time.Time) []*MoonTableElement {
	return c.CreateMoonTable(tGiven)
}

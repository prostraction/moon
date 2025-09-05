package moon

import (
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

type illumFunc func(tGiven time.Time, loc *time.Location) float64

func CurrentMoonPhase(tGiven time.Time, loc *time.Location, lang string) (float64, float64, float64, PhaseResp, PhaseResp, PhaseResp) {
	currentMoonIllumination, currentMoonIlluminationBefore, currentMoonIlluminationAfter := currentMoonPhaseCalc(tGiven, loc, GetCurrentMoonIllumination)
	dayBeginMoonIllumination, dayBeginMoonIlluminationBefore, dayBeginMoonIlluminationAfter := currentMoonPhaseCalc(tGiven, loc, GetDailyMoonIllumination)
	dayEndMoonIllumination, dayEndMoonIlluminationBefore, dayEndMoonIlluminationAfter := currentMoonPhaseCalc(tGiven.AddDate(0, 0, 1), loc, GetDailyMoonIllumination)

	moonPhaseCurrent := PhaseResp{}
	moonPhaseBegin := PhaseResp{}
	moonPhaseEnd := PhaseResp{}

	moonPhaseCurrent.Name, moonPhaseCurrent.Emoji = GetMoonPhase(currentMoonIlluminationBefore, currentMoonIllumination, currentMoonIlluminationAfter, lang)
	moonPhaseBegin.Name, moonPhaseBegin.Emoji = GetMoonPhase(dayBeginMoonIlluminationBefore, dayBeginMoonIllumination, dayBeginMoonIlluminationAfter, lang)
	moonPhaseEnd.Name, moonPhaseEnd.Emoji = GetMoonPhase(dayEndMoonIlluminationBefore, dayEndMoonIllumination, dayEndMoonIlluminationAfter, lang)

	return currentMoonIllumination, dayBeginMoonIllumination, dayEndMoonIllumination, moonPhaseCurrent, moonPhaseBegin, moonPhaseEnd
}

func currentMoonPhaseCalc(tGiven time.Time, loc *time.Location, calcF illumFunc) (float64, float64, float64) {
	moonIllumination := calcF(tGiven, loc)
	moonIlluminationBefore := calcF(tGiven.AddDate(0, 0, -1), loc)
	moonIlluminationAfter := calcF(tGiven.AddDate(0, 0, 1), loc)

	// in rare UTC-12 case they are equal
	if moonIllumination == moonIlluminationBefore {
		moonIlluminationBefore = calcF(tGiven.AddDate(0, 0, -2), loc)
	}
	// just in case
	if moonIllumination == moonIlluminationAfter {
		moonIlluminationAfter = calcF(tGiven.AddDate(0, 0, 2), loc)
	}

	return moonIllumination, moonIlluminationBefore, moonIlluminationAfter
}

func (c *Cache) GenerateMoonTable(tGiven time.Time) []*MoonTableElement {
	/*moonTable := CreateMoonTable(tGiven)
	moonDays := GetMoonDays(tGiven, moonTable)
	moonIllumination := GetDailyMoonIllumination(tGiven, nil)

	zodiacPosition := int((moonDays.Minutes()/Fminute*360)/30) / 30 % 12*/
	return c.CreateMoonTable(tGiven) //, moonDays, moonIllumination, getZodiacSign(zodiacPosition)
}

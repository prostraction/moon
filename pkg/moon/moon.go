package moon

import (
	"log"
	"time"
)

func (c *Cache) CurrentMoonDays(tGiven time.Time, loc *time.Location) (time.Duration, time.Duration, time.Duration, string) {
	if loc == nil {
		loc = time.UTC
	}

	dayBeginTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day(), 0, 0, 0, 0, loc)
	dayEndTime := time.Date(tGiven.Year(), tGiven.Month(), tGiven.Day()+1, 0, 0, 0, 0, loc)

	moonTable := c.CreateMoonTable(tGiven)
	beginMoonDays := GetMoonDays(dayBeginTime, moonTable)
	currentMoonDays := GetMoonDays(tGiven, moonTable)
	endMoonDays := GetMoonDays(dayEndTime, moonTable)

	zodiacPosition := int((currentMoonDays.Minutes()/Fminute*360.)/30.) / 30. % 12

	return beginMoonDays, currentMoonDays, endMoonDays, getZodiacSign(zodiacPosition)
}

type illumFunc func(tGiven time.Time, loc *time.Location) float64

func CurrentMoonPhase(tGiven time.Time, loc *time.Location) (float64, float64, float64, string, string, string, string, string, string) {
	currentMoonIllumination, currentMoonIlluminationBefore, currentMoonIlluminationAfter := currentMoonPhaseCalc(tGiven, loc, GetCurrentMoonIllumination)
	dayBeginMoonIllumination, dayBeginMoonIlluminationBefore, dayBeginMoonIlluminationAfter := currentMoonPhaseCalc(tGiven, loc, GetDailyMoonIllumination)
	dayEndMoonIllumination, dayEndMoonIlluminationBefore, dayEndMoonIlluminationAfter := currentMoonPhaseCalc(tGiven.Local().AddDate(0, 0, 1), loc, GetDailyMoonIllumination)

	moonPhaseCurrent, moonPhaseEmojiCurrent := GetMoonPhase(currentMoonIlluminationBefore, currentMoonIllumination, currentMoonIlluminationAfter)
	moonPhaseDayBegin, moonPhaseEmojiDayBegin := GetMoonPhase(dayBeginMoonIlluminationBefore, dayBeginMoonIllumination, dayBeginMoonIlluminationAfter)
	moonPhaseDayEnd, moonPhaseEmojiDayEnd := GetMoonPhase(dayEndMoonIlluminationBefore, dayEndMoonIllumination, dayEndMoonIlluminationAfter)

	return currentMoonIllumination, dayBeginMoonIllumination, dayEndMoonIllumination, moonPhaseCurrent, moonPhaseEmojiCurrent, moonPhaseDayBegin, moonPhaseEmojiDayBegin, moonPhaseDayEnd, moonPhaseEmojiDayEnd
}

func currentMoonPhaseCalc(tGiven time.Time, loc *time.Location, calcF illumFunc) (float64, float64, float64) {

	moonIllumination := calcF(tGiven, loc)
	log.Println(tGiven, moonIllumination)

	moonIlluminationBefore := calcF(tGiven.Local().AddDate(0, 0, -1), loc)
	log.Println(tGiven.Local().AddDate(0, 0, -1), moonIlluminationBefore)

	moonIlluminationAfter := calcF(tGiven.Local().AddDate(0, 0, 1), loc)
	log.Println(tGiven.Local().AddDate(0, 0, 1), moonIlluminationAfter)

	// in rare UTC-12 case they are equal
	if moonIllumination == moonIlluminationBefore {
		moonIlluminationBefore = calcF(tGiven.Local().AddDate(0, 0, -1), loc)
	}
	// just in case
	if moonIllumination == moonIlluminationAfter {
		moonIlluminationAfter = calcF(tGiven.Local().AddDate(0, 0, 1), loc)
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

func getZodiacSign(position int) string {
	if position >= 0 && position < len(signs) {
		return signs[position]
	}
	return ""
}

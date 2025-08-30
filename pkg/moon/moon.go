package moon

import (
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

func CurrentMoonDays(tGiven time.Time) (time.Duration, string) {
	moonTable := CreateMoonTable(tGiven)
	moonDays := GetMoonDays(tGiven, moonTable)
	return moonDays, "not working"
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
	moonIlluminationBefore := calcF(tGiven.Local().AddDate(0, 0, -1), loc)
	moonIlluminationAfter := calcF(tGiven.Local().AddDate(0, 0, 1), loc)

	// in rare UTC-12 case they are equal
	if moonIllumination == moonIlluminationBefore {
		moonIlluminationBefore = calcF(tGiven.Local().AddDate(0, 0, -2), loc)
	}
	// just in case
	if moonIllumination == moonIlluminationAfter {
		moonIlluminationAfter = calcF(tGiven.Local().AddDate(0, 0, 2), loc)
	}

	return moonIllumination, moonIlluminationBefore, moonIlluminationAfter
}

func Gen(tGiven time.Time) ([]*MoonTableElement, time.Duration, float64, string) {
	moonTable := CreateMoonTable(tGiven)
	moonDays := GetMoonDays(tGiven, moonTable)
	moonIllumination := GetDailyMoonIllumination(tGiven, nil)

	return moonTable, moonDays, moonIllumination, "not working" //getZodiacSign(zodiacPosition)
}

func getZodiacSign(position int) string {
	signs := []string{
		"Aries", "Taurus", "Gemini", "Cancer",
		"Leo", "Virgo", "Libra", "Scorpio",
		"Sagittarius", "Capricorn", "Aquarius", "Pisces",
	}
	return signs[position]
}

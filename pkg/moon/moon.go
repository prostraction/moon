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

func CurrentMoonPhase(tGiven time.Time) (float64, string, string) {
	moonIllumination := GetMoonIllumination(tGiven)
	moonIlluminationBefore := GetMoonIllumination(tGiven.Local().AddDate(0, 0, -1))
	moonIlluminationAfter := GetMoonIllumination(tGiven.Local().AddDate(0, 0, 1))

	moonPhase, moonPhaseEmoji := GetMoonPhase(moonIlluminationBefore, moonIllumination, moonIlluminationAfter)
	return moonIllumination, moonPhase, moonPhaseEmoji

}

func Gen(tGiven time.Time) ([]*MoonTableElement, time.Duration, float64, string) {

	moonTable := CreateMoonTable(tGiven)
	moonDays := GetMoonDays(tGiven, moonTable)
	moonIllumination := GetMoonIllumination(tGiven)

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
